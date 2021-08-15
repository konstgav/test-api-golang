package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"test-api-golang/graphql"
	"test-api-golang/mailserver"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		os.Setenv("MAILER_REMOTE_HOST", "")
		os.Setenv("MAILER_FROM", "")
		os.Setenv("MAILER_PASSWORD", "")
	}

	if err := CleanAndFillRepository(); err != nil {
		log.Println(err)
	}
	wg := new(sync.WaitGroup)
	wg.Add(4)

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

	wg.Wait()
}
