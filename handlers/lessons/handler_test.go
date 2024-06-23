package lessons

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/lessons"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateLesson(t *testing.T) {
	r := setupRouter()
	lessonService := lessons.NewLessonService(&MockLessonRepository{})
	lessonHandler := NewLessonHandler(lessonService)

	r.POST("/lessons", lessonHandler.CreateLesson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/lessons", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %v but got %v", http.StatusCreated, w.Code)
	}
}

func TestGetLessonByID(t *testing.T) {
	r := setupRouter()
	lessonService := lessons.NewLessonService(&MockLessonRepository{})
	lessonHandler := NewLessonHandler(lessonService)

	r.GET("/lessons/:id", lessonHandler.GetLessonByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lessons/60c72b2f9b1d8b6a8f8a53e0", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %v but got %v", http.StatusOK, w.Code)
	}
}

func TestUpdateLesson(t *testing.T) {
	r := setupRouter()
	lessonService := lessons.NewLessonService(&MockLessonRepository{})
	lessonHandler := NewLessonHandler(lessonService)

	r.PUT("/lessons/:id", lessonHandler.UpdateLesson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/lessons/60c72b2f9b1d8b6a8f8a53e0", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %v but got %v", http.StatusOK, w.Code)
	}
}

func TestDeleteLesson(t *testing.T) {
	r := setupRouter()
	lessonService := lessons.NewLessonService(&MockLessonRepository{})
	lessonHandler := NewLessonHandler(lessonService)

	r.DELETE("/lessons/:id", lessonHandler.DeleteLesson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/lessons/60c72b2f9b1d8b6a8f8a53e0", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected status %v but got %v", http.StatusNoContent, w.Code)
	}
}

type MockLessonRepository struct{}

func (m *MockLessonRepository) CreateLesson(lesson lessons.Lesson) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockLessonRepository) GetLessonByID(id primitive.ObjectID) (*lessons.Lesson, error) {
	return &lessons.Lesson{}, nil
}

func (m *MockLessonRepository) UpdateLesson(lesson lessons.Lesson) error {
	return nil
}

func (m *MockLessonRepository) DeleteLesson(id primitive.ObjectID) error {
	return nil
}
