package api

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abu-lang/abusim-core/schema"
)

// ActionType represents a type of valid Action
type ActionType int

const (
	ActionConfig    ActionType = iota
	ActionMemory    ActionType = iota
	ActionInput     ActionType = iota
	ActionDebugInfo ActionType = iota
	ActionDebugSet  ActionType = iota
	ActionDebugStep ActionType = iota
)

// Action represents an action that the API performs
type Action struct {
	Type    ActionType
	Payload interface{}
}

// ActionResponse represents a response to an Action
type ActionResponse struct {
	Error      bool
	StatusCode int
	Payload    interface{}
}

// Process waits for an Action, performs it and publish an ActionResponse
func Process(actions chan Action, responses chan ActionResponse, ends map[string]*schema.Endpoint) {
	// Forever...
	for {
		// ... I get an Action...
		action := <-actions
		// ... I execute the correct procedure based on its type and I publish the response
		switch action.Type {
		case ActionConfig:
			responses <- doConfigGet(action, ends)
		case ActionMemory:
			responses <- doMemoryGet(action, ends)
		case ActionInput:
			responses <- doInput(action, ends)
		case ActionDebugInfo:
			responses <- doDebugGet(action, ends)
		case ActionDebugSet:
			responses <- doDebugSet(action, ends)
		case ActionDebugStep:
			responses <- doDebugStep(action, ends)
		}
	}
}

func doConfigGet(action Action, ends map[string]*schema.Endpoint) ActionResponse {
	// I get the agent name...
	agentName := action.Payload.(string)
	// ... I send a configuration request...
	err := sendMessageByName(agentName, ends, &schema.EndpointMessage{
		Type:    schema.EndpointMessageTypeConfigREQ,
		Payload: struct{}{},
	})
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// ... and I receive the answer
	msg, err := receiveMessageByName(agentName, ends)
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	if msg.Type != schema.EndpointMessageTypeACK {
		err := errors.New("unexpected response")
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// I get the agent configuration from the answer...
	agent := msg.Payload.(schema.Agent)
	// ... I prepare the memory item strings...
	memory := []string{}
	for vartype, value := range agent.Memory {
		for name, initvalue := range value {
			memory = append(memory, strings.Join([]string{vartype, name, initvalue}, ":"))
		}
	}
	// ... and I respond with the configuration
	return ActionResponse{
		Error:      false,
		StatusCode: http.StatusOK,
		Payload: struct {
			Name             string   `json:"name"`
			MemoryController string   `json:"memorycontroller"`
			Memory           []string `json:"memory"`
			Rules            []string `json:"rules"`
			Endpoints        []string `json:"endpoints"`
			Tick             string   `json:"tick"`
		}{
			Name:             agent.Name,
			MemoryController: agent.MemoryController,
			Memory:           memory,
			Rules:            agent.Rules,
			Endpoints:        agent.Endpoints,
			Tick:             agent.Tick.String(),
		},
	}
}

func doMemoryGet(action Action, ends map[string]*schema.Endpoint) ActionResponse {
	// I get the agent name...
	agentName := action.Payload.(string)
	// ... I send a memory request...
	err := sendMessageByName(agentName, ends, &schema.EndpointMessage{
		Type:    schema.EndpointMessageTypeMemoryREQ,
		Payload: struct{}{},
	})
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// ... and I receive the answer
	msg, err := receiveMessageByName(agentName, ends)
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	if msg.Type != schema.EndpointMessageTypeACK {
		err := errors.New("unexpected response")
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// I get the state from the answer...
	state := msg.Payload.(schema.AgentState)
	// ... I prepare the memory...
	type mem struct {
		Bool    map[string]bool      `json:"bool"`
		Integer map[string]int64     `json:"integer"`
		Float   map[string]float64   `json:"float"`
		Text    map[string]string    `json:"text"`
		Time    map[string]time.Time `json:"time"`
	}
	m := mem{
		Bool:    state.Memory.Bool,
		Integer: state.Memory.Integer,
		Float:   state.Memory.Float,
		Text:    state.Memory.Text,
		Time:    state.Memory.Time,
	}
	// ... I prepare the pool...
	type poolElem struct {
		Resource string `json:"resource"`
		Value    string `json:"value"`
	}
	p := [][]poolElem{}
	for _, ruleActions := range state.Pool {
		poolActions := []poolElem{}
		for _, action := range ruleActions {
			poolActions = append(poolActions, poolElem(action))
		}
		p = append(p, poolActions)
	}
	// ... and I respond with the agent state
	return ActionResponse{
		Error:      false,
		StatusCode: http.StatusOK,
		Payload: struct {
			Name   string       `json:"name"`
			Memory mem          `json:"memory"`
			Pool   [][]poolElem `json:"pool"`
		}{
			Name:   agentName,
			Memory: m,
			Pool:   p,
		},
	}
}

func doInput(action Action, ends map[string]*schema.Endpoint) ActionResponse {
	// I get the payload...
	payload := action.Payload.(struct {
		agentName string
		actions   string
	})
	// ... I send an input request...
	err := sendMessageByName(payload.agentName, ends, &schema.EndpointMessage{
		Type:    schema.EndpointMessageTypeInputREQ,
		Payload: payload.actions,
	})
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// ... and I receive the answer
	msg, err := receiveMessageByName(payload.agentName, ends)
	if err != nil || msg.Type != schema.EndpointMessageTypeACK {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusInternalServerError,
			Payload:    err.Error(),
		}
	}
	// I get the eventual error from the answer...
	errInput := msg.Payload.(string)
	if errInput != "" {
		log.Println(errInput)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusBadRequest,
			Payload:    errInput,
		}
	}
	// and, if there is none, I respond affirmatively
	return ActionResponse{
		Error:      false,
		StatusCode: http.StatusOK,
		Payload: struct {
			Result string `json:"result"`
		}{
			Result: "ok",
		},
	}
}

func doDebugGet(action Action, ends map[string]*schema.Endpoint) ActionResponse {
	// I get the agent name...
	agentName := action.Payload.(string)
	// ... I send a debug request...
	err := sendMessageByName(agentName, ends, &schema.EndpointMessage{
		Type:    schema.EndpointMessageTypeDebugREQ,
		Payload: struct{}{},
	})
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// ... and I receive the answer
	msg, err := receiveMessageByName(agentName, ends)
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	if msg.Type != schema.EndpointMessageTypeACK {
		err := errors.New("unexpected response")
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// I get the state from the answer...
	dbgStatus := msg.Payload.(schema.AgentDebugStatus)
	// ... I prepare the status...
	type status struct {
		Paused    bool   `json:"paused"`
		Verbosity string `json:"verbosity"`
	}
	s := status{
		Paused:    dbgStatus.Paused,
		Verbosity: dbgStatus.Verbosity,
	}
	// ... and I respond with the agent debug status
	return ActionResponse{
		Error:      false,
		StatusCode: http.StatusOK,
		Payload: struct {
			Name   string `json:"name"`
			Status status `json:"status"`
		}{
			Name:   agentName,
			Status: s,
		},
	}
}

func doDebugSet(action Action, ends map[string]*schema.Endpoint) ActionResponse {
	// I get the payload...
	payload := action.Payload.(struct {
		agentName string
		paused    bool
		verbosity string
	})
	// ... I send an debug status change request...
	err := sendMessageByName(payload.agentName, ends, &schema.EndpointMessage{
		Type: schema.EndpointMessageTypeDebugChangeREQ,
		Payload: schema.AgentDebugStatus{
			Paused:    payload.paused,
			Verbosity: payload.verbosity,
		},
	})
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// ... and I receive the answer
	msg, err := receiveMessageByName(payload.agentName, ends)
	if err != nil || msg.Type != schema.EndpointMessageTypeACK {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusInternalServerError,
			Payload:    err.Error(),
		}
	}
	// Finally, I respond affirmatively
	return ActionResponse{
		Error:      false,
		StatusCode: http.StatusOK,
		Payload: struct {
			Result string `json:"result"`
		}{
			Result: "ok",
		},
	}
}

func doDebugStep(action Action, ends map[string]*schema.Endpoint) ActionResponse {
	// I get the agent name...
	agentName := action.Payload.(string)
	// ... I send a configuration request...
	err := sendMessageByName(agentName, ends, &schema.EndpointMessage{
		Type:    schema.EndpointMessageTypeDebugStepREQ,
		Payload: struct{}{},
	})
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// ... and I receive the answer
	msg, err := receiveMessageByName(agentName, ends)
	if err != nil {
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	if msg.Type != schema.EndpointMessageTypeACK {
		err := errors.New("unexpected response")
		log.Println(err)
		return ActionResponse{
			Error:      true,
			StatusCode: http.StatusNotFound,
			Payload:    err.Error(),
		}
	}
	// Finally, I respond affirmatively
	return ActionResponse{
		Error:      false,
		StatusCode: http.StatusOK,
		Payload: struct {
			Result string `json:"result"`
		}{
			Result: "ok",
		},
	}
}
