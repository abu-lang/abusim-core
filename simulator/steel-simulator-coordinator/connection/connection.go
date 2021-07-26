package connection

import (
	"log"
	"net"
	"steel-simulator-config/communication"
	"time"
)

func GetListener() net.Listener {
	listener, err := net.Listen("tcp4", ":5001")
	if err != nil {
		log.Fatalln(err)
	}
	return listener
}

func AcceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New agent connected from %s\n", conn.RemoteAddr().String())
	coord := communication.New(conn)
	clientName, err := getClientName(coord)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(clientName)
	for {
		time.Sleep(1 * time.Second)
	}
}

func getClientName(coord *communication.Coordinator) (string, error) {
	initMsg, err := coord.Read()
	if err != nil {
		return "", err
	}
	log.Println(initMsg.Payload)
	err = coord.Write(&communication.CoordinatorMessage{
		Type:    communication.CoordinatorMessageTypeACK,
		Payload: struct{}{},
	})
	if err != nil {
		return "", err
	}
}
