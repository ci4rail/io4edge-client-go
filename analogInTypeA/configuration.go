package analogInTypeA

import (
	"fmt"

	analogIn "github.com/ci4rail/io4edge-client-go/analogInTypeA/v1alpha1"
	"github.com/ci4rail/io4edge-client-go/functionblock"
)

type Configuration struct {
	SampleRate map[int]uint32
}

func (c *AnalogInTypeA) SetConfiguration(config Configuration) error {
	configurationCommand := analogIn.HardwareConfiguration{
		SampleRate: func(config Configuration) []*analogIn.HardwareConfigurationSampleRate {
			var analogSampleRate []*analogIn.HardwareConfigurationSampleRate
			for ch, rate := range config.SampleRate {
				analogSampleRate = append(analogSampleRate, &analogIn.HardwareConfigurationSampleRate{
					Channel:    uint32(ch),
					SampleRate: rate,
				})
			}
			return analogSampleRate
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
