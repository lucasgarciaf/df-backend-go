package availability

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Availability struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	InstructorID primitive.ObjectID `bson:"instructor_id,omitempty"`
	StartTime    time.Time          `bson:"start_time,omitempty"`
	EndTime      time.Time          `bson:"end_time,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
}
