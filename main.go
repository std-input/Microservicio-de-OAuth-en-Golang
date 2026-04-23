package main

import (
	"adbr.xx/auth_microservice/configs"
	"adbr.xx/auth_microservice/database"
	"adbr.xx/auth_microservice/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	configs.SetEnvironmentVariables()
	database.SetupDatabase()

	app := fiber.New()
	app.Use(logger.New())
	routes.Routes(app)

	app.Listen(":" + configs.Get("PORT"))
}
