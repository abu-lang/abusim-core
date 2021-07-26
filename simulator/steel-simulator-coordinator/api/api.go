package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/memory/{agentName}", HandleMemory)

	log.Fatal(http.ListenAndServe(":4000", router))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func HandleMemory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentName := vars["agentName"]
	fmt.Fprintln(w, "Agent name:", agentName)
}
