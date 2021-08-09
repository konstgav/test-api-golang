package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type EntityInterface interface {
}

type ListParametersInterface interface{}

type PaginationParameters struct {
	Page         int
	LimitPerPage int
}

type CrudListParameters struct {
	ListParametersInterface
	*PaginationParameters
}

type CrudRepositoryInterface interface {
	Find(id int) (EntityInterface, error)
	List(parameters ListParametersInterface) (EntityInterface, error)
	Create(item EntityInterface) (EntityInterface, error)
	Update(item EntityInterface) (EntityInterface, error)
	Delete(id int) error
}

type CrudRepository struct {
	CrudRepositoryInterface
}

func NewCrudRepository() *CrudRepository {
	return &CrudRepository{
		CrudRepositoryInterface: nil,
	}
}

var collection = ConnectDB()

func (c CrudRepository) Find(id int) (EntityInterface, error) {
	var product Product
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&product)
	return product, err
}

func (c CrudRepository) List(parameters ListParametersInterface) (EntityInterface, error) {
	var products []Product
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var product Product
		err := cur.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, product)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return products, err
}

func (c CrudRepository) Create(item EntityInterface) (EntityInterface, error) {
	result, err := collection.InsertOne(context.TODO(), item)
	return result, err
}

func (c CrudRepository) Update(item EntityInterface) (EntityInterface, error) {
	var product = item.(Product)
	filter := bson.M{"_id": product.ID}
	update := bson.D{
		{"$set", bson.D{
			{"_id", product.ID},
			{"name", product.Name},
			{"sku", product.Sku},
			{"type", product.Type},
			{"price", product.Price},
		}},
	}
	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)
	return item, err
}

func (c CrudRepository) Delete(id int) error {
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
