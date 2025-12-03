package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//go.mongodb.org/mongo-driver/mongo
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

	//TODO: "Ping()" ile DB bağlantısını kontrol ediyoruz.
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("MongoDB connected successfully")
	Client = client
	return client, nil
}
