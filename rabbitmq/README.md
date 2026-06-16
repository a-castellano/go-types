# RabbitMQ type

This type manages RabbitMQ config

## Required variables

The following env variables should be defined when using this type:

* **RABBITMQ_HOST** defines RabbitMQ host, it must be a string, its default value is "localhost" 
* **RABBITMQ_PORT** defines RabbitMQ port, it must be a valid port number, its default value is 5672 
* **RABBITMQ_USER** defines RabbitMQ user, its default value is "guest" 
* **RABBITMQ_PASSWORD** defines RabbitMQ user's password, its default value is "guest" 

## Logging

`Config` implements `slog.LogValuer`, so it can be passed directly to any `log/slog` logger. The password field is automatically redacted.

```go
config, err := rabbitmq.NewConfig()
if err != nil {
    log.Fatal(err)
}

slog.Info("RabbitMQ config loaded", "rabbitmq", config)
// Output: rabbitmq_host=localhost rabbitmq_port=5672 rabbitmq_user=guest rabbitmq_password=REDACTED
```
