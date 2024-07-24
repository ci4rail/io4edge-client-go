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
	"log"
	"net"
	"time"

	"github.com/ci4rail/io4edge-client-go/transport"
)

// UDPConnection manages the UDP connection and sorts data for the simulated sockets
type UDPConnection struct {
	netudp  *net.UDPConn
	sockets []*UDPSocket
	lis     *UDPListener
}

// UDPListener implements the Listener interface for UDP sockets
type UDPListener struct {
	socket chan *UDPSocket
}

// UDPSocket implements the Transport interface for UDP sockets
// It simulates a a connection to each UDP counterpart that its api behaves like a TCP connection.
type UDPSocket struct {
	conn        *UDPConnection
	remoteAddr  *net.UDPAddr
	readTimeout time.Duration
	readData    chan []byte
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
		sockets: []*UDPSocket{},
		lis:     lis,
	}

	go c.readFromUDPConnection()

	return lis, nil
}

// WaitForUDPSocketConnect waits for a client to connect to the UDP socket and returns a UDPSocket object
// There is no timeout
func (l *UDPListener) WaitForUDPSocketConnect() (*UDPSocket, error) {
	log.Printf("Waiting for UDP Socket connection")
	// wait on channel for client connection
	s := <-l.socket
	log.Printf("UDP Socket connected")
	return s, nil
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
		remoteAddr:  addr,
		readTimeout: 0,
		readData:    make(chan []byte),
	}

	// create UDPConnection
	c := &UDPConnection{
		netudp:  conn,
		sockets: []*UDPSocket{s},
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
		log.Printf("Waiting for UDP packet")
		msg := make([]byte, 65507)
		_, remoteAddr, err := c.netudp.ReadFromUDP(msg)
		if err != nil {
			return err
		}

		log.Printf("Received: %s", msg)

		for _, s := range c.sockets {
			if s.remoteAddr.String() == remoteAddr.String() {
				log.Printf("Received on socket for remote %s", remoteAddr.String())
				go func() {
					s.readData <- msg
				}()
				continue CONNREAD
			}
		}

		if c.lis == nil {
			log.Printf("No listener -> client connection")
			// no listener -> only one socket allowed -> drop packet
			continue CONNREAD
		} else {
			log.Printf("Received on new socket for remote %s", remoteAddr.String())
			// create new socket
			s := &UDPSocket{
				remoteAddr: remoteAddr,
				readData:   make(chan []byte),
			}
			go func() {
				s.readData <- msg
			}()
			c.sockets = append(c.sockets, s)
			c.lis.socket <- s
		}
	}
}

// Close closes the UDP socket
func (s *UDPSocket) Close() error {
	// remove socket from connection
	// todo mutex?
	for i, socket := range s.conn.sockets {
		if socket == s {
			s.conn.sockets = append(s.conn.sockets[:i], s.conn.sockets[i+1:]...)
			break
		}
	}

	if len(s.conn.sockets) == 0 {
		return s.conn.netudp.Close()
	}
	return nil
}

// Write writes data to the UDP socket
func (s *UDPSocket) Write(p []byte) (n int, err error) {
	n, err = s.conn.netudp.WriteToUDP(p, s.remoteAddr)
	return n, err
}

// Read reads data from the UDP socket
func (s *UDPSocket) Read(p []byte) (n int, err error) {
	select {
	case data := <-s.readData:
		n = copy(p, data)
		return n, nil
	case <-time.After(s.readTimeout):
		return 0, transport.ErrTimeout
	}
}
