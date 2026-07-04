# Envelope type

This type carries what travels between services on the wire: the domain payload plus the trace context needed to keep a distributed trace connected across service boundaries.

Unlike the other types in this repo, `Envelope` is not a configuration read from the environment. It is a data type that is serialized into the message payload, sent over a transport (RabbitMQ today, HTTP or Kafka tomorrow), and deserialized on the other side.

## What it is for

Distributed tracing needs the trace context to travel from the producer to the consumer, so the consumer continues the same trace instead of starting a new, disconnected one. This type is the shared contract for doing that in a transport-agnostic way.

An `Envelope` has two parts:

- `Carrier`: a `map[string]string` filled by the OpenTelemetry propagator. It holds the W3C propagation fields (`traceparent`, `tracestate`, `baggage`) as plain string key-values.
- `Body`: the serialized domain payload, opaque bytes as far as the envelope is concerned.

The whole envelope (carrier and body) is serialized as a single payload and put on the wire. Carrying the trace context inside the body, rather than in transport headers, means the same mechanism works over any transport with no per-transport carrier adapters. The distributed trace still forms correctly: the parent/child link is established when the consumer extracts the carrier and starts its span.

## Relationship with the OpenTelemetry SDK

This type is deliberately dependency-free (standard library only). It knows nothing about OpenTelemetry: the `Carrier` is just a `map[string]string`. Filling that map from the active context (`Inject`) and rebuilding a context from it (`Extract`) is the job of the propagation helpers in `go-services`, which are the ones that depend on the OpenTelemetry SDK. The wire contract lives here; the SDK wiring lives there.

## Usage

On the producer side, build the envelope with the injected carrier and the serialized domain payload, then `Marshal` it and hand the bytes to the transport. On the consumer side, `Unmarshal` the bytes back into an `Envelope`, extract the carrier to continue the trace, and deserialize the body into the domain type.

```go
// Producer: carrier comes from the propagator, payload is the serialized domain data
data, err := (&envelope.Envelope{Carrier: carrier, Body: payload}).Marshal()
if err != nil {
    log.Fatal(err)
}
// hand data to the transport

// Consumer
received, err := envelope.Unmarshal(data)
if err != nil {
    log.Fatal(err)
}
// use received.Carrier to continue the trace, received.Body as the domain payload
```

The wire format is JSON, kept as an implementation detail behind `Marshal` and `Unmarshal`, so it can change without touching callers. Because `Body` is `[]byte`, it travels base64-encoded inside that JSON; this is expected and does not affect the domain payload once it is deserialized on the consumer side.
