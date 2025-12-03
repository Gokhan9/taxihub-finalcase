package router

import (
	"bitaksi-finalcase/driver-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(DriverHandler *handlers.DriverHandler) *fiber.App {

	app := fiber.New()

	driverGroup := app.Group("/drivers")
	driverGroup.Post("/", DriverHandler.CreateDriver)
	driverGroup.Get("/", DriverHandler.GetAllDrivers)
	driverGroup.Get("/:id", DriverHandler.GetDriverByID)
	driverGroup.Put("/:id", DriverHandler.UpdateDriverByID)
	driverGroup.Delete("/:id", DriverHandler.DeleteDriverByID)

	return app
}
