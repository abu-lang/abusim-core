package connection

import (
	"log"
	"net"
	"steel-simulator-common/communication"
)

type ConnCoord struct {
	Conn  net.Conn
	Coord *communication.Coordinator
}

func GetListener() net.Listener {
	listener, err := net.Listen("tcp4", ":5001")
	if err != nil {
		log.Fatalln(err)
	}
	return listener
}

func AcceptLoop(listener net.Listener, agents map[string]*ConnCoord) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, agents)
	}
}

func handleConnection(conn net.Conn, agents map[string]*ConnCoord) {
	log.Printf("New agent connected from %s\n", conn.RemoteAddr().String())
	coord := communication.New(conn)
	agentName, err := getAgentName(coord)
	if err != nil {
		log.Println(err)
		return
	}
	agents[agentName] = &ConnCoord{
		Conn:  conn,
		Coord: coord,
	}
}

func getAgentName(coord *communication.Coordinator) (string, error) {
	initMsg, err := coord.Read()
	if err != nil {
		return "", err
	}
	err = coord.Write(&communication.CoordinatorMessage{
		Type:    communication.CoordinatorMessageTypeACK,
		Payload: struct{}{},
	})
	if err != nil {
		return "", err
	}
	return initMsg.Payload.(string), nil
}
