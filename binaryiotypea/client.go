package binaryiotypea

import (
	"errors"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
)

// Client represents a client for the templateModule
type Client struct {
	fbClient *functionblock.Client
}

// Configuration represents the configuration parameters of the templateModule
type Configuration struct {
	OutputFrittingMask    uint8
	OutputWatchdogMask    uint8
	OutputWatchdogTimeout uint32
}

// Description represents the describe response of the templateModule
type Description struct {
	NumberOfChannels int
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
	io4eClient, err := functionblock.NewClientFromUniversalAddress(addrOrService, "_io4edge_binaryIoTypeA._tcp", timeout)

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
		OutputFrittingMask:    uint32(config.OutputFrittingMask),
		OutputWatchdogMask:    uint32(config.OutputWatchdogMask),
		OutputWatchdogTimeout: config.OutputWatchdogTimeout,
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
		OutputFrittingMask:    uint8(res.OutputFrittingMask),
		OutputWatchdogMask:    uint8(res.OutputWatchdogMask),
		OutputWatchdogTimeout: res.OutputWatchdogTimeout,
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
		NumberOfChannels: int(res.NumberOfChannels),
	}
	return desc, err
}

// SetOutput sets a single output channel
// a "true" state turns on the output switch
func (c *Client) SetOutput(channel int, state bool) error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_Single{
			Single: &fspb.SetSingle{
				Channel: uint32(channel),
				State:   state,
			},
		},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// ExitErrorState tries to recover the binary output controller from error state
// The binary output controller enters error state when there is an overurrent condition for a long time
// In the error state, no outputs can be set and no inputs can be read.
// This call tells the binary output controller to try again. This call does however not wait
// if the recovery was successful or not.
func (c *Client) ExitErrorState(channel int, state bool) error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_ExitError{},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// SetAllOutputs sets all or a group of output channels
// states: binary coded map of outputs. 0 means switch off, 1 means switch on, LSB is Output0
// mask: defines which channels are affected by the set all command.
func (c *Client) SetAllOutputs(states uint8, mask uint8) error {
	fsCmd := &fspb.FunctionControlSet{

		Type: &fspb.FunctionControlSet_All{
			All: &fspb.SetAll{
				Values: uint32(states),
				Mask:   uint32(mask),
			},
		},
	}
	_, err := c.fbClient.FunctionControlSet(fsCmd)
	return err
}

// Input reads the state of the input pin of a single channel
// the returned value is false if the input level is below switching threshold, true otherwise
func (c *Client) Input(channel int) (bool, error) {
	fsCmd := &fspb.FunctionControlGet{
		Type: &fspb.FunctionControlGet_Single{
			Single: &fspb.GetSingle{
				Channel: uint32(channel),
			},
		},
	}
	any, err := c.fbClient.FunctionControlGet(fsCmd)
	if err != nil {
		return false, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return false, err
	}
	return res.GetSingle().State, nil
}

// AllInputs reads the state of all input pins defined by mask.
// each bit in the returned value corresponds to one channel, bit0 being channel 0.
// The bit is false if the input level is below switching threshold, true otherwise.
// Channels whose bit is cleared in mask are reported as 0
func (c *Client) AllInputs(mask uint8) (uint8, error) {
	fsCmd := &fspb.FunctionControlGet{
		Type: &fspb.FunctionControlGet_All{
			All: &fspb.GetAll{
				Mask: uint32(mask),
			},
		},
	}
	any, err := c.fbClient.FunctionControlGet(fsCmd)
	if err != nil {
		return 0, err
	}
	res := new(fspb.FunctionControlGetResponse)
	if err := any.UnmarshalTo(res); err != nil {
		return 0, err
	}
	return uint8(res.GetAll().Inputs), nil
}

// StartStream starts the stream on this connection
// channelFilterMask defines the watched channels. Only changes on those channels generate samples in the stream
func (c *Client) StartStream(genericConfig *functionblock.StreamConfiguration, channelFilterMask uint8) error {
	fsCmd := &fspb.StreamControlStart{
		ChannelFilterMask: uint32(channelFilterMask),
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
