//go:build integration_tests || unit_tests || redis_tests || redis_unit_tests

package redis

import (
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
		t.Errorf("NewConfig method without any env varible suited shouldn't fail, error was '%s'.", err.Error())
	} else {

		if config.Host != "localhost" {
			t.Errorf("redis config.Host should be \"localhost\" but was \"%s\".", config.Host)
		}

		if config.Port != 6379 {
			t.Errorf("redis config.Host should be 6379 but was %d.", config.Port)
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

func TestRedisConfigWithNeativePort(t *testing.T) {

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

func TestRedisConfigWithNeativeDatabase(t *testing.T) {

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
			t.Errorf("redis config.Database should be22 but was \"%d\".", config.Database)
		}
	}
}
