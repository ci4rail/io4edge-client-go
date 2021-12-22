package binaryIoTypeA

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationBasic(t *testing.T) {
	assert := assert.New(t)
	client := NewClient(nil)
	config := Configuration{
		OutputFritting:        map[int]bool{0: false, 1: true, 2: true, 3: false},
		OutputWatchdog:        map[int]bool{0: true, 1: true, 2: true, 3: true},
		OutputWatchdogTimeout: []int{-1},
	}
	err := client.SetConfiguration(config)
	assert.Nil(err)
}
