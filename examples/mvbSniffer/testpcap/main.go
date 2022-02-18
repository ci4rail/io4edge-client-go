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
	"log"
	"os"
	"time"

	"github.com/ci4rail/io4edge-client-go/examples/mvbSniffer/pcap/busshark"
	"github.com/ci4rail/io4edge-client-go/examples/mvbSniffer/pcap/pcap"
)

func main() {
	f, err := os.Create("dat.pcap")
	if err != nil {
		log.Fatalf("Can't create file %s", err)
	}
	defer f.Close()

	w, err := pcap.NewWriter(f)
	if err != nil {
		log.Fatalf("Can't create writer %s", err)
	}

	m := busshark.Pkt(1, 50, 1, 1, []byte{0x41, 0x8c})
	w.AddPacket(uint64(time.Now().UnixMicro()), m)

	m = busshark.Pkt(1, 80, 1, 2, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x21, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x31, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x41, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})
	w.AddPacket(uint64(time.Now().UnixMicro()), m)

}
