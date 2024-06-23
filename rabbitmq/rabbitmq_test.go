//go:build integration_tests || unit_tests || rabbitmq_tests || rabbitmq_unit_tests

package rabbitmq

import (
	"os"
	"testing"
)

func teardown() {
	os.Unsetenv("RABBITMQ_HOST")
	os.Unsetenv("RABBITMQ_PORT")
	os.Unsetenv("RABBITMQ_USER")
	os.Unsetenv("RABBITMQ_PASSWORD")
}

func TestRabbitmqConfigWithoutEnvVariables(t *testing.T) {

	defer teardown()

	config, err := NewConfig()

	if err != nil {
		t.Errorf("NewConfig method without any env varible suited shouldn't fail.")
	} else {
		if config.host != "localhost" {
			t.Errorf("Rabbitmq config.host should be \"localhost\" but was \"%s\".", config.host)
		}
		if config.port != 5672 {
			t.Errorf("Rabbitmq config.host should be 5672 but was %d.", config.port)
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

	defer teardown()

	os.Setenv("RABBITMQ_PORT", "invalidport")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue1(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_PORT", "65536")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue2(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_PORT", "0")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue3(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_PORT", "-1")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("NewConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithEnvVariables(t *testing.T) {

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
		if config.connectionString != "amqp://user:password@127.0.0.1:1123/" {
			t.Errorf("Rabbitmq config.connectionString should be \"amqp://user:password@127.0.0.1:1123/\" but was \"%s\".", config.connectionString)
		}
	}
}
