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

	"github.com/ci4rail/hw-inventory-go/serno"
	"github.com/ci4rail/io4edge-client-go/core"
)

func main() {
	// timeout to wait for the service to show up
	const timeoutService = 1 * time.Second
	// timeout to wait for the device to respond
	const timeoutCmd = 5 * time.Second

	if len(os.Args) != 2 {
		log.Fatalf("Usage: identify <device-address>\n")
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device at <address>
	c, err := core.NewClientFromService(address, timeoutService)
	if err != nil {
		log.Fatalf("Failed to create core client: %v\n", err)
	}

	// Get the active firmware version from the device
	fwID, err := c.IdentifyFirmware(timeoutCmd)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %s\n", fwID.Name, fwID.Version)

	// Get the hardware name and version from the device
	hwID, err := c.IdentifyHardware(timeoutCmd)
	if err != nil {
		log.Fatalf("Failed to identify hardware: %v\n", err)
	}

	// device reports its serial number into two 64 bit values. Convert it to UUID
	u, err := serno.UUIDFromInt(hwID.SerialNumber.Hi, hwID.SerialNumber.Lo)
	if err != nil {
		log.Fatalf("Failed to make uuid from serial: %v\n", err)
	}
	fmt.Printf("Hardware name: %s, serial: %s, rev: %d\n", hwID.RootArticle, u.String(), hwID.MajorVersion)
}
