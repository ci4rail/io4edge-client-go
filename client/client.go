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
	"time"

	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
	"google.golang.org/protobuf/proto"
)

// FunctionInfo is an interface to query properties of the io4edge function
type FunctionInfo interface {
	// NetAddress returns the IP address (or host name) and the default port of the function
	NetAddress() (host string, port int, err error)
	// FuncClass returns the class of the io4edge function: e.g. core/datalogger/controlio/ttynvt
	FuncClass() (class string, err error)
	// Security tells whether function channels use encryption (no/tls)
	Security() (security string, err error)
	// AuxPort returns the protocol of the aux port (tcp/udp) and the port
	// returns error if no aux port for function
	AuxPort() (protcol string, port int, err error)
	// AuxSchema returns the schema name of the aux channel
	// returns error if no aux port for function
	AuxSchemaID() (schemaID string, err error)
}

// Client represents a client for an io4edge function
type Client struct {
	Ch       *Channel
	FuncInfo FunctionInfo
}

// NewClient creates a new client for an io4edge function
func NewClient(c *Channel, funcInfo FunctionInfo) *Client {
	return &Client{Ch: c, FuncInfo: funcInfo}
}

// NewClientFromSocketAddress creates a new function client from a socket with the specified address.
func NewClientFromSocketAddress(address string) (*Client, error) {
	return newClientFromSocketAddress(address, NewFuncInfoDefault(address))
}

func newClientFromSocketAddress(address string, funcInfo FunctionInfo) (*Client, error) {
	t, err := socket.NewSocketConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ms := transport.NewFramedStreamFromTransport(t)
	ch := NewChannel(ms)
	c := NewClient(ch, funcInfo)

	return c, nil
}

// NewClientFromService creates a new function client from a socket with a address, which was acquired from the specified service.
// The timeout specifies the maximal time waiting for a service to show up.
func NewClientFromService(serviceAddr string, timeout time.Duration) (*Client, error) {
	instance, service, err := ParseInstanceAndService(serviceAddr)
	if err != nil {
		return nil, err
	}
	svcInfo, err := NewServiceInfo(instance, service, timeout)
	if err != nil {
		return nil, err
	}
	ipAddrPort := svcInfo.GetIPAddressPort()
	c, err := newClientFromSocketAddress(ipAddrPort, svcInfo)
	return c, err
}

// Command issues a command cmd to a channel, waits for the devices response and returns it in res
func (c *Client) Command(cmd proto.Message, res proto.Message, timeout time.Duration) error {
	err := c.Ch.WriteMessage(cmd)
	if err != nil {
		return err
	}
	err = c.Ch.ReadMessage(res, timeout)
	if err != nil {
		return errors.New("Failed to receive device response: " + err.Error())
	}
	return err
}

func (c *Client) ReadMessage(res proto.Message, timeout time.Duration) error {
	err := c.Ch.ReadMessage(res, timeout)
	if err != nil {
		return errors.New("Failed to read forever: " + err.Error())
	}
	return nil
}
