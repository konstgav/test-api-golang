package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseName = "product"
var CollectionName = "productmodel"
var MongoURI = "mongodb://mongo:27017"

func ConnectDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI(MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	collection := client.Database(DatabaseName).Collection(CollectionName)
	if collection == nil {
		log.Fatal("Collection is nil")
	}
	return collection
}
