package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type Reservation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ClusterID primitive.ObjectID `bson:"clusterID" json:"clusterID"`
	UserID    primitive.ObjectID `bson:"userID" json:"userID"`
	Nodes     int                `bson:"nodes" json:"nodes"`
	StartTime int64              `bson:"startTime" json:"startTime"`
	EndTime   int64              `bson:"endTime" json:"endTime"`
	IsExpired bool               `bson:"isExpired" json:"isExpired"`
}

func GetAllReservations() ([]Reservation, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			Find(mdbInstance.Ctx, bson.M{})
	}

	reservationCursor, err := runQuery[*mongo.Cursor](query)
	if err != nil {
		return nil, err
	}

	var reservations []Reservation
	if err = reservationCursor.All(mdbInstance.Ctx, &reservations); err != nil {
		return nil, err
	}

	for i := 0; i < len(reservations); i++ {
		r := reservations[i]
		r.IsExpired = CheckExpired(r)
		reservations[i] = r
	}

	return reservations, nil
}

func GetReservationById(_id primitive.ObjectID) (Reservation, error) {

	query := func() (*mongo.SingleResult, error) {
		singleRes := mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			FindOne(mdbInstance.Ctx, bson.M{"_id": _id})
		return singleRes, singleRes.Err()
	}

	reservationSingleRes, err := runQuery[*mongo.SingleResult](query)
	if err != nil {
		return Reservation{}, err
	}

	var reservation Reservation
	if err = reservationSingleRes.Decode(&reservation); err != nil {
		return Reservation{}, err
	}

	reservation.IsExpired = CheckExpired(reservation)

	return reservation, nil
}

func GetReservationsByClusterId(clusterID primitive.ObjectID) ([]Reservation, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			Find(mdbInstance.Ctx, bson.M{"clusterID": clusterID})
	}

	reservationCursor, err := runQuery[*mongo.Cursor](query)
	if err != nil {
		return nil, err
	}

	var reservations []Reservation
	if err = reservationCursor.All(mdbInstance.Ctx, &reservations); err != nil {
		return nil, err
	}

	for i := 0; i < len(reservations); i++ {
		r := reservations[i]
		r.IsExpired = CheckExpired(r)
		reservations[i] = r
	}

	return reservations, nil
}

func GetReservationsByUserId(userID primitive.ObjectID) ([]Reservation, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			Find(mdbInstance.Ctx, bson.M{"userID": userID})
	}

	reservationCursor, err := runQuery[*mongo.Cursor](query)
	if err != nil {
		return nil, err
	}

	var reservations []Reservation
	if err = reservationCursor.All(mdbInstance.Ctx, &reservations); err != nil {
		return nil, err
	}

	for i := 0; i < len(reservations); i++ {
		r := reservations[i]
		r.IsExpired = CheckExpired(r)
		reservations[i] = r
	}

	return reservations, nil
}

func AddReservation(reservation Reservation) (Reservation, error) {

	query := func() (*mongo.InsertOneResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			InsertOne(mdbInstance.Ctx, reservation)
	}

	if err := mdbInstance.Validate.Struct(reservation); err != nil {
		return Reservation{}, err
	}

	insertRes, err := runQuery[*mongo.InsertOneResult](query)
	if err != nil {
		return Reservation{}, err
	}

	reservation.ID = insertRes.InsertedID.(primitive.ObjectID)

	return reservation, nil
}

func EditReservation(_id primitive.ObjectID, reservation Reservation) (Reservation, error) {

	query := func() (*mongo.UpdateResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			UpdateOne(mdbInstance.Ctx,
				bson.M{"_id": _id},
				bson.M{"$set": reservation})
	}

	_, err := runQuery[*mongo.UpdateResult](query)
	if err != nil {
		return Reservation{}, err
	}

	return reservation, nil
}

func DeleteReservation(_id primitive.ObjectID) error {

	query := func() (*mongo.DeleteResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection("reservations").
			DeleteOne(mdbInstance.Ctx, bson.M{"_id": _id})
	}

	_, err := runQuery[*mongo.DeleteResult](query)
	if err != nil {
		return err
	}

	return nil
}

func CheckExpired(reservation Reservation) bool {
	return time.Now().After(time.UnixMicro(reservation.EndTime))
}
