package admins

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminRepository interface {
	CreateAdmin(admin Admin) (primitive.ObjectID, error)
	GetAdminByID(id primitive.ObjectID) (*Admin, error)
	GetAdminByEmail(email string) (*Admin, error)
	UpdateAdmin(admin Admin) error
	DeleteAdmin(id primitive.ObjectID) error
}
