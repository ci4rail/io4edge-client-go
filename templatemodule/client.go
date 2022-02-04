package templatemodule

import (
	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/templateModule/go/templateModule/v1alpha1"
)

// Client represents a client for the templateModule
type Client struct {
	fbClient *functionblock.Client
}

// Configuration represents the configuration parameters of the templateModule
type Configuration struct {
	SomeValue int
}

// NewClient creates a new new templateModule client from a functionBlock client
func NewClient(fbClient *functionblock.Client) *Client {
	return &Client{
		fbClient: fbClient,
	}
}

// ConfigurationSet configures the template function block
func (c *Client) ConfigurationSet(config *Configuration) error {
	fsCmd := &fspb.ConfigurationSet{}
	_, err := c.fbClient.ConfigurationSet(fsCmd)
	return err
}
