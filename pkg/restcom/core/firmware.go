package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/core"
)

type getFirmwareResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// IdentifyFirmware gets the firmware name and version from the device
func (c *Client) IdentifyFirmware(timeout time.Duration) (name string, version string, err error) {
	resp, err := c.requestMustBeOk("/firmware", http.MethodGet, nil, nil, timeout)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var id getFirmwareResponse
	err = json.NewDecoder(resp.Body).Decode(&id)

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
	return core.LoadFirmware(c, file, chunkSize, timeout, prog)
}

// LoadFirmwareFromBuffer loads a binary from a firmware package in memory to the device.
// Checks first if the firmware is compatible with the device.
// Checks then if the device's firmware version is the same
// timeout is for each chunk
func (c *Client) LoadFirmwareFromBuffer(buffer []byte, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	return core.LoadFirmwareFromBuffer(c, buffer, chunkSize, timeout, prog)
}

// LoadFirmwareBinaryFromFile loads new firmware from file into the device device
// timeout is for each chunk
func (c *Client) LoadFirmwareBinaryFromFile(file string, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	return core.LoadFirmwareBinaryFromFile(c, file, chunkSize, timeout, prog)
}

// LoadFirmwareBinary loads new firmware from r into the device device
func (c *Client) LoadFirmwareBinary(r *bufio.Reader, chunkSize uint, timeout time.Duration, prog func(bytes uint, msg string)) (restartingNow bool, err error) {
	totalBytes := uint(0)
	restartingNow = false

	data := make([]byte, chunkSize)

	for {
		var err error
		atEOF := false

		_, err = io.ReadFull(r, data)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// end of file reached
			atEOF = true
		} else if err != nil {
			return restartingNow, errors.New("read firmware failed: " + err.Error())
		}

		urlParams := map[string]string{
			"offset": fmt.Sprintf("%d", totalBytes),
			"last":   fmt.Sprintf("%t", atEOF),
		}

		try := 3

		for try = 3; try >= 0; try-- {
			// create io.reader from data
			body := bytes.NewReader(data)

			_, err = c.requestMustBeOk("/firmware", http.MethodPut, body, urlParams, timeout)
			if err == nil {
				break
			}
			prog(totalBytes, fmt.Sprintf("Error %s Retry...", err))
		}
		if try < 0 || err != nil {
			return restartingNow, errors.New("load firmware chunk command failed: " + err.Error())
		}

		totalBytes += uint(len(data))
		prog(totalBytes, "")

		if atEOF {
			restartingNow = true // TODO: response has no info if restart is needed
			break
		}
	}

	return restartingNow, nil
}
