package functionblock

import (
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	"github.com/docker/distribution/uuid"
	any "github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/runtime/protoiface"
)

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	BucketSamples     uint32
	KeepaliveInterval uint32
	BufferedSamples   uint32
}

// StreamControlStart returns the marshalled message to start the stream
func StreamControlStart(config *StreamConfiguration, fscmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(fscmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
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

// StreamControlStop returns the marshalled message to stop the stream
func StreamControlStop() (*fbv1.Command, error) {
	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_StreamControl{
			StreamControl: &fbv1.StreamControl{
				Action: &fbv1.StreamControl_Stop{
					Stop: &fbv1.StreamControlStop{},
				},
			},
		},
	}, nil
}
