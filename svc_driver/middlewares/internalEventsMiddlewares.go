package CustomMiddlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/incafox/golang-api/models/modelDriver"
)

func InternalEventsCheckDrivers() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var body ModelDriver.RequestLogin // requestLogin
		fmt.Println(body)
		err := c.BodyParser(&body)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}
		if len(body.Password) < 4 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"Error": "Invalid password."})
		}
		if len(body.Username) < 4 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"Error": "Invalid username."})
		}
		return c.Next()
	}
}
