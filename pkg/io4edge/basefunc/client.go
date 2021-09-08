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

// Package basefunc provides the API for the io4edge base functions
// i.e. firmware and hardware id management
package basefunc

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge"
)

// Client represents a client for the io4edge base function
type Client struct {
	ch *io4edge.Channel
}

// NewClient creates a new client for the base function
func NewClient(c *io4edge.Channel) (*Client, error) {
	return &Client{ch: c}, nil
}

// Command issues a command cmd to a base function channel, waits for the devices response and returns it in res
func (c *Client) Command(cmd *BaseFuncCommand, res *BaseFuncResponse, timeout time.Duration) error {
	err := c.ch.WriteMessage(cmd)
	if err != nil {
		return err
	}
	err = c.ch.ReadMessage(res, timeout)
	if err != nil {
		return errors.New("Failed to receive device response: " + err.Error())
	}
	if res.Status != BaseFuncStatus_OK {
		return errors.New("Device reported error status: " + res.Status.String())
	}
	return err
}
