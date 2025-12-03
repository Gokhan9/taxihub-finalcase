package main

import (
	"bitaksi-finalcase/api-gateway/internal/config"
	"bitaksi-finalcase/api-gateway/internal/middleware"
	"bitaksi-finalcase/api-gateway/internal/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config'te yükleme hatası: %v", err)
	}

	app := fiber.New(fiber.Config{})

	// -------------------------
	// 3) Global middleware'ler
	//    - Örnek JWT ve RateLimit middleware'leri ekleniyor
	// -------------------------
	// Rate limit (basit IP tabanlı limiter)
	app.Use(middleware.RateLimitMiddleware())

	//*router'u ekler ve kurar
	gwCfg := router.GatewayConfig{
		DriverServiceURL: cfg.DriverServiceURL,
		JWTSecret:        cfg.JWTSecret,
	}

	r := router.SetupRouter(app, gwCfg)

	log.Printf("API Gateway portu %s üzerinde çalışıyor.", cfg.Port)
	if err := r.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("API Gateway başlatılamıyor: %v", err)
	}

}
