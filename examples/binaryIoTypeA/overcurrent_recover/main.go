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
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	binio "github.com/ci4rail/io4edge-client-go/binaryiotypea"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
)

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address>  OR  %s <ip:port>", os.Args[0], os.Args[0])
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device
	var c *binio.Client
	var err error

	c, err = binio.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create binio client: %v\n", err)
	}

	// wait until output 0 is short circuit manually
	for {
		val, err := c.Input(0)
		if err != nil {
			log.Errorf("can't get input: %v", err)
		}

		err = c.SetOutput(0, true)

		if err != nil {
			if functionblock.HaveResponseStatus(err, v1alpha1.Status_HW_FAULT) {
				break
			} else {
				log.Errorf("can't set output: %v", err)
			}
		}
		fmt.Printf("Waiting for overcurrent. Current status is %t\n", val)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Detected overcurrent. Try recovery\n")

	for {
		err := c.ExitErrorState()
		if err != nil {
			log.Errorf("can't execute exit error state: %v", err)
		}
		err = c.SetOutput(0, true)

		if err != nil {
			log.Errorf("can't set output: %v", err)
		} else {
			break
		}

		fmt.Printf("Waiting for overcurrent to recover\n")

		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Overcurrent recovered\n")
}
