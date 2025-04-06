package controllers

import (
	"go-backend/database"
	"go-backend/middleware"
	"go-backend/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // opsiyonel, default "user"
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	user := models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(hash),
		Role:     data.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kullanıcı oluşturulamadı"})
	}

	return c.JSON(user)
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
