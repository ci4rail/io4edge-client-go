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

package socket

import (
	"net"
	"time"

	"github.com/ci4rail/io4edge-client-go/transport"
	log "github.com/sirupsen/logrus"
)

// ReadBufferSize is the size of the buffer for incoming messages of one connection
var ReadBufferSize = 10

// UDPConnection manages the UDP connection and sorts data for the simulated sockets
type UDPConnection struct {
	netudp  *net.UDPConn
	sockets map[string]*UDPSocket
	lis     *UDPListener
}

// UDPListener implements the Listener interface for UDP sockets
type UDPListener struct {
	socket chan *UDPSocket
	conn   *UDPConnection
}

// UDPSocket implements the Transport interface for UDP sockets
// It simulates a connection to each UDP counterpart that its api behaves like a TCP connection.
type UDPSocket struct {
	conn         *UDPConnection
	remoteAddr   *net.UDPAddr
	readDeadline time.Time
	readData     chan []byte
}

// NewUDPSocketListener creates a Listener on a socket on a UDP socket
// port should be the port to listen to e.g. ":9999"
func NewUDPSocketListener(port string) (*UDPListener, error) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	// create UDPConnection and UDPListener
	lis := &UDPListener{
		socket: make(chan *UDPSocket),
	}

	c := &UDPConnection{
		netudp:  conn,
		sockets: map[string]*UDPSocket{},
		lis:     lis,
	}

	lis.conn = c

	go c.readFromUDPConnection()

	return lis, nil
}

// WaitForUDPSocketConnect waits for a client to connect to the UDP socket and returns a UDPSocket object
// There is no timeout
func (l *UDPListener) WaitForUDPSocketConnect() (*UDPSocket, error) {
	log.Debugf("Waiting for UDP Socket connection")
	// wait on channel for client connection
	s := <-l.socket
	log.Debugf("UDP Socket connected")
	return s, nil
}

// Close closes the UDP Connection
func (l *UDPListener) Close() error {
	return l.conn.netudp.Close()
}

// NewUDPSocketConnection connects to a UDP server at address and returns a UDPSocket object
func NewUDPSocketConnection(address string) (*UDPSocket, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	// create UDPSocket
	s := &UDPSocket{
		remoteAddr:   addr,
		readDeadline: time.Time{},
		readData:     make(chan []byte, ReadBufferSize),
	}

	// create UDPConnection
	c := &UDPConnection{
		netudp:  conn,
		sockets: map[string]*UDPSocket{addr.String(): s},
		lis:     nil,
	}

	// set connection in socket
	s.conn = c

	go c.readFromUDPConnection()

	return s, nil
}

func (c *UDPConnection) readFromUDPConnection() error {
CONNREAD:
	for {
		log.Debugf("Waiting for UDP packet")
		tmp := make([]byte, 65507)
		n, remoteAddr, err := c.netudp.ReadFromUDP(tmp)
		if err != nil {
			return err
		}

		log.Debugf("Received: %s size: %d", tmp, n)
		// recreate message buffer with correct size
		msg := tmp[:n]

		// check if socket exists
		s, ok := c.sockets[remoteAddr.String()]
		if ok {
			log.Debugf("Received on socket for remote %s", remoteAddr.String())
			go func() {
				select {
				case s.readData <- msg:
				default:
					log.Debugf("Message buffer full -> drop oldest message")
					<-s.readData
					s.readData <- msg
				}
			}()
			continue CONNREAD
		}

		if c.lis == nil {
			log.Debugf("No listener -> client connection")
			// no listener -> only one socket allowed -> drop packet
			continue CONNREAD
		} else {
			log.Debugf("Received on new socket for remote %s", remoteAddr.String())
			// create new socket
			s := &UDPSocket{
				conn:         c,
				remoteAddr:   remoteAddr,
				readDeadline: time.Time{},
				readData:     make(chan []byte, ReadBufferSize),
			}
			go func() {
				select {
				case s.readData <- msg:
				default:
					log.Debugf("Message buffer full -> drop oldest message")
					<-s.readData
					s.readData <- msg
				}
			}()
			c.sockets[remoteAddr.String()] = s
			c.lis.socket <- s
		}
	}
}

// SetReadDeadline sets the read deadline for the UDP socket
func (s *UDPSocket) SetReadDeadline(t time.Time) error {
	s.readDeadline = t
	return nil
}

// Close closes the UDP socket
func (s *UDPSocket) Close() error {
	// remove socket from connection
	delete(s.conn.sockets, s.remoteAddr.String())

	// close connection if it is no server
	if s.conn.lis == nil {
		return s.conn.netudp.Close()
	}
	return nil
}

// Write writes data to the UDP socket
func (s *UDPSocket) Write(p []byte) (n int, err error) {
	if s.conn.lis == nil {
		// client connection
		n, err = s.conn.netudp.Write(p)
	} else {
		n, err = s.conn.netudp.WriteToUDP(p, s.remoteAddr)
	}

	return n, err
}

// Read reads data from the UDP socket
func (s *UDPSocket) Read(p []byte) (n int, err error) {
	// check if socket still exists
	_, ok := s.conn.sockets[s.remoteAddr.String()]
	if !ok {
		return 0, transport.ErrClosed
	}

	if s.readDeadline == (time.Time{}) {
		data := <-s.readData
		n = copy(p, data)
		return n, nil
	}
	timeout := time.Until(s.readDeadline)
	select {
	case data := <-s.readData:
		n = copy(p, data)
		// log.Debugf("Read data: %v", p)
		return n, nil
	case <-time.After(timeout):
		return 0, transport.ErrTimeout
	}

}
