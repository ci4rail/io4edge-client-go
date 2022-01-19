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

package cmd

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/core"
	"github.com/ci4rail/io4edge-client-go/internal/client"
	e "github.com/ci4rail/io4edge-client-go/internal/errors"

	"github.com/spf13/cobra"
)

var (
	chunkSize uint
)

var loadFirmwareCmd = &cobra.Command{
	Use:     "load-firmware FW_PKG",
	Aliases: []string{"load"},
	Short:   "Upload firmware package to device",
	Long: `Upload firmware package to device.
Example:
io4edge-cli load-firmware <firmware-package-file>`,
	Run:  loadFirmware,
	Args: cobra.ExactArgs(1),
}

func progressCb(bytes uint) {
	fmt.Printf("\r%d kBytes loaded.", bytes/1024)
}

func loadFirmware(cmd *cobra.Command, args []string) {
	file := args[0]
	c, err := client.NewCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	restartingNow, err := c.LoadFirmware(file, chunkSize, time.Duration(timeoutSecs)*time.Second, progressCb)
	e.ErrChk(err)

	readbackFirmwareID(c, restartingNow)
}

var loadRawFirmwareCmd = &cobra.Command{
	Use:     "load-raw-firmware FW_FILE",
	Aliases: []string{"load-raw"},
	Short:   "Upload firmware binary file to device",
	Long: `Upload firmware binary file to device.
NO COMPATIBILITY CHECKS! USE AT YOUR OWN RISK!
Example:
io4edge-client-go load-raw-firmware <firmware-file>`,
	Run:  loadRawFirmware,
	Args: cobra.ExactArgs(1),
}

func loadRawFirmware(cmd *cobra.Command, args []string) {
	file := args[0]
	c, err := client.NewCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	restartingNow, err := c.LoadFirmwareBinaryFromFile(file, chunkSize, time.Duration(timeoutSecs)*time.Second, progressCb)
	e.ErrChk(err)
	readbackFirmwareID(c, restartingNow)
}

func readbackFirmwareID(c *core.Client, restartingNow bool) {
	fmt.Println("\nFirmware load finished.")
	if restartingNow {
		fmt.Println("Let device restart...")
		// Must sleep here, because device needs approx. 1s to reset
		// If we reconnect too fast, we can still reach the device before it is reset
		time.Sleep(3 * time.Second)
		fmt.Println("Reconnect to restarted device...")
		var err error
		c, err = client.NewCliClient(deviceID, ipAddrPort)
		e.ErrChk(err)
	}
	fmt.Println("Reading back firmware id")
	fwName, fwVersion, err := c.IdentifyFirmware(time.Duration(timeoutSecs) * time.Second)
	e.ErrChk(err)
	fmt.Printf("Firmware name: %s, Version %s\n", fwName, fwVersion)
}

func init() {
	rootCmd.AddCommand(loadFirmwareCmd)
	rootCmd.AddCommand(loadRawFirmwareCmd)
	loadFirmwareCmd.PersistentFlags().UintVarP(&chunkSize, "chunksize", "c", 1024, "Size of chunk in bytes for loading")
	loadRawFirmwareCmd.PersistentFlags().UintVarP(&chunkSize, "chunksize", "c", 1024, "Size of chunk in bytes for loading")
}
