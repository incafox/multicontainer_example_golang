package CustomMiddlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func AuthRequired() func(c *fiber.Ctx) error {
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
