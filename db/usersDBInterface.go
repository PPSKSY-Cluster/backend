package db

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsers() ([]bson.M, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.Database(os.Getenv("DB_NAME")).Collection("users").Find(mdbInstance.Ctx, bson.M{})
	}

	usersCursor, err := RunQuery(query)
	if err != nil {
		return nil, err
	}

	var users []bson.M
	if err = usersCursor.All(mdbInstance.Ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
