package iou01

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationBasic(t *testing.T) {
	assert := assert.New(t)
	client := NewIou01Client()
	config := Configuration{
		Fritting:         []bool{false, true, true, false},
		AnalogSampleRate: []uint32{100, 200, 300, 400},
	}
	err := client.SetConfiguration(config)
	assert.Nil(err)
}
