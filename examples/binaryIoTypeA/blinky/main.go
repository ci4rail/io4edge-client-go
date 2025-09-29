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

	binio "github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/binaryiotypea"
)

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port>", os.Args[0])
	}
	address := os.Args[1]

	c, err := binio.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create binio client: %v\n", err)
	}

	desc, err := c.Describe()
	if err != nil {
		log.Fatalf("can't get description: %v", err)
	}

	// Run pattern through outputs
	i := 0
	for {
		err := c.SetOutput(prevOutput(i, desc.NumberOfChannels), false)
		if err != nil {
			log.Fatalf("can't switch off: %v", err)
		}
		err = c.SetOutput(i, true)
		if err != nil {
			log.Fatalf("can't switch on: %v", err)
		}
		time.Sleep(300 * time.Millisecond)
		i++
		if i == desc.NumberOfChannels {
			i = 0
		}
	}
}

func prevOutput(i int, numChannels int) int {
	if i == 0 {
		return numChannels - 1
	}
	return i - 1
}
