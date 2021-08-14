package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"test-api-golang/interfaces"

	"github.com/gorilla/mux"
)

type GorillaRouterInterface interface {
	InitRouter() *mux.Router
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the test CRUD API!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests(router *mux.Router, c interfaces.CrudControllerInterface) {
	router.HandleFunc("/", homePage)
	router.HandleFunc("/product", c.List).Methods("GET")
	router.HandleFunc("/product", c.Create).Methods("POST")
	router.HandleFunc("/product/{id}", c.Get).Methods("GET")
	router.HandleFunc("/product/{id}", c.Delete).Methods("DELETE")
	router.HandleFunc("/product/{id}", c.Update).Methods("PUT")
}

/*
func handleGraphqlRequests(router *mux.Router, c interfaces.GraphqlControllerInterface) {
	router.HandleFunc("/graphql", c.Graphql).Methods("GET")
}
*/
type router struct{}

func (router *router) InitRouter() *mux.Router {
	controller := ServiceContainer().InjectCrudController()
	r := mux.NewRouter().StrictSlash(true)
	handleRequests(r, controller)
	log.Println("Rest API v2.0 - Mux Routers")
	return r
}

var (
	m          *router
	routerOnce sync.Once
)

func GorillaRouter() GorillaRouterInterface {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}
