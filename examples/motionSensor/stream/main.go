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
	"flag"
	"fmt"
	"os"
	"time"

	"log"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge-client-go/motionsensor"
)

func main() {
	const timeout = 0 // use default timeout

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <mdns-service-address OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	runtime := flag.Uint("runtime", 10, "Seconds to receive data")
	sampleRate := flag.Float64("samplerate", 12.5, "Sample rate in Hz")
	fullScale := flag.Uint("fullscale", 2, "Full scale in g")
	highPass := flag.Bool("hp", false, "Enable high pass filter")
	bandWithRatio := flag.Uint("bwr", 2, "Bandwith ratio")
	fileName := flag.String("file", "out.csv", "Output file name")
	dumpToConsole := flag.Bool("console", false, "Dump data to console")
	lowLatency := flag.Bool("lowlatency", false, "Use stream low latency mode")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	address := flag.Arg(0)

	// Create a client object to work with the io4edge device
	c, err := motionsensor.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create canl2 client: %v\n", err)
	}

	// set configuration
	if err := c.UploadConfiguration(
		motionsensor.WithSampleRate(uint32(*sampleRate*1000.0)),
		motionsensor.WithFullScale(int32(*fullScale)),
		motionsensor.WithHighPassFilterEnable(*highPass),
		motionsensor.WithBandWidthRatio(int32(*bandWithRatio))); err != nil {
		log.Fatalf("Failed to set configuration: %v\n", err)
	}

	// start stream
	err = c.StartStream(
		functionblock.WithBucketSamples(10),
		functionblock.WithBufferedSamples(200),
		functionblock.WithLowLatencyMode(*lowLatency),
	)
	if err != nil {
		log.Fatalf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	streamToCsv(c, *fileName, time.Second*time.Duration(*runtime), *dumpToConsole)
}

func streamToCsv(c *motionsensor.Client, fileName string, duration time.Duration, dumpToConsole bool) {
	start := time.Now()

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'

	defer w.Flush()

	prevTs := uint64(0)

	nSamples := 0
	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			fmt.Printf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			fmt.Printf("got stream data seq=%d with %d samples\n", sd.Sequence, len(samples))

			for _, sample := range samples {
				nSamples++
				record := []string{
					fmt.Sprintf("%d", sample.Timestamp),
					fmt.Sprintf("%.6f", sample.X),
					fmt.Sprintf("%.6f", sample.Y),
					fmt.Sprintf("%.6f", sample.Z),
				}

				if dumpToConsole {
					fmt.Printf("t: %15d dt: %7d x: %.6f, y: %.6f, z: %.6f\n", sample.Timestamp, sample.Timestamp-prevTs, sample.X, sample.Y, sample.Z)
				}

				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to file", err)
				}
				prevTs = sample.Timestamp
			}
			fmt.Printf("wrote %d samples\n", nSamples)
		}
		if time.Since(start) > duration {
			return
		}
	}
}
