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
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	anain "github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/analogintypeb"
	fspb "github.com/ci4rail/io4edge_api/analogInTypeB/go/analogInTypeB/v1"
)

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <mdns-service-address> <sample-rate>  OR  %s <ip:port> <sample-rate>", os.Args[0], os.Args[0])
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device
	var c *anain.Client
	var err error

	sampleRate, err := strconv.ParseFloat(os.Args[2], 32)
	if err != nil {
		log.Fatalf("Can't convert sample rate: %v\n", err)
	}

	c, err = anain.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create anain client: %v\n", err)
	}

	// read the channel specification
	spec, err := c.Describe()
	if err != nil {
		log.Fatalf("Failed to describe: %v\n", err)
	}

	fmt.Printf("Channel groups:\n")
	totalChannels := 0
	for i, group := range spec {
		fmt.Printf(" Group %d: channels=%v sampleRates=%v gains=%v\n", i, group.Channels, group.SupportedSampleRates, group.SupportedGains)
		totalChannels += len(group.Channels)
	}
	fmt.Printf(" Total channels: %d\n", totalChannels)

	// set all channels to same sampleRate
	configs := make([]*fspb.ChannelConfig, totalChannels)
	for i := 0; i < totalChannels; i++ {
		configs[i] = &fspb.ChannelConfig{
			Channel:    int32(i),
			SampleRate: float32(sampleRate),
			Gain:       1,
		}
	}

	if err := c.UploadConfiguration(anain.WithChannelConfig(configs)); err != nil {
		log.Fatalf("Failed to set configuration: %v\n", err)
	}

	// read back configuration
	config, err := c.DownloadConfiguration()
	if err != nil {
		log.Fatalf("Failed to download configuration: %v\n", err)
	}
	fmt.Printf("Current configuration:\n")
	for _, chConf := range config.ChannelConfig {
		fmt.Printf(" Channel %d: sampleRate=%.1f gain=%d\n", chConf.Channel, chConf.SampleRate, chConf.Gain)
	}

	// start stream
	err = c.StartStream(
		anain.WithFBStreamOption(functionblock.WithBucketSamples(400)),
		anain.WithFBStreamOption(functionblock.WithBufferedSamples(800)),
		anain.WithChannelMask(1<<uint32(totalChannels)-1),
	)
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	readStreamFor(c, time.Second*10)
}

func readStreamFor(c *anain.Client, duration time.Duration) {
	start := time.Now()

	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			fmt.Printf("got stream data seq=%d ts=%d\n", sd.Sequence, sd.DeliveryTimestamp)

			for i, sample := range samples {
				fmt.Printf("  #%d: ts=%d ch %d", i, sample.Timestamp, sample.BaseChannel)
				for _, value := range sample.Value {
					fmt.Printf(" %.4f", value)
				}
				fmt.Printf("\n")
			}
		}
		if time.Since(start) > duration {
			return
		}
	}
}
