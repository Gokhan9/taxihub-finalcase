package handlers

import (
	"bitaksi-finalcase/driver-service/internal/dto"
	"bitaksi-finalcase/driver-service/internal/models"
	"bitaksi-finalcase/driver-service/internal/services"
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverHandler struct {
	service *services.DriverService
}

func NewDriverHandler(service *services.DriverService) *DriverHandler {
	return &DriverHandler{service: service}
}

// --------------------------------------  ENDPOINTS --------------------------------------

// POST
func (h *DriverHandler) CreateDriver(c *fiber.Ctx) error {

	var req dto.DriverCreateRequest

	// JSON parse
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "JSON geçerli değil",
			"details": err.Error(),
		})
	}

	if req.FirstName == "" || req.LastName == "" || req.Plate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ad, soyad ve plaka alanları doldurulmalıdır.",
		})
	}

	driver := &models.Driver{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Location: models.Location{
			Type:        "Point",
			Coordinates: []float64{req.Lon, req.Lat},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := h.service.CreateDriver(c.Context(), driver)
	if err != nil {

		// plate already exists
		if strings.Contains(err.Error(), "Plaka Kayıtlı") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// invalid plate format or missing fields
		if strings.Contains(err.Error(), "Geçersiz Plaka") || strings.Contains(err.Error(), "Doldurulmadı") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// GetDriverByID
func (h *DriverHandler) GetDriverByID(c *fiber.Ctx) error {

	id := c.Params("id")

	//*servis çağrısı
	driver, err := h.service.GetDriverByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if driver == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "sürücü bulunamadı",
		})
	}

	return c.Status(fiber.StatusOK).JSON(driver)
}

// UpdateDriverByID
func (h *DriverHandler) UpdateDriverByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.DriverUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updatedDriver, err := h.service.UpdateDriverByID(c.Context(), id, &req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "Driver not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(updatedDriver)
}

// DeleteDriverByID
func (h *DriverHandler) DeleteDriverByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if err := h.service.DeleteDriverByID(context.Background(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetNearbyTaxisHandler
func (h *DriverHandler) GetNearbyTaxisHandler(c *fiber.Ctx) error {

	// * lat,lon ve taxitype üzeriden query parametreleri oluşturulur.
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	taxiType := c.Query("taxiType")
	radiusStr := c.Query("radius", "6") // Default radius to 6 km

	// str-float dönüşümleri. Query parametreleri her zaman string gelir ve coğrafi hesaplamalar yaparken sayı gerekiyor.
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid latitude"})
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid longitude"})
	}

	radiusKm, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid radius"})
	}

	drivers, err := h.service.GetNearbyTaxis(c.Context(), lat, lon, radiusKm, taxiType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// DTO Listesi oluştur
	responseList := make([]dto.NearbyDriverResponse, 0)

	// for ile beraber "_, val" ile "driver" değerleri içinde driver'ın dLat ve dLon bilgisini dönüyoruz.
	for _, d := range drivers {
		distanceRound := math.Round(d.Distance*10) / 10 // * istek atınca "distanceKm" değeri "4.97007864" şeklinde dönüyordu, "math" paketi ile 1 ondalıklı basamak şeklinde çıktı almak için kullandım.
		// api'den gelecek isteğe dto ile istediğimiz değerleri dönüyoruz/dışarı açıyoruz.
		item := dto.NearbyDriverResponse{
			FirstName:  d.FirstName,
			LastName:   d.LastName,
			Plate:      d.Plate,
			DistanceKm: distanceRound,
		}

		responseList = append(responseList, item)
	}

	return c.JSON(responseList)
}

// GetAllDrivers
func (h *DriverHandler) GetAllDrivers(c *fiber.Ctx) error {

	// read query params
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	drivers, total, err := h.service.GetAllDrivers(context.Background(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
		"drivers":  drivers,
	})
}
