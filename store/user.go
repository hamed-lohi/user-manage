package store

import (
	"context"
	"fmt"

	"github.com/hamed-lohi/user-management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	collection *mongo.Collection
}

func NewUserStore(coll *mongo.Collection) *UserStore {
	return &UserStore{
		collection: coll,
	}
}

func (us *UserStore) GetUserList() (*[]model.User, error) {
	var users []model.User
	cursor, err := us.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.User
		err := cursor.Decode(&elem)
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}

		users = append(users, elem)

	}

	if err := cursor.Err(); err != nil {
		//log.Fatal(err)
		return nil, err
	}

	return &users, nil
}

func (us *UserStore) GetByID(id primitive.ObjectID) (*model.User, error) {
	var m model.User
	fmt.Println(id)
	if err := us.collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByEmail(e string) (*model.User, error) {
	var m model.User
	if err := us.collection.FindOne(context.TODO(), bson.D{{"email", e}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := us.collection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) Create(u *model.User) error {

	if _, err := us.collection.InsertOne(context.TODO(), u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Update(u *model.User) error {
	filter := bson.M{"_id": u.ID}

	if _, err := us.collection.ReplaceOne(context.TODO(), filter, u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Delete(id primitive.ObjectID) error {

	filter := bson.M{"_id": id}

	if _, err := us.collection.DeleteOne(context.TODO(), filter); err != nil {
		return err
	}
	return nil
}
