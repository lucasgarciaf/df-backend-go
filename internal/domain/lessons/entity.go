package lessons

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lesson struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CourseID     primitive.ObjectID `bson:"course_id,omitempty"`
	InstructorID primitive.ObjectID `bson:"instructor_id,omitempty"`
	Title        string             `bson:"title,omitempty"`
	Description  string             `bson:"description,omitempty"`
	Schedule     time.Time          `bson:"schedule,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
}
