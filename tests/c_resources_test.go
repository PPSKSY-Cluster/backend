package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/PPSKSY-Cluster/backend/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_cresources(t *testing.T) {
	// setup testing
	fmt.Print("Testing cresources... ")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	// create user
	user := db.User{Username: "foo", Password: "bar"}
	tokenStr, createdUser := createUserAndLogin(t, app, user)

	fooscResource := db.CResource{
		Name:                     "foos cresource",
		Description:              "this is foos first cresource",
		Nodes:                    2,
		Type:                     db.CustomRT,
		Admins:                   []primitive.ObjectID{createdUser.ID},
		BalancingAlg:             db.CustomLB,
		Reservations:             []primitive.ObjectID{},
		HighAvailability:         false,
		HighPerformanceComputing: false,
		OperatingSystem:          db.LinuxOS,
	}

	// execute tests
	createOneTest := TestReq{
		description:  "Create one resource (expect 201)",
		expectedCode: 201,
		route:        "/api/cresources/",
		method:       "POST",
		body:         fooscResource,
		expectedData: fooscResource,
	}
	createdCResource := executeTestReq[db.CResource](t, app, createOneTest, tokenStr)

	getAllTest := TestReq{
		description:  "Get all resources (expect 200)",
		expectedCode: 200,
		route:        "/api/cresources/",
		method:       "GET",
		body:         nil,
		expectedData: []db.CResource{createdCResource},
	}
	executeTestReq[[]db.CResource](t, app, getAllTest, tokenStr)

	getOneTest := TestReq{
		description:  "Get one resource (expect 200)",
		expectedCode: 200,
		route:        "/api/cresources/" + createdCResource.ID.Hex(),
		method:       "GET",
		body:         nil,
		expectedData: createdCResource,
	}
	executeTestReq[db.CResource](t, app, getOneTest, tokenStr)

	editedCResource := createdCResource
	editedCResource.Name = "edited cresource"
	editOneTest := TestReq{
		description:  "Edit one resource (expect 200)",
		expectedCode: 200,
		route:        "/api/cresources/" + createdCResource.ID.Hex(),
		method:       "PUT",
		body:         editedCResource,
		expectedData: editedCResource,
	}
	executeTestReq[db.CResource](t, app, editOneTest, tokenStr)

	deleteOneTest := TestReq{
		description:  "Delete one resource (expect 204)",
		expectedCode: 204,
		route:        "/api/cresources/" + createdCResource.ID.Hex(),
		method:       "DELETE",
		body:         nil,
		expectedData: nil,
	}
	executeTestReq[db.CResource](t, app, deleteOneTest, tokenStr)

}
