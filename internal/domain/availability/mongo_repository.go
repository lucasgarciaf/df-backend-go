package availability

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAvailabilityRepository struct {
	db *mongo.Collection
}

func NewMongoAvailabilityRepository(db *mongo.Database) *MongoAvailabilityRepository {
	return &MongoAvailabilityRepository{
		db: db.Collection("availability"),
	}
}

func (r *MongoAvailabilityRepository) CreateAvailability(availability Availability) (primitive.ObjectID, error) {
	availability.ID = primitive.NewObjectID()
	availability.CreatedAt = time.Now()
	availability.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), availability)
	return availability.ID, err
}

func (r *MongoAvailabilityRepository) GetAvailabilityByID(id primitive.ObjectID) (*Availability, error) {
	var availability Availability
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&availability)
	if err != nil {
		return nil, err
	}
	return &availability, nil
}

func (r *MongoAvailabilityRepository) UpdateAvailability(availability Availability) error {
	availability.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": availability.ID}, bson.M{"$set": availability})
	return err
}

func (r *MongoAvailabilityRepository) DeleteAvailability(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
