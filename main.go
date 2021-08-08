package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var collection = ConnectDB()

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the test CRUD API!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/products", returnAllProducts).Methods("GET")
	myRouter.HandleFunc("/product", createNewProduct).Methods("POST")
	myRouter.HandleFunc("/product/{id}", returnSingleProduct).Methods("GET")
	myRouter.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")
	myRouter.HandleFunc("/product/{id}", updateProduct).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []Product
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		GetError(err, w)
		return
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
	json.NewEncoder(w).Encode(products)
}

func returnSingleProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	var params = mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		GetError(err, w)
		return
	}
	filter := bson.M{"_id": id}
	err = collection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func createNewProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	result, err := collection.InsertOne(context.TODO(), product)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		GetError(err, w)
		return
	}
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	var params = mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		GetError(err, w)
		return
	}
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&product)
	update := bson.D{
		{"$set", bson.D{
			{"name", product.Name},
			{"sku", product.Sku},
			{"type", product.Type},
			{"price", product.Price},
		}},
	}
	err = collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)
	if err != nil {
		GetError(err, w)
		return
	}
	product.ID = id
	json.NewEncoder(w).Encode(product)
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
