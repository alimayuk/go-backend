package middleware

import (
	"fmt"
	"go-backend/database"
	"go-backend/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("gizli_anahtar")

func GenerateJWT(userID uint, role string, c *fiber.Ctx) (string, error) {
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"jti":     jti,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// 🧠 GÖZLE GÖRÜNÜR LOG
	fmt.Println("⏺ Session Ekleniyor:", jti)

	session := models.Session{
		UserID:    userID,
		TokenID:   jti,
		UserAgent: c.Get("User-Agent"),
		IP:        c.IP(),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		fmt.Println("❌ Session eklenemedi:", err)
	} else {
		fmt.Println("✅ Session başarıyla eklendi")
	}

	return signedToken, nil
}

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header eksik"})
		}

		fields := strings.Fields(tokenStr)
		if len(fields) == 2 && strings.ToLower(fields[0]) == "bearer" {
			tokenStr = fields[1]
		} else if len(fields) == 1 {
			// sadece token geldiyse → sorun yok, direkt kullan
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization formatı hatalı"})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token geçersiz"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		jti := claims["jti"].(string)

		// Session kontrolü
		var session models.Session
		if err := database.DB.Where("token_id = ? AND user_id = ?", jti, userID).First(&session).Error; err != nil || session.Revoked {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Oturum geçersiz"})
		}

		// ID'leri route'larda kullanmak için kaydet
		c.Locals("user_id", userID)
		c.Locals("jti", jti)
		role := claims["role"].(string)
		c.Locals("role", role) // 👈 Burada kaydetmezsen IsAdmin bulamaz
		return c.Next()
	}
}

func IsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin erişimi gerekiyor"})
		}
		return c.Next()
	}
}
