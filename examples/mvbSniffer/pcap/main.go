/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/examples/mvbSniffer/pcap/busshark"
	"github.com/ci4rail/io4edge-client-go/examples/mvbSniffer/pcap/pcap"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	sniffer "github.com/ci4rail/io4edge-client-go/mvbsniffer"
	fspb "github.com/ci4rail/io4edge_api/mvbSniffer/go/mvbSniffer/v1"
)

func timeDeltaSnifferToHost(snifferTs uint64) uint64 {
	hostTs := time.Now().UnixMicro()
	return uint64(hostTs) - snifferTs
}

func sampleLine(sample *fspb.Sample) int {
	if sample.Line == fspb.Line_A {
		return 1
	}
	return 2
}

func sampleFrameType(sample *fspb.Sample) int {
	if sample.FrameType == fspb.FrameType_MASTER {
		return 1
	} else if sample.FrameType == fspb.FrameType_SLAVE {
		return 2
	} else {
		return 0
	}
}

func readStreamFor(c *sniffer.Client, w *pcap.Writer, duration time.Duration) {
	start := time.Now()
	frameNumber := uint64(0)

	for {
		if time.Since(start) > duration {
			return
		}
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
			continue
		}

		timeDelta := timeDeltaSnifferToHost(sd.DeliveryTimestamp)

		samples := sd.FSData.GetSamples()
		fmt.Printf("got stream data seq=%d pkts=%d ts=%d td=%v\n", sd.Sequence, len(samples), sd.DeliveryTimestamp, timeDelta)

		for i, sample := range samples {
			m := busshark.Pkt(frameNumber, 50*sample.Timestamp, sampleLine(sample), sampleFrameType(sample), sample.Payload)

			if err := w.AddPacket(sample.Timestamp+timeDelta, m); err != nil {
				log.Errorf("pcap add packet faile: %v\n", err)
			}

			fmt.Printf("  #%d: %v\n", i, sample)
		}
	}
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port> <pcap-file>", os.Args[0])
	}
	address := os.Args[1]
	pcapFile := os.Args[2]

	c, err := sniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create anain client: %v\n", err)
	}

	f, err := os.Create(pcapFile)
	if err != nil {
		log.Fatalf("Can't create file %s", err)
	}
	defer f.Close()

	w, err := pcap.NewWriter(f)
	if err != nil {
		log.Fatalf("Can't create writer %s", err)
	}

	// start stream
	err = c.StartStream(&functionblock.StreamConfiguration{
		BucketSamples:     100,
		BufferedSamples:   200,
		KeepaliveInterval: 1000,
	})
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	readStreamFor(c, w, time.Second*10)
}
