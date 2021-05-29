package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/incafox/golang-api/database"
	"github.com/incafox/golang-api/routes"
	utilsApi "github.com/incafox/golang-api/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var bearerToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjM1NjYxMTUsInN1YiI6Imx1ZmZ5QGdtYWlsY29tIiwidXNlcm5hbWUiOiJsdWZmeSJ9.LJpWeDpKq0i32Ltj6PDYw6SCf4DiPZrrR8D0Idjq-qs"

func setupApp() *fiber.App {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(utilsApi.TypeLog().Error, "Error loading .env file")
	}
	os.Setenv("testenv", "true")
	database.InitDatabase()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	router.SetupRouterV1(app)
	return app
}

//3 located with less than 2km, 3 located with more than 2km
var mockDrivers = []string{
	`{"Username":"user1", "email":"user1@testing.com", "distance": 1.42}`,
	`{"username":"user2", "email":"user2@testing.com", "distance": 0.22}`,
	`{"username":"user3", "email":"user3@testing.com", "distance": 1.98}`,
	`{"username":"user4", "email":"user4@testing.com", "distance": 2.32}`,
	`{"username":"user5", "email":"user5@testing.com", "distance": 3.92}`,
	`{"username":"user6", "email":"user6@testing.com", "distance": 12.32}`,
}

type Driver struct {
	Email    string  `json:"email" `
	Username string  `json:"username" `
	Distance float64 `json:"distance" `
}

//Add fake drivers
func TestInternalChannel(t *testing.T) {
	app := setupApp()
	database.UseCollection("drivers").DeleteMany(context.Background(), bson.M{})
	for _, driverStr := range mockDrivers {
		var driver Driver
		json.Unmarshal([]byte(driverStr), &driver)
		body, _ := json.Marshal(driver)
		req := httptest.NewRequest("POST", "/svc_drive/api/v1/internal/newdriver", bytes.NewReader(body))
		req.Header.Set("Authorization", bearerToken)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 1000)
		if resp.StatusCode == fiber.StatusOK {
			t.Log("Correct status", resp.StatusCode)
		} else {
			t.Error("Wrong status", resp.StatusCode)
		}
	}
}

//get driver with pagination
//5 firts users
func TestGetDriversWithPagination(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest("GET", "/svc_drive/api/v1/driver/getdrivers?limit=5&page=0", nil)
	req.Header.Set("Authorization", bearerToken)
	resp, _ := app.Test(req, 1000)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	type Response struct {
		Drivers []bson.M `json:"drivers" `
	}
	var birds Response
	json.Unmarshal([]byte(string(body)), &birds)
	if resp.StatusCode == fiber.StatusAccepted && len(birds.Drivers) == 5 {
		t.Log("Correct status", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}
}

//get drivers with less than 2km (see the mocks)
func TestGetDriversByDistance(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest("GET", "/svc_drive/api/v1/driver/getdriversbydistance?distance=2", nil)
	req.Header.Set("Authorization", bearerToken)
	resp, _ := app.Test(req, 1000)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	type Response struct {
		Drivers []bson.M `json:"drivers" `
	}
	var response Response
	json.Unmarshal([]byte(string(body)), &response)
	if resp.StatusCode == fiber.StatusAccepted && len(response.Drivers) == 3 {
		t.Log("Correct status", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}
}
