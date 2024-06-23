package lessons

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLessonRepository struct {
	db *mongo.Collection
}

func NewMongoLessonRepository(db *mongo.Database) *MongoLessonRepository {
	return &MongoLessonRepository{
		db: db.Collection("lessons"),
	}
}

func (r *MongoLessonRepository) CreateLesson(lesson Lesson) (primitive.ObjectID, error) {
	lesson.ID = primitive.NewObjectID()
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), lesson)
	return lesson.ID, err
}

func (r *MongoLessonRepository) GetLessonByID(id primitive.ObjectID) (*Lesson, error) {
	var lesson Lesson
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&lesson)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *MongoLessonRepository) UpdateLesson(lesson Lesson) error {
	lesson.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": lesson.ID}, bson.M{"$set": lesson})
	return err
}

func (r *MongoLessonRepository) DeleteLesson(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
