package client

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

const (
	funcclass = "funcclass"
	security  = "security"
	auxport   = "auxport"
	auxschema = "auxschema"
)

var server *avahi.Server

// ServiceInfo stores the avahi service struct of a service to make information about the service available through getter functions
type ServiceInfo struct {
	service avahi.Service
}

// searchService searches for the service with the specified instance and service name
func searchService(instanceName string, serviceName string, timeout time.Duration) (avahi.Service, error) {
	var s avahi.Service

	sb, err := server.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, serviceName, "local", 0)
	if err != nil {
		return s, err
	}

	for {
		select {
		case s = <-sb.AddChannel:
			s, err = server.ResolveService(s.Interface, s.Protocol, s.Name,
				s.Type, s.Domain, avahi.ProtoUnspec, 0)
			if err != nil {
				return s, err
			}
			if s.Name == instanceName {
				return s, err
			}
		case <-time.After(timeout):
			err = errors.New("could not find instance or service")
			return s, err
		}
	}
}

// getTxtValueFromKey searches the given txt array for the given key and returns the corresponding value.
func getTxtValueFromKey(key string, txt [][]byte) (string, error) {
	for idx := range txt {
		field := strings.Split(string(txt[idx]), "=")
		if field[0] == key {
			return field[1], nil
		}
	}

	err := errors.New("could not find key: " + key)
	return "", err
}

// initAvahiServer creates a new avahi server and stores it in the server variable (only one server is needed to search for multiple services)
func initAvahiServer() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	server, err = avahi.ServerNew(conn)
	if err != nil {
		return err
	}

	return nil
}

// NewServiceInfo creates a new avahi server if necessary, browses interfaces for the specified mdns service and returns a service info object
// The service address consists of <instance_name>.<service_name>.<protocol>
// The instanceName should contain the instance name of the service address
// The serviceName should contain the service name of the service address together with the protocol
// The timeout specifies the time to wait for the service to show up
func NewServiceInfo(instanceName string, serviceName string, timeout time.Duration) (*ServiceInfo, error) {
	var svcInf ServiceInfo
	var err error

	if server == nil {
		err = initAvahiServer()
		if err != nil {
			return nil, err
		}
	}

	svcInf.service, err = searchService(instanceName, serviceName, timeout)
	if err != nil {
		return nil, err
	}

	return &svcInf, nil
}

// GetNetAddress gives the caller the ip address and port of the service
func (svcInf *ServiceInfo) GetNetAddress() (string, uint16) {
	return svcInf.service.Address, svcInf.service.Port
}

// GetFuncclass gives the caller the funcclass value of the service
func (svcInf *ServiceInfo) GetFuncclass() (string, error) {
	value, err := getTxtValueFromKey(funcclass, svcInf.service.Txt)
	return value, err
}

// GetSecurity gives the caller the security value of the service
func (svcInf *ServiceInfo) GetSecurity() (string, error) {
	value, err := getTxtValueFromKey(security, svcInf.service.Txt)
	return value, err
}

// GetAuxport gives the caller the auxport value of the service splitted into protocol and port
func (svcInf *ServiceInfo) GetAuxport() (string, int, error) {
	value, err := getTxtValueFromKey(auxport, svcInf.service.Txt)
	if err != nil {
		return "", 0, err
	}

	field := strings.Split(string(value), "-")
	if len(field) != 2 {
		err = errors.New("invalid txt field structure of field " + auxport)
		return "", 0, err
	}

	protocol := field[0]
	port, err := strconv.Atoi(field[1])
	return protocol, port, err
}

// GetAuxschema gives the caller the auxschema value of the service
func (svcInf *ServiceInfo) GetAuxschema() (string, error) {
	value, err := getTxtValueFromKey(auxschema, svcInf.service.Txt)
	return value, err
}
