package schema

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"time"
)

// Agent represents an agent config
type Agent struct {
	Name             string
	MemoryController string
	Memory           map[string]map[string]string
	Rules            []string
	Endpoints        []string
	Tick             time.Duration
}

// Deserialize returns the agent object from a string
func (a *Agent) Deserialize(str string) error {
	// I decode the Base64...
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	// ... and I decode the bytes into an agent
	b := bytes.Buffer{}
	b.Write(by)
	err = gob.NewDecoder(&b).Decode(&a)
	if err != nil {
		return err
	}
	return nil
}
