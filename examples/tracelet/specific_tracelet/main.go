package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/tracelet"
	proto "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

func catchTracelet1Stream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Infof("Received tracelet 1 location: %v", loc.Gnss)
	}
}

func catchTracelet2Stream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Infof("Received tracelet 2 location: %v", loc.Gnss)
	}
}

func catchTracelet3Stream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Infof("Received tracelet 3 location: %v", loc.Gnss)
	}
}

func main() {
	log.SetLevel(log.InfoLevel)
	server := tracelet.NewTraceletUDPServer(11002, time.Second*5)

	go catchTracelet1Stream(server.Subscribe("TRACELET-1"))
	go catchTracelet2Stream(server.Subscribe("TRACELET-2"))
	go catchTracelet3Stream(server.Subscribe("TRACELET-3"))
	server.ListenForConnections()
}
