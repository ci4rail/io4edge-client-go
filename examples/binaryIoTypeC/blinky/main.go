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
	"time"

	"log"

	"github.com/ci4rail/io4edge-client-go/binaryiotypec"
	biniopb "github.com/ci4rail/io4edge_api/binaryIoTypeC/go/binaryIoTypeC/v1alpha1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	channels := flag.Uint("channels", 4, "Number of channels to use")
	flag.Parse()

	numberOfChannels := int(*channels)

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	address := flag.Arg(0)

	// Create a client object to work with the io4edge device
	c, err := binaryiotypec.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	// configure the channels to use as outputs
	channnelConfig := make([]*biniopb.ChannelConfig, numberOfChannels)
	for i := 0; i < numberOfChannels; i++ {
		channnelConfig[i] = &biniopb.ChannelConfig{
			Channel:      int32(i),
			Mode:         biniopb.ChannelMode_BINARYIOTYPEC_OUTPUT_PUSH_PULL,
			InitialValue: false,
		}
	}

	if err := c.UploadConfiguration(
		binaryiotypec.WithChannelConfig(channnelConfig),
	); err != nil {
		log.Fatalf("Failed to upload configuration: %v\n", err)
	}

	// Run pattern through outputs
	i := 0
	for {
		prev := prevOutput(i, numberOfChannels)
		err := c.SetOutput(prev, false)
		if err != nil {
			log.Printf("can't switch off channel %d: %v", prev, err)
		}
		err = c.SetOutput(i, true)
		if err != nil {
			log.Printf("can't switch on channel %d: %v", i, err)
		}
		time.Sleep(300 * time.Millisecond)
		i++
		if i == numberOfChannels {
			i = 0
		}
	}
}

func prevOutput(i int, numChannels int) int {
	if i == 0 {
		return numChannels - 1
	}
	return i - 1
}
