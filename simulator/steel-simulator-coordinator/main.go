package main

import (
	"log"
	"steel-simulator-config/communication"
	"steel-simulator-coordinator/api"
	"steel-simulator-coordinator/connection"
)

func main() {
	agents := make(map[string]*communication.Coordinator)
	log.Println("Starting listener")
	listener := connection.GetListener()
	defer listener.Close()
	go connection.AcceptLoop(listener, agents)
	log.Println("Starting API")
	api.Serve(agents)
}
