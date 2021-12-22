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
	Fritting map[int]bool
}

func (c *Client) SetConfiguration(config Configuration) error {
	cmd := binio.ConfigurationControlSet{
		OutputFrittingMap: func(config Configuration) uint32 {
			var fritting uint32 = 0
			for ch, f := range config.Fritting {
				if f {
					fritting |= 1 << ch
				}
			}
			return fritting
		}(config),
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

func (c *Client) Describe() (*binio.ConfigurationControlDescribeResponse, error) {
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
