package bidetwriter

import (
	"bufio"
	"io"
)

// BidetWriter a writer that calls the provided Flushfunc() function after the io.Writer.Write() call.
type BidetWriter struct {
	Writer    io.Writer
	FlushFunc func() error
}

func NewBidetWriter(w *bufio.Writer) BidetWriter {
	return BidetWriter{
		Writer:    w,
		FlushFunc: w.Flush,
	}
}

func (bw BidetWriter) Close() error {
	return bw.FlushFunc()
}

func (bw BidetWriter) Write(b []byte) (n int, err error) {
	n, err = bw.Writer.Write(b)
	flushErr := bw.FlushFunc()
	// Return the Write() error first before we return the underlying flush error
	if err == nil {
		err = flushErr
	}
	return
}
