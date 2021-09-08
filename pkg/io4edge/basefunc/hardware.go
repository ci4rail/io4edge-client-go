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
