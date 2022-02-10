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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	anain "github.com/ci4rail/io4edge-client-go/analogintypea"
	"github.com/ci4rail/io4edge-client-go/functionblock"
)

func streamToCsv(c *anain.Client, fileName string, duration time.Duration) {
	start := time.Now()

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'

	defer w.Flush()

	nSamples := 0
	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			//fmt.Printf("got stream data seq=%d ts=%d\n", sd.Sequence, sd.DeliveryTimestamp)

			for _, sample := range samples {
				nSamples++
				record := []string{
					fmt.Sprintf("%d", sample.Timestamp),
					fmt.Sprintf("%.4f", sample.Value),
				}
				//record[1] = strings.Replace(record[1], ".", ",", 1)

				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to file", err)
				}
			}
			fmt.Printf("wrote %d samples\n", nSamples)
		}
		if time.Since(start) > duration {
			return
		}
	}
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port> <csv-file> <sample-rate>", os.Args[0])
	}
	address := os.Args[1]

	// Create a client object to work with the io4edge device
	var c *anain.Client
	var err error

	fileName := os.Args[2]
	sampleRate, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Can't convert sample rate: %v\n", err)
	}

	c, err = anain.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create anain client: %v\n", err)
	}

	// set sampleRate
	if err := c.UploadConfiguration(
		&anain.Configuration{
			SampleRate: uint32(sampleRate),
		}); err != nil {
		log.Fatalf("Failed to set configuration: %v\n", err)
	}

	// start stream. Trigger on changes of all channels
	err = c.StartStream(&functionblock.StreamConfiguration{
		BucketSamples:     100,
		BufferedSamples:   200,
		KeepaliveInterval: 1000,
	})
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	streamToCsv(c, fileName, time.Second*10)
}
