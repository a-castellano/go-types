//go:build integration_tests || unit_tests || rabbitmq_tests

package rabbitmq

import (
	"os"
	"testing"
)

func TestRabbitmqConfigWithoutEnvVariables(t *testing.T) {

	config, err := NewRabbitmqConfig()

	if err != nil {
		t.Errorf("NewRabbitmqConfig method without any env varible suited shouldn't fail.")
	} else {
		if config.Host != "localhost" {
			t.Errorf("RabbitmqConfig Host should be \"localhost\" but was \"%s\".", config.Host)
		}
		if config.Port != 5672 {
			t.Errorf("RabbitmqConfig Host should be 5672 but was %d.", config.Port)
		}
		if config.User != "guest" {
			t.Errorf("RabbitmqConfig User should be \"guest\" but was \"%s\".", config.User)
		}
		if config.Password != "guest" {
			t.Errorf("RabbitmqConfig Password should be \"guest\" but was \"%s\".", config.Password)
		}
	}
}

func TestRabbitmqConfigWithStringAsPortValue(t *testing.T) {

	os.Setenv("RABBITMQ_PORT", "invalidport")

	_, err := NewRabbitmqConfig()

	if err == nil {
		t.Errorf("NewRabbitmqConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue1(t *testing.T) {

	os.Setenv("RABBITMQ_PORT", "65536")

	_, err := NewRabbitmqConfig()

	if err == nil {
		t.Errorf("NewRabbitmqConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue2(t *testing.T) {

	os.Setenv("RABBITMQ_PORT", "0")

	_, err := NewRabbitmqConfig()

	if err == nil {
		t.Errorf("NewRabbitmqConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithInvalidPortValue3(t *testing.T) {

	os.Setenv("RABBITMQ_PORT", "-1")

	_, err := NewRabbitmqConfig()

	if err == nil {
		t.Errorf("NewRabbitmqConfig method with \"RABBITMQ_PORT\" env variable containing invalid value should fail.")
	}
}

func TestRabbitmqConfigWithEnvVariables(t *testing.T) {

	os.Setenv("RABBITMQ_HOST", "127.0.0.1")
	os.Setenv("RABBITMQ_PORT", "1123")
	os.Setenv("RABBITMQ_USER", "user")
	os.Setenv("RABBITMQ_PASSWORD", "password")

	config, err := NewRabbitmqConfig()

	if err != nil {
		t.Errorf("NewRabbitmqConfig method with valid env varibles suited shouldn't fail.")
	} else {
		if config.Host != "127.0.0.1" {
			t.Errorf("RabbitmqConfig Host should be \"127.0.0.1\" but was \"%s\".", config.Host)
		}
		if config.Port != 1123 {
			t.Errorf("RabbitmqConfig Host should be 1123 but was %d.", config.Port)
		}
		if config.User != "user" {
			t.Errorf("RabbitmqConfig User should be \"user\" but was \"%s\".", config.User)
		}
		if config.Password != "password" {
			t.Errorf("RabbitmqConfig Password should be \"password\" but was \"%s\".", config.Password)
		}
	}
}
