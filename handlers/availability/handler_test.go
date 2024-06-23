package availability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/availability"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateAvailability(t *testing.T) {
	r := setupRouter()
	availabilityService := availability.NewAvailabilityService(&MockAvailabilityRepository{})
	availabilityHandler := NewAvailabilityHandler(availabilityService)

	r.POST("/availability", availabilityHandler.CreateAvailability)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/availability", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %v but got %v", http.StatusCreated, w.Code)
	}
}

type MockAvailabilityRepository struct{}

func (m *MockAvailabilityRepository) CreateAvailability(availability availability.Availability) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockAvailabilityRepository) GetAvailabilityByID(id primitive.ObjectID) (*availability.Availability, error) {
	return &availability.Availability{}, nil
}

func (m *MockAvailabilityRepository) UpdateAvailability(availability availability.Availability) error {
	return nil
}

func (m *MockAvailabilityRepository) DeleteAvailability(id primitive.ObjectID) error {
	return nil
}
