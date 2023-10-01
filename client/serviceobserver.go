package client

import (
	"fmt"

	"github.com/gobwas/glob"
	"github.com/holoplot/go-avahi"
	log "github.com/sirupsen/logrus"
)

func runServiceBrowser(domain string, serviceType string, addServiceChan chan ServiceInfo, removeServiceChan chan ServiceInfo) {
	log.Debugf("starting service Browser for %s", serviceType)
	if serviceType == "" {
		return
	}
	sb, err := server.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, serviceType, domain, 0)
	if err != nil {
		log.Errorf("servicescan: can't start service browser for %s: %v", serviceType, err)
		return
	}
	for {
		var svcInf ServiceInfo
		select {
		case s := <-sb.AddChannel:
			log.Debugf("browser got add service %s if=%v prot=%v type=%v domain=%v", s.Name, s.Interface, s.Protocol, s.Type, s.Domain)
			go func() {
				s, err = server.ResolveService(s.Interface, s.Protocol, s.Name,
					s.Type, s.Domain, avahi.ProtoUnspec, 0)
				if err == nil {
					log.Debugf("browser got resolved service %s", s.Name)
					svcInf.service = s
					addServiceChan <- svcInf
				} else {
					log.Debugf("servicescan: can't resolve service %s: %v", s.Name, err)
				}
			}()
		case s := <-sb.RemoveChannel:
			log.Debugf("browser got rem service %s", s.Name)
			svcInf.service = s
			removeServiceChan <- svcInf
		}
	}
}

// ServiceObserver watches for added or removed services whose serviceTypes are matching the serviceNamePattern
// serviceNamePattern is compared as a glob pattern, i.e. if you want to observe service types beginning with "_io4edge", specify
// "_io4edge.*"; if you want to observe all services, specify "*".
//
// serviceObserver runs in a loop until one of the callbacks returns an error.
// serviceAdded and serviceRemoved are called when an instance of an observed service type is added or removed.
func ServiceObserver(serviceNamePattern string, serviceAdded func(ServiceInfo) error, serviceRemoved func(ServiceInfo) error) error {
	startedServiceBrowsers := make(map[string]struct{}, 0)
	addServiceChan := make(chan ServiceInfo)
	removeServiceChan := make(chan ServiceInfo)

	log.Debugf("starting service observer for %s", serviceNamePattern)
	g, err := glob.Compile(serviceNamePattern)
	if err != nil {
		err = fmt.Errorf("service name pattern pattern invalid: %v", err)
		return err
	}

	err = initAvahiServer()
	if err != nil {
		return err
	}

	stb, err := server.ServiceTypeBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, "local", 0)
	if err != nil {
		return err
	}

	for {
		select {
		case serviceType := <-stb.AddChannel:
			log.Debugf("ServiceObserver add %s", serviceType.Type)

			if g.Match(serviceType.Type) {
				_, present := startedServiceBrowsers[serviceType.Type]
				if !present {
					startedServiceBrowsers[serviceType.Type] = struct{}{}
					go runServiceBrowser(serviceType.Domain, serviceType.Type, addServiceChan, removeServiceChan)
				}
			}

		case serviceType := <-stb.RemoveChannel:
			log.Debugf("browser got remove service type service %s", serviceType.Type)

		case svcInf := <-addServiceChan:
			log.Debugf("ServiceObserver serviceAdded add %s", svcInf.service.Name)
			if err := serviceAdded(svcInf); err != nil {
				return err
			}
		case svcInf := <-removeServiceChan:
			if err := serviceRemoved(svcInf); err != nil {
				return err
			}
		}
	}
}
