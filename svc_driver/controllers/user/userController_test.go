package userController

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// go test -run -v Test_Handler
func Test_Handler(t *testing.T) {
	app := fiber.New()

	app.Get("/test", SetDriverDistance)

	resp, err := app.Test(httptest.NewRequest("GET", "/test", nil))

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
	if resp.StatusCode == fiber.StatusOK {
		t.Log("Correct status", resp.StatusCode)
	} else {
		t.Error("Wrong status", resp.StatusCode)
	}
}
