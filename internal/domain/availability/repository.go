package availability

import "go.mongodb.org/mongo-driver/bson/primitive"

type AvailabilityRepository interface {
	CreateAvailability(availability Availability) (primitive.ObjectID, error)
	GetAvailabilityByID(id primitive.ObjectID) (*Availability, error)
	UpdateAvailability(availability Availability) error
	DeleteAvailability(id primitive.ObjectID) error
}
