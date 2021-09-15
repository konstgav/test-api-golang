package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"test-api-golang/repository"
	"time"
)

var TestDatasetFilename = "test_dataset.json"

func CleanAndFillRepository() error {
	collection := repository.ConnectDB()
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
	var docs []interface{}
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
	err = collection.Database().Client().Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}
