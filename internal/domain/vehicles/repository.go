package vehicles

import "go.mongodb.org/mongo-driver/bson/primitive"

type VehicleRepository interface {
	CreateVehicle(vehicle Vehicle) (primitive.ObjectID, error)
	GetVehicleByID(id primitive.ObjectID) (*Vehicle, error)
	UpdateVehicle(vehicle Vehicle) error
	DeleteVehicle(id primitive.ObjectID) error
}
