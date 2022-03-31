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

// Package binaryiotypea provides the API for the io4edge binaryIoTypeA functionblock
package binaryiotypea

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
)

// Client represents a client for the binaryIoTypeA Module
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration()
type ConfigOption func(*fspb.ConfigurationSet)

// Configuration describes the current configuration of the binaryIoTypeA function.
// Returned by DownloadConfiguration()
type Configuration struct {
	// InputFrittingMask reflects on which inputs the fritting pulses are enabled
	InputFrittingMask uint8
	// OutputWatchdogMask reflects on which outputs the watchdog is enabled
	OutputWatchdogMask uint8
	// OutputWatchdogTimeoutMs reflects the watchdog timeout in ms
	OutputWatchdogTimeoutMs uint32
}

// Description represents the describe response of the binaryIoTypeA function
type Description struct {
	NumberOfChannels int
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new binaryIoTypeA client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_binaryIoTypeA._tcp).
// The timeout specifies the maximal time waiting for a service to show up. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_binaryIoTypeA._tcp", timeout)

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

// WithInputFritting may be passed to UploadConfiguration.
// mask defines on which inputs the fritting pulses shall be enabled (bit0=first IO).
func WithInputFritting(mask uint8) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.OutputFrittingMask = uint32(mask)
	}
}

// WithOutputWatchdog may be passed to UploadConfiguration.
// mask defines to which outputs the watchdog shall apply (bit0=first IO).
// timeoutMs defines the watchdog timeout in ms, it's the same for all selected outputs
func WithOutputWatchdog(mask uint8, timoutMs uint32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.OutputWatchdogMask = uint32(mask)
		c.OutputWatchdogTimeout = timoutMs
	}
}

// UploadConfiguration configures the binaryIoTypeA function block.
// Arguments may be one or more of the following functions:
//  - WithOutputWatchdog
//  - WithInputFritting
// Options that are not specified take default values.
func (c *Client) UploadConfiguration(opts ...ConfigOption) error {

	// set defaults
	fsCmd := &fspb.ConfigurationSet{
		OutputFrittingMask:    uint32(0x00),
		OutputWatchdogMask:    uint32(0x00),
		OutputWatchdogTimeout: 0,
	}

	for _, opt := range opts {
		opt(fsCmd)
	}
	_, err := c.fbClient.UploadConfiguration(fsCmd)
	return err
}

// DownloadConfiguration reads the binaryIoTypeA function block configuration
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
		InputFrittingMask:       uint8(res.OutputFrittingMask),
		OutputWatchdogMask:      uint8(res.OutputWatchdogMask),
		OutputWatchdogTimeoutMs: res.OutputWatchdogTimeout,
	}
	return cfg, err
}

// Describe reads the binaryIoTypeA function block description
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
		NumberOfChannels: int(res.NumberOfChannels),
	}
	return desc, err
}

// SetOutput sets a single output channel
// a "true" state turns on the output switch
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

// ExitErrorState tries to recover the binary output controller from error state.
//
// The binary output controller enters error state when there is an overurrent condition for a long time.
//
// In the error state, no outputs can be set; inputs can still be read.
// This call tells the binary output controller to try again. This call does however not wait
// if the recovery was successful or not.
func (c *Client) ExitErrorState() error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_ExitError{},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// SetAllOutputs sets all or a group of output channels
//
// states: binary coded map of outputs. 0 means switch off, 1 means switch on, LSB is Output0
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
// The returned value is false if the input level is below switching threshold, true otherwise
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
// The bit is false if the input level is below switching threshold, true otherwise.
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

// StreamConfigOption is a type to pass options to StartStream()
type StreamConfigOption func(*StreamConfiguration)

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	ChannelFilterMask uint8
	FBOptions         []functionblock.StreamConfigOption
}

// WithChannelFilterMask may be passed to StartStream.
//
// channelFilterMask defines the watched channels. Only changes on those channels generate samples in the stream
func WithChannelFilterMask(channelFilterMask uint8) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.ChannelFilterMask = channelFilterMask
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
//  - WithChannelFilterMask
//  - WithFBStreamOption(functionblock.WithXXXX(...))
// Options that are not specified take default values.
func (c *Client) StartStream(opts ...StreamConfigOption) error {
	config := &StreamConfiguration{
		ChannelFilterMask: 0xff,
	}
	for _, opt := range opts {
		opt(config)
	}

	fsCmd := &fspb.StreamControlStart{
		ChannelFilterMask: uint32(config.ChannelFilterMask),
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
