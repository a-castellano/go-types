package opentelemetry

import (
	"errors"
	"net/url"
	"os"
)

// ExporterType identifies which exporter family the telemetry pipeline will use.
type ExporterType int

// Exporter types. Stdout is the zero value, used when OTEL_EXPORTER_OTLP_ENDPOINT
// is not defined; OTLP is used when the endpoint is set.
const (
	Stdout ExporterType = iota
	OTLP
)

// Config is a type that holds the data required to configure OpenTelemetry.
type Config struct {
	appName      string // The name of the app, sourced from APP_NAME and used as the telemetry service.name
	enabled      bool   // Whether telemetry is active, sourced from ENABLE_TELEMETRY; defaults to false (opt-in)
	exporterType ExporterType
	exporterURL  string
}

// AppName returns the name of the app, used as the telemetry service.name.
func (c *Config) AppName() string {
	return c.appName
}

// Enabled returns whether telemetry is active.
func (c *Config) Enabled() bool {
	return c.enabled
}

// ExporterType returns the exporter type, derived from OTEL_EXPORTER_OTLP_ENDPOINT.
func (c *Config) ExporterType() ExporterType {
	return c.exporterType
}

// ExporterURL returns the OTLP endpoint URL; it is empty when the exporter type is Stdout.
func (c *Config) ExporterURL() string {
	return c.exporterURL
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := Config{}

	config.appName = os.Getenv("APP_NAME")
	if config.appName == "" {
		return nil, errors.New("env variable \"APP_NAME\" must be defined and have a value")
	}

	_, otelServiceNameVarDefined := os.LookupEnv("OTEL_SERVICE_NAME")
	// APP_NAME is the only accepted source for service.name, so OTEL_SERVICE_NAME must not be set

	if otelServiceNameVarDefined {
		return nil, errors.New("env variable \"OTEL_SERVICE_NAME\" cannot be defined. APP_NAME will be use to set that value")
	}

	_, otelResourceAttributesVarDefined := os.LookupEnv("OTEL_RESOURCE_ATTRIBUTES")
	// For the time being this variable is forbidden, its values will be managed if required

	if otelResourceAttributesVarDefined {
		return nil, errors.New("env variable \"OTEL_RESOURCE_ATTRIBUTES\" cannot be defined for the time being")
	}

	enableFlagValue, enableFlagDefined := os.LookupEnv("ENABLE_TELEMETRY")
	// ENABLE_TELEMETRY is the single source of truth for whether telemetry is active.
	// It is optional: when unset, enabled stays false (opt-in). Only "true" or "false" are accepted.
	if enableFlagDefined {
		if enableFlagValue != "true" && enableFlagValue != "false" {
			return nil, errors.New("env variable \"ENABLE_TELEMETRY\" valid values are only true or false")
		}
		config.enabled = enableFlagValue == "true"
	}

	if config.enabled {
		otelExporter, otelExporterVarDefined := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")

		config.exporterType = Stdout

		if otelExporterVarDefined {
			parsedURL, parseError := url.ParseRequestURI(otelExporter)
			if parseError != nil {
				return nil, errors.New("env variable \"OTEL_EXPORTER_OTLP_ENDPOINT\" content is not a valid endpoint")
			}
			if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
				return nil, errors.New("env variable \"OTEL_EXPORTER_OTLP_ENDPOINT\" scheme is not valid, only http and https are accepted")
			}
			if parsedURL.Hostname() == "" {
				return nil, errors.New("env variable \"OTEL_EXPORTER_OTLP_ENDPOINT\" hostname is empty")
			}

			config.exporterType = OTLP
			config.exporterURL = otelExporter
		}
	}

	return &config, nil
}
