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
	"bufio"
	"os"
	"time"

	e "github.com/ci4rail/io4edge-client-go/internal/errors"
	"github.com/spf13/cobra"
)

var readPartitionCmd = &cobra.Command{
	Use:     "read-partition PARTITION_NAME OUTPUT_FILE",
	Aliases: []string{"part", "rp"},
	Short:   "Read partition from device",
	Run:     readPartition,
	Args:    cobra.ExactArgs(2),
}

func readPartition(cmd *cobra.Command, args []string) {
	c, err := newCliClient(deviceID, ipAddrPort)
	e.ErrChk(err)

	partitionName := args[0]
	outputFile := args[1]

	f, err := os.Create(outputFile)
	e.ErrChk(err)
	defer f.Close()

	// open output file as bufio.Writer
	w := bufio.NewWriter(f)

	err = c.ReadPartition(time.Duration(timeoutSecs)*time.Second, partitionName, w)
	e.ErrChk(err)

}

func init() {
	rootCmd.AddCommand(readPartitionCmd)
}
