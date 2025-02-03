package stdio

import (
	"io"
	"testing"

	"github.com/Kilemonn/flow/testutil"
	"github.com/stretchr/testify/require"
)

// TestStdInReader ensure std in always returns nothing by default
func TestStdInReader(t *testing.T) {
	buffer := make([]byte, 10)
	reader, err := CreateStdInReader()
	require.NoError(t, err)
	n, _ := reader.Read(buffer)
	require.Equal(t, 0, n)
}

// TestCreateStdInReader_WithNameLine ensure that reading in reads the exact bytes, even hidden bytes
func TestCreateStdInReader_WithNameLine(t *testing.T) {
	expected := "TestCreateStdInReader_WithNameLine\n"

	testutil.WithBytesInStdIn(t, []byte(expected), func() {
		reader, err := CreateStdInReader()
		require.NoError(t, err)
		bytes, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, expected, string(bytes))
	})
}

// TestCreateStdInReader_WithoutNewLine ensure that reading in reads the exact bytes, and doesn't add hidden bytes
func TestCreateStdInReader_WithoutNewLine(t *testing.T) {
	expected := "TestCreateStdInReader_WithoutNewLine"

	testutil.WithBytesInStdIn(t, []byte(expected), func() {
		reader, err := CreateStdInReader()
		require.NoError(t, err)
		bytes, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, expected, string(bytes))
	})
}
