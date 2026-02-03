// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetServiceNowThirdPartySettings(t *testing.T) {
	tests := []struct {
		name           string
		projectName    string
		statusCode     int
		responseBody   string
		expectError    bool
		expectNil      bool
		expectedHost   string
		expectedUser   string
		expectedFields int
	}{
		{
			name:        "successful get",
			projectName: "test-project",
			statusCode:  200,
			responseBody: `{
				"host": "https://dev123456.service-now.com/",
				"sysparmQuery": "",
				"proxy": "",
				"serviceNowUser": "admin",
				"serviceNowPassword": "password123",
				"instanceField": "short_description",
				"instanceFieldRegex": "1",
				"timestampFormat": "yyyy-MM-dd HH:mm:ss",
				"clientId": "client-id-123",
				"clientSecret": "client-secret-456",
				"additionalFields": ["work_end", "priority"],
				"success": true
			}`,
			expectError:    false,
			expectNil:      false,
			expectedHost:   "https://dev123456.service-now.com/",
			expectedUser:   "admin",
			expectedFields: 2,
		},
		{
			name:         "settings not found",
			projectName:  "nonexistent-project",
			statusCode:   404,
			responseBody: `{"success": false}`,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "settings don't exist - success false",
			projectName:  "test-project",
			statusCode:   200,
			responseBody: `{"success": false}`,
			expectError:  false,
			expectNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.Header.Get("X-User-Name") != "test_user" {
					t.Errorf("Expected X-User-Name header")
				}
				if r.Header.Get("X-API-Key") != "test_key" {
					t.Errorf("Expected X-API-Key header")
				}

				// Check query parameters
				params := r.URL.Query()
				if params.Get("projectName") != tt.projectName {
					t.Errorf("Expected projectName=%s, got %s", tt.projectName, params.Get("projectName"))
			}
			if params.Get("cloudType") != "ServiceNow" {
				t.Errorf("Expected cloudType=ServiceNow, got %s", params.Get("cloudType"))
			}

			w.WriteHeader(tt.statusCode)
			_, _ = w.Write([]byte(tt.responseBody))
		}))
		defer server.Close()
			client, _ := NewClient(server.URL, "test_user", "test_key")
			settings, err := client.GetServiceNowThirdPartySettings(tt.projectName)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && settings != nil {
				t.Error("Expected nil settings but got non-nil")
			}
			if !tt.expectNil && settings == nil && !tt.expectError {
				t.Error("Expected non-nil settings but got nil")
			}
			if settings != nil {
				if settings.Host != tt.expectedHost {
					t.Errorf("Expected host=%s, got %s", tt.expectedHost, settings.Host)
				}
				if settings.ServiceNowUser != tt.expectedUser {
					t.Errorf("Expected user=%s, got %s", tt.expectedUser, settings.ServiceNowUser)
				}
				if len(settings.AdditionalFields) != tt.expectedFields {
					t.Errorf("Expected %d additional fields, got %d", tt.expectedFields, len(settings.AdditionalFields))
				}
			}
		})
	}
}

func TestCreateOrUpdateServiceNowThirdPartySettings(t *testing.T) {
	tests := []struct {
		name           string
		projectName    string
		settings       *ServiceNowThirdPartySettings
		statusCode     int
		responseBody   string
		expectError    bool
		expectedParams map[string]string
	}{
		{
			name:        "successful create/update",
			projectName: "test-project",
			settings: &ServiceNowThirdPartySettings{
				Host:               "https://dev123456.service-now.com/",
				SysparmQuery:       "",
				Proxy:              "",
				ServiceNowUser:     "admin",
				ServiceNowPassword: "password123",
				InstanceField:      "short_description",
				InstanceFieldRegex: "1",
				TimestampFormat:    "yyyy-MM-dd HH:mm:ss",
				ClientID:           "client-id-123",
				ClientSecret:       "client-secret-456",
				AdditionalFields:   []string{"work_end", "priority"},
			},
			statusCode:   200,
			responseBody: `{"success": true, "message": "Setting updated successfully"}`,
			expectError:  false,
			expectedParams: map[string]string{
				"projectName":        "test-project",
				"cloudType":          "ServiceNow",
				"host":               "https://dev123456.service-now.com/",
				"serviceNowUser":     "admin",
				"serviceNowPassword": "password123",
			},
		},
		{
			name:        "api error",
			projectName: "test-project",
			settings: &ServiceNowThirdPartySettings{
				Host:               "https://dev123456.service-now.com/",
				ServiceNowUser:     "admin",
				ServiceNowPassword: "password123",
				ClientID:           "client-id-123",
				ClientSecret:       "client-secret-456",
			},
			statusCode:   500,
			responseBody: `{"success": false, "message": "Internal server error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Check query parameters
				params := r.URL.Query()
				for key, expectedValue := range tt.expectedParams {
					if params.Get(key) != expectedValue {
						t.Errorf("Expected %s=%s, got %s", key, expectedValue, params.Get(key))
					}
				}

				// Verify additional fields is JSON encoded
			if additionalFields := params.Get("additionalFields"); additionalFields != "" {
				var fields []string
				if err := json.Unmarshal([]byte(additionalFields), &fields); err != nil {
					t.Errorf("additionalFields is not valid JSON: %v", err)
				}
			}

			w.WriteHeader(tt.statusCode)
			_, _ = w.Write([]byte(tt.responseBody))
		}))
		defer server.Close()
			client, _ := NewClient(server.URL, "test_user", "test_key")
			err := client.CreateOrUpdateServiceNowThirdPartySettings(tt.projectName, tt.settings)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestServiceNowThirdPartySettings_URLEncoding(t *testing.T) {
	settings := &ServiceNowThirdPartySettings{
		Host:               "https://dev123456.service-now.com/",
		ServiceNowPassword: "P@ssw0rd!#$",
		SysparmQuery:       "state=1^priority=2",
		AdditionalFields:   []string{"work_end", "priority"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		// Check that special characters are properly encoded
		if password := params.Get("serviceNowPassword"); password != settings.ServiceNowPassword {
			t.Errorf("Password not properly encoded. Expected %s, got %s", settings.ServiceNowPassword, password)
		}


	if query := params.Get("sysparmQuery"); query != settings.SysparmQuery {
		t.Errorf("sysparmQuery not properly encoded. Expected %s, got %s", settings.SysparmQuery, query)
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte(`{"success": true}`))
}))
defer server.Close()
	client, _ := NewClient(server.URL, "test_user", "test_key")
	err := client.CreateOrUpdateServiceNowThirdPartySettings("test-project", settings)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
