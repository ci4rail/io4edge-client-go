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
	if len(os.Args) != 2 {
		log.Fatalf("Usage: identify <device-address>\n")
	}
	address := os.Args[1]

	c, err := core.NewClientFromSocketAddress(address)
	if err != nil {
		log.Fatalf("Failed to create basefunc client: %v\n", err)
	}

	fwID, err := c.IdentifyFirmware(5 * time.Second)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %d.%d.%d\n", fwID.Name, fwID.MajorVersion, fwID.MinorVersion, fwID.PatchVersion)

	hwID, err := c.IdentifyHardware(5 * time.Second)
	if err != nil {
		log.Fatalf("Failed to identify hardware: %v\n", err)
	}

	u, err := serno.FromSerial(hwID.SerialNumber.Hi, hwID.SerialNumber.Lo)
	if err != nil {
		log.Fatalf("Failed to make uuid from serial: %v\n", err)
	}
	fmt.Printf("Hardware name: %s, serial: %s, rev: %d\n", hwID.RootArticle, u.String(), hwID.MajorVersion)
}
