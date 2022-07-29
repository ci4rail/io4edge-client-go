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
	fspb "github.com/ci4rail/io4edge_api/canL2/go/canL2/v1alpha1"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

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

	for i := 0; i < 10; i++ {
		frames := []*fspb.Frame{}

		for j := 0; j < 5; j++ {
			f := &fspb.Frame{
				MessageId:           uint32(i),
				Data:                []byte{byte(j)},
				ExtendedFrameFormat: false,
				RemoteFrame:         false,
			}
			frames = append(frames, f)
		}
		fmt.Printf("Send frames: %v\n", frames)
		err = c.SendFrames(frames)
		if err != nil {
			log.Printf("Send failed: %v\n", err)
		}
	}
}
