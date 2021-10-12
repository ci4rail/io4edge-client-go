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

// Package client provides the API for io4edge I/O devices
package client

import (
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
)

// Client represents a client for an io4edge function
type Client struct {
	ch *Channel
}

// NewClient creates a new client for an io4edge function
func NewClient(c *Channel) *Client {
	return &Client{ch: c}
}

// Command issues a command cmd to a channel, waits for the devices response and returns it in res
func (c *Client) Command(cmd proto.Message, res proto.Message, timeout time.Duration) error {
	err := c.ch.WriteMessage(cmd)
	if err != nil {
		return err
	}
	err = c.ch.ReadMessage(res, timeout)
	if err != nil {
		return errors.New("Failed to receive device response: " + err.Error())
	}
	return err
}
