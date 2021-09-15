package schema

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"
)

// AgentConfiguration represents an agent config
type AgentConfiguration struct {
	Name             string                       `json:"name"`
	MemoryController string                       `json:"memorycontroller"`
	Memory           map[string]map[string]string `json:"memory"`
	Rules            []string                     `json:"rules"`
	Endpoints        []string                     `json:"endpoints"`
	Tick             time.Duration                `json:"tick"`
}

// Serialize returns the agent as a string
func (a *AgentConfiguration) Serialize() (string, error) {
	// I encode the agent as bytes...
	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(a)
	if err != nil {
		return "", err
	}
	// ... and I encode them in Base64
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// Deserialize returns the agent object from a string
func (a *AgentConfiguration) Deserialize(str string) error {
	// I decode the Base64...
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	// ... and I decode the bytes into an agent
	b := bytes.Buffer{}
	b.Write(by)
	err = json.NewDecoder(&b).Decode(&a)
	if err != nil {
		return err
	}
	return nil
}
