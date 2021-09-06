package basefunc

import (
	"errors"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge"
	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/transport"
)

// NewClientFromSocketAddress creates a new base function client from a socket with the specified address.
func NewClientFromSocketAddress(address string) (*Client, error) {
	t, err := transport.NewSocketConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ms, err := transport.NewMsgStreamFromConnection(t)
	if err != nil {
		return nil, errors.New("can't create msg stream: " + err.Error())
	}

	ch, err := io4edge.NewChannel(ms)
	if err != nil {
		return nil, errors.New("can't create channel: " + err.Error())
	}
	c, err := NewClient(ch)
	if err != nil {
		return nil, errors.New("can't basefunc create client: " + err.Error())
	}
	return c, nil
}
