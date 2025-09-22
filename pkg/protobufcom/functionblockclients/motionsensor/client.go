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

// Package motionsensor provides the API for the io4edge motion sensor function block
package motionsensor

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/functionblock"
	fspb "github.com/ci4rail/io4edge_api/motionSensor/go/motionSensor/v1"
)

// Client represents a client for the motionSensor
type Client struct {
	fbClient *functionblock.Client
}

// ConfigOption is a type to pass options to UploadConfiguration()
type ConfigOption func(*fspb.ConfigurationSet)

// Configuration describes the current configuration of the motionSensor function.
// Returned by DownloadConfiguration()
type Configuration struct {
	SampleRateMilliHz    uint32
	FullScaleG           int32
	HighPassFilterEnable bool
	BandWidthRatio       int32
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new motionSensor client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_motionSensor._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_motionSensor._tcp", timeout)

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

// WithSampleRate may be passed to UploadConfiguration.
//
// sampleReate is the desired data rate in 1/1000 Hz
func WithSampleRate(sampleRateMilliHz uint32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.SampleRateMilliHz = sampleRateMilliHz
	}
}

// WithFullScale may be passed to UploadConfiguration.
//
// fullscale is the full scale acceleration in g
func WithFullScale(fullScale int32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.FullScaleG = fullScale
	}
}

// WithHighPassFilterEnable may be passed to UploadConfiguration.
func WithHighPassFilterEnable(enable bool) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.HighPassFilterEnable = enable
	}
}

// WithBandWidthRatio may be passed to UploadConfiguration.
//
// band width of low/hig-hpass as ratio of sample_rate
// .e.g. select 2 when the filterbandwith shall be sample_rate/2
func WithBandWidthRatio(ratio int32) ConfigOption {
	return func(c *fspb.ConfigurationSet) {
		c.BandWidthRatio = ratio
	}
}

// UploadConfiguration configures the analogInTypeA function block
// Arguments may be one or more of the following functions:
//   - WithSampleRate
//   - WithFullScale
//   - WithHighPassFilterEnable
//   - WithBandWidthRatio
//
// Options that are not specified take default value (SampleRate 12.5Hz, FullScale:2g, HighPassFilter: disable, BandWidthRatio: 2)
func (c *Client) UploadConfiguration(opts ...ConfigOption) error {
	fsCmd := &fspb.ConfigurationSet{
		SampleRateMilliHz:    12500,
		FullScaleG:           2,
		HighPassFilterEnable: false,
		BandWidthRatio:       2,
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
		SampleRateMilliHz:    res.SampleRateMillihz,
		FullScaleG:           res.FullScaleG,
		HighPassFilterEnable: res.HighPassFilterEnable,
		BandWidthRatio:       res.BandWidthRatio,
	}
	return cfg, err
}

// StartStream starts the stream on this connection.
// Arguments may be one or more of the functionblock.WithXXX() functions that
// may be passed to functionblock.StartStream()
// Options that are not specified take default values.
func (c *Client) StartStream(opts ...functionblock.StreamConfigOption) error {
	err := c.fbClient.StartStream(opts, &fspb.StreamControlStart{})
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
