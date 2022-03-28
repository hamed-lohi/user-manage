package user

import (
	"github.com/hamed-lohi/user-management/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	GetByID(primitive.ObjectID) (*model.User, error)

	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	GetUserList() (*[]model.User, error)
	Create(*model.User) error
	Update(*model.User) error
	Delete(id primitive.ObjectID) error
}
