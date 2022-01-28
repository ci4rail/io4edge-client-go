package templateModule

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	templateModule "github.com/ci4rail/io4edge_api/templateModule/go/templateModule/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Configuration struct {
}

func (c *Client) SetConfiguration(config Configuration) error {
	if c.connected {
		cmd := templateModule.ConfigurationControlSet{}

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
		if res.Status == functionblockV1.Status_WRONG_CLIENT {
			return fmt.Errorf("wrong client")
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
		cmd := templateModule.ConfigurationControlGet{}
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
		if res.Status == functionblockV1.Status_WRONG_CLIENT {
			return nil, fmt.Errorf("wrong client")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return nil, fmt.Errorf(res.Error.Error)
		}
		get := templateModule.ConfigurationControlResponse{}
		err = anypb.UnmarshalTo(res.GetConfigurationControl().FunctionSpecificConfigurationControlResponse, &get, proto.UnmarshalOptions{})
		if err != nil {
			return nil, err
		}
		ret := &Configuration{}
		return ret, nil
	}
	return nil, fmt.Errorf("not connected")
}

func (c *Client) Describe() (*templateModule.ConfigurationControlDescribeResponse, error) {
	if c.connected {
		cmd := templateModule.ConfigurationControlDescribe{}
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
		if res.Status == functionblockV1.Status_WRONG_CLIENT {
			return nil, fmt.Errorf("wrong client")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return nil, fmt.Errorf(res.Error.Error)
		}
		describe := templateModule.ConfigurationControlResponse{}
		err = anypb.UnmarshalTo(res.GetConfigurationControl().FunctionSpecificConfigurationControlResponse, &describe, proto.UnmarshalOptions{})
		if err != nil {
			return nil, err
		}

		return describe.GetDescribe(), nil
	}
	return nil, fmt.Errorf("not connected")
}
