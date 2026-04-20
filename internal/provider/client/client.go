// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

// DoRequestWithLicenseKeyHeader performs an HTTP request with X-License-Key header
// Some API endpoints (like holiday API) use X-License-Key instead of X-API-Key
func (c *Client) DoRequestWithLicenseKeyHeader(method, path string, body interface{}) ([]byte, int, error) {
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

	// Add authentication headers with X-License-Key
	req.Header.Set("X-User-Name", c.Username)
	req.Header.Set("X-License-Key", c.LicenseKey)
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

// DoRequestWithCookieAuth performs an HTTP request using Cookie-based auth.
// Some older /api/v1/ endpoints (e.g. logdedicatedmode) only accept Cookie: userName=...
// for authentication and ignore X-User-Name/X-API-Key headers.
// Credentials (userName, licenseKey) must be included in the path as query params by the caller.
func (c *Client) DoRequestWithCookieAuth(method, path string) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("userName=%s;", c.Username))

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

// DoMultipartFormRequest performs an HTTP request with multipart/form-data.
// fields contains plain-text form fields; fileParts contains file-like parts sent as file uploads.
func (c *Client) DoMultipartFormRequest(method, path string, fields map[string]string, fileParts map[string][]byte) ([]byte, int, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Add authentication fields
	if err := w.WriteField("userName", c.Username); err != nil {
		return nil, 0, fmt.Errorf("failed to write userName field: %w", err)
	}
	if err := w.WriteField("licenseKey", c.LicenseKey); err != nil {
		return nil, 0, fmt.Errorf("failed to write licenseKey field: %w", err)
	}

	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, 0, fmt.Errorf("failed to write field %s: %w", k, err)
		}
	}
	for fieldName, data := range fileParts {
		part, err := w.CreateFormFile(fieldName, fieldName+".json")
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create form file part: %w", err)
		}
		if _, err := part.Write(data); err != nil {
			return nil, 0, fmt.Errorf("failed to write file data for part %s: %w", fieldName, err)
		}
	}
	w.Close()

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, fullURL, &buf)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-User-Name", c.Username)
	req.Header.Set("X-API-Key", c.LicenseKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

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
