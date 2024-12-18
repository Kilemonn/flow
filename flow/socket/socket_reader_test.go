package socket

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/Kilemonn/flow/flow/testutil"
	"github.com/stretchr/testify/require"
)

func TestUDPReadAndWrite(t *testing.T) {
	reader, err := NewUDPSocketReader("127.0.0.1", 12345)
	require.NoError(t, err)
	writer, err := NewUDPSocketWriter("127.0.0.1", 12345)
	require.NoError(t, err)

	content := "TestUDPReadAndWrite"
	n, err := writer.Write([]byte(content))
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	b := make([]byte, len(content))
	// Since data is ready to read we are expecting this to return immediately
	expectedReadTime := 0 * time.Microsecond
	require.Less(t, expectedReadTime, SocketReadDeadline)
	testutil.TakesAtleast(t, expectedReadTime, func() {
		n, err = reader.Read(b)
	})

	require.NoError(t, err)
	require.Equal(t, len(content), n)

	require.Equal(t, content, string(b))
}

// Make sure that when we perform a UDP read and there is no data, that we do not hang forever and that we wait for the
// provided [SocketReadDeadline] to pass before continuing.
func TestUDPRead_NoData(t *testing.T) {
	reader, err := NewUDPSocketReader("127.0.0.1", 44321)
	require.NoError(t, err)

	b := make([]byte, 10)
	testutil.TakesAtleast(t, SocketReadDeadline, func() {
		n, err := reader.Read(b)
		require.Equal(t, err, io.EOF)
		require.Equal(t, 0, n)
	})
}

// Ensure that when attempting to accept a TCP connection but there is none, that the deadline
// kicks in and stops infinitely blocking the call, and also make sure no new connection is added.
func TestTCPRead_NoNewConnections(t *testing.T) {
	reader, err := NewTCPSocketReader("127.0.0.1", 54541)
	require.NoError(t, err)

	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())
	testutil.TakesAtleast(t, SocketReadDeadline, func() {
		reader.(*TCPTimeoutReader).acceptWaitingConnections()
	})
	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())
}

// Ensure pending TCP connections are accepted and added to the conns list
func TestTCPAcceptNewConnections(t *testing.T) {
	reader, err := NewTCPSocketReader("127.0.0.1", 54540)
	require.NoError(t, err)

	_, err = NewTCPSocketWriter("127.0.0.1", 54540)
	require.NoError(t, err)

	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())
	testutil.TakesAtleast(t, 0*time.Millisecond, func() {
		reader.(*TCPTimeoutReader).acceptWaitingConnections()
	})
	require.Equal(t, 1, reader.(*TCPTimeoutReader).connectionCount())

	_, err = NewTCPSocketWriter("127.0.0.1", 54540)
	require.NoError(t, err)

	testutil.TakesAtleast(t, 0*time.Millisecond, func() {
		reader.(*TCPTimeoutReader).acceptWaitingConnections()
	})
	require.Equal(t, 2, reader.(*TCPTimeoutReader).connectionCount())

}

// Make sure the call to read will accept new connections even if there is no data waiting to be read
func TestTCPRead_NewConnectionsNoData(t *testing.T) {
	reader, err := NewTCPSocketReader("127.0.0.1", 54539)
	require.NoError(t, err)

	_, err = NewTCPSocketWriter("127.0.0.1", 54539)
	require.NoError(t, err)

	_, err = NewTCPSocketWriter("127.0.0.1", 54539)
	require.NoError(t, err)

	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())

	b := make([]byte, 10)
	n, err := reader.Read(b)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)
	require.Equal(t, 2, reader.(*TCPTimeoutReader).connectionCount())
}

// Ensure that we read data from multiple connections if they have data.
// Make sure once all data is exhausted from all connections that we get an EOF back.
func TestTCPRead_DataFromMultipleSockets(t *testing.T) {
	reader, err := NewTCPSocketReader("127.0.0.1", 53213)
	require.NoError(t, err)

	content := "TestTCPRead_DataFromMultipleSockets"

	writer1, err := NewTCPSocketWriter("127.0.0.1", 53213)
	require.NoError(t, err)
	n, err := writer1.Write([]byte(content))
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	writer2, err := NewTCPSocketWriter("127.0.0.1", 53213)
	require.NoError(t, err)
	n, err = writer2.Write([]byte(content))
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())

	b := make([]byte, len(content))
	n, err = reader.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
	require.Equal(t, content, string(b))
	require.Equal(t, 2, reader.(*TCPTimeoutReader).connectionCount())

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

// Ensure that we remove the connection if its remote peer has closed the connection.
func TestTCPRead_RemoteConnectionCloses(t *testing.T) {
	reader, err := NewTCPSocketReader("127.0.0.1", 63321)
	require.NoError(t, err)

	writer1, err := NewTCPSocketWriter("127.0.0.1", 63321)
	require.NoError(t, err)

	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())
	reader.(*TCPTimeoutReader).acceptWaitingConnections()
	require.Equal(t, 1, reader.(*TCPTimeoutReader).connectionCount())

	err = writer1.(*net.TCPConn).Close()
	require.NoError(t, err)

	b := make([]byte, 10)
	n, err := reader.Read(b)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)

	require.Equal(t, 0, reader.(*TCPTimeoutReader).connectionCount())
}
