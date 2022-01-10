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

	"github.com/ci4rail/io4edge-client-go/binaryIoTypeA"
	binio "github.com/ci4rail/io4edge-client-go/binaryIoTypeA/v1alpha1"
)

func handleSample(sample *binio.Sample) {
	fmt.Println(sample.Timestamp, sample.Valid, sample.Channel, sample.Value)
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 3 {
		log.Fatalf("Usage: identify svc <mdns-service-address>  OR  identify ip <ip:port>")
	}
	addressType := os.Args[1]
	address := os.Args[2]

	// Create a client object to work with the io4edge device at <address>
	var c *binaryIoTypeA.Client
	var err error

	if addressType == "svc" {
		c, err = binaryIoTypeA.NewClientFromService(address, timeout)
	} else {
		c, err = binaryIoTypeA.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create binaryIoTypeA client: %v\n", err)
	}
	id, err := c.StartStream(0, handleSample)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Started stream with id %d\n", id)
	time.Sleep(60 * time.Second)
	err = c.StopStream(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Stopped stream with id %d\n", id)
}
