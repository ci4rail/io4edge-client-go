package functionblock

import (
	fbv1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
	"github.com/docker/distribution/uuid"
	any "github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/runtime/protoiface"
)

func ConfigurationControlSet(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_ConfigurationControl{
			ConfigurationControl: &fbv1.ConfigurationControl{
				Action: &fbv1.ConfigurationControl_FunctionSpecificConfigurationControlSet{
					FunctionSpecificConfigurationControlSet: anyCmd,
				},
			},
		},
	}, nil
}

func ConfigurationControlGet(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Context: &fbv1.Context{Value: uuid.Generate().String()},
		Type: &fbv1.Command_ConfigurationControl{
			ConfigurationControl: &fbv1.ConfigurationControl{
				Action: &fbv1.ConfigurationControl_FunctionSpecificConfigurationControlGet{
					FunctionSpecificConfigurationControlGet: anyCmd,
				},
			},
		},
	}, nil
}

func ConfigurationControlDescribe(cmd protoiface.MessageV1) (*fbv1.Command, error) {
	anyCmd, err := any.MarshalAny(cmd)
	if err != nil {
		return nil, err
	}
	ctx := uuid.Generate().String()
	return &fbv1.Command{
		Context: &fbv1.Context{Value: ctx},
		Type: &fbv1.Command_ConfigurationControl{
			ConfigurationControl: &fbv1.ConfigurationControl{
				Action: &fbv1.ConfigurationControl_FunctionSpecificConfigurationControlDescribe{
					FunctionSpecificConfigurationControlDescribe: anyCmd,
				},
			},
		},
	}, nil
}
