package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/*
Request(istek), router'a geldikten sonra güvenlik katmanından geçmesi gerekir. Biz bunu middleware içerisinde kontrol ediyoruz.
Gelen requestlerin authorization(yetkilendirme) yapılarını kontrol eder. "Bearer token" formatında bir format var mı kontrol eder.
Yetkisiz kişilerin servislere erişimini kısıtlaması nedeniyle önemli bir katman.
*/

func JWTMiddleware(secret string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		// secretten gelen değer boşsa, güvenlik burada devredışı kalmış olabilir.
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

		//token := parts[1]
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// İmza metodunun "hmac" olup olmadığını kontrol edeceğimiz aşama
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz veya Süresi Dolmuş Token"})
		}

		// token'ın süresi dolmamış, yukarıda ki kod bloğunu geçerse içerik bilgilerini context'e ekleyebiliriz..
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user", claims)
		}

		return c.Next()
	}
}
