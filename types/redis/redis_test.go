//go:build unit_tests || redis_unit_tests

package redis

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

var currentHost string
var currentHostDefined bool

var currentPort string
var currentPortDefined bool

var currentDatabase string
var currentDatabaseDefined bool

var currentPassword string
var currentPasswordDefined bool

var redis_host_env_variable string = "REDIS_HOST"
var redis_port_env_variable string = "REDIS_PORT"
var redis_database_env_variable string = "REDIS_DATABASE"
var redis_password_env_variable string = "REDIS_PASSWORD"

func setUp() {

	if envHost, found := os.LookupEnv(redis_host_env_variable); found {
		currentHost = envHost
		currentHostDefined = true
	} else {
		currentHostDefined = false
	}

	if envPort, found := os.LookupEnv(redis_port_env_variable); found {
		currentPort = envPort
		currentPortDefined = true
	} else {
		currentPortDefined = false
	}

	if envDatabase, found := os.LookupEnv(redis_database_env_variable); found {
		currentDatabase = envDatabase
		currentDatabaseDefined = true
	} else {
		currentDatabaseDefined = false
	}

	if envPassword, found := os.LookupEnv(redis_password_env_variable); found {
		currentPassword = envPassword
		currentPasswordDefined = true
	} else {
		currentPasswordDefined = false
	}

	os.Unsetenv(redis_host_env_variable)
	os.Unsetenv(redis_port_env_variable)
	os.Unsetenv(redis_database_env_variable)
	os.Unsetenv(redis_password_env_variable)

}

func teardown() {

	if currentHostDefined {
		os.Setenv(redis_host_env_variable, currentHost)
	} else {
		os.Unsetenv(redis_host_env_variable)
	}

	if currentPortDefined {
		os.Setenv(redis_port_env_variable, currentPort)
	} else {
		os.Unsetenv(redis_port_env_variable)
	}

	if currentDatabaseDefined {
		os.Setenv(redis_database_env_variable, currentDatabase)
	} else {
		os.Unsetenv(redis_database_env_variable)
	}

	if currentPasswordDefined {
		os.Setenv(redis_password_env_variable, currentPassword)
	} else {
		os.Unsetenv(redis_password_env_variable)
	}

}

func TestRedisConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method without any env variable set shouldn't fail, error was '%s'.", err.Error())
	} else {

		if config.Host != "localhost" {
			t.Errorf("redis config.Host should be \"localhost\" but was \"%s\".", config.Host)
		}

		if config.Port != 6379 {
			t.Errorf("redis config.Port should be 6379 but was %d.", config.Port)
		}

		if config.Database != 0 {
			t.Errorf("redis config.Database should be 0 but was \"%d\".", config.Database)
		}
	}
}

func TestRedisConfigWithInvalidPort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_port_env_variable, "invalidport")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_PORT\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithNegativePort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_port_env_variable, "-1")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_PORT\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithPortTooBig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_port_env_variable, "100000")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_PORT\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithInvalidDatabase(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_database_env_variable, "invalid")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_DATABASE\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithNegativeDatabase(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_database_env_variable, "-1")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_DATABASE\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithValidVariables(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_host_env_variable, "127.0.0.1")
	os.Setenv(redis_port_env_variable, "1234")
	os.Setenv(redis_database_env_variable, "2")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with valid config shouldn't fail, error was '%s'.", err.Error())
	} else {

		if config.Host != "127.0.0.1" {
			t.Errorf("redis config.Host should be \"127.0.0.1\" but was \"%s\".", config.Host)
		}
		if config.Port != 1234 {
			t.Errorf("redis config.Port should be 1234 but was %d.", config.Port)
		}

		if config.Database != 2 {
			t.Errorf("redis config.Database should be 2 but was \"%d\".", config.Database)
		}
	}
}

func TestLogValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_host_env_variable, "127.0.0.1")
	os.Setenv(redis_port_env_variable, "1234")
	os.Setenv(redis_database_env_variable, "2")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with valid config shouldn't fail, error was '%s'.", err.Error())
	} else {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		logger.Info("test log", "redis config", config)

		bufferLen := buf.Len()

		if bufferLen <= 0 {
			t.Errorf("TestLogValue has failed, buffer is empty")
		} else {
			var loggedData map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &loggedData); err != nil {
				t.Errorf("TestLogValue has failed, cannot unmarshal json log")
			} else {
				redisConfig := loggedData["redis config"].(map[string]interface{})
				databaseValue := redisConfig["database"].(float64)
				if databaseValue != 2 {
					t.Errorf("TestLogValue has failed, database should be 2 but it was %0.0f", databaseValue)
				}
			}

		}
	}
}

func TestLogValueWithPassword(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(redis_host_env_variable, "127.0.0.1")
	os.Setenv(redis_port_env_variable, "1234")
	os.Setenv(redis_database_env_variable, "2")
	os.Setenv(redis_password_env_variable, "ultraSecret")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with valid config shouldn't fail, error was '%s'.", err.Error())
	} else {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		logger.Info("test log", "redis config", config)

		bufferLen := buf.Len()

		if bufferLen <= 0 {
			t.Errorf("TestLogValue has failed, buffer is empty")
		} else {
			var loggedData map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &loggedData); err != nil {
				t.Errorf("TestLogValue has failed, cannot unmarshal json log")
			} else {
				redisConfig := loggedData["redis config"].(map[string]interface{})
				passwordValue := redisConfig["password"].(string)
				if passwordValue != "*****" {
					t.Errorf("TestLogValue has failed, password should be \"*****\" but it was \"%s\"", passwordValue)
				}
			}

		}
	}
}
