package db

import (
	"context"
	"fmt"
	"os"
	"time"

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
	return nil
}

func runQuery[T any](f func() (T, error)) (T, error) {
	return f()
}