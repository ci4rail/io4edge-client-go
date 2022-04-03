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

var scanTime uint

type scanResults struct {
	devices []client.ServiceInfo
}

func (r *scanResults) serviceAdded(s client.ServiceInfo) error {
	r.devices = append(r.devices, s)
	return nil
}

func (r *scanResults) serviceRemoved(s client.ServiceInfo) error {
	for i, t := range r.devices {
		if s.GetInstanceName() == t.GetInstanceName() {
			r.devices = append(r.devices[:i], r.devices[i+1:]...)
		}
	}
	return nil
}

func scan(cmd *cobra.Command, args []string) {
	result := &scanResults{
		make([]client.ServiceInfo, 0),
	}
	go func() {
		err := client.ServiceObserver(coreServiceType, result.serviceAdded, result.serviceRemoved)
		e.ErrChk(err)
	}()
	time.Sleep(time.Duration(scanTime) * time.Second)

	if len(result.devices) == 0 {
		fmt.Printf("No io4edge devices found\n")
	} else {
		fmt.Printf("Found %d io4edge devices\n", len(result.devices))
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Device ID", "IP", "Hardware"})

		for _, t := range result.devices {
			ip, _, _ := t.NetAddress()
			var rootArticle string
			c, err := newCliClient("", t.GetIPAddressPort())
			if err == nil {
				rootArticle, _, _, _ = c.IdentifyHardware(time.Duration(timeoutSecs) * time.Second)
			}
			if rootArticle == "" {
				rootArticle = "(Unknown)"
			}
			table.Append([]string{
				t.GetInstanceName(),
				ip,
				rootArticle,
			})
		}
		table.Render() // Send output
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().UintVarP(&scanTime, "scantime", "", 2, "scan time in seconds")
}
