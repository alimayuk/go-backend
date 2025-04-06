package routes

import (
	"go-backend/controllers"

	"go-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Go + Fiber projen çalışıyor!")
	})

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)

	api := app.Group("/api", middleware.Protected())
	todos := api.Group("/todos")
	todos.Get("/", controllers.GetTodos)
	todos.Get("/:id", controllers.GetTodoByID)
	todos.Put("/:id", controllers.UpdateTodo)
	todos.Post("/", controllers.CreateTodo)
	todos.Delete("/:id", controllers.DeleteTodo)
	api.Post("/logout", controllers.Logout)

	admin := api.Group("/admin", middleware.IsAdmin())
	admin.Delete("/dangerous", func(c *fiber.Ctx) error {
		return c.SendString("Sadece admin görebilir 😎")
	})
}
