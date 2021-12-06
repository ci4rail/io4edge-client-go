package iou01

import (
	"fmt"

	iou01v1 "github.com/ci4rail/io4edge-client-go/iou01/v1alpha1"
)

// SetBinaryChannel sets the binary channel to the given value
func (c *Iou01) SetBinaryChannel(channel int, state bool) error {

	cmd := iou01v1.FunctionControlSet{
		Type: &iou01v1.FunctionControlSet_Single{
			Single: &iou01v1.SetSingle{
				Channel: iou01v1.BinaryOutput(channel),
				State:   state,
			},
		},
	}
	fmt.Println(cmd)
	return nil
}

// SetAllBinaryChannels sets all binary channels to the given value from the bitmask
// LSB is channel 0. True: output is on, False, output is off
// For iou01 there are only four output channels. All other bits are ignored
func (c *Iou01) SetAllBinaryChannels(output uint32) error {
	cmd := iou01v1.FunctionControlSet{
		Type: &iou01v1.FunctionControlSet_All{
			All: &iou01v1.SetAll{
				Values: output,
			},
		},
	}
	fmt.Println(cmd)
	return nil
}
