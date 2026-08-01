package envelope

import (
	"encoding/base64"
	"encoding/json"
)

// Envelope is a type that carries a serialized domain payload together with the
// trace propagation context, so a distributed trace stays connected across
// service boundaries regardless of the transport it travels over.
type Envelope struct {
	Carrier map[string]string `json:"carrier"` // Propagation carrier filled by the OpenTelemetry propagator: traceparent, tracestate and baggage
	Body    []byte            `json:"body"`    // Serialized domain payload; opaque bytes as far as the envelope is concerned
}

// Marshal serializes the Envelope into its wire representation so it can be handed
// to any transport as a single []byte.
func (envelope *Envelope) Marshal() ([]byte, error) {
	return json.Marshal(envelope)
}

// Unmarshal is the function that validates and returns an Envelope instance from
// its wire representation, the inverse of Marshal.
func Unmarshal(data []byte) (*Envelope, error) {
	var envelope Envelope
	err := json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}
	return &envelope, err
}

// MarshalString serializes the Envelope and returns it base64-encoded, so it can be
// handed to a transport that moves text rather than bytes and needs a single safe
// value, such as an HTTP header.
//
// The alphabet is standard base64 with padding (encoding/base64.StdEncoding). It is
// part of the wire contract: producer and consumer must agree on it, so it cannot be
// swapped without breaking every service already speaking this protocol.
//
// The name is deliberately not MarshalText. Implementing encoding.TextMarshaler would
// change what json.Marshal produces for an Envelope, and since Marshal is itself
// json.Marshal, the two would call each other until the stack ran out.
func (envelope *Envelope) MarshalString() (string, error) {

	marshaledData, err := envelope.Marshal()

	return base64.StdEncoding.EncodeToString(marshaledData), err
}

// UnmarshalString validates and returns an Envelope instance from its base64 string
// representation, the inverse of MarshalString.
//
// The input must use the same standard base64 alphabet with padding that MarshalString
// produces; anything else, and anything that decodes to bytes which are not a valid
// envelope, is reported as an error.
func UnmarshalString(encodedString string) (*Envelope, error) {
	var emptyEnvelope Envelope
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return &emptyEnvelope, err
	}

	envelope, unmarshalError := Unmarshal(decodedBytes)

	if unmarshalError != nil {
		return &emptyEnvelope, unmarshalError
	}

	return envelope, nil
}
