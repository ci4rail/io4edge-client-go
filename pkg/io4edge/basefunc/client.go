package basefunc

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge"
)

// Client represents a client for the io4edge base function
type Client struct {
	ch *io4edge.Channel
}

// NewClient creates a new client for the base function
func NewClient(c *io4edge.Channel) (*Client, error) {
	return &Client{ch: c}, nil
}

// Command issues a command cmd to a base function channel, waits for the devices response and returns it in res
func (c *Client) Command(cmd *BaseFuncCommand, res *BaseFuncResponse, timeout time.Duration) error {
	err := c.ch.WriteMessage(cmd)
	if err != nil {
		return err
	}
	err = c.ch.ReadMessage(res, timeout)
	if err != nil {
		return errors.New("Failed to receive device response: " + err.Error())
	}
	if res.Status != BaseFuncStatus_OK {
		return errors.New("Device reported error status: " + res.Status.String())
	}
	return err
}
