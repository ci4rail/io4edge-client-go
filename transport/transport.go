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

package transport

import (
	"errors"
	"time"
)

// Transport is the interface used by message stream to communicate with the underlying transport layer
// e.g. tcp sockets
type Transport interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	SetReadDeadline(t time.Time) error
	Close() error
}

// TransportMsg is the interface used by message stream to communicate with the underlying transport layer
// e.g. udp sockets
type TransportMsg interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	SetReadDeadline(t time.Time) error
	Close() error
}

// Timeout error
var ErrTimeout = errors.New("socket timeout")
var ErrClosed = errors.New("socket closed")
