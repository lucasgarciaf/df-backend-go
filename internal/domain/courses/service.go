package courses

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) CreateCourse(course Course) (primitive.ObjectID, error) {
	course.ID = primitive.NewObjectID()
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()
	return s.repo.CreateCourse(course)
}

func (s *CourseService) GetCourseByID(id primitive.ObjectID) (*Course, error) {
	return s.repo.GetCourseByID(id)
}

func (s *CourseService) UpdateCourse(course Course) error {
	course.UpdatedAt = time.Now()
	return s.repo.UpdateCourse(course)
}

func (s *CourseService) DeleteCourse(id primitive.ObjectID) error {
	return s.repo.DeleteCourse(id)
}
