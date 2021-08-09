package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the test CRUD API!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests(router *mux.Router, c CrudControllerInterface) {
	router.HandleFunc("/", homePage)
	router.HandleFunc("/product", c.List).Methods("GET")
	router.HandleFunc("/product", c.Delete).Methods("POST")
	router.HandleFunc("/product/{id}", c.Get).Methods("GET")
	router.HandleFunc("/product/{id}", c.Delete).Methods("DELETE")
	router.HandleFunc("/product/{id}", c.Update).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	crudController := NewCrudController()
	router := mux.NewRouter().StrictSlash(true)
	handleRequests(router, crudController)
}
