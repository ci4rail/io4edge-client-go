/*
Copyright © 2024 Ci4Rail GmbH <engineering@ci4rail.com>

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

// Package pbcore provides the API for the io4edge core functions
// i.e. firmware and hardware id management
// for devices using the classic socket transport.
package pbcore

import (
	"errors"
	"time"

	pbchannelclient "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/channel"
	api "github.com/ci4rail/io4edge_api/io4edge/go/core_api/v1alpha2"
)

// Client represents a client for the io4edge core function
type Client struct {
	funcClient *pbchannelclient.Client
}

// NewClient creates a new client for the core function from the function client c
func NewClient(c *pbchannelclient.Client) *Client {
	return &Client{funcClient: c}
}

// NewClientFromSocketAddress creates a new core function client from a socket with the specified address.
func NewClientFromSocketAddress(address string) (*Client, error) {
	c, err := pbchannelclient.NewClientFromSocketAddress(address)
	if err != nil {
		return nil, errors.New("can't create function client: " + err.Error())
	}
	coreClient := NewClient(c)

	return coreClient, nil
}

// NewClientFromService creates a new core function client from a socket with a address, which was acquired from the specified service.
// The timeout specifies the maximal time waiting for a service to show up.
func NewClientFromService(serviceAddr string, timeout time.Duration) (*Client, error) {
	c, err := pbchannelclient.NewClientFromService(serviceAddr, timeout)

	if err != nil {
		return nil, err
	}
	coreClient := NewClient(c)

	return coreClient, nil
}

// Command issues a command cmd to a core function channel, waits for the devices response and returns it in res
func (c *Client) Command(cmd *api.CoreCommand, res *api.CoreResponse, timeout time.Duration) error {
	err := c.funcClient.Command(cmd, res, timeout)
	if err != nil {
		return err
	}
	if res.Status != api.Status_OK {
		return errors.New("Device reported error status: " + res.Status.String())
	}
	return err
}
