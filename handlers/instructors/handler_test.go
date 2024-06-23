package instructors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/instructors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateInstructor(t *testing.T) {
	r := setupRouter()
	instructorService := instructors.NewInstructorService(&MockInstructorRepository{})
	instructorHandler := NewInstructorHandler(instructorService)

	r.POST("/instructors", instructorHandler.CreateInstructor)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/instructors", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %v but got %v", http.StatusCreated, w.Code)
	}
}

type MockInstructorRepository struct{}

func (m *MockInstructorRepository) CreateInstructor(instructor instructors.Instructor) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockInstructorRepository) GetInstructorByID(id primitive.ObjectID) (*instructors.Instructor, error) {
	return &instructors.Instructor{}, nil
}

func (m *MockInstructorRepository) GetInstructorByEmail(email string) (*instructors.Instructor, error) {
	return &instructors.Instructor{}, nil
}

func (m *MockInstructorRepository) UpdateInstructor(instructor instructors.Instructor) error {
	return nil
}

func (m *MockInstructorRepository) DeleteInstructor(id primitive.ObjectID) error {
	return nil
}
