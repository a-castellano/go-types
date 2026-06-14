# slog type

This type manages [slog](https://pkg.go.dev/log/slog) logger config.

## Required variables

The following env variables should be defined when using this type:

- **APP_NAME** defines the application name, it must be a non-empty string, it has no default value and must be set
- **SLOG_LEVEL** defines the default log level, it must be one of "Debug", "Info", "Warn" or "Error", its default value is "Info"
- **SLOG_OUTPUT** defines the default log output, by default it is set to "stdout", but it can be set to a file path, or sterr. The service that uses this config will implemente the output managing logic.
- **SLOG_FORMAT** defines the log format, it must be either "JSON" or "plain", its default value is "JSON"
- **SLOG_ADD_SOURCE** defines whether file:line is added to logs, it must be either "true" or "false", its default value is "true"
