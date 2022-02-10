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
	"strconv"
	"time"

	"github.com/ci4rail/io4edge-client-go/internal/client"
	e "github.com/ci4rail/io4edge-client-go/internal/errors"

	"github.com/spf13/cobra"
)

var programHardwareIdentificationCmd = &cobra.Command{
	Use:     "program-hwid NAME MAJOR SERIAL",
	Aliases: []string{"hwid"},
	Short:   "Program new HW ID",
	Long: `Program new HW ID into device.
Example:
io4edge-cli program-hwid S101-IOU04 1 6ba7b810-9dad-11d1-80b4-00c04fd430c8`,
	Run:  programHardwareIdentification,
	Args: cobra.ExactArgs(3),
}

func programHardwareIdentification(cmd *cobra.Command, args []string) {
	name := args[0]
	major, err := strconv.Atoi(args[1])
	e.ErrChk(err)
	serial := args[2]

	c, err := client.NewCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	err = c.ProgramHardwareIdentification(name, uint32(major), serial, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)
	fmt.Println("Success. Read back programmed ID")
	identifyHardwareFromClient(c)
}

func init() {
	rootCmd.AddCommand(programHardwareIdentificationCmd)
}
