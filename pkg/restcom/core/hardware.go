package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ci4rail/io4edge-client-go/v2/pkg/core"
)

// IdentifyHardware gets the firmware name and version from the device
// TODO: support customer part number
func (c *Client) IdentifyHardware(timeout time.Duration) (*core.HardwareInventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.requestMustBeOk(ctx, "/hardware", http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var id map[string]any
	err = json.NewDecoder(resp.Body).Decode(&id)

	if err != nil {
		return nil, err
	}
	partNumber, ok := id["part_number"].(string)
	if !ok {
		return nil, fmt.Errorf("part_number missing or not a string")
	}
	serialNumber, ok := id["serial_number"].(string)
	if !ok {
		return nil, fmt.Errorf("serial_number missing or not a string")
	}
	majorVersion, ok := id["major_version"].(float64)
	if !ok {
		return nil, fmt.Errorf("major_version missing or not a number")
	}

	customExtension := make(map[string]string)
	// look for custom extension
	for k, v := range id {
		if k == "part_number" || k == "serial_number" || k == "major_version" {
			continue
		}
		if _, ok := v.(string); ok {
			customExtension[k] = v.(string)
		}
	}

	return &core.HardwareInventory{
		PartNumber:      partNumber,
		SerialNumber:    serialNumber,
		MajorVersion:    uint32(majorVersion),
		CustomExtension: customExtension,
	}, nil
}

// ProgramHardwareIdentification programs the hardware identification
func (c *Client) ProgramHardwareIdentification(i *core.HardwareInventory, timeout time.Duration) error {
	id := map[string]interface{}{
		"part_number":   i.PartNumber,
		"serial_number": i.SerialNumber,
		"major_version": int(i.MajorVersion),
	}
	for k, v := range i.CustomExtension {
		id[k] = v
	}

	body, err := json.Marshal(id)
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.requestMustBeOk(ctx, "/hardware", http.MethodPut, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to program hardware info: %w", err)
	}
	defer resp.Body.Close()
	return nil
}
