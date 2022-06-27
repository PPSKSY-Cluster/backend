package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MDB struct {
	Client     mongo.Client
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

var mdbInstance MDB

func InitDB() error {
	fmt.Println("Connecting to mongodb")

	var mdb MDB

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return err
	}
	mdb.Client = *mongoClient

	mdb.Ctx, mdb.CancelFunc = context.WithTimeout(context.Background(), 2000*time.Second)

	err = mdb.Client.Connect(mdb.Ctx)
	if err != nil {
		return err
	}

	err = mdb.Client.Ping(mdb.Ctx, nil)
	if err != nil {
		return err
	}

	mdbInstance = mdb

	if err = setupCollections(); err != nil {
		return err
	}

	return nil
}

func runQuery[T any](f func() (T, error)) (T, error) {
	return f()
}

func setupCollections() error {
	// create unique index over username
	if err := createIndex("users", "username", true); err != nil {
		return err
	}

	// setup validation
	singleRes := mdbInstance.Client.Database(os.Getenv("DB_NAME")).RunCommand(mdbInstance.Ctx, bson.D{
		{Key: "collMod", Value: "users"},
		{Key: "validator", Value: userValidator},
		{Key: "validationLevel", Value: "strict"},
	})
	if singleRes != nil {
		return singleRes.Err()
	}

	return nil
}

func createIndex(collectionName string, field string, unique bool) error {

	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(unique),
	}

	collection := mdbInstance.Client.Database(os.Getenv("DB_NAME")).Collection(collectionName)

	_, err := collection.Indexes().CreateOne(mdbInstance.Ctx, mod)
	if err != nil {
		return err
	}

	return nil
}
