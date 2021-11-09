package client

import (
	"errors"
	"strings"
)

// ParseInstanceAndService parses one string, which describes a service and split it up into instance name and service name
func ParseInstanceAndService(serviceAddr string) (string, string, error) {
	s := strings.Split(serviceAddr, ".")
	if len(s) < 3 {
		err := errors.New("service address not parseable (one of these are missing: instance, service, protocol)")
		return "", "", err
	}

	service := s[len(s)-2] + "." + s[len(s)-1]

	// instance name may contain dots. Repair them
	instance := s[0]
	for i := 1; i < len(s)-2; i++ {
		instance += "." + s[i]
	}

	return instance, service, nil
}
