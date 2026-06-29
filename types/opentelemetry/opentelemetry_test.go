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
	"appName":                {VariableName: "APP_NAME"},
	"otelServiceName":        {VariableName: "OTEL_SERVICE_NAME"},
	"otelResourceAttributes": {VariableName: "OTEL_RESOURCE_ATTRIBUTES"},
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
		t.Errorf("NewConfig method without any env variable set should fail, error was '%s'.", err.Error())
	} else {
		expectedError := "env variable \"APP_NAME\" must be defined and have a value"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
		}
	}
}

func TestOpenTelemetryConfigWithAppNameAndOtelServiceNameVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["otelServiceName"].VariableName, "MyApp")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"APP_NAME\" and \"OTEL_SERVICE_NAME\" env variables set should fail")
	} else {
		expectedError := "env variable \"OTEL_SERVICE_NAME\" cannot be defined. APP_NAME will be use to set that value"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
		}

	}
}

func TestOpenTelemetryConfigWithAppNameAndOtelResourceAttributesVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")
	os.Setenv(envVariables["otelResourceAttributes"].VariableName, "any=value")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"OTEL_RESOURCE_ATTRIBUTES\" env variable set should fail")
	} else {
		expectedError := "env variable \"OTEL_RESOURCE_ATTRIBUTES\" cannot be defined for the time being"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
		}

	}
}

func TestOpenTelemetryConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(envVariables["appName"].VariableName, "MyApp")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("TestOpenTelemetryConfig should not fail")
	} else {
		expectedAppName := "MyApp"
		if config.AppName != expectedAppName {
			t.Fatalf("Expected app name '%s' but got '%s'", expectedAppName, config.AppName)
		}

	}
}
