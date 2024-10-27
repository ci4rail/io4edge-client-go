package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ci4rail/io4edge-client-go/pkg/core"
)

type parameterValue struct {
	Value string `json:"value"`
}

type setParameterResponse struct {
	RebootRequired bool `json:"reboot_required"`
}

// SetPersistentParameter sets a persistent parameter
func (c *Client) SetPersistentParameter(name string, value string, timeout time.Duration) (bool, error) {
	// encode value
	body, err := json.Marshal(parameterValue{Value: value})
	if err != nil {
		return false, fmt.Errorf("failed to encode value: %w", err)
	}
	resp, err := c.requestMustBeOk("/parameter/"+name, "PUT", bytes.NewReader(body), nil)
	if err != nil {
		return false, fmt.Errorf("failed to set parameter: %w", err)
	}
	var r setParameterResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}
	return r.RebootRequired, nil
}

// GetPersistentParameter gets a persistent parameter
func (c *Client) GetPersistentParameter(name string, timeout time.Duration) (value string, err error) {
	resp, err := c.request("/parameter/"+name, "GET", nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to exec parameter get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			return "", &core.ParameterIsReadProtectedError{}
		}
		return "", fmt.Errorf("failed to get parameter %w", c.decodeErrorResponse(resp))
	}
	var p parameterValue
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return "", fmt.Errorf("failed to decode parameter: %w", err)
	}
	return p.Value, nil
}
