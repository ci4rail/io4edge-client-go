package tracelet

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/server"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

// Server represents a tracelet server
type Server struct {
	srv     *server.UDPServer
	timeout time.Duration
}

// Channel represents a connection to a tracelet
type Channel struct {
	ch         *client.Channel
	TraceletID string
	timeout    time.Duration
}

// NewTraceletServer creates a new tracelet server
func NewTraceletServer(port string, defaultTimeout time.Duration) *Server {
	// create server
	addr := ":" + port
	srv, err := server.NewServer(addr)
	if err != nil {
		log.Fatal("can't create server: " + err.Error())
		return nil
	}
	traceletServer := &Server{
		srv:     srv,
		timeout: defaultTimeout,
	}
	return traceletServer
}

// ManageConnections manages the connections to the tracelet server
func (s *Server) ManageConnections(c chan *Channel) {
	go func() {
		for {
			ch, err := s.srv.ManageConnections()
			if err != nil {
				log.Errorf("Error reading message: %v", err)
				continue
			}

			tracelet := &Channel{
				ch:         ch,
				TraceletID: "",
				timeout:    s.timeout,
			}
			log.Infof("Received new tracelet")
			c <- tracelet
			log.Printf("Sent new tracelet")
		}
	}()
}

// Close closes the tracelet server
func (s *Server) Close() {
	s.srv.Close()
}

// WriteData writes a single message to the tracelet channel
func (t *Channel) WriteData(msg *tracelet.ServerToTracelet) error {
	err := t.ch.WriteMessage(msg)
	if err != nil {
		log.Fatal("can't write to connection: " + err.Error())
		return err
	}
	return nil
}

// ReadData reads a single message from the tracelet channel
func (t *Channel) ReadData() (*tracelet.TraceletToServer, error) {
	msg := &tracelet.TraceletToServer{}
	err := t.ch.ReadMessage(msg, t.timeout)
	if err == transport.ErrTimeout {
		log.Printf("Read timeout -> close socket")
		t.Close()
		return nil, err
	} else if err != nil {
		log.Fatal("can't read from connection: " + err.Error())
		return nil, err
	}

	if t.TraceletID == "" {
		// store tracelet_id on first read
		t.TraceletID = msg.TraceletId
		log.Printf("Tracelet ID: %s", t.TraceletID)
	}
	return msg, nil
}

// ReadStream reads messages from the tracelet channel and sends them to the provided channel
func (t *Channel) ReadStream(stream chan *tracelet.TraceletToServer) {
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

// Close closes the tracelet channel
func (t *Channel) Close() {
	log.Printf("Closing connection")
	t.ch.Close()
}
