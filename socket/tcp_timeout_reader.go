package socket

import (
	"net"
	"slices"
	"time"

	"github.com/Kilemonn/flow/queuedreader"
)

type TCPTimeoutReader struct {
	Listener *net.TCPListener
	Conns    []*net.TCPConn
	indicies []int
}

// Close all connections then the listener. Only the first occurring error will be returned.
func (r TCPTimeoutReader) Close() error {
	var err error
	for _, c := range r.Conns {
		e := c.Close()
		if e != nil && err == nil {
			err = e
		}
	}

	e := r.Listener.Close()
	if e != nil && err == nil {
		err = e
	}

	return err
}

// Get the amount of active connections
func (r *TCPTimeoutReader) connectionCount() int {
	return len(r.Conns)
}

// Check if any incoming connections are pending to be accepted.
// This is naturally blocking, so there is a deadline set for [ScoketReadDeadline]
// before this function returns with no accepted connections.
func (r *TCPTimeoutReader) acceptWaitingConnections() {
	for {
		r.Listener.SetDeadline(time.Now().Add(SocketReadDeadline))
		conn, err := r.Listener.AcceptTCP()
		if err != nil {
			// We got an error and it IS a timeout so leave without error
			if e, ok := err.(net.Error); ok && e.Timeout() {
				return
			}
		}

		r.Conns = append(r.Conns, conn)
	}
}

// Removes connections from the connections list that have been marked for removal.
func (r *TCPTimeoutReader) removeClosedConnections() {
	if len(r.indicies) == 0 {
		return
	}

	slices.Sort(r.indicies)
	// Sort and then reverse iterate so we don't change any of the indicies of further forward elements when we remove them
	for _, i := range slices.Backward(r.indicies) {
		r.Conns = append(r.Conns[:i], r.Conns[i+1:]...)
	}
	r.indicies = []int(nil)
}

// Firstly calls [acceptWaitingConnections].
// Then wraps the read with a deadline to timeout the Read attempt if there is no incoming data.
// Timeout used is [SocketReadDeadline].
// Removes any connections that have been closed.
func (r *TCPTimeoutReader) Read(b []byte) (n int, err error) {
	// Firstly we need to accept any connections and add them to our connection list
	r.acceptWaitingConnections()
	defer r.removeClosedConnections()

	q := queuedreader.NewQueuedReader(r.Conns)
	q.SetPreReadHandlerFunc(func(conn *net.TCPConn) {
		conn.SetReadDeadline(time.Now().Add(SocketReadDeadline))
	})
	// EOF occurs when the remote closes the connection OR when there is no data to be read (depending on the reader)
	q.SetEOFHandlerFunc(func(i int, conn *net.TCPConn) {
		r.indicies = append(r.indicies, i)
	})

	return q.Read(b)
}
