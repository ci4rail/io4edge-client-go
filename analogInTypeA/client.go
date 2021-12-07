package analogInTypeA

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
)

// Client represents a client for the analogInTypeA function
type Client struct {
	funcClient *client.Client
}

func NewClient(c *client.Client) *Client {
	return &Client{
		funcClient: c,
	}
}

// NewClientFromSocketAddress creates a new analogInTypeA function client from a socket with the specified address.
func NewClientFromSocketAddress(address string) (*Client, error) {
	c, err := client.NewClientFromSocketAddress(address)
	if err != nil {
		return nil, errors.New("can't create function client: " + err.Error())
	}
	functionClient := NewClient(c)

	return functionClient, nil
}

// NewClientFromService creates a new core function client from a socket with a address, which was acquired from the specified service.
// The timeout specifies the maximal time waiting for a service to show up.
func NewClientFromService(serviceAddr string, timeout time.Duration) (*Client, error) {
	c, err := client.NewClientFromService(serviceAddr, timeout)

	if err != nil {
		return nil, err
	}
	functionClient := NewClient(c)

	return functionClient, nil
}
