package lessons

import "go.mongodb.org/mongo-driver/bson/primitive"

type LessonService struct {
	repo LessonRepository
}

func NewLessonService(repo LessonRepository) *LessonService {
	return &LessonService{repo: repo}
}

func (s *LessonService) CreateLesson(lesson Lesson) (primitive.ObjectID, error) {
	return s.repo.CreateLesson(lesson)
}

func (s *LessonService) GetLessonByID(id primitive.ObjectID) (*Lesson, error) {
	return s.repo.GetLessonByID(id)
}

func (s *LessonService) UpdateLesson(lesson Lesson) error {
	return s.repo.UpdateLesson(lesson)
}

func (s *LessonService) DeleteLesson(id primitive.ObjectID) error {
	return s.repo.DeleteLesson(id)
}
