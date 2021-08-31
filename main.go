package main

import (
	"log"
	"net/http"
	"sync"
	"test-api-golang/graphql"
	"test-api-golang/mailserver"
	"test-api-golang/rabbitmq"
)

func main() {
	if err := CleanAndFillRepository(); err != nil {
		log.Println(err)
	}
	wg := new(sync.WaitGroup)
	wg.Add(5)

	go func() {
		log.Println(http.ListenAndServe(":8080", GorillaRouter().InitRouter()))
		wg.Done()
	}()

	go func() {
		http.Handle("/graphql", graphql.CreateGrapQLHandler())
		log.Println(http.ListenAndServe(":5000", nil))
		wg.Done()
	}()

	go func() {
		mailserver.StartMailer()
		wg.Done()
	}()

	go func() {
		mailserver.MessageLoop()
		wg.Done()
	}()

	go func() {
		rabbitmq.RecieveMessages()
		wg.Done()
	}()
	wg.Wait()
}
