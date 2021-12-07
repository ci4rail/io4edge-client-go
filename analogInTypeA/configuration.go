package analogInTypeA

import (
	"fmt"

	analogIn "github.com/ci4rail/io4edge-client-go/analogInTypeA/v1alpha1"
	functionblock "github.com/ci4rail/io4edge-client-go/functionblock"
)

type Configuration struct {
	SampleRate map[int]uint32
}

func (c *Client) SetConfiguration(config Configuration) error {
	configurationCommand := analogIn.ConfigurationControlSet{
		SampleRate: func(config Configuration) []*analogIn.ConfigurationSampleRate {
			var analogSampleRate []*analogIn.ConfigurationSampleRate
			for ch, rate := range config.SampleRate {
				analogSampleRate = append(analogSampleRate, &analogIn.ConfigurationSampleRate{
					Channel:    uint32(ch),
					SampleRate: rate,
				})
			}
			return analogSampleRate
		}(config),
	}
	envelopeCmd, err := functionblock.ConfigurationControlSet(&configurationCommand)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", configurationCommand)
	fmt.Printf("%+v\n", envelopeCmd)
	return nil
}
