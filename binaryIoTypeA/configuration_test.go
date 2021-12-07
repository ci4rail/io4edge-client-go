package binaryIoTypeA

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationBasic(t *testing.T) {
	assert := assert.New(t)
	client := NewBinaryIoTypeAClient()
	config := Configuration{
		Fritting: map[int]bool{0: false, 1: true, 2: true, 3: false},
	}
	err := client.SetConfiguration(config)
	assert.Nil(err)
}
