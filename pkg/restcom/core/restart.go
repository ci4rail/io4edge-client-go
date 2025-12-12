package core

import (
	"context"
	"net/http"
	"time"
)

// Restart restarts the device
func (c *Client) Restart(timeout time.Duration) (restartingNow bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.requestMustBeOk(ctx, "/restart", http.MethodPost, nil, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}
