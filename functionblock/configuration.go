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

// UploadConfiguration executes the configuration set command on the device.
//
// fsCmd is the function specific configuration set object
//
// returns the function specific response as a protobuf any object
func (c *Client) UploadConfiguration(fsCmd proto.Message) (*anypb.Any, error) {
	cmd, err := configurationSetMessage(fsCmd)
	if err != nil {
		return nil, err
	}
	res, err := c.command(cmd)
	if err != nil {
		return nil, err
	}
	fsRes := res.GetConfiguration().GetFunctionSpecificConfigurationSet()
	return fsRes, nil
}

// DownloadConfiguration executes the configuration get command on the device.
//
// fsCmd is the function specific configuration get object
//
// returns the function specific response as a protobuf any object
func (c *Client) DownloadConfiguration(fsCmd proto.Message) (*anypb.Any, error) {
	cmd, err := configurationGetMessage(fsCmd)
	if err != nil {
		return nil, err
	}
	res, err := c.command(cmd)
	if err != nil {
		return nil, err
	}
	fsRes := res.GetConfiguration().GetFunctionSpecificConfigurationGet()
	return fsRes, nil
}

// Describe executes the configuration describe command on the device.
//
// fsCmd is the function specific configuration describe object
//
// returns the function specific response as a protobuf any object
func (c *Client) Describe(fsCmd proto.Message) (*anypb.Any, error) {
	cmd, err := configurationDescribeMessage(fsCmd)
	if err != nil {
		return nil, err
	}

	res, err := c.command(cmd)
	if err != nil {
		return nil, err
	}
	fsRes := res.GetConfiguration().GetFunctionSpecificConfigurationDescribe()
	return fsRes, nil
}

// configurationSetMessage returns a functionblock level ConfigurationSet Message with function specific fsCmd embedded
func configurationSetMessage(fsCmd proto.Message) (*fbv1.Command, error) {
	anyCmd, err := anypb.New(fsCmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Type: &fbv1.Command_Configuration{

			Configuration: &fbv1.Configuration{
				Action: &fbv1.Configuration_FunctionSpecificConfigurationSet{
					FunctionSpecificConfigurationSet: anyCmd,
				},
			},
		},
	}, nil
}

// configurationGetMessage returns a functionblock level ConfigurationGet Message with function specific fsCmd embedded
func configurationGetMessage(fsCmd proto.Message) (*fbv1.Command, error) {
	anyCmd, err := anypb.New(fsCmd)
	if err != nil {
		return nil, err
	}

	return &fbv1.Command{
		Type: &fbv1.Command_Configuration{
			Configuration: &fbv1.Configuration{
				Action: &fbv1.Configuration_FunctionSpecificConfigurationGet{
					FunctionSpecificConfigurationGet: anyCmd,
				},
			},
		},
	}, nil
}

// configurationDescribe returns a functionblock level ConfigurationDescribe Message with function specific fsCmd embedded
func configurationDescribeMessage(fsCmd proto.Message) (*fbv1.Command, error) {
	anyCmd, err := anypb.New(fsCmd)
	if err != nil {
		return nil, err
	}
	return &fbv1.Command{
		Type: &fbv1.Command_Configuration{
			Configuration: &fbv1.Configuration{
				Action: &fbv1.Configuration_FunctionSpecificConfigurationDescribe{
					FunctionSpecificConfigurationDescribe: anyCmd,
				},
			},
		},
	}, nil
}
