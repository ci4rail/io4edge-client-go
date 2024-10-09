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

	"github.com/ci4rail/io4edge-client-go/coreclient"
	"github.com/ci4rail/io4edge-client-go/socketcore"
)

func progressCb(bytes uint, msg string) {
	if msg != "" {
		fmt.Printf("\n%s\n", msg)
	}
	fmt.Printf("\r%d kBytes loaded.", bytes/1024)
}

func createClient(addressType string, address string, timeout time.Duration) coreclient.If {
	var c coreclient.If
	var err error

	if addressType == "svc" {
		c, err = socketcore.NewClientFromService(address, timeout)
	} else {
		c, err = socketcore.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create core client: %v\n", err)
	}
	return c
}

// Example Usage:
//   load svc S103-SIO01-14af003b-0591-4358-bf44-07d1d8af1cf5._io4edge-core._tcp fw-sio01-default-0.1.2.fwpkg
// OR
//   load ip 192.168.24.233:9999 fw-sio01-default-0.1.2.fwpkg

func main() {
	const timeout = 0 // use default timeout
	const chunkSize = 1024

	if len(os.Args) != 4 {
		log.Fatalf("Usage: load svc <mdns-service-address> <fwpkg> OR  load ip <ip:port> <fwpkg>")
	}
	addressType := os.Args[1]
	address := os.Args[2]
	file := os.Args[3]

	// Create a client object to work with the io4edge device at <address>
	c := createClient(addressType, address, timeout)

	// Load the firmware package into the device
	// Loading happens in chunks of <chunkSize>. 1024 should work with each device
	// <timeout> is the time to wait for responses from device
	restartingNow, err := c.LoadFirmware(file, chunkSize, timeout, progressCb)
	if err != nil {
		log.Fatalf("Failed to load firmware package: %v\n", err)
	}

	log.Printf("Load succeeded. Reading back firmware ID\n")

	if restartingNow {
		fmt.Println("Let device restart...")
		// Must sleep here, because device needs approx. 1s to reset
		// If we reconnect too fast, we can still reach the device before it is reset
		time.Sleep(3 * time.Second)

		// must create a new client, device has rebooted
		c = createClient(addressType, address, timeout)
		if err != nil {
			log.Fatalf("Failed to create core client: %v\n", err)
		}
	}

	// Get the now active firmware version from the device
	fwName, fwVersion, err := c.IdentifyFirmware(timeout)
	if err != nil {
		log.Fatalf("Failed to identify firmware: %v\n", err)
	}

	fmt.Printf("Firmware name: %s, Version %s\n", fwName, fwVersion)

}
