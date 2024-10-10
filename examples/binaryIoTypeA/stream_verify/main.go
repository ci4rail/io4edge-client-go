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

	binio "github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/binaryiotypea"
	fspb "github.com/ci4rail/io4edge_api/binaryIoTypeA/go/binaryIoTypeA/v1alpha1"
)

func manipulateOutputs(c *binio.Client, wg *sync.WaitGroup, quit chan bool) {
	go func() {
		var values uint8 = 0x00
		i := 0
		var direction int32 = 1
		for {
			select {
			case <-quit:
				wg.Done()
				return
			default:
				values += uint8(direction)
				//fmt.Printf("set:  %04b\n", values)
				err := c.SetAllOutputs(values, 0x0F)
				if err != nil {
					log.Printf("Failed to set all outputs: %v\n", err)
				}
				time.Sleep(time.Millisecond * 20) // 50 Hz is maximum output frequency
				i++
				if i%15 == 0 {
					direction *= -1
				}
			}
		}
	}()
}

func readStreamFor(c *binio.Client, duration time.Duration) {
	start := time.Now()
	values := uint8(0)
	var prevSample *fspb.Sample
	expectedValue := uint8(1)
	direction := 1

	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			fmt.Printf("got stream data seq=%d ts=%d\n", sd.Sequence, sd.DeliveryTimestamp)

			for i, sample := range samples {
				if prevSample != nil && sample.Timestamp != prevSample.Timestamp {
					fmt.Printf("     Reconstructed values ts=%d values=%04b(%d)\n", sample.Timestamp,
						values, values)

					if expectedValue != values {
						fmt.Printf("*** Expected value=%d, got %d", expectedValue, values)
						expectedValue = values
					}
					if expectedValue == 15 {
						direction = -1
					}
					if expectedValue == 0 {
						direction = 1
					}
					expectedValue += uint8(direction)
				}

				if sample.Value {
					values |= 1 << sample.Channel
				} else {
					values &= ^(1 << sample.Channel)
				}
				fmt.Printf("  #%d: ts=%d channel=%d %t\n", i, sample.Timestamp,
					sample.Channel, sample.Value)
				prevSample = sample

			}
		}
		if time.Since(start) > duration {
			return
		}
	}
}

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port>", os.Args[0])
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device
	var c *binio.Client
	var err error
	var quit chan bool = make(chan bool)
	var wg sync.WaitGroup = sync.WaitGroup{}

	c, err = binio.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create binio client: %v\n", err)
	}

	wg.Add(1)

	// start stream. Trigger on changes of all channels
	err = c.StartStream(binio.WithChannelFilterMask(0xf))
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	// manipulate outputs in background to force transitions
	manipulateOutputs(c, &wg, quit)

	readStreamFor(c, time.Second*10)

	quit <- true
	wg.Wait()
}
