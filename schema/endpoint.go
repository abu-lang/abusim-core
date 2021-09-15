package schema

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

// Endpoint represents an agent-coordinatior connection side
type Endpoint struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// New creates a new endpoint from a connection
func New(conn net.Conn) *Endpoint {
	// I create a reader and a writer for the connection...
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
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
	err = json.NewDecoder(&buf).Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Write sends a message
func (e *Endpoint) Write(msg *EndpointMessage) error {
	// I encode the message...
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(msg)
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
