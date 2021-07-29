package config

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"strings"
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

// NewAgent creates a new agent with the specified name and some default values
func NewAgent(name string) *Agent {
	// I return the agent
	return &Agent{
		Name:             name,
		MemoryController: "basic",
		Memory:           make(map[string]map[string]string),
		Rules:            []string{},
		Endpoints:        []string{},
		Tick:             time.Second,
	}
}

// AddMemoryItem parses a memory item string and adds it to the agent
func (a *Agent) AddMemoryItem(item string) error {
	// I split the item string into at most 3 pieces...
	parts := strings.SplitN(item, ":", 3)
	// ... I create a struct to hold the parsing result...
	memoryItem := struct {
		Type  string
		Name  string
		Value string
	}{
		Value: "",
	}
	// ... and I check whether I hava an initialization value
	switch len(parts) {
	case 3:
		// If I have it I assign it...
		memoryItem.Value = parts[2]
		fallthrough
	case 2:
		// ... and I also assign the remaining parts
		memoryItem.Type = parts[0]
		memoryItem.Name = parts[1]
	default:
		// If the item string did not contain any colon, I raise an error
		return fmt.Errorf("bad value in memory item \"%s\": unknown number of parts", item)
	}
	// I create the type map if necessary...
	if _, ok := a.Memory[memoryItem.Type]; !ok {
		a.Memory[memoryItem.Type] = make(map[string]string)
	}
	// ... and I add the item to the memory
	a.Memory[memoryItem.Type][memoryItem.Name] = memoryItem.Value
	return nil
}

// SetMemoryController sets a memory controller for the agent
func (a *Agent) SetMemoryController(controller string) {
	// If the value is set...
	if controller != "" {
		// ... I assign it
		a.MemoryController = controller
	}
}

// SetTick sets a tick for the agent
func (a *Agent) SetTick(tick string) error {
	// If the value is set...
	if tick != "" {
		// ... I parse the duration...
		tickDuration, err := time.ParseDuration(tick)
		if err != nil {
			return err
		}
		// ... and I assign it
		a.Tick = tickDuration
	}
	return nil
}

// Serialize returns the agent as a string
func (a *Agent) Serialize() (string, error) {
	// I encode the agent as bytes...
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(a)
	if err != nil {
		return "", err
	}
	// ... and I encode them in Base64
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
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
