package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type hardwareValue struct {
	PartNumber   string `json:"part_number"`
	SerialNumber string `json:"serial_number"`
	MajorVersion int    `json:"major_version"`
}

// IdentifyHardware gets the firmware name and version from the device
// TODO: support customer part number
func (c *Client) IdentifyHardware(timeout time.Duration) (name string, major uint32, serial string, err error) {
	resp, err := c.requestMustBeOk("/hardware", http.MethodGet, nil, nil)
	if err != nil {
		return "", 0, "", err
	}
	var id hardwareValue
	err = json.NewDecoder(resp.Body).Decode(&id)

	if err != nil {
		return "", 0, "", err
	}
	return id.PartNumber, uint32(id.MajorVersion), id.SerialNumber, nil
}

// ProgramHardwareIdentification programs the hardware identification
func (c *Client) ProgramHardwareIdentification(name string, major uint32, serial string, timeout time.Duration) error {
	body, err := json.Marshal(hardwareValue{
		PartNumber:   name,
		SerialNumber: serial,
		MajorVersion: int(major),
	})
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}
	_, err = c.requestMustBeOk("/hardware", http.MethodPut, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to program hardware info: %w", err)
	}
	return nil
}
