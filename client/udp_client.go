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
package client

import (
	"errors"

	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
)

// NewUDPClientFromSocketAddress creates a new function client from a socket with the specified address.
func NewUDPClientFromSocketAddress(address string) (*Client, error) {
	return newUDPClientFromSocketAddress(address, NewFuncInfoDefault(address))
}

func newUDPClientFromSocketAddress(address string, funcInfo FunctionInfo) (*Client, error) {
	ch, err := createUDPChannel(address)
	if err != nil {
		return nil, err
	}
	c := NewClient(ch, funcInfo)

	return c, nil
}

func createUDPChannel(address string) (*Channel, error) {
	t, err := socket.NewUDPSocketConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ms := transport.NewFrameHandshakeFromTransport(t)
	ch := NewChannel(ms)
	return ch, nil
}
