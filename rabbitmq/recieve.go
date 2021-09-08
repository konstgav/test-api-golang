package rabbitmq

import (
	"log"
	"net/http"
	"os"
	"time"
)

func RecieveMessages() {
	session := RabbitMQConnection{
		Name:        "task_queue",
		isReadyChan: make(chan bool),
		done:        make(chan bool),
	}
	rabbitmq_URI := os.Getenv("RABBITMQ_URI")
	if rabbitmq_URI == "" {
		panic("Environmental variable RABBITMQ_URI do not set")
	}
	go session.handleReconnect(rabbitmq_URI)
	<-session.isReadyChan
	session.worker()
}

func (session *RabbitMQConnection) worker() {
	msgs, err := session.Stream()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for d := range msgs {
		log.Printf("Recieved  %v %s", d.DeliveryTag, d.Body)
		success := CheckURL(string(d.Body))
		if success {
			log.Println("Correct link", d.DeliveryTag, string(d.Body))
			d.Ack(false)
		} else {
			log.Println("Incorrect link", d.DeliveryTag, string(d.Body))
			if time.Since(d.Timestamp) > time.Duration(messageTTL)*time.Second {
				d.Ack(false)
			} else {
				d.Nack(false, false)
			}
		}
	}
}

func CheckURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	status := resp.Status
	log.Println(status)
	return status[:2] == "20"
}
