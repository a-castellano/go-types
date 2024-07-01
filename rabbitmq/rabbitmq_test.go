//go:build integration_tests || unit_tests || rabbitmq_tests || rabbitmq_unit_tests

package rabbitmq

import (
	"os"
	"testing"
)

var currentHost string
var currentHostDefined bool

var currentPort string
var currentPortDefined bool

var currentUser string
var currentUserDefined bool

var currentPassword string
var currentPasswordDefined bool

func setUp() {

	if envHost, found := os.LookupEnv("RABBITMQ_HOST"); found {
		currentHost = envHost
		currentHostDefined = true
	} else {
		currentHostDefined = false
	}

	if envPort, found := os.LookupEnv("RABBITMQ_PORT"); found {
		currentPort = envPort
		currentPortDefined = true
	} else {
		currentPortDefined = false
	}

	if envUser, found := os.LookupEnv("RABBITMQ_USER"); found {
		currentUser = envUser
		currentUserDefined = true
	} else {
		currentUserDefined = false
	}

	if envPassword, found := os.LookupEnv("RABBITMQ_PASSWORD"); found {
		currentPassword = envPassword
		currentPasswordDefined = true
	} else {
		currentPasswordDefined = false
	}

	os.Unsetenv("RABBITMQ_HOST")
	os.Unsetenv("RABBITMQ_PORT")
	os.Unsetenv("RABBITMQ_DATABASE")
	os.Unsetenv("RABBITMQ_PASSWORD")

}

func teardown() {

	if currentHostDefined {
		os.Setenv("RABBITMQ_HOST", currentHost)
	} else {
		os.Unsetenv("RABBITMQ_HOST")
	}

	if currentPortDefined {
		os.Setenv("RABBITMQ_PORT", currentPort)
	} else {
		os.Unsetenv("RABBITMQ_PORT")
	}

	if currentUserDefined {
		os.Setenv("RABBITMQ_USER", currentUser)
	} else {
		os.Unsetenv("RABBITMQ_USER")
	}

	if currentPasswordDefined {
		os.Setenv("RABBITMQ_PASSWORD", currentPassword)
	} else {
		os.Unsetenv("RABBITMQ_PASSWORD")
	}

}

func TestRabbitmqConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method without any env varible suited shouldn't fail.")
	} else {
		if config.host != "localhost" {
			t.Errorf("Rabbitmq config.host should be \"localhost\" but was \"%s\".", config.host)
		}
		if config.port != 5672 {
			t.Errorf("Rabbitmq config.port should be 5672 but was %d.", config.port)
		}
		if config.user != "guest" {
			t.Errorf("Rabbitmq config.user should be \"guest\" but was \"%s\".", config.user)
		}
		if config.password != "guest" {
			t.Errorf("Rabbitmq config.password should be \"guest\" but was \"%s\".", config.password)
		}
	}
}

func TestRabbitmqConfigWithStringAsPortValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("RABBITMQ_PORT", "invalidport")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue1(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("RABBITMQ_PORT", "65536")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue2(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("RABBITMQ_PORT", "0")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue3(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("RABBITMQ_PORT", "-1")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("RABBITMQ_HOST", "127.0.0.1")
	os.Setenv("RABBITMQ_PORT", "1123")
	os.Setenv("RABBITMQ_USER", "user")
	os.Setenv("RABBITMQ_PASSWORD", "password")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with valid env varibles suited shouldn't fail.")
	} else {
		if config.host != "127.0.0.1" {
			t.Errorf("Rabbitmq config.host should be \"127.0.0.1\" but was \"%s\".", config.host)
		}
		if config.port != 1123 {
			t.Errorf("Rabbitmq config.host should be 1123 but was %d.", config.port)
		}
		if config.user != "user" {
			t.Errorf("Rabbitmq config.user should be \"user\" but was \"%s\".", config.user)
		}
		if config.password != "password" {
			t.Errorf("Rabbitmq config.password should be \"password\" but was \"%s\".", config.password)
		}
		if config.ConnectionString != "amqp://user:password@127.0.0.1:1123/" {
			t.Errorf("Rabbitmq config.ConnectionString should be \"amqp://user:password@127.0.0.1:1123/\" but was \"%s\".", config.ConnectionString)
		}
	}
}
