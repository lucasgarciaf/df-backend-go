package availability

import "go.mongodb.org/mongo-driver/bson/primitive"

type AvailabilityService struct {
	repo AvailabilityRepository
}

func NewAvailabilityService(repo AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{repo: repo}
}

func (s *AvailabilityService) CreateAvailability(availability Availability) (primitive.ObjectID, error) {
	return s.repo.CreateAvailability(availability)
}

func (s *AvailabilityService) GetAvailabilityByID(id primitive.ObjectID) (*Availability, error) {
	return s.repo.GetAvailabilityByID(id)
}

func (s *AvailabilityService) UpdateAvailability(availability Availability) error {
	return s.repo.UpdateAvailability(availability)
}

func (s *AvailabilityService) DeleteAvailability(id primitive.ObjectID) error {
	return s.repo.DeleteAvailability(id)
}
