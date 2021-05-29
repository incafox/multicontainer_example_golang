package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/incafox/golang-api/controllers/user"
	"github.com/incafox/golang-api/middlewares"
)

func SetupRouter(app *fiber.App) {

	//-----------MAIN ROUTE V1
	api := app.Group("/svc_auth/api/v1")

	//-----------USER ROUTES
	routeUser := api.Group("/user")
	routeUser.Post("/login", CustomMiddlewares.ValidateLoginFields(), driverController.DriverLogin)
	routeUser.Post("/register", CustomMiddlewares.ValidateRegisterFields(), driverController.DriverRegister)
	routeUser.Get("/checkauth", CustomMiddlewares.AuthRequired(), driverController.Restricted)
}

//svc_auth/api/v1/user/checkauth/
