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

// Restart performs a device restart
func (c *Client) Restart(timeout time.Duration) (restartingNow bool, err error) {
	cmd := &api.CoreCommand{
		Id: api.CommandId_RESTART,
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return false, err
	}
	return res.RestartingNow, nil
}
