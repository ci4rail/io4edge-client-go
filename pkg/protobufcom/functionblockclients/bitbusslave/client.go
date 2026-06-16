/*
Copyright © 2026 Ci4Rail GmbH <engineering@ci4rail.com>

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

// Package bitbusslave provides the API for the io4edge bitbusSlave functionblock
package bitbusslave

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	fspb "github.com/ci4rail/io4edge_api/bitbusSlave/go/bitbusSlave/v1"
)

// Client represents a client for the bitbusSlave module.
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration().
type ConfigOption func(*fspb.ConfigurationSet)

// State describes the current slave state.
type State struct {
	Mode             fspb.SlaveMode
	HavePendingTxMsg bool
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data.
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new bitbusSlave client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_bitbusSlave._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_bitbusSlave._tcp", timeout)
	if err != nil {
		return nil, err
	}
	return &Client{
		fbClient: io4eClient,
	}, nil
}

// Close terminates the underlying connection to the functionblock.
func (c *Client) Close() {
	c.fbClient.Close()
}

// WithSlaveAddress may be passed to UploadConfiguration.
func WithSlaveAddress(slaveAddress int32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.SlaveAddress = slaveAddress
	}
}

// WithMaxFrameLength may be passed to UploadConfiguration.
func WithMaxFrameLength(maxFrameLength int32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.MaxFrameLength = maxFrameLength
	}
}

// WithAppWDTimeoutMS may be passed to UploadConfiguration.
func WithAppWDTimeoutMS(timeoutMS int32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.AppWdTimeoutMs = timeoutMS
	}
}

// WithIdleResponse may be passed to UploadConfiguration.
func WithIdleResponse(idleResponse []byte) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.IdleResponse = idleResponse
	}
}

// WithBaud62500 may be passed to UploadConfiguration.
func WithBaud62500(enable bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.Baud_62500 = enable
	}
}

// UploadConfiguration configures the bitbusSlave function block.
// Arguments may be one or more of the following functions:
//   - WithSlaveAddress (mandatory)
//   - WithMaxFrameLength (optional, default 255)
//   - WithAppWDTimeoutMS (optional, default 5000)
//   - WithIdleResponse (optional, default nil)
//   - WithBaud62500 (optional, default false)
func (c *Client) UploadConfiguration(opts ...ConfigOption) error {
	fsCmd := &fspb.ConfigurationSet{}
	for _, opt := range opts {
		opt(fsCmd)
	}

	_, err := c.fbClient.UploadConfiguration(fsCmd)
	return err
}

// SetPreparedTxMsg sets the next application tx message sent when the master addresses this slave.
// The firmware rejects it if
//   - the length of bitbusInformation exceeds the configured MaxFrameLength
//   - the last prepared message has not been sent yet (i.e. there is still a pending tx message)
func (c *Client) SetPreparedTxMsg(bitbusInformation []byte) error {
	fsCmd := &fspb.FunctionControlSet{
		Type: &fspb.FunctionControlSet_TxMsg{
			TxMsg: &fspb.PreparedTxMsg{
				BitbusInformation: bitbusInformation,
			},
		},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// State reads the current mode and pending-tx status from the slave.
func (c *Client) State() (*State, error) {
	any, err := c.fbClient.FunctionControlGet(&fspb.FunctionControlGet{})
	if err != nil {
		return nil, err
	}

	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return nil, err
	}

	return &State{
		Mode:             res.Mode,
		HavePendingTxMsg: res.HavePendingTxMsg,
	}, nil
}

// StartStream starts the stream on this connection.
// Arguments may be one or more of the functionblock.WithXXX() functions that
// may be passed to functionblock.StartStream().
// Options that are not specified take default values.
func (c *Client) StartStream(opts ...functionblock.StreamConfigOption) error {
	err := c.fbClient.StartStream(opts, &fspb.StreamControlStart{})
	if err != nil {
		return err
	}
	return nil
}

// StopStream stops the stream on this connection.
func (c *Client) StopStream() error {
	return c.fbClient.StopStream()
}

// ReadStream reads the next stream data object from the buffer.
//
// Returns the meta data and the unmarshalled function specific stream data.
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
