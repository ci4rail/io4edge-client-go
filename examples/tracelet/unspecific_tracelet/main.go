package main

import (
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/tracelet"
	proto "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

func catchTraceletStream(stream chan *proto.TraceletToServer) {
	for {
		msg := <-stream
		loc := msg.GetLocation()
		log.Printf("Received location of %s: %v", msg.TraceletId, loc.Gnss)
	}
}

func traceletRegistration(register chan *tracelet.TraceletChannel, stream chan *proto.TraceletToServer) {
	for {
		ch := <-register
		log.Printf("Received new tracelet")

		ch.ReadStream(stream)
	}
}

func main() {
	register := make(chan *tracelet.TraceletChannel)
	stream := make(chan *proto.TraceletToServer)
	server := tracelet.NewTraceletServer("11002", time.Second*5)
	server.ManageConnections(register)

	go catchTraceletStream(stream)
	traceletRegistration(register, stream)
}
