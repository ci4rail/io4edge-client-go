package functionblock

import (
	"errors"
	"time"

	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

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

// StartStream starts the stream with configuration config, passing the function specific configuration from fscmd
func (c *Client) StartStream(config *StreamConfiguration, fsCmd proto.Message) error {
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
