package students

import (
	"context"
	"log"
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
	result, err := r.db.InsertOne(context.Background(), student)
	if err != nil {
		log.Printf("Failed to insert student: %v", err)
		return primitive.NilObjectID, err
	}
	log.Printf("Student inserted successfully: %v", student)
	log.Printf("InsertOne result: %v", result)
	return student.ID, err
}

func (r *MongoStudentRepository) GetStudentByID(id primitive.ObjectID) (*Student, error) {
	var student Student
	log.Println("Looking for student ID: ", id)
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *MongoStudentRepository) GetAllStudents() ([]Student, error) {
	var students []Student
	cursor, err := r.db.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var student Student
		if err := cursor.Decode(&student); err != nil {
			log.Printf("Failed to decode student: %v", err)
			return nil, err
		}
		students = append(students, student)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return students, nil
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
