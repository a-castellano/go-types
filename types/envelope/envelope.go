package envelope

import (
	"encoding/json"
)

type Envelope struct {
	Carrier map[string]string `json:"carrier"`
	Body    []byte            `json:"body"`
}

func (envelope *Envelope) Marshal() ([]byte, error) {

	return json.Marshal(envelope)
}

func Unmarshal(data []byte) (*Envelope, error) {
	var envelope Envelope
	err := json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}
	return &envelope, err
}
