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
	"time"

	e "github.com/ci4rail/io4edge-client-go/internal/errors"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart device",
	Run:   restart,
}

func restart(cmd *cobra.Command, args []string) {
	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	_, err = c.Restart(time.Duration(timeoutSecs) * time.Second)
	e.ErrChk(err)
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
