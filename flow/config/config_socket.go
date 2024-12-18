package config

import (
	"io"

	"github.com/Kilemonn/flow/flow/socket"
)

type ConfigSocket struct {
	ID       string
	Protocol string
	Port     uint16
	Address  string
}

func (c ConfigSocket) GetID() string {
	return c.ID
}

func (c ConfigSocket) Validate() error {
	return nil
}

func (c ConfigSocket) Reader() (io.ReadCloser, error) {
	return socket.CreateSocketReader(c.Protocol, c.Address, c.Port)
}

func (c ConfigSocket) Writer() (io.WriteCloser, error) {
	return socket.CreateSocketWriter(c.Protocol, c.Address, c.Port)
}
