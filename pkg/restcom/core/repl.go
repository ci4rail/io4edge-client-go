package core

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

// ReplCommand sends a REPL command to the device and returns the response
func (c *Client) ReplCommand(cmd string, timeout time.Duration) (response string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.requestMustBeOk(ctx, "/repl", http.MethodPost, bytes.NewReader([]byte(cmd)), nil)
	if err != nil {
		return "", fmt.Errorf("failed to execute repl command: %w", err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read repl response: %w", err)
	}
	return buf.String(), nil
}
