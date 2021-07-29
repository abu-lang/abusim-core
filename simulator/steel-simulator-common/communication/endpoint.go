package communication

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"net"
	"steel-simulator-common/config"

	"steel-lang/datastructure"
)

// Endpoint represent an agent-coordinatior connection side
type Endpoint struct {
	conn      net.Conn
	readwrite *bufio.ReadWriter
}

type EndpointMessageType int

const (
	EndpointMessageTypeACK       = iota
	EndpointMessageTypeINIT      = iota
	EndpointMessageTypeMemoryREQ = iota
	EndpointMessageTypeMemoryRES = iota
	EndpointMessageTypeInputREQ  = iota
	EndpointMessageTypeInputRES  = iota
	EndpointMessageTypeConfigREQ = iota
	EndpointMessageTypeConfigRES = iota
)

// EndpointMessage represent an agent-coordinatior message
type EndpointMessage struct {
	Type    EndpointMessageType
	Payload interface{}
}

// New creates a new endpoint from a connection
func New(conn net.Conn) *Endpoint {
	// I create a reader and a writer for the connection...
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	// ... I register the structures passed as payload in the messages...
	gob.Register(struct{}{})
	gob.Register(datastructure.Resources{})
	gob.Register(config.Agent{})
	// ... and I return the endpoint
	return &Endpoint{
		conn:      conn,
		readwrite: bufio.NewReadWriter(r, w),
	}
}

// Read expect a message and returns it
func (e *Endpoint) Read() (*EndpointMessage, error) {
	// I read the 4 bytes that contains the length of the message...
	h := make([]byte, 4)
	_, err := io.ReadFull(e.readwrite, h)
	if err != nil {
		return nil, err
	}
	// ... I get the message length...
	l := binary.BigEndian.Uint32(h)
	// ... I read the entire message...
	b := make([]byte, l)
	_, err = io.ReadFull(e.readwrite, b)
	if err != nil {
		return nil, err
	}
	// ... and I decode it
	buf := bytes.Buffer{}
	buf.Write(b)
	msg := &EndpointMessage{}
	err = gob.NewDecoder(&buf).Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Write sends a message
func (e *Endpoint) Write(msg *EndpointMessage) error {
	// I encode the message...
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return err
	}
	// ... I send the length...
	h := make([]byte, 4)
	binary.BigEndian.PutUint32(h, uint32(buf.Len()))
	_, err = e.readwrite.Write(h)
	if err != nil {
		return err
	}
	// ... and I send the message itself
	_, err = e.readwrite.Write(buf.Bytes())
	if err != nil {
		return err
	}
	// Finally, I flush the writer to ensure the message is gone
	e.readwrite.Flush()
	return nil
}

func (e *Endpoint) Close() {
	e.conn.Close()
}
