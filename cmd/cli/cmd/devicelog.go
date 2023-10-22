/*
Copyright © 2023 Ci4Rail GmbH <engineering@ci4rail.com>

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

var (
	logStreamTimeout uint
)

var deviceLogCmd = &cobra.Command{
	Use:   "log",
	Short: "stream device log",
	Run:   streamLogs,
}

func streamLogs(cmd *cobra.Command, args []string) {
	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	r, err := c.StreamLogs(time.Duration(logStreamTimeout)*time.Second, streamLogInfoCb)
	e.ErrChk(err)

	defer r.Close()

	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err != nil {
			e.ErrChk(err)
			break
		}
		fmt.Print(string(buf[:n]))
		//fmt.Printf("Logs: %d bytes\n", n)
	}
}

func streamLogInfoCb(info string) {
	fmt.Printf("%s", info)
}

func init() {
	rootCmd.AddCommand(deviceLogCmd)
	deviceLogCmd.PersistentFlags().UintVarP(&logStreamTimeout, "stream-timeout", "", 3, "Timeout in seconds for stream, reestablishes connection if no data is received for this time. 0 means no timeout.")
}
