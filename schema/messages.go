package schema

import (
	"encoding/json"
	"time"
)

type EndpointMessage struct {
	Type    EndpointMessageType `json:"type"`
	Payload interface{}         `json:"payload"`
}

func (m *EndpointMessage) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type EndpointMessageType `json:"type"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}
	m.Type = typ.Type

	switch typ.Type {
	case EndpointMessageTypeACK:
		m.Payload = &EndpointMessagePayloadACK{}
	case EndpointMessageTypeINIT:
		m.Payload = &EndpointMessagePayloadINIT{}
	case EndpointMessageTypeMemoryREQ:
		m.Payload = &EndpointMessagePayloadMemoryREQ{}
	case EndpointMessageTypeMemoryRES:
		m.Payload = &EndpointMessagePayloadMemoryRES{}
	case EndpointMessageTypeInputREQ:
		m.Payload = &EndpointMessagePayloadInputREQ{}
	case EndpointMessageTypeInputRES:
		m.Payload = &EndpointMessagePayloadInputRES{}
	case EndpointMessageTypeConfigREQ:
		m.Payload = &EndpointMessagePayloadConfigREQ{}
	case EndpointMessageTypeConfigRES:
		m.Payload = &EndpointMessagePayloadConfigRES{}
	case EndpointMessageTypeDebugREQ:
		m.Payload = &EndpointMessagePayloadDebugREQ{}
	case EndpointMessageTypeDebugRES:
		m.Payload = &EndpointMessagePayloadDebugRES{}
	case EndpointMessageTypeDebugChangeREQ:
		m.Payload = &EndpointMessagePayloadDebugChangeREQ{}
	case EndpointMessageTypeDebugChangeRES:
		m.Payload = &EndpointMessagePayloadDebugChangeRES{}
	case EndpointMessageTypeDebugStepREQ:
		m.Payload = &EndpointMessagePayloadDebugStepREQ{}
	case EndpointMessageTypeDebugStepRES:
		m.Payload = &EndpointMessagePayloadDebugStepRES{}
	}

	type tmp EndpointMessage // avoids infinite recursion
	return json.Unmarshal(data, (*tmp)(m))
}

type EndpointMessageType int

const (
	EndpointMessageTypeACK            = iota
	EndpointMessageTypeINIT           = iota
	EndpointMessageTypeMemoryREQ      = iota
	EndpointMessageTypeMemoryRES      = iota
	EndpointMessageTypeInputREQ       = iota
	EndpointMessageTypeInputRES       = iota
	EndpointMessageTypeConfigREQ      = iota
	EndpointMessageTypeConfigRES      = iota
	EndpointMessageTypeDebugREQ       = iota
	EndpointMessageTypeDebugRES       = iota
	EndpointMessageTypeDebugChangeREQ = iota
	EndpointMessageTypeDebugChangeRES = iota
	EndpointMessageTypeDebugStepREQ   = iota
	EndpointMessageTypeDebugStepRES   = iota
)

type EndpointMessagePayloadACK struct{}

type EndpointMessagePayloadINIT struct {
	Name string `json:"name"`
}

type EndpointMessagePayloadMemoryREQ struct{}
type EndpointMessagePayloadMemoryRES struct {
	Memory MemoryResources `json:"memory"`
	Pool   [][]PoolElem    `json:"pool"`
}

type EndpointMessagePayloadInputREQ struct {
	Input string `json:"input"`
}
type EndpointMessagePayloadInputRES struct {
	Error string `json:"error"`
}

type EndpointMessagePayloadConfigREQ struct{}
type EndpointMessagePayloadConfigRES struct {
	Agent AgentConfiguration `json:"agent"`
}

type EndpointMessagePayloadDebugREQ struct{}
type EndpointMessagePayloadDebugRES struct {
	Paused    bool   `json:"paused"`
	Verbosity string `json:"verbosity"`
}

type EndpointMessagePayloadDebugChangeREQ struct {
	Paused    bool   `json:"paused"`
	Verbosity string `json:"verbosity"`
}
type EndpointMessagePayloadDebugChangeRES struct{}

type EndpointMessagePayloadDebugStepREQ struct{}
type EndpointMessagePayloadDebugStepRES struct{}

// MemoryResources represents the resources of an agent
type MemoryResources struct {
	Bool    map[string]bool      `json:"bool"`
	Integer map[string]int64     `json:"integer"`
	Float   map[string]float64   `json:"float"`
	Text    map[string]string    `json:"text"`
	Time    map[string]time.Time `json:"time"`
}

// PoolElem represents a pool element
type PoolElem struct {
	Resource string `json:"res"`
	Value    string `json:"val"`
}
