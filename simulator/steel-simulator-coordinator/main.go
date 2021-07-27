package main

import (
	"log"
	"os"
	"os/signal"
	"steel-simulator-coordinator/api"
	"steel-simulator-coordinator/connection"
	"syscall"
)

func main() {
	agents := make(map[string]*connection.ConnCoord)
	setupCloseHandler(agents)
	log.Println("Starting listener")
	listener := connection.GetListener()
	defer listener.Close()
	go connection.AcceptLoop(listener, agents)
	log.Println("Starting API")
	api.Serve(agents)
}

func setupCloseHandler(agents map[string]*connection.ConnCoord) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		for _, agent := range agents {
			agent.Conn.Close()
		}
		os.Exit(0)
	}()
}
