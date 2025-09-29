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
	"os"
	"time"

	e "github.com/ci4rail/io4edge-client-go/v2/internal/errors"
	"github.com/ci4rail/io4edge-client-go/v2/pkg/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	paramFile string
)

var setPersistentParameterCmd = &cobra.Command{
	Use:     "set-parameter [NAME VALUE] [-f FILE]",
	Aliases: []string{"set-para", "set-persist"},
	Short:   "Set a persistent parameter.",
	Long: `Program a parameter into the non volatile storage (nvs) of the device.
It is only possible to set parameters for which the device already provides a place in the nvs.
Passing an empty value deletes the parameter.
It is possible to set multiple parameters at once by providing a YAML file with the parameters via the -f option.
Examples:
io4edge-cli -s S101-IOU04-USB-EXT-1 set-parameter wifi-ssid Ci4Rail-Guest
io4edge-cli -s S101-IOU04-USB-EXT-1 set-parameter -f wifi.yaml`,
	Run: setPersistentParameter,
	//Args: cobra.ExactArgs(2),
}

func setPersistentParameter(cmd *cobra.Command, args []string) {
	haveErrors := false
	rebootRequired := false

	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	if paramFile == "" {
		if len(args) != 2 {
			fmt.Printf("Error: set-parameter requires exactly two arguments: NAME VALUE\n")
			os.Exit(1)
		}
		name := args[0]
		value := args[1]

		rebootRequired, err = setAndVerifyParameter(c, name, value)
		if err != nil {
			fmt.Printf("Error setting parameter %s: %v\n", name, err)
			haveErrors = true
		}
	} else {
		data, err := os.ReadFile(paramFile)
		e.ErrChk(err)
		var params map[string]string
		err = yaml.Unmarshal(data, &params)
		e.ErrChk(err)

		fmt.Printf("Setting parameters...\n")
		for name, value := range params {
			reboot, err := setAndVerifyParameter(c, name, value)
			if err != nil {
				fmt.Printf("Error setting parameter %s: %v\n", name, err)
				haveErrors = true
			}
			if reboot {
				rebootRequired = true
			}
		}
	}
	if rebootRequired {
		fmt.Printf("Reboot is required to apply parameters\n")
	}
	if haveErrors {
		os.Exit(1)
	}
}

// setAndVerifyParameter sets a parameter and verifies it was set correctly
// returns rebootRequired and error
func setAndVerifyParameter(c core.If, name string, value string) (bool, error) {
	rebootRequired, err := c.SetPersistentParameter(name, value, time.Duration(timeoutSecs)*time.Second)
	if err != nil {
		return rebootRequired, err
	}
	if value != "" {
		value, err = c.GetPersistentParameter(name, time.Duration(timeoutSecs)*time.Second)
		if err != nil {
			if _, ok := err.(*core.ParameterIsReadProtectedError); ok {
				fmt.Printf("Couldn't read back parameter %s. May be it's read-only?\n", name)
			} else {
				return rebootRequired, fmt.Errorf("error reading back parameter %s: %v", name, err)
			}
		} else {
			fmt.Printf("%s = %s\n", name, value)
		}
	} else {
		fmt.Printf("Parameter %s deleted\n", name)
	}
	return rebootRequired, nil
}

func init() {
	rootCmd.AddCommand(setPersistentParameterCmd)
	setPersistentParameterCmd.Flags().StringVarP(&paramFile, "file", "f", "", "YAML file containing parameters to set")
}
