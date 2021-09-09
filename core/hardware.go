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

	api "github.com/ci4rail/io4edge-client-go/core/v1alpha1"
)

// IdentifyHardware gets the hardware inventory data from the device
func (c *Client) IdentifyHardware(timeout time.Duration) (*api.IdentifyHardwareResponse, error) {
	cmd := &api.CoreCommand{
		Id: api.CommandId_IDENTIFY_HARDWARE,
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return nil, err
	}
	return res.GetIdentifyHardware(), nil
}

// ProgramHardwareIdentification programs hardware inventory data into the device.
// Intended to be used during hardware manufacturing process only
func (c *Client) ProgramHardwareIdentification(id *api.ProgramHardwareIdentificationCommand, timeout time.Duration) error {
	cmd := &api.CoreCommand{
		Id: api.CommandId_PROGRAM_HARDWARE_IDENTIFICATION,
		Data: &api.CoreCommand_ProgramHardwareIdentification{
			ProgramHardwareIdentification: id,
		},
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return err
	}
	return nil
}
