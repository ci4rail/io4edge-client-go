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

package basefunc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	fwpkg "github.com/ci4rail/firmware-packaging-go"
)

// FirmwareAlreadyPresentError is returned by LoadFirmware as a dummy error
type FirmwareAlreadyPresentError struct {
}

// Error returns the error string for FirmwareAlreadyPresentError
func (e *FirmwareAlreadyPresentError) Error() string {
	return "Requested Firmware already present"
}

// IdentifyFirmware gets the firmware name and version from the device
func (c *Client) IdentifyFirmware(timeout time.Duration) (*ResIdentifyFirmware, error) {
	cmd := &BaseFuncCommand{
		Id: BaseFuncCommandId_IDENTIFY_FIRMWARE,
	}
	res := &BaseFuncResponse{}
	if err := c.Command(cmd, res, timeout); err != nil {
		return nil, err
	}
	return res.GetIdentifyFirmware(), nil
}

// LoadFirmware loads a binary from a firmware package to the device
//
// timeout is for each chunk
func (c *Client) LoadFirmware(file string, chunkSize uint, timeout time.Duration) error {

	pkg, err := fwpkg.NewFirmwarePackageConsumerFromFile(file)
	if err != nil {
		return err
	}
	manifest := pkg.Manifest()

	// get currently running firmware
	currentFWID, err := c.IdentifyFirmware(timeout)
	if err != nil {
		return err
	}

	// get devices hardware id
	hwID, err := c.IdentifyHardware(timeout)
	if err != nil {
		return err
	}

	// check compatibility
	err = AssertFirmwareIsCompatibleWithHardware(
		manifest.Compatibility.HW,
		manifest.Compatibility.MajorRevs,
		hwID.RootArticle,
		int(hwID.MajorVersion),
	)
	if err != nil {
		return err
	}

	// check if fw already running
	curVer := fmt.Sprintf("%d.%d.%d", currentFWID.MajorVersion, currentFWID.MinorVersion, currentFWID.PatchVersion)

	if strings.EqualFold(currentFWID.Name, manifest.Name) && curVer == manifest.Version {
		return &FirmwareAlreadyPresentError{}
	}

	fwFile := new(bytes.Buffer)
	err = pkg.File(fwFile)
	if err != nil {
		return err
	}

	err = c.LoadFirmwareBinary(bufio.NewReader(fwFile), chunkSize, timeout)
	if err != nil {
		return err
	}
	return err
}

// LoadFirmwareBinaryFromFile loads new firmware from file into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.New("cannot open file " + file + " : " + err.Error())
	}

	defer f.Close()

	r := bufio.NewReader(f)
	return c.LoadFirmwareBinary(r, chunkSize, timeout)
}

// LoadFirmwareBinary loads new firmware via r into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration) error {

	cmd := &BaseFuncCommand{
		Id: BaseFuncCommandId_LOAD_FIRMWARE_CHUNK,
		Data: &BaseFuncCommand_LoadFirmwareChunk{
			LoadFirmwareChunk: &CmdLoadFirmwareChunk{
				Data: make([]byte, chunkSize),
			},
		},
	}
	data := cmd.GetLoadFirmwareChunk().Data

	chunkNumber := uint32(0)

	for {
		atEOF := false

		_, err := r.Read(data)
		if err != nil {
			return errors.New("read firmware failed: " + err.Error())
		}

		// check if we are at EOF
		_, err = r.Peek(1)
		if err == io.EOF {
			atEOF = true
		}

		cmd.GetLoadFirmwareChunk().IsLastChunk = atEOF
		cmd.GetLoadFirmwareChunk().ChunkNumber = chunkNumber

		res := &BaseFuncResponse{}
		err = c.Command(cmd, res, timeout)
		if err != nil {
			return errors.New("load firmware chunk command failed: " + err.Error())
		}

		if atEOF {
			break
		}
		chunkNumber++
	}

	return nil
}

// AssertFirmwareIsCompatibleWithHardware checks if the firmware specified by fwHw and fwMajorRevs is compatible
// with hardware hwName, hwMajor
func AssertFirmwareIsCompatibleWithHardware(fwHw string, fwMajorRevs []int, hwName string, hwMajor int) error {
	if !strings.EqualFold(fwHw, hwName) {
		return errors.New("firmware " + fwHw + " not suitable for hardware " + hwName)
	}
	var ok = false
	for _, b := range fwMajorRevs {
		if hwMajor == b {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("firmware doesn't support hardware version %d", hwMajor)
	}
	return nil
}
