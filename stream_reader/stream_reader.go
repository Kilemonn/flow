package stream_reader

import (
	"bytes"
	"encoding/binary"
	"io"
)

type StreamReader[T any] struct {
	buff                bytes.Buffer
	Reader              io.Reader
	headerSize          int
	expectedMessageSize int
	headerToSizeFunc    func([]byte) int
}

func NewStreamReader[T any](reader io.Reader, headerSize int, headerToSize func([]byte) int) StreamReader[T] {
	return StreamReader[T]{
		Reader:              reader,
		buff:                bytes.Buffer{},
		headerSize:          headerSize,
		expectedMessageSize: 0,
		headerToSizeFunc:    headerToSize,
	}
}

func (r StreamReader[T]) Read() (t T, err error) {
	// Then we are not in the middle of reading a message and should check for headers
	if r.expectedMessageSize == 0 {
		b := make([]byte, 100)
		_, err := r.Reader.Read(b)
		if err != nil {
			return t, err
		}
		_, err = r.buff.Write(b)
		if err != nil {
			return t, err
		}

		if r.buff.Len() >= r.headerSize {
			header := make([]byte, r.headerSize)
			_, err = r.buff.Read(header)
			if err != nil {
				return t, err
			}

			r.expectedMessageSize = r.headerToSizeFunc(header)

			if r.buff.Len() >= r.expectedMessageSize {
				var res bytes.Buffer
				_, err = res.Write(header)
				if err != nil {
					return t, err
				}

				m := make([]byte, r.expectedMessageSize)
				_, err = r.buff.Read(m)
				if err != nil {
					return t, err
				}

				_, err = res.Write(m)
				if err != nil {
					return t, err
				}

				err = binary.Read(bytes.NewReader(res.Bytes()), binary.BigEndian, &t)
				return t, err
			}
		}
	} else {
		// We have enough bytes to construct and check another header
		if r.buff.Len() >= int(r.headerSize) {

		}
	}

	if r.buff.Len() >= int(r.headerSize) {

	}

	return t, nil
}

func (r StreamReader[T]) HeaderSize() int {
	return r.headerSize
}
