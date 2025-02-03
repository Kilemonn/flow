package ipc

import (
	"io"
	"time"

	"github.com/Kilemonn/flow/queuedreader"
	ipcClient "github.com/Kilemonn/go-ipc/client"
	ipcServer "github.com/Kilemonn/go-ipc/server"
)

const (
	IPCReadDeadline = 10 * time.Millisecond
)

type IPCReader struct {
	server  ipcServer.IPCServer
	clients []ipcClient.IPCClient
}

func (r IPCReader) Close() (err error) {
	for _, c := range r.clients {
		e := c.Close()
		if e != nil && err == nil {
			err = e
		}
	}

	e := r.server.Close()
	if e != nil && err == nil {
		err = e
	}
	return err
}

// Get the amount of active connections
func (r IPCReader) connectionCount() int {
	return len(r.clients)
}

// Check if any incoming connections are pending to be accepted.
// This is naturally blocking, so there is a deadline set for [IPCReadDeadline]
// before this function returns with no accepted connections.
func (r *IPCReader) acceptWaitingConnections() {
	for {
		client, err := r.server.Accept(IPCReadDeadline)
		if err != nil {
			return
		}
		client.ReadTimeout = IPCReadDeadline
		r.clients = append(r.clients, client)
	}
}

func (r *IPCReader) Read(b []byte) (n int, err error) {
	r.acceptWaitingConnections()

	q := queuedreader.NewQueuedReader(r.clients)
	return q.Read(b)
}

func NewIPCReader(ipcChannelName string) (io.ReadCloser, error) {
	server, err := ipcServer.NewIPCServer(ipcChannelName, &ipcServer.IPCServerConfig{Override: true})
	if err != nil {
		return nil, err
	}

	return &IPCReader{server: *server}, nil
}
