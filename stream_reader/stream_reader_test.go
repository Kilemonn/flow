package stream_reader

import (
	"encoding/binary"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/Kilemonn/flow/testutil"
	"github.com/stretchr/testify/require"
)

type MyType struct {
	I    int32
	Data [20]byte
}

func TestNewStreamWriter(t *testing.T) {
	testutil.WithTempFile(t, func(filepath string) {
		file, err := os.OpenFile(filepath, os.O_RDWR, fs.ModeType)
		require.NoError(t, err)
		defer file.Close()

		str := "a test str"
		initial := MyType{
			I: int32(len(str)),
		}
		copy(initial.Data[:], str)

		require.NoError(t, binary.Write(file, binary.BigEndian, initial))

		_, err = file.Seek(0, io.SeekStart)
		require.NoError(t, err)

		reader := NewStreamReader[MyType](file, 4, func(b []byte) int {
			// Returning 20, since the data portion is 20 bytes
			return int(20)
		})

		obj, err := reader.Read()
		require.NoError(t, err)

		require.Equal(t, initial.I, obj.I)
		require.Equal(t, len(initial.Data), len(obj.Data))
		require.Equal(t, initial.Data, obj.Data)

		obj, err = reader.Read()
		require.Error(t, err)
	})
}
