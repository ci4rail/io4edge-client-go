package binaryIoTypeA

import (
	"github.com/ci4rail/io4edge-client-go/functionblock"
	binIo "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("cmd: %+v\n", cmd)
	log.Debugf("envelopeCmd: %+v\n", envelopeCmd)
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
	log.Debugf("cmd: %+v\n", cmd)
	log.Debugf("envelopeCmd: %+v\n", envelopeCmd)
	return nil
}
