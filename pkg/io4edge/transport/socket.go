/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
