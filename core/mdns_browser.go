package core

import (
	"errors"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

// searchService searches for the service with the specified instance name
func searchService(sb *avahi.ServiceBrowser, srv *avahi.Server, name string, t int) (avahi.Service, error) {
	var s avahi.Service
	for {
		select {
		case s = <-sb.AddChannel:
			s, err := srv.ResolveService(s.Interface, s.Protocol, s.Name,
				s.Type, s.Domain, avahi.ProtoUnspec, 0)
			if err != nil {
				return s, err
			}
			if s.Name == name {
				return s, err
			}
		case <-time.After(time.Duration(t) * time.Second):
			err := errors.New("could not find instance or service")
			return s, err
		}
	}
}

// GetAddressFromService starts the mdns server and browse interfaces for mdns services with the specified service name.
// Sort out the service with the specified instance name.
// Return the ip address and port of the found service
// TODO get additional ports from txt field
func GetAddressFromService(instanceName string, serviceName string, timeout int) (string, uint16, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return "", 0, err
	}

	server, err := avahi.ServerNew(conn)
	if err != nil {
		return "", 0, err
	}

	sb, err := server.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, serviceName, "local", 0)
	if err != nil {
		return "", 0, err
	}

	service, err := searchService(sb, server, instanceName, timeout)

	return service.Address, service.Port, err
}
