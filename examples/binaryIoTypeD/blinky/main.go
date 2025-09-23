/*
Copyright © 2025 Ci4Rail GmbH <engineering@ci4rail.com>

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

	binio "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/binaryiotyped"
	biniopb "github.com/ci4rail/io4edge_api/binaryIoTypeD/go/binaryIoTypeD/v1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	channels := flag.Uint("channels", 8, "Number of channels to use")
	flag.Parse()

	numberOfChannels := int(*channels)

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	address := flag.Arg(0)

	// Create a client object to work with the io4edge device
	c, err := binio.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	// configure the first half of the channels to use as outputs, the rest as inputs
	channelConfig := make([]*biniopb.ChannelConfig, numberOfChannels)
	for i := 0; i < numberOfChannels; i++ {
		if i < numberOfChannels/2 {
			channelConfig[i] = &biniopb.ChannelConfig{
				Channel:                   int32(i),
				Mode:                      biniopb.ChannelMode_BINARYIOTYPED_OUTPUT_HIGH_ACTIVE,
				InitialValue:              false,
				OverloadRecoveryTimeoutMs: 50,
				WatchdogTimeoutMs:         2000,
			}
		} else {
			channelConfig[i] = &biniopb.ChannelConfig{
				Channel:        int32(i),
				Mode:           biniopb.ChannelMode_BINARYIOTYPED_INPUT_HIGH_ACTIVE,
				FrittingEnable: true,
			}
		}
	}

	if err := c.UploadConfiguration(
		binio.WithChannelConfig(channelConfig),
	); err != nil {
		log.Fatalf("Failed to upload configuration: %v\n", err)
	}

	errorCount := 0
	for {
		for channel := 0; channel < numberOfChannels/2; channel++ {
			for _, state := range []bool{true, false} {
				err := c.SetOutput(channel, state)
				if err != nil {
					log.Printf("can't switch channel %d to %v: %v", channel, state, err)
					errorCount++
					continue
				}
				time.Sleep(1000 * time.Millisecond)

				errorCount += checkChannels(c, 1<<channel+1<<(numberOfChannels/2))
				fmt.Printf("Errors so far: %d\n", errorCount)
			}
		}

	}
}

func checkChannels(c *binio.Client, chMask uint32) int {
	errorCount := 0
	inputs, diags, err := c.Inputs()
	if err != nil {
		log.Printf("can't read inputs: %v", err)
		return 1
	}
	if inputs != chMask {
		log.Printf("input mismatch: expected %08b, got %08b", chMask, inputs)
		errorCount++
	}
	for _, diag := range diags {
		if diag != uint32(biniopb.ChannelDiag_NoDiag) {
			log.Printf("diagnostic error: %v", diag)
			errorCount++
		}
	}
	return errorCount
}
