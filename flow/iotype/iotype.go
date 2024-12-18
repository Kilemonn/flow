package iotype

// We might still need this for now...
type IOType int

const (
	StdIO IOType = iota
	File
	Port
	Socket
	IPC
)
