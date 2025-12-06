package router

import (
	"bitaksi-finalcase/api-gateway/internal/handlers"
	"bitaksi-finalcase/api-gateway/internal/middleware"
	"bitaksi-finalcase/api-gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type GatewayConfig struct {
	DriverServiceURL string
	JWTSecret        string
}

/*
Fiber App'i işler ve route'ları kurar.
Konfigürasyon yüklendikten sonra ki gelen isteklerin(HTTP) yol haritasını belirleyen "trafik polisidir.." Kısaca istekleri dağıtır.
*/
func SetupRouter(app *fiber.App, cfg GatewayConfig) *fiber.App {

	// "/health" endpointi(adres) kullanarak servis&servislerin ayakta olup olmadığını bildirir.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	driverGroup := app.Group("/drivers") // Driver'dan gelen route'leri bir grup altında topluyoruz.

	// ProxyHandler önceki router içerisine "JWTSecret" ekledik. /drivers apilerine gelen istekler ilk önce "JWTMiddleware" içerisinde kontrol edilip "token" geçerliyse proxy'ye devam edecek.
	driverGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	// Proxy Handler: Tüm "/drivers" isteklerini DriverService'e ileten bir proxy.
	// Koruma sağlar(method&path), driver-service'e yönlendirme yapıyor.
	proxyHandler := handlers.NewProxyHandler(cfg.DriverServiceURL)

	// Endpointleri, proxy handler içerisine mapledik.
	// "/*" drivers altında yer alan tüm HTTP isteklerini yakalar ve bunları API gateway içerisinde yer alan handlers yapısında ki "Proxy" yapısına iletir. Bu sayede istemci arkada kaç tane servis çalışırsa çalışsın tek bir servis ile muhattap olur.
	driverGroup.All("/*", proxyHandler.Handle)
	driverGroup.All("", proxyHandler.Handle)

	app.All("/_forward/*", utils.ForwardRawRequest(cfg.DriverServiceURL))
	return app
}
