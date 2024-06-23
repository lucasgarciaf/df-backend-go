package admins

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAdminRepository struct {
	db *mongo.Collection
}

func NewMongoAdminRepository(db *mongo.Database) *MongoAdminRepository {
	return &MongoAdminRepository{
		db: db.Collection("admins"),
	}
}

func (r *MongoAdminRepository) CreateAdmin(admin Admin) (primitive.ObjectID, error) {
	admin.ID = primitive.NewObjectID()
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()
	_, err := r.db.InsertOne(context.Background(), admin)
	return admin.ID, err
}

func (r *MongoAdminRepository) GetAdminByID(id primitive.ObjectID) (*Admin, error) {
	var admin Admin
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *MongoAdminRepository) GetAdminByEmail(email string) (*Admin, error) {
	var admin Admin
	err := r.db.FindOne(context.Background(), bson.M{"email": email}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *MongoAdminRepository) UpdateAdmin(admin Admin) error {
	admin.UpdatedAt = time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": admin.ID}, bson.M{"$set": admin})
	return err
}

func (r *MongoAdminRepository) DeleteAdmin(id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
