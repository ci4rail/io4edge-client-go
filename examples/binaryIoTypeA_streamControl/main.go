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

	"github.com/ci4rail/io4edge-client-go/binaryiotypea"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	binio "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
)

func handleSample(sample *binio.Sample, sequenceNumber uint32) {
	fmt.Printf("\r")
	fmt.Printf("seq %d: %d: channel %d: %t\n", sequenceNumber, sample.Timestamp, sample.Channel, sample.Value)
}

func myrecover() {
	fmt.Println("Lost connecetion to io4edge device. Exiting now.")
}

func functionControl(c *binaryiotypea.Client, wg *sync.WaitGroup, quit chan bool) {
	go func() {
		var values uint32 = 0x00
		i := 0
		var direction int32 = 1
		for {
			select {
			case <-quit:
				wg.Done()
				return
			default:
				values += uint32(direction)
				fmt.Printf("set:  %04b\n", values)
				err := c.SetAll(values, 0x0F)
				if err != nil {
					log.Printf("Failed to set all channels: %v\n", err)
				}
				time.Sleep(time.Millisecond * 10)
				i++
				if i%15 == 0 {
					direction *= -1
				}
			}
		}
	}()
}

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
	var quit chan bool = make(chan bool)
	var wg sync.WaitGroup = sync.WaitGroup{}

	if addressType == "svc" {
		c, err = binaryiotypea.NewClientFromService(address, timeout)
	} else {
		c, err = binaryiotypea.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create binaryiotypea client: %v\n", err)
	}
	c.SetRecover(myrecover)
	wg.Add(1)
	functionControl(c, &wg, quit)
	config := &binaryiotypea.StreamConfiguration{
		Generic: functionblock.StreamConfiguration{
			BucketSamples:     10,
			KeepaliveInterval: 100000,
			BufferedSamples:   30,
		},
		ChannelFilterMask: 0xF,
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

	time.Sleep(10 * time.Second)
	err = c.StopStream()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Stopped stream")
	}

	quit <- true
	wg.Wait()

}
