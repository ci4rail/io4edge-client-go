package cmd

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/core"
	"github.com/ci4rail/io4edge-client-go/v2/pkg/iputil"
	pbcore "github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/core"
	restcore "github.com/ci4rail/io4edge-client-go/v2/pkg/restcom/core"
)

const (
	coreServiceType = "_io4edge-core._tcp"
	restCorePort    = 443
)

// newCliClientFromService creates the io4edge core client from the device address
func newCliClientFromService(deviceID string) (core.If, error) {
	serviceAddr := deviceID + "." + coreServiceType
	c, err := pbcore.NewClientFromService(serviceAddr, time.Duration(timeoutSecs)*time.Second)
	return c, err
}

// newCliClientFromIP creates the io4edge core client from the ip address and the port
func newCliClientFromIP(ipAddrPort string) (core.If, error) {
	var c core.If

	_, port, err := iputil.NetAddressSplit(ipAddrPort)
	if err != nil {
		return nil, err
	}
	if port%1000 == restCorePort {
		if password == "" {
			return nil, errors.New("password required for REST API")
		}
		c, err = restcore.NewClientFromSocketAddress(ipAddrPort, password)
	} else {
		c, err = pbcore.NewClientFromSocketAddress(ipAddrPort)
	}
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
