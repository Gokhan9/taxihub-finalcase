package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}

// Sabit tanımladık. Sebebi ise kodun başka yerinde bu kullanımların yanlış olmaması.
const (
	StatusAvailable = "Available" // Müsait, yolcu arıyor&bekliyor
	StatusBusy      = "Busy"      // Dolu, yolcu taşıyor
	StatusOffline   = "Offline"   // Aktif değil, uygulama kapalı
)

type Driver struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"` //Belgeyi mongo'ya "primitive _id" olarak kaydeder. - id(string) - objectID tipine çevirir.
	FirstName   string             `bson:"firstName" json:"firstName"`
	LastName    string             `bson:"lastName" json:"lastName"`
	Plate       string             `bson:"plate" json:"plate"`
	TaxiType    string             `bson:"taxiType" json:"taxiType"`
	CarBrand    string             `bson:"carBrand" json:"carBrand"`
	CarModel    string             `bson:"carModel" json:"carModel"`
	Location    Location           `bson:"location" json:"location"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	Distance    float64            `json:"distance,omitempty"`
	Status      string             `bson:"status" json:"status"`           // Sürücünün o anki durumu. json:"status" -> API'nin yanıtını döner, bson:"status" -> veritabanında bu ad ile saklanır.
	Rating      float64            `bson:"rating" json:"rating"`           // ortalama puan = 4.8 gibi. float döndüğünden dolayı.
	RatingCount int                `bson:"ratingCount" json:"ratingCount"` // kaç kişinin oy kullandığını tutar.
}
