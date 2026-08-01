# Envelope type

This type carries what travels between services on the wire: the domain payload plus the trace context needed to keep a distributed trace connected across service boundaries.

Unlike the other types in this repo, `Envelope` is not a configuration read from the environment. It is a data type that is serialized into the message payload, sent over a transport (RabbitMQ today, HTTP or Kafka tomorrow), and deserialized on the other side.

## What it is for

Distributed tracing needs the trace context to travel from the producer to the consumer, so the consumer continues the same trace instead of starting a new, disconnected one. This type is the shared contract for doing that in a transport-agnostic way.

An `Envelope` has two parts:

- `Carrier`: a `map[string]string` filled by the OpenTelemetry propagator. It holds the W3C propagation fields (`traceparent`, `tracestate`, `baggage`) as plain string key-values.
- `Body`: the serialized domain payload, opaque bytes as far as the envelope is concerned.

The whole envelope (carrier and body) is serialized as a single value and put on the wire. Keeping the carrier inside the envelope, rather than spreading it across transport-native headers, means the same mechanism works over any transport with no per-transport carrier adapters. The distributed trace still forms correctly: the parent/child link is established when the consumer extracts the carrier and starts its span.

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

## The string form

Some transports do not move bytes, they move text. An HTTP header value is the case that motivated this: the envelope travels in a dedicated header so the JSON contract of the endpoint stays untouched, and a header value has to be a text-safe string.

`MarshalString` and `UnmarshalString` are the symmetric pair for that, mirroring `Marshal` and `Unmarshal`:

```go
// Producer: same envelope, encoded as a single text-safe value
value, err := (&envelope.Envelope{Carrier: carrier, Body: payload}).MarshalString()
if err != nil {
    log.Fatal(err)
}
// hand value to the transport, for instance as an HTTP header value

// Consumer
received, err := envelope.UnmarshalString(value)
if err != nil {
    log.Fatal(err)
}
```

Which form to use: `Marshal` for a transport that carries bytes, such as a RabbitMQ message body; `MarshalString` for one that needs a text value, such as an HTTP header.

### The encoding is part of the contract

The string form is standard base64 with padding (`encoding/base64.StdEncoding`) over the same JSON representation. Unlike the JSON format, which is an implementation detail behind the functions, the alphabet is not free to change: producer and consumer are different services, and they must agree on it. Two encodings coexisting across services is the failure this pair exists to prevent, so decoders must not be written by hand against a different alphabet.

Standard base64 emits `+`, `/` and `=`. All three are valid inside an HTTP header value, which is what this form targets. None of them is safe in a URL query string or a filename without further escaping, so a caller putting the value there is responsible for escaping it.

### Where the boundary is

`MarshalString` returns a `string` and knows nothing about HTTP: no `net/http` import, no header names, no transport semantics. Putting the value into a request and reading it back out is the job of `go-services` or of the application. `go-types` stays a types library.

### An empty body needs no special case

When the payload travels as the HTTP body, the envelope carries only the carrier and `Body` is left empty. That needs no special handling: `nil` bytes round-trip as `{"carrier":{...},"body":null}` and come back as an empty `Body`.
