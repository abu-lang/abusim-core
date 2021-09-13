package schema

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"net"

	"github.com/abu-lang/goabu/memory"
)

// Endpoint represents an agent-coordinatior connection side
type Endpoint struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

type EndpointMessageType int

const (
	EndpointMessageTypeACK            = iota
	EndpointMessageTypeINIT           = iota
	EndpointMessageTypeMemoryREQ      = iota
	EndpointMessageTypeInputREQ       = iota
	EndpointMessageTypeConfigREQ      = iota
	EndpointMessageTypeDebugREQ       = iota
	EndpointMessageTypeDebugChangeREQ = iota
	EndpointMessageTypeDebugStepREQ   = iota
)

// EndpointMessage represents an agent-coordinatior message
type EndpointMessage struct {
	Type    EndpointMessageType
	Payload interface{}
}

// AgentState represents an agent state
type AgentState struct {
	Memory memory.Resources
	Pool   [][]PoolElem
}

// PoolElem represents an pool element
type PoolElem struct {
	Resource string
	Value    string
}

// AgentDebugStatus represents an agent debug status
type AgentDebugStatus struct {
	Paused    bool
	Verbosity string
}

// New creates a new endpoint from a connection
func New(conn net.Conn) *Endpoint {
	// I create a reader and a writer for the connection...
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	// ... I register the structures passed as payload in the messages...
	gob.Register(struct{}{})
	gob.Register(AgentState{})
	gob.Register(AgentDebugStatus{})
	gob.Register(Agent{})
	// ... and I return the endpoint
	return &Endpoint{
		conn:   conn,
		reader: r,
		writer: w,
	}
}

// Read expect a message and returns it
func (e *Endpoint) Read() (*EndpointMessage, error) {
	// I read the 4 bytes that contains the length of the message...
	h := make([]byte, 4)
	_, err := io.ReadFull(e.reader, h)
	if err != nil {
		return nil, err
	}
	// ... I get the message length...
	l := binary.BigEndian.Uint32(h)
	// ... I read the entire message...
	b := make([]byte, l)
	_, err = io.ReadFull(e.reader, b)
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
	_, err = e.writer.Write(h)
	if err != nil {
		return err
	}
	// ... and I send the message itself
	_, err = e.writer.Write(buf.Bytes())
	if err != nil {
		return err
	}
	// Finally, I flush the writer to ensure the message is gone
	e.writer.Flush()
	return nil
}

// Close closes the connection
func (e *Endpoint) Close() {
	e.conn.Close()
}
