package store

import (
	"fmt"

	"github.com/hamed-lohi/user-management/db"
	"github.com/hamed-lohi/user-management/model"
	"github.com/hamed-lohi/user-management/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	dbProvider *db.DBProvider
	collection *mongo.Collection
}

// Verify Interface Compliance
var _ user.Store = (*UserStore)(nil)

func NewUserStore(dp *db.DBProvider) *UserStore {
	//cli.Database("user_management_db").Collection("users")
	return &UserStore{
		dbProvider: dp,
		collection: dp.GetCollection(db.Users),
	}
}

func (us *UserStore) GetUserList() (*[]model.User, error) {

	coll := us.collection
	// coll := cli.GetCollection(db.Users)

	var users []model.User
	cursor, err := coll.Find(us.dbProvider.Context, bson.D{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(us.dbProvider.Context) {
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
	if err := us.collection.FindOne(us.dbProvider.Context, bson.D{{"_id", id}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByEmail(e string) (*model.User, error) {
	var m model.User
	if err := us.collection.FindOne(us.dbProvider.Context, bson.D{{"email", e}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := us.collection.FindOne(us.dbProvider.Context, bson.D{{"username", username}}).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) Create(u *model.User) error {

	if _, err := us.collection.InsertOne(us.dbProvider.Context, u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Update(u *model.User) error {
	filter := bson.M{"_id": u.ID}
	if _, err := us.collection.ReplaceOne(us.dbProvider.Context, filter, u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	if _, err := us.collection.DeleteOne(us.dbProvider.Context, filter); err != nil {
		return err
	}
	return nil
}
