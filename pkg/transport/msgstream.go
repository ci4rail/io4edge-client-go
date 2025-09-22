/*
Copyright © 2024 Ci4Rail GmbH <engineering@ci4rail.com>

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

package transport

import "time"

// MsgStream is the interface used by a Channel to exchange message frames with the transport layer
// e.g. socket, websocket...
type MsgStream interface {
	ReadMsg(timeout time.Duration) ([]byte, error)
	WriteMsg(payload []byte) (err error)
	Close() error
}
