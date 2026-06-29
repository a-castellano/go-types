package opentelemetry

import (
	"cmp"
	"errors"
	"os"
)

// Config is a type that defines required data for defining slog parameters
type Config struct {
	AppName string // The name of the app
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := Config{}

	config.AppName = cmp.Or(os.Getenv("APP_NAME"), "")
	if config.AppName == "" {
		return nil, errors.New("env variable \"APP_NAME\" must be defined and have a value")
	}

	_, OtelServiceNameVarDefined := os.LookupEnv("OTEL_SERVICE_NAME")

	if OtelServiceNameVarDefined {
		return nil, errors.New("env variable \"OTEL_SERVICE_NAME\" cannot be defined. APP_NAME will be use to set that value")

	}

	return &config, nil
}
