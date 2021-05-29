package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/incafox/golang-api/database"
	"github.com/incafox/golang-api/routes"
	utilsApi "github.com/incafox/golang-api/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(utilsApi.TypeLog().Error, "Error loading .env file")
	}
	database.InitDatabase()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	router.SetupRouterV1(app)
	err = app.Listen(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}
}
