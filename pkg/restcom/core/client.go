package core

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	userName  = "io4edge"
	urlPrefix = "/api/v1"
)

// Client represents a client for the io4edge core function via REST API
type Client struct {
	password   string
	httpClient *http.Client
	baseURL    *url.URL
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// New creates a new client for the io4edge core functions via REST API
func New(address string, password string) (*Client, error) {
	// ignore certificate errors, TODO: allow to configure
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{Transport: tr}
	urlStr := fmt.Sprintf("https://%s%s", address, urlPrefix)
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		password:   password,
		httpClient: client,
		baseURL:    baseURL,
	}
	return c, nil
}

// request sends a request to the device and returns the response.
// url is appended to the base URL and must start with a slash.
// body is the request body or nil.
// params are URL parameters or nil.
// caller must close resp.Body
func (c *Client) request(ctx context.Context, relPath string, verb string, body io.Reader, params map[string]string) (*http.Response, error) {
	full := *c.baseURL
	full.Path += relPath

	if params != nil {
		v := full.Query()
		for key, value := range params {
			v.Add(key, value)
		}
		full.RawQuery = v.Encode()
	}
	c.httpClient.Timeout = 0
	fmt.Printf("Requesting %s %s\n", verb, full.String())
	req, err := http.NewRequestWithContext(ctx, verb, full.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(userName, c.password)
	return c.httpClient.Do(req)
}

// requestMustBeOk sends a request to the device and returns the response body as
// bytes if the status code is 200
// see request for parameter description
// caller must close resp.Body if no error is returned
func (c *Client) requestMustBeOk(url string, verb string, body io.Reader, params map[string]string, timeout time.Duration) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.request(ctx, url, verb, body, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, c.decodeErrorResponse(resp)
	}
	return resp, nil
}

func (c *Client) decodeErrorResponse(resp *http.Response) error {
	detail := ""
	// check if response body is in json format
	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/json" {
		var errDetail errorResponse
		err := json.NewDecoder(resp.Body).Decode(&errDetail)
		if err == nil {
			detail = fmt.Sprintf(" (%s:%s)", errDetail.Code, errDetail.Message)
		}
	}

	return fmt.Errorf("unexpected status code %d%s", resp.StatusCode, detail)
}

// ReadPartition reads a partition from the device
func (c *Client) ReadPartition(timeout time.Duration, partitionName string, offset uint32, w *bufio.Writer, prog func(bytes uint, msg string)) (err error) {
	return fmt.Errorf("not implemented")
}

// ResetReason gets the reset reason
func (c *Client) ResetReason(timeout time.Duration) (reason string, err error) {
	return "", fmt.Errorf("not implemented")
}

// StreamLogs streams the logs from the device
func (c *Client) StreamLogs(streamTimeout time.Duration, infoCb func(msg string)) (io.ReadCloser, error) {
	return nil, fmt.Errorf("not implemented")
}

// Close closes the client. Here it does nothing
func (c *Client) Close() {
}
