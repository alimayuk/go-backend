package main

import (
	"go-backend/config"
	"go-backend/database"
	"go-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	database.Connect()
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":3000")
}
