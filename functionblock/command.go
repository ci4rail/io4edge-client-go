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

package functionblock

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
)

// Command issues a command cmd to a channel, waits for the devices to respond and returns the response
// timeout specifies how long to wait for response
func (c *Client) command(cmd *fbv1.Command) (*fbv1.Response, error) {

	// only one command may be pending per connection
	c.cmdMutex.Lock()
	defer c.cmdMutex.Unlock()

	cmd.Context = &fbv1.Context{Value: fmt.Sprintf("%d", c.cmdSeqNo)}

	log.Debug("sending command with context: ", cmd.Context)

	err := c.funcClient.Ch.WriteMessage(cmd)
	if err != nil {
		return nil, err
	}

	// tell response loop which seqNo is waiting
	c.waitingCmdSeqChan <- c.cmdSeqNo
	c.cmdSeqNo++

	var res *fbv1.Response

	// wait for response loop to wake us
	select {
	case r := <-c.responseChan:
		res = r
	case <-time.After(time.Duration(c.commandTimeout) * time.Second):
		err = errors.New("timeout waiting for command")
	}

	if err != nil {
		return nil, err
	}
	if res.Status != fbv1.Status_OK {
		err = responseErrorNew(cmd.String(), res)
		return res, err
	}

	return res, nil
}

func wakeupCommand(c *Client, res *fbv1.Response) {
	log.Debug("trying to wakeup command for context", res.Context.Value)
	// check if we're waking up the right command
	select {
	case context := <-c.waitingCmdSeqChan:

		log.Debug("got waiting context ", context)
		if fmt.Sprintf("%d", context) == res.Context.Value {
			select {
			case c.responseChan <- res:
				log.Debug("wakeup command for context", context)
			case <-time.After(1 * time.Second):
				log.Error("cannot wakeup command")
			}
		} else {
			log.Error("got response for wrong context")
		}
	case <-time.After(1 * time.Second):
		log.Error("timeout waiting for command context")
		return
	}
}

// readResponses read responses in background
func (c *Client) readResponses() {
	log.Debug("about to start go ReadResponses()")
	go func(c *Client) {
		//defer c.recover()
		log.Debug("ReadResponses running")
		for {
			log.Debug("ReadResponses loop")
			res := &fbv1.Response{}
			err := c.funcClient.ReadMessage(res, 0)
			if err != nil {
				break
			}
			switch res.Type.(type) {
			case *fbv1.Response_Stream:
				c.streamChan <- res.GetStream()

			default:
				go wakeupCommand(c, res)
			}
		}
		log.Error("ReadResponses exit")
	}(c)
}
