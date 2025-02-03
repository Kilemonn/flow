package socket

import (
	"io"
	"net"
	"time"
)

type UDPTimeoutReader struct {
	Conn *net.UDPConn
}

func (r UDPTimeoutReader) Close() error {
	return r.Conn.Close()
}

// Wraps the read with a deadline to timeout the Read attempt if there is no incoming data.
// Timeout used is [SocketReadDeadline].
func (r UDPTimeoutReader) Read(b []byte) (n int, err error) {
	r.Conn.SetReadDeadline(time.Now().Add(SocketReadDeadline))
	n, err = r.Conn.Read(b)
	if err != nil {
		// We got an error and it IS a timeout so leave without error
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// Return EOF here so the call from io.Copy doesn't permenantly loop
			return n, io.EOF
		}
	}
	return
}
