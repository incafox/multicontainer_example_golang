package eventController

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/incafox/golang-api/database"
	"github.com/incafox/golang-api/utils"
	"log"
)

type Driver struct {
	Email    string  `json:"email" `
	Username string  `json:"username" `
	Distance float64 `json:"distance" `
}

func AddNewDriver(c *fiber.Ctx) error {
	var body Driver //map[string]interface{}
	errParse := c.BodyParser(&body)
	if errParse != nil {
		log.Println(utilsApi.TypeLog().Error, "parsing error", errParse, body)
	}
	driverColl := database.UseCollection("drivers")
	_, err := driverColl.InsertOne(context.TODO(), body)
	if err != nil {
		log.Println(utilsApi.TypeLog().DatabaseError, "eventController", err)
	}
	return c.SendString("new driver added")
}
