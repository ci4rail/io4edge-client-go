package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ChangeAPIPassword changes the API password of the device
func (c *Client) ChangeAPIPassword(newPassword string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	passwdChange := map[string]string{
		"password": newPassword,
	}
	body, err := json.Marshal(passwdChange)
	if err != nil {
		return fmt.Errorf("failed to marshal password change request: %w", err)
	}
	resp, err := c.requestMustBeOk(ctx, "/users/io4edge/basic_auth", http.MethodPut, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to change API password: %w", err)
	}
	defer resp.Body.Close()
	return nil
}
