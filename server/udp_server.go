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

// Package client provides the API for io4edge I/O devices
package server

import (
	"errors"
	"log"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
)

// UDPServer represents a server for io4edge devices
type UDPServer struct {
	lis *socket.UDPListener
}

// NewServer creates a new function client from a socket with the specified address.
func NewServer(port string) (*UDPServer, error) {
	lis, err := socket.NewUDPSocketListener(port)
	if err != nil {
		log.Fatal("can't create listener: " + err.Error())
		return nil, errors.New("can't create listener: " + err.Error())
	}

	srv := &UDPServer{
		lis: lis,
	}

	return srv, nil
}

func (s *UDPServer) ManageConnections() (*client.Channel, error) {
	sock, err := s.lis.WaitForUDPSocketConnect()
	if err != nil {
		return nil, errors.New("Error reading message: " + err.Error())
	}

	fh := transport.NewFrameHandshakeFromTransport(sock)
	ch := client.NewChannel(fh)
	log.Printf("New channel created")

	return ch, nil
}

func (s *UDPServer) Close() {
	s.lis.Close()
}
