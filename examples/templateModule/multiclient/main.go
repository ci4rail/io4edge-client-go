package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge-client-go/templatemodule"
)

const (
	timeout = 5 * time.Second
)

func client(clientNum int, duration time.Duration, address string, wg *sync.WaitGroup) {
	numSamples := 0

	log.Printf("Client %d starting", clientNum)
	c, err := templatemodule.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("%d: Failed to create templateModule client: %v\n", clientNum, err)
	}
	log.Printf("Client %d connected", clientNum)

	err = c.UploadConfiguration(templatemodule.WithSampleRate(100))
	if err != nil {
		log.Errorf("%d: ConfigurationSet failed: %v\n", clientNum, err)
	}

	err = c.StartStream(
		templatemodule.WithModulo(1),
		templatemodule.WithFBStreamOption(functionblock.WithBucketSamples(30)),
		templatemodule.WithFBStreamOption(functionblock.WithBufferedSamples(50)),
		templatemodule.WithFBStreamOption(functionblock.WithKeepaliveInterval(2000)),
	)
	if err != nil {
		log.Errorf("%d: StartStream failed: %v\n", clientNum, err)
	}
	start := time.Now()

	for {
		sd, err := c.ReadStream(time.Second * 5)

		if err != nil {
			log.Errorf("%d: ReadStreamData failed: %v\n", clientNum, err)
		} else {

			samples := sd.FSData.GetSamples()
			fmt.Printf("%d: got stream data %d\n", clientNum, sd.Sequence)
			numSamples += len(samples)
			//fmt.Printf("got stream data %d samples: %v\n", sd.Sequence, samples)
		}
		if time.Since(start) > duration {
			break
		}
	}
	if err := c.StopStream(); err != nil {
		log.Errorf("%d: StopStream failed: %v\n", clientNum, err)
	}
	log.Printf("Client %d stopped. Got %d samples", clientNum, numSamples)
	wg.Done()
}

func testNConnections(address string) {
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		// this delay is needed because we have a limit in the maximum number of simultaneous connection attempts
		time.Sleep(100 * time.Millisecond)
		go client(i, 10*time.Second, address, &wg)
	}
	wg.Wait()
}

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR ip:port>", os.Args[0])
	}
	address := os.Args[1]
	//log.SetLevel(log.DebugLevel)

	testNConnections(address)
}
