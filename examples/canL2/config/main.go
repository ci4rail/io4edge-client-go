/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"log"

	"github.com/ci4rail/io4edge-client-go/canl2"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	bitratePtr := flag.Uint("bitrate", 500000, "CAN Bitrate")
	samplePointPtr := flag.Float64("samplepoint", 0.8, "CAN Sample Point (0.0-1.0)")
	sjwPtr := flag.Uint("sjw", 1, "CAN Synchronization Jump Width")
	listenOnlyPtr := flag.Bool("listenonly", false, "Listen only mode")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	address := flag.Arg(0)

	// Create a client object to work with the io4edge device
	c, err := canl2.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create canl2 client: %v\n", err)
	}

	// configure CAN
	err = c.UploadConfiguration(
		canl2.WithBitRate(uint32(*bitratePtr)),
		canl2.WithSamplePoint(float32(*samplePointPtr)),
		canl2.WithSJW(uint8(*sjwPtr)),
		canl2.WithListenOnly(*listenOnlyPtr),
	)
	if err != nil {
		log.Fatalf("Failed to configure canl2: %v\n", err)
	}

	fmt.Printf("Can configured!\n")
}
