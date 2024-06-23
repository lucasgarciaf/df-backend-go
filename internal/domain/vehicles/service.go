package vehicles

import "go.mongodb.org/mongo-driver/bson/primitive"

type VehicleService struct {
	repo VehicleRepository
}

func NewVehicleService(repo VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) CreateVehicle(vehicle Vehicle) (primitive.ObjectID, error) {
	return s.repo.CreateVehicle(vehicle)
}

func (s *VehicleService) GetVehicleByID(id primitive.ObjectID) (*Vehicle, error) {
	return s.repo.GetVehicleByID(id)
}

func (s *VehicleService) UpdateVehicle(vehicle Vehicle) error {
	return s.repo.UpdateVehicle(vehicle)
}

func (s *VehicleService) DeleteVehicle(id primitive.ObjectID) error {
	return s.repo.DeleteVehicle(id)
}
