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

// Package analogintypeb provides the API for the io4edge analogInTypeB functionblock
package analogintypeb

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/functionblock"
	fspb "github.com/ci4rail/io4edge_api/analogInTypeB/go/analogInTypeB/v1"
)

// Client represents a client for the analogInTypeB Module
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration()
type ConfigOption func(*fspb.ConfigurationSet)

// Configuration describes the current configuration of the analogInTypeB function.
// Returned by DownloadConfiguration()
type Configuration struct {
	// ChannelConfig describes the configuration of each channel
	ChannelConfig []*fspb.ChannelConfig
}

// ChannelGroupSpecification describes a group of channels that have the same configuration options
type ChannelGroupSpecification struct {
	Channels             []uint32  // channels in that group, e.g. [0,1] or [2,3,4,5]
	SupportedSampleRates []float32 // in Hz
	SupportedGains       []uint32  // e.g. 1, 2, 4, 8
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new analogInTypeB client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_analogInTypeB._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_analogInTypeB._tcp", timeout)

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

// UploadConfiguration configures the analogInTypeB function block
// Arguments may be one or more of the following functions:
//   - WithChannelConfig()
//
// Options that are not specified take default values.
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

// DownloadConfiguration reads the analogInTypeB function block configuration
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

// Describe reads the channels specification
func (c *Client) Describe() ([]*fspb.ChannelGroupSpecification, error) {
	fsCmd := &fspb.ConfigurationDescribe{}
	any, err := c.fbClient.Describe(fsCmd)
	if err != nil {
		return nil, err
	}
	res := new(fspb.ConfigurationDescribeResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return nil, err
	}
	return res.ChannelSpecification, err
}

// Values reads the current analog input level of all channels
//
// range -1 .. +1 (for min/max voltage or current)
func (c *Client) Values() ([]float32, error) {
	any, err := c.fbClient.FunctionControlGet(&fspb.FunctionControlGet{})
	if err != nil {
		return []float32{}, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return []float32{}, err
	}
	return res.Value, nil
}

// StreamConfigOption is a type to pass options to StartStream()
type StreamConfigOption func(*StreamConfiguration)

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	ChannelMask uint32
	FBOptions   []functionblock.StreamConfigOption
}

// WithChannelMask may be passed to StartStream.
//
// channelMask defines the channels to be streamed.
// Bit 0 corresponds to channel 0, bit 1 to channel 1, etc.
// E.g. channelMask = 5 (binary 0101) streams channels 0 and 2.
func WithChannelMask(channelMask uint32) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.ChannelMask = channelMask
	}
}

// WithFBStreamOption may be passed to StartStream.
//
// opt is one of the functions that may be passed to functionblock.StartStream, e.g. WithBucketSamples()
func WithFBStreamOption(opt functionblock.StreamConfigOption) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.FBOptions = append(c.FBOptions, opt)
	}
}

// StartStream starts the stream on this connection.
// Arguments may be one or more of the following functions:
//   - WithChannelMask
//   - WithFBStreamOption(functionblock.WithXXXX(...))
//
// Options that are not specified take default values.
func (c *Client) StartStream(opts ...StreamConfigOption) error {
	config := &StreamConfiguration{
		ChannelMask: 0xffffffff,
	}
	for _, opt := range opts {
		opt(config)
	}

	fsCmd := &fspb.StreamControlStart{
		ChannelMask: uint32(config.ChannelMask),
	}
	err := c.fbClient.StartStream(config.FBOptions, fsCmd)
	if err != nil {
		return err
	}
	return nil
}

// StopStream stops the stream on this connection
func (c *Client) StopStream() error {
	return c.fbClient.StopStream()
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
