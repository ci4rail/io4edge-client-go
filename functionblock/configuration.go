package functionblock

import (
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	"github.com/docker/distribution/uuid"
	any "github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/runtime/protoiface"
)

// ConfigurationSet writes a new configuration
func ConfigurationSet(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_Configuration{

			Configuration: &fbv1.Configuration{
				Action: &fbv1.Configuration_FunctionSpecificConfigurationSet{
					FunctionSpecificConfigurationSet: anyCmd,
				},
			},
		},
	}, nil
}

// ConfigurationGet reads the configuration
func ConfigurationGet(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_Configuration{
			Configuration: &fbv1.Configuration{
				Action: &fbv1.Configuration_FunctionSpecificConfigurationGet{
					FunctionSpecificConfigurationGet: anyCmd,
				},
			},
		},
	}, nil
}

// ConfigurationDescribe describes the unit
func ConfigurationDescribe(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}
	ctx := uuid.Generate().String()
	return &fbv1.Command{
		Context: &fbv1.Context{Value: ctx},
		Type: &fbv1.Command_Configuration{
			Configuration: &fbv1.Configuration{
				Action: &fbv1.Configuration_FunctionSpecificConfigurationDescribe{
					FunctionSpecificConfigurationDescribe: anyCmd,
				},
			},
		},
	}, nil
}
