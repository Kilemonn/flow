package serial

import (
	"time"

	goSerial "go.bug.st/serial"
)

type TimeoutPort struct {
	Port    goSerial.Port
	Timeout time.Duration
}

func NewTimeoutPort(port goSerial.Port, timeout time.Duration) TimeoutPort {
	return TimeoutPort{
		Port:    port,
		Timeout: timeout,
	}
}

// [io.Reader.Read]
func (p TimeoutPort) Read(b []byte) (int, error) {
	if p.Timeout > 0 {
		p.Port.SetReadTimeout(p.Timeout)
	}
	return p.Port.Read(b)
}

// [io.Writer.Write]
func (p TimeoutPort) Write(b []byte) (int, error) {
	return p.Port.Write(b)
}

// [io.Closer.Close]
func (p TimeoutPort) Close() error {
	return p.Port.Close()
}
