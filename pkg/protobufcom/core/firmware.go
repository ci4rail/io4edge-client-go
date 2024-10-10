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

package pbcore

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	fwpkg "github.com/ci4rail/firmware-packaging-go"
	"github.com/ci4rail/io4edge-client-go/pkg/core"
	api "github.com/ci4rail/io4edge_api/io4edge/go/core_api/v1alpha2"
)

// FirmwareAlreadyPresentError is returned by LoadFirmware as a dummy error
type FirmwareAlreadyPresentError struct {
}

// Error returns the error string for FirmwareAlreadyPresentError
func (e *FirmwareAlreadyPresentError) Error() string {
	return "Requested Firmware already present"
}

// IdentifyFirmware gets the firmware name and version from the device
func (c *Client) IdentifyFirmware(timeout time.Duration) (name string, version string, err error) {
	cmd := &api.CoreCommand{
		Id: api.CommandId_IDENTIFY_FIRMWARE,
	}
	res := &api.CoreResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return "", "", err
	}
	return res.GetIdentifyFirmware().Name, res.GetIdentifyFirmware().Version, nil
}

// LoadFirmware loads a binary from a firmware package to the device.
// Checks first if the firmware is compatible with the device.
// Checks then if the device's firmware version is the same
// timeout is for each chunk
func (c *Client) LoadFirmware(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	restartingNow = false

	pkg, err := fwpkg.NewFirmwarePackageConsumerFromFile(file)
	if err != nil {
		return restartingNow, err
	}
	manifest := pkg.Manifest()

	// get currently running firmware
	fwName, fwVersion, err := c.IdentifyFirmware(timeout)
	if err != nil {
		return restartingNow, err
	}

	// get devices hardware id
	rootArticle, majorVersion, _, err := c.IdentifyHardware(timeout)
	if err != nil {
		return restartingNow, err
	}

	// check compatibility
	err = core.AssertFirmwareIsCompatibleWithHardware(
		manifest.Compatibility.HW,
		manifest.Compatibility.MajorRevs,
		rootArticle,
		int(majorVersion),
	)
	if err != nil {
		return restartingNow, err
	}

	// check if fw already running
	if strings.EqualFold(fwName, manifest.Name) && fwVersion == manifest.Version {
		return restartingNow, &FirmwareAlreadyPresentError{}
	}

	fwFile := new(bytes.Buffer)
	err = pkg.File(fwFile)
	if err != nil {
		return restartingNow, err
	}

	restartingNow, err = c.LoadFirmwareBinary(bufio.NewReader(fwFile), chunkSize, timeout, prog)
	return restartingNow, err
}

// LoadFirmwareBinaryFromFile loads new firmware from file into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	f, err := os.Open(file)
	if err != nil {
		return false, errors.New("cannot open file " + file + " : " + err.Error())
	}

	defer f.Close()

	r := bufio.NewReader(f)
	return c.LoadFirmwareBinary(r, chunkSize, timeout, prog)
}

// LoadFirmwareBinary loads new firmware via r into the device device
// prog is a callback function that gets called after loading a chunk. The callback gets passed the number of bytes transferred yet
// timeout is for each chunk
func (c *Client) LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {

	cmd := &api.CoreCommand{
		Id: api.CommandId_LOAD_FIRMWARE_CHUNK,
		Data: &api.CoreCommand_LoadFirmwareChunk{
			LoadFirmwareChunk: &api.LoadFirmwareChunkCommand{
				Data: make([]byte, chunkSize),
			},
		},
	}
	data := cmd.GetLoadFirmwareChunk().Data

	chunkNumber := uint32(0)
	totalBytes := uint(0)
	restartingNow = false

	for {
		atEOF := false

		_, err := r.Read(data)
		if err != nil {
			return restartingNow, errors.New("read firmware failed: " + err.Error())
		}

		// check if we are at EOF
		_, err = r.Peek(1)
		if err == io.EOF {
			atEOF = true
		}

		cmd.GetLoadFirmwareChunk().IsLastChunk = atEOF
		cmd.GetLoadFirmwareChunk().ChunkNumber = chunkNumber

		try := 3
		var res *api.CoreResponse

		for try = 3; try >= 0; try-- {
			res = &api.CoreResponse{}
			err = c.Command(cmd, res, timeout)
			if err == nil {
				break
			}
			if res.Status != api.Status_OK {
				// don't retry on errors returned by the device
				break
			}
			prog(totalBytes, "Reestablish connection to device...")
			err = c.funcClient.RestartChannel()
			if err != nil {
				return restartingNow, errors.New("load firmware chunk command failed (can't restart connection): " + err.Error())
			}
		}
		if try < 0 || err != nil {
			return restartingNow, errors.New("load firmware chunk command failed: " + err.Error())
		}

		totalBytes += uint(len(data))
		prog(totalBytes, "")

		restartingNow = res.RestartingNow
		if atEOF {
			break
		}
		chunkNumber++
	}

	return restartingNow, nil
}
