package binaryIoTypeA

import (
	"fmt"

	binIo "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
)

// SetBinaryChannel sets the binary channel to the given value
func (c *Client) SetBinaryChannel(channel int, state bool) error {

	cmd := binIo.FunctionControlSet{
		Type: &binIo.FunctionControlSet_Single{
			Single: &binIo.SetSingle{
				Channel: uint32(channel),
				State:   state,
			},
		},
	}
	envelopeCmd, err := functionblock.ConfigurationControlSet(&cmd)
	if err != nil {
		return err
	}
	fmt.Println(cmd)
	fmt.Println(envelopeCmd)
	return nil
}

// SetAllBinaryChannels sets all binary channels to the given value from the bitmask
// LSB is channel 0. True: output is on, False, output is off
// For iou01 there are only four output channels. All other bits are ignored
func (c *Client) SetAllBinaryChannels(output uint32) error {
	cmd := binIo.FunctionControlSet{
		Type: &binIo.FunctionControlSet_All{
			All: &binIo.SetAll{
				Values: output,
			},
		},
	}
	envelopeCmd, err := functionblock.ConfigurationControlSet(&cmd)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", cmd)
	fmt.Printf("%+v\n", envelopeCmd)
	return nil
}
