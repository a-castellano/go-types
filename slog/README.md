# slog type

This type manages [slog](https://pkg.go.dev/log/slog) logger config.

## Required variables

The following env variables should be defined when using this type:

- **APP_NAME** defines the application name, it must be a non-empty string, it has no default value and must be set
- **SLOG_LEVEL** defines the default log level, it must be one of "Debug", "Info", "Warn" or "Error", its default value is "Info"
- **SLOG_FORMAT** defines the log format, it must be either "JSON" or "plain", its default value is "JSON"
- **SLOG_ADD_SOURCE** defines whether file:line is added to logs, it must be either "true" or "false", its default value is "true"

## About the output

For the time being I do not need to send logs to stderr or files, my apps will be running wrapped by systemd units that will manage final output.
