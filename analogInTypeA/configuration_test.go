package analogInTypeA

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationBasic(t *testing.T) {
	assert := assert.New(t)
	client := NewAnalogInTypeAClient()
	config := Configuration{
		SampleRate: map[int]uint32{0: 100, 1: 200, 2: 300, 3: 400},
	}
	err := client.SetConfiguration(config)
	assert.Nil(err)
}
