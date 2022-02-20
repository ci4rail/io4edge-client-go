package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge-client-go/mvbsniffer"
)

func generatePattern() string {
	cl := mvbsniffer.NewCommandList()

	cl.AddMasterFrame(0, false, 5, 0x1, 567)
	cl.AddMasterFrame(0, false, 3, 0x1, 888)

	cl.AddMasterFrame(0, false, 5, 0x1, 123)
	cl.AddSlaveFrame(0, false, 2, []uint8{0xaa, 0xbb, 0xcc, 0xdd})
	cl.AddMasterFrame(0, false, 5, 0x0, 200)
	cl.AddSlaveFrame(0, false, 2, []uint8{0x12, 0x34})
	cl.AddMasterFrame(0, false, 5, 0x1, 567)
	cl.AddMasterFrame(0, false, 400, 0x1, 888)

	return cl.StartGeneratorString(true)
}

func readStreamFor(c *mvbsniffer.Client, duration time.Duration) {
	const (
		StInit   = 0
		StFrm123 = 1
		StFrm200 = 2
	)

	state := StInit
	start := time.Now()
	prevTs := uint64(0)
	n := 0
	for {
		// read next bucket from stream
		sd, err := c.ReadStream(time.Second * 1)

		if err != nil {
			log.Errorf("ReadStreamData failed: %v\n", err)
		} else {
			samples := sd.FSData.GetSamples()
			//fmt.Printf("got stream data seq=%d ts=%d\n", sd.Sequence, sd.DeliveryTimestamp)

			for _, sample := range samples {
				//fmt.Printf("st=%d #%d: %v\n", state, i, sample)

				switch state {
				case StInit:
					switch sample.Address {
					case 123:
						state = StFrm200
					case 200:
						state = StFrm123
					default:
						log.Errorf("#%d Bad address received %d", n, sample.Address)
					}
				case StFrm123:
					if sample.Address != 123 {
						log.Errorf("#%d FRM123 Bad address received %d", n, sample.Address)
					} else {

						dt := sample.Timestamp - prevTs
						if dt < 410 || dt > 580 {
							log.Errorf("#%d FRM123 wrong dt %d (%v/%v)", n, dt, sample.Timestamp, prevTs)
						}
						state = StFrm200
					}
				case StFrm200:
					if sample.Address != 200 {
						log.Errorf("#%d FRM200 Bad address received %d", n, sample.Address)
					} else {
						dt := sample.Timestamp - prevTs
						if dt < 7 || dt > 60 {
							log.Errorf("#%d FRM200 wrong dt %d (%v/%v)", n, dt, sample.Timestamp, prevTs)
						}
						state = StFrm123
					}
				}
				prevTs = sample.Timestamp
				n++
			}
		}
		if time.Since(start) > duration {
			return
		}
	}
}

func main() {
	const timeout = 5 * time.Second

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <mdns-service-address OR <ip:port>>", os.Args[0])
	}
	address := os.Args[1]

	pattern := generatePattern()

	fmt.Printf("Generator pattern: %s\n", pattern)
	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create mvbsniffer client: %v\n", err)
	}

	// start stream
	err = c.StartStream(&functionblock.StreamConfiguration{
		BucketSamples:     100,
		BufferedSamples:   200,
		KeepaliveInterval: 1000,
	})
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	err = c.SendPattern(pattern)
	if err != nil {
		log.Errorf("SendPattern failed: %v\n", err)
	}

	readStreamFor(c, time.Second*10)
}
