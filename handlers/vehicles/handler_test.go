package vehicles

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/vehicles"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateVehicle(t *testing.T) {
	r := setupRouter()
	vehicleService := vehicles.NewVehicleService(&MockVehicleRepository{})
	vehicleHandler := NewVehicleHandler(vehicleService)

	r.POST("/vehicles", vehicleHandler.CreateVehicle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vehicles", nil) // Add appropriate body content for the request

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %v but got %v", http.StatusCreated, w.Code)
	}
}

type MockVehicleRepository struct{}

func (m *MockVehicleRepository) CreateVehicle(vehicle vehicles.Vehicle) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockVehicleRepository) GetVehicleByID(id primitive.ObjectID) (*vehicles.Vehicle, error) {
	return &vehicles.Vehicle{}, nil
}

func (m *MockVehicleRepository) UpdateVehicle(vehicle vehicles.Vehicle) error {
	return nil
}

func (m *MockVehicleRepository) DeleteVehicle(id primitive.ObjectID) error {
	return nil
}
