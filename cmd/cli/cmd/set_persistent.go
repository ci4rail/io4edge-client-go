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

	e "github.com/ci4rail/io4edge-client-go/internal/errors"
	"github.com/spf13/cobra"
)

var setPersistentParameterCmd = &cobra.Command{
	Use:     "set-parameter NAME VALUE",
	Aliases: []string{"set-para", "set-persist"},
	Short:   "Set a persistent parameter.",
	Long: `Program a parameter into the non volatile storage (nvs) of the device.
While the name is the key to the value. It is only possible to set parameters for which the device already provides a place in the nvs.
Example:
io4edge-cli -s S101-IOU04-USB-EXT-1 set-parameter wifi-ssid Ci4Rail-Guest`,
	Run:  setPersistentParameter,
	Args: cobra.ExactArgs(2),
}

func setPersistentParameter(cmd *cobra.Command, args []string) {
	name := args[0]
	value := args[1]

	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	err = c.SetPersistentParameter(name, value, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	value, err = c.GetPersistentParameter(name, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	fmt.Printf("Parameter %s was set to %s\n", name, value)
}

func init() {
	rootCmd.AddCommand(setPersistentParameterCmd)
}
