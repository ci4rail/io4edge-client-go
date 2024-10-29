package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertFirmwareIsCompatibleWithHardware(t *testing.T) {
	assert.Nil(t, AssertFirmwareIsCompatibleWithHardware("S101-IOU01", []int{1, 2, 3}, "s101-IOU01", 1))
	assert.NotNil(t, AssertFirmwareIsCompatibleWithHardware("S101-IOU01", []int{1, 2, 3}, "s101-IOU04", 1))
	assert.NotNil(t, AssertFirmwareIsCompatibleWithHardware("S101-IOU01", []int{1, 2, 3}, "s101-IOU01", 4))
}
