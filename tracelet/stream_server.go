package tracelet

import (
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/server"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

// Server represents a tracelet server
type Server struct {
	srv           *server.UDPServer
	timeout       time.Duration
	subscriptions map[string]chan *tracelet.TraceletToServer
}

// Channel represents a connection to a tracelet
type Channel struct {
	ch      *client.Channel
	timeout time.Duration
}

// NewTraceletServer creates a new tracelet server
func NewTraceletServer(port string, timeout time.Duration) *Server {
	// create server
	addr := ":" + port
	srv, err := server.NewServer(addr)
	if err != nil {
		log.Fatal("can't create server: " + err.Error())
		return nil
	}
	traceletServer := &Server{
		srv:           srv,
		timeout:       timeout,
		subscriptions: make(map[string]chan *tracelet.TraceletToServer),
	}
	return traceletServer
}

// ListenForConnections listens for connections to the tracelet server and starts reading stream
func (s *Server) ListenForConnections() {
	for {
		ch, err := s.srv.ListenForNextConnection()
		if err != nil {
			log.Errorf("Error reading message: %v", err)
			continue
		}

		tracelet := &Channel{
			ch:      ch,
			timeout: s.timeout,
		}
		log.Infof("Received new tracelet")
		tracelet.ReadStream(&s.subscriptions)

	}
}

// Close closes the tracelet server
func (s *Server) Close() {
	s.srv.Close()
}

// Subscribe subscribes to a tracelet pattern and returns a stream channel. All messages of
// tracelets whose ids matches the pattern will be sent to this channel.
func (s *Server) Subscribe(traceletPattern string) chan *tracelet.TraceletToServer {
	stream := make(chan *tracelet.TraceletToServer)
	s.subscriptions[traceletPattern] = stream
	return stream
}

// Unsubscribe unsubscribes from a tracelet pattern
func (s *Server) Unsubscribe(traceletPattern string) {
	delete(s.subscriptions, traceletPattern)
}

// ReadData reads a single message from the tracelet channel
func (t *Channel) readData() (*tracelet.TraceletToServer, error) {
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

	return msg, nil
}

// ReadStream reads messages from the tracelet channel and sends them to the provided channel
func (t *Channel) ReadStream(subscriptions *map[string]chan *tracelet.TraceletToServer) {
	go func() {
		for {
			msg, err := t.readData()
			if err == transport.ErrTimeout || err == transport.ErrClosed {
				return
			} else if err != nil {
				continue
			}
			// go through all subscriptions and search for matches
			for pattern, stream := range *subscriptions {
				// check if pattern matches tracelet id
				match, err := regexp.MatchString(pattern, msg.TraceletId)
				if err != nil {
					log.Errorf("Error matching pattern: %v", err)
				} else if match {
					stream <- msg
				}
			}
		}
	}()
}

// Close closes the tracelet channel
func (t *Channel) Close() {
	log.Printf("Closing connection")
	t.ch.Close()
}
