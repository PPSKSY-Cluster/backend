package db

import (
	"context"
	"errors"
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

type Query func() (_ interface{}, _ error)

func runQuery(f Query) (interface{}, error) {
	return f()
}

func runQueryToCursor(query Query) (*mongo.Cursor, error) {
	userI, err := runQuery(query)
	if err != nil {
		return &mongo.Cursor{}, err
	}
	user, ok := userI.(*mongo.Cursor)
	if !ok {
		return &mongo.Cursor{}, errors.New("Could not convert userI to *mongo.Cursor")
	}
	return user, nil
}

func runQueryToSingleRes(query Query) (*mongo.SingleResult, error) {
	userI, err := runQuery(query)
	if err != nil {
		return &mongo.SingleResult{}, err
	}
	user, ok := userI.(*mongo.SingleResult)
	if !ok {
		return &mongo.SingleResult{}, errors.New("Could not convert userI to *mongo.SingleResult")
	}
	return user, nil
}
