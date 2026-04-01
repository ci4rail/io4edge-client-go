/*
Copyright © 2026 Ci4Rail GmbH <engineering@ci4rail.com>

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

// On SQ1, connect IO1 with IO2 and run example with -channels 2

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"log"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	binio "github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/binaryiotyped"
	biniopb "github.com/ci4rail/io4edge_api/binaryIoTypeD/go/binaryIoTypeD/v1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	channels := flag.Uint("channels", 4, "Number of channels to use")
	runtime := flag.Uint("runtime", 10, "Runtime in seconds")
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

	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(1)
	var quit chan bool = make(chan bool)
	start := time.Now()
	err = c.StartStream(binio.WithFBStreamOption(functionblock.WithBucketSamples(100)),
		binio.WithFBStreamOption(functionblock.WithBufferedSamples(200)),
		binio.WithFBStreamOption(functionblock.WithLowLatencyMode(true)),
	)
	if err != nil {
		log.Fatalf("StartStream failed: %v\n", err)
	}
	fmt.Printf("Started stream, running for %d seconds\n", *runtime)

	// manipulate outputs in background to force transitions
	manipulateOutputs(c, numberOfChannels/2, &wg, quit)
	// stop manipulation when this function exits
	defer func() {
		quit <- true
		wg.Wait()
	}()

	firstTs := uint64(0)
	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 5)

		if err != nil {
			log.Fatalf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			log.Printf("got stream data seq=%d ts=%d samples=%d\n", sd.Sequence, sd.DeliveryTimestamp, len(samples))

			for i, sample := range samples {
				if firstTs == 0 {
					firstTs = sample.Timestamp
				}
				log.Printf("sample %d: relTs=%10dus values=b%0*b", i, sample.Timestamp-firstTs, numberOfChannels, sample.Inputs&((1<<numberOfChannels)-1))
				for ch, diag := range sample.Diag {
					if ch < numberOfChannels && diag != 0 {
						log.Printf(" diag ch %d=%d", ch, diag)
					}
				}
			}
		}
		if time.Since(start) > time.Second*time.Duration(*runtime) {
			break
		}
	}
}

func manipulateOutputs(c *binio.Client, numberOfChannels int, wg *sync.WaitGroup, quit chan bool) {
	go func() {
		i := uint32(0)
		chMask := uint32((1 << numberOfChannels) - 1)
		for {
			select {
			case <-quit:
				wg.Done()
				return
			default:
				fmt.Printf("Setting outputs to %08b\n", i)
				_, _, err := c.SetOutputs(i, chMask)
				if err != nil {
					log.Printf("can't set outputs: %v", err)
				}
				time.Sleep(time.Millisecond * 50)
				i++
			}
		}
	}()
}
