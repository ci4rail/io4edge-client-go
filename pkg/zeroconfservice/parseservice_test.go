package zeroconfservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInstanceAndServiceOk(t *testing.T) {
	instance, service, err := ParseInstanceAndService("iou04-usb-ext-4._io4edge-core._tcp")
	assert.Nil(t, err)
	assert.Equal(t, "iou04-usb-ext-4", instance)
	assert.Equal(t, "_io4edge-core._tcp", service)

	instance, service, err = ParseInstanceAndService("foo.bar.baz._io4edge-core._tcp")
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar.baz", instance)
	assert.Equal(t, "_io4edge-core._tcp", service)
}

func TestParseInstanceAndServiceIncomplete(t *testing.T) {
	_, _, err := ParseInstanceAndService("_io4edge-core._tcp")
	assert.NotNil(t, err)

	_, _, err = ParseInstanceAndService("")
	assert.NotNil(t, err)
}
