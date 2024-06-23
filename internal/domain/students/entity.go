package students

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName    string             `bson:"first_name,omitempty" json:"first_name"`
	LastName     string             `bson:"last_name,omitempty" json:"last_name"`
	Username     string             `bson:"username,omitempty" json:"username"`
	Email        string             `bson:"email,omitempty" json:"email"`
	Age          int                `bson:"age,omitempty" json:"age"`
	PasswordHash string             `bson:"password_hash,omitempty" json:"password_hash"`
	Role         string             `bson:"role,omitempty" json:"role"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}
