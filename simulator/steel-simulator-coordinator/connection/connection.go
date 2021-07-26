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
		c, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log.Printf("New agent connected from %s\n", c.RemoteAddr().String())
	coord := communication.New(c)
	initMsg, err := coord.Read()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(initMsg.Payload)
	err = coord.Write(&communication.CoordinatorMessage{
		Type:    communication.CoordinatorMessageTypeACK,
		Payload: struct{}{},
	})
	if err != nil {
		log.Println(err)
		return
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
