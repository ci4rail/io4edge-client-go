package basefunc

import (
	"time"
)

// IdentifyHardware gets the hardware inventory data from the device
func (c *Client) IdentifyHardware(timeout time.Duration) (*ResIdentifyHardware, error) {
	cmd := &BaseFuncCommand{
		Id: BaseFuncCommandId_IDENTIFY_HARDWARE,
	}
	res := &BaseFuncResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return nil, err
	}
	return res.GetIdentifyHardware(), nil
}

// ProgramHardwareIdentification programs hardware inventory data into the device.
// Intended to be used during hardware manufacturing process only
func (c *Client) ProgramHardwareIdentification(id *CmdProgramHardwareIdentification, timeout time.Duration) error {
	cmd := &BaseFuncCommand{
		Id: BaseFuncCommandId_PROGRAM_HARDWARE_IDENTIFICATION,
		Data: &BaseFuncCommand_ProgramHardwareIdentification{
			ProgramHardwareIdentification: id,
		},
	}
	res := &BaseFuncResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return err
	}
	return nil
}
