package templateModule

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	templateModule "github.com/ci4rail/io4edge_api/templateModule/go/templateModule/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func (c *Client) StreamStatus() bool {
	return c.streamStatus
}

type StreamConfiguration struct {
	KeepaliveInterval uint32
	BufferSize        uint32
}

var (
	defaultConfiguration = &StreamConfiguration{
		KeepaliveInterval: DefaultKeepaliveInterval,
		BufferSize:        DefaultBufferSize,
	}
)

func (c *Client) StreamDefaultConfiguration() *StreamConfiguration {
	return defaultConfiguration
}

func (c *Client) StartStream(config *StreamConfiguration, callback func(*templateModule.Sample, uint32)) error {
	if c.connected {
		if config == nil {
			config = defaultConfiguration
		}
		cmd := templateModule.StreamControlStart{
			KeepaliveInterval: config.KeepaliveInterval,
			BufferSize:        config.BufferSize,
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
		go func(c *Client, callback func(*templateModule.Sample, uint32)) {
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
						streamData := &templateModule.StreamData{}
						err = anypb.UnmarshalTo(res.FunctionSpecificStreamData, streamData, proto.UnmarshalOptions{})
						if err != nil {
							log.Error(err)
							continue
						}
						for _, sample := range streamData.Samples {
							if callback != nil {
								callback(sample, res.Sequence)
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
