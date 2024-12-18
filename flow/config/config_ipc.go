package config

import (
	"io"

	"github.com/Kilemonn/flow/flow/ipc"
)

type ConfigIPC struct {
	ID      string
	Channel string
}

func (c ConfigIPC) GetID() string {
	return c.ID
}

func (c ConfigIPC) Validate() error {
	return nil
}

func (c ConfigIPC) Reader() (io.ReadCloser, error) {
	return ipc.NewIPCReader(c.Channel)
}

func (c ConfigIPC) Writer() (io.WriteCloser, error) {
	return ipc.NewIPCWriter(c.Channel)
}
