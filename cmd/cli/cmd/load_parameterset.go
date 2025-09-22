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

	e "github.com/ci4rail/io4edge-client-go/internal/errors"
	"github.com/spf13/cobra"
)

var loadParameterSetCmd = &cobra.Command{
	Use:     "load-parameterset NAMESPACE FILE",
	Aliases: []string{"load-ps"},
	Short:   "Load a parameter set.",
	Long:    `Load parameter set of namespace from file.`,
	Run:     loadParameterSet,
	Args:    cobra.ExactArgs(2),
}

func loadParameterSet(cmd *cobra.Command, args []string) {
	namespace := args[0]
	filename := args[1]

	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	// read all bytes from file
	bytes, err := os.ReadFile(filename)
	e.ErrChk(err)

	bytes, err = c.LoadParameterSet(time.Duration(timeoutSecs)*time.Second, namespace, bytes)
	e.ErrChk(err)

	fmt.Println("Parameter set loaded. Response:")
	fmt.Println(string(bytes))
}

func init() {
	rootCmd.AddCommand(loadParameterSetCmd)
}
