// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
)

// logLabelMutex ensures log label operations are serialized to prevent race conditions
// when multiple projects are being created/updated simultaneously
var logLabelMutex sync.Mutex

// LogLabelSetting represents a log label configuration
type LogLabelSetting struct {
	ProjectName    string   `json:"projectName"`
	LabelType      string   `json:"labelType"`
	LogLabelString string   `json:"logLabelString"` // JSON array as string
	Labels         []string `json:"-"`              // Parsed from LogLabelString
}

// LogLabelsResponse represents the API response for log labels operations
type LogLabelsResponse struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message,omitempty"`
	Keywords map[string]interface{} `json:"keywords,omitempty"`
}

// GetLogLabels retrieves log labels for a project
func (c *Client) GetLogLabels(projectName, username string) (map[string]string, error) {
	params := url.Values{}
	params.Add("projectName", projectName)

	path := fmt.Sprintf("/api/external/v1/projectkeywords?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil // No log labels found
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get log labels: HTTP %d", statusCode)
	}

	var response LogLabelsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Parse the keywords map - convert arrays of objects to JSON strings
	// Use compact JSON format (no extra whitespace) for consistency
	result := make(map[string]string)
	if response.Keywords != nil {
		for key, value := range response.Keywords {
			// Marshal the value (array of objects) to JSON string with no indent
			if jsonBytes, err := json.Marshal(value); err == nil {
				result[key] = string(jsonBytes)
			}
		}
	}

	return result, nil
}

// CreateOrUpdateLogLabels creates or updates log labels for a project
func (c *Client) CreateOrUpdateLogLabels(projectName, username string, settings []*LogLabelSetting) error {
	// Lock to prevent race conditions when multiple projects are setting labels simultaneously
	logLabelMutex.Lock()
	defer logLabelMutex.Unlock()

	// API endpoint for log labels
	path := fmt.Sprintf("/api/external/v1/watch-tower-setting?projectName=%s&customerName=%s",
		url.QueryEscape(projectName),
		url.QueryEscape(username))

	// Process each log label setting
	for _, setting := range settings {
		// Prepare the request body with logLabelSettingCreate
		requestBody := map[string]interface{}{
			"logLabelSettingCreate": map[string]interface{}{
				"labelType":      setting.LabelType,
				"logLabelString": setting.LogLabelString,
			},
		}

		// Pass the map directly to DoRequest - it will marshal it
		body, statusCode, err := c.DoRequest("POST", path, requestBody)
		if err != nil {
			return err
		}

		if statusCode != 200 {
			return fmt.Errorf("failed to create/update log label: HTTP %d - %s", statusCode, string(body))
		}

		// Check response
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			// If we can't parse but got 200, assume success
			continue
		}

		// Some APIs return success boolean, others don't
		if success, ok := response["success"].(bool); ok && !success {
			if msg, ok := response["message"].(string); ok {
				return fmt.Errorf("log label configuration failed: %s", msg)
			}
			return fmt.Errorf("log label configuration failed")
		}
	}

	return nil
}

// DeleteLogLabels removes log labels for a project
// Note: The API doesn't have a direct delete endpoint, so we set empty arrays
func (c *Client) DeleteLogLabels(projectName, username string, labelTypes []string) error {
	path := fmt.Sprintf("/api/external/v1/watch-tower-setting?projectName=%s&customerName=%s",
		url.QueryEscape(projectName),
		url.QueryEscape(username))

	// For each label type, set an empty array
	for _, labelType := range labelTypes {
		requestBody := map[string]interface{}{
			"logLabelSettingCreate": map[string]interface{}{
				"labelType":      labelType,
				"logLabelString": "[]", // Empty JSON array
			},
		}

		bodyJSON, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal log label delete request: %w", err)
		}

		body, statusCode, err := c.DoRequest("POST", path, bodyJSON)
		if err != nil {
			return err
		}

		// 200 or 204 are acceptable
		if statusCode != 200 && statusCode != 204 {
			return fmt.Errorf("failed to delete log label: HTTP %d - %s", statusCode, string(body))
		}
	}

	return nil
}

// MapLabelTypeToAPIField maps label type to API field name
func MapLabelTypeToAPIField(labelType string) string {
	// Based on the Terraform module mapping
	mapping := map[string]string{
		"whitelist":           "whitelist",
		"trainingWhitelist":   "trainingWhitelist",
		"blacklist":           "trainingBlacklistLabels",
		"featurelist":         "featurelist",
		"incidentlist":        "incidentlist",
		"triagelist":          "triagelist",
		"anomalyFeature":      "anomalyFeatureLabels",
		"dataFilter":          "dataFilterLabels",
		"patternName":         "patternNameLabels",
		"instanceName":        "instanceNameLabels",
		"dataQualityCheck":    "dataQualityCheckLabels",
		"extractionBlacklist": "extractionBlacklist",
	}

	if apiField, ok := mapping[labelType]; ok {
		return apiField
	}
	return labelType // Return as-is if not in mapping
}
