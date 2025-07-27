package rabbitmq

import (
	"cmp"
	"errors"
	"os"
	"strconv"
)

// Config is a type that defines required data for connecting to RabbitMQ server
type Config struct {
	host             string
	port             int
	user             string
	password         string
	ConnectionString string
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := new(Config)

	// Get host from RABBITMQ_HOST env variable
	config.host = cmp.Or(os.Getenv("RABBITMQ_HOST"), "localhost")

	// Get user from RABBITMQ_USER env variable
	config.user = cmp.Or(os.Getenv("RABBITMQ_USER"), "guest")

	// Get password from RABBITMQ_PASSWORD env variable
	config.password = cmp.Or(os.Getenv("RABBITMQ_PASSWORD"), "guest")

	// Get port from RABBITMQ_PORT env variable and validate its value
	var portAtoiErr error
	config.port, portAtoiErr = strconv.Atoi(cmp.Or(os.Getenv("RABBITMQ_PORT"), "5672"))

	if portAtoiErr != nil {
		return config, portAtoiErr
	}

	if config.port <= 0 || config.port >= 65536 {
		return config, errors.New("RabbitMQ port value must be between 1 and 65535")
	}

	config.ConnectionString = "amqp://" + config.user + ":" + config.password + "@" + config.host + ":" + strconv.Itoa(config.port) + "/"
	return config, nil
}
