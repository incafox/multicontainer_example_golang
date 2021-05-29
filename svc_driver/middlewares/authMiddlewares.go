package CustomMiddlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/incafox/golang-api/models/modelDriver"
	"github.com/incafox/golang-api/utils"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func authDriverMiddleware() func(c *fiber.Ctx) error {
	def := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unautorized",
			})
		},
		SigningKey: []byte("secret"),
	})
	return def
}

func LoginMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env")
		}
		svc_auth_host := os.Getenv("SVC_AUTH_HOST")
		svc_auth_port := os.Getenv("SVC_AUTH_PORT")
		var body ModelDriver.RequestLogin
		errParse := c.BodyParser(&body)
		if errParse != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}
		json_data, err := json.Marshal(body)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post("http://"+svc_auth_host+":"+svc_auth_port+utilsApi.Svc_auth_login_route, "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyResponse, errRead := ioutil.ReadAll(resp.Body)
		if errRead != nil {
			log.Fatalln(err)
		}
		log.Println(string(bodyResponse))
		var token ModelDriver.RequestToken
		json.Unmarshal([]byte(bodyResponse), &token)
		//fmt.Println("token .... ", token.Token)
		c.Locals("token", token.Token)
		if resp.StatusCode != fiber.StatusAccepted {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Next()
	}
}

func RegisterMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env")
		}
		svc_auth_host := os.Getenv("SVC_AUTH_HOST")
		svc_auth_port := os.Getenv("SVC_AUTH_PORT")
		var body ModelDriver.RequestLogin
		errParse := c.BodyParser(&body)
		if errParse != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}
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
		log.Println(string(bodyResponse))
		var token ModelDriver.RequestToken
		json.Unmarshal([]byte(bodyResponse), &token)
		//fmt.Println("token .... ", token.Token)
		c.Locals("token", token.Token)
		if resp.StatusCode == fiber.StatusUnauthorized {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}

func CheckTokenStatusMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env")
		}
		svc_auth_host := os.Getenv("SVC_AUTH_HOST")
		svc_auth_port := os.Getenv("SVC_AUTH_PORT")
		if os.Getenv("testenv") == "true" {
			svc_auth_host = "localhost"
		}
		dataAuth := c.Request().Header.Peek("Authorization")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://"+svc_auth_host+":"+svc_auth_port+utilsApi.Svc_auth_check_token_route, nil)
		req.Header.Set("Authorization", string(dataAuth))
		log.Println(utilsApi.TypeLog().Svc_Auth_Info, string(dataAuth))
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		bodyResponse, errRead := ioutil.ReadAll(res.Body)
		if errRead != nil {
			log.Fatalln(err)
		}
		log.Println(utilsApi.TypeLog().Svc_Auth_Info, string(bodyResponse), res.StatusCode)
		var result map[string]interface{}
		json.Unmarshal([]byte(string(bodyResponse)), &result)
		c.Locals("username", result["driver"])
		if res.StatusCode != fiber.StatusOK {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token")
		}
		return c.Next()
	}
}
