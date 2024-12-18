package serial

import (
	"fmt"
	"io"
	"slices"

	goSerial "go.bug.st/serial"
)

func CreatePort(channel string, mode goSerial.Mode) (io.ReadWriteCloser, error) {
	return CreateReadWriter(channel, mode)
}

func CreateReadWriter(portName string, mode goSerial.Mode) (io.ReadWriteCloser, error) {
	ports := GetSerialPorts()
	if !slices.Contains(ports, portName) {
		return nil, fmt.Errorf("failed to find port with name [%s]", portName)
	}

	port, err := OpenSerialConnection(portName, mode)
	if err != nil {
		return nil, err
	}

	return port, err

}
