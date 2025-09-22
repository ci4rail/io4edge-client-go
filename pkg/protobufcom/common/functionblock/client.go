/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

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

// Package functionblock provides the API for the io4edge function blocks
package functionblock

import (
	"sync"
	"time"

	pbchannelclient "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/channel"
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
)

// Client represents a client for a generic functionblock
type Client struct {
	funcClient     *pbchannelclient.Client
	commandTimeout int
	streamChan     chan *fbv1.StreamData
	cmdSeqNo       uint32
	cmdMutex       sync.Mutex
	// command/response handshake
	waitingCmdSeqChan chan uint32 // waiting sequence number
	responseChan      chan *fbv1.Response
}

func newClient(c *pbchannelclient.Client) *Client {
	client := &Client{
		funcClient:        c,
		commandTimeout:    5,
		streamChan:        make(chan *fbv1.StreamData, 100),
		waitingCmdSeqChan: make(chan uint32),
		responseChan:      make(chan *fbv1.Response),
	}
	client.readResponses()
	return client
}

// NewClientFromUniversalAddress creates a new functionblock client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is the instance name of an mdns service.
//
// If service is non-empty and addrOrService is a mdns instance name, it is appended to the addrOrService.
// .e.g. if addrOrService is "iou01-sn01-binio" and service is "_io4edge_binaryIoTypeA._tcp", the mdns instance
// name "iou01-sn01-binio._io4edge_binaryIoTypeA._tcp" is used.
//
// The timeout specifies the maximal time waiting for a service to show up. If 0, use default timeout. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, service string, timeout time.Duration) (*Client, error) {
	io4eClient, err := pbchannelclient.NewClientFromUniversalAddress(addrOrService, service, timeout)

	if err != nil {
		return nil, err
	}
	c := newClient(io4eClient)

	return c, nil
}

// Close terminates the underlying connection to the functionblock
func (c *Client) Close() {
	c.funcClient.Close()
}
