package binaryIoTypeA

import (
	"fmt"
	"time"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Configuration struct {
	OutputFritting        int32
	OutputWatchdog        int32
	OutputWatchdogTimeout int32
}

func (c *Client) SetConfiguration(config Configuration) error {
	if c.connected {
		cmd := binio.ConfigurationControlSet{
			OutputFrittingMask:    config.OutputFritting,
			OutputWatchdogMask:    config.OutputWatchdog,
			OutputWatchdogTimeout: config.OutputWatchdogTimeout,
		}

		envelopeCmd, err := functionblock.ConfigurationControlSet(&cmd)
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
		if res.Status == functionblockV1.Status_ERROR {
			return fmt.Errorf(res.Error.Error)
		}
		return nil
	}
	return fmt.Errorf("not connected")
}

func (c *Client) GetConfiguration() (*Configuration, error) {
	if c.connected {
		cmd := binio.ConfigurationControlGet{}
		envelopeCmd, err := functionblock.ConfigurationControlGet(&cmd)
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
		if res.Status == functionblockV1.Status_ERROR {
			return nil, fmt.Errorf(res.Error.Error)
		}
		get := binio.ConfigurationControlResponse{}
		err = anypb.UnmarshalTo(res.GetConfigurationControl().FunctionSpecificConfigurationControlResponse, &get, proto.UnmarshalOptions{})
		if err != nil {
			return nil, err
		}
		ret := &Configuration{
			OutputFritting:        int32(get.GetGet().OutputFrittingMask) & 0xFF,
			OutputWatchdog:        int32(get.GetGet().OutputWatchdogMask) & 0xFF,
			OutputWatchdogTimeout: get.GetGet().OutputWatchdogTimeout,
		}
		return ret, nil
	}
	return nil, fmt.Errorf("not connected")
}

func (c *Client) Describe() (*binio.ConfigurationControlDescribeResponse, error) {
	if c.connected {
		cmd := binio.ConfigurationControlDescribe{}
		envelopeCmd, err := functionblock.ConfigurationControlDescribe(&cmd)
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
		if res.Status == functionblockV1.Status_ERROR {
			return nil, fmt.Errorf(res.Error.Error)
		}
		describe := binio.ConfigurationControlResponse{}
		err = anypb.UnmarshalTo(res.GetConfigurationControl().FunctionSpecificConfigurationControlResponse, &describe, proto.UnmarshalOptions{})
		if err != nil {
			return nil, err
		}

		return describe.GetDescribe(), nil
	}
	return nil, fmt.Errorf("not connected")
}
