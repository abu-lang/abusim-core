package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"steel-lang/datastructure"
	"steel-simulator-common/communication"
	"steel-simulator-coordinator/connection"
	"time"

	"github.com/gorilla/mux"
)

func Serve(agents map[string]*connection.ConnCoord) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/memory/{agentName}", GetHandleMemory(agents)).Methods(http.MethodGet, http.MethodPost)

	log.Fatal(http.ListenAndServe(":4000", router))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func GetHandleMemory(agents map[string]*connection.ConnCoord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentName := vars["agentName"]
		switch r.Method {
		case http.MethodGet:
			err := sendMessageByName(agentName, agents, &communication.CoordinatorMessage{
				Type:    communication.CoordinatorMessageTypeMemoryREQ,
				Payload: struct{}{},
			})
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			msg, err := receiveMessageByName(agentName, agents)
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			if msg.Type != communication.CoordinatorMessageTypeMemoryRES {
				err := errors.New("unexpected response")
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			resources := msg.Payload.(datastructure.Resources)
			writeResponse(w, http.StatusOK, struct {
				Bool    map[string]bool        `json:"bool"`
				Integer map[string]int64       `json:"integer"`
				Float   map[string]float64     `json:"float"`
				Text    map[string]string      `json:"text"`
				Time    map[string]time.Time   `json:"time"`
				Other   map[string]interface{} `json:"other"`
			}{
				Bool:    resources.Bool,
				Integer: resources.Integer,
				Float:   resources.Float,
				Text:    resources.Text,
				Time:    resources.Time,
				Other:   resources.Other,
			})
		case http.MethodPost:
			req := struct {
				Actions string `json:"actions"`
			}{}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			err = sendMessageByName(agentName, agents, &communication.CoordinatorMessage{
				Type:    communication.CoordinatorMessageTypeInputREQ,
				Payload: req.Actions,
			})
			if err != nil {
				log.Println(err)
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			msg, err := receiveMessageByName(agentName, agents)
			if err != nil || msg.Type != communication.CoordinatorMessageTypeInputRES {
				log.Println(err)
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			errInput := msg.Payload.(string)
			if errInput != "" {
				log.Println(errInput)
				writeError(w, http.StatusBadRequest, errInput)
				return
			}
			writeResponse(w, http.StatusOK, struct {
				Result string `json:"result"`
			}{
				Result: "ok",
			})
		}
	}
}

func sendMessageByName(agentName string, agents map[string]*connection.ConnCoord, message *communication.CoordinatorMessage) error {
	agent, ok := agents[agentName]
	if !ok {
		return fmt.Errorf("unknown agent \"%s\"", agentName)
	}
	err := agent.Coord.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func receiveMessageByName(agentName string, agents map[string]*connection.ConnCoord) (*communication.CoordinatorMessage, error) {
	agent, ok := agents[agentName]
	if !ok {
		return nil, fmt.Errorf("unknown agent \"%s\"", agentName)
	}
	msg, err := agent.Coord.Read()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func writeError(w http.ResponseWriter, h int, e string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(h)
	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: e,
	})
	w.Write(b)
}

func writeResponse(w http.ResponseWriter, h int, p interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(h)
	b, _ := json.Marshal(p)
	w.Write(b)
}
