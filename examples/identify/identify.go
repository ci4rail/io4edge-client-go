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
	const timeout = 5 * time.Second

	if len(os.Args) != 2 {
		log.Fatalf("Usage: identify <device-address>\n")
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device at <address>
	c, err := core.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	// Get the active firmware version from the device
	fwID, err := c.IdentifyFirmware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %d.%d.%d\n", fwID.Name, fwID.MajorVersion, fwID.MinorVersion, fwID.PatchVersion)

	// Get the hardware name and version from the device
	hwID, err := c.IdentifyHardware(timeout)
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
