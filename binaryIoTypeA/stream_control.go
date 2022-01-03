package binaryIoTypeA

import (
	"fmt"
	"time"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
)

func (c *Client) StartStream(enableMask int, callback func(*binio.Sample)) (int, error) {
	cmd := binio.StreamControlStart{
		EnableMask: uint32(enableMask),
	}

	envelopeCmd, err := functionblock.StreamControlStart(&cmd)
	if err != nil {
		return -1, err
	}
	res := &functionblockV1.Response{}
	err = c.funcClient.Command(envelopeCmd, res, time.Second*5)
	if err != nil {
		return -1, err
	}
	if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
		return -1, fmt.Errorf("not implemented")
	}
	if res.Status == functionblockV1.Status_ERROR {
		return -1, fmt.Errorf(res.Error.Error)
	}
	id := int(res.GetStreamControl().Id)
	go func(quit chan bool, callback func(*binio.Sample)) {
		fmt.Println("Started stream for id:", id)
		for {
			select {
			case <-quit:
				return
			default:
				res := &binio.Sample{}
				err := c.funcClient.ReadMessage(res, 0)
				if err != nil {
					fmt.Println(err)
					continue
				}
				callback(res)
			}
		}
	}(c.streamClientChannels[id], callback)
	return id, nil
}

func (c *Client) StopStream(id int) error {
	cmd, err := functionblock.StreamControlStop(id)
	if err != nil {
		return err
	}
	res := &functionblockV1.Response{}
	err = c.funcClient.Command(cmd, res, time.Second*5)
	if err != nil {
		return err
	}
	if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
		return fmt.Errorf("not implemented")
	}
	if res.Status == functionblockV1.Status_ERROR {
		return fmt.Errorf(res.Error.Error)
	}
	return nil
}
