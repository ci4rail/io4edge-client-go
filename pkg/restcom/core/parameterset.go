package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetParameterSet gets the parameter set from the device
func (c *Client) GetParameterSet(timeout time.Duration, namespace string) ([]byte, error) {
	resp, err := c.requestMustBeOk(parameterSetURL(namespace), http.MethodGet, nil, nil, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameter set: %w", err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	return bytes, nil
}

// LoadParameterSet loads the parameter set to the device
func (c *Client) LoadParameterSet(timeout time.Duration, namespace string, data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)

	resp, err := c.requestMustBeOk(parameterSetURL(namespace), http.MethodPut, reader, nil, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to load parameter set: %w", err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	return bytes, nil
}

func parameterSetURL(namespace string) string {
	return fmt.Sprintf("/%s/parameterset", namespace)
}
