package courses

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCourseRepository struct {
	db *mongo.Collection
}

func NewMongoCourseRepository(db *mongo.Database) *MongoCourseRepository {
	return &MongoCourseRepository{
		db: db.Collection("courses"),
	}
}

func (r *MongoCourseRepository) CreateCourse(course Course) (primitive.ObjectID, error) {
	course.ID = primitive.NewObjectID()
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), course)
	return course.ID, err
}

func (r *MongoCourseRepository) GetCourseByID(id primitive.ObjectID) (*Course, error) {
	var course Course
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *MongoCourseRepository) UpdateCourse(course Course) error {
	course.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": course.ID}, bson.M{"$set": course})
	return err
}

func (r *MongoCourseRepository) DeleteCourse(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
