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

type Coordinator struct {
	readwrite *bufio.ReadWriter
}

type CoordinatorMessageType int

const (
	CoordinatorMessageTypeACK       = iota
	CoordinatorMessageTypeINIT      = iota
	CoordinatorMessageTypeMemoryREQ = iota
	CoordinatorMessageTypeMemoryRES = iota
	CoordinatorMessageTypeInputREQ  = iota
	CoordinatorMessageTypeInputRES  = iota
	CoordinatorMessageTypeConfigREQ = iota
	CoordinatorMessageTypeConfigRES = iota
)

type CoordinatorMessage struct {
	Type    CoordinatorMessageType
	Payload interface{}
}

func New(conn net.Conn) *Coordinator {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	gob.Register(struct{}{})
	gob.Register(datastructure.Resources{})
	gob.Register(config.Agent{})

	return &Coordinator{
		readwrite: bufio.NewReadWriter(r, w),
	}
}

func (c *Coordinator) Read() (*CoordinatorMessage, error) {
	h := make([]byte, 4)
	_, err := io.ReadFull(c.readwrite, h)
	if err != nil {
		return nil, err
	}
	l := binary.BigEndian.Uint32(h)
	b := make([]byte, l)
	_, err = io.ReadFull(c.readwrite, b)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	buf.Write(b)
	msg := &CoordinatorMessage{}
	err = gob.NewDecoder(&buf).Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Coordinator) Write(msg *CoordinatorMessage) error {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return err
	}
	h := make([]byte, 4)
	binary.BigEndian.PutUint32(h, uint32(buf.Len()))
	_, err = c.readwrite.Write(h)
	if err != nil {
		return err
	}
	_, err = c.readwrite.Write(buf.Bytes())
	if err != nil {
		return err
	}
	c.readwrite.Flush()
	return nil
}
