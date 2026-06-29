package opentelemetry

import (
	"errors"
	"os"
)

// Config is a type that holds the data required to configure OpenTelemetry.
type Config struct {
	AppName string // The name of the app, sourced from APP_NAME and used as the telemetry service.name
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := Config{}

	config.AppName = os.Getenv("APP_NAME")
	if config.AppName == "" {
		return nil, errors.New("env variable \"APP_NAME\" must be defined and have a value")
	}

	_, OtelServiceNameVarDefined := os.LookupEnv("OTEL_SERVICE_NAME")
	// APP_NAME is the only accepted source for service.name, so OTEL_SERVICE_NAME must not be set

	if OtelServiceNameVarDefined {
		return nil, errors.New("env variable \"OTEL_SERVICE_NAME\" cannot be defined. APP_NAME will be use to set that value")
	}

	_, OtelResourceAttributesVarDefined := os.LookupEnv("OTEL_RESOURCE_ATTRIBUTES")
	// For the time being this variable is forbidden, its values will be managed if required

	if OtelResourceAttributesVarDefined {
		return nil, errors.New("env variable \"OTEL_RESOURCE_ATTRIBUTES\" cannot be defined for the time being")
	}

	return &config, nil
}
