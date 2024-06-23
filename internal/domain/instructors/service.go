package instructors

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
	InstructorID primitive.ObjectID `json:"instructor_id"`
	Role         string             `json:"role"`
	jwt.RegisteredClaims
}

type InstructorService struct {
	repo InstructorRepository
}

func NewInstructorService(repo InstructorRepository) *InstructorService {
	return &InstructorService{repo: repo}
}

func (s *InstructorService) CreateInstructor(instructor Instructor, password string) (primitive.ObjectID, error) {
	// Create user in Keycloak
	err := s.createUserInKeycloak(instructor, password)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Create user in MongoDB
	existingInstructor, _ := s.repo.GetInstructorByEmail(instructor.Email)
	if existingInstructor != nil {
		return primitive.NilObjectID, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, err
	}
	instructor.PasswordHash = string(hashedPassword)
	instructor.Role = "instructor" // Ensure Role is set
	instructor.CreatedAt = time.Now()
	instructor.UpdatedAt = time.Now()

	return s.repo.CreateInstructor(instructor)
}

func (s *InstructorService) Authenticate(email, password string) (string, error) {
	instructor, err := s.repo.GetInstructorByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(instructor.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.createToken(instructor.ID, instructor.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *InstructorService) createToken(instructorID primitive.ObjectID, role string) (string, error) {
	claims := &AuthClaims{
		InstructorID: instructorID,
		Role:         role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

func (s *InstructorService) GetInstructorByID(id primitive.ObjectID) (*Instructor, error) {
	return s.repo.GetInstructorByID(id)
}

func (s *InstructorService) UpdateInstructor(instructor Instructor) error {
	instructor.UpdatedAt = time.Now()
	return s.repo.UpdateInstructor(instructor)
}

func (s *InstructorService) DeleteInstructor(id primitive.ObjectID) error {
	return s.repo.DeleteInstructor(id)
}

func (s *InstructorService) Logout(refreshToken string) error {
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

func (s *InstructorService) createUserInKeycloak(instructor Instructor, password string) error {
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
		Username:  instructor.Username,
		FirstName: instructor.FirstName,
		LastName:  instructor.LastName,
		Email:     instructor.Email,
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
		RealmRoles: []string{"instructor"}, // Assign "instructor" role
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
