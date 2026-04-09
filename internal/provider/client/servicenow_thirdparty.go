// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ServiceNowThirdPartySettings represents ServiceNow third-party settings
type ServiceNowThirdPartySettings struct {
	Host                 string   `json:"host"`
	SysparmQuery         string   `json:"sysparmQuery"`
	Proxy                string   `json:"proxy"`
	ServiceNowUser       string   `json:"serviceNowUser"`
	ServiceNowPassword   string   `json:"serviceNowPassword"`
	InstanceField        string   `json:"instanceField"`
	InstanceFieldRegex   string   `json:"instanceFieldRegex"`
	TimestampFormat      string   `json:"timestampFormat"`
	ClientID             string   `json:"clientId"`
	ClientSecret         string   `json:"clientSecret"`
	AdditionalFields     []string `json:"additionalFields"`
	DefaultFields        []string `json:"defaultFields"`
	Fields               []string `json:"fields"`
	ComponentNameRule    string   `json:"componentNameRule"`
	ServiceNowImportFlag bool     `json:"serviceNowImportFlag"`
}

// ServiceNowThirdPartyResponse represents the API response
type ServiceNowThirdPartyResponse struct {
	Success bool                          `json:"success"`
	Message string                        `json:"message,omitempty"`
	Data    *ServiceNowThirdPartySettings `json:"-"` // Embedded in root
}

// doServiceNowThirdPartyRequest performs an HTTP request with X-License-Key header for ServiceNow third-party API
func (c *Client) doServiceNowThirdPartyRequest(method, path string) ([]byte, int, error) {
	reqURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers - ServiceNow third-party API uses X-License-Key
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

// GetServiceNowThirdPartySettings retrieves ServiceNow third-party settings for a project
func (c *Client) GetServiceNowThirdPartySettings(projectName string) (*ServiceNowThirdPartySettings, error) {
	params := url.Values{}
	params.Add("projectName", projectName)
	params.Add("cloudType", "ServiceNow")
	params.Add("tzOffset", "-18000000") // Default offset

	path := fmt.Sprintf("/api/external/v1/thirdpartysetting?%s", params.Encode())
	body, statusCode, err := c.doServiceNowThirdPartyRequest("GET", path)
	if err != nil {
		return nil, fmt.Errorf("failed to get ServiceNow third-party settings: %w", err)
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil // Configuration doesn't exist
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get ServiceNow third-party settings: HTTP %d, body: %s", statusCode, string(body))
	}

	// Parse the response - the settings are embedded in the root object
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if success field exists and is true
	if success, ok := rawResponse["success"].(bool); ok && !success {
		// Return nil if settings don't exist
		return nil, nil
	}

	// Parse into settings struct
	settings := &ServiceNowThirdPartySettings{}
	if host, ok := rawResponse["host"].(string); ok {
		settings.Host = host
	}
	if sysparmQuery, ok := rawResponse["sysparmQuery"].(string); ok {
		settings.SysparmQuery = sysparmQuery
	}
	if proxy, ok := rawResponse["proxy"].(string); ok {
		settings.Proxy = proxy
	}
	if serviceNowUser, ok := rawResponse["serviceNowUser"].(string); ok {
		settings.ServiceNowUser = serviceNowUser
	}
	if serviceNowPassword, ok := rawResponse["serviceNowPassword"].(string); ok {
		settings.ServiceNowPassword = serviceNowPassword
	}
	if instanceField, ok := rawResponse["instanceField"].(string); ok {
		settings.InstanceField = instanceField
	}
	if instanceFieldRegex, ok := rawResponse["instanceFieldRegex"].(string); ok {
		settings.InstanceFieldRegex = instanceFieldRegex
	}
	if timestampFormat, ok := rawResponse["timestampFormat"].(string); ok {
		settings.TimestampFormat = timestampFormat
	}
	if clientID, ok := rawResponse["clientId"].(string); ok {
		settings.ClientID = clientID
	}
	if clientSecret, ok := rawResponse["clientSecret"].(string); ok {
		settings.ClientSecret = clientSecret
	}
	if componentNameRule, ok := rawResponse["componentNameRule"].(string); ok {
		settings.ComponentNameRule = componentNameRule
	}
	if serviceNowImportFlag, ok := rawResponse["serviceNowImportFlag"].(bool); ok {
		settings.ServiceNowImportFlag = serviceNowImportFlag
	}

	// Parse array fields
	if additionalFields, ok := rawResponse["additionalFields"].([]interface{}); ok {
		for _, field := range additionalFields {
			if fieldStr, ok := field.(string); ok {
				settings.AdditionalFields = append(settings.AdditionalFields, fieldStr)
			}
		}
	}
	if defaultFields, ok := rawResponse["defaultFields"].([]interface{}); ok {
		for _, field := range defaultFields {
			if fieldStr, ok := field.(string); ok {
				settings.DefaultFields = append(settings.DefaultFields, fieldStr)
			}
		}
	}
	if fields, ok := rawResponse["fields"].([]interface{}); ok {
		for _, field := range fields {
			if fieldStr, ok := field.(string); ok {
				settings.Fields = append(settings.Fields, fieldStr)
			}
		}
	}

	return settings, nil
}

// CreateOrUpdateServiceNowThirdPartySettings creates or updates ServiceNow third-party settings
func (c *Client) CreateOrUpdateServiceNowThirdPartySettings(projectName string, settings *ServiceNowThirdPartySettings) error {
	params := url.Values{}
	params.Add("projectName", projectName)
	params.Add("cloudType", "ServiceNow")
	params.Add("host", settings.Host)
	params.Add("proxy", settings.Proxy)
	params.Add("serviceNowUser", settings.ServiceNowUser)
	params.Add("serviceNowPassword", settings.ServiceNowPassword)
	params.Add("clientId", settings.ClientID)
	params.Add("clientSecret", settings.ClientSecret)
	params.Add("sysparmQuery", settings.SysparmQuery)
	params.Add("timestampFormat", settings.TimestampFormat)
	params.Add("instanceField", settings.InstanceField)
	params.Add("instanceFieldRegex", settings.InstanceFieldRegex)
	params.Add("useHostName", "false")
	params.Add("componentNameRule", settings.ComponentNameRule)
	if settings.ServiceNowImportFlag {
		params.Add("serviceNowImportFlag", "true")
	} else {
		params.Add("serviceNowImportFlag", "false")
	}

	// Encode additional fields as JSON array
	if settings.AdditionalFields != nil {
		additionalFieldsJSON, err := json.Marshal(settings.AdditionalFields)
		if err != nil {
			return fmt.Errorf("failed to marshal additional fields: %w", err)
		}
		params.Add("additionalFields", string(additionalFieldsJSON))
	} else {
		params.Add("additionalFields", "[]")
	}

	path := fmt.Sprintf("/api/external/v1/thirdpartysetting?%s", params.Encode())
	body, statusCode, err := c.doServiceNowThirdPartyRequest("POST", path)
	if err != nil {
		return fmt.Errorf("failed to create/update ServiceNow third-party settings: %w", err)
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to create/update ServiceNow third-party settings: HTTP %d, body: %s", statusCode, string(body))
	}

	var response ServiceNowThirdPartyResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("failed to create/update ServiceNow third-party settings: %s", response.Message)
	}

	return nil
}
