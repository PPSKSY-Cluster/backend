package db

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID   primitive.ObjectID `bson:"_id" 'json:"_id"`
	Name string             `bson:"name" 'json:"name"`
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

func GetUserById(_id primitive.ObjectID) (User, error) {

	query := func() (interface{}, error) {
		singleRes := mdbInstance.Client.Database(os.Getenv("DB_NAME")).Collection("users").FindOne(mdbInstance.Ctx, bson.M{"_id": _id})
		return singleRes, singleRes.Err()
	}

	userSingleRes, err := runQueryToSingleRes(query)
	if err != nil {
		return User{}, err
	}

	var user User
	if err = userSingleRes.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func AddUser(user User) (User, error) {

	query := func() (interface{}, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			InsertOne(mdbInstance.Ctx, user)
	}

	_, err := runQuery(query)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
