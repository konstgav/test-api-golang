package main

import (
	"log"
	"net/http"
	"sync"
	"test-api-golang/graphql"
)

func main() {
	if err := CleanAndFillRepository(); err != nil {
		log.Println(err)
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		log.Println(http.ListenAndServe(":8080", GorillaRouter().InitRouter()))
		wg.Done()
	}()

	go func() {
		http.Handle("/graphql", graphql.CreateGrapQLHandler())
		log.Println(http.ListenAndServe(":5000", nil))
		wg.Done()
	}()

	wg.Wait()
}
