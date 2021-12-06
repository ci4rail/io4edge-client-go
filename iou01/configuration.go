package iou01

import (
	"fmt"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	iou01v1 "github.com/ci4rail/io4edge-client-go/iou01/v1alpha1"
)

type Configuration struct {
	Fritting         []bool
	AnalogSampleRate []uint32
}

func (c *Iou01) SetConfiguration(config Configuration) error {
	configurationCommand := iou01v1.HardwareConfiguration{
		AnalogSampleRate: func(config Configuration) []*iou01v1.HardwareConfigurationAnalogSampleRate {
			var analogSampleRate []*iou01v1.HardwareConfigurationAnalogSampleRate
			for i := range config.AnalogSampleRate {
				analogSampleRate = append(analogSampleRate, &iou01v1.HardwareConfigurationAnalogSampleRate{
					Channel:    iou01v1.AnalogChannel(i),
					SampleRate: config.AnalogSampleRate[i],
				})
			}
			return analogSampleRate
		}(config),
		BinaryOutputFrittingMap: func(config Configuration) uint32 {
			var fritting uint32 = 0
			for i := 0; i < len(config.Fritting); i++ {
				if config.Fritting[i] {
					fritting |= 1 << uint32(i)
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
