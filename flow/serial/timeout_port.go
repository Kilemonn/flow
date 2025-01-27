package serial

import (
	"io"
	"time"

	goSerial "go.bug.st/serial"
)

type TimeoutPort struct {
	Port    goSerial.Port
	Timeout time.Duration
}

func NewTimeoutPort(port goSerial.Port, timeout time.Duration) (TimeoutPort, error) {
	err := port.SetReadTimeout(timeout)
	return TimeoutPort{
		Port:    port,
		Timeout: timeout,
	}, err
}

// [io.Reader.Read]
func (p TimeoutPort) Read(b []byte) (n int, err error) {
	n, err = p.Port.Read(b)

	// TODO: The library should fix this, https://github.com/bugst/go-serial/issues/141
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return n, err
}

// [io.Writer.Write]
func (p TimeoutPort) Write(b []byte) (int, error) {
	return p.Port.Write(b)
}

// [io.Closer.Close]
func (p TimeoutPort) Close() error {
	return p.Port.Close()
}
