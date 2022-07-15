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
	tokenStr, _ := createUserAndLogin(t, app, user)

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

	// check wrong jwt token
	checkWrongJWTTest := TestReq{
		description:  "Use token check route with wrong token (expect 401)",
		expectedCode: 401,
		route:        "/api/token-check",
		method:       "POST",
		body:         nil,
		expectedData: nil,
	}

	_ = executeTestReq[db.User](t, app, checkWrongJWTTest, "not.a.token")

	// check correct jwt
	checkCorrectJWTTest := TestReq{
		description:  "Use token check route with correct token (expect 200)",
		expectedCode: 200,
		route:        "/api/token-check",
		method:       "POST",
		body:         nil,
		expectedData: nil,
	}

	_ = executeTestReq[db.User](t, app, checkCorrectJWTTest, tokenStr)

	// request with correct jwt is made in users_test.go
}
