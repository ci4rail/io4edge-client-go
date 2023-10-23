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

package core

import (
	"bufio"
	"time"

	api "github.com/ci4rail/io4edge_api/io4edge/go/core_api/v1alpha2"
)

// ReadPartition reads a partition from the device and writes it to the given writer
func (c *Client) ReadPartition(timeout time.Duration, partitionName string, offset uint32, w *bufio.Writer, prog func(bytes uint, msg string)) (err error) {

	for {
		cmd := &api.CoreCommand{
			Id: api.CommandId_READ_PARTITION_CHUNK,
			Data: &api.CoreCommand_ReadPartitionChunk{
				ReadPartitionChunk: &api.ReadPartitionChunkCommand{
					PartName: partitionName,
					Offset:   offset,
				},
			},
		}
		res := &api.CoreResponse{}
		if err := c.Command(cmd, res, timeout); err != nil {
			return err
		}
		chunk := res.GetReadPartitionChunk()
		len := len(chunk.Data)
		if len == 0 {
			// no more data in partition
			break
		}
		_, err = w.Write(chunk.Data)
		if err != nil {
			return err
		}
		offset += uint32(len)
		prog(uint(offset), "")
	}
	return nil
}
