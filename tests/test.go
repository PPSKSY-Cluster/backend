package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/PPSKSY-Cluster/backend/api"
	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const TIMEOUT int = -1

func setupTestApplication() (*fiber.App, error) {
	// Server setup : load env file, setup db, setup router
	if err := godotenv.Load("./.env"); err != nil {
		return nil, err
	}
	if err := db.InitDB(); err != nil {
		return nil, err
	}

	app, err := api.InitRouter()
	if err != nil {
		return nil, err
	}

	return app, nil
}

type TestReq struct{
	description  string
	expectedCode int
	route        string
	method       string
	body         interface{}
}

// the provided user will be created, logged in 
// and the jwt token will be returned
func createUserAndLogin(t assert.TestingT, app *fiber.App, user db.User) string {
	createReq := TestReq{
		description:  "Create one user (expect 201)",
		expectedCode: 201,
		route:        "/api/users/",
		method:       "POST",
		body:         user,
	}

	executeTestReq[db.User](t, app, createReq, "")

	loginReq := TestReq{
		description:  "Login the previously created user (expect 200)",
		expectedCode: 200,
		route:        "/api/login/",
		method:       "POST",
		body:         user,
	}

	type TokenRes struct {
		Token string `json: "token"`
	}

	token := executeTestReq[TokenRes](t, app, loginReq, "")
	bearerStr := "Bearer " + token.Token
	
	return bearerStr
}

// executes a test request with the given params and 
// returns the unmarshaled response body for the given type T
// if you don't need authentication leave bearerToken empty
func executeTestReq[T any](t assert.TestingT, app *fiber.App, test TestReq, bearerToken string) T{
	fmt.Printf("\t%s\n", test.description)

	// io reader from body
	bodyBytes, _ := json.Marshal(test.body)
	body := bytes.NewReader(bodyBytes)

	// construct request with body and header
	req := httptest.NewRequest(test.method, test.route, body)
	header := http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {bearerToken},
	}
	req.Header = header

	// execute request
	res, err := app.Test(req, TIMEOUT)
	if err != nil {
		panic(err)
	}
	
	// test result
	assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

	var data T
	buffer, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(buffer, &data)
	return data
}
