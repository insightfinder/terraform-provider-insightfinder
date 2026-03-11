// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GlobalKBSetting represents the global knowledge base settings for a system
type GlobalKBSetting struct {
	SystemID                       string                   `json:"systemId,omitempty"`
	EnableGlobalKnowledgeBase      bool                     `json:"enableGlobalKnowledgeBase"`
	CompositeValidThreshold        int64                    `json:"compositeValidThreshold"`
	TimelineTopK                   int64                    `json:"timelineTopK"`
	EnableIgnoreInstancePrediction bool                     `json:"enableIgnoreInstancePrediction"`
	PredictionSource               int64                    `json:"predictionSource"`
	ShareSystemType                int64                    `json:"shareSystemType"`
	ActionExecutionTime            int64                    `json:"actionExecutionTime"`
	AutoFixValidationWindow        int64                    `json:"autoFixValidationWindow"`
	FilterSelfToSelf               bool                     `json:"filterSelfToSelf"`
	RuleSourceType                 int64                    `json:"ruleSourceType"`
	// SatelliteSystemSet is returned by the GET API (read format)
	SatelliteSystemSet  []SatelliteSystemSetEntry `json:"satelliteSystemSet,omitempty"`
	// SatelliteSystemList is sent in the POST API (write format)
	SatelliteSystemList []SatelliteSystem         `json:"satelliteSystemList,omitempty"`
}

// SatelliteSystemSetEntry represents one entry in the GET satelliteSystemSet array
type SatelliteSystemSetEntry struct {
	SystemPartitionKey SatelliteSystemPartitionKey `json:"systemPartitionKey"`
	Replay             bool                        `json:"replay"`
}

// SatelliteSystemPartitionKey is the nested key inside SatelliteSystemSetEntry
type SatelliteSystemPartitionKey struct {
	UserName   string `json:"userName"`
	SystemName string `json:"systemName"` // This is the system ID hash
	EnvName    string `json:"envName"`
}

// SatelliteSystem represents one entry in the POST satelliteSystemList array
type SatelliteSystem struct {
	SystemID string `json:"systemId"`
	UserName string `json:"userName"`
}

// SatelliteSystemSetToList converts GET-format satelliteSystemSet entries to POST-format satelliteSystemList entries
func SatelliteSystemSetToList(set []SatelliteSystemSetEntry) []SatelliteSystem {
	list := make([]SatelliteSystem, 0, len(set))
	for _, entry := range set {
		list = append(list, SatelliteSystem{
			SystemID: entry.SystemPartitionKey.SystemName,
			UserName: entry.SystemPartitionKey.UserName,
		})
	}
	return list
}

// IncidentPredictionSetting represents incident prediction settings for a system
type IncidentPredictionSetting struct {
	SystemID                     string  `json:"systemId,omitempty"`
	RuleActiveThreshold          float64 `json:"ruleActiveThreshold"`
	RuleInactiveThreshold        float64 `json:"ruleInactiveThreshold"`
	RuleActiveCondition          int64   `json:"ruleActiveCondition"`
	FalsePositiveTolerance       int64   `json:"falsePositiveTolerance"`
	KBTrainingLength             int64   `json:"kbTrainingLength"`
	Tolerance                    float64 `json:"tolerance"`
	EnableInsensitiveRuleMatching bool   `json:"enableInsensitiveRuleMatching"`
}

// HealthViewKey represents the key structure in health view settings
type HealthViewKey struct {
	SystemPartitionKey struct {
		UserName   string `json:"userName"`
		SystemName string `json:"systemName"`
		EnvName    string `json:"envName"`
	} `json:"systemPartitionKey"`
	Replay bool `json:"replay"`
}

// HealthViewSetting represents the health view / notifications settings for a system
type HealthViewSetting struct {
	Key                                   HealthViewKey          `json:"key"`
	Order                                 int64                  `json:"order"`
	HideFlag                              bool                   `json:"hideFlag"`
	AggregationInterval                   int64                  `json:"aggregationInterval"`
	PredictionEmail                       string                 `json:"predictionEmail"`
	AlertHealthScore                      float64                `json:"alertHealthScore"`
	LastAlertTimestamp                    int64                  `json:"lastAlertTimestamp"`
	EnableSplunkExport                    bool                   `json:"enableSplunkExport"`
	MetricTotal99Percentile               float64                `json:"metricTotal99Percentile"`
	LogTotal99Percentile                  float64                `json:"logTotal99Percentile"`
	SplunkExportTimestamp                 int64                  `json:"splunkExportTimestamp"`
	AlertFrequency                        int64                  `json:"alertFrequency"`
	EmailDampeningPeriod                  int64                  `json:"emailDampeningPeriod"`
	AlertsEmailDampeningPeriod            int64                  `json:"alertsEmailDampeningPeriod"`
	PredictionEmailDampeningPeriod        int64                  `json:"predictionEmailDampeningPeriod"`
	EnableSystemDownEmailAlert            bool                   `json:"enableSystemDownEmailAlert"`
	OnlySendWithRCA                       bool                   `json:"onlySendWithRCA"`
	EnableIncidentPredictionEmailAlert    bool                   `json:"enableIncidentPredictionEmailAlert"`
	EnableIncidentDetectionEmailAlert     bool                   `json:"enableIncidentDetectionEmailAlert"`
	EnableAlertsEmail                     bool                   `json:"enableAlertsEmail"`
	EnableHealthEmailAlert                bool                   `json:"enableHealthEmailAlert"`
	AlertEmail                            string                 `json:"alertEmail"`
	HealthAlertEmail                      string                 `json:"healthAlertEmail"`
	IncidentDetectionEmail                string                 `json:"incidentDetectionEmail"`
	EnableRootCauseEmailAlert             bool                   `json:"enableRootCauseEmailAlert"`
	RootCauseEmail                        string                 `json:"rootCauseEmail"`
	IncidentCountThreshold                map[string]int64       `json:"incidentCountThreshold,omitempty"`
	AssignmentMap                         map[string]any         `json:"assignmentMap"`
	IncidentDampeningWindow               int64                  `json:"incidentDampeningWindow"`
	SystemID                              string                 `json:"systemId,omitempty"`
	ID                                    string                 `json:"id,omitempty"`
}

// GetGlobalKBSetting retrieves the global knowledge base setting for a system
func (c *Client) GetGlobalKBSetting(systemID string) (*GlobalKBSetting, error) {
	systemIDsJSON, err := json.Marshal([]string{systemID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal systemIds: %w", err)
	}

	params := url.Values{}
	params.Add("customerName", c.Username)
	params.Add("systemIds", string(systemIDsJSON))

	path := fmt.Sprintf("/api/external/v1/globalknowledgebasesetting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get global KB setting: HTTP %d - %s", statusCode, string(body))
	}

	var settings []GlobalKBSetting
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse global KB setting: %w", err)
	}

	if len(settings) == 0 {
		return nil, nil
	}

	result := settings[0]
	result.SystemID = systemID
	return &result, nil
}

// SetGlobalKBSetting updates the global knowledge base setting for a system
func (c *Client) SetGlobalKBSetting(systemID string, setting *GlobalKBSetting) error {
	setting.SystemID = systemID

	settingModels := []GlobalKBSetting{*setting}
	settingModelsJSON, err := json.Marshal(settingModels)
	if err != nil {
		return fmt.Errorf("failed to marshal settingModels: %w", err)
	}

	formData := url.Values{}
	formData.Set("customerName", c.Username)
	formData.Set("settingModels", string(settingModelsJSON))
	formData.Set("systemName", systemID)

	body, statusCode, err := c.DoFormRequest("POST", "/api/external/v1/globalknowledgebasesetting", formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set global KB setting: HTTP %d - %s", statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set global KB setting: %s", msg)
		}
		return fmt.Errorf("failed to set global KB setting")
	}

	return nil
}

// GetIncidentPredictionSetting retrieves the incident prediction setting for a system
func (c *Client) GetIncidentPredictionSetting(systemID string) (*IncidentPredictionSetting, error) {
	systemIDsJSON, err := json.Marshal([]string{systemID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal systemIds: %w", err)
	}

	params := url.Values{}
	params.Add("customerName", c.Username)
	params.Add("systemIds", string(systemIDsJSON))

	path := fmt.Sprintf("/api/external/v2/IncidentPredictionSetting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get incident prediction setting: HTTP %d - %s", statusCode, string(body))
	}

	var settings []IncidentPredictionSetting
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse incident prediction setting: %w", err)
	}

	if len(settings) == 0 {
		return nil, nil
	}

	result := settings[0]
	result.SystemID = systemID
	return &result, nil
}

// SetIncidentPredictionSetting updates the incident prediction setting for a system
func (c *Client) SetIncidentPredictionSetting(systemID string, setting *IncidentPredictionSetting) error {
	setting.SystemID = systemID

	systemIDsJSON, err := json.Marshal([]string{systemID})
	if err != nil {
		return fmt.Errorf("failed to marshal systemIds: %w", err)
	}

	settingModels := []IncidentPredictionSetting{*setting}
	settingModelsJSON, err := json.Marshal(settingModels)
	if err != nil {
		return fmt.Errorf("failed to marshal settingModels: %w", err)
	}

	formData := url.Values{}
	formData.Set("customerName", c.Username)
	formData.Set("systemIds", string(systemIDsJSON))
	formData.Set("settingModels", string(settingModelsJSON))
	formData.Set("systemName", systemID)

	body, statusCode, err := c.DoFormRequest("POST", "/api/external/v2/IncidentPredictionSetting", formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set incident prediction setting: HTTP %d - %s", statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set incident prediction setting: %s", msg)
		}
		return fmt.Errorf("failed to set incident prediction setting")
	}

	return nil
}

// GetHealthViewSetting retrieves the health view / notifications setting for a specific system
func (c *Client) GetHealthViewSetting(systemID string) (*HealthViewSetting, error) {
	params := url.Values{}
	params.Add("customerName", c.Username)

	path := fmt.Sprintf("/api/external/v2/healthviewsetting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get health view settings: HTTP %d - %s", statusCode, string(body))
	}

	// Response is a map of systemId -> HealthViewSetting
	var allSettings map[string]HealthViewSetting
	if err := json.Unmarshal(body, &allSettings); err != nil {
		return nil, fmt.Errorf("failed to parse health view settings: %w", err)
	}

	setting, ok := allSettings[systemID]
	if !ok {
		return nil, nil
	}

	setting.SystemID = systemID
	setting.ID = systemID
	return &setting, nil
}

// SetHealthViewSetting updates the notifications settings for a specific system.
// It fetches current settings for all systems, updates only the target system, then posts all back.
func (c *Client) SetHealthViewSetting(systemID string, updates *HealthViewSetting) error {
	// GET current settings for all systems
	params := url.Values{}
	params.Add("customerName", c.Username)

	path := fmt.Sprintf("/api/external/v2/healthviewsetting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to get existing health view settings: HTTP %d - %s", statusCode, string(body))
	}

	var allSettingsMap map[string]HealthViewSetting
	if err := json.Unmarshal(body, &allSettingsMap); err != nil {
		return fmt.Errorf("failed to parse health view settings: %w", err)
	}

	// Get or create the setting for our system
	current, exists := allSettingsMap[systemID]
	if !exists {
		// Initialize with defaults matching the API structure
		current = HealthViewSetting{
			AssignmentMap: make(map[string]any),
		}
		current.Key.SystemPartitionKey.UserName = c.Username
		current.Key.SystemPartitionKey.SystemName = systemID
		current.Key.SystemPartitionKey.EnvName = "All"
	}

	// Apply only the notification fields from updates
	current.Order = updates.Order
	current.HideFlag = updates.HideFlag
	current.AggregationInterval = updates.AggregationInterval
	current.EnableSplunkExport = updates.EnableSplunkExport
	current.IncidentCountThreshold = updates.IncidentCountThreshold
	current.AssignmentMap = updates.AssignmentMap
	current.PredictionEmail = updates.PredictionEmail
	current.AlertHealthScore = updates.AlertHealthScore
	current.AlertFrequency = updates.AlertFrequency
	current.EmailDampeningPeriod = updates.EmailDampeningPeriod
	current.AlertsEmailDampeningPeriod = updates.AlertsEmailDampeningPeriod
	current.PredictionEmailDampeningPeriod = updates.PredictionEmailDampeningPeriod
	current.EnableSystemDownEmailAlert = updates.EnableSystemDownEmailAlert
	current.OnlySendWithRCA = updates.OnlySendWithRCA
	current.EnableIncidentPredictionEmailAlert = updates.EnableIncidentPredictionEmailAlert
	current.EnableIncidentDetectionEmailAlert = updates.EnableIncidentDetectionEmailAlert
	current.EnableAlertsEmail = updates.EnableAlertsEmail
	current.EnableHealthEmailAlert = updates.EnableHealthEmailAlert
	current.AlertEmail = updates.AlertEmail
	current.HealthAlertEmail = updates.HealthAlertEmail
	current.IncidentDetectionEmail = updates.IncidentDetectionEmail
	current.EnableRootCauseEmailAlert = updates.EnableRootCauseEmailAlert
	current.RootCauseEmail = updates.RootCauseEmail
	current.IncidentDampeningWindow = updates.IncidentDampeningWindow
	current.SystemID = systemID
	current.ID = systemID

	// Build settings array: updated target system + all other systems unchanged
	settingsArray := []HealthViewSetting{current}
	for id, s := range allSettingsMap {
		if id != systemID {
			s.SystemID = id
			s.ID = id
			settingsArray = append(settingsArray, s)
		}
	}

	settingsJSON, err := json.Marshal(settingsArray)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	formData := url.Values{}
	formData.Set("systemName", systemID)
	formData.Set("settings", string(settingsJSON))
	formData.Set("customerName", c.Username)

	respBody, statusCode, err := c.DoFormRequest("POST", "/api/external/v2/healthviewsetting", formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set health view settings: HTTP %d - %s", statusCode, string(respBody))
	}

	var response map[string]any
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set health view settings: %s", msg)
		}
		return fmt.Errorf("failed to set health view settings")
	}

	return nil
}
