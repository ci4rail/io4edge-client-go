package functionblock

import (
	functionblockv1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
	"github.com/docker/distribution/uuid"
	any "github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/runtime/protoiface"
)

func CreateCommand(cmd protoiface.MessageV1) (*functionblockv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}

	return &functionblockv1.Command{
		Context: &functionblockv1.Context{
			Value: uuid.Generate().String(),
		},
		Type: &functionblockv1.Command_Configuration{
			Configuration: &functionblockv1.ConfigurationControl{
				FunctionSpecificHardwareConfiguration: anyCmd,
			},
		},
	}, nil
}
