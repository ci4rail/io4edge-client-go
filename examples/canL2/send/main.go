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
	fspb "github.com/ci4rail/io4edge_api/canL2/go/canL2/v1alpha1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	buckets := flag.Int("buckets", 10, "number of buckets to send")
	messages := flag.Int("messages", 5, "number of messages per bucket")
	extended := flag.Bool("ext", false, "use extended frames")
	rtr := flag.Bool("rtr", false, "use rtr frames")
	gap := flag.Int("gap", 0, "gap between buckets (ms)")
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

	for i := 0; i < *buckets; i++ {
		frames := []*fspb.Frame{}

		for j := 0; j < *messages; j++ {
			f := &fspb.Frame{
				MessageId:           uint32(0x100 + (i & 0xFF)),
				Data:                []byte{},
				ExtendedFrameFormat: *extended,
				RemoteFrame:         *rtr,
			}
			len := j % 8
			for k := 0; k < len; k++ {
				f.Data = append(f.Data, byte(j))
			}

			frames = append(frames, f)
		}

		start := time.Now()
		err = c.SendFrames(frames)
		elapsed := time.Since(start)
		fmt.Printf("Send of %d frames took %s\n", len(frames), elapsed)

		if err != nil {
			log.Printf("Send failed: %v\n", err)
		}
		if *gap != 0 {
			time.Sleep(time.Millisecond * time.Duration(*gap))
		}
	}
}
