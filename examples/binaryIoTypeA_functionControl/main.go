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
	"sync"
	"time"

	"github.com/ci4rail/io4edge-client-go/binaryIoTypeA"
)

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

	quit := make(chan interface{})
	var wg sync.WaitGroup

	////////////////////////////////////////////////////////////////////////////
	// fmt.Println("Set Single example")
	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	quit <- true
	// }()

	// wg.Add(1)
	// go func() {
	// 	state := false
	// 	for {
	// 		select {
	// 		case <-quit:
	// 			wg.Done()
	// 			return
	// 		default:
	// 			state = !state
	// 			fmt.Printf("SetSingle(0, %v)\n", !state)
	// 			err = c.SetSingle(0, !state)
	// 			if err != nil {
	// 				log.Fatalf("Failed to set single channel: %v\n", err)
	// 			}
	// 			time.Sleep(20 * time.Millisecond)
	// 			readState, err := c.GetSingle(0)
	// 			if err != nil {
	// 				log.Fatalf("Failed to get single channel: %v\n", err)
	// 			}
	// 			fmt.Printf("GetSingle(0): %t)\n", readState)

	// 			fmt.Printf("SetSingle(1, %v)\n", state)
	// 			err = c.SetSingle(1, state)
	// 			if err != nil {
	// 				log.Fatalf("Failed to set single channel: %v\n", err)
	// 			}
	// 			time.Sleep(20 * time.Millisecond)
	// 			readState, err = c.GetSingle(1)
	// 			if err != nil {
	// 				log.Fatalf("Failed to get single channel: %v\n", err)
	// 			}
	// 			fmt.Printf("GetSingle(1): %t)\n", readState)

	// 			fmt.Printf("SetSingle(2, %v)\n", !state)
	// 			err = c.SetSingle(2, !state)
	// 			if err != nil {
	// 				log.Fatalf("Failed to set single channel: %v\n", err)
	// 			}
	// 			time.Sleep(20 * time.Millisecond)
	// 			readState, err = c.GetSingle(2)
	// 			if err != nil {
	// 				log.Fatalf("Failed to get single channel: %v\n", err)
	// 			}
	// 			fmt.Printf("GetSingle(2): %t)\n", readState)

	// 			fmt.Printf("SetSingle(3, %v)\n", state)
	// 			err = c.SetSingle(3, state)
	// 			if err != nil {
	// 				log.Fatalf("Failed to set single channel: %v\n", err)
	// 			}
	// 			time.Sleep(20 * time.Millisecond)
	// 			readState, err = c.GetSingle(3)
	// 			if err != nil {
	// 				log.Fatalf("Failed to get single channel: %v\n", err)
	// 			}
	// 			fmt.Printf("GetSingle(3): %t)\n", readState)
	// 			fmt.Println()
	// 			time.Sleep(500 * time.Millisecond)
	// 		}
	// 	}
	// }()
	// wg.Wait()

	////////////////////////////////////////////////////////////////////////////
	fmt.Println("Set All example modifiying values bitmask")
	go func() {
		time.Sleep(1000000 * time.Second)
		quit <- true
	}()

	wg.Add(1)
	go func() {
		var values uint32 = 0x00
		for {
			select {
			case <-quit:
				wg.Done()
				return
			default:
				for i := 0; i < 4; i++ {
					values = setBit(values, i)
					fmt.Printf("set:  %04b\n", values)
					err := c.SetAll(values, 0x0F)
					if err != nil {
						log.Printf("Failed to set all channels: %v\n", err)
					}
					time.Sleep(time.Millisecond * 100)
					readState, err := c.GetAll(0xFF)
					if err != nil {
						log.Printf("Failed to get all channels: %v\n", err)
					}
					fmt.Printf("read: %04b\n", readState)
					time.Sleep(time.Millisecond * 100)
				}
				for i := 3; i >= 0; i-- {
					values = clearBit(values, i)
					fmt.Printf("set:  %04b\n", values)
					err := c.SetAll(values, 0x0F)
					if err != nil {
						log.Printf("Failed to set all channels: %v\n", err)
					}
					time.Sleep(time.Millisecond * 100)
					readState, err := c.GetAll(0xFF)
					if err != nil {
						log.Printf("Failed to get all channels: %v\n", err)
					}
					fmt.Printf("read: %04b\n", readState)
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	}()
	wg.Wait()

	////////////////////////////////////////////////////////////////////////////
	fmt.Println("Set All example modifiying filter bitmask")
	go func() {
		time.Sleep(10 * time.Second)
		quit <- true
	}()

	wg.Add(1)
	go func() {
		var mask uint32 = 0x0F
		for {
			select {
			case <-quit:
				wg.Done()
				return
			default:
				for i := 0; i < 3; i++ {
					mask = mask >> 1
					err := c.SetAll(0x00, 0x0F)
					if err != nil {
						log.Fatalf("Failed to set all channels: %v\n", err)
					}
					err = c.SetAll(0x0F, mask)
					fmt.Printf("mask: %04b\n", mask)
					if err != nil {
						log.Fatalf("Failed to set all channels: %v\n", err)
					}
					time.Sleep(time.Millisecond * 500)
				}
				for i := 0; i < 3; i++ {
					err := c.SetAll(0x00, 0x0F)
					if err != nil {
						log.Fatalf("Failed to set all channels: %v\n", err)
					}
					mask = mask << 1
					fmt.Printf("mask: %04b\n", mask)
					err = c.SetAll(0x0F, mask)
					if err != nil {
						log.Fatalf("Failed to set all channels: %v\n", err)
					}
					time.Sleep(time.Millisecond * 500)
				}
			}
		}
	}()
	wg.Wait()

	// Turn off everything
	err = c.SetAll(0x00, 0x0F)
	if err != nil {
		log.Fatalf("Failed to set all channels: %v\n", err)
	}
}

func setBit(n uint32, pos int) uint32 {
	n |= (1 << pos)
	return n
}

func clearBit(n uint32, pos int) uint32 {
	var mask uint32 = ^(1 << pos)
	n &= mask
	return n
}
