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
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	binio "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/binaryiotypeb"
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

	// Set all outputs
	for values := uint8(0); values < 4; values++ {
		err = c.SetAllOutputs(values, 0x3)
		if err != nil {
			log.Fatalf("Failed to set outputs: %v\n", err)
		}
		log.Printf("Set outputs to %x", values)
		time.Sleep(1500 * time.Millisecond)
	}

	// Reset outputs
	err = c.SetAllOutputs(0x0, 0x3)
	if err != nil {
		log.Fatalf("Failed to set outputs: %v\n", err)
	}

	time.Sleep(1000 * time.Millisecond)

	// Set a single output
	err = c.SetOutput(0, true)
	if err != nil {
		log.Fatalf("Failed to set output: %v\n", err)
	}

	time.Sleep(1000 * time.Millisecond)

	// Reset output
	err = c.SetOutput(0, false)
	if err != nil {
		log.Fatalf("Failed to set output: %v\n", err)
	}
}
