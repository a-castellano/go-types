//go:build unit_tests || notification_unit_tests

package notification

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNotificationWithEmptyDestination(t *testing.T) {

	_, err := NewNotification("", "")

	if err == nil {
		t.Fatalf("TestNotificationWithEmptyDestination should fail.")
	}

	expectedError := "notification property cannot be empty: destination"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithEmptyDestination error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotificationWithEmptyMessage(t *testing.T) {

	_, err := NewNotification("toSomeone", "")

	if err == nil {
		t.Fatalf("TestNotificationWithEmptyMessage should fail.")
	}

	expectedError := "notification property cannot be empty: message"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithEmptyMessage error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotificationWithInvalidLevel(t *testing.T) {

	_, err := NewNotification("toSomeone", "some message", WithLevel(2000))

	if err == nil {
		t.Fatalf("TestNotificationWithInvalidLevel should fail.")
	}

	expectedError := "invalid notification level: 2000"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithInvalidLevel error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotification(t *testing.T) {

	destination := "toSomeone"
	message := "some message"
	level := Warning

	notification, err := NewNotification(destination, message, WithLevel(level))

	if err != nil {
		t.Fatalf("TestNotification should not fail with valid parameters, error was \"%s\"", err.Error())
	}

	if notification.Destination() != destination {
		t.Fatalf("notification destination should be \"%s\" instead of \"%s\"", destination, notification.Destination())
	}
	if notification.Message() != message {
		t.Fatalf("notification message should be \"%s\" instead of \"%s\"", message, notification.Message())
	}
	if notification.Level() != level {
		t.Fatalf("notification level should be \"%d\" instead of \"%d\"", level, notification.Level())
	}

}

func TestLogValue(t *testing.T) {

	destination := "toSomeone"
	message := "some message"
	level := Warning

	notification, err := NewNotification(destination, message, WithLevel(level))

	if err != nil {
		t.Fatalf("TestNotification should not fail with valid parameters, error was \"%s\"", err.Error())
	}
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	logger.Info("test log", "notification", notification)

	bufferLen := buf.Len()

	if bufferLen <= 0 {
		t.Fatalf("TestLogValue has failed, buffer is empty")
	}

	var loggedData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &loggedData); err != nil {
		t.Fatalf("TestLogValue has failed, cannot unmarshal json log")
	}

	smtpConfig := loggedData["notification"].(map[string]interface{})
	destinationValue := smtpConfig["destination"].(string)
	if destinationValue != destination {
		t.Fatalf("TestLogValue has failed, destination should be \"%s\" but it was \"%s\"", destination, destinationValue)
	}
}
