package tests

import (
	"os"
	"testing"

	"github.com/PPSKSY-Cluster/backend/db"
)

func Test_users(t *testing.T) {
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}
	defer db.DropDB(os.Getenv("DB_NAME"))

	user := db.User{Username: "foo", Password: "bar"}
	tokenStr, createdUser := createUserAndLogin(t, app, user)

	editUser := db.User{Username: "bar"}

	tests := []TestReq{
		{
			description:  "Get all users (expect 200)",
			expectedCode: 200,
			route:        "/api/users/",
			method:       "GET",
			body:         nil,
		},
		{
			description:  "Get one user (expect 200)",
			expectedCode: 200,
			route:        "/api/users/" + createdUser.ID.Hex(),
			method:       "GET",
			body:         nil,
		},
		{
			description:  "Edit one user (expect 200)",
			expectedCode: 200,
			route:        "/api/users/" + createdUser.ID.Hex(),
			method:       "PUT",
			body:         editUser,
		},
		{
			description:  "Delete one user (expect 204)",
			expectedCode: 204,
			route:        "/api/users/" + createdUser.ID.Hex(),
			method:       "DELETE",
			body:         nil,
		},
	}

	for _, test := range tests {
		executeTestReq[[]db.User](t, app, test, tokenStr)
	}
}
