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

// Package bitbussniffer provides the API for the io4edge bitbusSniffer function block.
package bitbussniffer

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	fspb "github.com/ci4rail/io4edge_api/bitbusSniffer/go/bitbusSniffer/v1"
)

// Client represents a client for the bitbusSniffer module.
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration.
type ConfigOption func(*fspb.ConfigurationSet)

// StreamData contains the stream metadata and the unmarshalled function-specific data.
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new bitbusSniffer client from addrOrService.
// If addrOrService has the form "host:port", it creates the client from that address.
// Otherwise, addrOrService is treated as the instance name of an mDNS service (without
// _io4edge_bitbusSniffer._tcp). timeout is the maximum time to wait for a service to
// appear. A zero timeout uses the default and is ignored for "host:port" addresses.
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	fbClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_bitbusSniffer._tcp", timeout)
	if err != nil {
		return nil, err
	}

	return &Client{fbClient: fbClient}, nil
}

// Close terminates the underlying connection to the function block.
func (c *Client) Close() {
	c.fbClient.Close()
}

// WithIgnoreCRC configures whether frames with an invalid CRC are captured.
func WithIgnoreCRC(ignore bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.IgnoreCrc = ignore
	}
}

// WithBaud62500 selects 62500 baud when enabled and 375000 baud otherwise.
func WithBaud62500(enable bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.Baud_62500 = enable
	}
}

// WithAddressFilter configures the address filter bit mask.
// The mask contains 32 bytes, with one bit for each Bitbus address.
func WithAddressFilter(addressFilter []byte) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.AddressFilter = addressFilter
	}
}

// WithMinFrameLength configures the minimum captured frame length.
func WithMinFrameLength(minFrameLength int32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.MinFrameLength = minFrameLength
	}
}

// WithPrepareSender enables support for sending frames.
func WithPrepareSender(enable bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.PrepareSender = enable
	}
}

// WithLoopbackEnable enables internal loopback and disables external bus activity.
func WithLoopbackEnable(enable bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.LoopbackEnable = enable
	}
}

// WithFullDuplex keeps the receiver enabled while the sender is transmitting.
func WithFullDuplex(enable bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.FullDuplex = enable
	}
}

// UploadConfiguration configures the bitbusSniffer function block.
// Arguments may be one or more of the following functions:
//   - WithIgnoreCRC
//   - WithBaud62500
//   - WithAddressFilter
//   - WithMinFrameLength
//   - WithPrepareSender
//   - WithLoopbackEnable
//   - WithFullDuplex
//
// Unspecified options use their protobuf zero values.
func (c *Client) UploadConfiguration(opts ...ConfigOption) error {
	fsCmd := &fspb.ConfigurationSet{}
	for _, opt := range opts {
		opt(fsCmd)
	}

	_, err := c.fbClient.UploadConfiguration(fsCmd)
	return err
}

// SendFrame sends a Bitbus frame. The frame contains the address byte, control
// byte, and optional information bytes. Sending requires a configuration with
// WithPrepareSender(true).
func (c *Client) SendFrame(bitbusFrame []byte) error {
	fsCmd := &fspb.FunctionControlSet{BitbusFrame: bitbusFrame}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// StartStream starts the stream on this connection.
// Arguments may be functionblock stream configuration options. Unspecified
// options retain their defaults.
func (c *Client) StartStream(opts ...functionblock.StreamConfigOption) error {
	return c.fbClient.StartStream(opts, &fspb.StreamControlStart{})
}

// StopStream stops the stream on this connection.
func (c *Client) StopStream() error {
	return c.fbClient.StopStream()
}

// ReadStream reads the next stream data object from the buffer.
//
// Each entry in StreamData.FSData.Samples contains:
//   - Timestamp: microseconds since device startup; it is not synchronized with
//     the client's clock.
//   - Flags: a bit mask composed of fspb.Sample_Flags values. Sample_bad_crc,
//     Sample_frames_lost, and Sample_buf_overrun indicate a bad CRC, frames lost
//     between the FPGA and microcontroller, and a buffer overrun between the
//     microcontroller and host, respectively.
//   - BitbusFrame: the received frame. Byte 0 is the address, byte 1 is the
//     control byte, and bytes 2 through n are the information field.
func (c *Client) ReadStream(timeout time.Duration) (*StreamData, error) {
	genericSD, err := c.fbClient.ReadStream(timeout)
	if err != nil {
		return nil, err
	}

	fsSD := new(fspb.StreamData)
	if err := genericSD.FSData.UnmarshalTo(fsSD); err != nil {
		return nil, fmt.Errorf("unmarshal bitbusSniffer stream data: %w", err)
	}

	return &StreamData{
		StreamDataMeta: genericSD.StreamDataMeta,
		FSData:         fsSD,
	}, nil
}
