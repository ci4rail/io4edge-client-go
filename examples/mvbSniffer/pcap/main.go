/*
Copyright © 2024 Ci4Rail GmbH <engineering@ci4rail.com>

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
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/examples/mvbSniffer/pcap/busshark"
	"github.com/ci4rail/io4edge-client-go/examples/mvbSniffer/pcap/pcap"
	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/functionblock"
	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/mvbsniffer"
	fspb "github.com/ci4rail/io4edge_api/mvbSniffer/go/mvbSniffer/v1"
)

// TODO: This is not 100% correct, as deliverytimestamp is ESP time and sample timestamps are from FPGA
func timeDeltaSnifferToHost(snifferTs uint64) uint64 {
	hostTs := time.Now().UnixMicro()
	return uint64(hostTs) - snifferTs
}

func readStreamFor(c *mvbsniffer.Client, w *pcap.Writer, duration time.Duration) {
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

		telegramCollection := sd.FSData.GetEntry()
		fmt.Printf("got stream data seq=%d pkts=%d ts=%d td=%v\n", sd.Sequence, len(telegramCollection), sd.DeliveryTimestamp, timeDelta)

		for i, telegram := range telegramCollection {
			// generate fake master packet
			m := busshark.Pkt(frameNumber, 50*telegram.Timestamp,
				int(telegram.Line)+1, 1 /*Master*/, []byte{byte(uint8(telegram.Type)<<4 + uint8(telegram.Address>>8)), byte(telegram.Address & 0xff)})

			if err := w.AddPacket(telegram.Timestamp+timeDelta, m); err != nil {
				log.Errorf("pcap add packet faile: %v\n", err)
			}
			frameNumber++

			if telegram.State&uint32(fspb.Telegram_kTimedOut) == 0 {
				m = busshark.Pkt(frameNumber, 50*telegram.Timestamp, int(telegram.Line)+1, 2 /*Slave*/, telegram.Data)

				if err := w.AddPacket(telegram.Timestamp+timeDelta, m); err != nil {
					log.Errorf("pcap add packet faile: %v\n", err)
				}

				frameNumber++
			}

			if telegram.State&uint32(fspb.Telegram_kMissedMVBFrames|fspb.Telegram_kMissedTelegrams) != 0 {
				fmt.Printf("  #%d: %v\n", i, telegram)
			}
		}
	}
}

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port> <pcap-file> <recordingtime-seconds>", os.Args[0])
	}
	address := os.Args[1]
	pcapFile := os.Args[2]
	recordingTimeSeconds, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Failed to convert recording time: %v\n", err)
	}
	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
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
	err = c.StartStream(
		mvbsniffer.WithFilterMask(mvbsniffer.FilterMask{
			// receive any telegram, except timed out frames
			FCodeMask:             0xFFFF,
			Address:               0x0000,
			Mask:                  0x0000,
			IncludeTimedoutFrames: true,
		}),
		mvbsniffer.WithFBStreamOption(functionblock.WithBucketSamples(300)),
		mvbsniffer.WithFBStreamOption(functionblock.WithBufferedSamples(600)),
	)
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	readStreamFor(c, w, time.Second*time.Duration(recordingTimeSeconds))
}
