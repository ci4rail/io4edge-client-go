package functionblock

import (
	fbv1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
	"github.com/docker/distribution/uuid"
	any "github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/runtime/protoiface"
)

func StreamControlStart(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_StreamControl{
			StreamControl: &fbv1.StreamControl{
				Action: &fbv1.StreamControl_FunctionSpecificStreamControlStart{
					FunctionSpecificStreamControlStart: anyCmd,
				},
			},
		},
	}, nil
}

func StreamControlStop() (*fbv1.Command, error) {
	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_StreamControl{
			StreamControl: &fbv1.StreamControl{
				Action: &fbv1.StreamControl_Stop_{},
			},
		},
	}, nil
}
