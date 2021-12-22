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
	envelopeCmd, err := functionblock.FunctionControlSet(&cmd, string(cmd.Type.(*binio.FunctionControlSet_Single).Single.ProtoReflect().Descriptor().FullName()))
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

func (c *Client) SetAll(values uint32, mask uint32) error {
	cmd := binio.FunctionControlSet{
		Type: &binio.FunctionControlSet_All{
			All: &binio.SetAll{
				Values: values,
				Mask:   mask,
			},
		},
	}
	envelopeCmd, err := functionblock.FunctionControlSet(&cmd, string(cmd.Type.(*binio.FunctionControlSet_All).All.ProtoReflect().Descriptor().FullName()))
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
