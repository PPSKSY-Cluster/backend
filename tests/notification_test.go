package tests

import (
	"fmt"
	"github.com/PPSKSY-Cluster/backend/db"
	"os"
	"testing"
	"time"
)

func Test_notifications(t *testing.T) {
	fmt.Println("Testing reservations... ")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	start := time.Now()
	end := start.Add(time.Hour * 2920)

	user := db.User{
		Username: "foo",
		Password: "bar",
		Mail:     "foo@bar.com",
		Type:     db.AdminUT,
	}
	tokenStr, createdUser := createUserAndLogin(t, app, user, true)

	foosResource := db.CResource{
		Name:            "foos cresource",
		Description:     "this is foos first cresource",
		Nodes:           2,
		OperatingSystem: db.LinuxOS,
		Owner:           createdUser.ID,
	}

	createOnecResourceTest := TestReq{
		description:  "Create one resource (expect 201)",
		expectedCode: 201,
		route:        "/api/cresources/",
		method:       "POST",
		body:         foosResource,
		expectedData: foosResource,
	}
	createdCResource := executeTestReq[db.CResource](t, app, createOnecResourceTest, tokenStr)

	foosReservation := db.Reservation{
		ClusterID: createdCResource.ID,
		Nodes:     2,
		UserID:    createdUser.ID,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		IsExpired: false,
	}

	//Full reservation for no notification to be sent
	createFullReservationTest := TestReq{
		description:  "Create one Reservation (expect 201)",
		expectedCode: 201,
		route:        "/api/reservations/",
		method:       "POST",
		body:         foosReservation,
		expectedData: foosReservation,
	}
	createdReservation := executeTestReq[db.Reservation](t, app, createFullReservationTest, tokenStr)

	foosNotification := db.ReservationNotification{
		ClusterID: createdCResource.ID,
		UserID:    createdUser.ID,
	}

	createNotificationTest := TestReq{
		description:  "Create notification for full resource (expect 201)",
		expectedCode: 201,
		route:        "/api/notifications",
		method:       "POST",
		body:         foosNotification,
		expectedData: foosNotification,
	}
	executeTestReq[db.ReservationNotification](t, app, createNotificationTest, tokenStr)

	//notification set to false as no email is expected to have been sent
	checkMailNotification(t, createdUser.Mail, false)

	updateReservation := createdReservation
	updateReservation.EndTime = time.Now().Unix()
	editReservationTest := TestReq{
		description:  "Edit one reservation (expect 200)",
		expectedCode: 200,
		route:        "/api/reservations/" + createdReservation.ID.Hex(),
		method:       "PUT",
		body:         updateReservation,
		expectedData: updateReservation,
	}
	executeTestReq[db.Reservation](t, app, editReservationTest, tokenStr)
	checkMailNotification(t, createdUser.Mail, true)
}
