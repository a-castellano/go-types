package notification

import (
	"errors"
	"fmt"
	"log/slog"
)

// Level represents the severity of a Notification.
type Level int

// Severity levels, ordered from least to most severe. Info is the zero value,
// so it is the level any Notification built without WithLevel defaults to.
const (
	Info Level = iota
	Warning
	Error

	// levelSentinel marks the end of the valid levels: WithLevel validates
	// against it, so it must always be the last value in this block. It is
	// not a real level — do not use it as one.
	levelSentinel
)

// ErrInvalidLevel is returned when a given Level is not one of the defined
// severity levels. Callers can detect it with errors.Is.
var ErrInvalidLevel = errors.New("invalid notification level")

// ErrEmptyProperty is returned when a required Notification property is empty.
// The wrapped message names the offending property. Callers can detect it with
// errors.Is.
var ErrEmptyProperty = errors.New("notification property cannot be empty")

// Notification is a type that defines a message to be delivered to a
// destination with a severity level.
type Notification struct {
	destination string // Where the notification is delivered
	message     string // Notification content
	level       Level  // Severity level, Info unless set via WithLevel
}

// Option configures a Notification during construction. Options are applied in
// order by NewNotification and may return an error to reject invalid values.
type Option func(*Notification) error

// WithLevel returns an Option that sets the notification severity level. It
// fails with ErrInvalidLevel if level is not one of the defined levels.
func WithLevel(level Level) Option {
	return func(notification *Notification) error {
		if level < Info || level >= levelSentinel {
			return fmt.Errorf("%w: %d", ErrInvalidLevel, level)
		}
		notification.level = level
		return nil
	}
}

// NewNotification is the function that validates and returns Notification
// instance. Destination and message are required and cannot be empty; the
// severity level defaults to Info unless WithLevel is passed.
func NewNotification(destination string, message string, opts ...Option) (Notification, error) {

	if destination == "" {
		return Notification{}, fmt.Errorf("%w: destination", ErrEmptyProperty)
	}

	if message == "" {
		return Notification{}, fmt.Errorf("%w: message", ErrEmptyProperty)
	}

	notification := Notification{
		destination: destination,
		message:     message,
	}

	for _, opt := range opts {
		if err := opt(&notification); err != nil {
			return Notification{}, err
		}
	}

	return notification, nil
}

// Destination returns where the notification is delivered.
func (notification *Notification) Destination() string {
	return notification.destination
}

// Message returns the notification content.
func (notification *Notification) Message() string {
	return notification.message
}

// Level returns the notification severity level.
func (notification *Notification) Level() Level {
	return notification.level
}

// LogValue allows to log Notification omitting the message content: only the
// destination and the severity level are logged.
func (notification Notification) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("destination", notification.destination),
		slog.Int("level", int(notification.level)),
	)
}
