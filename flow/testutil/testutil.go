package testutil

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// WithTempFile helper function to create a temp file before calling a method that accepts the temp file name
// the temp file is removed after this function finishes
func WithTempFile(t *testing.T, testFunc func(string)) {
	temp, err := os.CreateTemp("", "*")
	require.Nil(t, err)
	defer os.Remove(temp.Name())

	testFunc(temp.Name())
}

// CaptureStdout captures and returns an os.File that contains all content written to stdout during the provided test function
// stdout is returned to normal after this function
func CaptureStdout(t *testing.T, testFunc func()) *os.File {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)
	defer writer.Close()

	// Revert stdout after the end of this function
	defer func(v *os.File) { os.Stdout = v }(os.Stdout)
	os.Stdout = writer

	testFunc()

	return reader
}

// WithBytesInStdIn pre-load stdin with the provided bytes before running the provided test
// reverts std in after the test is complete
func WithBytesInStdIn(t *testing.T, bytes []byte, testFunc func()) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	n, err := writer.Write(bytes)
	require.Nil(t, err)
	require.Equal(t, len(bytes), n)
	writer.Close()

	// Revert stdin after the end of this function
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = reader

	testFunc()
}

// TakesAtleast asserts that the provided func thats atleast the provided duration or longer to complete
func TakesAtleast(t *testing.T, duration time.Duration, testFunc func()) {
	start := time.Now()

	testFunc()

	diff := time.Since(start)
	require.GreaterOrEqual(t, diff, duration)
}

func GetUDPPort(conn *net.UDPConn) uint16 {
	if conn != nil {
		if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			return addr.AddrPort().Port()
		}
	}

	return 0
}

func GetTCPPort(conn *net.TCPListener) uint16 {
	if conn != nil {
		if addr, ok := conn.Addr().(*net.TCPAddr); ok {
			return addr.AddrPort().Port()
		}
	}

	return 0
}
