package students

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Nerzal/gocloak"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lucasgarciaf/df-backend-go/config"
	"github.com/lucasgarciaf/df-backend-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("error in sercvice.go = invalid email or password")
	ErrEmailExists        = errors.New("error in sercvice.go = email already exists")
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
	// Check if the email already exists
	existingStudent, _ := s.repo.GetStudentByEmail(student.Email)
	if existingStudent != nil {
		log.Printf("Email already exists: %s", student.Email)
		return primitive.NilObjectID, ErrEmailExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return primitive.NilObjectID, err
	}
	student.PasswordHash = string(hashedPassword)
	student.Role = "student"
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()

	// Create the student in MongoDB
	studentID, err := s.repo.CreateStudent(student)
	if err != nil {
		log.Printf("Error creating student in MongoDB: %v", err)
		return primitive.NilObjectID, err
	}

	// Create the student in Keycloak
	err = s.createUserInKeycloak(student, password)
	if err != nil {
		// Rollback MongoDB creation in case of Keycloak error
		s.repo.DeleteStudent(studentID)
		log.Printf("Error creating student in Keycloak: %v", err)
		return primitive.NilObjectID, err
	}

	return studentID, nil
}

func (s *StudentService) Authenticate(email, password string) (string, error) {
	client := gocloak.NewClient(config.KeycloakURL)
	token, err := client.Login(config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm, email, password)
	if err != nil {
		log.Printf("Login to Keycloak failed: %v", err)
		return "", ErrInvalidCredentials
	}

	return token.AccessToken, nil
}

func (s *StudentService) AuthenticateWithKeycloak(email, password string) (string, error) {
	// Login to Keycloak
	token, err := utils.GetKeycloakToken(email, password)
	if err != nil {
		log.Printf("Login to Keycloak failed: %v", err)
		return "", err
	}

	return token, nil
}

func (s *StudentService) GetStudentByID(id primitive.ObjectID) (*Student, error) {
	return s.repo.GetStudentByID(id)
}

func (s *StudentService) GetAllStudents() ([]Student, error) {
	return s.repo.GetAllStudents()
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

	url := fmt.Sprintf("%s/admin/realms/%s/users", config.KeycloakURL, config.KeycloakRealm)
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

	// Keycloak does not return a body on successful user creation, so we use the Location header to get the user ID.
	location := resp.Header.Get("Location")
	if location == "" {
		log.Printf("Failed to get user location header from Keycloak response")
		return fmt.Errorf("failed to get user location header from Keycloak response")
	}

	// Extract user ID from the Location header
	segments := strings.Split(location, "/")
	userID := segments[len(segments)-1]
	log.Printf("User ID from Keycloak: %s", userID)

	// Assign the "student" role to the user
	if err := s.assignRoleToUser(userID, "student"); err != nil {
		log.Printf("Failed to assign role to user: %v", err)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (s *StudentService) assignRoleToUser(userID, roleName string) error {
	url := fmt.Sprintf("%s/admin/realms/%s/roles/%s", config.KeycloakURL, config.KeycloakRealm, roleName)
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

	assignRoleURL := fmt.Sprintf("%s/admin/realms/%s/users/%s/role-mappings/realm", config.KeycloakURL, config.KeycloakRealm, userID)
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
