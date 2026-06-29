package opentelemetry

import (
	"errors"
	"os"
)

// Config is a type that holds the data required to configure OpenTelemetry.
type Config struct {
	AppName string // The name of the app, sourced from APP_NAME and used as the telemetry service.name
	Enabled bool   // Whether telemetry is active, sourced from ENABLE_TELEMETRY; defaults to false (opt-in)
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

	enableFlagValue, enableFlagDefined := os.LookupEnv("ENABLE_TELEMETRY")
	// ENABLE_TELEMETRY is the single source of truth for whether telemetry is active.
	// It is optional: when unset, Enabled stays false (opt-in). Only "true" or "false" are accepted.
	if enableFlagDefined {
		if enableFlagValue != "true" && enableFlagValue != "false" {
			return nil, errors.New("env variable \"ENABLE_TELEMETRY\" valid values are only true or false")
		}
		config.Enabled = enableFlagValue == "true"
	}
	return &config, nil
}
