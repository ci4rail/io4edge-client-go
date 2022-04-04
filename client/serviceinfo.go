package client

import (
	"errors"
	"strconv"
	"strings"
	"sync"
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

type avahiServer interface {
	ServiceBrowserNew(iface, protocol int32, serviceType string, domain string, flags uint32) (*avahi.ServiceBrowser, error)
	ServiceTypeBrowserNew(iface, protocol int32, domain string, flags uint32) (*avahi.ServiceTypeBrowser, error)
	ResolveService(iface, protocol int32, name, serviceType, domain string, aprotocol int32, flags uint32) (reply avahi.Service, err error)
}

var server avahiServer

// ServiceInfo stores the avahi service struct of a service to make information about the service available through getter functions
type ServiceInfo struct {
	service avahi.Service
}

type serviceFromInstance map[string]ServiceInfo

type observerContext struct {
	foundInstances serviceFromInstance
	channels       []chan ServiceInfo
}

// +-----------------------------+
// |         observers           |
// +=============================+
// |   key: service name         |
// +-----------------------------+
// |           value:            |
// | +-------------------------+ |
// | |    observer context     | |
// | +=========================+ |
// | | +---------------------+ | |
// | | |   found instances   | | |
// | | +=====================+ | |
// | | | key: instance name  | | |
// | | +---------------------+ | |
// | | | value: service info | | |
// | | +---------------------+ | |
// | +-------------------------+ |
// | |     channels[]          | |
// | +-------------------------+ |
// +-----------------------------+
var observers = make(map[string]*observerContext)

// mutex to lock access to observers
var observersMutex sync.Mutex

func (o *observerContext) newInstanceFound(s ServiceInfo) error {
	observersMutex.Lock()
	defer observersMutex.Unlock()

	o.foundInstances[s.service.Name] = s

	for _, channel := range o.channels {
		select {
		case channel <- s:
		default:
		}
	}

	return nil
}

func (o *observerContext) instanceDisappeared(s ServiceInfo) error {
	observersMutex.Lock()
	defer observersMutex.Unlock()

	_, ok := o.foundInstances[s.service.Name]
	if ok {
		delete(o.foundInstances, s.service.Name)
	}

	return nil
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

	if server == nil {
		conn, err := dbus.SystemBus()
		if err != nil {
			return err
		}

		server, err = avahi.ServerNew(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeChannel(serviceName string, channel chan ServiceInfo) error {
	observersMutex.Lock()
	defer observersMutex.Unlock()
	for idx, c := range observers[serviceName].channels {
		if c == channel {
			observers[serviceName].channels = append(observers[serviceName].channels[:idx], observers[serviceName].channels[idx+1:]...)
			return nil
		}
	}

	err := errors.New("error: could not find channel in channelMap")
	return err
}

// GetServiceInfo creates a new avahi server if necessary and starts a new service observer instance if
// no one is already running for the given mdns service name.
// If a service with the correct service address (consisting of <instance_name>.<service_name>.<protocol>) is
// found within the specified timeout, a service info object is returned.
// The instanceName should contain the instance name of the service address
// The serviceName should contain the service name of the service address together with the protocol
// The timeout specifies the time to wait for the service to show up
func GetServiceInfo(instanceName string, serviceName string, timeout time.Duration) (*ServiceInfo, error) {
	var svcInf ServiceInfo
	var err error
	var channel chan ServiceInfo
	startObserver := false

	err = initAvahiServer()
	if err != nil {
		return nil, err
	}

	/* Avoid concurrent access to observerList */
	observersMutex.Lock()

	/* check instance map for service name and instance name */
	observer, exists := observers[serviceName]
	if exists {
		svcInf, exists = observer.foundInstances[instanceName]
		if exists {
			return &svcInf, nil
		}
	} else {
		observers[serviceName] = new(observerContext)
		observers[serviceName].foundInstances = make(serviceFromInstance)
		startObserver = true
	}

	/* create channel to get service info object when service observer found service */
	channel = make(chan ServiceInfo)
	observers[serviceName].channels = append(observers[serviceName].channels, channel)
	defer removeChannel(serviceName, channel)

	if startObserver {
		/* start service observer and pass observer context as context */
		o := observers[serviceName]
		go ServiceObserver(serviceName, o.newInstanceFound, o.instanceDisappeared)
	}
	observersMutex.Unlock()

	for {
		select {
		case svcInf = <-channel:
			if svcInf.GetInstanceName() == instanceName {
				return &svcInf, nil
			}
		case <-time.After(timeout):
			err = errors.New("error: could not find instance or service (" + instanceName + "." + serviceName + "): timeout")
			return nil, err
		}
	}
}

// NetAddress gives the caller the ip address and port of the service
func (svcInf *ServiceInfo) NetAddress() (string, int, error) {
	return svcInf.service.Address, int(svcInf.service.Port), nil
}

// FuncClass gives the caller the funcclass value of the service
func (svcInf *ServiceInfo) FuncClass() (string, error) {
	value, err := getTxtValueFromKey(funcclass, svcInf.service.Txt)
	return value, err
}

// Security gives the caller the security value of the service
func (svcInf *ServiceInfo) Security() (string, error) {
	value, err := getTxtValueFromKey(security, svcInf.service.Txt)
	return value, err
}

func auxPort(txt string) (string, int, error) {
	field := strings.Split(string(txt), "-")
	if len(field) != 2 {
		return "", 0, errors.New("invalid txt field structure of field " + auxport)
	}

	protocol := field[0]
	port, err := strconv.Atoi(field[1])
	return protocol, port, err
}

// AuxPort gives the caller the auxport value of the service protocol and port
func (svcInf *ServiceInfo) AuxPort() (string, int, error) {
	value, err := getTxtValueFromKey(auxport, svcInf.service.Txt)
	if err != nil {
		return "", 0, err
	}

	protocol, port, err := auxPort(value)
	if err != nil {
		return "", 0, err
	}
	if protocol != "tcp" && protocol != "udp" {
		return "", 0, errors.New("no aux port")
	}
	return protocol, port, nil
}

// AuxSchemaID gives the caller the auxschema value of the service
func (svcInf *ServiceInfo) AuxSchemaID() (string, error) {
	value, err := getTxtValueFromKey(auxschema, svcInf.service.Txt)
	if err != nil {
		return "", err
	}
	if value == "" || value == "not_avail" {
		return "", errors.New("no aux schema")
	}

	return value, nil
}

// GetIPAddressPort gets the ip address and port from the given service info object.
// It passes the caller the ip address and the port separated by ":" together in a string.
func (svcInf *ServiceInfo) GetIPAddressPort() string {
	ipAddress, port, _ := svcInf.NetAddress()
	ipAddrPort := ipAddress + ":" + strconv.Itoa(port)
	return ipAddrPort
}

// GetInstanceName gets the instance name of the given service info object.
func (svcInf *ServiceInfo) GetInstanceName() string {
	return svcInf.service.Name
}

// GetServiceType gets the service type of the given service info object
func (svcInf *ServiceInfo) GetServiceType() string {
	return svcInf.service.Type
}
