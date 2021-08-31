package repository

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseName = "product"
var CollectionName = "productmodel"

func ConnectDB() *mongo.Collection {
	mongo_URI := os.Getenv("MONGO_URI")
	log.Println(os.Environ())
	if mongo_URI == "" {
		panic("Environmental variable MONGO_URI do not set")
	}
	clientOptions := options.Client().ApplyURI(mongo_URI)
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
