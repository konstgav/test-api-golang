package rabbitmq

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	messageTTL     = 90 // time to live for message in queue
	resendTime     = 15 // time to return message in queue
	MaxOutstanding = 5  // Limit the number of workers for checking URLs
)

type URLCheck struct {
	Link string `json:"link"`
}

type RabbitMQConnection struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func CreateConnection() *RabbitMQConnection {
	ch, q, err := CreateQueue(messageTTL)
	go recieve_msg(messageTTL, resendTime)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &RabbitMQConnection{Channel: ch, Queue: q}
}

func (rbt RabbitMQConnection) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var entity URLCheck
	json.NewDecoder(r.Body).Decode(&entity)
	log.Println(entity.Link)
	RabbitMQSendMessage(rbt.Channel, rbt.Queue, []byte(entity.Link))
}

func CreateQueue(messageTTL int) (*amqp.Channel, amqp.Queue, error) {
	// TODO: move credentials & URL to environmental variables
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	//defer ch.Close()
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{"x-message-ttl": messageTTL * 1000}, // arguments
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return ch, q, err
}

func RabbitMQSendMessage(ch *amqp.Channel, q amqp.Queue, body []byte) error {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}
	log.Printf("Sent %s", string(body))
	return nil
}

func recieve_msg(messageTTL int, resendTime int) {
	ch, q, err := CreateQueue(messageTTL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	forever := make(chan bool)
	for i := 0; i < MaxOutstanding; i++ {
		go worker(ch, q, resendTime)
	}
	<-forever
}

func worker(ch *amqp.Channel, q amqp.Queue, resendTime int) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for d := range msgs {
		log.Printf("Recieved  %v %s", d.DeliveryTag, d.Body)
		t0 := time.Now()
		success := CheckURL(string(d.Body))
		t1 := time.Now()
		checkURLDuration := t1.Sub(t0)
		timeToWait := time.Duration(resendTime)*time.Second - checkURLDuration
		if success {
			log.Println("Correct link", d.DeliveryTag, string(d.Body))
			d.Ack(false)
		} else {
			log.Println(timeToWait)
			time.Sleep(timeToWait)
			log.Println("Incorrect link", d.DeliveryTag, string(d.Body))
			d.Nack(false, true)
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
