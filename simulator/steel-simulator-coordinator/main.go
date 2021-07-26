package main

import (
	"log"
	"steel-simulator-coordinator/api"
	"steel-simulator-coordinator/connection"
)

func main() {
	log.Println("Starting listener")
	listener := connection.GetListener()
	defer listener.Close()
	go connection.AcceptLoop(listener)
	log.Println("Starting API")
	api.Serve()
}
