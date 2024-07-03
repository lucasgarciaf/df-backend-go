package students

import (
	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/students"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

type MockStudentRepository struct{}

func (m *MockStudentRepository) CreateStudent(student students.Student) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *MockStudentRepository) GetStudentByID(id primitive.ObjectID) (*students.Student, error) {
	return &students.Student{}, nil
}

func (m *MockStudentRepository) GetAllStudents() ([]students.Student, error) {
	return []students.Student{}, nil
}

func (m *MockStudentRepository) GetStudentByEmail(email string) (*students.Student, error) {
	return &students.Student{}, nil
}

func (m *MockStudentRepository) UpdateStudent(student students.Student) error {
	return nil
}

func (m *MockStudentRepository) DeleteStudent(id primitive.ObjectID) error {
	return nil
}
