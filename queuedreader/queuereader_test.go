package queuedreader

import (
	"io"
	"net"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Ensure the prehandler is set properly and called on Read().
func TestRead_PreReadHandlerCalled(t *testing.T) {
	called := false
	q := NewQueuedReader([]*os.File{os.Stdin})
	require.Nil(t, q.preReadFunc)
	q.SetPreReadHandlerFunc(func(r *os.File) {
		called = true
	})
	require.NotNil(t, q.preReadFunc)

	b := make([]byte, 1)
	n, err := q.Read(b)
	require.Error(t, err)
	require.Equal(t, 0, n)

	require.True(t, called)
}

// Ensure that on EOF that the EOF func is called.
func TestRead_OnEOF(t *testing.T) {
	called := false
	q := NewQueuedReader([]*os.File{os.Stdin})
	require.Nil(t, q.onEOFFunc)
	q.SetEOFHandlerFunc(func(i int, r *os.File) {
		called = true
	})
	require.NotNil(t, q.onEOFFunc)

	b := make([]byte, 1)
	n, err := q.Read(b)
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)

	require.True(t, called)
}

// Ensure that on a timeout error that the timeout handler is called
func TestRead_OnTimeout(t *testing.T) {
	address, err := netip.ParseAddr("127.0.0.1")
	require.NoError(t, err)
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(address, 12235))
	conn, err := net.ListenUDP("udp", udpAddr)
	require.NoError(t, err)
	defer conn.Close()

	called := false
	q := NewQueuedReader([]*net.UDPConn{conn})
	q.SetPreReadHandlerFunc(func(r *net.UDPConn) {
		r.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	})
	require.Nil(t, q.onTimeoutFunc)
	q.SetTimeoutHandlerFunc(func(i int, r *net.UDPConn) {
		called = true
	})
	require.NotNil(t, q.onTimeoutFunc)

	b := make([]byte, 1)
	n, err := q.Read(b)
	require.Error(t, err)
	require.Equal(t, 0, n)

	require.True(t, called)
}
