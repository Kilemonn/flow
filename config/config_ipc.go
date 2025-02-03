package config

import (
	"io"

	"github.com/Kilemonn/flow/ipc"
)

type ConfigIPC struct {
	ID      string
	Channel string
}

// [ConfigModel.GetID]
func (c ConfigIPC) GetID() string {
	return c.ID
}

// [ConfigModel.Validate]
func (c ConfigIPC) Validate() error {
	return nil
}

// [ConfigModel.Reader]
func (c ConfigIPC) Reader() (io.ReadCloser, error) {
	return ipc.NewIPCReader(c.Channel)
}

// [ConfigModel.Writer]
func (c ConfigIPC) Writer() (io.WriteCloser, error) {
	return ipc.NewIPCWriter(c.Channel)
}
