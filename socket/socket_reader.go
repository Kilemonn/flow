package socket

import (
	"fmt"
	"io"
	"net"
	"net/netip"
	"strings"
	"time"
)

const (
	SocketReadDeadline = 10 * time.Millisecond
)

func CreateSocketReader(protocol string, addr string, port uint16) (io.ReadCloser, error) {
	if strings.ToLower(protocol) == "tcp" {
		return NewTCPSocketReader(addr, port)
	} else if strings.ToLower(protocol) == "udp" {
		return NewUDPSocketReader(addr, port)
	} else {
		return nil, fmt.Errorf("invalid protocol provided [%s]", protocol)
	}
}

func NewUDPSocketReader(addr string, port uint16) (io.ReadCloser, error) {
	address, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(address, port))
	conn, err := net.ListenUDP("udp", udpAddr)
	return UDPTimeoutReader{Conn: conn}, err
}

func NewTCPSocketReader(addr string, port uint16) (io.ReadCloser, error) {
	address, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	tcpAddr := net.TCPAddrFromAddrPort(netip.AddrPortFrom(address, port))
	listener, err := net.ListenTCP("tcp", tcpAddr)
	return &TCPTimeoutReader{Listener: listener}, err
}
