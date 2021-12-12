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
		fmt.Println(err)
		return err
	}
	fmt.Println(functionblockV1.Status_name[int32(res.Status)])
	return nil
}

func (c *Client) Describe() error {
	cmd := binio.ConfigurationControlDescribe{}
	envelopeCmd, err := functionblock.ConfigurationControlDescribe(&cmd)
	if err != nil {
		return err
	}
	res := &functionblockV1.Response{}
	err = c.funcClient.Command(envelopeCmd, res, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(functionblockV1.Status_name[int32(res.Status)])
	describe := binio.ConfigurationControlResponse{}
	err = anypb.UnmarshalTo(res.GetConfigurationControl().FunctionSpecificConfigurationControlResponse, &describe, proto.UnmarshalOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Number of channels: %d\n", describe.GetDescribe().GetNumberOfChannels())
	return nil
}
