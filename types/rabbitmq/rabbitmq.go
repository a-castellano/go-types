package rabbitmq

import (
	"cmp"
	"errors"
	"log/slog"
	"net/url"
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
		return nil, portAtoiErr
	}

	if config.port <= 0 || config.port >= 65536 {
		return nil, errors.New("RabbitMQ port value must be between 1 and 65535")
	}

	connectionURL := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(config.user, config.password),
		Host:   config.host + ":" + strconv.Itoa(config.port),
		Path:   "/",
	}

	config.ConnectionString = connectionURL.String()
	return config, nil
}

// LogValue allows to log conection URL masking password value
func (config Config) LogValue() slog.Value {
	connectionURL := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(config.user, "xxxxx"),
		Host:   config.host + ":" + strconv.Itoa(config.port),
		Path:   "/",
	}

	return slog.GroupValue(
		slog.String("url", connectionURL.String()),
	)
}
