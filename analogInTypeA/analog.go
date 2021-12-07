package analogInTypeA

import (
	"fmt"

	anaIn "github.com/ci4rail/io4edge-client-go/analogInTypeA/v1alpha1"
)

type AnalogInTypeA struct {
}

func NewAnalogInTypeAClient() *AnalogInTypeA {
	return &AnalogInTypeA{}
}

// SetBinaryChannel sets the binary channel to the given value
func (c *AnalogInTypeA) GetChannel(channel int) (*anaIn.Sample, error) {

	cmd := anaIn.FunctionControlGet{
		Type: &anaIn.FunctionControlGet_Single{
			Single: &anaIn.GetSingle{
				Channel: uint32(channel),
			},
		},
	}
	fmt.Println(cmd)
	return nil, nil
}
