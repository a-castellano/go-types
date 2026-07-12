//go:build unit_tests || notification_unit_tests

package notification

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNotificationWithEmptyDestination(t *testing.T) {

	_, err := NewNotification("", "", "")

	if err == nil {
		t.Fatalf("TestNotificationWithEmptyDestination should fail.")
	}

	expectedError := "notification property cannot be empty: destination"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithEmptyDestination error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotificationWithEmptyTitle(t *testing.T) {

	_, err := NewNotification("toSomeone", "", "")

	if err == nil {
		t.Fatalf("TestNotificationWithEmptyTitle should fail.")
	}

	expectedError := "notification property cannot be empty: title"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithEmptyTitle error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotificationWithEmptyMessage(t *testing.T) {

	_, err := NewNotification("toSomeone", "title", "")

	if err == nil {
		t.Fatalf("TestNotificationWithEmptyMessage should fail.")
	}

	expectedError := "notification property cannot be empty: message"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithEmptyMessage error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotificationWithInvalidLevel(t *testing.T) {

	_, err := NewNotification("toSomeone", "title", "message", WithLevel(2000))

	if err == nil {
		t.Fatalf("TestNotificationWithInvalidLevel should fail.")
	}

	expectedError := "invalid notification level: unknown(2000)"

	if err.Error() != expectedError {
		t.Fatalf("TestNotificationWithInvalidLevel error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestNotification(t *testing.T) {

	destination := "toSomeone"
	title := "Title"
	message := "some message"
	level := Warning

	notification, err := NewNotification(destination, title, message, WithLevel(level))

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
	title := "Title"
	message := "some message"
	level := Warning

	notification, err := NewNotification(destination, title, message, WithLevel(level))

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

	loggedNotification := loggedData["notification"].(map[string]interface{})
	destinationValue := loggedNotification["destination"].(string)
	if destinationValue != destination {
		t.Fatalf("TestLogValue has failed, destination should be \"%s\" but it was \"%s\"", destination, destinationValue)
	}
	levelValue := loggedNotification["level"].(string)
	if levelValue != level.String() {
		t.Fatalf("TestLogValue has failed, level should be \"%s\" but it was \"%s\"", level.String(), levelValue)
	}

}

func TestLevelStringer(t *testing.T) {
	if Info.String() != "info" {
		t.Fatalf("Info level string should be \"info\" but it was \"%s\"", Info.String())
	}
	if Warning.String() != "warning" {
		t.Fatalf("Warning level string should be \"warning\" but it was \"%s\"", Warning.String())
	}
	if Error.String() != "error" {
		t.Fatalf("Error level string should be \"error\" but it was \"%s\"", Error.String())
	}

}
