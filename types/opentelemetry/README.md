# OpenTelemetry type

This type manages OpenTelemetry configs.

Unlike the other config types, this one mostly validates the environment rather than reading values: the service identity comes from `APP_NAME` (the same variable `slog` already requires), and every other standard `OTEL_*` variable is left to the OpenTelemetry SDK to read on its own. The SDK never derives the service name from `APP_NAME`; building the Resource attribute is the SDK startup code's job, not this type's.

## Required variables

The following env variable must be defined when using this type:

- `APP_NAME` defines the application name, it must be a non-empty string, it has no default value and must be set. Its value is the single source for the telemetry `service.name`.

## Forbidden variables

These variables must not be defined; `NewConfig` returns an error if any of them is present:

- `OTEL_SERVICE_NAME` is rejected because the service name comes from `APP_NAME`, so the logger and the traces never diverge.
- `OTEL_RESOURCE_ATTRIBUTES` is rejected for the time being; its values will be managed by this type if and when they are required.

## Everything else

Any other `OTEL_*` variable (exporter, endpoint, protocol, sampler, ...) is intentionally not handled here. The SDK reads them directly from the environment, so there is no point duplicating them in this type.
