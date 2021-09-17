package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"test-api-golang/interfaces"
	"test-api-golang/oauth"
	"test-api-golang/rabbitmq"

	"github.com/gorilla/mux"
)

type GorillaRouterInterface interface {
	InitRouter() *mux.Router
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the test CRUD API!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests(router *mux.Router, c interfaces.CrudControllerInterface, googleAuth *oauth.GoogleAuth) {
	router.HandleFunc("/", HomePage)
	router.Handle("/product", googleAuth.AuthMiddleware(c.List)).Methods("GET")
	router.Handle("/product", oauth.IsAuthorized(c.Create)).Methods("POST")
	router.Handle("/product/{id}", googleAuth.AuthMiddleware(c.Get)).Methods("GET")
	router.Handle("/product/{id}", oauth.IsAuthorized(c.Delete)).Methods("DELETE")
	router.Handle("/product/{id}", oauth.IsAuthorized(c.Update)).Methods("PUT")
	router.HandleFunc("/test-post", oauth.PostRequestToProductApp).Methods("GET")
	router.HandleFunc("/authorize", googleAuth.Authorize).Methods("GET")
	router.HandleFunc("/oauth2callback", googleAuth.Oauth2callback).Methods("GET")
	router.Handle("/check", googleAuth.AuthMiddleware(http.HandlerFunc(googleAuth.Check))).Methods("GET")
	router.HandleFunc("/redis/{id}", c.Get).Methods("GET")
}

type router struct{}

func (router *router) InitRouter() *mux.Router {
	controller := ServiceContainer().InjectCrudController()
	grpcClientController := ServiceContainer().InjectGrpcClientController()
	googleAuth, err := oauth.CreateGoogleAuth()
	if err != nil {
		log.Println("Could not create Google Auth config")
	}
	r := mux.NewRouter().StrictSlash(true)
	handleRequests(r, controller, googleAuth)
	rbt := rabbitmq.NewConnection()
	r.HandleFunc("/rabbitmq", rbt.SendMessage).Methods("POST")
	r.HandleFunc("/sendmail", grpcClientController.SendMail).Methods("POST")
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
