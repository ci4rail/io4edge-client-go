package binaryiotypea

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	binio "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
	functionblockV1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// StreamStatus returns true if the stream is up
func (c *Client) StreamStatus() bool {
	return c.streamStatus
}

// StreamConfiguration defines the configuration of a stream
type StreamConfiguration struct {
	ChannelFilterMask uint32
	KeepaliveInterval uint32
	BufferSize        uint32
}

var (
	defaultConfiguration = &StreamConfiguration{
		ChannelFilterMask: DefaultChannelFilterMask,
		KeepaliveInterval: DefaultKeepaliveInterval,
		BufferSize:        DefaultBufferSize,
	}
)

// StreamDefaultConfiguration returns the default stream configuration
func (c *Client) StreamDefaultConfiguration() *StreamConfiguration {
	return defaultConfiguration
}

// StartStream starts the stream
func (c *Client) StartStream(config *StreamConfiguration, callback func(*binio.Sample, uint32)) error {
	if c.connected {
		if config == nil {
			config = defaultConfiguration
		}
		cmd := binio.StreamControlStart{
			ChannelFilterMask: config.ChannelFilterMask,
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
		if res.Status != functionblockV1.Status_OK {
			return fmt.Errorf(res.Error.Error)
		}
		go func(c *Client, callback func(*binio.Sample, uint32)) {
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

// StopStream stops the stream
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
		if res.Status != functionblockV1.Status_OK {
			return fmt.Errorf(res.Error.Error)
		}

		c.streamClientStopChannel <- true
		log.Debug("Stopped stream")
		return nil
	}
	return fmt.Errorf("not connected")
}
