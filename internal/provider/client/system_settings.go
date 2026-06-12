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
	SystemID                       string `json:"systemId,omitempty"`
	EnableGlobalKnowledgeBase      bool   `json:"enableGlobalKnowledgeBase"`
	CompositeValidThreshold        int64  `json:"compositeValidThreshold"`
	TimelineTopK                   int64  `json:"timelineTopK"`
	EnableIgnoreInstancePrediction bool   `json:"enableIgnoreInstancePrediction"`
	PredictionSource               int64  `json:"predictionSource"`
	ShareSystemType                int64  `json:"shareSystemType"`
	ActionExecutionTime            int64  `json:"actionExecutionTime"`
	AutoFixValidationWindow        int64  `json:"autoFixValidationWindow"`
	FilterSelfToSelf               bool   `json:"filterSelfToSelf"`
	RuleSourceType                 int64  `json:"ruleSourceType"`
	// SatelliteSystemSet is returned by the GET API (read format)
	SatelliteSystemSet []SatelliteSystemSetEntry `json:"satelliteSystemSet,omitempty"`
	// SatelliteSystemList is sent in the POST API (write format)
	SatelliteSystemList []SatelliteSystem `json:"satelliteSystemList,omitempty"`
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
	SystemID                      string  `json:"systemId,omitempty"`
	RuleActiveThreshold           float64 `json:"ruleActiveThreshold"`
	RuleInactiveThreshold         float64 `json:"ruleInactiveThreshold"`
	RuleActiveCondition           int64   `json:"ruleActiveCondition"`
	FalsePositiveTolerance        int64   `json:"falsePositiveTolerance"`
	KBTrainingLength              int64   `json:"kbTrainingLength"`
	Tolerance                     float64 `json:"tolerance"`
	EnableInsensitiveRuleMatching bool    `json:"enableInsensitiveRuleMatching"`
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

// CustomConsolidationRule represents one custom incident consolidation rule
type CustomConsolidationRule struct {
	ProjectEntries    []ProjectEntry     `json:"projectEntries"`
	FieldCorrelations []FieldCorrelation `json:"fieldCorrelations"`
}

// ProjectEntry represents a project with its matching conditions in a custom consolidation rule
type ProjectEntry struct {
	ProjectName string      `json:"projectName"`
	Conditions  []Condition `json:"conditions"`
}

// Condition represents a matching condition within a project entry
type Condition struct {
	Type    string `json:"type"`
	Keyword string `json:"keyword"`
}

// FieldCorrelation represents a set of correlated field keys across projects
type FieldCorrelation struct {
	ProjectFieldKeys []ProjectFieldKey `json:"projectFieldKeys"`
}

// ProjectFieldKey represents one project's field key in a field correlation
type ProjectFieldKey struct {
	ProjectName string  `json:"projectName"`
	Type        string  `json:"type"`
	FieldKey    *string `json:"fieldKey"`
}

// MetricLogConsolidationConfig represents a metric-to-log project consolidation mapping
type MetricLogConsolidationConfig struct {
	MetricProjectName string   `json:"metricProjectName"`
	LogProjectName    string   `json:"logProjectName"`
	FieldKeys         []string `json:"fieldKeys"`
}

// HealthViewSetting represents the health view / notifications settings for a system
type HealthViewSetting struct {
	Key                                 HealthViewKey                  `json:"key"`
	Order                               int64                          `json:"order"`
	HideFlag                            bool                           `json:"hideFlag"`
	AggregationInterval                 int64                          `json:"aggregationInterval"`
	PredictionEmail                     string                         `json:"predictionEmail"`
	AlertHealthScore                    float64                        `json:"alertHealthScore"`
	LastAlertTimestamp                  int64                          `json:"lastAlertTimestamp"`
	EnableSplunkExport                  bool                           `json:"enableSplunkExport"`
	MetricTotal99Percentile             float64                        `json:"metricTotal99Percentile"`
	LogTotal99Percentile                float64                        `json:"logTotal99Percentile"`
	SplunkExportTimestamp               int64                          `json:"splunkExportTimestamp"`
	AlertFrequency                      int64                          `json:"alertFrequency"`
	EmailDampeningPeriod                int64                          `json:"emailDampeningPeriod"`
	AlertsEmailDampeningPeriod          int64                          `json:"alertsEmailDampeningPeriod"`
	PredictionEmailDampeningPeriod      int64                          `json:"predictionEmailDampeningPeriod"`
	EnableSystemDownEmailAlert          bool                           `json:"enableSystemDownEmailAlert"`
	OnlySendWithRCA                     bool                           `json:"onlySendWithRCA"`
	EnableIncidentPredictionEmailAlert  bool                           `json:"enableIncidentPredictionEmailAlert"`
	EnableIncidentDetectionEmailAlert   bool                           `json:"enableIncidentDetectionEmailAlert"`
	EnableAlertsEmail                   bool                           `json:"enableAlertsEmail"`
	EnableHealthEmailAlert              bool                           `json:"enableHealthEmailAlert"`
	AlertEmail                          string                         `json:"alertEmail"`
	HealthAlertEmail                    string                         `json:"healthAlertEmail"`
	IncidentDetectionEmail              string                         `json:"incidentDetectionEmail"`
	EnableRootCauseEmailAlert           bool                           `json:"enableRootCauseEmailAlert"`
	RootCauseEmail                      string                         `json:"rootCauseEmail"`
	IncidentCountThreshold              map[string]int64               `json:"incidentCountThreshold,omitempty"`
	AssignmentMap                       map[string]any                 `json:"assignmentMap"`
	IncidentDampeningWindow             int64                          `json:"incidentDampeningWindow"`
	TicketOpenTime                      int64                          `json:"ticketOpenTime"`
	ComponentLevelIncidentConsolidation bool                           `json:"componentLevelIncidentConsolidation"`
	EnabledConsolidationAlgorithms      []string                       `json:"enabledConsolidationAlgorithms"`
	MaxNotificationDelayTolerance       int64                          `json:"maxNotificationDelayTolerance"`
	ProjectLevelDampeningWindows        []ProjectLevelDampeningWindow  `json:"projectLevelDampeningWindows"`
	CustomConsolidationRules            []CustomConsolidationRule      `json:"customConsolidationRules"`
	MetricLogConsolidationConfigs       []MetricLogConsolidationConfig `json:"metricLogConsolidationConfigs"`
	SystemID                            string                         `json:"systemId,omitempty"`
	ID                                  string                         `json:"id,omitempty"`
}

// ProjectLevelDampeningWindow represents a project-level dampening window entry
type ProjectLevelDampeningWindow struct {
	SourceProject  string `json:"ps"`
	TargetProject  string `json:"pt"`
	SourceCustomer string `json:"cs"`
	TargetCustomer string `json:"ct"`
	Duration       int64  `json:"d"`
}

// SystemDownSetting represents system down notification settings for a system
type SystemDownSetting struct {
	EnableSystemDownEmailAlert bool     `json:"enableSystemDownEmailAlert"`
	EmailDampeningPeriod       int64    `json:"emailDampeningPeriod"`
	EmailSet                   []string `json:"emailSet,omitempty"`
}

// systemDownSettingModel is the POST body model for system down settings
type systemDownSettingModel struct {
	EnableSystemDownEmailAlert bool     `json:"enableSystemDownEmailAlert"`
	EmailDampeningPeriod       int64    `json:"emailDampeningPeriod"`
	SystemID                   string   `json:"systemId"`
	EmailSet                   []string `json:"emailSet,omitempty"`
}

// InsightsReportSetting represents insights report notification settings (daily + weekly)
type InsightsReportSetting struct {
	EnableDailyInsightsReport  bool     `json:"enableDailyInsightsReport"`
	EmailSet                   []string `json:"emailSet"`
	WeeklyEmailSet             []string `json:"weeklyEmailSet"`
	EnableWeeklyInsightsReport bool     `json:"enableWeeklyInsightsReport"`
}

// insightsReportSettingModel is the POST body model for insights report settings
type insightsReportSettingModel struct {
	SystemID             string   `json:"systemId"`
	EmailSet             []string `json:"emailSet,omitempty"`
	EnableInsightsReport bool     `json:"enableInsightsReport"`
}

// InstanceDownSetting represents instance down notification settings for a project
type InstanceDownSetting struct {
	ProjectName              string   `json:"projectName"`
	InstanceDownEnable       bool     `json:"instanceDownEnable"`
	InstanceDownDampening    int64    `json:"instanceDownDampening"`
	InstanceDownThreshold    int64    `json:"instanceDownThreshold"`
	InstanceDownReportNumber int64    `json:"instanceDownReportNumber"`
	InstanceDownEmails       []string `json:"instanceDownEmails"`
}

// GetSystemDownSetting retrieves system down notification settings for a system
func (c *Client) GetSystemDownSetting(systemID string) (*SystemDownSetting, error) {
	systemIDsJSON, err := json.Marshal([]string{systemID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal systemIds: %w", err)
	}

	params := url.Values{}
	params.Add("customerName", c.Username)
	params.Add("systemIds", string(systemIDsJSON))

	path := fmt.Sprintf("/api/external/v2/systemdownsetting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get system down setting: HTTP %d - %s", statusCode, string(body))
	}

	var settings []SystemDownSetting
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse system down setting: %w", err)
	}

	if len(settings) == 0 {
		return nil, nil
	}

	return &settings[0], nil
}

// SetSystemDownSetting updates system down notification settings for a system
func (c *Client) SetSystemDownSetting(systemID string, setting *SystemDownSetting) error {
	models := []systemDownSettingModel{{
		EnableSystemDownEmailAlert: setting.EnableSystemDownEmailAlert,
		EmailDampeningPeriod:       setting.EmailDampeningPeriod,
		SystemID:                   systemID,
		EmailSet:                   setting.EmailSet,
	}}

	modelsJSON, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal settingModels: %w", err)
	}

	formData := url.Values{}
	formData.Set("customerName", c.Username)
	formData.Set("settingModels", string(modelsJSON))

	body, statusCode, err := c.DoFormRequest("POST", "/api/external/v2/systemdownsetting", formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set system down setting: HTTP %d - %s", statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set system down setting: %s", msg)
		}
		return fmt.Errorf("failed to set system down setting")
	}

	return nil
}

// GetInsightsReportSetting retrieves insights report notification settings for a system
func (c *Client) GetInsightsReportSetting(systemID string) (*InsightsReportSetting, error) {
	systemIDsJSON, err := json.Marshal([]string{systemID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal systemIds: %w", err)
	}

	params := url.Values{}
	params.Add("customerName", c.Username)
	params.Add("systemIds", string(systemIDsJSON))

	path := fmt.Sprintf("/api/external/v1/insightsreportsetting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get insights report setting: HTTP %d - %s", statusCode, string(body))
	}

	var settings []InsightsReportSetting
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse insights report setting: %w", err)
	}

	if len(settings) == 0 {
		return nil, nil
	}

	return &settings[0], nil
}

// SetInsightsReportSetting updates insights report settings (daily or weekly) for a system
func (c *Client) SetInsightsReportSetting(systemID string, emailSet []string, enableInsightsReport bool, isDaily bool) error {
	models := []insightsReportSettingModel{{
		SystemID:             systemID,
		EmailSet:             emailSet,
		EnableInsightsReport: enableInsightsReport,
	}}

	modelsJSON, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal settingModels: %w", err)
	}

	isDailyStr := "false"
	if isDaily {
		isDailyStr = "true"
	}

	formData := url.Values{}
	formData.Set("customerName", c.Username)
	formData.Set("isDaily", isDailyStr)
	formData.Set("settingModels", string(modelsJSON))
	formData.Set("systemId", systemID)

	body, statusCode, err := c.DoFormRequest("POST", "/api/external/v1/insightsreportsetting", formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set insights report setting: HTTP %d - %s", statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set insights report setting: %s", msg)
		}
		return fmt.Errorf("failed to set insights report setting")
	}

	return nil
}

// GetInstanceDownSetting retrieves instance down notification settings for a specific project
func (c *Client) GetInstanceDownSetting(projectName string) (*InstanceDownSetting, error) {
	params := url.Values{}
	params.Add("operation", "display")
	params.Add("projectName", projectName)
	params.Add("customerName", c.Username)

	path := fmt.Sprintf("/api/external/v1/projects/update?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get instance down setting for project %s: HTTP %d - %s", projectName, statusCode, string(body))
	}

	var response struct {
		Data struct {
			ProjectModelAllJSON struct {
				ProjectName              string   `json:"projectName"`
				InstanceDownThreshold    int64    `json:"instanceDownThreshold"`
				InstanceDownEnable       bool     `json:"instanceDownEnable"`
				InstanceDownEmails       []string `json:"instanceDownEmails"`
				InstanceDownDampening    int64    `json:"instanceDownDampening"`
				InstanceDownReportNumber int64    `json:"instanceDownReportNumber"`
			} `json:"projectModelAllJSON"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse instance down setting: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("failed to get instance down setting for project %s", projectName)
	}

	proj := response.Data.ProjectModelAllJSON
	emails := proj.InstanceDownEmails
	if emails == nil {
		emails = []string{}
	}

	return &InstanceDownSetting{
		ProjectName:              proj.ProjectName,
		InstanceDownEnable:       proj.InstanceDownEnable,
		InstanceDownDampening:    proj.InstanceDownDampening,
		InstanceDownThreshold:    proj.InstanceDownThreshold,
		InstanceDownReportNumber: proj.InstanceDownReportNumber,
		InstanceDownEmails:       emails,
	}, nil
}

// SetInstanceDownSetting updates instance down notification settings for a specific project
func (c *Client) SetInstanceDownSetting(setting *InstanceDownSetting) error {
	emails := setting.InstanceDownEmails
	if emails == nil {
		emails = []string{}
	}

	emailsJSON, err := json.Marshal(emails)
	if err != nil {
		return fmt.Errorf("failed to marshal instanceDownEmails: %w", err)
	}

	instanceDownEnableStr := "false"
	if setting.InstanceDownEnable {
		instanceDownEnableStr = "true"
	}

	formData := url.Values{}
	formData.Set("operation", "updateprojsettings")
	formData.Set("projectName", setting.ProjectName)
	formData.Set("instanceDownEnable", instanceDownEnableStr)
	formData.Set("instanceDownDampening", fmt.Sprintf("%d", setting.InstanceDownDampening))
	formData.Set("instanceDownThreshold", fmt.Sprintf("%d", setting.InstanceDownThreshold))
	formData.Set("instanceDownReportNumber", fmt.Sprintf("%d", setting.InstanceDownReportNumber))
	formData.Set("instanceDownEmails", string(emailsJSON))

	body, statusCode, err := c.DoFormRequest("POST", "/api/external/v1/projects/update", formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set instance down setting for project %s: HTTP %d - %s", setting.ProjectName, statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	// Response has shape: {"alertSetting": {"success": true}, "instanceDownThreshold": {"success": true}}
	// Check that no sub-key has success=false
	for key, val := range response {
		if subMap, ok := val.(map[string]any); ok {
			if success, ok := subMap["success"].(bool); ok && !success {
				return fmt.Errorf("failed to set instance down setting for project %s (key: %s)", setting.ProjectName, key)
			}
		}
	}

	return nil
}

// MiscellaneousSettings holds the miscellaneous system framework settings
type MiscellaneousSettings struct {
	LongTerm                             bool  `json:"longTerm"`
	ShouldAutoShare                      bool  `json:"shouldAutoShare"`
	RootCauseReverseEntryFilterThreshold int64 `json:"rootCauseReverseEntryFilterThreshold"`
	EnableCompositeTimeline              bool  `json:"enableCompositeTimeline"`
}

// SystemFrameworkSetting holds the fields updated via the systemFrameworkSetting operation
type SystemFrameworkSetting struct {
	ShouldAutoShare                      bool  `json:"shouldAutoShare"`
	RootCauseReverseEntryFilterThreshold int64 `json:"rootCauseReverseEntryFilterThreshold"`
	EnableCompositeTimeline              bool  `json:"enableCompositeTimeline"`
}

// systemFrameworkRawEntry is the internal parse format for each element of ownSystemArr
type systemFrameworkRawEntry struct {
	SystemKey struct {
		UserName        string `json:"userName"`
		SystemName      string `json:"systemName"`
		EnvironmentName string `json:"environmentName"`
	} `json:"systemKey"`
	Order         int64  `json:"order"`
	HideFlag      bool   `json:"hideFlag"`
	LongTerm      bool   `json:"longTerm"`
	SystemSetting string `json:"systemSetting"`
}

// systemSettingInner is the inner JSON string stored in the systemSetting field of systemFrameworkRawEntry
type systemSettingInner struct {
	ShouldAutoShare                      bool  `json:"shouldAutoShare"`
	RootCauseReverseEntryFilterThreshold int64 `json:"rootCauseReverseEntryFilterThreshold"`
	EnableCompositeTimeline              bool  `json:"enableCompositeTimeline"`
}

// getSystemFrameworkEntry fetches the raw system framework entry for a specific system
func (c *Client) getSystemFrameworkEntry(systemID string) (*systemFrameworkRawEntry, error) {
	params := url.Values{}
	params.Add("customerName", c.Username)
	params.Add("needDetail", "false")
	params.Add("tzOffset", "-18000000")

	path := fmt.Sprintf("/api/external/v1/systemframework?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get system framework: HTTP %d - %s", statusCode, string(body))
	}

	var response struct {
		OwnSystemArr []string `json:"ownSystemArr"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse system framework response: %w", err)
	}

	for _, entryStr := range response.OwnSystemArr {
		var entry systemFrameworkRawEntry
		if err := json.Unmarshal([]byte(entryStr), &entry); err != nil {
			continue
		}
		if entry.SystemKey.SystemName == systemID {
			return &entry, nil
		}
	}

	return nil, nil
}

// GetMiscellaneousSettings retrieves the miscellaneous system framework settings for a specific system
func (c *Client) GetMiscellaneousSettings(systemID string) (*MiscellaneousSettings, error) {
	entry, err := c.getSystemFrameworkEntry(systemID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	result := &MiscellaneousSettings{
		LongTerm: entry.LongTerm,
	}

	if entry.SystemSetting != "" {
		var inner systemSettingInner
		if err := json.Unmarshal([]byte(entry.SystemSetting), &inner); err == nil {
			result.ShouldAutoShare = inner.ShouldAutoShare
			result.RootCauseReverseEntryFilterThreshold = inner.RootCauseReverseEntryFilterThreshold
			result.EnableCompositeTimeline = inner.EnableCompositeTimeline
		}
	}

	return result, nil
}

// SetLongTermSetting updates the longTerm flag for a system, preserving the current order value.
func (c *Client) SetLongTermSetting(systemID string, longTerm bool) error {
	currentOrder := int64(0)
	entry, err := c.getSystemFrameworkEntry(systemID)
	if err == nil && entry != nil {
		currentOrder = entry.Order
	}

	type configEntry struct {
		SystemName      string `json:"systemName"`
		EnvironmentName string `json:"environmentName"`
		Order           int64  `json:"order"`
		LongTerm        bool   `json:"longTerm"`
	}

	config := []configEntry{{
		SystemName:      systemID,
		EnvironmentName: "All",
		Order:           currentOrder,
		LongTerm:        longTerm,
	}}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	params := url.Values{}
	params.Set("operation", "hideOrOrderOrLongTerm")
	params.Set("customerName", c.Username)
	params.Set("config", string(configJSON))
	params.Set("systemName", systemID)

	path := fmt.Sprintf("/api/external/v1/systemframework?%s", params.Encode())
	body, statusCode, err := c.DoRequest("POST", path, nil)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set longTerm setting: HTTP %d - %s", statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set longTerm setting: %s", msg)
		}
		return fmt.Errorf("failed to set longTerm setting")
	}

	return nil
}

// SetSystemFrameworkSetting updates shouldAutoShare, rootCauseReverseEntryFilterThreshold, and enableCompositeTimeline for a system
func (c *Client) SetSystemFrameworkSetting(systemID string, settings *SystemFrameworkSetting) error {
	type systemKeyType struct {
		UserName        string `json:"userName"`
		SystemName      string `json:"systemName"`
		EnvironmentName string `json:"environmentName"`
	}

	systemKey := systemKeyType{
		UserName:        c.Username,
		SystemName:      systemID,
		EnvironmentName: "All",
	}

	systemKeyJSON, err := json.Marshal(systemKey)
	if err != nil {
		return fmt.Errorf("failed to marshal systemKey: %w", err)
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal systemFrameworkSetting: %w", err)
	}

	params := url.Values{}
	params.Set("operation", "systemFrameworkSetting")
	params.Set("systemKey", string(systemKeyJSON))
	params.Set("systemFrameworkSetting", string(settingsJSON))
	params.Set("systemName", systemID)
	params.Set("customerName", c.Username)

	path := fmt.Sprintf("/api/external/v1/systemframework?%s", params.Encode())
	body, statusCode, err := c.DoRequest("POST", path, nil)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to set system framework setting: HTTP %d - %s", statusCode, string(body))
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if success, ok := response["success"].(bool); ok && !success {
		if msg, ok := response["message"].(string); ok {
			return fmt.Errorf("failed to set system framework setting: %s", msg)
		}
		return fmt.Errorf("failed to set system framework setting")
	}

	return nil
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
	current.TicketOpenTime = updates.TicketOpenTime
	current.ComponentLevelIncidentConsolidation = updates.ComponentLevelIncidentConsolidation
	current.EnabledConsolidationAlgorithms = updates.EnabledConsolidationAlgorithms
	current.MaxNotificationDelayTolerance = updates.MaxNotificationDelayTolerance
	current.ProjectLevelDampeningWindows = updates.ProjectLevelDampeningWindows
	current.CustomConsolidationRules = updates.CustomConsolidationRules
	current.MetricLogConsolidationConfigs = updates.MetricLogConsolidationConfigs
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
