package controllers

import (
	"go-backend/database"
	"go-backend/middleware"
	"go-backend/models"
	"go-backend/requests"
	"go-backend/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateAuth = validator.New()

func Register(c *fiber.Ctx) error {
	var body requests.RegisterRequest

	// JSON parse et
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}

	// Validasyon yap
	if err := validateAuth.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": utils.FormatValidationErrors(err),
		})
	}

	// Şifreyi hashle
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)

	// Yeni kullanıcı oluştur
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
		Role:     body.Role, // veya body'den al
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kullanıcı oluşturulamadı"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Kayıt başarılı",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Şifre yanlış"})
	}

	token, err := middleware.GenerateJWT(user.ID, user.Role, c) // sadece user.ID gönderiliyor
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token oluşturulamadı"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func Logout(c *fiber.Ctx) error {
	jti := c.Locals("jti").(string)

	if err := database.DB.Model(&models.Session{}).Where("token_id = ?", jti).Update("revoked", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Çıkış yapılamadı"})
	}

	return c.JSON(fiber.Map{"message": "Çıkış yapıldı"})
}
