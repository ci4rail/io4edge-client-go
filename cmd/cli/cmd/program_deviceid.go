/*
Copyright © 2024 Ci4Rail GmbH <engineering@ci4rail.com>

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

	e "github.com/ci4rail/io4edge-client-go/internal/errors"
	"github.com/spf13/cobra"
)

var programDeviceIdentificationCmd = &cobra.Command{
	Use:     "program-devid INSTANCE_NAME",
	Aliases: []string{"devid"},
	Short:   "Program default mdns instance name",
	Long: `Program default mdns instance name in the flash of the device.
After a restart of the device, it will show up with this name in the mdns browser.
Example:
io4edge-cli -i 192.168.200.1:9999 program-devid S101-IOU04-USB-EXT-1`,
	Run:  programDeviceIdentification,
	Args: cobra.ExactArgs(1),
}

func programDeviceIdentification(cmd *cobra.Command, args []string) {
	name := "device-id"
	value := args[0]

	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	err = c.SetPersistentParameter(name, value, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	value, err = c.GetPersistentParameter(name, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	fmt.Printf("Device id was set to %s\n", value)
	fmt.Printf("Restart of the device required to apply the new device id.\n")
}

func init() {
	rootCmd.AddCommand(programDeviceIdentificationCmd)
}
