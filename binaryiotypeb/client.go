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

// Package binaryiotypeb provides the API for the io4edge binaryIoTypeB functionblock
package binaryiotypeb

import (
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/binaryIoTypeB/go/binaryIoTypeB/v1alpha1"
)

// Client represents a client for the binaryIoTypeB Module
type Client struct {
	fbClient *functionblock.Client
}

// Description represents the describe response of the binaryIoTypeB function
type Description struct {
	Channels []*fspb.ChannelConfig
}

// NewClientFromUniversalAddress creates a new binaryIoTypeB client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_binaryIoTypeB._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_binaryIoTypeB._tcp", timeout)

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

// Describe reads the binaryIoTypeB function block description
func (c *Client) Describe() (*Description, error) {
	fsCmd := &fspb.ConfigurationDescribe{}
	any, err := c.fbClient.Describe(fsCmd)
	if err != nil {
		return nil, err
	}
	res := new(fspb.ConfigurationDescribeResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return nil, err
	}

	desc := &Description{
		Channels: res.ChannelConfig,
	}
	return desc, err
}

// SetOutput sets a single output channel
// a true activates the output
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

// SetAllOutputs sets all or a group of output channels
//
// states: binary coded map of outputs. 0 means deactivate output channel, 1 means activate output channel, LSB is Output0
//
// mask: defines which channels are affected by the set all command.
func (c *Client) SetAllOutputs(states uint8, mask uint8) error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_All{
			All: &fspb.SetAll{
				Values: uint32(states),
				Mask:   uint32(mask),
			},
		},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// Input reads the state of the input pin of a single channel.
//
// The returned value is false if the input channel is inactive, true if the input channel is active
func (c *Client) Input(channel int) (bool, error) {
	fsCmd := &fspb.FunctionControlGet{
		Type: &fspb.FunctionControlGet_Single{
			Single: &fspb.GetSingle{
				Channel: uint32(channel),
			},
		},
	}
	any, err := c.fbClient.FunctionControlGet(fsCmd)
	if err != nil {
		return false, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return false, err
	}
	return res.GetSingle().State, nil
}

// AllInputs reads the state of all input pins defined by mask.
//
// Each bit in the returned value corresponds to one channel, bit0 being channel 0.
// The bit is false if the input channel is inactive, true if the input channel is active
// Channels whose bit is cleared in mask are reported as 0
func (c *Client) AllInputs(mask uint8) (uint8, error) {
	fsCmd := &fspb.FunctionControlGet{
		Type: &fspb.FunctionControlGet_All{
			All: &fspb.GetAll{
				Mask: uint32(mask),
			},
		},
	}
	any, err := c.fbClient.FunctionControlGet(fsCmd)
	if err != nil {
		return 0, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return 0, err
	}
	return uint8(res.GetAll().Inputs), nil
}
