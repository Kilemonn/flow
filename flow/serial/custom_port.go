package serial

import (
	"io"

	goSerial "go.bug.st/serial"
)

// Wraps the [goSerial.Port] because its Read() function does not return EOF on timeout. This causes problems with [io.Copy]
// when this is used as a reader, since no error is returned it will hang.
// https://github.com/bugst/go-serial/issues/141
type CustomPort struct {
	Port goSerial.Port
}

func NewCustomPort(port goSerial.Port) CustomPort {
	return CustomPort{
		Port: port,
	}
}

// [io.Reader.Read]
func (p CustomPort) Read(b []byte) (n int, err error) {
	n, err = p.Port.Read(b)

	// TODO: The library should fix this, https://github.com/bugst/go-serial/issues/141
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return n, err
}

// [io.Writer.Write]
func (p CustomPort) Write(b []byte) (int, error) {
	return p.Port.Write(b)
}

// [io.Closer.Close]
func (p CustomPort) Close() error {
	return p.Port.Close()
}
