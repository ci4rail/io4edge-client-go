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

	"github.com/ci4rail/io4edge-client-go/canl2"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/canL2/go/canL2/v1alpha1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	acceptanceCodePtr := flag.Uint("acceptancecode", 0x00000000, "CAN Filter Acceptance Code")
	acceptanceMaskPtr := flag.Uint("acceptancemask", 0x00000000, "CAN Filter Acceptance Mask")
	runtimePtr := flag.Uint("runtime", 10, "Seconds to receive data")
	lowLatency := flag.Bool("lowlatency", false, "Use stream low latency mode")

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

	// start stream
	err = c.StartStream(
		canl2.WithFilter(uint32(*acceptanceCodePtr), uint32(*acceptanceMaskPtr)),
		canl2.WithFBStreamOption(functionblock.WithBucketSamples(100)),
		canl2.WithFBStreamOption(functionblock.WithBufferedSamples(200)),
		canl2.WithFBStreamOption(functionblock.WithLowLatencyMode(*lowLatency)),
	)
	if err != nil {
		log.Fatalf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	readStreamFor(c, time.Duration(*runtimePtr)*time.Second)
}

func readStreamFor(c *canl2.Client, duration time.Duration) {
	start := time.Now()

	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 5)

		if err != nil {
			log.Printf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.Samples
			fmt.Printf("got stream data with %d samples\n", len(samples))

			for _, s := range samples {
				fmt.Printf("  %s\n", dumpSample(s))
			}
		}

		if time.Since(start) > duration {
			return
		}
	}
}

func dumpSample(sample *fspb.Sample) string {
	var s string

	s = fmt.Sprintf("@%010d us: ", sample.Timestamp)
	if sample.IsDataFrame {
		f := sample.Frame
		s += "ID:"
		if f.ExtendedFrameFormat {
			s += fmt.Sprintf("%08x", f.MessageId)
		} else {
			s += fmt.Sprintf("%03x", f.MessageId)
		}
		if f.RemoteFrame {
			s += " R"
		}
		s += " DATA:"
		for _, b := range f.Data {
			s += fmt.Sprintf("%02x ", b)
		}
		s += " "
	}
	s += "ERROR:" + sample.Error.String()
	s += " STATE:" + sample.ControllerState.String()

	return s
}
