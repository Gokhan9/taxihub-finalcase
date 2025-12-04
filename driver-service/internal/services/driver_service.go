package services

import (
	"context"
	"errors"
	"log"
	"regexp"

	"bitaksi-finalcase/driver-service/internal/dto"
	"bitaksi-finalcase/driver-service/internal/models"
	"bitaksi-finalcase/driver-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Service layer, repository fonksiyonlarını çağırarak basit bir "business validation" işlemi yapıyoruz.
CRUD işlemlerinde ekstra kontrol ekledik (örn. boş alan kontrolü, update validation).
Handler’lar sadece HTTP request ve response ile ilgilenecek..
*/

type DriverService struct {
	repo *repository.DriverRepository
}

func NewDriverService(repo *repository.DriverRepository) *DriverService {
	return &DriverService{
		repo: repo,
	}
}

// CreateDriver adds a new driver after basic validation.
var ErrPlateExists = errors.New("Plaka Kayıtlı")
var plateRegex = regexp.MustCompile(`^[0-9]{2}[A-Z]{1,3}[0-9]{2,4}$`)

func (s *DriverService) CreateDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {

	if driver.FirstName == "" || driver.LastName == "" || driver.Plate == "" {
		log.Println("validation failed: missing fields.")
		return nil, errors.New("Ad, Soyad  ve plaka alanlarının doldurulması zorunludur.")
	}

	// Plate format check
	if !plateRegex.MatchString(driver.Plate) {
		log.Println("validation failed: invalid plate formats.")
		return nil, errors.New("Geçersiz Plaka Formatı. Örn:34ABC123")
	}

	// duplicate check
	exists, err := s.repo.IsPlateExists(ctx, driver.Plate)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPlateExists
	}

	return s.repo.CreateDriver(ctx, driver)
}

// GetDriverByID retrieves a driver by ID, handling string-to-ObjectID conversion.
func (s *DriverService) GetDriverByID(ctx context.Context, idStr string) (*models.Driver, error) {

	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, errors.New("Geçersiz ID")
	}

	driver, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if driver == nil {
		return nil, nil
	}

	return driver, nil
}

// UpdateDriverByID updates non-nil fields of a driver.
func (s *DriverService) UpdateDriverByID(ctx context.Context, id string, req *dto.DriverUpdateRequest) (*models.Driver, error) {
	// Güncellenecek alanları tutmak için boş bir map oluşturuyoruz.
	// bson.M MongoDB için bir map tipidir ve alan-değer çiftlerini temsil eder.
	update := bson.M{}

	// Eğer FirstName boş değilse, update map'ine ekle
	if req.FirstName != nil {
		update["firstName"] = *req.FirstName
	}
	if req.LastName != nil {
		update["lastName"] = *req.LastName
	}
	if req.Plate != nil {
		update["plate"] = *req.Plate
	}

	if req.Location != nil {
		update["location"] = bson.M{
			"type":        "Point",
			"coordinates": []float64{req.Location.Lon, req.Location.Lat},
		}
	}

	updatedDriver, err := s.repo.UpdateDriverByID(ctx, id, update)
	if err != nil {
		return nil, err
	}

	return updatedDriver, nil
}

// DeleteDriverByID removes a driver by ID.
func (s *DriverService) DeleteDriverByID(ctx context.Context, id string) error {
	return s.repo.DeleteDriverByID(ctx, id)
}

// GetNearbyTaxis retrieves a list of nearby drivers based on location and taxi type.
func (s *DriverService) GetNearbyTaxis(ctx context.Context, lat, lon, radiusKm float64, taxiType string) ([]*models.Driver, error) {
	drivers, err := s.repo.FindNearbyDrivers(ctx, lat, lon, radiusKm, taxiType)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}

// GetAllDrivers retrieves all drivers with pagination.
func (s *DriverService) GetAllDrivers(ctx context.Context, page, pageSize int) ([]*models.Driver, int64, error) {
	return s.repo.GetAllDrivers(ctx, page, pageSize)
}

func (s *DriverService) UpdateDriverStatus(ctx context.Context, id string, status string) error {

	// business rule ekleyerek status'un geçerli olup olmadığına göre kontrol yapıyoruz.
	if status != models.StatusAvailable && status != models.StatusBusy && status != models.StatusOffline {
		return errors.New("Geçersiz durum bilgisi. (Available, Busy ya da Offline) olmalıdır.")
	}

	return s.repo.UpdateDriverStatus(ctx, id, status)
}
