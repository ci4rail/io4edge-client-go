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
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/functionblock"
	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/mvbsniffer"
	mvbpb "github.com/ci4rail/io4edge_api/mvbSniffer/go/mvbSniffer/v1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	runtime := flag.Uint("runtime", 10, "Seconds to receive data")
	useInternalGenerator := flag.Bool("gen", false, "Use internal generator")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	address := flag.Arg(0)

	// Create a client object to work with the io4edge device
	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	if *useInternalGenerator {
		startInternalGenerator(c)
	} else {
		// ensure we use the external input in case the internal generator has been selected before
		err := c.SendPattern(c.SelectExternalInputString())
		if err != nil {
			log.Fatalf("SendPattern failed: %v\n", err)
		}
	}

	// start stream
	err = c.StartStream(
		mvbsniffer.WithFilterMask(mvbsniffer.FilterMask{
			// receive any telegram, except timed out frames
			FCodeMask:             0xFFFF,
			Address:               0x0000,
			Mask:                  0x0000,
			IncludeTimedoutFrames: false,
		}),
		mvbsniffer.WithFBStreamOption(functionblock.WithBucketSamples(100)),
		mvbsniffer.WithFBStreamOption(functionblock.WithBufferedSamples(200)),
	)
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	readStreamFor(c, time.Second*time.Duration(*runtime))
}

func readStreamFor(c *mvbsniffer.Client, duration time.Duration) {
	start := time.Now()
	prevTs := uint64(0)
	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			telegramCollection := sd.FSData.GetEntry()

			for _, telegram := range telegramCollection {
				dt := uint64(0)
				if prevTs != 0 {
					dt = telegram.Timestamp - prevTs
				}
				prevTs = telegram.Timestamp

				if telegram.State != uint32(mvbpb.Telegram_kSuccessful) {
					if telegram.State&uint32(mvbpb.Telegram_kTimedOut) != 0 {
						log.Errorf("No slave frame has been received to a master frame\n")
					}
					if telegram.State&uint32(mvbpb.Telegram_kMissedMVBFrames) != 0 {
						log.Errorf("one or more MVB frames are lost in the device since the last telegram\n")
					}
					if telegram.State&uint32(mvbpb.Telegram_kMissedTelegrams) != 0 {
						log.Errorf("one or more telegrams are lost\n")
					}
				}

				fmt.Printf("dt=%d %s\n", dt, telegramToString(telegram))
			}
		}
		if time.Since(start) > duration {
			return
		}
	}
}

func telegramToString(t *mvbpb.Telegram) string {
	s := fmt.Sprintf("addr=%06x, ", t.Address)
	s += fmt.Sprintf("%s, ", mvbpb.Telegram_Type_name[int32(t.Type)])
	if len(t.Data) > 0 {
		s += "data="
		for i := 0; i < len(t.Data); i++ {
			s += fmt.Sprintf("%02x ", t.Data[i])
		}
		s += ", "
	}
	return s
}

func startInternalGenerator(c *mvbsniffer.Client) {
	// start internal generator
	cl := mvbsniffer.NewCommandList()

	// Master frame: address 0x123, processdata 32 bit, wait 30us after frame
	cl.AddMasterFrame(0, false, 30, 1, 0x123)
	// Slave frame: respond to master, wait 200us after frame
	cl.AddSlaveFrame(0, false, 200, []uint8{0x01, 0x02, 0x03, 0x04})

	// Master frame: address 0x456, processdata 16 bit, wait 30us after frame
	cl.AddMasterFrame(0, false, 30, 0, 0x456)
	// Slave frame: respond to master, wait 800us after frame
	cl.AddSlaveFrame(0, false, 800, []uint8{0xaa, 0xbb})

	commandString := cl.StartGeneratorString(true)
	err := c.SendPattern(commandString)
	if err != nil {
		log.Fatalf("SendPattern failed: %v\n", err)
	}
}
