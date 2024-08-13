package main

import (
	"time"

	log "github.com/sirupsen/logrus"

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

func main() {
	log.SetLevel(log.InfoLevel)
	server := tracelet.NewTraceletServer("11002", time.Second*5)
	go server.ListenForConnections()

	catchTraceletStream(server.Subscribe(".*"))
}
