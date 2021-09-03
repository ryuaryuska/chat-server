package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnection() *mongo.Database {
	mongoURL := os.Getenv("MONGO_HOST")

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
	dbName := os.Getenv("MONGO_NAME")
	return client.Database(dbName)
}
