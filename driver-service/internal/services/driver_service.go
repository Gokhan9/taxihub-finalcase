package services

import (
	"context"
	"errors"

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

// ! Basit bir Validation işlemi driver ekliyoruz. Eklemeler yapılabilir.
func (s *DriverService) CreateDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {

	if driver.FirstName == "" || driver.LastName == "" || driver.Plate == "" {
		return nil, errors.New("Ad, soyad  ve plaka alanlarının doldurulması zorunludur.")
	}
	return s.repo.CreateDriver(ctx, driver)
}

// TODO: ID'ye göre Driver getirir. Eklemeler yapılabilir. "string" → "ObjectID" dönüşümü yalnızca service katmanında yapılmalı
func (s *DriverService) GetDriverByID(ctx context.Context, idStr string) (*models.Driver, error) {

	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, errors.New("geçersiz id")
	}

	driver, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if driver == nil {
		return nil, nil
	}

	return driver, nil

	//return s.repo.GetDriverByID(ctx, id)
}

// TODO: Driver'da ki dolu alanları güncelleme
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

// TODO: Driver'ları siler.
func (s *DriverService) DeleteDriverByID(ctx context.Context, id string) error {
	return s.repo.DeleteDriverByID(ctx, id)
}

// TODO: GetNearbyTaxis retrieves a list of nearby drivers based on location and taxi type.
func (s *DriverService) GetNearbyTaxis(ctx context.Context, lat, lon, radiusKm float64, taxiType string) ([]*models.Driver, error) {
	drivers, err := s.repo.FindNearbyDrivers(ctx, lat, lon, radiusKm, taxiType)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}

// TODO: Tüm driverları listeliyoruz.
func (s *DriverService) GetAllDrivers(ctx context.Context, page, pageSize int) ([]*models.Driver, int64, error) {
	return s.repo.GetAllDrivers(ctx, page, pageSize)
}
