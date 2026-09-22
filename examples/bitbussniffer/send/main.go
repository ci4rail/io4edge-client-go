// Copyright © 2026 Ci4Rail GmbH <engineering@ci4rail.com>
// SPDX-License-Identifier: Apache-2.0

// Send transmits one frame through an already configured Bitbus sniffer.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/protobufcom/functionblockclients/bitbussniffer"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <mdns-service-name OR ip:port> <address> <control> [information-hex]\n", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() < 3 || flag.NArg() > 4 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	address, err := strconv.ParseUint(args[1], 0, 8)
	if err != nil {
		return fmt.Errorf("invalid address byte: %w", err)
	}
	control, err := strconv.ParseUint(args[2], 0, 8)
	if err != nil {
		return fmt.Errorf("invalid control byte: %w", err)
	}
	var information []byte
	if len(args) == 4 {
		compact := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, args[3])
		information, err = hex.DecodeString(compact)
		if err != nil {
			return fmt.Errorf("invalid information hex: %w", err)
		}
	}
	c, err := bitbussniffer.NewClientFromUniversalAddress(args[0], 0)
	if err != nil {
		return err
	}
	defer c.Close()
	// Configure prepare_sender using dumpstream first; preserve its active configuration.
	if err := c.SendFrame(append([]byte{byte(address), byte(control)}, information...)); err != nil {
		return err
	}
	fmt.Println("frame sent")
	return nil
}
