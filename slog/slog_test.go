//go:build integration_tests || unit_tests || slog_tests || slog_unit_tests

package slog

import (
	formerslog "log/slog"
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

func TestSlogConfigWithoutEnvVariables(t *testing.T) {

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

func TestSlogConfigWithAppNameVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with \"APP_NAME\" env variable set shouldn't fail, error was '%s'.", err.Error())
	} else {
		// check default values
		if config.AppName != "MyApp" {

			t.Errorf("NewConfig method with \"APP_NAME\" env variable set should configure config.AppName with that value, actual value was '%s'.", config.AppName)
		}

		if config.DefaultLevel != formerslog.LevelInfo {
			t.Errorf("NewConfig default level should be %d but it was %d.", formerslog.LevelInfo, config.DefaultLevel)
		}
	}
}

func TestSlogConfigWithInvalidLogLevelValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")
	os.Setenv(slog_defaultlevel_env_variable, "invalid")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with an invalid \"SLOG_LEVEL\" value set should fail")
	} else {
		expectedError := "log level defined by `SLOG_LEVEL` variable only accepts the following values: \"Debug\", \"Info\", \"Warn\" or \"Error\". \"invalid\" is not a valid value."
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
		}
	}
}

func TestSlogConfigWithInvalidFormatValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")
	os.Setenv(slog_format_env_variable, "invalid")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with an invalid \"SLOG_FORMAT\" value set should fail")
	} else {
		expectedError := "log format defined by `SLOG_FORMAT` variable only accepts the values \"JSON\" or \"plain\", \"invalid\" is not a valid value"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
		}
	}
}

func TestSlogConfigWithInvalidAddSourceValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")
	os.Setenv(slog_addsource_env_variable, "invalid")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with an invalid \"SLOG_ADD_SOURCE\" value set should fail")
	} else {
		expectedError := "add log source property defined by `SLOG_ADD_SOURCE` variable only accepts the values \"true\" or \"false\", \"invalid\" is not a valid value"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
		}
	}
}

func TestSlogConfigLevelDebug(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")
	os.Setenv(slog_defaultlevel_env_variable, "Debug")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method call where we are testing log levels shouldn't fail, error was '%s'.", err.Error())
	} else {
		expectedLevel := formerslog.LevelDebug
		if config.DefaultLevel != expectedLevel {
			t.Errorf("NewConfig default level should be %d but it was %d.", expectedLevel, config.DefaultLevel)
		}
	}
}

func TestSlogConfigLevelWarn(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")
	os.Setenv(slog_defaultlevel_env_variable, "Warn")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method call where we are testing log levels shouldn't fail, error was '%s'.", err.Error())
	} else {
		expectedLevel := formerslog.LevelWarn
		if config.DefaultLevel != expectedLevel {
			t.Errorf("NewConfig default level should be %d but it was %d.", expectedLevel, config.DefaultLevel)
		}
	}
}

func TestSlogConfigLevelError(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(slog_appname_env_variable, "MyApp")
	os.Setenv(slog_defaultlevel_env_variable, "Error")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method call where we are testing log levels shouldn't fail, error was '%s'.", err.Error())
	} else {
		expectedLevel := formerslog.LevelError
		if config.DefaultLevel != expectedLevel {
			t.Errorf("NewConfig default level should be %d but it was %d.", expectedLevel, config.DefaultLevel)
		}
	}
}
