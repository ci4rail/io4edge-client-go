package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/mvbsniffer"
)

func generatePattern() string {
	cl := mvbsniffer.NewCommandList()

	// req 32 bits process data from address 123
	cl.AddMasterFrame(0, false, 5, 0x1, 567)
	cl.AddMasterFrame(0, false, 3, 0x1, 888)
	cl.AddMasterFrame(0, false, 5, 0x1, 123)

	cl.AddSlaveFrame(0, false, 100, []uint8{0xaa, 0xbb, 0xcc, 0xdd})

	return cl.StartGeneratorString(true)
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR <ip:port>>", os.Args[0])
	}
	address := os.Args[1]

	pattern := generatePattern()

	fmt.Printf("Generator pattern: %s\n", pattern)
	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create mvbSniffer client: %v\n", err)
	}

	err = c.SendPattern(pattern)
	if err != nil {
		log.Errorf("SendPattern failed: %v\n", err)
	}
}
