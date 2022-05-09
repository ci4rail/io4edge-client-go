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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	anain "github.com/ci4rail/io4edge-client-go/analogintypea"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	fspb "github.com/ci4rail/io4edge_api/analogInTypeA/go/analogInTypeA/v1alpha1"
)

type entry struct {
	channel int
	sample  fspb.Sample
}

func readStream(c *anain.Client, channel int, quitChan chan bool, sampleChan chan entry, wg *sync.WaitGroup) {
	for {
		select {
		case <-quitChan:
			wg.Done()
			return
		default:
		}

		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			//fmt.Printf("got stream data seq=%d ts=%d\n", sd.Sequence, sd.DeliveryTimestamp)

			for _, sample := range samples {
				sampleChan <- entry{
					channel: channel,
					sample:  *sample,
				}
			}
		}
	}
}

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 6 {
		log.Fatalf("Usage: %s <channel1-address> <channel2-address> <csv-file> <sample-rate> <runtime>", os.Args[0])
	}
	address := [2]string{os.Args[1], os.Args[2]}

	fileName := os.Args[3]
	sampleRate, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatalf("Can't convert sample rate: %v\n", err)
	}
	runtime, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Fatalf("Can't convert runtime: %v\n", err)
	}

	// Write CSV
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'

	defer w.Flush()
	sampleChan := make(chan entry)
	quitChan := make(chan bool)
	var wg sync.WaitGroup = sync.WaitGroup{}

	for channel := 0; channel < 2; channel++ {

		c, err := anain.NewClientFromUniversalAddress(address[channel], timeout)
		if err != nil {
			log.Fatalf("Failed to create anain client: %v\n", err)
		}

		// set sampleRate
		if err := c.UploadConfiguration(anain.WithSampleRate(uint32(sampleRate))); err != nil {
			log.Fatalf("Failed to set configuration: %v\n", err)
		}

		// start stream
		err = c.StartStream(
			functionblock.WithBucketSamples(100),
			functionblock.WithBufferedSamples(200),
		)
		if err != nil {
			log.Errorf("StartStream failed: %v\n", err)
		}

		fmt.Printf("Started stream on channel %d\n", channel)

		go readStream(c, channel, quitChan, sampleChan, &wg)
		wg.Add(1)
	}
	start := time.Now()
	numSamples := 0
	for {
		entry := <-sampleChan
		record := []string{
			fmt.Sprintf("%d", entry.channel),
			fmt.Sprintf("%d", entry.sample.Timestamp),
			fmt.Sprintf("%.4f", entry.sample.Value),
		}
		if time.Since(start) > time.Second*time.Duration(runtime) {
			break
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		numSamples++
	}
	fmt.Printf("Wrote %d samples\n", numSamples)
}
