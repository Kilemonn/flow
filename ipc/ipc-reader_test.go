package ipc

import (
	"io"
	"testing"
	"time"

	"github.com/Kilemonn/flow/testutil"
	"github.com/stretchr/testify/require"
)

// Ensure that when attempting to accept a new IPC connection when there is none, that the deadline
// kicks in and stops infinitely blocking the call, and also make sure no new connection is added.
func TestIPCRead_NoNewConnections(t *testing.T) {
	reader, err := NewIPCReader("TestIPCRead_NoNewConnections")
	require.NoError(t, err)
	defer reader.Close()

	require.Equal(t, 0, reader.(*IPCReader).connectionCount())
	testutil.TakesAtleast(t, IPCReadDeadline, func() {
		reader.(*IPCReader).acceptWaitingConnections()
	})
	require.Equal(t, 0, reader.(*IPCReader).connectionCount())
}

// Ensure pending IPC connections are accepted and added to the conns list
func TestIPCPAcceptNewConnections(t *testing.T) {
	reader, err := NewIPCReader("TestIPCPAcceptNewConnections")
	require.NoError(t, err)
	defer reader.Close()

	w, err := NewIPCWriter("TestIPCPAcceptNewConnections")
	require.NoError(t, err)
	defer w.Close()

	require.Equal(t, 0, reader.(*IPCReader).connectionCount())
	testutil.TakesAtleast(t, 0*time.Millisecond, func() {
		reader.(*IPCReader).acceptWaitingConnections()
	})
	require.Equal(t, 1, reader.(*IPCReader).connectionCount())

	w2, err := NewIPCWriter("TestIPCPAcceptNewConnections")
	require.NoError(t, err)
	defer w2.Close()

	testutil.TakesAtleast(t, 0*time.Millisecond, func() {
		reader.(*IPCReader).acceptWaitingConnections()
	})
	require.Equal(t, 2, reader.(*IPCReader).connectionCount())
}

// Make sure the call to read will accept new connections even if there is no data waiting to be read
func TestIPCRead_NewConnectionsNoData(t *testing.T) {
	reader, err := NewIPCReader("TestIPCRead_NewConnectionsNoData")
	require.NoError(t, err)
	defer reader.Close()

	w, err := NewIPCWriter("TestIPCRead_NewConnectionsNoData")
	require.NoError(t, err)
	defer w.Close()

	w2, err := NewIPCWriter("TestIPCRead_NewConnectionsNoData")
	require.NoError(t, err)
	defer w2.Close()

	require.Equal(t, 0, reader.(*IPCReader).connectionCount())

	b := make([]byte, 10)
	n, err := reader.Read(b)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)
	require.Equal(t, 2, reader.(*IPCReader).connectionCount())
}

// Ensure that we read data from multiple connections if they have data.
// Make sure once all data is exhausted from all connections that we get an EOF back.
func TestIPCRead_DataFromMultipleSockets(t *testing.T) {
	reader, err := NewIPCReader("TestIPCRead_DataFromMultipleSockets")
	require.NoError(t, err)
	defer reader.Close()

	content := "TestTCPRead_DataFromMultipleSockets"

	writer1, err := NewIPCWriter("TestIPCRead_DataFromMultipleSockets")
	require.NoError(t, err)
	defer writer1.Close()

	n, err := writer1.Write([]byte(content))
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	writer2, err := NewIPCWriter("TestIPCRead_DataFromMultipleSockets")
	require.NoError(t, err)
	defer writer2.Close()

	n, err = writer2.Write([]byte(content))
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	require.Equal(t, 0, reader.(*IPCReader).connectionCount())

	b := make([]byte, len(content))
	n, err = reader.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
	require.Equal(t, content, string(b))
	require.Equal(t, 2, reader.(*IPCReader).connectionCount())

	b = make([]byte, len(content))
	n, err = reader.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
	require.Equal(t, content, string(b))

	b = make([]byte, len(content))
	n, err = reader.Read(b)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)
}
