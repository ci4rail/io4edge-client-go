package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/mvbsniffer"
)

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <mdns-service-address OR <ip:port>> <pattern>", os.Args[0])
	}
	address := os.Args[1]
	pattern := os.Args[2]

	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create mvbSniffer client: %v\n", err)
	}

	err = c.SendPattern(pattern)
	if err != nil {
		log.Errorf("SendPattern failed: %v\n", err)
	}
}
