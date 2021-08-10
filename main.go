package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	log.Fatal(http.ListenAndServe(":5000", GorillaRouter().InitRouter()))
}
