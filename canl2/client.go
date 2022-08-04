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

// Package canl2 provides the API for the io4edeg template functionblock
package canl2

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/canL2/go/canL2/v1alpha1"
)

// Client represents a client for the canL2
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration()
type ConfigOption func(*fspb.ConfigurationSet)

// Configuration describes the current configuration of the canL2 function.
// Returned by DownloadConfiguration()
type Configuration struct {
	BitRate     uint32
	SamplePoint float32
	SJW         uint8
	ListenOnly  bool
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new canL2 client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_canL2._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_canL2._tcp", timeout)

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

// WithBitRate may be passed to UploadConfiguration.
//
// bitRate defines the bit rate in bis/s
func WithBitRate(bitRate uint32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.Baud = bitRate
	}
}

// WithSamplePoint may be passed to UploadConfiguration.
//
// Sample Point from 0.0-1.0 - basis to calculate tseg1 and tseg2
func WithSamplePoint(samplePoint float32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.SamplePoint = int32(samplePoint * 1000)
	}
}

// WithSJW may be passed to UploadConfiguration.
//
// Synchronization Jump Width
func WithSJW(sjw uint8) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.Sjw = int32(sjw)
	}
}

// WithListenOnly may be passed to UploadConfiguration.
//
// if activated it is not possible to send frames to the bus (and no ACK is sent by the CAN controller)
func WithListenOnly(listenOnly bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.ListenOnly = listenOnly
	}
}

// UploadConfiguration configures the analogInTypeA function block
// Arguments may be one or more of the following functions:
//  - WithBitRate
//  - WithSamplePoint
//  - WithSJW
//  - WithListenOnly
// Options that are not specified take default value (Bitrate 500000, SamplePoint 0.8, SJW 1, ListenOnly false)
func (c *Client) UploadConfiguration(opts ...ConfigOption) error {
	fsCmd := &fspb.ConfigurationSet{
		Baud:        500000,
		SamplePoint: 800,
		Sjw:         1,
		ListenOnly:  false,
	}

	for _, opt := range opts {
		opt(fsCmd)
	}

	_, err := c.fbClient.UploadConfiguration(fsCmd)
	return err
}

// DownloadConfiguration reads the template function block configuration
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
		BitRate:     res.Baud,
		SamplePoint: float32(res.SamplePoint / 1000),
		SJW:         uint8(res.Sjw),
		ListenOnly:  res.ListenOnly,
	}
	return cfg, err
}

// SendFrames sends a slice of frames to the CAN bus
// if the queue on the device is not large enough to contain all frames,
// send nothing and return temporarily unavailable error
func (c *Client) SendFrames(frames []*fspb.Frame) error {
	fsCmd := &fspb.FunctionControlSet{
		Frame: frames,
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// GetCtrlState returns the current state of the CAN controller
// The returned values is one of the following:
//  fspb == "github.com/ci4rail/io4edge_api/canL2/go/canL2/v1alpha1"
// 	fspb.ControllerState_CAN_BUS_OFF
// 	fspb.ControllerState_CAN_ERROR_PASSIVE
// 	fspb.ControllerState_CAN_ERROR_WARNING
func (c *Client) GetCtrlState() (uint32, error) {
	fsCmd := &fspb.FunctionControlGet{}
	any, err := c.fbClient.FunctionControlGet(fsCmd)
	if err != nil {
		return 0, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return 0, err
	}
	return uint32(res.ControllerState), nil
}

// StreamConfigOption is a type to pass options to StartStream()
type StreamConfigOption func(*StreamConfiguration)

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	acceptanceCode uint32
	acceptanceMask uint32
	FBOptions      []functionblock.StreamConfigOption
}

// WithFilter may be passed to StartStream.
//
// acceptanceCode and acceptanceMask define the filter
func WithFilter(acceptanceCode uint32, acceptanceMask uint32) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.acceptanceCode = acceptanceCode
		c.acceptanceMask = acceptanceMask
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
//  - WithFilter
//  - WithFBStreamOption(functionblock.WithXXXX(...))
// Options that are not specified take default values.
func (c *Client) StartStream(opts ...StreamConfigOption) error {
	config := &StreamConfiguration{
		acceptanceCode: 0x0,
		acceptanceMask: 0x0,
	}
	for _, opt := range opts {
		opt(config)
	}
	fsCmd := &fspb.StreamControlStart{
		AcceptanceCode: config.acceptanceCode,
		AcceptanceMask: config.acceptanceMask,
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

// ReadStream reads the next stream data object from the buffer
// returns the meta data and the unmarshalled function specific stream data
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
