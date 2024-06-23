package rabbitmq

import (
	"cmp"
	"context"
	"errors"
	"os"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config is a type that defines required data for connecting to RabbitMQ server
type Config struct {
	host             string
	port             int
	user             string
	password         string
	connectionString string
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := new(Config)

	// Get Host from RABBITMQ_HOST env variable
	config.host = cmp.Or(os.Getenv("RABBITMQ_HOST"), "localhost")

	// Get User from RABBITMQ_USER env variable
	config.user = cmp.Or(os.Getenv("RABBITMQ_USER"), "guest")

	// Get Password from RABBITMQ_PASSWORD env variable
	config.password = cmp.Or(os.Getenv("RABBITMQ_PASSWORD"), "guest")

	// Get Port from RABBITMQ_PORT env variable and validate its value
	var portAtoiErr error
	config.port, portAtoiErr = strconv.Atoi(cmp.Or(os.Getenv("RABBITMQ_PORT"), "5672"))

	if portAtoiErr != nil {
		return config, portAtoiErr
	}

	if config.port <= 0 || config.port >= 65536 {
		return config, errors.New("RabbitMQ portvalue must be between 1 and 65535")
	}

	config.connectionString = "amqp://" + config.user + ":" + config.password + "@" + config.host + ":" + strconv.Itoa(config.port) + "/"
	return config, nil
}

// SendMessage sends a message through queueName
func (rabbitmqConfig Config) SendMessage(queueName string, message []byte) error {

	conn, errDial := amqp.Dial(rabbitmqConfig.connectionString)
	if errDial != nil {
		return errDial
	}
	defer conn.Close()

	channel, errChannel := conn.Channel()

	if errChannel != nil {
		return errChannel
	}

	defer channel.Close()

	queue, errQueue := channel.QueueDeclare(

		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments

	)

	if errQueue != nil {
		return errQueue
	}

	// send message

	err := channel.Publish(

		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})

	if err != nil {
		return err
	}

	return nil
}

// ReceiveMessages stores messages in channel until it is closed using context
func (rabbitmqConfig Config) ReceiveMessages(ctx context.Context, queueName string, messages chan<- []byte, errors chan<- error) {

	conn, errDial := amqp.Dial(rabbitmqConfig.connectionString)

	if errDial != nil {
		errors <- errDial
		return
	}

	defer conn.Close()

	channel, errChannel := conn.Channel()

	if errChannel != nil {
		errors <- errChannel
		return
	}

	_, errQueue := channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if errQueue != nil {
		errors <- errQueue
		return
	}

	errChannelQos := channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if errChannelQos != nil {
		errors <- errChannelQos
		return
	}

	messagesToReceive, errMessageReceived := channel.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if errMessageReceived != nil {
		errors <- errMessageReceived
		return
	}

	for receivedMessage := range messagesToReceive {

		messages <- receivedMessage.Body
		select {
		case <-ctx.Done(): //exit function
			errors <- nil
			return
		default:
			continue
		}
	}

	errors <- nil
}
