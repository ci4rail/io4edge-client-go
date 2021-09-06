package io4edge

import (
	"time"

	"google.golang.org/protobuf/proto"
)

// Client represents a client for the io4edge base function
type Client struct {
	ch *Channel
}

// NewClient creates a new client for the base function
func NewClient(c *Channel) (*Client, error) {
	return &Client{ch: c}, nil
}

// Command issues a command cmd to a channel, waits for the devices response and returns it in res
func (c *Channel) Command(cmd proto.Message, res proto.Message, timeout time.Duration) error {
	err := c.WriteMessage(cmd)
	if err != nil {
		return err
	}
	err = c.ReadMessage(res, timeout)
	if err != nil {
		return err
	}
	return err
}
