package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	fwpkg "github.com/ci4rail/firmware-packaging-go"
)

// contains functions for firmware handling shared between REST and proto clients

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

// LoadFirmware loads a binary from a firmware package to the device.
// Checks first if the firmware is compatible with the device.
// Checks then if the device's firmware version is the same
// timeout is for each chunk
func LoadFirmware(c If, file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
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
	i, err := c.IdentifyHardware(timeout)
	if err != nil {
		return restartingNow, err
	}

	// check compatibility
	err = AssertFirmwareIsCompatibleWithHardware(
		manifest.Compatibility.HW,
		manifest.Compatibility.MajorRevs,
		i.PartNumber,
		int(i.MajorVersion),
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
func LoadFirmwareBinaryFromFile(c If, file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	f, err := os.Open(file)
	if err != nil {
		return false, errors.New("cannot open file " + file + " : " + err.Error())
	}

	defer f.Close()

	r := bufio.NewReader(f)
	return c.LoadFirmwareBinary(r, chunkSize, timeout, prog)
}
