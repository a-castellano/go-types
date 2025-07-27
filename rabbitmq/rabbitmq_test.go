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

var rabbitmq_host_env_variable string = "RABBITMQ_HOST"
var rabbitmq_port_env_variable string = "RABBITMQ_PORT"
var rabbitmq_user_env_variable string = "RABBITMQ_USER"
var rabbitmq_database_env_variable string = "RABBITMQ_DATABASE"
var rabbitmq_password_env_variable string = "RABBITMQ_PASSWORD"

func setUp() {

	if envHost, found := os.LookupEnv(rabbitmq_host_env_variable); found {
		currentHost = envHost
		currentHostDefined = true
	} else {
		currentHostDefined = false
	}

	if envPort, found := os.LookupEnv(rabbitmq_port_env_variable); found {
		currentPort = envPort
		currentPortDefined = true
	} else {
		currentPortDefined = false
	}

	if envUser, found := os.LookupEnv(rabbitmq_user_env_variable); found {
		currentUser = envUser
		currentUserDefined = true
	} else {
		currentUserDefined = false
	}

	if envPassword, found := os.LookupEnv(rabbitmq_password_env_variable); found {
		currentPassword = envPassword
		currentPasswordDefined = true
	} else {
		currentPasswordDefined = false
	}

	os.Unsetenv(rabbitmq_host_env_variable)
	os.Unsetenv(rabbitmq_port_env_variable)
	os.Unsetenv(rabbitmq_database_env_variable)
	os.Unsetenv(rabbitmq_password_env_variable)

}

func teardown() {

	if currentHostDefined {
		os.Setenv(rabbitmq_host_env_variable, currentHost)
	} else {
		os.Unsetenv(rabbitmq_host_env_variable)
	}

	if currentPortDefined {
		os.Setenv(rabbitmq_port_env_variable, currentPort)
	} else {
		os.Unsetenv(rabbitmq_port_env_variable)
	}

	if currentUserDefined {
		os.Setenv(rabbitmq_user_env_variable, currentUser)
	} else {
		os.Unsetenv(rabbitmq_user_env_variable)
	}

	if currentPasswordDefined {
		os.Setenv(rabbitmq_password_env_variable, currentPassword)
	} else {
		os.Unsetenv(rabbitmq_password_env_variable)
	}

}

func TestRabbitmqConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method without any env variable set shouldn't fail.")
	} else {
		if config.host != "localhost" {
			t.Errorf("RabbitMQ config.host should be \"localhost\" but was \"%s\".", config.host)
		}
		if config.port != 5672 {
			t.Errorf("RabbitMQ config.port should be 5672 but was %d.", config.port)
		}
		if config.user != "guest" {
			t.Errorf("RabbitMQ config.user should be \"guest\" but was \"%s\".", config.user)
		}
		if config.password != "guest" {
			t.Errorf("RabbitMQ config.password should be \"guest\" but was \"%s\".", config.password)
		}
	}
}

func TestRabbitmqConfigWithStringAsPortValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(rabbitmq_port_env_variable, "invalidport")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue1(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(rabbitmq_port_env_variable, "65536")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue2(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(rabbitmq_port_env_variable, "0")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue3(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(rabbitmq_port_env_variable, "-1")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv(rabbitmq_host_env_variable, "127.0.0.1")
	os.Setenv(rabbitmq_port_env_variable, "1123")
	os.Setenv(rabbitmq_user_env_variable, "user")
	os.Setenv(rabbitmq_password_env_variable, "password")

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method with valid env variables set shouldn't fail.")
	} else {
		if config.host != "127.0.0.1" {
			t.Errorf("RabbitMQ config.host should be \"127.0.0.1\" but was \"%s\".", config.host)
		}
		if config.port != 1123 {
			t.Errorf("RabbitMQ config.port should be 1123 but was %d.", config.port)
		}
		if config.user != "user" {
			t.Errorf("RabbitMQ config.user should be \"user\" but was \"%s\".", config.user)
		}
		if config.password != "password" {
			t.Errorf("RabbitMQ config.password should be \"password\" but was \"%s\".", config.password)
		}
		if config.ConnectionString != "amqp://user:password@127.0.0.1:1123/" {
			t.Errorf("RabbitMQ config.ConnectionString should be \"amqp://user:password@127.0.0.1:1123/\" but was \"%s\".", config.ConnectionString)
		}
	}
}
