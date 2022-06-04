package db

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name"`
}

func GetAllUsers() ([]User, error) {

	query := func() (interface{}, error) {
		return mdbInstance.Client.Database(os.Getenv("DB_NAME")).Collection("users").Find(mdbInstance.Ctx, bson.M{})
	}

	usersCursor, err := runQueryToCursor(query)
	if err != nil {
		return nil, err
	}

	var users []User
	if err = usersCursor.All(mdbInstance.Ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
