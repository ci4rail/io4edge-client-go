package core

import (
	"fmt"
	"strconv"
	"time"
)

// getIPAddressPort gets the ip address and port from the given service info object.
// It passes the caller the ip address and the port separated by ":" together in a string.
func getIPAddressPort(svcInfo *ServiceInfo) string {
	ipAddress, port := svcInfo.GetNetAddress()
	ipAddrPort := ipAddress + ":" + strconv.Itoa(int(port))
	return ipAddrPort
}

// NewClientFromService creates a new base function client from a socket with a address, which was acquired from the specified service.
// The timeout specifies the maximal time waiting for a service to show up.
func NewClientFromService(instance string, service string, timeout time.Duration) (*Client, error) {
	svcInfo, err := NewServiceInfo(instance, service, timeout)
	if err != nil {
		return nil, err
	}
	ipAddrPort := getIPAddressPort(svcInfo)
	fmt.Println("Try to connect with ", ipAddrPort)
	c, err := NewClientFromSocketAddress(ipAddrPort)
	return c, err
}
