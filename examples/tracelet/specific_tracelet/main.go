package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/tracelet"
	"github.com/ci4rail/io4edge-client-go/transport"
	proto "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

var traceletMap = map[string]chan *proto.TraceletToServer{
	"TRACELET-1": make(chan *proto.TraceletToServer),
	"TRACELET-2": make(chan *proto.TraceletToServer),
	"TRACELET-3": make(chan *proto.TraceletToServer),
}

func catchTracelet1Stream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Printf("Received tracelet 1 location: %v", loc.Gnss)
	}
}

func catchTracelet2Stream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Printf("Received tracelet 2 location: %v", loc.Gnss)
	}
}

func catchTracelet3Stream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Printf("Received tracelet 3 location: %v", loc.Gnss)
	}
}

func traceletRegistration(register chan *tracelet.Channel) {
	for {
		ch := <-register
		log.Printf("Received new tracelet")
		// initiate first read to get tracelet_id from the first message
		msg, err := ch.ReadData()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		stream, ok := traceletMap[ch.TraceletID]
		// its one of the tracelets we are interested in
		if ok {
			stream <- msg
			ch.ReadStream(stream)
		} else {
			// serve other channels anyway
			go func() {
				for {
					// drop data of irrelevant tracelets
					_, err := ch.ReadData()
					if err == transport.ErrClosed {
						log.Printf("Channel closed")
						return
					}
				}
			}()
		}
	}
}

func main() {
	register := make(chan *tracelet.Channel)
	server := tracelet.NewTraceletServer("11001", time.Second*5)
	server.ManageConnections(register)

	go catchTracelet1Stream(traceletMap["TRACELET-1"])
	go catchTracelet2Stream(traceletMap["TRACELET-2"])
	go catchTracelet3Stream(traceletMap["TRACELET-3"])

	traceletRegistration(register)
}
