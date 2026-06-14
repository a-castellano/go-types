//go:build integration_tests || unit_tests || slog_tests || slog_unit_tests

package slog

import (
	"os"
	"testing"
)

var currentDefaultLevel string
var currentDefaultLevelDefined bool

var currentFormat string
var currentFormatDefined bool

var currentAddSource string
var currentAddSourceDefined bool

var currentAppName string
var currentAppNameDefined bool

var slog_defaultlevel_env_variable string = "SLOG_LEVEL"
var slog_format_env_variable string = "SLOG_FORMAT"
var slog_addsource_env_variable string = "SLOG_ADD_SOURCE"
var slog_appname_env_variable string = "APP_NAME"

func setUp() {

	if envDefaultLevel, found := os.LookupEnv(slog_defaultlevel_env_variable); found {
		currentDefaultLevel = envDefaultLevel
		currentDefaultLevelDefined = true
	} else {
		currentDefaultLevelDefined = false
	}

	if envFormat, found := os.LookupEnv(slog_format_env_variable); found {
		currentFormat = envFormat
		currentFormatDefined = true
	} else {
		currentFormatDefined = false
	}

	if envAddSource, found := os.LookupEnv(slog_addsource_env_variable); found {
		currentAddSource = envAddSource
		currentAddSourceDefined = true
	} else {
		currentAddSourceDefined = false
	}

	if envAppName, found := os.LookupEnv(slog_appname_env_variable); found {
		currentAppName = envAppName
		currentAppNameDefined = true
	} else {
		currentAppNameDefined = false
	}

	os.Unsetenv(slog_defaultlevel_env_variable)
	os.Unsetenv(slog_format_env_variable)
	os.Unsetenv(slog_addsource_env_variable)
	os.Unsetenv(slog_appname_env_variable)

}

func teardown() {

	if currentDefaultLevelDefined {
		os.Setenv(slog_defaultlevel_env_variable, currentDefaultLevel)
	} else {
		os.Unsetenv(slog_defaultlevel_env_variable)
	}

	if currentFormatDefined {
		os.Setenv(slog_format_env_variable, currentFormat)
	} else {
		os.Unsetenv(slog_format_env_variable)
	}

	if currentAddSourceDefined {
		os.Setenv(slog_addsource_env_variable, currentAddSource)
	} else {
		os.Unsetenv(slog_addsource_env_variable)
	}

	if currentAppNameDefined {
		os.Setenv(slog_appname_env_variable, currentAppName)
	} else {
		os.Unsetenv(slog_appname_env_variable)
	}

}

func TestRedisConfigWithoutEnvVariables(t *testing.T) {

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
