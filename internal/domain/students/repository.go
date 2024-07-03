package students

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentRepository interface {
	CreateStudent(student Student) (primitive.ObjectID, error)
	GetStudentByID(id primitive.ObjectID) (*Student, error)
	GetStudentByEmail(email string) (*Student, error)
	UpdateStudent(student Student) error
	DeleteStudent(id primitive.ObjectID) error
	GetAllStudents() ([]Student, error)
}
