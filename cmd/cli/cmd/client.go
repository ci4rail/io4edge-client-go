package cmd

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/coreclient"
	"github.com/ci4rail/io4edge-client-go/socketcore"
)

const (
	coreServiceType = "_io4edge-core._tcp"
)

// newCliClientFromService creates the io4edge core client from the device address
func newCliClientFromService(deviceID string) (coreclient.If, error) {
	serviceAddr := deviceID + "." + coreServiceType
	c, err := socketcore.NewClientFromService(serviceAddr, time.Duration(timeoutSecs)*time.Second)
	return c, err
}

// newCliClientFromIP creates the io4edge core client from the ip address and the port
func newCliClientFromIP(ipAddrPort string) (coreclient.If, error) {
	c, err := socketcore.NewClientFromSocketAddress(ipAddrPort)
	return c, err
}

// newCliClient creates the io4edge cli client from either the ip address and port or from the device address,
// depending on which parameter is given.
func newCliClient(deviceID string, ipAddrPort string) (coreclient.If, error) {
	var c coreclient.If
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
