package transport

// MsgStream is the interface used by a Channel to exchange message frames with the transport layer
// e.g. socket, websocket...
type MsgStream interface {
	ReadMsg() (payload []byte, err error)
	WriteMsg(payload []byte) (err error)
	Close() error
}
