package driverController

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
	"time"
)

func DriverLogin(c *fiber.Ctx) error {
	err_env := godotenv.Load()
	if err_env != nil {
		log.Fatal("Error loading .env file")
	}
	secret := os.Getenv("SECRET")
	utilsApi.DefaultEnvField(secret, "secret")
	var body ModelDriver.RequestLogin
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
	}
	driverRequested, exists := ModelDriver.GetDriverByUsername(body.Username, body.Password)
	fmt.Println("ga", driverRequested.Username, driverRequested.Password, exists)
	fmt.Println(driverRequested.Username == "")

	if exists == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "driver don't exists",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = driverRequested.Username
	claims["sub"] = driverRequested.Email
	claims["exp"] = time.Now().Add(time.Hour * 2021).Unix()
	s, err2 := token.SignedString([]byte("secret"))
	log.Println(err2)
	if err2 != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"token": s,
		"user": struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}{
			Id:    driverRequested.Id,
			Email: driverRequested.Email,
		},
	})
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)
	return c.JSON(fiber.Map{"driver": name})
}

func DriverRegister(c *fiber.Ctx) error {
	var body ModelDriver.RequestRegister
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
	}
	errRegister := ModelDriver.CreateDriver(body.Username, body.Email, body.Password)
	if errRegister != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "driver creation failed",
			"details": errRegister.Error(),
		})
	}
	if os.Getenv("testenv") != "true" {
		_ = sendEventToDriver(c)
	}
	return c.Status(fiber.StatusCreated).SendString("registered")
}

func sendEventToDriver(c *fiber.Ctx) error {
	var body ModelDriver.RequestRegister
	err := c.BodyParser(&body)
	if err != nil {
		log.Println(err)
	}
	json_data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.Post("http://svc_drive_go:5000/svc_drive/api/v1/internal/newdriver", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	bodyResponse, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Println(err)
	}
	log.Println(string(bodyResponse), resp.Status)
	return nil
}
