package instructors

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoInstructorRepository struct {
	db *mongo.Collection
}

func NewMongoInstructorRepository(db *mongo.Database) *MongoInstructorRepository {
	return &MongoInstructorRepository{
		db: db.Collection("instructors"),
	}
}

func (r *MongoInstructorRepository) CreateInstructor(instructor Instructor) (primitive.ObjectID, error) {
	instructor.ID = primitive.NewObjectID()
	instructor.CreatedAt = time.Now()
	instructor.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), instructor)
	return instructor.ID, err
}

func (r *MongoInstructorRepository) GetInstructorByID(id primitive.ObjectID) (*Instructor, error) {
	var instructor Instructor
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&instructor)
	if err != nil {
		return nil, err
	}
	return &instructor, nil
}

func (r *MongoInstructorRepository) GetInstructorByEmail(email string) (*Instructor, error) {
	var instructor Instructor
	err := r.db.FindOne(context.Background(), bson.M{"email": email}).Decode(&instructor)
	if err != nil {
		return nil, err
	}
	return &instructor, nil
}

func (r *MongoInstructorRepository) UpdateInstructor(instructor Instructor) error {
	instructor.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": instructor.ID}, bson.M{"$set": instructor})
	return err
}

func (r *MongoInstructorRepository) DeleteInstructor(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
