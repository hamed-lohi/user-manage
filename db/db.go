package db

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/hamed-lohi/user-management/model"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	*mongo.Client
}

func New() *MongoClient {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &MongoClient{client}
}

func (cl *MongoClient) Dispose() {
	if err := cl.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (cl *MongoClient) Seed() {

	coll := cl.Database("user_management_db").Collection("users")

	var result model.User
	err := coll.FindOne(context.TODO(), bson.D{{"username", "Admin"}}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			admin := &model.User{
				//ID:       ,
				Username: "Admin",
				Email:    "admin@gmail.com",
				//Password: "aaaa",
				Bio:   new(string),
				Roles: []model.Role{model.Admin},
			}
			admin.SetPassword("aaa")
			coll.InsertOne(context.TODO(), admin)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("found document %v", result)

	// if err := cl.Disconnect(context.TODO()); err != nil {
	// 	panic(err)
	// }
}
