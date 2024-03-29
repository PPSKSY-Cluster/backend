package db

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserType int

const (
	UserUT       UserType = 0
	AdminUT      UserType = 1
	SuperAdminUT UserType = 2
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Mail     string             `bson:"mail" json:"mail"`
	Username string             `bson:"username" json:"username" validate:"required,min=3,max=20"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
	Type     UserType           `bson:"type" json:"type"`
}

var userDefaultProjection = bson.M{"password": 0}

func GetAllUsers() ([]User, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			Find(mdbInstance.Ctx, bson.M{}, options.Find().SetProjection(userDefaultProjection))
	}

	usersCursor, err := runQuery[*mongo.Cursor](query)
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

	query := func() (*mongo.SingleResult, error) {
		singleRes := mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			FindOne(mdbInstance.Ctx, bson.M{"_id": _id}, options.FindOne().SetProjection(userDefaultProjection))
		return singleRes, singleRes.Err()
	}

	userSingleRes, err := runQuery[*mongo.SingleResult](query)
	if err != nil {
		return User{}, err
	}

	var user User
	if err = userSingleRes.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserWithCredentials(username string) (User, error) {
	query := func() (*mongo.SingleResult, error) {
		singleRes := mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			FindOne(mdbInstance.Ctx, bson.M{"username": username})
		return singleRes, singleRes.Err()
	}

	userSingleRes, err := runQuery[*mongo.SingleResult](query)
	if err != nil {
		return User{}, err
	}

	var user User
	if err := userSingleRes.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func AddUser(user User) (User, error) {
	user.Type = UserUT // may be edited later by other users

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

	return user, nil
}

func EditUser(_id primitive.ObjectID, user User) (User, error) {

	query := func() (*mongo.UpdateResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			UpdateOne(mdbInstance.Ctx, bson.M{"_id": _id}, bson.M{"$set": user})
	}

	updateRes, err := runQuery[*mongo.UpdateResult](query)
	if err != nil && updateRes.ModifiedCount != 1 {
		return User{}, err
	}

	user.ID = _id
	return user, nil
}

func DeleteUser(_id primitive.ObjectID) error {

	query := func() (*mongo.DeleteResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("users").
			DeleteOne(mdbInstance.Ctx, bson.M{"_id": _id})
	}

	_, err := runQuery[*mongo.DeleteResult](query)
	if err != nil {
		return err
	}

	return nil
}
