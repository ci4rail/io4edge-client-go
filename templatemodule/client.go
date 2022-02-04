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
	SampleRate uint32
}

// Description represents the describe response of the templateModule
type Description struct {
	Ident string
}

// NewClient creates a new new templateModule client from a functionBlock client
func NewClient(fbClient *functionblock.Client) *Client {
	return &Client{
		fbClient: fbClient,
	}
}

// ConfigurationSet configures the template function block
func (c *Client) ConfigurationSet(config *Configuration) error {
	fsCmd := &fspb.ConfigurationSet{
		SampleRate: config.SampleRate,
	}
	_, err := c.fbClient.ConfigurationSet(fsCmd)
	return err
}

// ConfigurationGet reads the template function block configuration
func (c *Client) ConfigurationGet() (*Configuration, error) {
	fsCmd := &fspb.ConfigurationGet{}
	any, err := c.fbClient.ConfigurationGet(fsCmd)
	if err != nil {
		return nil, err
	}
	res := new(fspb.ConfigurationGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return nil, err
	}
	cfg := &Configuration{
		SampleRate: res.SampleRate,
	}
	return cfg, err
}

// ConfigurationDescribe reads the template function block description
func (c *Client) ConfigurationDescribe() (*Description, error) {
	fsCmd := &fspb.ConfigurationDescribe{}
	any, err := c.fbClient.ConfigurationDescribe(fsCmd)
	if err != nil {
		return nil, err
	}
	res := new(fspb.ConfigurationDescribeResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return nil, err
	}
	desc := &Description{
		Ident: res.Ident,
	}
	return desc, err
}

// SetCounter sets the templates module counter in the device
func (c *Client) SetCounter(value uint32) error {
	fsCmd := &fspb.FunctionControlSet{
		Value: value,
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// GetCounter reads the templates module counter from the device
func (c *Client) GetCounter() (uint32, error) {
	fsCmd := &fspb.FunctionControlGet{}
	any, err := c.fbClient.FunctionControlGet(fsCmd)
	if err != nil {
		return 0, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return 0, err
	}
	return res.Value, nil
}
