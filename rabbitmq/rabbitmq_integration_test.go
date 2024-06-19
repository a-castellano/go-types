//go:build integration_tests || rabbitmq_tests

package rabbitmq

import (
	"os"
	"testing"
)

func TestRabbitmqFailedConnection(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_HOST", "127.0.0.1")
	os.Setenv("RABBITMQ_PORT", "1123")
	os.Setenv("RABBITMQ_USER", "user")
	os.Setenv("RABBITMQ_PASSWORD", "password")

	config, _ := NewConfig()
	queueName := "test"
	testString := []byte("This is a Test")

	dial_error := config.SendMessage(queueName, testString)

	if dial_error == nil {
		t.Errorf("TestRabbitmqFailedConnection should fail.")
	}
}
