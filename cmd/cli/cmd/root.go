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
	e "github.com/ci4rail/io4edge-client-go/v2/internal/errors"

	"github.com/spf13/cobra"
)

var (
	deviceID    string
	ipAddrPort  string
	timeoutSecs int
	password    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "io4edge-cli",
	Short: "io4edge cli",
	Long: `Command line interface to communicate with Ci4Rail io4edge devices
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e.Er(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&deviceID, "device id", "d", "", "Distinct designation of the device (mdns instance name of the the device)")
	rootCmd.PersistentFlags().StringVarP(&ipAddrPort, "ip address", "i", "", "IP address of io4edge devices with port e.g. 192.168.200.1:9999")
	rootCmd.PersistentFlags().IntVarP(&timeoutSecs, "timeout", "t", 30, "Timeout in seconds to wait for device responses")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "core_io4edge", "Password for the REST API")
}
