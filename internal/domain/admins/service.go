package admins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lucasgarciaf/df-backend-go/config"
	"github.com/lucasgarciaf/df-backend-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
)

type AuthClaims struct {
	AdminID primitive.ObjectID `json:"admin_id"`
	Role    string             `json:"role"`
	jwt.RegisteredClaims
}

type AdminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) CreateAdmin(admin Admin, password string) (primitive.ObjectID, error) {
	// Create user in Keycloak
	err := s.createUserInKeycloak(admin, password)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Create user in MongoDB
	existingAdmin, _ := s.repo.GetAdminByEmail(admin.Email)
	if existingAdmin != nil {
		return primitive.NilObjectID, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, err
	}
	admin.PasswordHash = string(hashedPassword)
	admin.Role = "admin" // Ensure Role is set
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	return s.repo.CreateAdmin(admin)
}

func (s *AdminService) Authenticate(email, password string) (string, error) {
	admin, err := s.repo.GetAdminByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.createToken(admin.ID, admin.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AdminService) createToken(adminID primitive.ObjectID, role string) (string, error) {
	claims := &AuthClaims{
		AdminID: adminID,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

func (s *AdminService) GetAdminByID(id primitive.ObjectID) (*Admin, error) {
	return s.repo.GetAdminByID(id)
}

func (s *AdminService) GetAdminByEmail(email string) (*Admin, error) {
	return s.repo.GetAdminByEmail(email)
}

func (s *AdminService) UpdateAdmin(admin Admin) error {
	admin.UpdatedAt = time.Now()
	return s.repo.UpdateAdmin(admin)
}

func (s *AdminService) DeleteAdmin(id primitive.ObjectID) error {
	return s.repo.DeleteAdmin(id)
}

func (s *AdminService) Logout(refreshToken string) error {
	logoutURL := fmt.Sprintf("%sprotocol/openid-connect/logout", config.KEYCLOAK_URL)

	data := map[string]string{
		"client_id":     config.KEYCLOAK_CLIENT_ID,
		"client_secret": config.KEYCLOAK_CLIENT_SECRET,
		"refresh_token": refreshToken,
	}

	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}

	req, err := http.NewRequest("POST", logoutURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to logout from Keycloak: %s", resp.Status)
	}

	return nil
}

func (s *AdminService) createUserInKeycloak(admin Admin, password string) error {
	keycloakUser := struct {
		Username    string `json:"username"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		Email       string `json:"email"`
		Enabled     bool   `json:"enabled"`
		Credentials []struct {
			Type      string `json:"type"`
			Value     string `json:"value"`
			Temporary bool   `json:"temporary"`
		} `json:"credentials"`
		RealmRoles []string `json:"realmRoles"`
	}{
		Username:  admin.Username,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Email:     admin.Email,
		Enabled:   true,
		Credentials: []struct {
			Type      string `json:"type"`
			Value     string `json:"value"`
			Temporary bool   `json:"temporary"`
		}{
			{
				Type:      "password",
				Value:     password,
				Temporary: false,
			},
		},
		RealmRoles: []string{"admin"}, // Assign "admin" role
	}

	userJSON, err := json.Marshal(keycloakUser)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%sadmin/realms/%s/users", config.KEYCLOAK_URL, config.KEYCLOAK_REALM)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		return err
	}

	token, err := utils.GetClientToken()
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Use a valid client access token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create user in Keycloak: %s", resp.Status)
	}

	return nil
}
