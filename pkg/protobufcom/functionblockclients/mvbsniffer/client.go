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
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/functionblock"
	fspb "github.com/ci4rail/io4edge_api/mvbSniffer/go/mvbSniffer/v1"
)

// Client represents a client for the mvbSniffer Module
type Client struct {
	fbClient *functionblock.Client
}

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.TelegramCollection
}

// FilterMask defines a specific filter for MVB telegrams
type FilterMask struct {
	// MVB f_codes filter mask. Each bit corresponds to a specific f_code, bit 0=fcode-0, bit 1=fcode-1 etc
	FCodeMask uint16
	// Address to compare
	Address uint16
	// mask for comparison. Only bits set to one are compared against address
	Mask uint16
	// whether to include frames without slave response
	IncludeTimedoutFrames bool
}

// StreamFilter defines the MVB filter to be applied to a stream
// Refer to firmware documentation for max. numeber of Masks supported
type StreamFilter struct {
	Masks []FilterMask
}

// NewClientFromUniversalAddress creates a new mvbSniffer client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of a mdns service (without _io4edge_mvbSniffer._tcp).
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_mvbSniffer._tcp", timeout)

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

// StreamConfigOption is a type to pass options to StartStream()
type StreamConfigOption func(*StreamConfiguration)

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	FilterMask []*fspb.FilterMask
	FBOptions  []functionblock.StreamConfigOption
}

// WithFilterMask may be passed once or multiple times to StartStream.
func WithFilterMask(mask FilterMask) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.FilterMask = append(c.FilterMask, &fspb.FilterMask{
			FCodeMask:             uint32(mask.FCodeMask),
			Address:               uint32(mask.Address),
			Mask:                  uint32(mask.Mask),
			IncludeTimedoutFrames: mask.IncludeTimedoutFrames,
		})
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
//   - WithFilterMask
//   - WithFBStreamOption(functionblock.WithXXXX(...))
//
// Options that are not specified take default values.
func (c *Client) StartStream(opts ...StreamConfigOption) error {
	config := &StreamConfiguration{}

	for _, opt := range opts {
		opt(config)
	}

	err := c.fbClient.StartStream(config.FBOptions, &fspb.StreamControlStart{
		Filter: config.FilterMask,
	})

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

	fsSD := new(fspb.TelegramCollection)
	if err := genericSD.FSData.UnmarshalTo(fsSD); err != nil {
		return nil, fmt.Errorf("can't unmarshall samples: %v", err)
	}

	sd := &StreamData{
		StreamDataMeta: genericSD.StreamDataMeta,
		FSData:         fsSD,
	}
	return sd, nil
}
