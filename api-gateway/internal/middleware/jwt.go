package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

/*
Request(istek), router'a geldikten sonra güvenlik katmanından geçmesi gerekir. Biz bunu middleware içerisinde kontrol ediyoruz.
Gelen requestlerin authorization(yetkilendirme) yapılarını kontrol eder. "Bearer token" formatında bir format var mı kontrol eder.
Yetkisiz kişilerin servislere erişimini kısıtlaması nedeniyle önemli bir katman.
*/

func JWTMiddleware(secret string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Next()
		}

		// Get isteği üzerinden yetkilendirme kontrolü ve sonrasında bir str formatında json err.
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
