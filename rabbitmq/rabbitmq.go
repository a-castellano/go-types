package rabbitmq

import (
	"cmp"
	"errors"
	"os"
	"strconv"
)

// RabbitmqConfig defines required data for connecting to RabbitMQ server

type config struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewRabbitmqConfig() (*config, error) {
	config := new(config)

	// Get Host from RABBITMQ_HOST env variable
	config.Host = cmp.Or(os.Getenv("RABBITMQ_HOST"), "localhost")

	// Get Port from RABBITMQ_PORT env variable and validate its value
	var portAtoiErr error
	config.Port, portAtoiErr = strconv.Atoi(cmp.Or(os.Getenv("RABBITMQ_PORT"), "5672"))
	if portAtoiErr != nil {
		return config, portAtoiErr
	} else {
		if config.Port <= 0 || config.Port >= 65536 {
			return config, errors.New("RabbitMQ portvalue must be between 1 and 65535")
		}
	}

	// Get User from RABBITMQ_USER env variable
	config.User = cmp.Or(os.Getenv("RABBITMQ_USER"), "guest")

	// Get Password from RABBITMQ_PASSWORD env variable
	config.Password = cmp.Or(os.Getenv("RABBITMQ_PASSWORD"), "guest")

	return config, nil
}
