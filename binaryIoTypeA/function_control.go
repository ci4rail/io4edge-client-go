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

func (c *Client) SetSingle(channel uint, state bool) error {
	if c.connected {
		cmd := binio.FunctionControlSet{
			Type: &binio.FunctionControlSet_Single{
				Single: &binio.SetSingle{
					Channel: uint32(channel),
					State:   state,
				},
			},
		}
		envelopeCmd, err := functionblock.FunctionControlSet(&cmd)
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
		return nil
	}
	return fmt.Errorf("not connected")
}

func (c *Client) SetAll(values uint32, mask uint32) error {
	if c.connected {
		cmd := binio.FunctionControlSet{
			Type: &binio.FunctionControlSet_All{
				All: &binio.SetAll{
					Values: values,
					Mask:   mask,
				},
			},
		}
		envelopeCmd, err := functionblock.FunctionControlSet(&cmd)
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
		return nil
	}
	return fmt.Errorf("not connected")
}

func (c *Client) GetSingle(channel uint) (bool, error) {
	if c.connected {
		cmd := binio.FunctionControlGet{
			Type: &binio.FunctionControlGet_Single{
				Single: &binio.GetSingle{
					Channel: uint32(channel),
				},
			},
		}
		envelopeCmd, err := functionblock.FunctionControlGet(&cmd)
		if err != nil {
			return false, err
		}
		res, err := c.Command(envelopeCmd, time.Second*5)
		if err != nil {
			return false, err
		}
		if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
			return false, fmt.Errorf("not implemented")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return false, fmt.Errorf(res.Error.Error)
		}
		get := binio.FunctionControlResponse{}
		err = anypb.UnmarshalTo(res.GetFunctionControl().FunctionSpecificFunctionControlResponse, &get, proto.UnmarshalOptions{})
		if err != nil {
			return false, err
		}
		return get.GetGetSingle().State, nil
	}
	return false, fmt.Errorf("not connected")
}

func (c *Client) GetAll(mask uint32) (uint32, error) {
	if c.connected {
		cmd := binio.FunctionControlGet{
			Type: &binio.FunctionControlGet_All{
				All: &binio.GetAll{
					Mask: mask,
				},
			},
		}
		envelopeCmd, err := functionblock.FunctionControlGet(&cmd)
		if err != nil {
			return 0, err
		}
		res, err := c.Command(envelopeCmd, time.Second*5)
		if err != nil {
			return 0, err
		}
		if res.Status == functionblockV1.Status_NOT_IMPLEMENTED {
			return 0, fmt.Errorf("not implemented")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return 0, fmt.Errorf(res.Error.Error)
		}
		get := binio.FunctionControlResponse{}
		err = anypb.UnmarshalTo(res.GetFunctionControl().FunctionSpecificFunctionControlResponse, &get, proto.UnmarshalOptions{})
		if err != nil {
			return 0, err
		}
		return get.GetGetAll().Inputs, nil
	}
	return 0, fmt.Errorf("not connected")
}
