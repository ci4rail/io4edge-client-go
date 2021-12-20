package binaryIoTypeA

import (
	"fmt"
	"time"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
)

func (c *Client) SetSingle(channel uint, state bool) error {
	cmd := binio.FunctionControlSet{
		Type: &binio.FunctionControlSet_Single{
			Single: &binio.SetSingle{
				Channel: uint32(channel),
				State:   state,
			},
		},
	}
	single := binio.SetSingle{}
	envelopeCmd, err := functionblock.FunctionControlSet(&cmd, string(single.ProtoReflect().Descriptor().FullName()))
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

func (c *Client) SetAll(channel uint, values uint32, mask uint32) error {
	cmd := binio.FunctionControlSet{
		Type: &binio.FunctionControlSet_All{
			All: &binio.SetAll{
				Values: values,
				Mask:   mask,
			},
		},
	}
	envelopeCmd, err := functionblock.FunctionControlSet(&cmd, "setAll")
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
