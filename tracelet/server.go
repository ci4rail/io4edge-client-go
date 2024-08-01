package tracelet

import (
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/server"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

// TraceletServer represents a tracelet server
type TraceletServer struct {
	srv     *server.UDPServer
	timeout time.Duration
}

// TraceletChannel represents a connection to a tracelet
type TraceletChannel struct {
	ch          *client.Channel
	Tracelet_id string
	timeout     time.Duration
}

// NewTraceletServer creates a new tracelet server
func NewTraceletServer(port string, defaultTimeout time.Duration) *TraceletServer {
	// create server
	addr := ":" + port
	srv, err := server.NewServer(addr)
	if err != nil {
		log.Fatal("can't create server: " + err.Error())
		return nil
	}
	traceletServer := &TraceletServer{
		srv:     srv,
		timeout: defaultTimeout,
	}
	return traceletServer
}

// ManageConnections manages the connections to the tracelet server
func (s *TraceletServer) ManageConnections(c chan *TraceletChannel) {
	go func() {
		for {
			ch, err := s.srv.ManageConnections()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			tracelet := &TraceletChannel{
				ch:          ch,
				Tracelet_id: "",
				timeout:     s.timeout,
			}
			log.Printf("Received new tracelet")
			c <- tracelet
			log.Printf("Sent new tracelet")
		}
	}()
}

// WriteData writes a single message to the tracelet channel
func (t *TraceletChannel) WriteData(msg *tracelet.ServerToTracelet) error {
	err := t.ch.WriteMessage(msg)
	if err != nil {
		log.Fatal("can't write to connection: " + err.Error())
		return err
	}
	return nil
}

// ReadData reads a single message from the tracelet channel
func (t *TraceletChannel) ReadData() (*tracelet.TraceletToServer, error) {
	msg := &tracelet.TraceletToServer{}
	err := t.ch.ReadMessage(msg, t.timeout)
	if err == transport.ErrTimeout {
		log.Printf("connection timed out -> close socket")
		t.ch.Close()
		return nil, err
	} else if err != nil {
		log.Fatal("can't read from connection: " + err.Error())
		return nil, err
	}

	if t.Tracelet_id == "" {
		// store tracelet_id on first read
		t.Tracelet_id = msg.TraceletId
		log.Printf("Tracelet ID: %s", t.Tracelet_id)
	}
	return msg, nil
}

// ReadStream reads messages from the tracelet channel and sends them to the provided channel
func (t *TraceletChannel) ReadStream(stream chan *tracelet.TraceletToServer) {
	go func() {
		for {
			msg, err := t.ReadData()
			if err == transport.ErrTimeout || err == transport.ErrClosed {
				return
			} else if err != nil {
				continue
			}
			stream <- msg
		}
	}()
}

func (t *TraceletChannel) Close() {
	log.Printf("Closing connection")
	t.ch.Close()
}
