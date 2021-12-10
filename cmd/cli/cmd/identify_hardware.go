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

var identifyHardwareCmd = &cobra.Command{
	Use:     "identify-hardware",
	Aliases: []string{"id-hw", "hw"},
	Short:   "Get hardware infos from device",
	Run:     identifyHardware,
}

func identifyHardware(cmd *cobra.Command, args []string) {
	c, err := client.NewCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)
	identifyHardwareFromClient(c)
}

func identifyHardwareFromClient(c *core.Client) {
	rootArticle, majorVersion, serialNumber, err := c.IdentifyHardware(time.Duration(timeoutSecs) * time.Second)
	e.ErrChk(err)

	fmt.Printf("Hardware name: %s, rev: %d, serial: %s\n", rootArticle, majorVersion, serialNumber)
}

func init() {
	rootCmd.AddCommand(identifyHardwareCmd)
}