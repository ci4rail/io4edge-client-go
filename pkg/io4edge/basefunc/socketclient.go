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

package basefunc

import (
	"errors"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge"
	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/transport"
)

// NewClientFromSocketAddress creates a new base function client from a socket with the specified address.
func NewClientFromSocketAddress(address string) (*Client, error) {
	t, err := transport.NewSocketConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ms, err := transport.NewMsgStreamFromConnection(t)
	if err != nil {
		return nil, errors.New("can't create msg stream: " + err.Error())
	}

	ch, err := io4edge.NewChannel(ms)
	if err != nil {
		return nil, errors.New("can't create channel: " + err.Error())
	}
	c, err := NewClient(ch)
	if err != nil {
		return nil, errors.New("can't basefunc create client: " + err.Error())
	}
	return c, nil
}
