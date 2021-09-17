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
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/core"
)

func main() {
	const timeout = 5 * time.Second
	const chunkSize = 1024

	if len(os.Args) != 3 {
		log.Fatalf("Usage: load  <device-address> <fwpkg>\n")
	}
	address := os.Args[1]
	file := os.Args[2]

	// Create a client object to work with the io4edge device at <address>
	c, err := core.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create core client: %v\n", err)
	}

	// Load the firmware package into the device
	// Loading happens in chunks of <chunkSize>. 1024 should work with each device
	// <timeout> is the time to wait for responses from device
	restartingNow, err := c.LoadFirmware(file, chunkSize, timeout)
	if err != nil {
		log.Fatalf("Failed to load firmware package: %v\n", err)
	}

	log.Printf("Load succeeded. Reading back firmware ID\n")

	if restartingNow {
		// must create a new client, device has rebooted
		c, err = core.NewClientFromSocketAddress(address)
		if err != nil {
			log.Fatalf("Failed to create core client: %v\n", err)
		}
	}

	// Get the now active firmware version from the device
	fwID, err := c.IdentifyFirmware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %s\n", fwID.Name, fwID.Version)

}
