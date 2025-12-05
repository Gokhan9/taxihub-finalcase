package handlers

import (
	"bitaksi-finalcase/driver-service/internal/dto"
	"bitaksi-finalcase/driver-service/internal/models"
	"bitaksi-finalcase/driver-service/internal/services"
	"context"
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverHandler struct {
	service *services.DriverService
}

func NewDriverHandler(service *services.DriverService) *DriverHandler {
	return &DriverHandler{service: service}
}

// --------------------------------------  ENDPOINTS --------------------------------------
// CreateDriver godoc
// @Summary Yeni Driver Oluştur
// @Description Sisteme yeni bir taksi sürücüsü kaydeder.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param request body dto.DriverCreateRequest true "Sürücü Bilgileri"
// @Success 201 {object} models.Driver
// @Failure 400 {object} map[string]string "Hatalı İstek"
// @Failure 409 {object} map[string]string "Plaka zaten kayıtlı"
// @Router /drivers [post]
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
		if errors.Is(err, services.ErrPlateExists) {
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
// @Summary ID'ye göre Driver Detayı Getir
// @Description ID'si verilen driver'ın bilgilerini detaylı bir şekilde dönen metod.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID"
// @Success 200 {object} models.Driver
// @Failure 400 {object} map[string]string "Geçersiz ID"
// @Failure 401 {object} map[string]string "Sürücü bulunamadı"
// @Router /drivers/{id} [get]
// GET
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
// @Summary ID'ye göre Driver Günceller
// @Description Driver'ın ad, soyad, araç ve plaka bilgilerini günceller.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID"
// @Param request body dto.DriverUpdateRequest true "Driver güncelleme bilgileri"
// @Success 200 {object} models.Driver
// @Failure 400 {object} map[string]string "Hatalı İstek - StatusBadRequest"
// @Failure 404 {object} map[string]string "Sürücü Bulunamadı - StatusNotFound"
// @Router /driver/{id} [put]
// PUT
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
// @Summary ID'ye göre Driver Siler
// @Description API'den gelen isteğe göre ID ile eşleşen Driver'ı kalıcı olarak silme işlemi yapar.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Sürücü ID"
// @Success 204
// @Failure 400 {object} map[string]string "Geçersiz ID"
// @Failure 404 {object} map[string]string "Sürücü Bulunamadı"
// @Router /drivers/{id} [delete]
// DELETE
func (h *DriverHandler) DeleteDriverByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if err := h.service.DeleteDriverByID(context.Background(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetNearbyTaxisHandler godoc
// @Summary Yakındaki Taksileri Getir
// @Description Belirtilen konum ve yarıçaptaki Müsait (Available) taksileri listeler.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param lat query number true "Enlem (Latitude)"
// @Param lon query number true "Boylam (Longitude)"
// @Param radius query number false "Yarıçap (km) - Default: 6"
// @Success 200 {array} dto.NearbyDriverResponse
// @Router /drivers/nearby [get]
// GET
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

	// dto listesi oluşturduk
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

// GetAllDrivers godoc
// @Summary Tüm Driverları Listele
// @Description Kayıtlı olan tüm driverları sayfalı yani (pagination) şeklinde listeler.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param page query int false "Page Number (Default: 1)"
// @Param pageSize query int false "Page Size (Default: 20)"
// @Success 200 {object} map[string]interface{} "drivers: [], total: int, page: int"
// @Failure 500 {object} map[string]string "Sunucu Hatası - Status Internal Server Error"
// @Router /drivers [get]
// GET
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

// UpdateStatus godoc
// @Summary Driver'ın Status Bilgilerini Günceller
// @Description Driver'ın uygulamayı başlatması, müşteri olması veya olmaması durumunu (available, busy or offline) günceller.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Sürücü ID"
// @Param request body dto.UpdateStatusRequest true "Durum Bilgisi"
// @Success 200 {object} map[string]interface{} "message ve status durumunu döner."
// @Failure 400 {object} map[string]string "geçersiz durum bilgisi"
// @Router /drivers/{id}/status [patch]
// PATCH
func (h *DriverHandler) UpdateStatus(c *fiber.Ctx) error {

	id := c.Params("id")

	var req dto.UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}

	// servis çağrısı
	err := h.service.UpdateDriverStatus(c.Context(), id, req.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sürücünün durumu güncellendi",
		"status":  req.Status,
	})
}

// RateDriver godoc
// @Summary Sürücüyü Puanla
// @Description Bir sürücüye 1-5 arası puan verir.
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Sürücü ID"
// @Param request body dto.RateDriverRequest true "Puan Bilgisi"
// @Success 200 {object} map[string]interface{}
// @Router /drivers/{id}/rate [post]
// POST
func (h *DriverHandler) RateDriver(c *fiber.Ctx) error {

	id := c.Params("id")

	var req dto.RateDriverRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Geçersiz JSON"})
	}

	updatedDriver, err := h.service.RateDriver(c.Context(), id, req.Score)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(bson.M{
		"message":     "Puan Verildi.",
		"rating":      updatedDriver.Rating,
		"ratingCount": updatedDriver.RatingCount,
	})
}

// UpdateLocation godoc
// @Summary Driver'dan Location Bilgisi Alma
// @Description Sürücünün hareket halinde konumunu güncelleyecek bir metod.
// @Tags Driver
// @Accept json
// @Product json
// @Param id path string true "Sürücü ID"
// @Param request body dto.UpdateLocationRequest true "Konum Bilgisi"
// @Success 200 {object} map[string]string "Başarı Mesajı"
// @Failure 400 {object} map[string]string "Hatalı İstek"
// @Failure 404 {object} map[string]string "Sürücü Bulunamadı"
// @Router /drivers/{id}/location [put]
// PUT
func (h *DriverHandler) UpdateLocation(c *fiber.Ctx) error {

	id := c.Params("id")
	var req dto.UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("JSON parsing hatası:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz JSON Formatı.",
		})
	}

	log.Printf("Servis Çağırılıyor. ID: %s, Lat: %f, Lon: %f\n", id, req.Lat, req.Lon)

	err := h.service.UpdateDriverLocation(c.Context(), id, req)
	if err != nil {
		log.Println("Service Hatası:", err)
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"err": "Sürücü Bulunamadı"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Println("UpdateLocation başarılı bir şekilde çalışıyor.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Konum Başarıyla Güncellendi.",
	})
}
