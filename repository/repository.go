package repository

import (
	"context"
	"fmt"
	"log"
	"test-api-golang/interfaces"
	"test-api-golang/model"

	"github.com/mitchellh/mapstructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaginationParameters struct {
	Page         int
	LimitPerPage int
}

type CrudListParameters struct {
	interfaces.ListParametersInterface
	*PaginationParameters
}

type CrudRepository struct {
	interfaces.CrudRepositoryInterface
	collection *mongo.Collection
}

func NewCrudRepository(collection *mongo.Collection) *CrudRepository {
	return &CrudRepository{
		collection: collection,
	}
}

func (c CrudRepository) Find(id int) (interfaces.EntityInterface, error) {
	var product model.Product
	filter := bson.M{"_id": id}
	err := c.collection.FindOne(context.TODO(), filter).Decode(&product)
	return product, err
}

func (c CrudRepository) List(parameters interfaces.ListParametersInterface) (interfaces.EntityInterface, error) {
	var products []model.Product
	cur, err := c.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var product model.Product
		err := cur.Decode(&product)
		if err != nil {
			log.Println(err)
		}
		products = append(products, product)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return products, err
}

func (c CrudRepository) Create(item interfaces.EntityInterface) (interfaces.EntityInterface, error) {
	result, err := c.collection.InsertOne(context.TODO(), item)
	return result, err
}

func (c CrudRepository) Update(id int, item interfaces.EntityInterface) (interfaces.EntityInterface, error) {
	data := item.(map[string]interface{})
	var product model.Product
	mapstructure.Decode(data, &product)
	fmt.Println(product.Type, product)
	filter := bson.M{"_id": id}
	update := bson.D{
		{"$set", bson.D{
			{"name", product.Name},
			{"sku", product.Sku},
			{"type", product.Type},
			{"price", product.Price},
		}},
	}
	err := c.collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)
	return item, err
}

func (c CrudRepository) Delete(id int) error {
	filter := bson.M{"_id": id}
	_, err := c.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
