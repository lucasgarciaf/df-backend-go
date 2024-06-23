package instructors

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instructor struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `json:"username"`
	FirstName    string             `json:"firstName"`
	LastName     string             `json:"lastName"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"-"`
	Role         string             `json:"role"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
