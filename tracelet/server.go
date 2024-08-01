package tracelet

import (
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/server"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

type TraceletServer struct {
	srv     *server.UDPServer
	timeout time.Duration
}

type TraceletChannel struct {
	ch          *client.Channel
	Tracelet_id string
	timeout     time.Duration
}

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

func (s *TraceletServer) ManageConnections(c chan *TraceletChannel) {
	go func() {
		for {
			ch, err := s.srv.ManageConnections()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// todo mutex?
			tracelet := &TraceletChannel{
				// srv:         s,
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

func (t *TraceletChannel) WriteData(msg *tracelet.ServerToTracelet) error {
	err := t.ch.WriteMessage(msg)
	if err != nil {
		log.Fatal("can't write to connection: " + err.Error())
		return err
	}
	return nil
}

func (t *TraceletChannel) TestTimeout() {
	msg := &tracelet.TraceletToServer{}
	for {
		log.Printf("empty buffer")
		err := t.ch.ReadMessage(msg, time.Second*2)
		if err == transport.ErrTimeout {
			log.Printf("connection timed out")
			break
		}
	}
	time.Sleep(2 * time.Second)
	err := t.ch.ReadMessage(msg, time.Second*3)
	if err == transport.ErrTimeout {
		log.Printf("connection timed out")
	}
	err = t.ch.ReadMessage(msg, 0)
	if err == transport.ErrTimeout {
		log.Printf("connection timed out")
	}
	log.Printf("TestTimeout done")
}

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

	// update tracelet id on read
	t.Tracelet_id = msg.TraceletId
	// if t.Tracelet_id == "" {
	// 	// store tracelet_id on first read
	// 	t.Tracelet_id = msg.TraceletId
	// 	log.Printf("Tracelet ID: %s", t.Tracelet_id)
	// }
	return msg, nil
}

func (t *TraceletChannel) ReadStream(stream chan *tracelet.TraceletToServer_Location) {
	go func() {
		for {
			msg, err := t.ReadData()
			if err == transport.ErrTimeout || err == transport.ErrClosed {
				return
			} else if err != nil {
				continue
			}

			loc := msg.GetLocation()
			stream <- loc
		}
	}()
}

func (t *TraceletChannel) Close() {
	log.Printf("Closing connection")
	t.ch.Close()
}
