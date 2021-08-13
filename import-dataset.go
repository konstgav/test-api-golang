package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"test-api-golang/model"
	"time"
)

var TestDatasetFilename = "test_dataset.json"

func CleanAndFillRepository() error {
	collection := ConnectDB()
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	if err := collection.Drop(ctx); err != nil {
		log.Println("error connection")
		return err
	}
	byteValues, err := ioutil.ReadFile(TestDatasetFilename)
	if err != nil {
		return err
	}
	log.Println("byteValues:", string(byteValues))
	var docs []model.Product
	err = json.Unmarshal(byteValues, &docs)
	if err != nil {
		return err
	}
	for _, doc := range docs {
		_, insertErr := collection.InsertOne(ctx, doc)
		if insertErr != nil {
			return (err)
		}
	}
	collection.Database().Client().Disconnect(ctx)
	return nil
}
