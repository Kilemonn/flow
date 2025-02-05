package config

import (
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/Kilemonn/flow/serial"
	goSerial "go.bug.st/serial"
)

type ConfigPort struct {
	ID      string
	Channel string
	Mode    goSerial.Mode
	// The resolved and connected port, in a scenario where we call validate
	Port        *serial.CustomPort `json:"-"`
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
		return fmt.Errorf("failed to open connection to port with comm [%s] and ID [%s] with error: [%s]", c.Channel, c.GetID(), err.Error())
	}

	if c.ReadTimeout > 0 {
		err = port.SetReadTimeout(time.Millisecond * time.Duration(c.ReadTimeout))
		if err != nil {
			return fmt.Errorf("failed to set timeout on serial port connection with comm [%s] and ID [%s] with error: [%s]", c.Channel, c.GetID(), err.Error())
		}
	}

	customPort := serial.NewCustomPort(port)
	c.Port = &customPort
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
