package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mdb struct {
	Client     mongo.Client
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

func InitMongoDB() (*mdb, error) {
	fmt.Println("Connecting to mongodb")

	var mdb mdb

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return nil, err
	}
	mdb.Client = *mongoClient

	mdb.Ctx, mdb.CancelFunc = context.WithTimeout(context.Background(), 2000*time.Second)

	err = mdb.Client.Connect(mdb.Ctx)
	if err != nil {
		return nil, err
	}

	err = mdb.Client.Ping(mdb.Ctx, nil)
	if err != nil {
		return nil, err
	}

	return &mdb, nil
}
