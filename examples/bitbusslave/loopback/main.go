// Copyright © 2026 Ci4Rail GmbH <engineering@ci4rail.com>
// SPDX-License-Identifier: Apache-2.0

// Loopback exercises a slave and simulated master on one device without wiring.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/common/functionblock"
	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/bitbusslave"
	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/bitbussniffer"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <slave-endpoint> <sniffer-endpoint>\n", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, flag.Arg(0), flag.Arg(1)); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// exchange skips exactly one local echo: an RR reply may equal the command.
func exchange(sniffer *bitbussniffer.Client, control byte, information []byte) ([]byte, error) {
	command := append([]byte{1, control}, information...)
	if err := sniffer.SendFrame(command); err != nil {
		return nil, err
	}
	deadline := time.Now().Add(3 * time.Second)
	echoSeen := false
	for time.Now().Before(deadline) {
		sd, err := sniffer.ReadStream(time.Until(deadline))
		if err != nil {
			return nil, fmt.Errorf("await slave response: %w", err)
		}
		for _, sample := range sd.FSData.Samples {
			if sample.Flags != 0 {
				return nil, fmt.Errorf("sniffer flags: 0x%x", sample.Flags)
			}
			frame := sample.BitbusFrame
			if bytes.Equal(frame, command) && !echoSeen {
				echoSeen = true
				continue
			}
			if len(frame) >= 2 && frame[0] == 1 {
				return frame, nil
			}
		}
	}
	return nil, fmt.Errorf("no slave response")
}

func run(ctx context.Context, slaveAddress, snifferAddress string) error {
	sniffer, err := bitbussniffer.NewClientFromUniversalAddress(snifferAddress, 0)
	if err != nil {
		return err
	}
	defer sniffer.Close()
	slave, err := bitbusslave.NewClientFromUniversalAddress(slaveAddress, 0)
	if err != nil {
		return err
	}
	defer slave.Close()
	if err := configureLoopback(sniffer, slave); err != nil {
		return err
	}
	opts := []functionblock.StreamConfigOption{
		functionblock.WithBucketSamples(1), functionblock.WithBufferedSamples(64),
		functionblock.WithKeepaliveInterval(1000), functionblock.WithLowLatencyMode(true),
	}
	if err := sniffer.StartStream(opts...); err != nil {
		return err
	}
	defer sniffer.StopStream()
	if err := slave.StartStream(opts...); err != nil {
		return err
	}
	defer slave.StopStream()

	if err := handshake(sniffer); err != nil {
		return err
	}
	return runMessages(ctx, sniffer, slave)
}

func configureLoopback(sniffer *bitbussniffer.Client, slave *bitbusslave.Client) error {
	if err := sniffer.UploadConfiguration(
		bitbussniffer.WithBaud62500(false),
		bitbussniffer.WithAddressFilter(bytes.Repeat([]byte{0xff}, 32)),
		bitbussniffer.WithMinFrameLength(0),
		bitbussniffer.WithPrepareSender(true),
		bitbussniffer.WithLoopbackEnable(true),
		bitbussniffer.WithFullDuplex(true),
	); err != nil {
		return fmt.Errorf("configure sniffer: %w", err)
	}
	if err := slave.UploadConfiguration(
		bitbusslave.WithSlaveAddress(1),
		bitbusslave.WithBaud62500(false),
		bitbusslave.WithAppWDTimeoutMS(5000),
		bitbusslave.WithIdleResponse([]byte{0}),
	); err != nil {
		return fmt.Errorf("configure slave: %w", err)
	}
	return nil
}

func handshake(sniffer *bitbussniffer.Client) error {
	// DISC can initially receive FRMR if a previous session left the slave in NRM.
	response, err := exchange(sniffer, 0x53, nil)
	if err != nil {
		return err
	}
	if bytes.Equal(response, []byte{1, 0x97}) {
		response, err = exchange(sniffer, 0x53, nil)
		if err != nil {
			return err
		}
	}
	if !bytes.Equal(response, []byte{1, 0x73}) {
		return fmt.Errorf("DISC failed: %x", response)
	}
	response, err = exchange(sniffer, 0x93, nil)
	if err != nil {
		return err
	}
	if !bytes.Equal(response, []byte{1, 0x73}) {
		return fmt.Errorf("SNRM failed: %x", response)
	}
	return nil
}

func runMessages(ctx context.Context, sniffer *bitbussniffer.Client, slave *bitbusslave.Client) error {
	master := masterSession{sniffer: sniffer}
	nextSlave, nextMaster, nextStatus := time.Now(), time.Now(), time.Now()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for ctx.Err() == nil {
		now := time.Now()
		if !now.Before(nextSlave) {
			queued, err := queueSlaveMessage(slave)
			if err != nil {
				return err
			}
			if queued {
				nextSlave = now.Add(200 * time.Millisecond)
			}
		}
		sending := !now.Before(nextMaster)
		if err := master.poll(sending); err != nil {
			return err
		}
		if sending {
			nextMaster = now.Add(500 * time.Millisecond)
		}
		if err := printSlaveMessages(slave); err != nil {
			return err
		}
		if !now.Before(nextStatus) {
			if err := printSlaveStatus(slave); err != nil {
				return err
			}
			nextStatus = now.Add(time.Second)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
	return nil
}

func queueSlaveMessage(slave *bitbusslave.Client) (bool, error) {
	state, err := slave.State()
	if err != nil {
		return false, err
	}
	if state.HavePendingTxMsg {
		return false, nil
	}
	if err := slave.SetPreparedTxMsg([]byte("slave message")); err != nil {
		return false, err
	}
	return true, nil
}

// masterSession tracks the sequence numbers across master polls.
type masterSession struct {
	sniffer *bitbussniffer.Client
	ns      byte // Next master transmit sequence.
	nr      byte // Expected slave receive sequence.
}

func (m *masterSession) poll(sending bool) error {
	control := (m.nr << 5) | 0x11 // RR poll.
	var information []byte
	if sending {
		control = (m.nr << 5) | 0x10 | (m.ns << 1)
		information = []byte("master message")
	}
	response, err := exchange(m.sniffer, control, information)
	if err != nil {
		return err
	}
	if sending {
		m.ns = (m.ns + 1) % 8
	}
	return m.handleResponse(response)
}

func (m *masterSession) handleResponse(response []byte) error {
	control := response[1]
	if control&1 == 0 {
		if (control>>1)&7 != m.nr {
			return fmt.Errorf("unexpected slave sequence: %x", response)
		}
		m.nr = (m.nr + 1) % 8
		fmt.Printf("Master received: %q\n", response[2:])
	} else if control&0x0f != 0x01 {
		return fmt.Errorf("unexpected slave response: %x", response)
	}
	if (control>>5)&7 != m.ns {
		return fmt.Errorf("slave did not acknowledge master message: %x", response)
	}
	if control&1 == 0 {
		return m.acknowledge()
	}
	return nil
}

func (m *masterSession) acknowledge() error {
	// ACK before queuing another message. RNR prevents another idle I frame.
	response, err := exchange(m.sniffer, (m.nr<<5)|0x15, nil)
	if err != nil {
		return err
	}
	if !bytes.Equal(response, []byte{1, (m.ns << 5) | 0x11}) {
		return fmt.Errorf("unexpected RNR response: %x", response)
	}
	return nil
}

func printSlaveMessages(slave *bitbusslave.Client) error {
	// A short positive timeout drains queued samples without a zero-timeout select race.
	for {
		sd, err := slave.ReadStream(time.Millisecond)
		if err != nil {
			// The function-block API currently returns an untyped timeout error.
			if err.Error() == "timeout waiting for stream data" {
				return nil
			}
			return err
		}
		for _, sample := range sd.FSData.Samples {
			fmt.Printf("Slave received: %q\n", sample.BitbusInformation)
		}
	}
}

func printSlaveStatus(slave *bitbusslave.Client) error {
	state, err := slave.State()
	if err != nil {
		return err
	}
	fmt.Printf("Slave status: %s, pending TX=%t\n", state.Mode, state.HavePendingTxMsg)
	return nil
}
