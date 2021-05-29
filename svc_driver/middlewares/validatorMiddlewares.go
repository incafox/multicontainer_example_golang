package CustomMiddlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/incafox/golang-api/models/modelDriver"
	"strconv"
	"strings"
)

func ValidateLoginFields() func(c *fiber.Ctx) error {
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

func ValidateRegisterFields() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var body ModelDriver.RequestRegister
		err := c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}
		if len(body.Password) < 4 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"Error": "Password should be more than 7 characters."})
		}

		if strings.Contains(body.Email, "@") == false {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"Error": "Introduce valid email with xxx@xxx format."})
		}

		if len(body.Username) < 4 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"Error": "Username should be more than 7 characters"})
		}
		return c.Next()
	}
}

func ValidateQuerySetAvailability() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		available := c.Query("available", "")
		if available == "" || len(available) != 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		val, err := strconv.Atoi(available)
		if err != nil || val > 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		return c.Next()
	}
}

func ValidateQuerySetDistance() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		distance := c.Query("distance", "")
		if distance == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		_, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		return c.Next()
	}
}

func ValidateQueryAllDrivers() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		page := c.Query("page", "")
		limit := c.Query("limit", "")
		if page == "" && limit == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		_, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		_, err = strconv.Atoi(limit)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		return c.Next()
	}
}
func ValidateQueryDriversByDistance() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		distance := c.Query("distance", "2")
		available := c.Query("available", "1")
		if distance == "" && available == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		_, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		availableInt, err := strconv.Atoi(available)
		if err != nil || len(available) != 1 || availableInt > 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
		}
		return c.Next()
	}
}
