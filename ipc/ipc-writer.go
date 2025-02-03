package ipc

import (
	"io"

	ipcClient "github.com/Kilemonn/go-ipc/client"
)

func NewIPCWriter(ipcChannelName string) (io.WriteCloser, error) {
	return ipcClient.NewIPCClient(ipcChannelName)
}
