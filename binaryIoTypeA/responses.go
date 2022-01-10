package binaryIoTypeA

import (
	"fmt"
	"time"

	functionblockV1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
)

func (c *Client) ReadResponse() {
	fmt.Println("about to start go ReadResponse()")
	go func(c *Client) {
		fmt.Println("ReadResponse running")
		for {
			fmt.Println("ReadResponse loop")
			res := &functionblockV1.Response{}
			err := c.funcClient.ReadMessage(res, time.Millisecond*1000)
			fmt.Println(res)
			if err != nil {
				continue
			}
			if res.Context != nil {
				fmt.Println("received response for context:", res.Context.Value)
				c.responses.Store(res.Context.Value, res)
			} else {
				c.streamData <- res.GetStream()
			}
		}
	}(c)
}

func (c *Client) GetResponse(context string) *functionblockV1.Response {
	res, ok := c.responses.LoadAndDelete(context)
	if !ok {
		return nil
	}
	return res.(*functionblockV1.Response)
}

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
