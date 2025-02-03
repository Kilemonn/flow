package socket

import (
	"fmt"
	"io"
	"net"
	"net/netip"
	"strings"
)

func CreateSocketWriter(protocol string, addr string, port uint16) (io.WriteCloser, error) {
	if strings.ToLower(protocol) == "tcp" {
		return NewTCPSocketWriter(addr, port)
	} else if strings.ToLower(protocol) == "udp" {
		return NewUDPSocketWriter(addr, port)
	} else {
		return nil, fmt.Errorf("invalid protocol provided [%s]", protocol)
	}
}

func NewUDPSocketWriter(addr string, port uint16) (io.WriteCloser, error) {
	address, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(address, port))
	return net.DialUDP("udp", nil, udpAddr)
}

func NewTCPSocketWriter(addr string, port uint16) (io.WriteCloser, error) {
	address, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	tcpAddr := net.TCPAddrFromAddrPort(netip.AddrPortFrom(address, port))
	return net.DialTCP("tcp", nil, tcpAddr)
}
