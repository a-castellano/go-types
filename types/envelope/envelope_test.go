//go:build unit_tests || envelope_unit_tests

package envelope

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"maps"
	"testing"
)

type book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func TestMarshal(t *testing.T) {
	body := []byte{0x00, 0xFF, 0x01}
	carrier := map[string]string{"traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}
	envelope := Envelope{Carrier: carrier, Body: body}

	_, err := envelope.Marshal()
	if err != nil {
		t.Fatalf("TestMarshal should not fail, error was '%s'", err.Error())
	}
}

func TestUnmarshalFailedData(t *testing.T) {
	data := []byte{0x00, 0xFF, 0x01}
	envelope, err := Unmarshal(data)

	if err == nil {
		t.Fatalf("TestUnmarshalFailedData should fail, encoded data has no sense")
	} else {
		if envelope != nil {
			t.Fatalf("TestUnmarshalFailedData envelope should be nil as encoded data has no sense")
		}
	}
}

func TestMarshalUnmarshal(t *testing.T) {

	formerBook := book{Title: "J'irai cracher sur vos tombes", Author: "Vernon Sullivan", Pages: 200}

	encodedFormerBook, err := json.Marshal(formerBook)
	if err != nil {
		t.Fatalf("book Marshal should not fail, error was '%s'", err.Error())
	}

	carrier := map[string]string{"traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}
	envelope := Envelope{Carrier: carrier, Body: encodedFormerBook}

	marshaledEnvelope, err := envelope.Marshal()
	if err != nil {
		t.Fatalf("TestMarshalUnmarshal should not fail when marshaledEnvelope is generated, error was '%s'", err.Error())
	}

	receivedEnvelope, err := Unmarshal(marshaledEnvelope)
	if err != nil {
		t.Fatalf("TestMarshalUnmarshal should not fail when envelope is unmashaled again, error was '%s'", err.Error())
	}

	var receivedBook book
	err = json.Unmarshal(receivedEnvelope.Body, &receivedBook)
	if err != nil {
		t.Fatalf("TestMarshalUnmarshal should not fail when reveived book is unmashaled, error was '%s'", err.Error())
	}

	if formerBook != receivedBook {
		t.Fatalf("formerBook and receivedBook should be equal")
	}

	if !maps.Equal(carrier, receivedEnvelope.Carrier) {
		t.Fatalf("carrier should not change after marshal unmarshal operation")

	}
}

func TestMarshalUnmarshalNilCarrier(t *testing.T) {

	// Telemetry disabled: the producer builds an envelope with no carrier, only the body.
	body := []byte("payload with no trace context")
	envelope := Envelope{Carrier: nil, Body: body}

	marshaledEnvelope, err := envelope.Marshal()
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalNilCarrier should not fail when marshaledEnvelope is generated, error was '%s'", err.Error())
	}

	receivedEnvelope, err := Unmarshal(marshaledEnvelope)
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalNilCarrier should not fail when envelope is unmashaled again, error was '%s'", err.Error())
	}

	if !bytes.Equal(body, receivedEnvelope.Body) {
		t.Fatalf("body should not change after marshal unmarshal operation")
	}

	if len(receivedEnvelope.Carrier) != 0 {
		t.Fatalf("carrier should stay empty after round-trip when telemetry is disabled")
	}
}

func TestMarshalString(t *testing.T) {
	body := []byte{0x00, 0xFF, 0x01}
	carrier := map[string]string{"traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}
	envelope := Envelope{Carrier: carrier, Body: body}

	_, err := envelope.MarshalString()
	if err != nil {
		t.Fatalf("TestMarshalString should not fail, error was '%s'", err.Error())
	}
}

func TestUnmarshalStringFailedData(t *testing.T) {
	data := []byte{0x00, 0xFF, 0x01}
	encodedData := base64.StdEncoding.EncodeToString(data)
	_, err := UnmarshalString(encodedData)

	if err == nil {
		t.Fatalf("TestUnmarshalStringFailedData should fail, encoded data has no sense")
	}
}

func TestUnmarshalStringFailedString(t *testing.T) {
	encodedData := "ññññññ"
	_, err := UnmarshalString(encodedData)

	if err == nil {
		t.Fatalf("TestUnmarshalStringFailedString should fail, string is not base64 compatible")
	}
}

func TestMarshalUnmarshalString(t *testing.T) {

	formerBook := book{Title: "J'irai cracher sur vos tombes", Author: "Vernon Sullivan", Pages: 200}

	encodedFormerBook, err := json.Marshal(formerBook)
	if err != nil {
		t.Fatalf("book Marshal should not fail, error was '%s'", err.Error())
	}

	carrier := map[string]string{"traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}
	envelope := Envelope{Carrier: carrier, Body: encodedFormerBook}

	marshaledEnvelopeString, err := envelope.MarshalString()
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalString should not fail when marshaledEnvelope is generated, error was '%s'", err.Error())
	}

	receivedEnvelope, err := UnmarshalString(marshaledEnvelopeString)
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalString should not fail when envelope is unmashaled again, error was '%s'", err.Error())
	}

	var receivedBook book
	err = json.Unmarshal(receivedEnvelope.Body, &receivedBook)
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalString should not fail when reveived book is unmashaled, error was '%s'", err.Error())
	}

	if formerBook != receivedBook {
		t.Fatalf("formerBook and receivedBook should be equal")
	}

	if !maps.Equal(carrier, receivedEnvelope.Carrier) {
		t.Fatalf("carrier should not change after marshal unmarshal operation")

	}
}

// TestMarshalUnmarshalStringEmptyBody covers the HTTP case, where the domain payload
// travels as the request body and the envelope carries only the trace context, so Body
// is left empty. This test was written by an AI agent (Claude).
func TestMarshalUnmarshalStringEmptyBody(t *testing.T) {

	carrier := map[string]string{"traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}
	envelope := Envelope{Carrier: carrier}

	marshaledEnvelopeString, err := envelope.MarshalString()
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalStringEmptyBody should not fail when marshaledEnvelopeString is generated, error was '%s'", err.Error())
	}

	receivedEnvelope, err := UnmarshalString(marshaledEnvelopeString)
	if err != nil {
		t.Fatalf("TestMarshalUnmarshalStringEmptyBody should not fail when envelope is unmarshaled again, error was '%s'", err.Error())
	}

	if len(receivedEnvelope.Body) != 0 {
		t.Fatalf("body should stay empty after marshal unmarshal operation")
	}

	if !maps.Equal(carrier, receivedEnvelope.Carrier) {
		t.Fatalf("carrier should not change after marshal unmarshal operation")
	}
}
