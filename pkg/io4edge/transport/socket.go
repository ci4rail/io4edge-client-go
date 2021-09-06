package transport

import (
	"net"
)

// NewSocketListener creates a Listener on a socket on a TCP socket
// port should be the port to listen to e.g. ":9999"
// pass the Listener to WaitForConnect
func NewSocketListener(port string) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// WaitForSocketConnect waits for a client to connect to the TCP socket and returns the TCP connection
// There is no timeout
func WaitForSocketConnect(l *net.TCPListener) (*net.TCPConn, error) {

	conn, err := l.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewSocketConnection connects to a TCP server at address and returns the TCP connection
func NewSocketConnection(address string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewMsgStreamFromConnection creates a message stream from TCP connection
func NewMsgStreamFromConnection(conn *net.TCPConn) (*FramedStream, error) {
	return &FramedStream{
		trans: conn,
	}, nil
}
