package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/PPSKSY-Cluster/backend/db"
)

func Test_auth(t *testing.T) {
	fmt.Print("Testing auth...")
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	// create user and get a viable token
	user := db.User{Username: "foo", Password: "bar"}
	createUserAndLogin(t, app, user)

	// use wrong pw
	wrongPWUser := user
	wrongPWUser.Password = "foo"

	badLoginTest := TestReq{
		description:  "Try login with wrong Password (expect 401)",
		expectedCode: 401,
		route:        "/api/login/",
		method:       "POST",
		body:         wrongPWUser,
		expectedData: nil,
	}

	_ = executeTestReq[db.User](t, app, badLoginTest, "")

	// request with wrong jwt
	getAllUsersWrongJWTTest := TestReq{
		description:  "Attempt to get users with an invalid jwt (expect 401)",
		expectedCode: 401,
		route:        "/api/users",
		method:       "GET",
		body:         nil,
		expectedData: nil,
	}

	_ = executeTestReq[db.User](t, app, getAllUsersWrongJWTTest, "not.a.token")

	// request with correct jwt is made in users_test.go
}
