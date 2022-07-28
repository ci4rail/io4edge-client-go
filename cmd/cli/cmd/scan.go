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

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	e "github.com/ci4rail/io4edge-client-go/internal/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for io4edge devices",
	Run:   scan,
}

var (
	scanTime            uint
	enableShowFunctions bool
)

type device struct {
	core         *client.ServiceInfo
	functions    map[string]*client.ServiceInfo
	hardwareName string
	serial       string
}

type scanResults struct {
	devices map[string]*device // key: ip address as string
}

func (r *scanResults) numFoundIo4edgeDevices() int {
	n := 0
	for _, d := range r.devices {
		if d.core != nil {
			n++
		}
	}
	return n
}

func (r *scanResults) serviceAdded(s client.ServiceInfo) error {
	ip, _, _ := s.NetAddress()
	d, known := r.devices[ip]
	if !known {
		r.devices[ip] = &device{
			functions: make(map[string]*client.ServiceInfo, 0),
		}
		d = r.devices[ip]
	}
	if s.GetServiceType() == coreServiceType {
		d.core = &s
	} else {
		_, functionKnown := d.functions[s.GetInstanceName()]
		if !functionKnown {
			d.functions[s.GetInstanceName()] = &s
		}
	}

	return nil
}

func (r *scanResults) serviceRemoved(s client.ServiceInfo) error {
	// ignore removals for now
	return nil
}

func scan(cmd *cobra.Command, args []string) {
	result := &scanResults{
		make(map[string]*device, 0),
	}
	go func() {
		err := client.ServiceObserver("*", result.serviceAdded, result.serviceRemoved)
		e.ErrChk(err)
	}()
	time.Sleep(time.Duration(scanTime) * time.Second)

	if result.numFoundIo4edgeDevices() == 0 {
		fmt.Printf("No io4edge devices found\n")
	} else {
		devicesStr := "device"
		if result.numFoundIo4edgeDevices() > 1 {
			devicesStr += "s"
		}
		fmt.Printf("Found %d io4edge %s\n", result.numFoundIo4edgeDevices(), devicesStr)
		for _, d := range result.devices {

			if d.core != nil {
				c, err := newCliClient("", d.core.GetIPAddressPort())
				if err == nil {
					d.hardwareName, _, d.serial, _ = c.IdentifyHardware(time.Duration(timeoutSecs) * time.Second)
				}
				if d.hardwareName == "" {
					d.hardwareName = "(Unknown Hardware)"
				}
				if d.serial == "" {
					d.serial = "(Unknown Serial)"
				}
			}
		}

		if enableShowFunctions {
			for ip, d := range result.devices {

				if d.core != nil {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Service Type", "Service Name", "Port"})
					fmt.Printf("\n%s, %s, %s, %s\n", d.core.GetInstanceName(), ip, d.hardwareName, d.serial)
					for _, f := range d.functions {
						_, port, _ := f.NetAddress()
						table.Append([]string{f.GetServiceType(), f.GetInstanceName(), strconv.Itoa(port)})
					}
					table.Render() // Send output
				}
			}
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Device ID", "IP", "Hardware", "Serial"})
			for ip, d := range result.devices {

				if d.core != nil {
					table.Append([]string{d.core.GetInstanceName(), ip, d.hardwareName, d.serial})
				}
			}
			table.Render() // Send output
		}
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().UintVarP(&scanTime, "scantime", "", 2, "scan time in seconds")
	scanCmd.PersistentFlags().BoolVarP(&enableShowFunctions, "functions", "f", false, "show device sub functions")
}
