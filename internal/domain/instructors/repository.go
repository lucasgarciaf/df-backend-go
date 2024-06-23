package instructors

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructorRepository interface {
	CreateInstructor(instructor Instructor) (primitive.ObjectID, error)
	GetInstructorByEmail(email string) (*Instructor, error)
	GetInstructorByID(id primitive.ObjectID) (*Instructor, error)
	UpdateInstructor(instructor Instructor) error
	DeleteInstructor(id primitive.ObjectID) error
}
