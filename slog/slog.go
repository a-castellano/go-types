package slog

import (
	"cmp"
	"errors"
	"fmt"
	formerslog "log/slog"
	"os"
)

// Config is a type that defines required data for defining slog parameters
type Config struct {
	DefaultLevel formerslog.Level // Specifies log default level, Info for example
	Format       string           // Log format JSON or plain
	Output       string           // Log output, default is stdout, the service will manage the output if a file is set instead of strdout or sterr
	AddSource    bool             // Adds file:line to logs
	AppName      string           // The name of the app
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := Config{}

	// Get log level from SLOG_LEVEL `env`` variable
	defaultLevel := cmp.Or(os.Getenv("SLOG_LEVEL"), "Info")
	// define config.DefaultLevel from defaultLevel value
	switch defaultLevel {
	case "Debug":
		config.DefaultLevel = formerslog.LevelDebug
	case "Info":
		config.DefaultLevel = formerslog.LevelInfo
	case "Warn":
		config.DefaultLevel = formerslog.LevelWarn
	case "Error":
		config.DefaultLevel = formerslog.LevelError
	default:
		return nil, fmt.Errorf("log level defined by `SLOG_LEVEL` variable only accepts the following values: \"Debug\", \"Info\", \"Warn\" or \"Error\". \"%s\" is not a valid value.", defaultLevel)
	}

	// Get format from SLOG_FORMAT env variable, default value is JSON
	config.Format = cmp.Or(os.Getenv("SLOG_FORMAT"), "JSON")
	if config.Format != "JSON" && config.Format != "plain" {
		return nil, fmt.Errorf("log format defined by `SLOG_FORMAT` variable only accepts the values \"JSON\" or \"plain\", \"%s\" is not a valid value", config.Format)
	}

	// Get AddSource value from SLOG_ADD_SOURCE `env` variable
	addSource := cmp.Or(os.Getenv("SLOG_ADD_SOURCE"), "true")
	if addSource != "false" && addSource != "true" {
		return nil, fmt.Errorf("add log source property defined by `SLOG_ADD_SOURCE` variable only accepts the values \"true\" or \"false\", \"%s\" is not a valid value", addSource)
	}
	config.AddSource = addSource == "true"

	config.AppName = cmp.Or(os.Getenv("APP_NAME"), "")
	if config.AppName == "" {
		return nil, errors.New("env variable \"APP_NAME\" must be defined and have a value")
	}
	config.Output = cmp.Or(os.Getenv("SLOG_OUTPUT"), "stdout")

	return &config, nil
}
