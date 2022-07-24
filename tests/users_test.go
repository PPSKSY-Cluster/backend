package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/PPSKSY-Cluster/backend/db"
)

func Test_users(t *testing.T) {
	// setup testing
	fmt.Print("Testing users... ")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	// create user
	user := db.User{Username: "foo", Password: "bar"}
	tokenStr, createdUser := createUserAndLogin(t, app, user, false)

	editUser := db.User{ID: createdUser.ID, Username: "bar"}

	// execute tests
	getAllTest := TestReq{
		description:  "Get all users (expect 200)",
		expectedCode: 200,
		route:        "/api/users/",
		method:       "GET",
		body:         nil,
		expectedData: []db.User{createdUser},
	}
	executeTestReq[[]db.User](t, app, getAllTest, tokenStr)

	getOneTest := TestReq{
		description:  "Get one user (expect 200)",
		expectedCode: 200,
		route:        "/api/users/" + createdUser.ID.Hex(),
		method:       "GET",
		body:         nil,
		expectedData: createdUser,
	}
	executeTestReq[db.User](t, app, getOneTest, tokenStr)

	editOneTest := TestReq{
		description:  "Edit one user (expect 200)",
		expectedCode: 200,
		route:        "/api/users/" + createdUser.ID.Hex(),
		method:       "PUT",
		body:         editUser,
		expectedData: editUser,
	}
	executeTestReq[db.User](t, app, editOneTest, tokenStr)

	deleteOneTest := TestReq{
		description:  "Delete one user (expect 204)",
		expectedCode: 204,
		route:        "/api/users/" + createdUser.ID.Hex(),
		method:       "DELETE",
		body:         nil,
		expectedData: nil,
	}
	executeTestReq[db.User](t, app, deleteOneTest, tokenStr)
}
