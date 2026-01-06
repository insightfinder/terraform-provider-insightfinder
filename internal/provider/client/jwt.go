// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// SystemFramework represents a system configuration
type SystemFramework struct {
	SystemKey         SystemKey              `json:"systemKey"`
	SystemID          string                 `json:"systemId"`
	SystemName        string                 `json:"systemName"`
	SystemDisplayName string                 `json:"systemDisplayName"`
	SystemSetting     string                 `json:"systemSetting"`
	EnvironmentArr    []string               `json:"environmentArr"`
	Settings          map[string]interface{} `json:"-"` // Parsed from SystemSetting
}

// SystemKey represents the key structure in system framework
type SystemKey struct {
	UserName        string `json:"userName"`
	SystemName      string `json:"systemName"` // This is the actual system ID (hash)
	EnvironmentName string `json:"environmentName"`
}

// UnmarshalJSON implements custom JSON unmarshaling for SystemFramework
func (s *SystemFramework) UnmarshalJSON(data []byte) error {
	type Alias SystemFramework
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// SystemFrameworkResponse represents the API response for system framework operations
type SystemFrameworkResponse struct {
	Success        bool     `json:"success"`
	Message        string   `json:"message,omitempty"`
	OwnSystemArr   []string `json:"ownSystemArr,omitempty"`
	EnvironmentArr []string `json:"environmentArr,omitempty"`
	ShareSystemArr []string `json:"shareSystemArr,omitempty"`
}

// JWTConfig represents JWT configuration for a system
type JWTConfig struct {
	SystemName string `json:"systemName"`
	SystemID   string `json:"systemId"`
	JWTSecret  string `json:"jwtSecret"`
	JWTType    int    `json:"jwtType"` // 1 for system-level JWT
}

// GetSystemFramework retrieves system framework configuration
func (c *Client) GetSystemFramework(username string, needDetail bool) (*SystemFrameworkResponse, error) {
	params := url.Values{}
	params.Add("customerName", username)
	if needDetail {
		params.Add("needDetail", "true")
	} else {
		params.Add("needDetail", "false")
	}
	params.Add("tzOffset", "0")

	path := fmt.Sprintf("/api/external/v1/systemframework?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get system framework: HTTP %d", statusCode)
	}

	var response SystemFrameworkResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetJWTConfig retrieves JWT configuration for a specific system
func (c *Client) GetJWTConfig(systemName, username string) (*JWTConfig, error) {
	normalizedName := strings.TrimSpace(systemName)
	if normalizedName == "" {
		return nil, fmt.Errorf("system name is required to fetch JWT configuration")
	}

	resolvedIDs, err := c.ResolveSystemNameToIDs([]string{normalizedName}, username)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, nil
		}
		return nil, err
	}

	if len(resolvedIDs) == 0 {
		return nil, nil
	}

	targetID := strings.TrimSpace(resolvedIDs[0])
	if targetID == "" {
		return nil, fmt.Errorf("system '%s' returned empty identifier", normalizedName)
	}

	response, err := c.GetSystemFramework(username, true)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil // No systems found
	}

	systems := make([]string, 0, len(response.OwnSystemArr)+len(response.ShareSystemArr))
	systems = append(systems, response.OwnSystemArr...)
	systems = append(systems, response.ShareSystemArr...)

	var matchingSystem *SystemFramework

searchLoop:
	for _, systemStr := range systems {
		var system SystemFramework
		if err := json.Unmarshal([]byte(systemStr), &system); err != nil {
			continue
		}

		idCandidates := []string{
			strings.TrimSpace(system.SystemKey.SystemName),
			strings.TrimSpace(system.SystemID),
			strings.TrimSpace(system.SystemName),
		}

		for _, candidate := range idCandidates {
			if candidate == "" {
				continue
			}
			if strings.EqualFold(candidate, targetID) {
				systemCopy := system
				matchingSystem = &systemCopy
				break searchLoop
			}
		}
	}

	if matchingSystem == nil {
		return nil, nil
	}

	var settings map[string]interface{}
	if trimmed := strings.TrimSpace(matchingSystem.SystemSetting); trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &settings); err != nil {
			return nil, fmt.Errorf("failed to parse system settings for '%s': %w", normalizedName, err)
		}
	}

	displayName := normalizedName
	if resolvedNames, err := c.ResolveSystemIDsToNames([]string{targetID}, username); err == nil {
		if len(resolvedNames) > 0 && strings.TrimSpace(resolvedNames[0]) != "" {
			displayName = strings.TrimSpace(resolvedNames[0])
		}
	}

	jwtConfig := &JWTConfig{
		SystemName: displayName,
		SystemID:   targetID,
		JWTType:    1,
	}

	if settings != nil {
		if secret, ok := settings["systemLevelJWTSecret"].(string); ok {
			jwtConfig.JWTSecret = secret
		}
		if jwtType, ok := settings["jwtType"].(float64); ok {
			jwtConfig.JWTType = int(jwtType)
		}
	}

	return jwtConfig, nil
}

// CreateOrUpdateJWTConfig creates or updates JWT configuration for a system
func (c *Client) CreateOrUpdateJWTConfig(config *JWTConfig, username string) error {
	if config == nil {
		return fmt.Errorf("config is required")
	}

	if strings.TrimSpace(config.SystemID) == "" {
		trimmedName := strings.TrimSpace(config.SystemName)
		if trimmedName == "" {
			return fmt.Errorf("either system_id or system_name must be provided")
		}

		ids, err := c.ResolveSystemNameToIDs([]string{trimmedName}, username)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return fmt.Errorf("system '%s' not found", trimmedName)
		}
		config.SystemID = strings.TrimSpace(ids[0])
		if config.SystemID == "" {
			return fmt.Errorf("system '%s' returned empty identifier", trimmedName)
		}
	}

	// Prepare the systemKey JSON
	systemKey := map[string]interface{}{
		"userName":        username,
		"systemName":      config.SystemID,
		"environmentName": "All",
	}
	systemKeyJSON, err := json.Marshal(systemKey)
	if err != nil {
		return fmt.Errorf("failed to marshal system key: %w", err)
	}

	// Prepare the systemFrameworkSetting JSON
	systemFrameworkSetting := map[string]interface{}{
		"systemLevelJWTSecret": config.JWTSecret,
		"jwtType":              config.JWTType,
	}
	systemFrameworkSettingJSON, err := json.Marshal(systemFrameworkSetting)
	if err != nil {
		return fmt.Errorf("failed to marshal system framework setting: %w", err)
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("operation", "systemFrameworkSetting")
	formData.Set("systemKey", string(systemKeyJSON))
	formData.Set("systemFrameworkSetting", string(systemFrameworkSettingJSON))

	path := "/api/external/v1/systemframework?tzOffset=0"
	body, statusCode, err := c.DoFormRequest("POST", path, formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to configure JWT: HTTP %d - %s", statusCode, string(body))
	}

	// Check if response indicates success
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		// If we can't parse the response but got 200, assume success
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("JWT configuration failed: %s", msg)
		}
		return fmt.Errorf("JWT configuration failed")
	}

	return nil
}

// DeleteJWTConfig removes JWT configuration from a system
func (c *Client) DeleteJWTConfig(config *JWTConfig, username string) error {
	// To delete, we set an empty JWT secret
	emptyConfig := &JWTConfig{
		SystemName: config.SystemName,
		SystemID:   config.SystemID,
		JWTSecret:  "",
		JWTType:    0,
	}
	return c.CreateOrUpdateJWTConfig(emptyConfig, username)
}
