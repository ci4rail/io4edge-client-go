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
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge-client-go/templatemodule"
	templatemodulePb "github.com/ci4rail/io4edge_api/templateModule/go/templateModule/v1alpha1"
)

func handleSample(sample *templatemodulePb.Sample, sequenceNumber uint32) {
	fmt.Printf("\r")
	fmt.Printf("seq %d: %d: %d\n", sequenceNumber, sample.Timestamp, sample.Value)
}

func myrecover() {
	fmt.Println("Lost connecetion to io4edge device. Exiting now.")
}

func functionControl(c *templatemodule.Client, wg *sync.WaitGroup, quit chan bool) {
	go func() {
		var value uint32 = 0
		for {
			select {
			case <-quit:
				wg.Done()
				return
			default:
				value += 1000
				fmt.Printf("set:  %d\n", value)
				err := c.Set(value)
				if err != nil {
					log.Printf("Failed to set all channels: %v\n", err)
				}
				time.Sleep(time.Millisecond * 1000)
			}
		}
	}()
}

func configurationControl(c *templatemodule.Client) {
	err := c.SetConfiguration(templatemodule.Configuration{})
	if err != nil {
		fmt.Printf("Failed to set configuration: %v\n", err)
	}

	readConfig, err := c.GetConfiguration()
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
	} else {
		fmt.Printf("readConfig: %v\n", readConfig)
	}

	describe, err := c.Describe()
	if err != nil {
		fmt.Printf("Failed to config describe: %v\n", err)
	} else {
		fmt.Printf("Describe: %v", describe)
	}
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 3 {
		log.Fatalf("Usage: identify svc <mdns-service-address>  OR  identify ip <ip:port>")
	}
	addressType := os.Args[1]
	address := os.Args[2]

	// Create a client object to work with the io4edge device at <address>
	var c *templatemodule.Client
	var err error
	var quit chan bool = make(chan bool)
	var wg sync.WaitGroup = sync.WaitGroup{}

	if addressType == "svc" {
		c, err = templatemodule.NewClientFromService(address, timeout)
	} else {
		c, err = templatemodule.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create templatemodule client: %v\n", err)
	}

	c.SetRecover(myrecover)

	configurationControl(c)

	wg.Add(1)
	functionControl(c, &wg, quit)
	config := &templatemodule.StreamConfiguration{
		Generic: functionblock.StreamConfiguration{
			BucketSamples:     10,
			KeepaliveInterval: 100000,
			BufferedSamples:   30,
		},
	}
	err = c.StartStream(config, handleSample)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Started stream")
	time.Sleep(10 * time.Second)
	err = c.StopStream()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Stopped stream")
	}
	time.Sleep(1 * time.Second)
	err = c.StartStream(config, handleSample)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(40 * time.Second)
	err = c.StopStream()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Stopped stream")
	}

	quit <- true
	wg.Wait()

}
