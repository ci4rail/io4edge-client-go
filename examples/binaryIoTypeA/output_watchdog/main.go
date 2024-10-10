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
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	binio "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/binaryiotypea"
)

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port>", os.Args[0])
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device
	var c *binio.Client
	var err error

	c, err = binio.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create binio client: %v\n", err)
	}

	// set 2 seconds watchdog on output0, but not on other outputs
	if err := c.UploadConfiguration(binio.WithOutputWatchdog(0x1, 2000)); err != nil {
		log.Fatalf("Failed to set configuration: %v\n", err)
	}

	// switch on output 0 and 1
	err = c.SetAllOutputs(0x3, 0xF)
	if err != nil {
		log.Fatalf("can't set outputs: %v", err)
	}

	time.Sleep(1500 * time.Millisecond)
	// output 0 and 1 must still be on
	values, _ := c.AllInputs(0x3)
	if values != 0x3 {
		log.Fatalf("Outputs did not turn on\n")
	}

	time.Sleep(700 * time.Millisecond)
	// now, the watchdog should have turned off output 0
	values, _ = c.AllInputs(0x3)
	if values != 0x2 {
		log.Fatalf("Watchdog did not turn off output 0\n")
	}
	fmt.Printf("Ok, watchdog turned off output 0\n")

	// now retrigger outputs just in time
	for i := 0; i < 10; i++ {
		_ = c.SetAllOutputs(0x3, 0xF)
		time.Sleep(1500 * time.Millisecond)

		values, _ := c.AllInputs(0xF)
		if values != 0x3 {
			log.Fatalf("Watchdog did turn off output 0 unintentionally\n")
		}
		fmt.Printf("inputs: %04b\n", values)
	}
	fmt.Printf("Ok, triggering in time succeeded\n")
}
