package rabbitmq

import (
	"cmp"
	"errors"
	"os"
	"strconv"
)

// Config is a type that defines required data for connecting to RabbitMQ server
type Config struct {
	host     string
	port     int
	user     string
	password string
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

	return config, nil
}

// SendMessage sends a message through queueName
func (rabbitmqConfig Config) SendMessage(queueName string, message []byte) error {

	connectionString := "amqp://" + rabbitmqConfig.User + ":" + rabbitmqConfig.Password + "@" + rabbitmqConfig.Host + ":" + strconv.Itoa(rabbitmqConfig.Port) + "/"

	conn, errDial := amqp.Dial(connectionString)
	defer conn.Close()

	if errDial != nil {
		return errDial
	}

	channel, errChannel := conn.Channel()
	defer channel.Close()
	if errChannel != nil {
		return errChannel
	}

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
			Body:         encodedNotification,
		})

	if err != nil {
		return err
	}
	return nil
}
