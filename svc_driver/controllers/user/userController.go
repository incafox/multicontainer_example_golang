package userController

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/incafox/golang-api/models/modelDriver"
	utilsApi "github.com/incafox/golang-api/utils"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	//"github.com/joho/godotenv"
)

type requestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DriverAfterLogin(c *fiber.Ctx) error {
	claimToken := c.Locals("token")
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "successfully logged", "token": claimToken})
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func CreateDriver(c *fiber.Ctx) error {
	var body ModelDriver.RequestRegister
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
	}
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env")
	}
	svc_auth_host := os.Getenv("SVC_AUTH_HOST")
	svc_auth_port := os.Getenv("SVC_AUTH_PORT")

	json_data, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post("http://"+svc_auth_host+":"+svc_auth_port+utilsApi.Svc_auth_register_route, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyResponse, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Fatalln(err)
	}
	log.Println(string(bodyResponse), resp.Status)
	var token ModelDriver.RequestToken
	json.Unmarshal([]byte(bodyResponse), &token)
	fmt.Println("token .... ", token.Token)
	c.Locals("token", token.Token)
	if resp.StatusCode != fiber.StatusCreated {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	return c.Status(fiber.StatusCreated).SendString("registered")
}

func Test(c *fiber.Ctx) error {
	return c.SendString(" hola  ")
}

func GetAllDrivers(c *fiber.Ctx) error {
	page := c.Query("page", "")
	limit := c.Query("limit", "")
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	drivers, status := ModelDriver.GetAllDriversFromMongo(pageInt, limitInt)
	log.Println(status)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"page":    page,
		"limit":   limit,
		"drivers": drivers,
	})
}

func SetDriverAvailability(c *fiber.Ctx) error {
	available := c.Query("available", "")
	availableInt, _ := strconv.Atoi(available)
	username := c.Locals("username")
	str := fmt.Sprintf("%v", username)
	ModelDriver.SetDriverAvailability(str, availableInt)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"username": username,
	})
}

func SetDriverDistance(c *fiber.Ctx) error {
	distance := c.Query("distance", "")
	distanceInt, _ := strconv.ParseFloat(distance, 64)
	username := c.Locals("username")
	str := fmt.Sprintf("%v", username)
	ModelDriver.SetDriverDistance(str, distanceInt)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"username": username,
	})
}
func GetDriversByDistance(c *fiber.Ctx) error {
	distance := c.Query("distance", "2")
	available := c.Query("available", "1")
	distanceInt, _ := strconv.ParseFloat(distance, 64)
	availableInt, _ := strconv.Atoi(available)
	drivers, _ := ModelDriver.GetAllDriversByDistanceFromMongo(distanceInt, availableInt)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"distance":  distance,
		"available": available,
		"drivers":   drivers,
	})
}
