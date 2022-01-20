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

func (c *Client) Set(value uint32) error {
	if c.connected {
		cmd := templateModule.FunctionControlSet{
			Value: value,
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

func (c *Client) Get() (uint32, error) {
	if c.connected {
		cmd := templateModule.FunctionControlGet{}
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
		if res.Status == functionblockV1.Status_WRONG_CLIENT {
			return 0, fmt.Errorf("wrong client")
		}
		if res.Status == functionblockV1.Status_ERROR {
			return 0, fmt.Errorf(res.Error.Error)
		}
		get := templateModule.FunctionControlResponse{}
		err = anypb.UnmarshalTo(res.GetFunctionControl().FunctionSpecificFunctionControlResponse, &get, proto.UnmarshalOptions{})
		if err != nil {
			return 0, err
		}
		return get.GetGet().Value, nil
	}
	return 0, fmt.Errorf("not connected")
}
