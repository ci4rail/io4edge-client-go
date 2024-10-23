package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	fwpkg "github.com/ci4rail/firmware-packaging-go"
	"github.com/ci4rail/io4edge-client-go/pkg/core"
)

type getFirmwareResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// IdentifyFirmware gets the firmware name and version from the device
func (c *Client) IdentifyFirmware(timeout time.Duration) (name string, version string, err error) {
	resp, err := c.requestMustBeOk("/firmware", http.MethodGet, nil)
	if err != nil {
		return "", "", err
	}
	var id getFirmwareResponse

	err = json.Unmarshal(resp, &id)
	if err != nil {
		return "", "", err
	}
	return id.Name, id.Version, nil
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
		return restartingNow, &core.FirmwareAlreadyPresentError{}
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
func (c *Client) LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	return false, fmt.Errorf("not implemented")
}

// LoadFirmwareBinary loads new firmware from r into the device device
func (c *Client) LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	return false, fmt.Errorf("not implemented")
}
