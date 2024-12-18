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

// [ConfigModel.GetID]
func (c ConfigSocket) GetID() string {
	return c.ID
}

// [ConfigModel.Validate]
func (c ConfigSocket) Validate() error {
	return nil
}

// [ConfigModel.Reader]
func (c ConfigSocket) Reader() (io.ReadCloser, error) {
	return socket.CreateSocketReader(c.Protocol, c.Address, c.Port)
}

// [ConfigModel.Writer]
func (c ConfigSocket) Writer() (io.WriteCloser, error) {
	return socket.CreateSocketWriter(c.Protocol, c.Address, c.Port)
}
