package transport

import (
	"net"
)

// NewListener creates a listener using the address and scheme
// Returns a listener and error
func NewListener(addr, scheme string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(scheme, addr)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP(scheme, tcpAddr)
	if err != nil {
		panic(err)
	}

	// We could here secure the connection via TLS...
	return l, nil
}
