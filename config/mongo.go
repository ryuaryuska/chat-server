package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnection() *mongo.Database {
	mongoURL := "mongodb://ryu:password@192.168.1.16:27018"

	clientOptions := options.Client()
	clientOptions.ApplyURI(mongoURL)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db connected")
	dbName := "topindopay"
	return client.Database(dbName)
}
