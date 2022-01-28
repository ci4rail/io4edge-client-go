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

	binaryiotypea "github.com/ci4rail/io4edge-client-go/binaryiotypea"
)

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 3 {
		log.Fatalf("Usage: identify svc <mdns-service-address>  OR  identify ip <ip:port>")
	}
	addressType := os.Args[1]
	address := os.Args[2]

	// Create a client object to work with the io4edge device at <address>
	var c *binaryiotypea.Client
	var err error

	if addressType == "svc" {
		c, err = binaryiotypea.NewClientFromService(address, timeout)
	} else {
		c, err = binaryiotypea.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create binaryiotypea client: %v\n", err)
	}

	err = c.SetConfiguration(binaryiotypea.Configuration{
		OutputFritting:        -1,
		OutputWatchdog:        -1,
		OutputWatchdogTimeout: 11000,
	})
	if err != nil {
		fmt.Printf("Failed to set configuration: %v\n", err)
	}

	readConfig, err := c.GetConfiguration()
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
	} else {
		fmt.Printf("OutputFritting: %v\n", readConfig.OutputFritting)
		fmt.Printf("OutputWatchdog: %v\n", readConfig.OutputWatchdog)
		fmt.Printf("OutputWatchdogTimeout: %v\n", readConfig.OutputWatchdogTimeout)
	}

	err = c.SetConfiguration(binaryiotypea.Configuration{
		OutputFritting:        0x05,
		OutputWatchdog:        0,
		OutputWatchdogTimeout: 100,
	})
	if err != nil {
		fmt.Printf("Failed to set configuration: %v\n", err)
	}

	readConfig, err = c.GetConfiguration()
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
	} else {
		fmt.Printf("OutputFritting: %v\n", readConfig.OutputFritting)
		fmt.Printf("OutputWatchdog: %v\n", readConfig.OutputWatchdog)
		fmt.Printf("OutputWatchdogTimeout: %v\n", readConfig.OutputWatchdogTimeout)
	}

	err = c.SetConfiguration(binaryiotypea.Configuration{
		OutputFritting:        0x06,
		OutputWatchdog:        0x07,
		OutputWatchdogTimeout: 1250,
	})
	if err != nil {
		fmt.Printf("Failed to set configuration: %v\n", err)
	}

	readConfig, err = c.GetConfiguration()
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
	} else {
		fmt.Printf("OutputFritting: %v\n", readConfig.OutputFritting)
		fmt.Printf("OutputWatchdog: %v\n", readConfig.OutputWatchdog)
		fmt.Printf("OutputWatchdogTimeout: %v\n", readConfig.OutputWatchdogTimeout)
	}

	describe, err := c.Describe()
	if err != nil {
		fmt.Printf("Failed to config describe: %v\n", err)
	} else {
		fmt.Println("Describe: Number of channels: ", describe.NumberOfChannels)
	}
}
