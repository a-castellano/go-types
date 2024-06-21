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

func TestRabbitmqSendMessage(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_HOST", "rabbitmq")
	os.Setenv("RABBITMQ_PORT", "5672")
	os.Setenv("RABBITMQ_USER", "guest")
	os.Setenv("RABBITMQ_PASSWORD", "guest")

	config, _ := NewConfig()
	queueName := "test"
	testString := []byte("This is a Test")

	sendError := config.SendMessage(queueName, testString)

	if sendError != nil {
		t.Errorf("TestRabbitmqSendMessage shouldn't fail. Error was '%s'.", sendError.Error())
	}
}
