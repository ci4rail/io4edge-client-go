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

	if len(os.Args) < 4 {
		log.Fatalf("Usage: parameter svc <mdns-service-address> <param_name> [<param_value>]  OR  parameter ip <ip:port> <param_name> [<param_value>]")
	}
	addressType := os.Args[1]
	address := os.Args[2]
	name := os.Args[3]

	var value string

	if len(os.Args) > 4 {
		value = os.Args[4]
	}

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

	if len(value) > 0 {
		err = c.SetPersistantParameter(name, value, timeout)
	} else {
		value, err = c.GetPersistantParameter(name, timeout)
		if err == nil {
			fmt.Printf("Parameter value: %s\n", value)
		}
	}
	if err != nil {
		log.Fatalf("Command failed: %v\n", err)
	}
}
