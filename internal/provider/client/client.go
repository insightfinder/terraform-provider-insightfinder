// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the InsightFinder API client
type Client struct {
	BaseURL    string
	Username   string
	LicenseKey string
	HTTPClient *http.Client
}

// NewClient creates a new InsightFinder API client
func NewClient(baseURL, username, licenseKey string) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if licenseKey == "" {
		return nil, fmt.Errorf("license key cannot be empty")
	}

	return &Client{
		BaseURL:    baseURL,
		Username:   username,
		LicenseKey: licenseKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

// DoRequest performs an HTTP request with authentication headers
func (c *Client) DoRequest(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	req.Header.Set("X-User-Name", c.Username)
	req.Header.Set("X-API-Key", c.LicenseKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// DoFormRequest performs an HTTP request with form data
func (c *Client) DoFormRequest(method, path string, formData url.Values) ([]byte, int, error) {
	// Add authentication to form data
	formData.Set("userName", c.Username)
	formData.Set("licenseKey", c.LicenseKey)

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers (some endpoints might use these)
	req.Header.Set("X-User-Name", c.Username)
	req.Header.Set("X-API-Key", c.LicenseKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// APIError represents an error response from the API
type APIError struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode,omitempty"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error (code %d): %s", e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("API error (code %d)", e.ErrorCode)
}
