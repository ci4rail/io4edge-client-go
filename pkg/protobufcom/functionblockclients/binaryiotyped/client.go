/*
Copyright © 2025 Ci4Rail GmbH <engineering@ci4rail.com>

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

// Package binaryiotyped provides the API for the io4edge binaryIoTypeD functionblock
package binaryiotyped

import (
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	fspb "github.com/ci4rail/io4edge_api/binaryIoTypeD/go/binaryIoTypeD/v1"
)

// Client represents a client for the binaryIoTypeD Module
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration()
type ConfigOption func(*fspb.ConfigurationSet)

// Configuration describes the current configuration of the binaryIoTypeD function.
// Returned by DownloadConfiguration()
type Configuration struct {
	// ChannelConfig describes the configuration of each channel
	ChannelConfig []*fspb.ChannelConfig
}

// NewClientFromUniversalAddress creates a new binaryIoTypeD client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_binaryIoTypeD._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_binaryIoTypeD._tcp", timeout)

	if err != nil {
		return nil, err
	}
	return &Client{
		fbClient: io4eClient,
	}, nil
}

// Close terminates the underlying connection to the functionblock
func (c *Client) Close() {
	c.fbClient.Close()
}

// WithChannelConfig may be passed to UploadConfiguration.
// each entry describes the configuration of one channel.
// Undescribed channels remain unchanged.
func WithChannelConfig(ch []*fspb.ChannelConfig) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.ChannelConfig = ch
	}
}

// UploadConfiguration configures the binaryIoTypeD function block.
// Arguments may be one or more of the following functions:
//   - WithChannelConfig
//
// Options that are not specified remain unchanged.
func (c *Client) UploadConfiguration(opts ...ConfigOption) error {

	// set defaults
	fsCmd := &fspb.ConfigurationSet{
		ChannelConfig: []*fspb.ChannelConfig{},
	}

	for _, opt := range opts {
		opt(fsCmd)
	}
	_, err := c.fbClient.UploadConfiguration(fsCmd)
	return err
}

// DownloadConfiguration reads the binaryIoTypeD function block configuration
func (c *Client) DownloadConfiguration() (*Configuration, error) {
	fsCmd := &fspb.ConfigurationGet{}
	any, err := c.fbClient.DownloadConfiguration(fsCmd)
	if err != nil {
		return nil, err
	}
	res := new(fspb.ConfigurationGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return nil, err
	}
	cfg := &Configuration{
		ChannelConfig: res.ChannelConfig,
	}
	return cfg, err
}

// SetOutput sets a single output channel
// a "true" state turns on the output switch, a "false" state turns it off.
func (c *Client) SetOutput(channel int, state bool) error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_Single{
			Single: &fspb.SetSingle{
				Channel: uint32(channel),
				State:   state,
			},
		},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// SetOutputs sets all or a group of output channels
//
// states: binary coded map of outputs. 0 means switch off, 1 means switch on, LSB is Channel0
//
// mask: binary coded map of outputs to be set. 0 means do not change, 1 means change, LSB is Channel0
func (c *Client) SetOutputs(states uint32, mask uint32) error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_All{
			All: &fspb.SetAll{
				Values: states,
				Mask:   mask,
			},
		},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// Inputs gets the state of all channels, regardless whether they are configured as input or output.
//
// Each bit in the returned "inputs" corresponds to one channel, bit0 being channel 0.
// The bit is false if the pin level is inactive, or true if active.
// diag is a slice with bitfields containing diagnostic bits for each channel.Each bit in the returned state corresponds to one channel, bit0 being channel 0.
func (c *Client) Inputs() (states uint32, diag []uint32, err error) {
	any, err := c.fbClient.FunctionControlGet(&fspb.FunctionControlGet{})
	if err != nil {
		return 0, nil, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return 0, nil, err
	}
	return res.Inputs, res.Diag, nil
}
