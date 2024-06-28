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

func setUp() {

	if envHost, found := os.LookupEnv("REDIS_HOST"); found {
		currentHost = envHost
		currentHostDefined = true
	} else {
		currentHostDefined = false
	}

	if envPort, found := os.LookupEnv("REDIS_PORT"); found {
		currentPort = envPort
		currentPortDefined = true
	} else {
		currentPortDefined = false
	}

	if envDatabase, found := os.LookupEnv("REDIS_DATABASE"); found {
		currentDatabase = envDatabase
		currentDatabaseDefined = true
	} else {
		currentDatabaseDefined = false
	}

	if envPassword, found := os.LookupEnv("REDIS_PASSWORD"); found {
		currentPassword = envPassword
		currentPasswordDefined = true
	} else {
		currentPasswordDefined = false
	}

	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_DATABASE")
	os.Unsetenv("REDIS_PASSWORD")

}

func teardown() {

	if currentHostDefined {
		os.Setenv("REDIS_HOST", currentHost)
	}

	if currentPortDefined {
		os.Setenv("REDIS_PORT", currentPort)
	}

	if currentDatabaseDefined {
		os.Setenv("REDIS_DATABASE", currentDatabase)
	}
	if currentPasswordDefined {
		os.Setenv("REDIS_PASSWORD", currentPassword)
	}

}

func TestRedisConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method without any env varible suited shouldn't fail, error was '%s'.", err.Error())
	} else {

		if config.host != "localhost" {
			t.Errorf("redis config.host should be \"localhost\" but was \"%s\".", config.host)
		}

		if config.port != 6379 {
			t.Errorf("redis config.host should be 6379 but was %d.", config.port)
		}

		if config.database != 0 {
			t.Errorf("redis config.database should be 0 but was \"%d\".", config.database)
		}
	}
}

func TestRedisConfigWithInvalidPort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("REDIS_PORT", "invalidport")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_PORT\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithInvalidDatabase(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("REDIS_DATABASE", "-1")
	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"REDIS_DATABASE\" env variable containing invalid value should fail.")
	}

}

func TestRedisConfigWithValidVariables(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("REDIS_HOST", "127.0.0.1")
	os.Setenv("REDIS_PORT", "1234")
	os.Setenv("REDIS_DATABASE", "2")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with valid config shouldn't fail, error was '%s'.", err.Error())
	} else {

		if config.host != "127.0.0.1" {
			t.Errorf("redis config.host should be \"127.0.0.1\" but was \"%s\".", config.host)
		}
		if config.port != 1234 {
			t.Errorf("redis config.port should be 1234 but was %d.", config.port)
		}

		if config.database != 2 {
			t.Errorf("redis config.database should be22 but was \"%d\".", config.database)
		}
	}
}
