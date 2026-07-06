# RabbitMQ type

This type manages RabbitMQ configs, defining the connection to a RabbitMQ server.

## Exported API

The package exposes a `Config` struct created through `NewConfig()`:

```go
config, err := rabbitmq.NewConfig()
```

`NewConfig()` reads the environment variables listed below, validates them and returns a `*Config` (or an error if the port is not a valid port number).

Connection details (host, port, user and password) are kept private. The only exported field is the final product built from them:

| Field              | Type     | Description                                                                              |
| ------------------ | -------- | ---------------------------------------------------------------------------------------- |
| `ConnectionString` | `string` | AMQP URI ready to be passed to an AMQP client, e.g. `amqp://guest:guest@localhost:5672/` |

The URI is assembled with `net/url`, so credentials containing URL-reserved characters (`@`, `/`, `:`, ...) are percent-encoded and always produce a valid URI.

## Environment variables

All variables are optional; each one has a default value:

- **RABBITMQ_HOST** defines RabbitMQ host, it must be a string, its default value is "localhost"
- **RABBITMQ_PORT** defines RabbitMQ port, it must be a valid port number between 1 and 65535, its default value is 5672
- **RABBITMQ_USER** defines RabbitMQ user, its default value is "guest"
- **RABBITMQ_PASSWORD** defines RabbitMQ user's password, its default value is "guest"

## Logging

`Config` implements `slog.LogValuer`, so it can be passed directly to any `log/slog` logger. It logs the connection URL with the password masked as `xxxxx`.

```go
config, err := rabbitmq.NewConfig()
if err != nil {
    log.Fatal(err)
}

slog.Info("RabbitMQ config loaded", "rabbitmq", config)
// Output: rabbitmq.url=amqp://guest:xxxxx@localhost:5672/
```
