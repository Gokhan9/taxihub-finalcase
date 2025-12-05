package middleware

import (
	"net"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter) //IP - Rate Limiter haritalaması yapılır. Memory üzerinde her IP adresine ait bir rate limiter tutuluyor.
var mu sync.Mutex                             // Harita erişim çakışmalarının önüne geçmek ve önlemek için mutex.

/*
Genel kullanım amacı IP isteklerinden gelen request(istek) sayısını sınırlar.. (örnek: saniyede 5 istek etc..)
Middleware kapsamında önemi ise servisler üzerinde aşırı yükleme ve kötü niyetli yazılımlara karşı tam koruma sağlanan ara yazılım katmanı olarakta bilinir.
*/
func getVisitor(ip string) *rate.Limiter {

	mu.Lock() // Harita üzerinde concurrency(eşzamanlı) erişimi koyulur.
	defer mu.Unlock()

	lim, exists := visitors[ip]
	if !exists {
		lim = rate.NewLimiter(5, 10) // Harita üzerinde her ip için yeni bir rate limiter oluşturur
		visitors[ip] = lim
	}

	return lim
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			mu.Lock() // 5 istek sonrası goroutine devreye girer ve locklama yapar.
			for ip, l := range visitors {
				_ = l
				_ = ip
			}
			mu.Unlock()
		}
	}()
}

// Fiber Middleware
func RateLimitMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		ip := c.IP() // Request'in gerçek IP adresini almak.
		if ip == "" {
			host, _, err := net.SplitHostPort(c.Context().RemoteAddr().String())
			if err == nil {
				ip = host
			}
		}

		lim := getVisitor(ip)
		// Rate limit kontrolü yapılır, Limit doluysa(5), 429 döner.
		if !lim.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limit'i aştınız"})
		}

		return c.Next()

	}

}
