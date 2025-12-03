package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(secret string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Next()
		}

		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header eksik."})
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header formatı: bearer <token>"})
		}

		token := parts[1]

		//TODO: Token doğrulaması yapılır.
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token geçersiz"})
		}

		return c.Next()
	}
}
