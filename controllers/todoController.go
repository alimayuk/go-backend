package controllers

import (
	"go-backend/database"
	"go-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetTodos(c *fiber.Ctx) error {
	todos := []models.Todo{}
	if err := database.DB.Find(&todos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Todo'lar alınamadı",
		})
	}
	return c.Status(fiber.StatusOK).JSON(todos)
}

func GetTodoByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID boş olamaz",
		})
	}
	todo := models.Todo{}
	if err := database.DB.First(&todo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo bulunamadı",
		})
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

func CreateTodo(c *fiber.Ctx) error {
	todo := new(models.Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Todo oluşturulamadı",
		})
	}
	if todo.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Başlık boş olamaz",
		})
	}
	if err := database.DB.Create(todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Todo oluşturulamadı",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(todo)

}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID boş olamaz",
		})
	}
	todo := models.Todo{}
	if err := database.DB.First(&todo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo bulunamadı",
		})
	}
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Todo güncellenemedi",
		})
	}
	if err := database.DB.Save(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Todo güncellenemedi",
		})
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID boş olamaz",
		})
	}
	todo := models.Todo{}
	if err := database.DB.Where("id = ?", id).Delete(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Todo silinemedi",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo başarıyla silindi",
	})
}
