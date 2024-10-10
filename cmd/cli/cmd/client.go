package cmd

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/core"
	pbcore "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/core"
)

const (
	coreServiceType = "_io4edge-core._tcp"
)

// newCliClientFromService creates the io4edge core client from the device address
func newCliClientFromService(deviceID string) (core.If, error) {
	serviceAddr := deviceID + "." + coreServiceType
	c, err := pbcore.NewClientFromService(serviceAddr, time.Duration(timeoutSecs)*time.Second)
	return c, err
}

// newCliClientFromIP creates the io4edge core client from the ip address and the port
func newCliClientFromIP(ipAddrPort string) (core.If, error) {
	c, err := pbcore.NewClientFromSocketAddress(ipAddrPort)
	return c, err
}

// newCliClient creates the io4edge cli client from either the ip address and port or from the device address,
// depending on which parameter is given.
func newCliClient(deviceID string, ipAddrPort string) (core.If, error) {
	var c core.If
	var err error

	if deviceID != "" {
		c, err = newCliClientFromService(deviceID)
	} else if ipAddrPort != "" {
		c, err = newCliClientFromIP(ipAddrPort)
	} else {
		err = errors.New("no device specified (either device id or ip address required)")
	}

	return c, err
}
