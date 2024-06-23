package courses

import "go.mongodb.org/mongo-driver/bson/primitive"

type CourseRepository interface {
	CreateCourse(course Course) (primitive.ObjectID, error)
	GetCourseByID(id primitive.ObjectID) (*Course, error)
	UpdateCourse(course Course) error
	DeleteCourse(id primitive.ObjectID) error
}
