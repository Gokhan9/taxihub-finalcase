package router

import (
	_ "bitaksi-finalcase/driver-service/docs"
	"bitaksi-finalcase/driver-service/internal/handlers"

	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRouter(DriverHandler *handlers.DriverHandler) *fiber.App {

	app := fiber.New()

	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	driverGroup := app.Group("/drivers")
	driverGroup.Post("/", DriverHandler.CreateDriver)
	driverGroup.Post("/:id/rate", DriverHandler.RateDriver)
	driverGroup.Get("/", DriverHandler.GetAllDrivers)
	driverGroup.Get("/nearby", DriverHandler.GetNearbyTaxisHandler) // New route for nearby taxis
	driverGroup.Get("/:id", DriverHandler.GetDriverByID)
	driverGroup.Put("/:id", DriverHandler.UpdateDriverByID)
	driverGroup.Put("/:id/location", DriverHandler.UpdateLocation)
	driverGroup.Patch("/:id/status", DriverHandler.UpdateStatus)
	driverGroup.Delete("/:id", DriverHandler.DeleteDriverByID)

	return app
}
