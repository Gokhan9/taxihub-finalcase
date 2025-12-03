package router

import (
	"bitaksi-finalcase/api-gateway/internal/handlers"
	"bitaksi-finalcase/api-gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type GatewayConfig struct {
	DriverServiceURL string
	JWTSecret        string
}

/*
TODO: Fiber App'i işler ve route'ları kurar.
*/
func SetupRouter(app *fiber.App, cfg GatewayConfig) *fiber.App {

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	driverGroup := app.Group("/drivers") // Driver'dan gelen route'leri bir grup altında topluyoruz.

	/*
		TODO: Proxy Handler: Tüm "/drivers" isteklerini DriverService'e ileten bir proxy.
		TODO: Koruma sağlar(method&path), driver-service'e yönlendirme yapıyor.
	*/
	proxyHandler := handlers.NewProxyHandler(cfg.DriverServiceURL)

	// TODO: Endpointleri, proxy handler içerisine mapledik.
	driverGroup.All("/*", proxyHandler.Handle)
	driverGroup.All("", proxyHandler.Handle)

	app.All("/_forward/*", utils.ForwardRawRequest(cfg.DriverServiceURL))
	return app
}
