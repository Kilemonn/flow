package bidetwriter

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/Kilemonn/flow/testutil"
	"github.com/stretchr/testify/require"
)

// Test base case where we are trying to write data through a writer that is buffered
// Ensure that performing write does not immediately get flushed by default
func TestBidetWriter_WithBufferedWriter(t *testing.T) {
	data := "TestBidetWriter_WithBufferedWriter"

	testutil.WithTempFile(t, func(filepath string) {
		file, err := os.OpenFile(filepath, os.O_APPEND, os.ModeType)
		require.NoError(t, err)
		defer file.Close()

		writer := bufio.NewWriter(file)
		require.Less(t, len(data), writer.Available())
		writer.Write([]byte(data))

		pos, err := file.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(0), pos)

		// After calling flush, the data is available
		writer.Flush()
		pos, err = file.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(len(data)), pos)
	})
}

// TestBidetWriter_WithBidetWriter ensuring that the flush method is called on every write even if the
// underlying writer's internal buffer is no where near full
func TestBidetWriter_WithBidetWriter(t *testing.T) {
	data := "TestBidetWriter_WithBidetWriter"

	testutil.WithTempFile(t, func(filepath string) {
		file, err := os.OpenFile(filepath, os.O_APPEND, os.ModeType)
		require.NoError(t, err)
		defer file.Close()

		bufferedWriter := bufio.NewWriter(file)
		writer := NewBidetWriter(bufferedWriter)
		require.Less(t, len(data), bufferedWriter.Available())
		writer.Write([]byte(data))

		pos, err := file.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(len(data)), pos)
	})
}
