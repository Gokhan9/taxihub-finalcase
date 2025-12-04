package router

import (
	"bitaksi-finalcase/driver-service/internal/handlers"

	_ "bitaksi-finalcase/driver-service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRouter(DriverHandler *handlers.DriverHandler) *fiber.App {

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	driverGroup := app.Group("/drivers")
	driverGroup.Post("/", DriverHandler.CreateDriver)
	driverGroup.Post("/:id/rate", DriverHandler.RateDriver)
	driverGroup.Get("/", DriverHandler.GetAllDrivers)
	driverGroup.Get("/nearby", DriverHandler.GetNearbyTaxisHandler) // New route for nearby taxis
	driverGroup.Get("/:id", DriverHandler.GetDriverByID)
	driverGroup.Put("/:id", DriverHandler.UpdateDriverByID)
	driverGroup.Patch("/:id/status", DriverHandler.UpdateStatus)
	driverGroup.Delete("/:id", DriverHandler.DeleteDriverByID)

	return app
}
