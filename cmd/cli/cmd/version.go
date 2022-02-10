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

	"runtime/debug"

	"github.com/ci4rail/io4edge-client-go/internal/version"
	"github.com/spf13/cobra"
)

const (
	io4edge_api_package = "github.com/ci4rail/io4edge_api"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information and quit",
	Long: `Print version information and quit
This command displays version information for the io4edge-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			panic("Can't read BuildInfo")
		}
		fmt.Printf("io4edge-cli version: %s\n", version.Version)
		for _, dep := range buildInfo.Deps {
			if dep.Path == io4edge_api_package {
				fmt.Printf("io4edge api version: %s\n", dep.Version)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
