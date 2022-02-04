package templatemodule

import (
	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/templateModule/go/templateModule/v1alpha1"
)

type Client struct {
	fbClient *functionblock.Client
}

type Configuration struct {
	SomeValue int
}

func NewClient(fbClient *functionblock.Client) *Client {
	return &Client{
		fbClient: fbClient,
	}
}

func (c *Client) ConfigurationSet(config *Configuration) error {
	fsCmd := &fspb.ConfigurationSet{}
	_, err := c.fbClient.ConfigurationSet(fsCmd)
	return err
}
