package lessons

import "go.mongodb.org/mongo-driver/bson/primitive"

type LessonRepository interface {
	CreateLesson(lesson Lesson) (primitive.ObjectID, error)
	GetLessonByID(id primitive.ObjectID) (*Lesson, error)
	UpdateLesson(lesson Lesson) error
	DeleteLesson(id primitive.ObjectID) error
}
