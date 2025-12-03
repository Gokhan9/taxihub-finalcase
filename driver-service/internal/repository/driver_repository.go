package repository

import (
	"bitaksi-finalcase/driver-service/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Repository katmanı ile MongoDB koleksiyonuna erişim sağlıyoruz.
CreateDriver, GetDriverByID, UpdateDriver, DeleteDriver ve GetAllDrivers CRUD işlemlerini kapsar.
Tüm tarih alanlarını şuan burada yönetiyoruz. (CreatedAt, UpdatedAt).
*/

type DriverRepository struct {
	collection *mongo.Collection
}

// TODO: Yeni bir Driver Repository oluşturduk.
func NewDriverRepository(db *mongo.Database) *DriverRepository {
	return &DriverRepository{
		collection: db.Collection("drivers"),
	}
}

// TODO: Yeni Driver ekler.
func (r *DriverRepository) CreateDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {

	if driver.ID.IsZero() {
		driver.ID = primitive.NewObjectID()
	}

	now := time.Now()
	if driver.CreatedAt.IsZero() {
		driver.CreatedAt = now
	}
	driver.UpdatedAt = now

	res, err := r.collection.InsertOne(ctx, driver)
	if err != nil {
		return nil, err
	}

	// Insert sonrası ID set
	driver.ID = res.InsertedID.(primitive.ObjectID)

	// Artık driver struct'ı tam ve hazır → DB'den tekrar almaya gerek yok
	return driver, nil
}

// TODO: GetByID
func (r *DriverRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Driver, error) {

	var driver models.Driver

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&driver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &driver, nil
}

// TODO: Driver'ı Günceller.
func (r *DriverRepository) UpdateDriverByID(ctx context.Context, id string, update bson.M) (*models.Driver, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	//* “updatedAt alanını güncelleme”
	update["updatedAt"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	//*MongoDB'de ID'ye göre Update. "$set" ile istediğimiz alan güncellenir. "FindOneAndUpdate", ID bulamazsa ErrNoDocuments döndürüyor
	result := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": update}, opts)

	var updatedDriver models.Driver
	if err := result.Decode(&updatedDriver); err != nil {
		return nil, err
	}

	return &updatedDriver, nil
}

// TODO: Driver siler.
func (r *DriverRepository) DeleteDriverByID(ctx context.Context, id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// TODO: Tüm Driver'ları listeler.
func (r *DriverRepository) GetAllDrivers(ctx context.Context, page, pageSize int) ([]*models.Driver, int64, error) {

	//?MongoDB pagination yaklaşımı:
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	//*toplam count sayısı
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	//*Driverları listele
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var drivers []*models.Driver
	for cursor.Next(ctx) {
		var driver models.Driver
		if err := cursor.Decode(&driver); err != nil {
			return nil, 0, err
		}
		drivers = append(drivers, &driver)
	}
	return drivers, total, nil
}
