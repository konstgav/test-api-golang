package rabbitmq

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	messageTTL     = 90              // time to live for message in queue
	resendTime     = 15              // time to return message in queue
	MaxOutstanding = 5               // Limit the number of workers for checking URLs
	reconnectDelay = 5 * time.Second // When reconnecting to the server after connection failure
	reInitDelay    = 2 * time.Second // When setting up the channel after a channel exception
	resendDelay    = 5 * time.Second // When resending messages the server didn't confirm
)

type URLCheck struct {
	Link string `json:"link"`
}

type RabbitMQConnection struct {
	Connection      *amqp.Connection
	Channel         *amqp.Channel
	Queue           amqp.Queue
	isReady         bool
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReadyChan     chan bool
	Name            string
}

func NewConnection() *RabbitMQConnection {
	session := RabbitMQConnection{
		Name: "task_queue",
		done: make(chan bool),
	}
	rabbitmq_URI := os.Getenv("RABBITMQ_URI")
	if rabbitmq_URI == "" {
		panic("Environmental variable RABBITMQ_URI do not set")
	}
	go session.handleReconnect(rabbitmq_URI)
	return &session
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (session *RabbitMQConnection) handleReconnect(addr string) {
	for {
		session.isReady = false
		log.Println("Attempting to connect RabbitMQ server")

		conn, err := session.connect(addr)

		if err != nil {
			log.Println("Failed to connect RabbitMQ server. Retrying...")

			select {
			case <-session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := session.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (session *RabbitMQConnection) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	session.changeConnection(conn)
	log.Println("Connected to RabbitMQ server!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (session *RabbitMQConnection) handleReInit(conn *amqp.Connection) bool {
	for {
		session.isReady = false

		err := session.init(conn)

		if err != nil {
			log.Println("Failed to initialize RabbitMQ channel. Retrying...")

			select {
			case <-session.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-session.done:
			return true
		case <-session.notifyConnClose:
			log.Println("RabbitMQ connection closed. Reconnecting...")
			return false
		case <-session.notifyChanClose:
			log.Println("RabbitMQ channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (session *RabbitMQConnection) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	err = ch.Confirm(false)

	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		session.Name,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		amqp.Table{"x-message-ttl": messageTTL * 1000}, // Arguments
	)

	if err != nil {
		return err
	}

	session.changeChannel(ch)
	session.isReady = true
	session.isReadyChan <- true
	log.Println("RabbitMQ setup!")

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (session *RabbitMQConnection) changeConnection(connection *amqp.Connection) {
	session.Connection = connection
	session.notifyConnClose = make(chan *amqp.Error)
	session.Connection.NotifyClose(session.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (session *RabbitMQConnection) changeChannel(channel *amqp.Channel) {
	session.Channel = channel
	session.notifyChanClose = make(chan *amqp.Error)
	session.notifyConfirm = make(chan amqp.Confirmation, 1)
	session.Channel.NotifyClose(session.notifyChanClose)
	session.Channel.NotifyPublish(session.notifyConfirm)
}

// Push will push data onto the queue, and wait for a confirm.
// If no confirms are received until within the resendTimeout,
// it continuously re-sends messages until a confirm is received.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (session *RabbitMQConnection) Push(data []byte) error {
	if !session.isReady {
		return errors.New("failed to push: not connected")
	}
	for {
		err := session.UnsafePush(data)
		if err != nil {
			log.Println("Push failed. Retrying...")
			select {
			case <-session.done:
				return errors.New("rabbitmq session is shutting down")
			case <-time.After(resendDelay):
			}
			continue
		}
		select {
		case confirm := <-session.notifyConfirm:
			if confirm.Ack {
				log.Println("Push confirmed!")
				return nil
			}
		case <-time.After(resendDelay):
		}
		log.Println("Push didn't confirm. Retrying...")
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (session *RabbitMQConnection) UnsafePush(data []byte) error {
	if !session.isReady {
		return errors.New("not connected to a rabbitmq server")
	}
	return session.Channel.Publish(
		"",           // Exchange
		session.Name, // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

// Stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (session *RabbitMQConnection) Stream() (<-chan amqp.Delivery, error) {
	if !session.isReady {
		return nil, errors.New("not connected to a rabbitmq server")
	}
	return session.Channel.Consume(
		session.Name,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

// Close will cleanly shutdown the channel and connection.
func (session *RabbitMQConnection) Close() error {
	if !session.isReady {
		return errors.New("already closed: not connected to the server")
	}
	err := session.Channel.Close()
	if err != nil {
		return err
	}
	err = session.Connection.Close()
	if err != nil {
		return err
	}
	close(session.done)
	session.isReady = false
	return nil
}

func (session *RabbitMQConnection) SendMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit rabbit endpoint")
	var entity URLCheck
	json.NewDecoder(r.Body).Decode(&entity)
	err := session.UnsafePush([]byte(entity.Link))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Sent ", entity.Link)
}
