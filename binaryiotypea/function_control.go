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

// SetSingle sets a single binary output
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
		if res.Status != functionblockV1.Status_OK {
			return fmt.Errorf(res.Error.Error)
		}
		return nil
	}
	return fmt.Errorf("not connected")
}

// SetAll sets all binary outputs in mask
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
		if res.Status != functionblockV1.Status_OK {
			return fmt.Errorf(res.Error.Error)
		}
		return nil
	}
	return fmt.Errorf("not connected")
}

// GetSingle get the value of channel
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
		if res.Status != functionblockV1.Status_OK {
			return false, fmt.Errorf(res.Error.Error)
		}
		get := binio.FunctionControlGetResponse{}
		err = anypb.UnmarshalTo(res.GetFunctionControl().GetFunctionSpecificFunctionControlGet(), &get, proto.UnmarshalOptions{})
		if err != nil {
			return false, err
		}
		return get.GetSingle().State, nil
	}
	return false, fmt.Errorf("not connected")
}

// GetAll gets the value of all channels in mask
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
		if res.Status != functionblockV1.Status_OK {
			return 0, fmt.Errorf(res.Error.Error)
		}
		get := binio.FunctionControlGetResponse{}
		err = anypb.UnmarshalTo(res.GetFunctionControl().GetFunctionSpecificFunctionControlGet(), &get, proto.UnmarshalOptions{})
		if err != nil {
			return 0, err
		}
		return get.GetAll().Inputs, nil
	}
	return 0, fmt.Errorf("not connected")
}
