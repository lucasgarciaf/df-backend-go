package students

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/students"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateStudent(t *testing.T) {
	r := setupRouter()
	studentService := students.NewStudentService(&MockStudentRepository{})
	studentHandler := NewStudentHandler(studentService)

	r.POST("/students", studentHandler.CreateStudent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/students", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %v but got %v", http.StatusCreated, w.Code)
	}
}

type MockStudentRepository struct{}

func (m *MockStudentRepository) CreateStudent(student students.Student) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockStudentRepository) GetStudentByID(id primitive.ObjectID) (*students.Student, error) {
	return &students.Student{}, nil
}

func (m *MockStudentRepository) GetStudentByEmail(email string) (*students.Student, error) {
	return &students.Student{}, nil
}

func (m *MockStudentRepository) UpdateStudent(student students.Student) error {
	return nil
}

func (m *MockStudentRepository) DeleteStudent(id primitive.ObjectID) error {
	return nil
}
