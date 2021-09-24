package core

import (
	"fmt"
	"strconv"
)

// getIPAddressPort gets the ip address and port from a service specified by the given service name and instance name.
// It passes the caller the ip address and the port separated by ":" together in a string
func getIPAddressPort(instanceName string, serviceName string, timeout int) (string, error) {
	ipAddress, port, err := GetAddressFromService(instanceName, serviceName, timeout)
	if err != nil {
		return "", err
	}
	ipAddrPort := ipAddress + ":" + strconv.Itoa(int(port))

	return ipAddrPort, err
}

// NewClientFromService creates a new base function client from a socket with a address, which was acquired from the specified service
func NewClientFromService(instance string, service string) (*Client, error) {
	ipAddrPort, err := getIPAddressPort(instance, service, 3)
	if err != nil {
		return nil, err
	}
	fmt.Println("Try to connect with ", ipAddrPort)
	c, err := NewClientFromSocketAddress(ipAddrPort)
	return c, err
}
