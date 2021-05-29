package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/incafox/golang-api/database"
	ModelDriver "github.com/incafox/golang-api/models/modelDriver"
	"github.com/incafox/golang-api/routes"
	"github.com/joho/godotenv"
)

var bearerToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjM1NjYxMTUsInN1YiI6Imx1ZmZ5QGdtYWlsY29tIiwidXNlcm5hbWUiOiJsdWZmeSJ9.LJpWeDpKq0i32Ltj6PDYw6SCf4DiPZrrR8D0Idjq-qs"

var mockDriver string = `{"username":"user1", "email":"user1@testing.com", "password": "awsebs" }`
var mockNonRegisteredDriver string = `{"username":"user1231", "email":"user1231@testing.com", "password": "xyz1231sde" }`

func setupApp() *fiber.App {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.InitDB(true) //testenv = true
	os.Setenv("testenv", "true")
	database.Seed()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	router.SetupRouter(app)
	return app
}

//try to register the same username and email twice
func Test_Register(t *testing.T) {
	app := setupApp()
	var driver ModelDriver.RequestRegister
	json.Unmarshal([]byte(mockDriver), &driver)
	body, _ := json.Marshal(driver)
	_, errDB := database.GetDBConn().Exec(context.Background(), "delete from Drivers where username=$1;", driver.Username)
	if errDB != nil {
		log.Println("Failed deleting test driver data ", errDB)
	}
	req := httptest.NewRequest("POST", "/svc_auth/api/v1/user/register", bytes.NewReader(body))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 1000)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == fiber.StatusCreated {
		t.Log("Correct status", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}

	//second try with same username and email should be rejected
	req = httptest.NewRequest("POST", "/svc_auth/api/v1/user/register", bytes.NewReader(body))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, 1000)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == fiber.StatusUnprocessableEntity {
		t.Log("Correct status, cannot register same username or email", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}
}

/*
  login into registered user in previus test
  1st try : using non registered user credentials
	2nd try : using registered user from previous test
*/
func Test_Login_Handler(t *testing.T) {
	app := setupApp()
	var driver ModelDriver.RequestLogin
	json.Unmarshal([]byte(mockNonRegisteredDriver), &driver)
	body, _ := json.Marshal(driver)
	req := httptest.NewRequest("POST", "/svc_auth/api/v1/user/login", bytes.NewReader(body))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	if resp.StatusCode == fiber.StatusUnauthorized {
		t.Log("Correct status", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}

	json.Unmarshal([]byte(mockDriver), &driver)
	body, _ = json.Marshal(driver)
	req = httptest.NewRequest("POST", "/svc_auth/api/v1/user/login", bytes.NewReader(body))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, 1000)
	if resp.StatusCode == fiber.StatusAccepted {
		t.Log("Correct status", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}
}
