package binaryiotypea

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	functionblockV1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
)

// ReadResponses read responses in background
func (c *Client) ReadResponses() {
	log.Debug("about to start go ReadResponses()")
	go func(c *Client) {
		defer c.recover()
		log.Debug("ReadResponses running")
		for {
			select {
			case <-c.readResponsesStopChan:
				log.Debug("Exiting ReadResponses")
				return
			default:
				if c.responsePending > 0 || c.streamRunning {
					log.Debug("ReadResponses loop")
					res := &functionblockV1.Response{}
					err := c.funcClient.ReadMessage(res, time.Second*time.Duration(c.streamKeepaliveInterval))
					log.Debug("err:", err)
					if err != nil {
						panic(err)
					}
					if res.Status == functionblockV1.Status_OK {
						if res.Context != nil {
							log.Debug("received response for context:", res.Context.Value)
							c.responses.Store(res.Context.Value, res)
						} else {
							c.streamStatus = true
							c.streamData <- res.GetStream()
						}
					} else if res.Status == functionblockV1.Status_WRONG_CLIENT {
						c.responses.Store(res.Context.Value, res)
					} else {
						fmt.Println("received error response:", res.Error.Error)
					}

				}
			}
		}
	}(c)
}

// GetResponse ???
func (c *Client) GetResponse(context string) *functionblockV1.Response {
	res, ok := c.responses.LoadAndDelete(context)
	if !ok {
		return nil
	}
	return res.(*functionblockV1.Response)
}

// WaitForResponse waits for a response
func (c *Client) WaitForResponse(context string, timeout time.Duration) (*functionblockV1.Response, error) {
	timeoutChan := make(chan bool, 1)
	if timeout > 0 {
		go func() {
			time.Sleep(timeout)
			timeoutChan <- true
		}()
	}
	for {
		select {
		default:
			res := c.GetResponse(context)
			if res != nil {
				return res, nil
			}
		case <-timeoutChan:
			return nil, fmt.Errorf("timeout")
		}
	}
}