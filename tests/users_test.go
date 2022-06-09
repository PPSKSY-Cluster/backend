package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/PPSKSY-Cluster/backend/handlers"
	"github.com/stretchr/testify/assert"
)

func Test_userListHandler(t *testing.T) {
	const route string = "/api/users/"
	const method string = "GET"
	const timeout int = 5

	tests := []struct {
		description  string
		expectedCode int
	}{
		{
			description:  "Get all users (expect 200)",
			expectedCode: 200,
		},
	}

	app, err := SetupTestApplication()
	if err != nil {
		t.Error(err)
	}
	app.Get(route, handlers.UserListHandler())

	// run tests
	for _, test := range tests {
		fmt.Printf("\t%s\n", test.description)
		testReq := httptest.NewRequest(method, route, nil)

		res, err := app.Test(testReq, timeout)
		if err != nil {
			panic(err)
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
	}
}
