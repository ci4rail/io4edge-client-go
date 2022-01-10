package binaryIoTypeA

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	fbv1 "github.com/ci4rail/io4edge-client-go/functionblock/v1alpha1"
)

type ClientInterface interface {
	NewClientFromSocketAddress(address string) (*Client, error)
	NewClientFromService(serviceAddr string, timeout time.Duration) (*Client, error)
}

// Client represents a client for the binaryIoTypeA function
type Client struct {
	funcClient           *client.Client
	streamClientChannels map[int]chan bool
	responses            sync.Map
	streamData           chan *fbv1.StreamData
}

func NewClient(c *client.Client) *Client {
	client := &Client{
		funcClient:           c,
		streamClientChannels: make(map[int]chan bool),
		responses:            sync.Map{},
		streamData:           make(chan *fbv1.StreamData, 100),
	}
	client.ReadResponse()
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
	fmt.Println("sending command with context: ", cmd.Context)
	err := c.funcClient.Ch.WriteMessage(cmd)
	if err != nil {
		return nil, err
	}
	res := &fbv1.Response{}
	res, err = c.WaitForResponse(cmd.Context.Value, timeout)
	if err != nil {
		return nil, err
	}
	return res, nil
}
