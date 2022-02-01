package binaryiotypea

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	binio "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
	functionblockV1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Configuration is the type for set/get configuration
type Configuration struct {
	OutputFritting        int32
	OutputWatchdog        int32
	OutputWatchdogTimeout int32
}

// SetConfiguration sets configuration
func (c *Client) SetConfiguration(config Configuration) error {
	if c.connected {
		cmd := binio.ConfigurationSet{
			OutputFrittingMask:    config.OutputFritting,
			OutputWatchdogMask:    config.OutputWatchdog,
			OutputWatchdogTimeout: config.OutputWatchdogTimeout,
		}

		envelopeCmd, err := functionblock.ConfigurationSet(&cmd)
		if err != nil {
			return err
		}
		res := &functionblockV1.Response{}
		err = c.funcClient.Command(envelopeCmd, res, time.Second*5)
		if err != nil {
			return err
		}
		if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
			return fmt.Errorf("not implemented")
		}
		if res.Status != functionblockV1.Status_OK {
			fmt.Printf("Status %d\n", res.Status)
			return fmt.Errorf(res.Error.Error)
		}
		return nil
	}
	return fmt.Errorf("not connected")
}

// GetConfiguration gets configuration
func (c *Client) GetConfiguration() (*Configuration, error) {
	if c.connected {
		cmd := binio.ConfigurationGet{}
		envelopeCmd, err := functionblock.ConfigurationGet(&cmd)
		if err != nil {
			return nil, err
		}
		res := &functionblockV1.Response{}
		err = c.funcClient.Command(envelopeCmd, res, time.Second*5)
		if err != nil {
			return nil, err
		}
		if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
			return nil, fmt.Errorf("not implemented")
		}
		if res.Status != functionblockV1.Status_OK {
			return nil, fmt.Errorf(res.Error.Error)
		}
		get := binio.ConfigurationGetResponse{}
		err = anypb.UnmarshalTo(res.GetConfiguration().GetFunctionSpecificConfigurationGet(), &get, proto.UnmarshalOptions{})
		if err != nil {
			return nil, err
		}
		ret := &Configuration{
			OutputFritting:        int32(get.OutputFrittingMask) & 0xFF,
			OutputWatchdog:        int32(get.OutputWatchdogMask) & 0xFF,
			OutputWatchdogTimeout: get.OutputWatchdogTimeout,
		}
		return ret, nil
	}
	return nil, fmt.Errorf("not connected")
}

// Describe the function
func (c *Client) Describe() (*binio.ConfigurationDescribeResponse, error) {
	if c.connected {
		cmd := binio.ConfigurationDescribe{}
		envelopeCmd, err := functionblock.ConfigurationDescribe(&cmd)
		if err != nil {
			return nil, err
		}
		res := &functionblockV1.Response{}
		err = c.funcClient.Command(envelopeCmd, res, time.Second*5)
		if err != nil {
			return nil, err
		}

		if res.Status != functionblockV1.Status_OK {
			return nil, fmt.Errorf(res.Error.Error)
		}
		describe := binio.ConfigurationDescribeResponse{}
		err = anypb.UnmarshalTo(res.GetConfiguration().GetFunctionSpecificConfigurationDescribe(), &describe, proto.UnmarshalOptions{})
		if err != nil {
			return nil, err
		}

		return &describe, nil
	}
	return nil, fmt.Errorf("not connected")
}
