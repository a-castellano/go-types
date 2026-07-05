# SMTP type

This type manages SMTP configs, defining the connection to a SMTP server used for sending email.

## Required variables

The following env variables must be defined when using this type, none of them has a default value:

* **SMTP_FROM** defines the sender email username (without domain), it must be a string
* **SMTP_DOMAIN** defines the sender email domain, it must be a string
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
// Output: from=notifications domain=example.com host=smtp.example.com port=587 user=mailer password=***** validate_tls=true
```
