package functionblock

import (
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func functionControlSetMessage(fsCmd proto.Message) (*fbv1.Command, error) {
	anyCmd, err := anypb.New(fsCmd)
	if err != nil {
		return nil, err
	}
	return &fbv1.Command{
		Type: &fbv1.Command_FunctionControl{
			FunctionControl: &fbv1.FunctionControl{
				Action: &fbv1.FunctionControl_FunctionSpecificFunctionControlSet{
					FunctionSpecificFunctionControlSet: anyCmd,
				},
			},
		},
	}, nil
}

func functionControlGetMessage(fsCmd proto.Message) (*fbv1.Command, error) {
	anyCmd, err := anypb.New(fsCmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Type: &fbv1.Command_FunctionControl{
			FunctionControl: &fbv1.FunctionControl{
				Action: &fbv1.FunctionControl_FunctionSpecificFunctionControlGet{
					FunctionSpecificFunctionControlGet: anyCmd,
				},
			},
		},
	}, nil
}
