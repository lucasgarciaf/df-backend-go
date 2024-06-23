package students

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	StudentID primitive.ObjectID `json:"student_id"`
	Role      string             `json:"role"`
	jwt.RegisteredClaims
}

type StudentService struct {
	repo StudentRepository
}

func NewStudentService(repo StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) CreateStudent(student Student, password string) (primitive.ObjectID, error) {
	// Create user in Keycloak
	err := s.createUserInKeycloak(student, password)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Create user in MongoDB
	existingStudent, _ := s.repo.GetStudentByEmail(student.Email)
	if existingStudent != nil {
		return primitive.NilObjectID, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, err
	}
	student.PasswordHash = string(hashedPassword)
	student.Role = "student" // Ensure Role is set
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()

	return s.repo.CreateStudent(student)
}

func (s *StudentService) Authenticate(email, password string) (string, error) {
	student, err := s.repo.GetStudentByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.createToken(student.ID, student.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *StudentService) createToken(studentID primitive.ObjectID, role string) (string, error) {
	claims := &AuthClaims{
		StudentID: studentID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

func (s *StudentService) GetStudentByID(id primitive.ObjectID) (*Student, error) {
	return s.repo.GetStudentByID(id)
}

func (s *StudentService) UpdateStudent(student Student) error {
	student.UpdatedAt = time.Now()
	return s.repo.UpdateStudent(student)
}

func (s *StudentService) DeleteStudent(id primitive.ObjectID) error {
	return s.repo.DeleteStudent(id)
}

func (s *StudentService) createUserInKeycloak(student Student, password string) error {
	// Create user in Keycloak
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
		Username:  student.Username,
		FirstName: student.FirstName,
		LastName:  student.LastName,
		Email:     student.Email,
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
		RealmRoles: []string{"student"},
	}

	userJSON, err := json.Marshal(keycloakUser)
	if err != nil {
		log.Printf("Failed to marshal user JSON: %v", err)
		return err
	}

	url := fmt.Sprintf("%sadmin/realms/%s/users", config.KEYCLOAK_URL, config.KEYCLOAK_REALM)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	token, err := utils.GetClientToken()
	if err != nil {
		log.Printf("Failed to get client token: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to do request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to create user in Keycloak: %s", resp.Status)
		return fmt.Errorf("failed to create user in Keycloak: %s", resp.Status)
	}

	// Get the user ID from Keycloak
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response from Keycloak: %v", err)
		return fmt.Errorf("failed to decode response from Keycloak: %w", err)
	}

	// Assign the "student" role to the user
	if err := s.assignRoleToUser(result.ID, "student"); err != nil {
		log.Printf("Failed to assign role to user: %v", err)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (s *StudentService) assignRoleToUser(userID, roleName string) error {
	// Get the role ID
	url := fmt.Sprintf("%sadmin/realms/%s/roles/%s", config.KEYCLOAK_URL, config.KEYCLOAK_REALM, roleName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for role ID: %v", err)
		return err
	}

	token, err := utils.GetClientToken()
	if err != nil {
		log.Printf("Failed to get client token: %v", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute request for role ID: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get role ID from Keycloak: %s", resp.Status)
		return fmt.Errorf("failed to get role ID from Keycloak: %s", resp.Status)
	}

	var role struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		log.Printf("Failed to decode role response: %v", err)
		return fmt.Errorf("failed to decode role response: %w", err)
	}
	log.Printf("Role ID for %s: %s", roleName, role.ID)

	// Assign the role to the user
	assignRoleURL := fmt.Sprintf("%sadmin/realms/%s/users/%s/role-mappings/realm", config.KEYCLOAK_URL, config.KEYCLOAK_REALM, userID)
	roleMapping := []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		{ID: role.ID, Name: role.Name},
	}

	roleMappingJSON, err := json.Marshal(roleMapping)
	if err != nil {
		log.Printf("Failed to marshal role mapping JSON: %v", err)
		return err
	}

	req, err = http.NewRequest("POST", assignRoleURL, bytes.NewBuffer(roleMappingJSON))
	if err != nil {
		log.Printf("Failed to create request for assigning role: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Failed to execute request for assigning role: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("Failed to assign role to user in Keycloak: %s", resp.Status)
		return fmt.Errorf("failed to assign role to user in Keycloak: %s", resp.Status)
	}

	log.Printf("Successfully assigned role %s to user %s", roleName, userID)
	return nil
}
