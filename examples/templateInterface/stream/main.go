package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/common/functionblock"
	"github.com/ci4rail/io4edge-client-go/pkg/protobufcom/functionblockclients/templateinterface"
)

func main() {
	const timeout = 0 // use default timeout

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address>  OR  %s <ip:port>", os.Args[0], os.Args[0])
	}
	address := os.Args[1]
	//log.SetLevel(log.DebugLevel)

	// Create a client object to work with the io4edge device at <address>
	var c *templateinterface.Client
	var err error

	c, err = templateinterface.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create templateinterface client: %v\n", err)
	}

	err = c.UploadConfiguration(templateinterface.WithSampleRate(40))
	if err != nil {
		log.Errorf("ConfigurationSet failed: %v\n", err)
	}

	// provoke error
	err = c.UploadConfiguration(templateinterface.WithSampleRate(100000))
	if err != nil {
		log.Errorf("ConfigurationSet failed: %v\n", err)
	}

	cfg, err := c.DownloadConfiguration()
	if err != nil {
		log.Errorf("ConfigurationGet failed: %v\n", err)
	}
	fmt.Printf("Configuration: %v\n", cfg)

	desc, err := c.Describe()
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

	err = c.StartStream(
		templateinterface.WithModulo(4),
		templateinterface.WithFBStreamOption(functionblock.WithBucketSamples(40)),
		templateinterface.WithFBStreamOption(functionblock.WithBufferedSamples(1000)),
		templateinterface.WithFBStreamOption(functionblock.WithKeepaliveInterval(2000)),
	)
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	for {
		sd, err := c.ReadStream(time.Second * 5)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {

			samples := sd.FSData.GetSamples()
			fmt.Printf("got stream data %d samples: %v\n", sd.Sequence, samples)
		}
	}
}
