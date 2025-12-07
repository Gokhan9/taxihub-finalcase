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
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "bitaksi-finalcase/driver-service/docs"
)

// @title Bitaksi Driver Service API
// @version 1.0
// @description Driver service içerisinde sürücü kayıt, arama, getirme, silme, durum güncelleme, status durumu, puanlamaya göre hesaplama yapan işlemleri yönetir.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http
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

	log.Println("MongoDB Bağlantısı Başarılı!")

	// VERİTABANI
	db := mongoClient.Database(cfg.MongoDBName)

	// TODO: repo, service ve handler'ı birbirine bağlıyoruz.
	driverRepo := repository.NewDriverRepository(db)
	driverRepo.EnSureIndexes(context.Background()) // Ensure unique indexes (plate)
	driverService := services.NewDriverService(driverRepo)
	driverHandler := handlers.NewDriverHandler(driverService)

	fiberApp := router.SetupRouter(driverHandler)

	// --- GRACEFUL SHUTDOWN KATMANI ----

	c := make(chan os.Signal, 1)                    // Bir kanal oluşturduk ve gelebilecek kapanma sinyallerini alacak.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Yakalamak istediğimiz sinyalleri (CTRL+C - SIGTERM) "Notify" ile tanımladık.

	// * Goroutine
	go func() {
		log.Printf("Driver Service %s Portunda Çalışıyor.", cfg.Port) // Sunucumuzun belirttiğimiz portta çalıştığını doğrulayan str metin.
		if err := fiberApp.Listen("0.0.0.0:" + cfg.Port); err != nil {
			log.Panicf("Fiber sunucusu başlatılamadı: %v", err) // Fiber sunucusunun başlatıldığı sırada karşılaşılacak problemde "panic" oluşturup goroutine sayesine uygulamayı durdurur.
		}
	}()

	<-c // Kapanma sinyali gelene kadar uygulama burada bloklanır, kapanma sinyali geldiği senaryoda alt satırlara girer.
	log.Println("KAPANMA SİNYALİ. Graceful shutdown başlatılıyor.")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second) // kapanma işlemini "5 saniye" olarak sınırladık, 5 saniye içerisinde bitmezse context iptal olacak.
	defer cancelShutdown()                                                                  // işlem biter, context serbest bırakılr.

	// fiber sunucusu kapatılması
	if err := fiberApp.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("sunucu kapatılırken hata oluştu: %v", err)
	}

	// mongodb bağlantısının kapatılması.
	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		log.Printf("mongodb bağlantısı kapatılırken hata oluştu: %v", err)
	}

	log.Println("Uygulama başarılı bir şekilde sonlandırıldı.")
}
