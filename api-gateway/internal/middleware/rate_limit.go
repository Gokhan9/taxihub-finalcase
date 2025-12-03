package middleware

import (
	"net"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {

	mu.Lock()
	defer mu.Unlock()

	lim, exists := visitors[ip]
	if !exists {
		lim = rate.NewLimiter(5, 10)
		visitors[ip] = lim
	}

	return lim
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			mu.Lock()
			for ip, l := range visitors {
				_ = l
				_ = ip
			}
			mu.Unlock()
		}
	}()
}

func RateLimitMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if ip == "" {
			host, _, err := net.SplitHostPort(c.Context().RemoteAddr().String())
			if err == nil {
				ip = host
			}
		}

		lim := getVisitor(ip)
		if !lim.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limit'i aştınız"})
		}

		return c.Next()

	}

}
