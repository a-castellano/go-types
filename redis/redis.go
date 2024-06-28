package redis

import (
	"cmp"
	"errors"
	"os"
	"strconv"
)

// Config is a type that defines required data for connecting to Redis server
type Config struct {
	host     string
	port     int
	password string
	database int
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := Config{}

	// Get host from REDIS_HOST env variable
	config.host = cmp.Or(os.Getenv("REDIS_HOST"), "localhost")

	// Get port from  REDIS_PORT env variable
	var portAtoiErr error
	config.port, portAtoiErr = strconv.Atoi(cmp.Or(os.Getenv("REDIS_PORT"), "6379"))

	if portAtoiErr != nil {
		return nil, portAtoiErr
	}

	if config.port <= 0 || config.port >= 65536 {
		return nil, errors.New("Redis port value must be between 1 and 65535")
	}

	// Get database from REDIS_DATABASE
	var databaseAtoiErr error
	config.database, databaseAtoiErr = strconv.Atoi(cmp.Or(os.Getenv("REDIS_DATABASE"), "0"))

	if databaseAtoiErr != nil {
		return nil, databaseAtoiErr
	}

	if config.database < 0 {
		return nil, errors.New("Redis database value must be a positive integer")
	}

	// Get Password from REDIS_PASSWORD env variable
	config.password = cmp.Or(os.Getenv("REDIS_PASSWORD"), "")

	return &config, nil
}
