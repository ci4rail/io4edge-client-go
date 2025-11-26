package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertFirmwareIsCompatibleWithHardware(t *testing.T) {
	assert.Nil(t, AssertFirmwareIsCompatibleWithHardware("S101-IOU01", []int{1, 2, 3}, "s101-IOU01", 1))
	assert.NotNil(t, AssertFirmwareIsCompatibleWithHardware("S101-IOU01", []int{1, 2, 3}, "s101-IOU04", 1))
	assert.NotNil(t, AssertFirmwareIsCompatibleWithHardware("S101-IOU01", []int{1, 2, 3}, "s101-IOU01", 4))
	assert.Nil(t, AssertFirmwareIsCompatibleWithHardware("s103-sio06", []int{1, 2}, "S103-SIO06-02-00001", 2))
	assert.NotNil(t, AssertFirmwareIsCompatibleWithHardware("s103-sio06-02-0001", []int{1, 2}, "S103-SIO06-02", 2))
}
