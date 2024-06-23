//go:build integration_tests || rabbitmq_tests

package rabbitmq

import (
	"context"
	"os"
	"testing"
	"time"
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

func TestRabbitmqReceiveMessageFromInvalidConfig(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_HOST", "rabbitmq")
	os.Setenv("RABBITMQ_PORT", "5672")
	os.Setenv("RABBITMQ_USER", "guest")
	os.Setenv("RABBITMQ_PASSWORD", "badPassword")

	config, _ := NewConfig()
	queueName := "test"
	messagesReceived := make(chan []byte)

	ctx, _ := context.WithCancel(context.Background())

	receiveErrors := make(chan error)

	go config.ReceiveMessages(ctx, queueName, messagesReceived, receiveErrors)

	time.Sleep(3 * time.Second)

	receiveError := <-receiveErrors

	if receiveError == nil {
		t.Errorf("TestRabbitmqReceiveMessageFromInvalidConfig should fail.")
	}
}

func TestRabbitmqReceiveMessage(t *testing.T) {

	defer teardown()

	os.Setenv("RABBITMQ_HOST", "rabbitmq")
	os.Setenv("RABBITMQ_PORT", "5672")
	os.Setenv("RABBITMQ_USER", "guest")
	os.Setenv("RABBITMQ_PASSWORD", "guest")

	config, _ := NewConfig()
	queueName := "test"
	messagesReceived := make(chan []byte)

	ctx, cancel := context.WithCancel(context.Background())

	receiveErrors := make(chan error)

	go config.ReceiveMessages(ctx, queueName, messagesReceived, receiveErrors)

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)

	select {
	case receivedError := <-receiveErrors:
		t.Errorf("TestRabbitmqReceiveMessage shouldn't fail, errorwas \"%s\".", receivedError.Error())
	case messageReceived := <-messagesReceived:
		stringReceived := string(messageReceived)
		if stringReceived != "This is a Test" {
			t.Errorf("Received Message shold be \"This is a Test\", not \"%s\".", stringReceived)
		}
	default:
		t.Errorf("TestRabbitmqReceiveMessage shold return a message or an error")
	}
}
