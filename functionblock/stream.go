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

package functionblock

import (
	"errors"
	"time"

	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// StreamConfigOption is a type to pass options to StartStream()
type StreamConfigOption func(*StreamConfiguration)

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	BucketSamples     uint32
	KeepaliveInterval uint32
	BufferedSamples   uint32
}

// StreamDataMeta contains meta information about a Stream Data message
type StreamDataMeta struct {
	DeliveryTimestamp uint64
	Sequence          uint32
}

// StreamData contains the meta data of the stream and the function specific message
type StreamData struct {
	StreamDataMeta
	FSData *anypb.Any // function specific data
}

// WithBucketSamples may be passed to StartStream.
//
// numSamples define the max number of samples per message in StreamData
func WithBucketSamples(numSamples uint32) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.BucketSamples = numSamples
	}
}

// WithBufferedSamples may be passed to StartStream.
//
// numSamples define the number of samples buffered in the device for that stream
func WithBufferedSamples(numSamples uint32) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.BufferedSamples = numSamples
	}
}

// WithKeepaliveInterval may be passed to StartStream.
//
// timeMS defines the time in ms after which the devices buffer is flushed and sent to
// the client, even if the number of buffered samples is less than BucketSamples
func WithKeepaliveInterval(timeMS uint32) StreamConfigOption {
	return func(c *StreamConfiguration) {
		c.KeepaliveInterval = timeMS
	}
}

// StartStream starts the stream with configuration config, passing the function specific configuration from fscmd
func (c *Client) StartStream(opts []StreamConfigOption, fsCmd proto.Message) error {
	config := &StreamConfiguration{
		BucketSamples:     25,
		KeepaliveInterval: 1000,
		BufferedSamples:   50,
	}
	for _, opt := range opts {
		opt(config)
	}
	cmd, err := streamControlStartMessage(config, fsCmd)
	if err != nil {
		return err
	}
	_, err = c.command(cmd)
	if err != nil {
		return err
	}
	return nil
}

// StopStream stops the stream
func (c *Client) StopStream() error {
	cmd, err := streamControlStopMessage()
	if err != nil {
		return err
	}
	_, err = c.command(cmd)
	if err != nil {
		return err
	}
	log.Debug("Stopped stream")
	return nil
}

// ReadStream reads the next stream data object from the buffer
func (c *Client) ReadStream(timeout time.Duration) (*StreamData, error) {
	select {
	case d := <-c.streamChan:
		log.Debug("got stream data seq", d.Sequence)

		sd := &StreamData{
			StreamDataMeta: StreamDataMeta{
				DeliveryTimestamp: d.DeliveryTimestampUs,
				Sequence:          d.Sequence,
			},
			FSData: d.FunctionSpecificStreamData,
		}
		return sd, nil

	case <-time.After(timeout):
		// TODO: Specific error code!
		log.Warn("ReadStreamData timeout")
		return nil, errors.New("timeout waiting for stream data")
	}
}

// streamControlStartMessage returns the message to start the stream
func streamControlStartMessage(config *StreamConfiguration, fsCmd proto.Message) (*fbv1.Command, error) {
	anyCmd, err := anypb.New(fsCmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Type: &fbv1.Command_StreamControl{
			StreamControl: &fbv1.StreamControl{
				Action: &fbv1.StreamControl_Start{
					Start: &fbv1.StreamControlStart{
						BucketSamples:                      config.BucketSamples,
						KeepaliveInterval:                  config.KeepaliveInterval,
						BufferedSamples:                    config.BufferedSamples,
						FunctionSpecificStreamControlStart: anyCmd,
					},
				},
			},
		},
	}, nil
}

// streamControlStopMessage returns the message to stop the stream
func streamControlStopMessage() (*fbv1.Command, error) {
	return &fbv1.Command{
		Type: &fbv1.Command_StreamControl{
			StreamControl: &fbv1.StreamControl{
				Action: &fbv1.StreamControl_Stop{
					Stop: &fbv1.StreamControlStop{},
				},
			},
		},
	}, nil
}
