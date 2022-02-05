package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge-client-go/templatemodule"
)

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 3 {
		log.Fatalf("Usage: identify svc <mdns-service-address>  OR  identify ip <ip:port>")
	}
	addressType := os.Args[1]
	address := os.Args[2]

	// Create a client object to work with the io4edge device at <address>
	var c *templatemodule.Client
	var fbClient *functionblock.Client
	var err error

	if addressType == "svc" {
		fbClient, err = functionblock.NewClientFromService(address, timeout)
	} else {
		fbClient, err = functionblock.NewClientFromSocketAddress(address)
	}
	if err != nil {
		log.Fatalf("Failed to create functionblock client: %v\n", err)
	}
	c = templatemodule.NewClient(fbClient)

	err = c.ConfigurationSet(&templatemodule.Configuration{SampleRate: 8000})
	if err != nil {
		log.Errorf("ConfigurationSet failed: %v\n", err)
	}

	// provoke error
	err = c.ConfigurationSet(&templatemodule.Configuration{SampleRate: 100000})
	if err != nil {
		log.Errorf("ConfigurationSet failed: %v\n", err)
	}

	cfg, err := c.ConfigurationGet()
	if err != nil {
		log.Errorf("ConfigurationGet failed: %v\n", err)
	}
	fmt.Printf("Configuration: %v\n", cfg)

	desc, err := c.ConfigurationDescribe()
	if err != nil {
		log.Errorf("ConfigurationDescribe failed: %v\n", err)
	}
	fmt.Printf("Description: %v\n", desc)

	err = c.SetCounter(1234)
	if err != nil {
		log.Errorf("SetCounter failed: %v\n", err)
	}

	cnt, err := c.GetCounter()
	if err != nil {
		log.Errorf("SetCounter failed: %v\n", err)
	}
	fmt.Printf("counter: %d\n", cnt)

	err = c.StartStream(&functionblock.StreamConfiguration{
		BucketSamples:     400,
		BufferedSamples:   1000,
		KeepaliveInterval: 2000,
	}, 4)
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	for {
		sd, err := fbClient.ReadStreamData(time.Second * 5)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			fmt.Printf("got stream data %d\n", sd.Sequence)
		}
	}
}
