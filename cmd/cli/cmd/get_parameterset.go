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

var getParameterSetCmd = &cobra.Command{
	Use:     "get-parameterset NAMESPACE",
	Aliases: []string{"get-ps"},
	Short:   "Get a parameter set.",
	Long:    `Read and dump parameter set from namespace.`,
	Run:     getParameterSet,
	Args:    cobra.ExactArgs(1),
}

func getParameterSet(cmd *cobra.Command, args []string) {
	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	bytes, err := c.GetParameterSet(time.Duration(timeoutSecs)*time.Second, args[0])
	e.ErrChk(err)

	fmt.Println(string(bytes))
}

func init() {
	rootCmd.AddCommand(getParameterSetCmd)
}
