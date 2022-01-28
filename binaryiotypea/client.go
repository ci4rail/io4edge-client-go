package binaryiotypea

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// ClientInterface provides the interface to establish the socket
type ClientInterface interface {
	NewClientFromSocketAddress(address string) (*Client, error)
	NewClientFromService(serviceAddr string, timeout time.Duration) (*Client, error)
}

// Client represents a client for the binaryIoTypeA function
type Client struct {
	funcClient              *client.Client
	streamClientStopChannel chan bool
	readResponsesStopChan   chan bool
	responses               sync.Map
	streamData              chan *fbv1.StreamData
	responsePending         int
	streamRunning           bool
	streamStatus            bool
	customRecover           func()
	connected               bool
	streamKeepaliveInterval uint32
}

// NewClient creates a new client and waits for responses
func NewClient(c *client.Client) *Client {
	client := &Client{
		funcClient:              c,
		streamClientStopChannel: make(chan bool),
		readResponsesStopChan:   make(chan bool),
		responses:               sync.Map{},
		streamData:              make(chan *fbv1.StreamData, 200),
		responsePending:         0,
		streamRunning:           false,
		streamStatus:            false,
		customRecover:           nil,
		connected:               true,
		streamKeepaliveInterval: DefaultKeepaliveInterval,
	}
	client.ReadResponses()
	return client
}

// NewClientFromSocketAddress creates a new binaryIoTypeA function client from a socket with the specified address.
func NewClientFromSocketAddress(address string) (*Client, error) {
	c, err := client.NewClientFromSocketAddress(address)
	if err != nil {
		return nil, errors.New("can't create function client: " + err.Error())
	}
	binClient := NewClient(c)

	return binClient, nil
}

// NewClientFromService creates a new core function client from a socket with a address, which was acquired from the specified service.
// The timeout specifies the maximal time waiting for a service to show up.
func NewClientFromService(serviceAddr string, timeout time.Duration) (*Client, error) {
	c, err := client.NewClientFromService(serviceAddr, timeout)

	if err != nil {
		return nil, err
	}
	binClient := NewClient(c)

	return binClient, nil
}

// Command issues a command cmd to a channel, waits for the devices response and returns it in res
func (c *Client) Command(cmd *fbv1.Command, timeout time.Duration) (*fbv1.Response, error) {
	log.Debug("sending command with context: ", cmd.Context)
	c.responsePending++
	err := c.funcClient.Ch.WriteMessage(cmd)
	if err != nil {
		return nil, err
	}
	res, err := c.WaitForResponse(cmd.Context.Value, timeout)
	if err != nil {
		return nil, err
	}
	c.responsePending--
	return res, nil
}

// SetRecover sets a function to inform the client about a panic event like timeout.
func (c *Client) SetRecover(customRecover func()) {
	c.customRecover = customRecover
}

func (c *Client) recover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered in binaryIoTypeA Client", r)
	}
	c.streamClientStopChannel <- true
	c.connected = false
	c.streamStatus = false
	if c.customRecover != nil {
		c.customRecover()
	}
}
