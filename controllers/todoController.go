package controllers

import (
	"fmt"
	"go-backend/database"
	"go-backend/models"
	"go-backend/requests"
	"go-backend/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateTodo = validator.New()

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
	var body requests.TodoRequest

	// Form-data ile gönderilen text alanları BodyParser ile parse edebiliriz
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}

	if err := validateTodo.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": utils.FormatValidationErrors(err),
		})
	}

	var imagePath string
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		imagePath = "uploads/" + filename
		if err := c.SaveFile(file, imagePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Görsel kaydedilemedi"})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Görsel dosyası zorunludur",
		})
	}

	todo := models.Todo{
		Title:     body.Title,
		IsDone:    body.IsDone,
		ImagePath: imagePath,
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Todo kaydedilemedi"})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID boş olamaz"})
	}

	var todo models.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo bulunamadı"})
	}

	// Form'dan gelen yeni title ve is_done
	title := c.FormValue("title")
	isDone := c.FormValue("is_done") == "true"

	// Yeni dosya yüklenmişse:
	imagePath := todo.ImagePath // varsayılan eski dosya

	if file, err := c.FormFile("image"); err == nil && file != nil {
		// ✅ Yeni resmi yükle
		newPath, err := utils.UploadImage(c, "image")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Görsel yüklenemedi"})
		}

		// ✅ Eski resmi sil
		utils.DeleteImage(todo.ImagePath)
		imagePath = newPath
	}

	todo.Title = title
	todo.IsDone = isDone
	todo.ImagePath = imagePath

	if err := database.DB.Save(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Todo güncellenemedi"})
	}

	return c.Status(200).JSON(todo)
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID boş olamaz",
		})
	}

	var todo models.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo bulunamadı",
		})
	}

	// Görsel varsa sil
	if todo.ImagePath != "" {
		if err := utils.DeleteImage(todo.ImagePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Görsel silinemedi",
			})
		}
	}

	if err := database.DB.Delete(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Todo silinemedi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo başarıyla silindi",
	})
}
