package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/PPSKSY-Cluster/backend/db"
)

func Test_reservations(t *testing.T) {
	fmt.Println("Testing reservations... ")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	user := db.User{Username: "foo", Password: "bar", Type: db.AdminUT}
	tokenStr, createdUser := createUserAndLogin(t, app, user, true)

	start := time.Now()
	end := start.Add(time.Second * 30)
	send := start.Add(time.Second * 27)

	fooscResource := db.CResource{
		Name:            "foos cresource",
		Description:     "this is foos first cresource",
		Nodes:           10,
		OperatingSystem: db.LinuxOS,
		Owner:           createdUser.ID,
	}

	createOnecResourceTest := TestReq{
		description:  "Create one resource (expect 201)",
		expectedCode: 201,
		route:        "/api/cresources/",
		method:       "POST",
		body:         fooscResource,
		expectedData: fooscResource,
	}
	createdCResource := executeTestReq[db.CResource](t, app, createOnecResourceTest, tokenStr)

	foosReservation := db.Reservation{
		ClusterID: createdCResource.ID,
		Nodes:     5,
		UserID:    createdUser.ID,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		IsExpired: false,
	}

	endBeforeStartReservation := db.Reservation{
		ClusterID: createdCResource.ID,
		Nodes:     5,
		UserID:    createdUser.ID,
		StartTime: end.Unix(),
		EndTime:   start.Unix(),
		IsExpired: false,
	}

	insufficienNodesReservation := db.Reservation{
		ClusterID: createdCResource.ID,
		Nodes:     10,
		UserID:    createdUser.ID,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		IsExpired: false,
	}

	start = start.Truncate(time.Hour * 24)
	end = end.Truncate(time.Hour * 24)

	expiredReservation := db.Reservation{
		ClusterID: createdCResource.ID,
		Nodes:     4,
		UserID:    createdUser.ID,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		IsExpired: false,
	}

	createOneReservationTest := TestReq{
		description:  "Create one Reservation (expect 201)",
		expectedCode: 201,
		route:        "/api/reservations/",
		method:       "POST",
		body:         foosReservation,
		expectedData: foosReservation,
	}
	createdReservation := executeTestReq[db.Reservation](t, app, createOneReservationTest, tokenStr)
	checkMailTime(t, send)

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
	editedReservation.Nodes = 6
	editOneTest := TestReq{
		description:  "Edit one reservation (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/" + createdReservation.ID.Hex(),
		method:       "PUT",
		body:         editedReservation,
		expectedData: editedReservation,
	}
	executeTestReq[db.Reservation](t, app, editOneTest, tokenStr)

	endBeforeStartTest := TestReq{
		description:  "Add one reservation where end < start",
		expectedCode: 400,
		route:        "/api/reservations/",
		method:       "POST",
		body:         endBeforeStartReservation,
		expectedData: nil,
	}
	_ = executeTestReq[db.Reservation](t, app, endBeforeStartTest, tokenStr)

	insufficienNodesTest := TestReq{
		description:  "Add one reservation with more nodes than the cluster has available",
		expectedCode: 400,
		route:        "/api/reservations",
		method:       "POST",
		body:         insufficienNodesReservation,
		expectedData: nil,
	}
	_ = executeTestReq[db.Reservation](t, app, insufficienNodesTest, tokenStr)

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
