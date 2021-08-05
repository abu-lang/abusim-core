package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"steel-simulator-common/communication"
	"steel-simulator-common/config"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Serve serves the API on the API port
func Serve(ends map[string]*communication.Endpoint) {
	// I create a router for the API and I set the handlers...
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/config/{agentName}", GetHandleConfig(ends)).Methods(http.MethodGet)
	router.HandleFunc("/memory/{agentName}", GetHandleMemory(ends)).Methods(http.MethodGet, http.MethodPost)
	// ... I set up the CORS middleware...
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{"POST", "GET"},
		AllowedHeaders: []string{"Accept", "content-type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})
	// ... and I serve the CORS decorated API
	log.Fatal(http.ListenAndServe(":4000", c.Handler(router)))
}

// HandleIndex handles the queries to the main page
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	// I set the content type...
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// ... I write the I'm a Teapot status code :) ...
	w.WriteHeader(http.StatusTeapot)
	// ... and I write the a welcome message
	b, _ := json.Marshal(struct {
		Welcome string `json:"welcome"`
	}{
		Welcome: "Welcome to steel-coordinator API!",
	})
	w.Write(b)
}

// GetHandleConfig returns an handler for the configuration method
func GetHandleConfig(ends map[string]*communication.Endpoint) http.HandlerFunc {
	// I return the handler, decorated with the list of endpoints
	return func(w http.ResponseWriter, r *http.Request) {
		// I get the agent name from the query...
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		// ... I send a configuration request...
		err := sendMessageByName(agentName, ends, &communication.EndpointMessage{
			Type:    communication.EndpointMessageTypeConfigREQ,
			Payload: struct{}{},
		})
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		// ... and I receive the answer
		msg, err := receiveMessageByName(agentName, ends)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if msg.Type != communication.EndpointMessageTypeConfigRES {
			err := errors.New("unexpected response")
			log.Println(err)
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		// I get the agent configuration from the answer...
		agent := msg.Payload.(config.Agent)
		// ... I prepare the memory item strings...
		memory := []string{}
		for vartype, value := range agent.Memory {
			for name, initvalue := range value {
				memory = append(memory, strings.Join([]string{vartype, name, initvalue}, ":"))
			}
		}
		// ... and I respond with the configuration
		writeResponse(w, http.StatusOK, struct {
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
		})

	}
}

// GetHandleMemory returns an handler for the memory method
func GetHandleMemory(ends map[string]*communication.Endpoint) http.HandlerFunc {
	// I return the handler, decorated with the list of endpoints
	return func(w http.ResponseWriter, r *http.Request) {
		// I get the agent name from the query...
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		// ... and I check what do I have to do
		switch r.Method {
		// If I need to retrieve the memory...
		case http.MethodGet:
			// ... I send a memory request...
			err := sendMessageByName(agentName, ends, &communication.EndpointMessage{
				Type:    communication.EndpointMessageTypeMemoryREQ,
				Payload: struct{}{},
			})
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			// ... and I receive the answer
			msg, err := receiveMessageByName(agentName, ends)
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			if msg.Type != communication.EndpointMessageTypeMemoryRES {
				err := errors.New("unexpected response")
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			// I get the state from the answer...
			state := msg.Payload.(communication.AgentState)
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
			writeResponse(w, http.StatusOK, struct {
				Name   string       `json:"name"`
				Memory mem          `json:"memory"`
				Pool   [][]poolElem `json:"pool"`
			}{
				Name:   agentName,
				Memory: m,
				Pool:   p,
			})
		// If I need to do an input...
		case http.MethodPost:
			// ... I parse the request body to extract the input payload...
			req := struct {
				Actions string `json:"actions"`
			}{}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			// ... I send an input request...
			err = sendMessageByName(agentName, ends, &communication.EndpointMessage{
				Type:    communication.EndpointMessageTypeInputREQ,
				Payload: req.Actions,
			})
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			// ... and I receive the answer
			msg, err := receiveMessageByName(agentName, ends)
			if err != nil || msg.Type != communication.EndpointMessageTypeInputRES {
				log.Println(err)
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// I get the eventual error from the answer...
			errInput := msg.Payload.(string)
			if errInput != "" {
				log.Println(errInput)
				writeError(w, http.StatusBadRequest, errInput)
				return
			}
			// and, if there is none, I respond affirmatively
			writeResponse(w, http.StatusOK, struct {
				Result string `json:"result"`
			}{
				Result: "ok",
			})
		}
	}
}

// sendMessageByName sends a message to an agent, given its name
func sendMessageByName(agentName string, ends map[string]*communication.Endpoint, message *communication.EndpointMessage) error {
	// I check if the agent exists...
	end, ok := ends[agentName]
	if !ok {
		return fmt.Errorf("unknown agent \"%s\"", agentName)
	}
	// ... and I write the message
	err := end.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// receiveMessageByName receives a message from an agent, given its name
func receiveMessageByName(agentName string, ends map[string]*communication.Endpoint) (*communication.EndpointMessage, error) {
	// I check if the agent exists...
	end, ok := ends[agentName]
	if !ok {
		return nil, fmt.Errorf("unknown agent \"%s\"", agentName)
	}
	// ... and I read the message
	msg, err := end.Read()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// writeError writes a JSON error
func writeError(w http.ResponseWriter, h int, e string) {
	// I set the content type...
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// ... I write the specified status code...
	w.WriteHeader(h)
	// ... and I write the error
	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: e,
	})
	w.Write(b)
}

// writeResponse writes a JSON response
func writeResponse(w http.ResponseWriter, h int, p interface{}) {
	// I set the content type...
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// ... I write the specified status code...
	w.WriteHeader(h)
	// ... and I write the response
	b, _ := json.Marshal(p)
	w.Write(b)
}
