package core

import (
	"errors"
	"log"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

func searchService(sb *avahi.ServiceBrowser, srv *avahi.Server, name string, t int) (avahi.Service, error) {
	var s avahi.Service
	for {
		select {
		case s = <-sb.AddChannel:
			s, err := srv.ResolveService(s.Interface, s.Protocol, s.Name,
				s.Type, s.Domain, avahi.ProtoUnspec, 0)
			if err != nil {
				log.Fatalf("ResolveService() failed: %v", err)
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

func GetAddressFromService(instanceName string, serviceName string, timeout int) (string, uint16, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Cannot get system bus: %v", err)
	}

	server, err := avahi.ServerNew(conn)
	if err != nil {
		log.Fatalf("Avahi new failed: %v", err)
	}

	sb, err := server.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, serviceName, "local", 0)
	if err != nil {
		log.Fatalf("ServiceBrowserNew() failed: %v", err)
	}

	service, err := searchService(sb, server, instanceName, timeout)

	return service.Address, service.Port, err
}
