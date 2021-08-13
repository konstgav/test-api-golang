package main

import (
	"log"
	"net/http"
)

func main() {
	if err := CleanAndFillRepository(); err != nil {
		log.Println(err)
	}
	log.Println("Rest API v2.0 - Mux Routers")
	log.Println(http.ListenAndServe(":8080", GorillaRouter().InitRouter()))
}
