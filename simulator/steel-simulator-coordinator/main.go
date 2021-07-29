package main

import (
	"log"
	"os"
	"os/signal"
	"steel-simulator-common/communication"
	"steel-simulator-coordinator/api"
	"steel-simulator-coordinator/endpoint"
	"syscall"
)

func main() {
	// I create a map for the endpoints...
	ends := make(map[string]*communication.Endpoint)
	// ... I set up the handler to close the connections...
	setupCloseHandler(ends)
	// ... I listen for connection...
	log.Println("Starting listener")
	listener := endpoint.GetListener()
	defer listener.Close()
	// ... I handle the incoming connections...
	go endpoint.HandleConnections(listener, ends)
	// ... and I serve the API
	log.Println("Starting API")
	api.Serve(ends)
}

// setupCloseHandler waits for a SIGTERM and then closes all the connections
func setupCloseHandler(ends map[string]*communication.Endpoint) {
	// I register for the SIGTERMs...
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// ... and I run a goroutine to handle their arrival
	go func() {
		// I block until a SIGTERM...
		<-c
		// ... I close all the connections...
		for _, end := range ends {
			end.Close()
		}
		// ... and I exit
		os.Exit(0)
	}()
}
