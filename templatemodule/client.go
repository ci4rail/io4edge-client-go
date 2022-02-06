package templatemodule

import (
	"errors"
	"time"

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

// StreamData contains the meta data of the stream and the unmarshalled function specific data
type StreamData struct {
	functionblock.StreamDataMeta
	FSData *fspb.StreamData
}

// NewClientFromUniversalAddress creates a new templateModule client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is a mnds service name.
// The timeout specifies the maximal time waiting for a service to show up. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, timeout)

	if err != nil {
		return nil, err
	}
	return &Client{
		fbClient: io4eClient,
	}, nil
}

// UploadConfiguration configures the template function block
func (c *Client) UploadConfiguration(config *Configuration) error {
	fsCmd := &fspb.ConfigurationSet{
		SampleRate: config.SampleRate,
	}
	_, err := c.fbClient.UploadConfiguration(fsCmd)
	return err
}

// DownloadConfiguration reads the template function block configuration
func (c *Client) DownloadConfiguration() (*Configuration, error) {
	fsCmd := &fspb.ConfigurationGet{}
	any, err := c.fbClient.DownloadConfiguration(fsCmd)
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

// Describe reads the template function block description
func (c *Client) Describe() (*Description, error) {
	fsCmd := &fspb.ConfigurationDescribe{}
	any, err := c.fbClient.Describe(fsCmd)
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

// StartStream starts the stream on this connection
func (c *Client) StartStream(genericConfig *functionblock.StreamConfiguration, increment uint32) error {
	fsCmd := &fspb.StreamControlStart{
		SampleIncrement: increment,
	}
	err := c.fbClient.StartStream(genericConfig, fsCmd)
	if err != nil {
		return err
	}
	return nil
}

// StopStream stops the stream on this connection
func (c *Client) StopStream() error {
	return c.fbClient.StopStream()
}

// ReadStream reads the next stream data object from the buffer
// returns the meta data and the unmarshalled function specific stream data
func (c *Client) ReadStream(timeout time.Duration) (*StreamData, error) {
	genericSD, err := c.fbClient.ReadStream(timeout)
	if err != nil {
		return nil, err
	}

	fsSD := new(fspb.StreamData)
	if err := genericSD.FSData.UnmarshalTo(fsSD); err != nil {
		return nil, errors.New("can't unmarshall samples")
	}

	sd := &StreamData{
		StreamDataMeta: genericSD.StreamDataMeta,
		FSData:         fsSD,
	}
	return sd, nil
}
