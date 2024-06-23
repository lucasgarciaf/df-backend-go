package students

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStudentRepository struct {
	db *mongo.Collection
}

func NewMongoStudentRepository(db *mongo.Database) *MongoStudentRepository {
	return &MongoStudentRepository{
		db: db.Collection("students"),
	}
}

func (r *MongoStudentRepository) CreateStudent(student Student) (primitive.ObjectID, error) {
	student.ID = primitive.NewObjectID()
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), student)
	return student.ID, err
}

func (r *MongoStudentRepository) GetStudentByID(id primitive.ObjectID) (*Student, error) {
	var student Student
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *MongoStudentRepository) GetStudentByEmail(email string) (*Student, error) {
	var student Student
	err := r.db.FindOne(context.Background(), bson.M{"email": email}).Decode(&student)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *MongoStudentRepository) UpdateStudent(student Student) error {
	student.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": student.ID}, bson.M{"$set": student})
	return err
}

func (r *MongoStudentRepository) DeleteStudent(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
