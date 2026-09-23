// Copyright © 2026 Ci4Rail GmbH <engineering@ci4rail.com>
// SPDX-License-Identifier: Apache-2.0

// Dumpstream prints captured Bitbus frames and their status flags.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/bitbussniffer"
	pb "github.com/ci4rail/io4edge_api/bitbusSniffer/go/bitbusSniffer/v1"
)

func main() {
	lowLatency := flag.Bool("lowlatency", false, "Use stream low latency mode")
	prepareSender := flag.Bool("prepare_sender", false, "Enable sending frames")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] <mdns-service-name OR ip:port>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, flag.Arg(0), *lowLatency, *prepareSender); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, address string, lowLatency, prepareSender bool) error {
	c, err := bitbussniffer.NewClientFromUniversalAddress(address, 0)
	if err != nil {
		return err
	}
	defer c.Close()
	filter := make([]byte, 32)
	filter[0] = 0x07 // Capture addresses 0, 1 and 2.
	if err := c.UploadConfiguration(
		bitbussniffer.WithIgnoreCRC(true),
		bitbussniffer.WithBaud62500(false),
		bitbussniffer.WithAddressFilter(filter),
		bitbussniffer.WithMinFrameLength(1),
		bitbussniffer.WithPrepareSender(prepareSender),
	); err != nil {
		return fmt.Errorf("configure sniffer: %w", err)
	}
	if err := c.StartStream(
		functionblock.WithBucketSamples(20),
		functionblock.WithBufferedSamples(60),
		functionblock.WithKeepaliveInterval(3000),
		functionblock.WithLowLatencyMode(lowLatency),
	); err != nil {
		return err
	}
	defer c.StopStream()
	for ctx.Err() == nil {
		sd, err := c.ReadStream(3 * time.Second)
		if err != nil {
			// The function-block API currently returns an untyped timeout error.
			if err.Error() == "timeout waiting for stream data" {
				fmt.Println("Timeout while reading stream")
				continue
			}
			return err
		}
		fmt.Printf("%s: Received %d samples, seq=%d, ts=%.6f\n",
			time.Now().Format(time.RFC3339Nano), len(sd.FSData.Samples), sd.Sequence,
			float64(sd.DeliveryTimestamp)/1e6)
		for _, sample := range sd.FSData.Samples {
			fmt.Println(sampleString(sample))
		}
	}
	return nil
}

func sampleString(sample *pb.Sample) string {
	s := fmt.Sprintf("%10d us: ", sample.Timestamp)
	frame := sample.BitbusFrame
	if len(frame) >= 2 {
		s += fmt.Sprintf("ADDR: 0x%02X CTRL: 0x%02X INFO: % X", frame[0], frame[1], frame[2:])
	} else {
		s += fmt.Sprintf("SHORT FRAME: % X", frame)
	}
	if sample.Flags&uint32(pb.Sample_bad_crc) != 0 {
		s += " CRC_ERR"
	}
	if sample.Flags&uint32(pb.Sample_frames_lost) != 0 {
		s += " LOST"
	}
	if sample.Flags&uint32(pb.Sample_buf_overrun) != 0 {
		s += " BUF_OVERRUN"
	}
	return s
}
