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

package client

import (
	"time"

	"github.com/ci4rail/io4edge-client-go/transport"
	"google.golang.org/protobuf/proto"
)

// ChannelIf is a interface for the Channel
type ChannelIf interface {
	Close()
	WriteMessage(m proto.Message) error
	ReadMessage(m proto.Message, timeout time.Duration) error
}

// Channel represents a connection between the host and the device used to exchange protobuf messages
type Channel struct {
	ms transport.MsgStream
}

// NewChannel creates a new channel using the transport mechanism in t
func NewChannel(ms transport.MsgStream) *Channel {
	return &Channel{ms: ms}
}

// Close terminates the message stream
func (c *Channel) Close() {
	c.ms.Close()
}

// WriteMessage encodes m using protobuf and sends the encoded value through the message stream
func (c *Channel) WriteMessage(m proto.Message) error {
	payload, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return c.ms.WriteMsg(payload)
}

// ReadMessage waits until Timeout for a new message in transport stream and decodes it via protobuf
// timeout of 0 waits forever
func (c *Channel) ReadMessage(m proto.Message, timeout time.Duration) error {
	msg, err := c.ms.ReadMsg(timeout)

	if err != nil {
		return err
	}

	return proto.Unmarshal(msg, m)
}
