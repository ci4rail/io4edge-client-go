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

	e "github.com/ci4rail/io4edge-client-go/v2/internal/errors"
	"github.com/spf13/cobra"
)

var getPersistentParameterCmd = &cobra.Command{
	Use:     "get-parameter NAME",
	Aliases: []string{"get-para", "get-persist"},
	Short:   "Get a persistent parameter.",
	Long: `Read a parameter from the non volatile storage (nvs) of the device.
Example:
io4edge-cli -s S101-IOU04-USB-EXT-1 get-parameter wifi-ssid`,
	Run:  getPersistentParameter,
	Args: cobra.ExactArgs(1),
}

func getPersistentParameter(cmd *cobra.Command, args []string) {
	name := args[0]

	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	value, err := c.GetPersistentParameter(name, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	fmt.Printf("Read parameter name: %s, value: %s\n", name, value)
}

func init() {
	rootCmd.AddCommand(getPersistentParameterCmd)
}
