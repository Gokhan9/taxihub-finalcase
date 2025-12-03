package main

import (
	"bitaksi-finalcase/driver-service/internal/config"
	"bitaksi-finalcase/driver-service/internal/handlers"
	"bitaksi-finalcase/driver-service/internal/repository"
	"bitaksi-finalcase/driver-service/internal/router"
	"bitaksi-finalcase/driver-service/internal/services"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	//* Config entegrasyonu
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config dosyası yüklenemedi: %v", err)
	}

	//*MongoDB Bağlantısı
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = cfg.MongoURI //? fallback
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB bağlantısı kurulamadı: %v", err)
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping başarısız: %v", err)
	}

	log.Println("MongoDB bağlantısı başarılı")

	// VERİTABANI
	db := mongoClient.Database(cfg.MongoDBName)

	// TODO: repo, service ve handler'ı birbirine bağlıyoruz.
	driverRepo := repository.NewDriverRepository(db)
	driverService := services.NewDriverService(driverRepo)
	driverHandler := handlers.NewDriverHandler(driverService)

	fiberApp := router.SetupRouter(driverHandler)

	log.Printf("Driver Service %s portunda çalışıyor.", cfg.Port)
	if err := fiberApp.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Fiber başlatılamadı: %v", err)
	}
}
