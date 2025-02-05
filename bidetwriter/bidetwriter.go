package bidetwriter

import (
	"bufio"
	"io"
	"os"
)

// BidetWriter a writer that calls the provided Flushfunc() function after the io.Writer.Write() call.
type BidetWriter struct {
	writeCloser io.Closer
	Writer      io.Writer
	FlushFunc   func() error
}

// Takes a write closer and will close the file if Close() is called.
func NewBidetWriter(w io.WriteCloser) BidetWriter {
	bufWriter := bufio.NewWriter(w)
	return BidetWriter{
		writeCloser: w,
		Writer:      bufWriter,
		FlushFunc:   bufWriter.Flush,
	}
}

func (bw BidetWriter) Close() error {
	if bw.writeCloser != os.Stdout {
		return bw.writeCloser.Close()
	}
	return nil
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
