// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ServiceNowConfig represents ServiceNow integration configuration
type ServiceNowConfig struct {
	Account         string   `json:"account"`
	ServiceHost     string   `json:"service_host"`
	Password        string   `json:"password"`
	Proxy           string   `json:"proxy,omitempty"`
	DampeningPeriod int      `json:"dampening_period"`
	AppID           string   `json:"app_id,omitempty"`
	AppKey          string   `json:"app_key,omitempty"`
	AuthType        string   `json:"auth_type,omitempty"`
	SystemIDs       []string `json:"system_ids"`
	SystemNames     []string `json:"system_names,omitempty"`
	Options         []string `json:"options"`
	ContentOption   []string `json:"content_option"`
}

// ServiceNowResponse represents the API response for ServiceNow operations
type ServiceNowResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetServiceNowConfig retrieves ServiceNow integration configuration
func (c *Client) GetServiceNowConfig(account, serviceHost, username string) (*ServiceNowConfig, error) {
	params := url.Values{}
	params.Add("tzOffset", "0")
	params.Add("account", account)
	params.Add("customerName", username)
	params.Add("serviceProvider", "ServiceNow")
	params.Add("operation", "display")
	params.Add("service_host", serviceHost)

	path := fmt.Sprintf("/api/external/v1/service-integration?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil // Configuration doesn't exist
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get ServiceNow config: HTTP %d", statusCode)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the key exists - if not, configuration doesn't exist
	if _, ok := response["key"]; !ok {
		return nil, nil // Configuration doesn't exist
	}

	// Parse the configuration
	config := &ServiceNowConfig{
		Account:     account,
		ServiceHost: serviceHost,
	}

	// Extract configuration fields from response (root level)
	if pwd, ok := response["password"].(string); ok {
		config.Password = pwd
	}
	if dampening, ok := response["dampeningPeriod"].(float64); ok {
		config.DampeningPeriod = int(dampening)
	}
	if appID, ok := response["appId"].(string); ok {
		config.AppID = appID
	}
	if appKey, ok := response["appKey"].(string); ok {
		config.AppKey = appKey
	}
	if authType, ok := response["authType"].(string); ok && authType != "" {
		config.AuthType = strings.ToLower(authType)
	}

	// Parse serviceNowIntegrationConfig JSON string
	if integrationConfigStr, ok := response["serviceNowIntegrationConfig"].(string); ok && integrationConfigStr != "" {
		var integrationConfig map[string]interface{}
		if err := json.Unmarshal([]byte(integrationConfigStr), &integrationConfig); err == nil {
			if systemIDs, ok := integrationConfig["systemIds"].([]interface{}); ok {
				for _, id := range systemIDs {
					if idStr, ok := id.(string); ok {
						config.SystemIDs = append(config.SystemIDs, idStr)
					}
				}
			}
			if contentOpt, ok := integrationConfig["contentOption"].([]interface{}); ok {
				for _, opt := range contentOpt {
					if optStr, ok := opt.(string); ok {
						config.ContentOption = append(config.ContentOption, optStr)
					}
				}
			}
			if systemNames, ok := integrationConfig["systemNames"].([]interface{}); ok {
				for _, name := range systemNames {
					if nameStr, ok := name.(string); ok {
						config.SystemNames = append(config.SystemNames, nameStr)
					}
				}
			}
		}
	}

	// Parse options JSON array string
	if optionsStr, ok := response["options"].(string); ok && optionsStr != "" {
		var options []string
		if err := json.Unmarshal([]byte(optionsStr), &options); err == nil {
			config.Options = options
		}
	}

	if len(config.SystemNames) == 0 && len(config.SystemIDs) > 0 {
		if names, err := c.ResolveSystemIDsToNames(config.SystemIDs, username); err == nil {
			config.SystemNames = names
		}
	}

	return config, nil
}

// CreateOrUpdateServiceNowConfig creates or updates ServiceNow integration
func (c *Client) CreateOrUpdateServiceNowConfig(config *ServiceNowConfig, username string, verify bool) error {
	// Format system IDs as JSON array string
	systemIDsJSON, err := json.Marshal(config.SystemIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal system IDs: %w", err)
	}

	// Format options as JSON array string
	optionsJSON, err := json.Marshal(config.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	// Format content options as JSON array string
	contentOptionJSON, err := json.Marshal(config.ContentOption)
	if err != nil {
		return fmt.Errorf("failed to marshal content options: %w", err)
	}

	if config.AuthType == "" {
		config.AuthType = "basic"
	}
	// Prepare form data
	formData := url.Values{}
	if verify {
		formData.Set("verify", "true")
	}
	formData.Set("operation", "ServiceNow")
	formData.Set("service_host", config.ServiceHost)
	formData.Set("proxy", config.Proxy)
	formData.Set("account", config.Account)
	formData.Set("password", config.Password)
	formData.Set("dampeningPeriod", fmt.Sprintf("%d", config.DampeningPeriod))
	formData.Set("appId", config.AppID)
	formData.Set("appKey", config.AppKey)
	formData.Set("auth_type", config.AuthType)
	formData.Set("customerName", username)
	formData.Set("systemIds", string(systemIDsJSON))
	formData.Set("options", string(optionsJSON))
	formData.Set("contentOption", string(contentOptionJSON))

	path := "/api/external/v1/service-integration"
	body, statusCode, err := c.DoFormRequest("POST", path, formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to configure ServiceNow: HTTP %d - %s", statusCode, string(body))
	}

	// Check if response indicates success
	var response ServiceNowResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// If we can't parse the response but got 200, assume success
		return nil
	}

	if !response.Success {
		if response.Message != "" {
			return fmt.Errorf("ServiceNow configuration failed: %s", response.Message)
		}
		return fmt.Errorf("ServiceNow configuration failed")
	}

	return nil
}

// DeleteServiceNowConfig removes ServiceNow integration
func (c *Client) DeleteServiceNowConfig(account, serviceHost, username string) error {
	serviceHost = strings.TrimSpace(serviceHost)
	if serviceHost == "" {
		return fmt.Errorf("service_host is required for deletion")
	}

	serviceID := fmt.Sprintf("ServiceNow:%s:%s", account, serviceHost)

	formData := url.Values{}
	formData.Set("serviceProvider", "PagerDuty")
	formData.Set("operation", "delete")
	formData.Set("service_id", serviceID)
	formData.Set("serviceOwner", username)
	formData.Set("customerName", username)

	path := "/api/external/v1/service-integration"
	body, statusCode, err := c.DoFormRequest("POST", path, formData)
	if err != nil {
		return err
	}

	// 200 or 404 are both acceptable for deletion
	if statusCode != 200 && statusCode != 404 {
		return fmt.Errorf("failed to delete ServiceNow config: HTTP %d - %s", statusCode, string(body))
	}

	return nil
}

// ResolveSystemNameToIDs converts system names to system IDs
func (c *Client) ResolveSystemNameToIDs(systemNames []string, username string) ([]string, error) {
	if len(systemNames) == 0 {
		return []string{}, nil
	}

	systemFramework, err := c.GetSystemFramework(username, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get system framework: %w", err)
	}

	if systemFramework == nil {
		return nil, fmt.Errorf("no systems found")
	}

	systems := make([]string, 0, len(systemFramework.OwnSystemArr)+len(systemFramework.ShareSystemArr))
	systems = append(systems, systemFramework.OwnSystemArr...)
	systems = append(systems, systemFramework.ShareSystemArr...)

	if len(systems) == 0 {
		return nil, fmt.Errorf("no systems found")
	}

	nameToID := make(map[string]string)

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

		var resolvedID string
		for _, candidate := range idCandidates {
			if candidate != "" {
				resolvedID = candidate
				break
			}
		}
		if resolvedID == "" {
			continue
		}

		nameCandidates := []string{
			strings.TrimSpace(system.SystemDisplayName),
			strings.TrimSpace(system.SystemName),
		}

		for _, candidate := range nameCandidates {
			if candidate == "" {
				continue
			}

			normalized := strings.ToLower(candidate)
			if _, exists := nameToID[normalized]; !exists {
				nameToID[normalized] = resolvedID
			}
		}
	}

	if len(nameToID) == 0 {
		return nil, fmt.Errorf("no systems found")
	}

	systemIDs := make([]string, 0, len(systemNames))
	missing := make([]string, 0)

	for _, name := range systemNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			return nil, fmt.Errorf("system name cannot be empty")
		}

		normalized := strings.ToLower(trimmed)
		if id, ok := nameToID[normalized]; ok {
			systemIDs = append(systemIDs, id)
		} else {
			missing = append(missing, trimmed)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("system(s) not found: %s", strings.Join(missing, ", "))
	}

	return systemIDs, nil
}

// ResolveSystemIDsToNames converts system IDs to system names
func (c *Client) ResolveSystemIDsToNames(systemIDs []string, username string) ([]string, error) {
	if len(systemIDs) == 0 {
		return []string{}, nil
	}

	systemFramework, err := c.GetSystemFramework(username, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get system framework: %w", err)
	}

	if systemFramework == nil {
		return nil, fmt.Errorf("no systems found")
	}

	systems := make([]string, 0, len(systemFramework.OwnSystemArr)+len(systemFramework.ShareSystemArr))
	systems = append(systems, systemFramework.OwnSystemArr...)
	systems = append(systems, systemFramework.ShareSystemArr...)

	if len(systems) == 0 {
		return nil, fmt.Errorf("no systems found")
	}

	systemNames := make([]string, 0, len(systemIDs))
	idToName := make(map[string]string)

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

		var resolvedID string
		for _, candidate := range idCandidates {
			if candidate != "" {
				resolvedID = candidate
				break
			}
		}
		if resolvedID == "" {
			continue
		}

		displayName := strings.TrimSpace(system.SystemDisplayName)
		if displayName == "" {
			displayName = strings.TrimSpace(system.SystemName)
		}
		if displayName == "" {
			continue
		}

		idToName[resolvedID] = displayName
	}

	for _, id := range systemIDs {
		trimmedID := strings.TrimSpace(id)
		if trimmedID == "" {
			continue
		}

		if name, ok := idToName[trimmedID]; ok {
			systemNames = append(systemNames, name)
		} else {
			systemNames = append(systemNames, trimmedID)
		}
	}

	return systemNames, nil
}
