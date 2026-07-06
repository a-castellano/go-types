# SMTP type

This type manages SMTP configs, defining the connection to a SMTP server used for sending email.

## Exported API

The package exposes a `Config` struct created through `NewConfig()`:

```go
config, err := smtp.NewConfig()
```

`NewConfig()` reads the environment variables listed below, validates them and returns a `*Config` (or an error describing the first missing or invalid variable).

`Config` fields are private and immutable after construction; values are read through getters:

| Method | Returns | Description |
|---|---|---|
| `From()` | `string` | Sender email address, from `SMTP_FROM` |
| `Address()` | `string` | SMTP server address in `host:port` form, built from `SMTP_HOST` and `SMTP_PORT` |
| `Username()` | `string` | SMTP authentication username |
| `Password()` | `string` | SMTP authentication password |
| `ValidateTLS()` | `bool` | Whether TLS certificates are validated |

## Required variables

The following env variables must be defined when using this type, none of them has a default value:

* **SMTP_FROM** defines the sender email address, it must be a valid email address
* **SMTP_HOST** defines SMTP server hostname, it must be a string
* **SMTP_PORT** defines SMTP server port, it must be a valid port number between 1 and 65535
* **SMTP_USERNAME** defines SMTP authentication username
* **SMTP_PASSWORD** defines SMTP authentication password

## Optional variables

* **SMTP_VALIDATE_TLS** defines whether TLS certificates are validated, its default value is "true". Any value other than "true" disables validation.

## Logging

`Config` implements `slog.LogValuer`, so it can be passed directly to any `log/slog` logger. The password field is always masked.

```go
config, err := smtp.NewConfig()
if err != nil {
    log.Fatal(err)
}

slog.Info("SMTP config loaded", "smtp", config)
// Output: smtp.from=notifications@example.com smtp.address=smtp.example.com:587 smtp.user=mailer smtp.password=***** smtp.validate_tls=true
```
