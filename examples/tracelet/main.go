package main

import (
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/tracelet"
	protoTrace "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

func main() {
	clients := make(map[string]*tracelet.TraceletChannel)
	clients_known := make(map[string]*tracelet.TraceletChannel)
	register := make(chan *tracelet.TraceletChannel)
	location := make(chan *protoTrace.TraceletToServer_Location)
	server := tracelet.NewTraceletServer("11002", time.Second*5)
	server.ManageConnections(register)

	// for {
	// 	ch := <-register
	// 	log.Printf("Received new tracelet")
	// 	// initiate first read to get tracelet_id from the first message
	// 	msg, err := ch.ReadData()
	// 	if err != nil {
	// 		log.Printf("Error reading message: %v", err)
	// 		continue
	// 	} else if msg.GetType().(*protoTrace.TraceletToServer_Location_) != nil {
	// 		loc := msg.GetLocation()
	// 		log.Printf("Received location: %v", loc.Gnss)
	// 	}

	// 	log.Printf("Start test timeout1")
	// 	time.Sleep(10 * time.Second)
	// 	log.Printf("Start test timeout2")
	// 	ch.TestTimeout()
	// }

	go func() {
		for {
			ch := <-register
			log.Printf("Received new tracelet")
			// initiate first read to get tracelet_id from the first message
			ch.ReadData()
			if ch.Tracelet_id == "TRACELET-1" {
				clients_known[ch.Tracelet_id] = ch
				ch.ReadStream(location)
			} else {
				clients[ch.Tracelet_id] = ch
			}
		}
	}()

	go func() {
		for {
			loc := <-location
			log.Printf("Received location: %+v", loc.Gnss)
		}
	}()

	for {
		for _, ch := range clients {
			msg, err := ch.ReadData()
			if err != nil {
				log.Printf("Error reading message: %v", err)
			} else {
				if msg.GetType().(*protoTrace.TraceletToServer_Location_) != nil {
					loc := msg.GetLocation()
					log.Printf("Received location: %v", loc)
				}
				if msg.GetType().(*protoTrace.TraceletToServer_Status) != nil {
					sta := msg.GetStatus()
					log.Printf("Received status: %v", sta)
				}
			}
		}
	}
}
