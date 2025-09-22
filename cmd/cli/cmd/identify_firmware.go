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

var identifyFirmwareCmd = &cobra.Command{
	Use:     "identify-firmware",
	Aliases: []string{"id-fw", "fw"},
	Short:   "Get firmware infos from device",
	Run:     identifyFirmware,
}

func identifyFirmware(cmd *cobra.Command, args []string) {
	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)
	fwName, fwVersion, err := c.IdentifyFirmware(time.Duration(timeoutSecs) * time.Second)
	e.ErrChk(err)
	fmt.Printf("Firmware name: %s, Version %s\n", fwName, fwVersion)
}

func init() {
	rootCmd.AddCommand(identifyFirmwareCmd)
}
