package functionblock

import (
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

// StreamCallback is the type for the callback function to receive stream data
type StreamCallback func(*fbv1.StreamData)

// StreamStart starts the stream with configuration config, passing the function specific configuration from fscmd
func (c *Client) StreamStart(config *StreamConfiguration, fscmd anypb.Any) error {
	cmd, err := StreamControlStartMessage(config, &fscmd)
	if err != nil {
		return err
	}
	_, err = c.command(cmd)
	if err != nil {
		return err
	}
	return nil
}

// StreamStop stops the stream
func (c *Client) StreamStop() error {
	cmd, err := StreamControlStopMessage()
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

// StreamControlStartMessage returns the marshalled message to start the stream
func StreamControlStartMessage(config *StreamConfiguration, fsCmd proto.Message) (*fbv1.Command, error) {
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

// StreamControlStopMessage returns the marshalled message to stop the stream
func StreamControlStopMessage() (*fbv1.Command, error) {
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
