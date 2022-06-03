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
	client     mongo.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func InitMongoDB() (*mdb, error) {
	fmt.Println("Connecting to mongodb")

	var mdb mdb

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return nil, err
	}
	mdb.client = *mongoClient

	mdb.ctx, mdb.cancelFunc = context.WithTimeout(context.Background(), 2000*time.Second)

	err = mdb.client.Connect(mdb.ctx)
	if err != nil {
		return nil, err
	}

	return &mdb, nil
}
