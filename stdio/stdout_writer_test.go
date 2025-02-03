package stdio

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/Kilemonn/flow/testutil"
	"github.com/stretchr/testify/require"
)

func TestMultiWriter(t *testing.T) {
	// Two empty buffers
	var foo, bar bytes.Buffer

	// Create a multi writer
	mw := io.MultiWriter(&foo, &bar)

	// Write message into multi writer
	fmt.Fprintln(mw, "Multi writer test")

	// Optional: verfiy data stored in buffers
	require.Equal(t, foo.String(), bar.String())
}

// TestStdoutWriter ensure that using StdIO will write to std out
func TestStdoutWriter(t *testing.T) {
	expected := "My name is Mr Cow"
	written := testutil.CaptureStdout(t, func() {
		writer, err := CreateStdOutWriter()
		require.NoError(t, err)
		writer.Write([]byte(expected))
	})
	read, err := io.ReadAll(written)
	require.NoError(t, err)
	require.Equal(t, expected, string(read))
}
