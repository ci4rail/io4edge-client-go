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
	"flag"
	"fmt"
	"os"
	"time"

	"log"

	binio "github.com/ci4rail/io4edge-client-go/binaryiotypec"
	biniopb "github.com/ci4rail/io4edge_api/binaryIoTypeC/go/binaryIoTypeC/v1alpha1"
)

// Hardwaresetup: Connect channel 0-3 to channel 8-11 and channel 4-7 to channel 12-15
func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s <mdns-service-address OR ip:port>\n", os.Args[0])
		os.Exit(1)
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	address := flag.Arg(0)

	// Create a client object to work with the io4edge device
	c, err := binio.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	// configure the channels 0-7 to use as outputs and channels 8-15 to use as inputs
	channnelConfig := make([]*biniopb.ChannelConfig, 16)
	for i := 0; i < 16; i++ {
		var mode biniopb.ChannelMode
		if i > 7 {
			mode = biniopb.ChannelMode_BINARYIOTYPEC_INPUT_TYPE_1_3
		} else {
			mode = biniopb.ChannelMode_BINARYIOTYPEC_OUTPUT_PUSH_PULL
		}
		channnelConfig[i] = &biniopb.ChannelConfig{
			Channel:      int32(i),
			Mode:         mode,
			InitialValue: false,
		}
	}

	if err := c.UploadConfiguration(
		binio.WithChannelConfig(channnelConfig),
	); err != nil {
		log.Fatalf("Failed to upload configuration: %v\n", err)
	}

	outputs := uint32(0)
	// genrate a pattern on the outputs and read back the inputs
	for {
		err := c.SetAllOutputs(outputs, 0x00ff)
		if err != nil {
			log.Printf("can't set outputs: %v", err)
		}
		log.Printf("Set Outputs: %08b", outputs)
		time.Sleep(time.Millisecond * 200)
		inputs, diags, err := c.AllInputs()
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		for channel, diag := range diags {
			if diag&uint32(biniopb.ChannelDiag_NoSupplyVoltage) != 0 {
				log.Printf("Group of channel %d has no power", channel)
			}
			if diag&uint32(biniopb.ChannelDiag_CurrentLimit) != 0 {
				log.Printf("Detected overcurrent on channel %d", channel)
			}
			if diag&uint32(biniopb.ChannelDiag_Overload) != 0 {
				log.Printf("Detected overload on channel %d", channel)
			}
			if diag&uint32(biniopb.ChannelDiag_SupplyUndervoltage) != 0 {
				log.Printf("Detected undervoltage on channel %d", channel)
			}
			if diag&uint32(biniopb.ChannelDiag_SupplyOvervoltage) != 0 {
				log.Printf("Detected overvoltage on channel %d", channel)
			}
		}
		inputs = inputs >> 8
		log.Printf("Get Inputs:  %08b", inputs)
		if outputs != inputs {
			log.Printf("Error: Inputs don't match outputs")
		}
		outputs++
		if outputs > 0xff {
			outputs = 0
		}
	}
}
