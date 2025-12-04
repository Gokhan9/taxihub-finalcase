package router

import (
	"bitaksi-finalcase/driver-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(DriverHandler *handlers.DriverHandler) *fiber.App {

	app := fiber.New()

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
