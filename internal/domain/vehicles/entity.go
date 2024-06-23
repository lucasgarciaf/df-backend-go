package vehicles

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Make         string             `bson:"make,omitempty"`
	Model        string             `bson:"model,omitempty"`
	Year         int                `bson:"year,omitempty"`
	LicensePlate string             `bson:"license_plate,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
}
