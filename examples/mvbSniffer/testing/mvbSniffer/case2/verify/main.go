package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/ci4rail/io4edge-client-go/mvbsniffer"
	fspb "github.com/ci4rail/io4edge_api/mvbSniffer/go/mvbSniffer/v1"
)

func readStreamFor(c *mvbsniffer.Client, duration time.Duration) {
	const (
		StInit = 0
		StFrm2 = 1
		StFrm3 = 2
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
			samples := sd.FSData.GetEntry()
			//fmt.Printf("got stream data seq=%d ts=%d\n", sd.Sequence, sd.DeliveryTimestamp)

			for _, sample := range samples {
				//fmt.Printf("st=%d #%d: %v\n", state, i, sample)
				if sample.State != uint32(fspb.Telegram_kSuccessful) {
					log.Errorf("#%d: %v\n", n, sample)
				}
				switch state {
				case StInit:
					switch sample.Address {
					case 2:
						state = StFrm3
					case 3:
						state = StFrm2
					default:
						log.Errorf("#%d Bad address received %d", n, sample.Address)
					}
				case StFrm2:
					if sample.Address != 2 {
						log.Errorf("#%d FRM2 Bad address received %d", n, sample.Address)
					} else {

						dt := sample.Timestamp - prevTs
						if dt < 2000 || dt > 2200 {
							log.Errorf("#%d FRM2 wrong dt %d (%v/%v)", n, dt, sample.Timestamp, prevTs)
						}
						if !bytes.Equal(sample.Data, []uint8{0x00, 0x00}) {
							log.Errorf("#%d FRM2 wrong bytes %v", n, sample.Data)
						}
						if sample.Line != fspb.Telegram_kA {
							log.Errorf("#%d FRM2 wrong line %v", n, sample.Line)
						}

						state = StFrm3
					}
				case StFrm3:
					if sample.Address != 3 {
						log.Errorf("#%d FRM3 Bad address received %d", n, sample.Address)
					} else {
						dt := sample.Timestamp - prevTs
						if dt < 1000 || dt > 1050 {
							log.Errorf("#%d FRM3 wrong dt %d (%v/%v)", n, dt, sample.Timestamp, prevTs)
						}
						if !bytes.Equal(sample.Data, []uint8{0x00, 0x00}) {
							log.Errorf("#%d FRM3 wrong bytes %v", n, sample.Data)
						}
						if sample.Line != fspb.Telegram_kB {
							log.Errorf("#%d FRM3 wrong line %v", n, sample.Line)
						}
						state = StFrm2
					}
				}
				prevTs = sample.Timestamp
				n++
			}
		}
		if time.Since(start) > duration {
			fmt.Printf("%d frames received\n", n)
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

	c, err := mvbsniffer.NewClientFromUniversalAddress(address, timeout)
	if err != nil {
		log.Fatalf("Failed to create mvbsniffer client: %v\n", err)
	}

	// start stream
	err = c.StartStream(&functionblock.StreamConfiguration{
		BucketSamples:     100,
		BufferedSamples:   200,
		KeepaliveInterval: 1000,
	}, mvbsniffer.StreamFilter{
		Masks: []mvbsniffer.FilterMask{
			{FCodeMask: 0x0001, Address: 0x0000, Mask: 0xFFC},
		},
	})
	if err != nil {
		log.Errorf("StartStream failed: %v\n", err)
	}

	fmt.Println("Started stream")

	readStreamFor(c, time.Second*10)
}
