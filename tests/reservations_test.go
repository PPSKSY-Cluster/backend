package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/PPSKSY-Cluster/backend/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_reservations(t *testing.T) {
	fmt.Println("Testing reservations... ")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	//TODO: Change the mail endpoint!
	user := db.User{Username: "foo", Password: "bar"}
	tokenStr, createdUser := createUserAndLogin(t, app, user, false)

	start := time.Now()
	end := start.Add(time.Minute * 3)

	//TODO: Actually think about how to test mail functionality -> Tests quit running, before task is scheduled
	diff := (float64(end.Unix()) - float64(start.Unix())) * 0.9
	mailTime := time.Unix(start.Unix()+int64(diff), 0)

	fmt.Println("The time the mail should arrive at: " + mailTime.String())

	foosReservation := db.Reservation{
		ClusterID: primitive.NewObjectID(),
		Nodes:     10,
		UserID:    createdUser.ID,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		IsExpired: false,
	}

	start = start.Truncate(time.Hour * 24)
	end = end.Truncate(time.Hour * 72)

	expiredReservation := db.Reservation{
		ClusterID: primitive.NewObjectID(),
		Nodes:     10,
		UserID:    createdUser.ID,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		IsExpired: false,
	}

	createOneTest := TestReq{
		description:  "Create one Reservation (expect 201)",
		expectedCode: 201,
		route:        "/api/reservations/",
		method:       "POST",
		body:         foosReservation,
		expectedData: foosReservation,
	}
	createdReservation := executeTestReq[db.Reservation](t, app, createOneTest, tokenStr)

	getAllTest := TestReq{
		description:  "Get all reservations (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/",
		method:       "GET",
		body:         nil,
		expectedData: []db.Reservation{createdReservation},
	}
	executeTestReq[[]db.Reservation](t, app, getAllTest, tokenStr)

	getAllClusterTest := TestReq{
		description:  "Get all reservations with given cluster id (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/clusters/" + createdReservation.ClusterID.Hex(),
		method:       "GET",
		body:         nil,
		expectedData: []db.Reservation{createdReservation},
	}
	executeTestReq[[]db.Reservation](t, app, getAllClusterTest, tokenStr)

	getAllUserTest := TestReq{
		description:  "Get all reservations with given uid (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/users/" + createdReservation.UserID.Hex(),
		method:       "GET",
		body:         nil,
		expectedData: []db.Reservation{createdReservation},
	}
	executeTestReq[[]db.Reservation](t, app, getAllUserTest, tokenStr)

	getOneTest := TestReq{
		description:  "Get one reservation (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/" + createdReservation.ID.Hex(),
		method:       "GET",
		body:         nil,
		expectedData: createdReservation,
	}
	executeTestReq[db.Reservation](t, app, getOneTest, tokenStr)

	editedReservation := createdReservation
	editedReservation.Nodes = 15
	editOneTest := TestReq{
		description:  "Edit one reservation (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/" + createdReservation.ID.Hex(),
		method:       "PUT",
		body:         editedReservation,
		expectedData: editedReservation,
	}
	executeTestReq[db.Reservation](t, app, editOneTest, tokenStr)

	deleteOneTest := TestReq{
		description:  "Delete one reservation (expect 204)",
		expectedCode: 204,
		route:        "/api/reservations/" + createdReservation.ID.Hex(),
		method:       "DELETE",
		body:         nil,
		expectedData: nil,
	}
	executeTestReq[db.Reservation](t, app, deleteOneTest, tokenStr)

	expiredAddTest := TestReq{
		description:  "Add one already expired reservation (expect 201)",
		expectedCode: 201,
		route:        "/api/reservations/",
		method:       "POST",
		body:         expiredReservation,
		expectedData: expiredReservation,
	}
	expiredReservation = executeTestReq[db.Reservation](t, app, expiredAddTest, tokenStr)

	expiredReservation.IsExpired = true
	expiredGetTest := TestReq{
		description:  "Retrieve expired reservation with isExpired now set to true (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/" + expiredReservation.ID.Hex(),
		method:       "GET",
		body:         nil,
		expectedData: expiredReservation,
	}
	executeTestReq[db.Reservation](t, app, expiredGetTest, tokenStr)
}
