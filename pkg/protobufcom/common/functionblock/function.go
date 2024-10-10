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
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// FunctionControlSet executes the function control set command on the device.
//
// fsCmd is the function specific command object
//
// returns the function specific response as a protobuf any object
func (c *Client) FunctionControlSet(fsCmd proto.Message) (*anypb.Any, error) {
	cmd, err := functionControlSetMessage(fsCmd)
	if err != nil {
		return nil, err
	}
	res, err := c.command(cmd)
	if err != nil {
		return nil, err
	}
	fsRes := res.GetFunctionControl().GetFunctionSpecificFunctionControlSet()
	return fsRes, nil
}

// FunctionControlGet executes the function control get command on the device.
//
// fsCmd is the function specific command object
//
// returns the function specific response as a protobuf any object
func (c *Client) FunctionControlGet(fsCmd proto.Message) (*anypb.Any, error) {
	cmd, err := functionControlGetMessage(fsCmd)
	if err != nil {
		return nil, err
	}
	res, err := c.command(cmd)
	if err != nil {
		return nil, err
	}
	fsRes := res.GetFunctionControl().GetFunctionSpecificFunctionControlGet()
	return fsRes, nil
}

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
