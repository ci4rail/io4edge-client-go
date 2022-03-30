package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/holoplot/go-avahi"
	"github.com/stretchr/testify/assert"
)

func TestGetTxtValueFromKey(t *testing.T) {
	txt := [][]byte{[]byte("func=abc"), []byte("delta=def")}

	v, err := getTxtValueFromKey("func", txt)
	assert.Nil(t, err)
	assert.Equal(t, "abc", v)

	v, err = getTxtValueFromKey("delta", txt)
	assert.Nil(t, err)
	assert.Equal(t, "def", v)

	_, err = getTxtValueFromKey("gamma", txt)
	assert.NotNil(t, err)
}

func TestAuxPort(t *testing.T) {
	protocol, port, err := auxPort("tcp-10000")
	assert.Nil(t, err)
	assert.Equal(t, protocol, "tcp")
	assert.Equal(t, port, 10000)

	_, _, err = auxPort("10000")
	assert.NotNil(t, err)

	_, _, err = auxPort("tcp-1A000")
	assert.NotNil(t, err)

	_, _, err = auxPort("not_avail")
	assert.NotNil(t, err)

}

type SimServer struct {
	scenarioPlayer func(*avahi.ServiceBrowser)
}

func scenario1(sb *avahi.ServiceBrowser) {
	time.Sleep(time.Microsecond * 100)
	sb.AddChannel <- avahi.Service{
		Protocol: 0,
		Type:     "_foo._tcp",
		Name:     "fooA",
	}
	time.Sleep(time.Microsecond * 50)
	sb.AddChannel <- avahi.Service{
		Protocol: 1,
		Type:     "_foo._tcp",
		Name:     "fooB",
	}

	time.Sleep(time.Microsecond * 100)

	sb.RemoveChannel <- avahi.Service{
		Protocol: 1,
		Type:     "_foo._tcp",
		Name:     "fooB",
	}
	sb.RemoveChannel <- avahi.Service{
		Protocol: 0,
		Type:     "_foo._tcp",
		Name:     "fooA",
	}
}

func (s *SimServer) ServiceBrowserNew(iface, protocol int32, serviceType string, domain string, flags uint32) (*avahi.ServiceBrowser, error) {
	fmt.Printf("ServiceBrowserNew %s\n", serviceType)
	sb := new(avahi.ServiceBrowser)
	sb.AddChannel = make(chan avahi.Service)
	sb.RemoveChannel = make(chan avahi.Service)
	go s.scenarioPlayer(sb)
	return sb, nil
}

func (s *SimServer) ResolveService(iface, protocol int32, name, serviceType, domain string, aprotocol int32, flags uint32) (reply avahi.Service, err error) {
	svc := avahi.Service{
		Protocol: protocol,
		Name:     name,
		Type:     serviceType,
		Domain:   domain,
	}
	// use aprotocol to select the address...
	switch protocol {
	case 0:
		svc.Address = "192.168.0.1"
		svc.Port = 1000
	case 1:
		svc.Address = "192.168.0.2"
		svc.Port = 1001
	default:
		svc.Address = "192.168.0.99"
		svc.Port = 9999
	}

	return svc, nil
}

type cbRecordEntry struct {
	ts     time.Time
	added  bool
	svcInf ServiceInfo
}
type callbackRecord struct {
	records []cbRecordEntry
}

func addCbRecord(context interface{}, added bool, svcInf ServiceInfo) {
	context.(*callbackRecord).records = append(context.(*callbackRecord).records, cbRecordEntry{
		ts:     time.Now(),
		added:  added,
		svcInf: svcInf,
	})
}
func addCb(svcInf ServiceInfo, context interface{}) error {
	addCbRecord(context, true, svcInf)
	return nil
}

func removeCb(svcInf ServiceInfo, context interface{}) error {
	addCbRecord(context, false, svcInf)
	return nil
}

func TestServiceObserver(t *testing.T) {
	server = &SimServer{
		scenarioPlayer: scenario1,
	}
	r := &callbackRecord{}
	go ServiceObserver("_foo._tcp", r, addCb, removeCb)
	time.Sleep(time.Millisecond * 300)

	assert.Equal(t, len(r.records), 4)

	assert.Equal(t, r.records[0].added, true)
	assert.Equal(t, r.records[0].svcInf.service.Name, "fooA")
	assert.Equal(t, r.records[0].svcInf.service.Address, "192.168.0.1")

	assert.Equal(t, r.records[1].added, true)
	assert.Equal(t, r.records[1].svcInf.service.Name, "fooB")
	assert.Equal(t, r.records[1].svcInf.service.Address, "192.168.0.2")

	assert.Equal(t, r.records[2].added, false)
	assert.Equal(t, r.records[2].svcInf.service.Name, "fooB")
	assert.Equal(t, r.records[3].added, false)
	assert.Equal(t, r.records[3].svcInf.service.Name, "fooA")
}

func TestGetServiceInfo1(t *testing.T) {
	server = &SimServer{
		scenarioPlayer: scenario1,
	}
	start := time.Now()
	svcInf, err := GetServiceInfo("fooA", "_foo._tcp", time.Second)
	assert.Nil(t, err)
	assert.Equal(t, svcInf.service.Name, "fooA")
	assert.Equal(t, svcInf.service.Address, "192.168.0.1")
	assert.Less(t, time.Since(start), time.Millisecond*300)

	_, err = GetServiceInfo("barA", "_bar._tcp", time.Second)
	assert.NotNil(t, err)

	svcInf, err = GetServiceInfo("fooB", "_foo._tcp", time.Second)
	assert.Nil(t, err)
	assert.Equal(t, svcInf.service.Name, "fooB")
	assert.Equal(t, svcInf.service.Address, "192.168.0.2")
	assert.Less(t, time.Since(start), time.Millisecond*30)

}
