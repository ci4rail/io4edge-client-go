package analogInTypeA

import (
	"fmt"

	anaIn "github.com/ci4rail/io4edge-client-go/analogInTypeA/v1alpha1"
	functionblock "github.com/ci4rail/io4edge-client-go/functionblock"
)

type AnalogInTypeA struct {
}

func NewAnalogInTypeAClient() *AnalogInTypeA {
	return &AnalogInTypeA{}
}

// GetChannel gets the current value of a given analog channel
func (c *AnalogInTypeA) GetChannel(channel int) (*anaIn.Sample, error) {

	cmd := anaIn.FunctionControlGet{
		Type: &anaIn.FunctionControlGet_Single{
			Single: &anaIn.GetSingle{
				Channel: uint32(channel),
			},
		},
	}
	envelopeCmd, err := functionblock.ConfigurationControlSet(&cmd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", cmd)
	fmt.Printf("%+v\n", envelopeCmd)
	return nil, nil
}
