package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lon float64 `bson:"lon" json:"lon"`
}

type Driver struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Plate     string             `bson:"plate" json:"plate"`
	TaxiType  string             `bson:"taxiType" json:"taxiType"`
	CarBrand  string             `bson:"carBrand" json:"carBrand"`
	CarModel  string             `bson:"carModel" json:"carModel"`
	Location  Location           `bson:"location" json:"location"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
