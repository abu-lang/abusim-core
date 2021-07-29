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
	ends := make(map[string]*communication.Endpoint)
	setupCloseHandler(ends)
	log.Println("Starting listener")
	listener := endpoint.GetListener()
	defer listener.Close()
	go endpoint.HandleConnections(listener, ends)
	log.Println("Starting API")
	api.Serve(ends)
}

func setupCloseHandler(ends map[string]*communication.Endpoint) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		for _, end := range ends {
			end.Close()
		}
		os.Exit(0)
	}()
}
