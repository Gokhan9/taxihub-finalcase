package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

/*
Bu fonksiyon MONGO_URI environment variable’ından alınacak URI ile çağrılacak
*/
func Connect(uri string) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	//TODO: "Ping()" İLE DB BAĞLANTISINI KONTROL EDİYORUZ.
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("MongoDB connected successfully")
	Client = client

	// *CREATING 2dsphere INDEX FOR THE 'location' FIELD IN THE 'drivers' COLLECTION..
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: "2dsphere"}},
	}
	_, err = client.Database("bitaksi").Collection("drivers").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("Error creating 2dsphere index: %v\n", err)
		return nil, err
	}
	fmt.Println("2dsphere index created successfully for 'drivers.location'")

	return client, nil
}
