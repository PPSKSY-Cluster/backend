package tests

import (
	"testing"

	"github.com/PPSKSY-Cluster/backend/db"
)

func Test_users(t *testing.T) {
	app, err := setupTestApplication()
	if err != nil {
		t.Error(err)
	}

	user := db.User{Username: "foo", Password: "bar"}
	tokenStr := createUserAndLogin(t, app, user)

	tests := []TestReq{
		{
			description:  "Get all users (expect 200)",
			expectedCode: 200,
			route:        "/api/users/",
			method:       "GET",
			body:         nil,
		},
	}

	for _, test := range tests {
		executeTestReq[[]db.User](t, app, test, tokenStr)
	}
}
