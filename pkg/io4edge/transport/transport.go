package transport

// Transport is the interface used by Channel to communicate with the underlying stream
// e.g. socket, websocket...
type Transport interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
}
