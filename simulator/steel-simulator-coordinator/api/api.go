package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"steel-simulator-config/communication"

	"github.com/gorilla/mux"
)

func Serve(agents map[string]*communication.Coordinator) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/memory/{agentName}", GetHandleMemory(agents)).Methods(http.MethodGet, http.MethodPost)

	log.Fatal(http.ListenAndServe(":4000", router))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func GetHandleMemory(agents map[string]*communication.Coordinator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		switch r.Method {
		case http.MethodGet:
			err := sendMessageByName(agentName, agents, &communication.CoordinatorMessage{
				Type:    communication.CoordinatorMessageTypeACK,
				Payload: "GET",
			})
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			writeResponse(w, http.StatusOK, struct {
				Method string `json:"method"`
				Agent  string `json:"agent"`
			}{
				Method: "GET",
				Agent:  agentName,
			})
		case http.MethodPost:
			err := sendMessageByName(agentName, agents, &communication.CoordinatorMessage{
				Type:    communication.CoordinatorMessageTypeACK,
				Payload: "POST",
			})
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			writeResponse(w, http.StatusOK, struct {
				Method string `json:"method"`
				Agent  string `json:"agent"`
			}{
				Method: "POST",
				Agent:  agentName,
			})
		}
	}
}

func sendMessageByName(agentName string, agents map[string]*communication.Coordinator, message *communication.CoordinatorMessage) error {
	agent, ok := agents[agentName]
	if !ok {
		return fmt.Errorf("unknown agent \"%s\"", agentName)
	}
	err := agent.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, h int, e string) {
	w.WriteHeader(h)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: e,
	})
	w.Write(b)
}

func writeResponse(w http.ResponseWriter, h int, p interface{}) {
	w.WriteHeader(h)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, _ := json.Marshal(p)
	w.Write(b)
}
