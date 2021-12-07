package binaryIoTypeA

import (
	"fmt"

	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
)

type Configuration struct {
	Fritting map[int]bool
}

func (c *BinaryIoTypeA) SetConfiguration(config Configuration) error {
	configurationCommand := binio.HardwareConfiguration{
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
	envelopeCmd, err := functionblock.CreateCommand(&configurationCommand)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", configurationCommand)
	fmt.Printf("%+v\n", envelopeCmd)
	return nil
}
