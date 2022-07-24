//go:build test
// +build test

// this file will not be included unless it's built for testing
package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

// this function is very dangerous because it drops the entire database
// (which is a functionality we'd like to have for running tests)
func DropDB(dbName string) {
	mdbInstance.Client.Database(dbName).Drop(mdbInstance.Ctx)
}

// this allows to automatically create admins and super admins
// for testing (just the first line of AddUser() is missing)
func AddUserWithType(user User) (User, error) {
	query := func() (*mongo.InsertOneResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			InsertOne(mdbInstance.Ctx, user)
	}

	if err := mdbInstance.Validate.Struct(user); err != nil {
		return User{}, err
	}

	insertRes, err := runQuery[*mongo.InsertOneResult](query)
	if err != nil {
		return User{}, err
	}

	user.ID = insertRes.InsertedID.(primitive.ObjectID)
	user.Password = ""

	return user, nil
}
