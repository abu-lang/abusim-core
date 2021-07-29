package endpoint

import (
	"log"
	"net"
	"steel-simulator-common/communication"
)

// GetListener returns a listener on the control port
func GetListener() net.Listener {
	// I create a TCP listener on the specified port...
	listener, err := net.Listen("tcp4", ":5001")
	if err != nil {
		log.Fatalln(err)
	}
	// ... and I return it
	return listener
}

// HandleConnections handles the incoming connections from agents
func HandleConnections(listener net.Listener, ends map[string]*communication.Endpoint) {
	// I loop...
	for {
		// ... I accept an incoming connection...
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// ... and I handle it
		go handleConnection(conn, ends)
	}
}

// handleConnection handles a single incoming connection
func handleConnection(conn net.Conn, ends map[string]*communication.Endpoint) {
	log.Printf("New agent connected from %s\n", conn.RemoteAddr().String())
	// I create a new endpoint...
	end := communication.New(conn)
	// ... I receive the initialization message...
	initMsg, err := end.Read()
	if err != nil {
		log.Println(err)
		return
	}
	// ... and I acknowledge it
	err = end.Write(&communication.EndpointMessage{
		Type:    communication.EndpointMessageTypeACK,
		Payload: struct{}{},
	})
	if err != nil {
		log.Println(err)
		return
	}
	// Finally, I add the endpoint to the endpoints pool
	ends[initMsg.Payload.(string)] = end
}
