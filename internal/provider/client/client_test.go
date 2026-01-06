// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		username    string
		licenseKey  string
		expectError bool
	}{
		{
			name:        "valid client",
			baseURL:     "https://test.insightfinder.com",
			username:    "test_user",
			licenseKey:  "test_key",
			expectError: false,
		},
		{
			name:        "empty base URL",
			baseURL:     "",
			username:    "test_user",
			licenseKey:  "test_key",
			expectError: true,
		},
		{
			name:        "empty username",
			baseURL:     "https://test.insightfinder.com",
			username:    "",
			licenseKey:  "test_key",
			expectError: true,
		},
		{
			name:        "empty license key",
			baseURL:     "https://test.insightfinder.com",
			username:    "test_user",
			licenseKey:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.username, tt.licenseKey)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if client != nil {
					t.Error("Expected nil client on error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if client == nil {
					t.Error("Expected client, got nil")
				}
				if client != nil {
					if client.BaseURL != tt.baseURL {
						t.Errorf("Expected BaseURL '%s', got '%s'", tt.baseURL, client.BaseURL)
					}
					if client.Username != tt.username {
						t.Errorf("Expected Username '%s', got '%s'", tt.username, client.Username)
					}
					if client.LicenseKey != tt.licenseKey {
						t.Errorf("Expected LicenseKey '%s', got '%s'", tt.licenseKey, client.LicenseKey)
					}
					if client.HTTPClient == nil {
						t.Error("Expected HTTPClient to be initialized")
					}
				}
			}
		})
	}
}

func TestDoRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		serverResponse string
		serverStatus   int
		expectError    bool
	}{
		{
			name:           "successful GET request",
			method:         "GET",
			path:           "/api/test",
			body:           nil,
			serverResponse: `{"status":"success"}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful POST request",
			method:         "POST",
			path:           "/api/test",
			body:           map[string]string{"key": "value"},
			serverResponse: `{"status":"created"}`,
			serverStatus:   http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "server error",
			method:         "GET",
			path:           "/api/error",
			body:           nil,
			serverResponse: `{"error":"internal server error"}`,
			serverStatus:   http.StatusInternalServerError,
			expectError:    false, // DoRequest doesn't error on non-2xx status
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify headers
				if r.Header.Get("X-User-Name") == "" {
					t.Error("Expected X-User-Name header")
				}
				if r.Header.Get("X-API-Key") == "" {
					t.Error("Expected X-API-Key header")
				}

				// Verify method and path
				if r.Method != tt.method {
					t.Errorf("Expected method '%s', got '%s'", tt.method, r.Method)
				}
				if r.URL.Path != tt.path {
					t.Errorf("Expected path '%s', got '%s'", tt.path, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create client
			client, err := NewClient(server.URL, "test_user", "test_key")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Execute request
			respBody, statusCode, err := client.DoRequest(tt.method, tt.path, tt.body)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if statusCode != tt.serverStatus {
					t.Errorf("Expected status code %d, got %d", tt.serverStatus, statusCode)
				}
				if string(respBody) != tt.serverResponse {
					t.Errorf("Expected response body '%s', got '%s'", tt.serverResponse, string(respBody))
				}
			}
		})
	}
}

func TestDoFormRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		formData       url.Values
		serverResponse string
		serverStatus   int
		expectError    bool
	}{
		{
			name:   "successful form POST",
			method: "POST",
			path:   "/api/form",
			formData: url.Values{
				"field1": []string{"value1"},
				"field2": []string{"value2"},
			},
			serverResponse: `{"status":"success"}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify headers
				if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
					t.Error("Expected Content-Type: application/x-www-form-urlencoded")
				}

				// Parse form
				if err := r.ParseForm(); err != nil {
					t.Errorf("Failed to parse form: %v", err)
				}

				// Verify authentication fields were added
				if r.FormValue("userName") != "test_user" {
					t.Error("Expected userName in form data")
				}
				if r.FormValue("licenseKey") != "test_key" {
					t.Error("Expected licenseKey in form data")
				}

				// Verify custom form fields
				for key, values := range tt.formData {
					if formValue := r.FormValue(key); formValue != values[0] {
						t.Errorf("Expected form field '%s' with value '%s', got '%s'", key, values[0], formValue)
					}
				}

				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create client
			client, err := NewClient(server.URL, "test_user", "test_key")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Execute request
			respBody, statusCode, err := client.DoFormRequest(tt.method, tt.path, tt.formData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if statusCode != tt.serverStatus {
					t.Errorf("Expected status code %d, got %d", tt.serverStatus, statusCode)
				}
				if string(respBody) != tt.serverResponse {
					t.Errorf("Expected response body '%s', got '%s'", tt.serverResponse, string(respBody))
				}
			}
		})
	}
}

func TestClientHTTPTimeout(t *testing.T) {
	client, err := NewClient("https://test.insightfinder.com", "test_user", "test_key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client.HTTPClient.Timeout == 0 {
		t.Error("Expected HTTP client timeout to be set")
	}
}
