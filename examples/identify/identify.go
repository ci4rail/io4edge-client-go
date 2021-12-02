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

	if len(os.Args) != 3 {
		log.Fatalf("Usage: identify svc <mdns-service-address>  OR  identify ip <ip:port>")
	}
	addressType := os.Args[1]
	address := os.Args[2]

	// Create a client object to work with the io4edge device at <address>
	var c *core.Client
	var err error

	if addressType == "svc" {
		c, err = core.NewClientFromService(address, timeout)
	} else {
		c, err = core.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create core client: %v\n", err)
	}

	// Get the active firmware version from the device
	fwName, fwVersion, err := c.IdentifyFirmware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %s\n", fwName, fwVersion)

	// Get the hardware name and version from the device
	rootArticle, majorVersion, serialNumber, err := c.IdentifyHardware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify hardware: %v\n", err)
	}

	fmt.Printf("Hardware name: %s, serial: %s, rev: %d\n", rootArticle, serialNumber, majorVersion)
}
