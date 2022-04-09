package db

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const _dbName = "user_management_db"

type Table string

const (
	Users     Table = "users"
	Customers Table = "customers"
)

type DBProvider struct {
	Client     *mongo.Client
	Context    context.Context
	CancelFunc context.CancelFunc
}

//var _ mongo.Client = (*MongoClient)(nil)

func NewDBProvider() *DBProvider {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	// if err != nil {
	// 	panic(err)
	// }

	// Get Client, Context, CalcelFunc and
	// err from connect method.
	client, ctx, cancel, err := connect(uri) // "mongodb://localhost:27017"
	if err != nil {
		panic(err)
	}

	mdb := &DBProvider{client, ctx, cancel}

	// Ping mongoDB with Ping method
	mdb.ping()

	//Seed(mdb)

	return mdb
}

func (cl *DBProvider) GetCollection(collName Table) *mongo.Collection {
	return cl.Client.Database(_dbName).Collection(string(collName))
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), 30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func (dp *DBProvider) close() {

	// CancelFunc to cancel to context
	defer dp.CancelFunc() //cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := dp.Client.Disconnect(dp.Context); err != nil {
			panic(err)
		}
		fmt.Println("****** Disposing")
	}()
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func (dp *DBProvider) ping() error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := dp.Client.Ping(dp.Context, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}
