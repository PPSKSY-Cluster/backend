package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"

	"github.com/PPSKSY-Cluster/backend/api"
	"github.com/PPSKSY-Cluster/backend/auth"
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

	if err := auth.InitAuth(); err != nil {
		return nil, err
	}

	app, err := api.InitRouter()
	if err != nil {
		return nil, err
	}

	return app, nil
}

type TestReq struct {
	description  string
	expectedCode int
	route        string
	method       string
	body         interface{}
	expectedData interface{}
}

type Mail struct {
	Id          int       `json:"id"`
	FromAdress  string    `json:"fromAdress"`
	ToAdress    string    `json:"toAdress"`
	Subject     string    `json:"subject"`
	ReceivedOn  time.Time `json:"receivedOn"`
	RawData     string    `json:"rawData"`
	Attachments []string  `json:"attachments"`
}

// the provided user will be created, logged in
// and the jwt token will be returned
func createUserAndLogin(t assert.TestingT, app *fiber.App, user db.User) (string, db.User) {
	expectUser := user
	expectUser.Password = ""
	createReq := TestReq{
		description:  "Create one user (expect 201)",
		expectedCode: 201,
		route:        "/api/users/",
		method:       "POST",
		body:         user,
		expectedData: expectUser,
	}

	createdUser := executeTestReq[db.User](t, app, createReq, "")

	loginReq := TestReq{
		description:  "Login the previously created user (expect 200)",
		expectedCode: 200,
		route:        "/api/login/",
		method:       "POST",
		body:         user,
	}

	type LoginRes struct {
		User  db.User `json:"user"`
		Token string  `json:"token"`
	}

	token := executeTestReq[LoginRes](t, app, loginReq, "")
	compare(t, expectUser, token.User)

	bearerStr := "Bearer " + token.Token
	return bearerStr, createdUser
}

// executes a test request with the given params and
// returns the unmarshaled response body for the given type T
// if you don't need authentication leave bearerToken empty
func executeTestReq[T any](t assert.TestingT, app *fiber.App, test TestReq, bearerToken string) T {
	fmt.Printf("\n%s\n\t", test.description)

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

	// compare each non-zero-valued struct field in
	// the expected struct, to the corresponding field
	// in the returned data (if a slice is expected/returned
	// loop over expected and compare to data)
	if test.expectedData != nil {
		compare(t, test.expectedData, data)
	}

	return data
}

func checkMailTime(t assert.TestingT, send time.Time) {
	time.Sleep(time.Second * 30)

	mail := checkMail()
	diff := mail.ReceivedOn.Local().Sub(send)
	assert.Less(t, diff, time.Second*1)

	deleteMail()
}

func checkMailNotification(t assert.TestingT, to string, notification bool) {
	time.Sleep(time.Second * 15)

	email := checkMail()

	subject := ""
	if notification {
		subject = "Cluster Reservierung wieder mÃ¶glich " //Look into changing this!
	}

	assert.Equal(t, subject, email.Subject)
	deleteMail()
}

func checkMail() Mail {
	var mails []Mail

	r, err := http.Get("http://localhost:5080/api/email")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&mails); err != nil {
		panic(err)
	}

	if len(mails) == 0 {
		return Mail{}
	}

	return mails[0]
}

func deleteMail() {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://localhost:5080/api/email", nil)
	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		panic(err)
	}
}

func compare(t assert.TestingT, expectedData interface{}, data interface{}) {
	expectedV := reflect.ValueOf(expectedData)
	dataV := reflect.ValueOf(data)

	if expectedV.Type().Kind() != dataV.Type().Kind() {
		t.Errorf("Expected %s, got %s", expectedV.Type().Name(), dataV.Type().Name())
	}

	switch expectedV.Type().Kind() {
	case reflect.Slice:
		compareSlice(t, expectedV, dataV)
	case reflect.Struct:
		compareStruct(t, expectedV, dataV)
	}
}

func compareSlice(t assert.TestingT, expected reflect.Value, actual reflect.Value) {
	for i := 0; i < expected.Len(); i++ {
		compareStruct(t, expected.Index(i), actual.Index(i))
	}
}

func compareStruct(t assert.TestingT, expected reflect.Value, actual reflect.Value) {
	for i := 0; i < expected.NumField(); i++ {
		expectF := expected.Field(i)
		if !expectF.IsZero() {
			dataF := actual.Field(i)
			assert.Equal(t, expectF.Interface(), dataF.Interface())
		}
	}
}
