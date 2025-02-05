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
		// Needed to add os.O_WRONLY for linux
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeType)
		require.NoError(t, err)
		defer file.Close()

		writer := bufio.NewWriter(file)
		require.Less(t, len(data), writer.Available())
		_, err = writer.Write([]byte(data))
		require.NoError(t, err)

		pos, err := file.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(0), pos)

		// After calling flush, the data is available
		require.NoError(t, writer.Flush())
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
		// Needed to add os.O_WRONLY for linux
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeType)
		require.NoError(t, err)
		defer file.Close()

		writer := NewBidetWriter(file)
		defer writer.Close()
		require.Less(t, len(data), writer.Writer.(*bufio.Writer).Available())
		_, err = writer.Write([]byte(data))
		require.NoError(t, err)

		pos, err := file.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(len(data)), pos)
	})
}

// Make sure that the BidetWriter does not close the underlying writer if is it [os.Stdout]
func TestBidetWriter_DoesntCloseWithStdout(t *testing.T) {
	writer := NewBidetWriter(os.Stdout)
	writer.Close()
	data := "TestBidetWriter_DoesntCloseWithStdout"
	n, err := os.Stdout.WriteString(data)

	require.NoError(t, err)
	require.Equal(t, len(data), n)
}
