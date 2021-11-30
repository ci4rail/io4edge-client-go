/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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

package core

import (
	"time"

	api "github.com/ci4rail/io4edge-client-go/core/v1alpha2"
)

// SetPersistantParameter writes a persistant parameter into the device
func (c *Client) SetPersistantParameter(name string, value string, timeout time.Duration) error {
	cmd := &api.CoreCommand{
		Id: api.CommandId_SET_PERSISTANT_PARAMETER,
		Data: &api.CoreCommand_SetPersistantParameter{
			SetPersistantParameter: &api.SetPersistantParameterCommand{
				Name:  name,
				Value: value,
			},
		},
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return err
	}
	return nil
}

// GetPersistantParameter reads a persistant parameter from the device
func (c *Client) GetPersistantParameter(name string, timeout time.Duration) (value string, err error) {
	cmd := &api.CoreCommand{
		Id: api.CommandId_GET_PERSISTANT_PARAMETER,
		Data: &api.CoreCommand_GetPersistantParameter{
			GetPersistantParameter: &api.GetPersistantParameterCommand{
				Name: name,
			},
		},
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return "", err
	}
	return res.GetPersistantParameter().Value, nil
}
