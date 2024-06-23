package courses

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/courses"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateCourse(t *testing.T) {
	r := setupRouter()
	courseService := courses.NewCourseService(&MockCourseRepository{})
	courseHandler := NewCourseHandler(courseService)

	r.POST("/courses", courseHandler.CreateCourse)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/courses", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %v but got %v", http.StatusCreated, w.Code)
	}
}

type MockCourseRepository struct{}

func (m *MockCourseRepository) CreateCourse(course courses.Course) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockCourseRepository) GetCourseByID(id primitive.ObjectID) (*courses.Course, error) {
	return &courses.Course{}, nil
}

func (m *MockCourseRepository) UpdateCourse(course courses.Course) error {
	return nil
}

func (m *MockCourseRepository) DeleteCourse(id primitive.ObjectID) error {
	return nil
}
