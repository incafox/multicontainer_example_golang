package router

import (
	"github.com/gofiber/fiber/v2"
	eventController "github.com/incafox/golang-api/controllers/events"
	userController "github.com/incafox/golang-api/controllers/user"
	CustomMiddlewares "github.com/incafox/golang-api/middlewares"
)

func SetupRouterV1(app *fiber.App) {
	/*        Main Routes        */
	api := app.Group("/svc_drive/api/v1")
	routeUser := api.Group("/driver")
	routeUser.Post("/login", CustomMiddlewares.ValidateLoginFields(), CustomMiddlewares.LoginMiddleware(), userController.DriverAfterLogin)
	routeUser.Post("/register", CustomMiddlewares.ValidateRegisterFields(), userController.CreateDriver)
	routeUser.Post("/setavailability", CustomMiddlewares.CheckTokenStatusMiddleware(), CustomMiddlewares.ValidateQuerySetAvailability(), userController.SetDriverAvailability)
	routeUser.Post("/setdistance", CustomMiddlewares.CheckTokenStatusMiddleware(), CustomMiddlewares.ValidateQuerySetDistance(), userController.SetDriverDistance)
	routeUser.Get("/getdrivers", CustomMiddlewares.CheckTokenStatusMiddleware(), CustomMiddlewares.ValidateQueryAllDrivers(), userController.GetAllDrivers)
	routeUser.Get("/getdriversbydistance", CustomMiddlewares.CheckTokenStatusMiddleware(), CustomMiddlewares.ValidateQueryDriversByDistance(), userController.GetDriversByDistance)
	//you should replace this with some internal queue service
	routeInternal := api.Group("/internal")
	routeInternal.Post("/newdriver", eventController.AddNewDriver)
}
