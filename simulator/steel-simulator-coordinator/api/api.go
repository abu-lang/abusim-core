package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"steel-simulator-common/communication"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Serve serves the API on the API port
func Serve(ends map[string]*communication.Endpoint) {
	// I create the channels to serialize the actions...
	actions := make(chan Action)
	responses := make(chan ActionResponse)
	// ... I create a router for the API and I set the handlers...
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/config/{agentName}", GetHandleConfig(actions, responses)).Methods(http.MethodGet)
	router.HandleFunc("/memory/{agentName}", GetHandleMemory(actions, responses)).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/debug/{agentName}", GetHandleDebug(actions, responses)).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/debug/{agentName}/step", GetHandleDebugStep(actions, responses)).Methods(http.MethodPost)
	// ... I set up the CORS middleware...
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost", "http://localhost:*"},
		AllowedMethods: []string{"POST", "GET"},
		AllowedHeaders: []string{"Accept", "content-type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})
	// ... I run the action processing function...
	go Process(actions, responses, ends)
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
func GetHandleConfig(actions chan Action, responses chan ActionResponse) http.HandlerFunc {
	// I return the handler, decorated with the list of endpoints
	return func(w http.ResponseWriter, r *http.Request) {
		// I get the agent name from the query...
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		// ... and I add a new action to process
		actions <- Action{
			Type:    ActionConfig,
			Payload: agentName,
		}
		// I get the response and I return it
		writeActionResponse(w, <-responses)
	}
}

// GetHandleMemory returns an handler for the memory method
func GetHandleMemory(actions chan Action, responses chan ActionResponse) http.HandlerFunc {
	// I return the handler, decorated with the list of endpoints
	return func(w http.ResponseWriter, r *http.Request) {
		// I get the agent name from the query...
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		// ... and I check what do I have to do
		switch r.Method {
		// If I need to retrieve the memory...
		case http.MethodGet:
			// ... I add a new action to process
			actions <- Action{
				Type:    ActionMemory,
				Payload: agentName,
			}
		// If I need to do an input...
		case http.MethodPost:
			// ... I parse the request body to extract the input payload...
			type request struct {
				Actions string `json:"actions"`
			}
			req := request{}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			// ... and I add a new action to process
			actions <- Action{
				Type: ActionInput,
				Payload: struct {
					agentName string
					actions   string
				}{
					agentName,
					req.Actions,
				},
			}
		}
		// I get the response and I return it
		writeActionResponse(w, <-responses)
	}
}

// GetHandleDebug returns an handler for the debug method
func GetHandleDebug(actions chan Action, responses chan ActionResponse) http.HandlerFunc {
	// I return the handler, decorated with the list of endpoints
	return func(w http.ResponseWriter, r *http.Request) {
		// I get the agent name from the query...
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		// ... and I check what do I have to do
		switch r.Method {
		// If I need to retrieve the debug state...
		case http.MethodGet:
			// ... I add a new action to process
			actions <- Action{
				Type:    ActionDebugInfo,
				Payload: agentName,
			}
		// If I need to change the debug status...
		case http.MethodPost:
			// ... I parse the request body to extract the status payload...
			type request struct {
				Paused    bool   `json:"paused"`
				Verbosity string `json:"verbosity"`
			}
			req := request{}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			actions <- Action{
				Type: ActionDebugSet,
				Payload: struct {
					agentName string
					paused    bool
					verbosity string
				}{
					agentName,
					req.Paused,
					req.Verbosity,
				},
			}
		}
		// I get the response and I return it
		writeActionResponse(w, <-responses)
	}
}

// GetHandleDebugStep returns an handler for the debug step method
func GetHandleDebugStep(actions chan Action, responses chan ActionResponse) http.HandlerFunc {
	// I return the handler, decorated with the list of endpoints
	return func(w http.ResponseWriter, r *http.Request) {
		// I get the agent name from the query...
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		// ... and I add a new action to process
		actions <- Action{
			Type:    ActionDebugStep,
			Payload: agentName,
		}
		// I get the response and I return it
		writeActionResponse(w, <-responses)
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

// writeActionResponse writes an error or response
func writeActionResponse(w http.ResponseWriter, res ActionResponse) {
	// I check whether the Action returned an error or a response and I write it
	if res.Error {
		writeError(w, res.StatusCode, res.Payload.(string))
	} else {
		writeResponse(w, res.StatusCode, res.Payload)
	}
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
