package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/basefunc"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: identify <device-address>\n")
	}
	address := os.Args[1]

	c, err := basefunc.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	fwID, err := c.IdentifyFirmware(5 * time.Second)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %d.%d.%d\n", fwID.Name, fwID.MajorVersion, fwID.MinorVersion, fwID.PatchVersion)

	hwID, err := c.IdentifyHardware(5 * time.Second)
	if err != nil {
		log.Fatalf("Failed to identify hardware: %v\n", err)
	}
	fmt.Printf("Hardware name: %s, serial: %16x-%16x, rev: %d\n", hwID.RootArticle, hwID.SerialNumber.Hi, hwID.SerialNumber.Lo, hwID.MajorVersion)
}
