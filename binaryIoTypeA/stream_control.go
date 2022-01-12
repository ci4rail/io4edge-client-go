package binaryIoTypeA

import (
	"fmt"
	"time"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func (c *Client) StreamStatus() bool {
	return c.streamStatus
}

type StreamConfiguration struct {
	ChannelFilterMask uint32
	KeepaliveInterval uint32
}

var (
	defaultConfiguration = &StreamConfiguration{
		ChannelFilterMask: DefaultChannelFilterMask,
		KeepaliveInterval: DefaultKeepaliveInterval,
	}
)

func (c *Client) StreamDefaultConfiguration() *StreamConfiguration {
	return defaultConfiguration
}

func (c *Client) StartStream(config *StreamConfiguration, callback func(*binio.Sample)) error {
	if c.connected {
		if config == nil {
			config = defaultConfiguration
		}
		cmd := binio.StreamControlStart{
			ChannelFilterMask: config.ChannelFilterMask,
			KeepaliveInterval: config.KeepaliveInterval,
		}

		envelopeCmd, err := functionblock.StreamControlStart(&cmd)
		if err != nil {
			return err
		}
		res, err := c.Command(envelopeCmd, time.Second*5)
		if err != nil {
			return err
		}
		if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
			return fmt.Errorf("not implemented")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return fmt.Errorf(res.Error.Error)
		}
		go func(c *Client, callback func(*binio.Sample)) {
			c.streamRunning = true
			for {
				select {
				case <-c.streamClientStopChannel:
					c.streamRunning = false
					return
				default:
					select {
					case res := <-c.streamData:
						log.Debugf("Received message: %+v\n", res)
						streamData := &binio.StreamData{}
						err = anypb.UnmarshalTo(res.FunctionSpecificStreamData, streamData, proto.UnmarshalOptions{})
						if err != nil {
							log.Error(err)
							continue
						}
						for _, sample := range streamData.Samples {
							if callback != nil {
								callback(sample)
							}
						}
					case <-time.After(time.Millisecond * 100):
						time.Sleep(time.Millisecond * 100)
						continue
					}
				}
			}
		}(c, callback)
		return nil
	}
	return fmt.Errorf("not connected")
}

func (c *Client) StopStream() error {
	if c.connected {
		cmd, err := functionblock.StreamControlStop()
		if err != nil {
			return err
		}
		res, err := c.Command(cmd, time.Second*5)
		if err != nil {
			return err
		}
		if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
			return fmt.Errorf("not implemented")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return fmt.Errorf(res.Error.Error)
		}

		c.streamClientStopChannel <- true
		log.Debug("Stopped stream")
		return nil
	}
	return fmt.Errorf("not connected")
}
