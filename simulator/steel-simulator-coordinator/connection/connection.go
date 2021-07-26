package connection

import (
	"fmt"
	"log"
	"net"
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
	fmt.Printf("New agent connected from %s\n", c.RemoteAddr().String())
}
