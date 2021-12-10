package binaryIoTypeA

import (
	"fmt"
	"time"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	functionblockV1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
)

type Configuration struct {
	Fritting map[int]bool
}

func (c *Client) SetConfiguration(config Configuration) error {
	cmd := binio.ConfigurationControlSet{
		OutputFrittingMap: func(config Configuration) uint32 {
			var fritting uint32 = 0
			for ch, f := range config.Fritting {
				if f {
					fritting |= 1 << ch
				}
			}
			return fritting
		}(config),
	}
	envelopeCmd, err := functionblock.ConfigurationControlSet(&cmd)
	if err != nil {
		return err
	}
	res := &functionblockV1.Response{}
	err = c.funcClient.Command(envelopeCmd, res, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(functionblockV1.Status_name[int32(res.Status)])
	return nil
}
