package functionblock

import (
	"sync"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
)

// Client represents a client for a generic functionblock
type Client struct {
	funcClient     *client.Client
	customRecover  func()
	commandTimeout int
	//streamKeepaliveInterval uint32
	streamChan chan *fbv1.StreamData

	cmdSeqNo uint32
	cmdMutex sync.Mutex
	// command/response handshake
	waitingCmdSeqChan chan uint32 // waiting sequence number
	responseChan      chan *fbv1.Response
}

func newClient(c *client.Client) *Client {
	client := &Client{
		funcClient:     c,
		customRecover:  nil,
		commandTimeout: 5,
		streamChan:     make(chan *fbv1.StreamData, 100),
		//streamKeepaliveInterval: 10,
		waitingCmdSeqChan: make(chan uint32),
		responseChan:      make(chan *fbv1.Response),
	}
	//log.SetLevel(log.DebugLevel)
	client.readResponses()
	return client
}

// NewClientFromUniversalAddress creates a new functionblock client from addrOrService.
// If addrOrService is of the form "host:port", it creates the client from that host/port,
// otherwise it assumes addrOrService is a mnds service name.
// The timeout specifies the maximal time waiting for a service to show up. Not used for "host:port"
func NewClientFromUniversalAddress(addrOrService string, timeout time.Duration) (*Client, error) {
	io4eClient, err := client.NewClientFromUniversalAddress(addrOrService, timeout)

	if err != nil {
		return nil, err
	}
	c := newClient(io4eClient)

	return c, nil
}

// SetRecover sets a function to inform the client about a panic event like timeout.
func (c *Client) SetRecover(customRecover func()) {
	c.customRecover = customRecover
}

// func (c *Client) recover() {
// 	if r := recover(); r != nil {
// 		fmt.Println("Recovered in binaryIoTypeA Client", r)
// 	}

// 	if c.customRecover != nil {
// 		c.customRecover()
// 	}
// }-uu80ö---.
