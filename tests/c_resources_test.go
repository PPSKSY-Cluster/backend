package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/PPSKSY-Cluster/backend/auth"
	"github.com/PPSKSY-Cluster/backend/db"
)

func Test_cresources(t *testing.T) {
	// setup testing
	fmt.Print("Testing cresources... ")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	// create user with admin rights
	user := db.User{Username: "foo", Password: "bar", Type: db.AdminUT}

	userWithHashPw := user
	userWithHashPw.Password, _ = auth.HashPW(user.Password)

	createdUser, _ := db.AddUserWithType(userWithHashPw)
	_, tokenStr, _ := auth.CheckCredentials()(user.Username, user.Password)

	fooscResource := db.CResource{
		Name:            "foos cresource",
		Description:     "this is foos first cresource",
		Nodes:           2,
		OperatingSystem: db.LinuxOS,
		Owner:           createdUser.ID,
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
