package vehicles

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoVehicleRepository struct {
	db *mongo.Collection
}

func NewMongoVehicleRepository(db *mongo.Database) *MongoVehicleRepository {
	return &MongoVehicleRepository{
		db: db.Collection("vehicles"),
	}
}

func (r *MongoVehicleRepository) CreateVehicle(vehicle Vehicle) (primitive.ObjectID, error) {
	vehicle.ID = primitive.NewObjectID()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), vehicle)
	return vehicle.ID, err
}

func (r *MongoVehicleRepository) GetVehicleByID(id primitive.ObjectID) (*Vehicle, error) {
	var vehicle Vehicle
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&vehicle)
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *MongoVehicleRepository) UpdateVehicle(vehicle Vehicle) error {
	vehicle.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": vehicle.ID}, bson.M{"$set": vehicle})
	return err
}

func (r *MongoVehicleRepository) DeleteVehicle(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
