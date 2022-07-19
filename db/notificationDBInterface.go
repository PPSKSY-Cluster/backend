package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type ReservationNotification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ClusterID primitive.ObjectID `bson:"clusterID" json:"clusterID"`
	UserID    primitive.ObjectID `bson:"userID" json:"userID"`
}

func GetAllNotifications() ([]ReservationNotification, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservation_notification").
			Find(mdbInstance.Ctx, bson.M{})
	}

	notificationCursor, err := runQuery[*mongo.Cursor](query)
	if err != nil {
		return nil, err
	}

	var notifications []ReservationNotification
	if err = notificationCursor.All(mdbInstance.Ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func AddNotification(notification ReservationNotification) (ReservationNotification, error) {

	query := func() (*mongo.InsertOneResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservation_notification").
			InsertOne(mdbInstance.Ctx, notification)
	}

	if err := mdbInstance.Validate.Struct(notification); err != nil {
		return ReservationNotification{}, err
	}

	insertRes, err := runQuery[*mongo.InsertOneResult](query)
	if err != nil {
		return ReservationNotification{}, err
	}

	notification.ID = insertRes.InsertedID.(primitive.ObjectID)

	return notification, nil
}

func DeleteNotification(_id primitive.ObjectID) error {

	query := func() (*mongo.DeleteResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservation_notification").
			DeleteOne(mdbInstance.Ctx, bson.M{"_id": _id})
	}

	_, err := runQuery[*mongo.DeleteResult](query)
	if err != nil {
		return err
	}

	return nil
}
