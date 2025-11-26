/*
Copyright © 2025 Ci4Rail GmbH <engineering@ci4rail.com>

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

	log "github.com/sirupsen/logrus"

	anain "github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/analogintypeb"
)

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address> OR <ip:port>", os.Args[0])
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device
	c, err := anain.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create anain client: %v\n", err)
	}

	// read current values
	values, err := c.Values()
	if err != nil {
		log.Fatalf("Can't get value: %v\n", err)
	}
	for i, val := range values {
		fmt.Printf("Current value channel %d: %.4f\n", i, val)
	}
}
