package envelope

import (
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
