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

	api "github.com/ci4rail/io4edge_api/io4edge/go/core_api/v1alpha2"
)

// SetPersistentParameter writes a persistent parameter into the device
func (c *Client) SetPersistentParameter(name string, value string, timeout time.Duration) error {
	cmd := &api.CoreCommand{
		Id: api.CommandId_SET_PERSISTENT_PARAMETER,
		Data: &api.CoreCommand_SetPersistentParameter{
			SetPersistentParameter: &api.SetPersistentParameterCommand{
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

// GetPersistentParameter reads a persistent parameter from the device
func (c *Client) GetPersistentParameter(name string, timeout time.Duration) (value string, err error) {
	cmd := &api.CoreCommand{
		Id: api.CommandId_GET_PERSISTENT_PARAMETER,
		Data: &api.CoreCommand_GetPersistentParameter{
			GetPersistentParameter: &api.GetPersistentParameterCommand{
				Name: name,
			},
		},
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return "", err
	}
	return res.GetPersistentParameter().Value, nil
}
