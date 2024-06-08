package rabbitmq

import (
	"cmp"
	"os"
)

// RabbitmqConfig defines required data for connecting to RabbitMQ server

type RabbitmqConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewRabbitmqConfig() (*RabbitmqConfig, error) {
	config := new(RabbitmqConfig)
	// Get Host from RABBITMQ_HOST env variable
	config.Host = cmp.Or(os.Getenv("RABBITMQ_HOST"), "localhost")
	return config, nil
}
