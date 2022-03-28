package store

import (
	"context"
	"fmt"

	"github.com/hamed-lohi/user-management/db"
	"github.com/hamed-lohi/user-management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore struct {
	// client     *db.MongoClient
	// collection *mongo.Collection
}

func NewUserStore() *UserStore {
	//cli.Database("user_management_db").Collection("users")
	return &UserStore{
		// client:     db.New(),
		// collection: client.,
	}
}

func (us *UserStore) GetUserList() (*[]model.User, error) {

	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	var users []model.User
	cursor, err := coll.Find(context.TODO(), bson.D{})
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
	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	var m model.User
	fmt.Println(id)
	if err := coll.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByEmail(e string) (*model.User, error) {
	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	var m model.User
	if err := coll.FindOne(context.TODO(), bson.D{{"email", e}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	var m model.User
	if err := coll.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) Create(u *model.User) error {
	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	if _, err := coll.InsertOne(context.TODO(), u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Update(u *model.User) error {

	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	filter := bson.M{"_id": u.ID}
	if _, err := coll.ReplaceOne(context.TODO(), filter, u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Delete(id primitive.ObjectID) error {
	cli := db.New()
	defer cli.Dispose()
	coll := cli.GetCollection(db.Users)

	filter := bson.M{"_id": id}
	if _, err := coll.DeleteOne(context.TODO(), filter); err != nil {
		return err
	}
	return nil
}
