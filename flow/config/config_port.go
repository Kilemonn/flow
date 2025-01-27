package config

import (
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/Kilemonn/flow/flow/serial"
	goSerial "go.bug.st/serial"
)

type ConfigPort struct {
	ID      string
	Channel string
	Mode    goSerial.Mode
	// The resolved and connected port, in a scenario where we call validate
	Port        *serial.TimeoutPort `json:"-"`
	ReadTimeout int
}

// [ConfigModel.GetID]
func (c ConfigPort) GetID() string {
	return c.ID
}

// [ConfigModel.Validate]
func (c ConfigPort) Validate() error {
	ports := serial.GetSerialPorts()
	if !slices.Contains(ports, c.Channel) {
		return fmt.Errorf("no port with name [%s] is available/connected", c.Channel)
	}
	return nil
}

func (c *ConfigPort) Open() error {
	port, err := serial.OpenSerialConnection(c.Channel, c.Mode)
	if err != nil {
		return fmt.Errorf("failed to open connection to port with comm [%s] and ID [%s] with error: [%s]", c.Channel, c.ID, err.Error())
	}
	timeoutPort := serial.NewTimeoutPort(port, (time.Millisecond * time.Duration(c.ReadTimeout)))
	c.Port = &timeoutPort

	return nil
}

// [ConfigModel.Reader]
func (c *ConfigPort) Reader() (io.ReadCloser, error) {
	if c.Port == nil {
		err := c.Open()
		if err != nil {
			return nil, err
		}
	}
	return *c.Port, nil
}

// [ConfigModel.Writer]
func (c *ConfigPort) Writer() (io.WriteCloser, error) {
	if c.Port == nil {
		err := c.Open()
		if err != nil {
			return nil, err
		}
	}
	return *c.Port, nil
}
