//go:build unit_tests || opentelemetry_unit_tests

package opentelemetry

import (
	"os"
	"testing"
)

type envVariable struct {
	Value        string
	IsDefined    bool
	VariableName string
}

var envVariables = map[string]envVariable{
	"appName":                  {VariableName: "APP_NAME"},
	"telemetryEnabled":         {VariableName: "ENABLE_TELEMETRY"},
	"otelServiceName":          {VariableName: "OTEL_SERVICE_NAME"},
	"otelResourceAttributes":   {VariableName: "OTEL_RESOURCE_ATTRIBUTES"},
	"otelExporterOTLPEndpoint": {VariableName: "OTEL_EXPORTER_OTLP_ENDPOINT"},
}

func setUp() {

	for key, variable := range envVariables {

		if envValue, found := os.LookupEnv(variable.VariableName); found {
			variable.Value = envValue
			variable.IsDefined = true
		} else {
			variable.IsDefined = false
		}

		os.Unsetenv(variable.VariableName)

		envVariables[key] = variable
	}

}

func teardown() {

	for _, variable := range envVariables {
		if variable.IsDefined {
			os.Setenv(variable.VariableName, variable.Value)
		} else {
			os.Unsetenv(variable.VariableName)
		}
	}
}

func TestOpenTelemetryConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("NewConfig method without any env variable set should fail, error was '%s'.", err.Error())
	}
	expectedError := "env variable \"APP_NAME\" must be defined and have a value"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}
}

func TestOpenTelemetryConfigWithappNameAndOtelServiceNameVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["otelServiceName"].VariableName, "MyApp")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("NewConfig method with \"APP_NAME\" and \"OTEL_SERVICE_NAME\" env variables set should fail")
	}

	expectedError := "env variable \"OTEL_SERVICE_NAME\" cannot be defined. APP_NAME will be use to set that value"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}
}

func TestOpenTelemetryConfigWithappNameAndOtelResourceAttributesVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["otelResourceAttributes"].VariableName, "any=value")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("NewConfig method with \"OTEL_RESOURCE_ATTRIBUTES\" env variable set should fail")
	}

	expectedError := "env variable \"OTEL_RESOURCE_ATTRIBUTES\" cannot be defined for the time being"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

}

// TestOpenTelemetryConfig covers the happy path with only APP_NAME set: NewConfig
// succeeds and, since ENABLE_TELEMETRY is unset, telemetry defaults to disabled.
func TestOpenTelemetryConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestOpenTelemetryConfig should not fail")
	}

	expectedappName := "MyApp"

	if config.AppName() != expectedappName {
		t.Fatalf("Expected app name '%s' but got '%s'", expectedappName, config.appName)
	}

	if config.Enabled() {
		t.Fatalf("TestOpenTelemetryConfig should come with opentelemetry disabled")
	}
}

// TestOpenTelemetryEnabledFail covers an invalid ENABLE_TELEMETRY value: only
// "true" or "false" are accepted, anything else must make NewConfig fail.
func TestOpenTelemetryEnabledFail(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["telemetryEnabled"].VariableName, "fail")

	_, err := NewConfig()
	if err == nil {
		t.Fatalf("NewConfig method with invalid value for \"ENABLE_TELEMETRY\" env variable should fail")
	}

	expectedError := "env variable \"ENABLE_TELEMETRY\" valid values are only true or false"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

}

// TestOpenTelemetryEnabledConfig covers ENABLE_TELEMETRY set to "true": the flag
// is the single source of truth and Enabled must reflect it.
func TestOpenTelemetryEnabledConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["telemetryEnabled"].VariableName, "true")

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestOpenTelemetryEnabledConfig should not fail")
	}
	if !config.Enabled() {

		t.Fatalf("TestOpenTelemetryEnabledConfig should come with opentelemetry enabled")
	}
	if config.ExporterType() != Stdout {

		t.Errorf("TestOpenTelemetryEnabledConfig exporterType should be Stdout")
	}

}

func TestOpenTelemetryEnabledInvalidExporterURL(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["telemetryEnabled"].VariableName, "true")
	os.Setenv(envVariables["otelExporterOTLPEndpoint"].VariableName, "12")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestOpenTelemetryEnabledInvalidExporterURL should fail")
	}
	expectedError := "env variable \"OTEL_EXPORTER_OTLP_ENDPOINT\" content is not a valid endpoint"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

}

func TestOpenTelemetryEnabledInvalidExporterURLSchema(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["telemetryEnabled"].VariableName, "true")
	os.Setenv(envVariables["otelExporterOTLPEndpoint"].VariableName, "jttl://localhost:21321")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestOpenTelemetryEnabledInvalidExporterURLSchema should fail")
	}
	expectedError := "env variable \"OTEL_EXPORTER_OTLP_ENDPOINT\" schema is not a valid, only http and https are accepted"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

}

func TestOpenTelemetryEnabledInvalidExporterURLHost(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["telemetryEnabled"].VariableName, "true")
	os.Setenv(envVariables["otelExporterOTLPEndpoint"].VariableName, "http://:21321")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestOpenTelemetryEnabledInvalidExporterURLHost should fail")
	}
	expectedError := "env variable \"OTEL_EXPORTER_OTLP_ENDPOINT\" hostname is empty"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

}

func TestOpenTelemetryValidExporterURL(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["telemetryEnabled"].VariableName, "true")

	ExporterURL := "http://localhost:12345"

	os.Setenv(envVariables["otelExporterOTLPEndpoint"].VariableName, ExporterURL)

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestOpenTelemetryValidExporterURL should not fail")
	}
	if !config.Enabled() {

		t.Fatalf("TestOpenTelemetryValidExporterURL should come with opentelemetry enabled")
	}
	if config.ExporterType() != OTLP {

		t.Errorf("TestOpenTelemetryValidExporterURL exporterType should be OTLP")
	}
	if config.ExporterURL() != ExporterURL {

		t.Errorf("TestOpenTelemetryValidExporterURL URL should be \"%s\", it was \"%s\"", ExporterURL, config.ExporterURL())
	}

}
