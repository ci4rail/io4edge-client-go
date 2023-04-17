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
	"sort"
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
	functions    map[string]*client.ServiceInfo // key: instance name
	ip           string
	hardwareName string
	serial       string
}

type byIP []*device
type byPort []*client.ServiceInfo

type scanResults struct {
	devices     map[string]*device // key: ip address as string
	scanRunning bool
}

func scan(cmd *cobra.Command, args []string) {
	result := &scanResults{
		make(map[string]*device, 0),
		true,
	}
	go func() {
		err := client.ServiceObserver("*", result.serviceAdded, result.serviceRemoved)
		e.ErrChk(err)
	}()
	time.Sleep(time.Duration(scanTime) * time.Second)
	result.scanRunning = false

	devices := result.sortDevicesByIP()
	addDevicesHwInfo(devices)

	if len(devices) == 0 {
		fmt.Printf("No io4edge devices found\n")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		setCommonTableOptions(table)

		if enableShowFunctions {
			table.SetHeader([]string{"Device ID", "Service Type", "Service Name", "IP:Port"})
			for _, d := range devices {

				// output core service info
				_, port, _ := d.core.NetAddress()
				table.Append([]string{d.core.GetInstanceName(), d.core.GetServiceType(), d.core.GetInstanceName(), fmt.Sprintf("%s:%d", d.ip, port)})

				functions := sortServicesByPort(d.functions)
				for _, f := range functions {
					_, port, _ := f.NetAddress()
					table.Append([]string{"", f.GetServiceType(), f.GetInstanceName(), fmt.Sprintf("%s:%d", d.ip, port)})
				}
			}
		} else {
			table.SetHeader([]string{"Device ID", "IP", "Hardware", "Serial"})
			for _, d := range devices {
				table.Append([]string{d.core.GetInstanceName(), d.ip, d.hardwareName, d.serial})
			}
		}
		table.Render() // Send output
	}
}

func setCommonTableOptions(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
}

func addDevicesHwInfo(devices []*device) {
	for _, d := range devices {

		c, err := newCliClient("", d.core.GetIPAddressPort())
		if err == nil {
			d.hardwareName, _, d.serial, _ = c.IdentifyHardware(time.Duration(timeoutSecs) * time.Second)
		}
		if d.hardwareName == "" {
			d.hardwareName = "(Unknown)"
		}
		if d.serial == "" {
			d.serial = "(Unknown)"
		}
	}
}

func (a byIP) Len() int           { return len(a) }
func (a byIP) Less(i, j int) bool { return a[i].core.GetIPAddressPort() < a[j].core.GetIPAddressPort() }
func (a byIP) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// sortDevicesByIP returns a slice of devices sorted by IP address, only devices with a core service are included
func (r *scanResults) sortDevicesByIP() []*device {
	devices := make([]*device, 0)
	for _, d := range r.devices {
		if d.core != nil {
			devices = append(devices, d)
		}
	}
	// sort devices according to IP address
	sort.Sort(byIP(devices))

	return devices
}

func (r *scanResults) serviceAdded(s client.ServiceInfo) error {
	if r.scanRunning {
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
			d.ip = ip
		} else {
			_, functionKnown := d.functions[s.GetInstanceName()]
			if !functionKnown {
				d.functions[s.GetInstanceName()] = &s
			}
		}
	}
	return nil
}

func (r *scanResults) serviceRemoved(s client.ServiceInfo) error {
	// ignore removals for now
	return nil
}

func (a byPort) Len() int { return len(a) }

func (a byPort) Less(i, j int) bool {
	_, portA, _ := a[i].NetAddress()
	_, portB, _ := a[j].NetAddress()
	return portA < portB
}

func (a byPort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func sortServicesByPort(svc map[string]*client.ServiceInfo) []*client.ServiceInfo {
	services := make([]*client.ServiceInfo, 0)
	for _, s := range svc {
		services = append(services, s)
	}
	sort.Sort(byPort(services))
	return services
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().UintVarP(&scanTime, "scantime", "", 2, "scan time in seconds")
	scanCmd.PersistentFlags().BoolVarP(&enableShowFunctions, "functions", "f", false, "show device sub functions")
}
