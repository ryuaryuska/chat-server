package main

import (
	"chat-app/model"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoConn *mongo.Database

var collectionName = "messages"
var ctx = context.Background()

func getMessageCode(code string) string {
	messageCollection := mongoConn.Collection("messageCode")
	var msg model.CodeMsg

	messageCollection.FindOne(ctx, bson.M{"code": code}).Decode(&msg)
	return msg.Message
}

func getPreviousMsg(room string) []Message {
	messageCollection := mongoConn.Collection(collectionName)

	csr, err := messageCollection.Find(ctx, bson.M{"target": room})
	if err != nil {
		log.Fatal(err)
	}

	messages := make([]Message, 0)

	for csr.Next(ctx) {
		var message Message
		err := csr.Decode(&message)
		if err != nil {
			log.Fatal(err)
		}

		messages = append(messages, message)
	}

	return messages
}

func countMsg() int64 {
	messageCollection := mongoConn.Collection(collectionName)

	total, _ := messageCollection.CountDocuments(ctx, bson.M{})

	return total
}

func insertToDb(msg Message) {
	messageCollection := mongoConn.Collection(collectionName)

	_, err := messageCollection.InsertOne(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}
}
