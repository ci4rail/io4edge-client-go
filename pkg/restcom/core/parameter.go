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
// name is the name of the parameter, it can be in the form "namespace.parameter" or just "parameter"
func (c *Client) SetPersistentParameter(name string, value string, timeout time.Duration) (bool, error) {
	// encode value
	body, err := json.Marshal(parameterValue{Value: value})
	if err != nil {
		return false, fmt.Errorf("failed to encode value: %w", err)
	}
	resp, err := c.requestMustBeOk(parameterUrlFromName(name), http.MethodPut, bytes.NewReader(body), nil)
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
// name is the name of the parameter, it can be in the form "namespace.parameter" or just "parameter"
func (c *Client) GetPersistentParameter(name string, timeout time.Duration) (value string, err error) {
	resp, err := c.request(parameterUrlFromName(name), "GET", nil, nil)
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

// splitNameSpaceAndParameter splits a name in form "namespace.param" into namespace and parameter
// If no namespace is given, the namespace is empty
func splitNameSpaceAndParameter(name string) (string, string) {
	s := name
	if i := bytes.IndexByte([]byte(name), '.'); i >= 0 {
		return s[:i], s[i+1:]
	}
	return "", s
}

// parameterUrlFromName returns the URL for a parameter (with optional namespace)
func parameterUrlFromName(name string) string {
	namespace, parameter := splitNameSpaceAndParameter(name)
	if namespace == "" {
		return "/parameter/" + parameter
	}
	return "/" + namespace + "/parameter/" + parameter
}