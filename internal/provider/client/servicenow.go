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
	Account               string     `json:"account"`
	ServiceHost           string     `json:"service_host"`
	Password              string     `json:"password"`
	Proxy                 string     `json:"proxy,omitempty"`
	DampeningPeriod       int        `json:"dampening_period"`
	AppID                 string     `json:"app_id,omitempty"`
	AppKey                string     `json:"app_key,omitempty"`
	AuthType              string     `json:"auth_type,omitempty"`
	SystemIDs             []string   `json:"system_ids"`
	SystemNames           []string   `json:"system_names,omitempty"`
	Options               []string   `json:"options"`
	ContentOption         []string   `json:"content_option"`
	ServiceNowField       string     `json:"service_now_field,omitempty"`
	ContentSource         string     `json:"content_source,omitempty"`
	TriggerWindowInMills  int64      `json:"trigger_window_in_mills,omitempty"`
	EnableFeedbackCollect bool       `json:"enable_feedback_collect,omitempty"`
	EnableTicketCreation  bool       `json:"enable_ticket_creation,omitempty"`
	TableMapping          [][]string `json:"table_mapping,omitempty"`
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
	params.Add("serviceProvider", "ServiceNow")
	params.Add("tzOffset", "-14400000")

	path := fmt.Sprintf("/api/external/v1/system/externalServlies/list?%s", params.Encode())
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

	var response struct {
		ExtServiceAllInfo []map[string]interface{} `json:"extServiceAllInfo"`
		Success           bool                     `json:"success"`
		Message           string                   `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return nil, nil
	}

	// Find the entry matching account and service_host
	normalizeHost := func(h string) string {
		return strings.TrimRight(strings.TrimSpace(h), "/")
	}
	var entry map[string]interface{}
	for _, info := range response.ExtServiceAllInfo {
		entryAccount, _ := info["account"].(string)
		entryHost, _ := info["service_host"].(string)
		if strings.EqualFold(strings.TrimSpace(entryAccount), strings.TrimSpace(account)) &&
			normalizeHost(entryHost) == normalizeHost(serviceHost) {
			entry = info
			break
		}
	}

	if entry == nil {
		return nil, nil // Not found
	}

	config := &ServiceNowConfig{
		Account:     account,
		ServiceHost: serviceHost,
	}

	if pwd, ok := entry["password"].(string); ok {
		config.Password = pwd
	}
	if dampening, ok := entry["dampeningPeriod"].(float64); ok {
		config.DampeningPeriod = int(dampening)
	}
	if appID, ok := entry["appId"].(string); ok {
		config.AppID = appID
	}
	if appKey, ok := entry["appKey"].(string); ok {
		config.AppKey = appKey
	}
	if proxy, ok := entry["proxy"].(string); ok {
		config.Proxy = proxy
	}
	if serviceNowField, ok := entry["serviceNowField"].(string); ok {
		config.ServiceNowField = serviceNowField
	}
	if contentSource, ok := entry["contentSource"].(string); ok {
		config.ContentSource = contentSource
	}
	if enableFeedback, ok := entry["enableServiceNowFeedbackCollect"].(bool); ok {
		config.EnableFeedbackCollect = enableFeedback
	}
	if enableTicket, ok := entry["enableTicketCreation"].(bool); ok {
		config.EnableTicketCreation = enableTicket
	}

	// Determine auth type from appId/appKey presence
	if config.AppID != "" && config.AppKey != "" {
		config.AuthType = "oauth"
	} else {
		config.AuthType = "basic"
	}

	// Parse configs JSON string for systemIds, contentOption, and triggerWindowInMills
	if configsStr, ok := entry["configs"].(string); ok && configsStr != "" {
		var configs map[string]interface{}
		if err := json.Unmarshal([]byte(configsStr), &configs); err == nil {
			if systemIDs, ok := configs["systemIds"].([]interface{}); ok {
				for _, id := range systemIDs {
					if idStr, ok := id.(string); ok {
						config.SystemIDs = append(config.SystemIDs, idStr)
					}
				}
			}
			if contentOpt, ok := configs["contentOption"].([]interface{}); ok {
				for _, opt := range contentOpt {
					if optStr, ok := opt.(string); ok {
						config.ContentOption = append(config.ContentOption, optStr)
					}
				}
			}
			if triggerWindow, ok := configs["triggerWindowInMills"].(float64); ok {
				config.TriggerWindowInMills = int64(triggerWindow)
			}
			// contentSource in configs takes priority if root level is empty
			if config.ContentSource == "" {
				if cs, ok := configs["contentSource"].(string); ok {
					config.ContentSource = cs
				}
			}
			// enableFeedbackCollect in configs only fills in if root-level key was absent
			// (root uses "enableServiceNowFeedbackCollect", configs uses "enableFeedbackCollect")
			if !config.EnableFeedbackCollect {
				if enableFeedback, ok := configs["enableFeedbackCollect"].(bool); ok {
					config.EnableFeedbackCollect = enableFeedback
				}
			}
			// enableTicketCreation in configs only fills in if not already set at root
			if !config.EnableTicketCreation {
				if enableTicket, ok := configs["enableTicketCreation"].(bool); ok {
					config.EnableTicketCreation = enableTicket
				}
			}
		}
	}

	// Parse options JSON array string
	if optionsStr, ok := entry["options"].(string); ok && optionsStr != "" {
		var options []string
		if err := json.Unmarshal([]byte(optionsStr), &options); err == nil {
			config.Options = options
		}
	}

	// Parse tableMapping array of [projectName, tableName] pairs
	if tableMappingRaw, ok := entry["tableMapping"].([]interface{}); ok {
		for _, row := range tableMappingRaw {
			if rowSlice, ok := row.([]interface{}); ok && len(rowSlice) == 2 {
				project, _ := rowSlice[0].(string)
				table, _ := rowSlice[1].(string)
				config.TableMapping = append(config.TableMapping, []string{project, table})
			}
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
	formData.Set("stored_account", config.Account)
	formData.Set("storedHost", config.ServiceHost)
	contentSource := config.ContentSource
	if contentSource == "" {
		contentSource = "work_notes"
	}
	formData.Set("contentSource", contentSource)
	if !verify {
		formData.Set("systemIds", string(systemIDsJSON))
		formData.Set("options", string(optionsJSON))
		formData.Set("contentOption", string(contentOptionJSON))
		if config.ServiceNowField != "" {
			formData.Set("serviceNowField", config.ServiceNowField)
		}
		if config.TriggerWindowInMills > 0 {
			formData.Set("triggerWindowInMills", fmt.Sprintf("%d", config.TriggerWindowInMills))
		}
		formData.Set("enableServiceNowFeedbackCollect", fmt.Sprintf("%t", config.EnableFeedbackCollect))
		formData.Set("enableTicketCreation", fmt.Sprintf("%t", config.EnableTicketCreation))
	}

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

// UpdateServiceNowTableMapping updates the table mapping for a ServiceNow integration
func (c *Client) UpdateServiceNowTableMapping(account, serviceHost, username string, tableMapping [][]string) error {
	mappingJSON, err := json.Marshal(tableMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal table mapping: %w", err)
	}

	formData := url.Values{}
	formData.Set("operation", "ServiceNow")
	formData.Set("customerName", username)
	formData.Set("serviceProvider", "ServiceNow")
	formData.Set("account", account)
	formData.Set("service_host", serviceHost)
	formData.Set("mappingList", string(mappingJSON))

	path := "/api/external/v1/service-integration?tzOffset=-14400000"
	body, statusCode, err := c.DoFormRequest("PUT", path, formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to update ServiceNow table mapping: HTTP %d - %s", statusCode, string(body))
	}

	var response ServiceNowResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if !response.Success {
		return fmt.Errorf("ServiceNow table mapping update failed: %s", response.Message)
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
