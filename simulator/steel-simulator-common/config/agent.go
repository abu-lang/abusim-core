package config

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Image     string
	Namespace string
	Agents    map[string]Agent
}

type Agent struct {
	Name             string
	MemoryController string
	Memory           map[string]map[string][]string
	Rules            []string
	Endpoints        []string
	Tick             time.Duration
}

func NewAgent(name string) *Agent {
	return &Agent{
		Name:             name,
		MemoryController: "basic",
		Memory:           make(map[string]map[string][]string),
		Rules:            nil,
		Endpoints:        nil,
		Tick:             time.Second,
	}
}

func (a *Agent) SetMemoryController(memorycontroller string) {
	if memorycontroller != "" {
		a.MemoryController = memorycontroller
	}
}

func (a *Agent) AddMemoryItem(item string) error {
	parts := strings.Split(item, ":")
	var memoryItem struct {
		Type   string
		Name   string
		Values []string
	}
	switch len(parts) {
	case 3:
		memoryItem.Values = strings.Split(parts[2], ",")
		fallthrough
	case 2:
		memoryItem.Type = parts[0]
		memoryItem.Name = parts[1]
	default:
		return fmt.Errorf("bad value in memory item \"%s\": unknown number of parts", item)
	}
	if _, ok := a.Memory[memoryItem.Type]; !ok {
		a.Memory[memoryItem.Type] = make(map[string][]string)
	}
	a.Memory[memoryItem.Type][memoryItem.Name] = memoryItem.Values
	return nil
}

func (a *Agent) AddRule(rule string) {
	a.Rules = append(a.Rules, rule)
}

func (a *Agent) SetEndpoints(endpoints []string) {
	a.Endpoints = endpoints
}

func (a *Agent) SetTick(tick string) {
	if tick != "" {
		tickDuration, err := time.ParseDuration(tick)
		if err != nil {
			return
		}
		a.Tick = tickDuration
	}
}

func (a *Agent) Serialize() (string, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(a)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func (a *Agent) Deserialize(str string) error {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	b.Write(by)
	err = gob.NewDecoder(&b).Decode(&a)
	if err != nil {
		return err
	}
	return nil
}
