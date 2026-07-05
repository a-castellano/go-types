package notification

import (
	"errors"
	"fmt"
)

type Level int

const (
	Info Level = iota
	Warning
	Error

	levelSentinel
)

var ErrInvalidLevel = errors.New("invalid notification level")
var ErrEmptyProperty = errors.New("notification property cannot be empty")

type Notification struct {
	destination string
	message     string
	level       Level
}

type Option func(*Notification) error

func WithLevel(level Level) Option {
	return func(notification *Notification) error {
		if level < Info || level >= levelSentinel {
			return fmt.Errorf("%w: %d", ErrInvalidLevel, level)
		}
		notification.level = level
		return nil
	}
}

func NewNotification(destination string, message string, opts ...Option) (Notification, error) {

	if destination == "" {
		return Notification{}, fmt.Errorf("%w: %s", ErrEmptyProperty, "destination")
	}

	if message == "" {
		return Notification{}, fmt.Errorf("%w: %s", ErrEmptyProperty, "message")
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

func (notification *Notification) Destination() string {
	return notification.destination
}

func (notification *Notification) Message() string {
	return notification.message
}

func (notification *Notification) Level() Level {
	return notification.level
}
