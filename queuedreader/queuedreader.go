package queuedreader

import (
	"errors"
	"io"
	"net"
)

// An error handler function that receives the index of the reader when the error occurred and the read object itself
type ErrorHandlerFunc[R io.Reader] func(int, R)

// A reader that sequentially attempts to read from the list of [io.Reader], similar to the [io.MultiReader].
// A Pre-read handler function, on EOF handler and on timeout handler function can be provided which will be called
// accordingly.
// Note, that when an EOF or timeout occurs, the reader will attempt to read from the next [io.Reader] until a
// result is returned OR another error occurs (not EOF or timeout).
type QueuedReader[R io.Reader] struct {
	readers []R

	preReadFunc   func(R)
	onEOFFunc     ErrorHandlerFunc[R]
	onTimeoutFunc ErrorHandlerFunc[R]
}

// NewQueuedReader creates a new [QueuedReader] from the provided slice of [io.Reader]s.
func NewQueuedReader[R io.Reader](readers []R) QueuedReader[R] {
	return QueuedReader[R]{
		readers: readers,

		preReadFunc:   nil,
		onEOFFunc:     nil,
		onTimeoutFunc: nil,
	}
}

// SetEOFHandlerFunc sets the function called when EOF occurs (note that this will not stop the read in loop)
func (q *QueuedReader[R]) SetEOFHandlerFunc(handler ErrorHandlerFunc[R]) {
	q.onEOFFunc = handler
}

// SetTimeoutHandlerFunc sets the function called when a timeout occurs (note that this will not stop the
// read in loop)
func (q *QueuedReader[R]) SetTimeoutHandlerFunc(handler ErrorHandlerFunc[R]) {
	q.onTimeoutFunc = handler
}

// SetPreReadHandlerFunc sets the function called before a read attempt (for network connections you can set
// deadline configuration or other things here)
func (q *QueuedReader[R]) SetPreReadHandlerFunc(handler func(R)) {
	q.preReadFunc = handler
}

// Read will iterate over the stored [io.Reader]s and return the first that performs a successful read (without EOF),
// or on the first non-EOF and non-Timeout error, or [io.EOF] will be returned if all [io.Reader]s timeout or return
// [io.EOF].
func (q QueuedReader[R]) Read(b []byte) (int, error) {
	for i, r := range q.readers {
		if q.preReadFunc != nil {
			q.preReadFunc(r)
		}
		n, err := r.Read(b)

		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// Call handler and continue to next reader
				if q.onTimeoutFunc != nil {
					q.onTimeoutFunc(i, r)
				}
			} else if errors.Is(err, io.EOF) {
				// EOF occurs when the remote closes the connection OR when there is no data to be read (depending on the reader)
				// Call handler and continue to next reader
				if q.onEOFFunc != nil {
					q.onEOFFunc(i, r)
				}
			} else {
				// On other errors, make sure we return immediately to the caller
				return n, err
			}
		} else {
			// If there is no error, return the read number of bytes to the caller
			return n, err
		}
	}
	return 0, io.EOF
}
