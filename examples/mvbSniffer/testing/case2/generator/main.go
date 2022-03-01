package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/mvbsniffer"
)

func errChk(err error) {
	if err != nil {
		panic(err)
	}
}

func generatePattern() string {
	cl := mvbsniffer.NewCommandList()

	errChk(cl.AddMasterFrame(7500, false, 43, 3, 11))

	errChk(cl.AddMasterFrame(7500, false, 43, 9, 528))
	errChk(cl.AddMasterFrame(7500, false, 43, 9, 528))
	errChk(cl.AddMasterFrame(7500, false, 800, 9, 272))

	errChk(cl.AddMasterFrame(7500, false, 3, 0, 2))
	errChk(cl.AddSlaveFrame(7500, false, 3, []uint8{0x00, 0x00}))

	errChk(cl.AddMasterFrame(7500, false, 43, 9, 528))
	errChk(cl.AddMasterFrame(7500, false, 43, 9, 528))
	errChk(cl.AddMasterFrame(7500, false, 800, 9, 272))

	errChk(cl.AddMasterFrame(7500, false, 3, 0, 3))
	errChk(cl.AddSlaveFrame(7500, true, 3, []uint8{0x00, 0x00}))

	errChk(cl.AddMasterFrame(7500, false, 43, 9, 528))
	errChk(cl.AddMasterFrame(7500, false, 43, 9, 528))
	errChk(cl.AddMasterFrame(7500, false, 800, 9, 272))

	errChk(cl.AddMasterFrame(7500, false, 43, 4, 6))

	return cl.StartGeneratorString(false) // external loop
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR <ip:port>>", os.Args[0])
	}
	address := os.Args[1]

	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create mvbsniffer client: %v\n", err)
	}

	// ensure pattern is stopped
	errChk(c.SendPattern(c.StopGeneratorString()))

	pattern := generatePattern()
	fmt.Printf("Generator pattern: %s\n", pattern)
	errChk(c.SendPattern(pattern))
}
