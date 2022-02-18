/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

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

// Package mvbsniffer provides the API for the io4edge mvbSniffer functionblock
package mvbsniffer

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/mvbSniffer/go/mvbSniffer/v1"
)

// Client represents a client for the mvbSniffer Module
type Client struct {
	fbClient *functionblock.Client
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new mvbSniffer client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_mvbSniffer._tcp).
// The timeout specifies the maximal time waiting for a service to show up. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_mvbSniffer._tcp", timeout)

	if err != nil {
		return nil, err
	}
	return &Client{
		fbClient: io4eClient,
	}, nil
}

// StartStream starts the stream on this connection.
func (c *Client) StartStream(genericConfig *functionblock.StreamConfiguration) error {
	err := c.fbClient.StartStream(genericConfig, &fspb.StreamControlStart{})
	if err != nil {
		return err
	}
	return nil
}

// StopStream stops the stream on this connection
func (c *Client) StopStream() error {
	return c.fbClient.StopStream()
}

// SendPattern sends a string to the self-test generator
func (c *Client) SendPattern(s string) error {
	fsCmd := &fspb.FunctionControlSet{
		GeneratorPattern: s,
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// ReadStream reads the next stream data object from the buffer.
//
// Returns the meta data and the unmarshalled function specific stream data
func (c *Client) ReadStream(timeout time.Duration) (*StreamData, error) {
	genericSD, err := c.fbClient.ReadStream(timeout)
	if err != nil {
		return nil, err
	}

	fsSD := new(fspb.StreamData)
	if err := genericSD.FSData.UnmarshalTo(fsSD); err != nil {
		return nil, errors.New("can't unmarshall samples")
	}

	sd := &StreamData{
		StreamDataMeta: genericSD.StreamDataMeta,
		FSData:         fsSD,
	}
	return sd, nil
}
