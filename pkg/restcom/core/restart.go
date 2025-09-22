package core

import (
	"net/http"
	"time"
)

// Restart restarts the device
func (c *Client) Restart(timeout time.Duration) (restartingNow bool, err error) {
	resp, err := c.requestMustBeOk("/restart", http.MethodPost, nil, nil, timeout)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}
