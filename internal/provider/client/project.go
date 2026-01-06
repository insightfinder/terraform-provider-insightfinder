// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ProjectConfig represents a project configuration
type ProjectConfig struct {
	ProjectName         string                 `json:"projectName"`
	ProjectDisplayName  string                 `json:"projectDisplayName,omitempty"`
	SystemName          string                 `json:"systemName"`
	DataType            string                 `json:"dataType"`
	InstanceType        string                 `json:"instanceType"`
	ProjectCloudType    string                 `json:"projectCloudType"`
	InsightAgentType    string                 `json:"insightAgentType,omitempty"`
	ProjectCreationType string                 `json:"projectCreationType,omitempty"`
	CValue              int                    `json:"cValue,omitempty"`
	PValue              float64                `json:"pValue,omitempty"`
	Settings            map[string]interface{} `json:"settings,omitempty"`
}

// Settings fields struct
type ProjectSettings struct {
	DailyModelSpan            int    `json:"dailyModelSpan,omitempty"`
	KeywordFeatureNumber      int    `json:"keywordFeatureNumber,omitempty"`
	MaxLogModelSize           int    `json:"maxLogModelSize,omitempty"`
	ModelKeywordSetting       int    `json:"modelKeywordSetting,omitempty"`
	NlpFlag                   bool   `json:"nlpFlag,omitempty"`
	ProjectModelFlag          bool   `json:"projectModelFlag,omitempty"`
	MaximumThreads            int    `json:"maximumThreads,omitempty"`
	LogDetectionMinCount      int    `json:"logDetectionMinCount,omitempty"`
	LogDetectionSize          int    `json:"logDetectionSize,omitempty"`
	MaximumDetectionWaitTime  int    `json:"maximumDetectionWaitTime,omitempty"`
	KeywordSetting            int    `json:"keywordSetting,omitempty"`
	LogPatternLimitLevel      int    `json:"logPatternLimitLevel,omitempty"`
	NormalEventCausalFlag     bool   `json:"normalEventCausalFlag,omitempty"`
	SimilaritySensitivity     string `json:"similaritySensitivity,omitempty"`
	CollectAllRareEventsFlag  bool   `json:"collectAllRareEventsFlag,omitempty"`
	RareEventAlertThresholds  int    `json:"rareEventAlertThresholds"`
	LogAnomalyEventBaseScore  string `json:"logAnomalyEventBaseScore,omitempty"`
	RareNumberLimit           int    `json:"rareNumberLimit,omitempty"`
	WhitelistNumberLimit      int    `json:"whitelistNumberLimit,omitempty"`
	NewPatternNumberLimit     int    `json:"newPatternNumberLimit,omitempty"`
	HotNumberLimit            int    `json:"hotNumberLimit,omitempty"`
	ColdNumberLimit           int    `json:"coldNumberLimit,omitempty"`
	RareAnomalyType           int    `json:"rareAnomalyType,omitempty"`
	HotEventThreshold         int    `json:"hotEventThreshold,omitempty"`
	ColdEventThreshold        int    `json:"coldEventThreshold,omitempty"`
	DisableLogCompressEvent   bool   `json:"disableLogCompressEvent,omitempty"`
	EnableHotEvent            bool   `json:"enableHotEvent,omitempty"`
	HotEventCalmDownPeriod    int    `json:"hotEventCalmDownPeriod,omitempty"`
	InstanceDownEnable        bool   `json:"instanceDownEnable,omitempty"`
	AnomalySamplingInterval   int    `json:"anomalySamplingInterval,omitempty"`
	HotEventDetectionMode     int    `json:"hotEventDetectionMode,omitempty"`
	AnomalyDetectionMode      int    `json:"anomalyDetectionMode,omitempty"`
	PrettyJSONConvertorFlag   bool   `json:"prettyJsonConvertorFlag,omitempty"`
	ZoneNameKey               string `json:"zoneNameKey"`
	MultiLineFlag             bool   `json:"multiLineFlag,omitempty"`
	FeatureOutlierSensitivity string `json:"featureOutlierSensitivity,omitempty"`
	BaseValueSetting          struct {
		IsSourceProject       bool          `json:"isSourceProject,omitempty"`
		MappingKeys           []interface{} `json:"mappingKeys,omitempty"`
		BaseValueKeys         []interface{} `json:"baseValueKeys,omitempty"`
		MetricProjects        []interface{} `json:"metricProjects,omitempty"`
		AdditionalMetricNames []interface{} `json:"additionalMetricNames,omitempty"`
	} `json:"baseValueSetting,omitempty"`
	CdfSetting                         []interface{} `json:"cdfSetting,omitempty"`
	DisableModelKeywordStatsCollection bool          `json:"disableModelKeywordStatsCollection,omitempty"`
	InstanceConvertFlag                bool          `json:"instanceConvertFlag,omitempty"`
	NewAlertFlag                       bool          `json:"newAlertFlag,omitempty"`
	IsGroupingByInstance               bool          `json:"isGroupingByInstance,omitempty"`
	LogLabelSettingCreate              struct {
	} `json:"logLabelSettingCreate,omitempty"`
	FeatureOutlierThreshold float64 `json:"featureOutlierThreshold,omitempty"`
	LogLabelSettingRemove   struct {
	} `json:"logLabelSettingRemove,omitempty"`
	LogToMetricCreate struct {
	} `json:"logToMetricCreate,omitempty"`
	LogToMetricDelete struct {
	} `json:"logToMetricDelete,omitempty"`
	LogJSONTypeUpdate struct {
	} `json:"logJsonTypeUpdate,omitempty"`
	IsTracePrompt        bool          `json:"isTracePrompt,omitempty"`
	LogToLogSettingList  []interface{} `json:"logToLogSettingList,omitempty"`
	LlmEvaluationSetting struct {
		IsHallucinationEvaluation     bool `json:"isHallucinationEvaluation,omitempty"`
		IsAnswerRelevantEvaluation    bool `json:"isAnswerRelevantEvaluation,omitempty"`
		IsLogicConsistencyEvaluation  bool `json:"isLogicConsistencyEvaluation,omitempty"`
		IsFactualInaccuracyEvaluation bool `json:"isFactualInaccuracyEvaluation,omitempty"`
		IsMaliciousPromptEvaluation   bool `json:"isMaliciousPromptEvaluation,omitempty"`
		IsToxicityEvaluation          bool `json:"isToxicityEvaluation,omitempty"`
		IsPiiPhiLeakageEvaluation     bool `json:"isPiiPhiLeakageEvaluation,omitempty"`
		IsTopicGuardrailsEvaluation   bool `json:"isTopicGuardrailsEvaluation,omitempty"`
		IsToneDetectionEvaluation     bool `json:"isToneDetectionEvaluation,omitempty"`
		IsAnomalousOutliersEvaluation bool `json:"isAnomalousOutliersEvaluation,omitempty"`
		ShowSafetyTemplate            bool `json:"showSafetyTemplate,omitempty"`
		IsGenderBiasEvaluation        bool `json:"isGenderBiasEvaluation,omitempty"`
		IsRacialBiasEvaluation        bool `json:"isRacialBiasEvaluation,omitempty"`
		IsSocioeconomicBiasEvaluation bool `json:"isSocioeconomicBiasEvaluation,omitempty"`
		IsCulturalBiasEvaluation      bool `json:"isCulturalBiasEvaluation,omitempty"`
		IsReligiousBiasEvaluation     bool `json:"isReligiousBiasEvaluation,omitempty"`
		IsPoliticalBiasEvaluation     bool `json:"isPoliticalBiasEvaluation,omitempty"`
		IsDisabilityBiasEvaluation    bool `json:"isDisabilityBiasEvaluation,omitempty"`
		IsAgeBiasEvaluation           bool `json:"isAgeBiasEvaluation,omitempty"`
	} `json:"llmEvaluationSetting,,omitempty"`
	IsEdgeBrain                          bool          `json:"isEdgeBrain,omitempty"`
	ProjectName                          string        `json:"projectName,omitempty"`
	CValue                               int           `json:"cValue,omitempty"`
	PValue                               float64       `json:"pValue,omitempty"`
	IncidentPredictionWindow             int           `json:"incidentPredictionWindow,omitempty"`
	MinIncidentPredictionWindow          int           `json:"minIncidentPredictionWindow,omitempty"`
	IncidentRelationSearchWindow         int           `json:"incidentRelationSearchWindow,omitempty"`
	IncidentPredictionEventLimit         int           `json:"incidentPredictionEventLimit,omitempty"`
	RootCauseCountThreshold              int           `json:"rootCauseCountThreshold,omitempty"`
	RootCauseProbabilityThreshold        float64       `json:"rootCauseProbabilityThreshold,omitempty"`
	CompositeRCALimit                    int           `json:"compositeRCALimit,omitempty"`
	RootCauseLogMessageSearchRange       int           `json:"rootCauseLogMessageSearchRange,omitempty"`
	CausalPredictionSetting              int           `json:"causalPredictionSetting,omitempty"`
	CausalMinDelay                       string        `json:"causalMinDelay,omitempty"`
	RootCauseRankSetting                 int           `json:"rootCauseRankSetting,omitempty"`
	MaximumRootCauseResultSize           int           `json:"maximumRootCauseResultSize,omitempty"`
	MultiHopSearchLevel                  int           `json:"multiHopSearchLevel,omitempty"`
	AvgPerIncidentDowntimeCost           float64       `json:"avgPerIncidentDowntimeCost,omitempty"`
	PredictionRuleActiveCondition        int           `json:"predictionRuleActiveCondition,omitempty"`
	PredictionRuleFalsePositiveThreshold int           `json:"predictionRuleFalsePositiveThreshold,omitempty"`
	PredictionRuleActiveThreshold        float64       `json:"predictionRuleActiveThreshold,omitempty"`
	PredictionRuleInactiveThreshold      float64       `json:"predictionRuleInactiveThreshold,omitempty"`
	PredictionProbabilityThreshold       float64       `json:"predictionProbabilityThreshold,omitempty"`
	AlertHourlyCost                      float64       `json:"alertHourlyCost,omitempty"`
	AlertAverageTime                     int           `json:"alertAverageTime,omitempty"`
	IgnoreInstanceForKB                  bool          `json:"ignoreInstanceForKB,omitempty"`
	ShowInstanceDown                     bool          `json:"showInstanceDown,omitempty"`
	RetentionTime                        int           `json:"retentionTime,omitempty"`
	UBLRetentionTime                     int           `json:"UBLRetentionTime,omitempty"`
	TrainingFilter                       bool          `json:"trainingFilter,omitempty"`
	SharedUsernames                      []interface{} `json:"sharedUsernames,omitempty"`
	MultiHopSearchLimit                  string        `json:"multiHopSearchLimit,omitempty"`
	ProjectDisplayName                   string        `json:"projectDisplayName,omitempty"`
	EnableNewAlertEmail                  bool          `json:"enableNewAlertEmail,omitempty"`
	ProjectTimeZone                      string        `json:"projectTimeZone,omitempty"`
	PredictionCountThreshold             int           `json:"predictionCountThreshold,omitempty"`
	SamplingInterval                     int           `json:"samplingInterval,omitempty"`
	MinValidModelSpan                    int           `json:"minValidModelSpan,omitempty"`
	MaxWebHookRequestSize                int           `json:"maxWebHookRequestSize,omitempty"`
	WebhookURL                           string        `json:"webhookUrl,omitempty"`
	WebhookHeaderList                    []interface{} `json:"webhookHeaderList,omitempty"`
	WebhookTypeSetStr                    string        `json:"webhookTypeSetStr,omitempty"`
	WebhookBlackListSetStr               string        `json:"webhookBlackListSetStr,omitempty"`
	WebhookCriticalKeywordSetStr         string        `json:"webhookCriticalKeywordSetStr,omitempty"`
	WebhookAlertDampening                int           `json:"webhookAlertDampening,omitempty"`
	Proxy                                string        `json:"proxy,omitempty"`
	NewPatternRange                      int           `json:"newPatternRange,omitempty"`
	LargeProject                         bool          `json:"largeProject,omitempty"`
	EnableAnomalyScoreEscalation         bool          `json:"enableAnomalyScoreEscalation,omitempty"`
	EscalationAnomalyScoreThreshold      string        `json:"escalationAnomalyScoreThreshold,omitempty"`
	IgnoreAnomalyScoreThreshold          string        `json:"ignoreAnomalyScoreThreshold,omitempty"`
	EnableStreamDetection                bool          `json:"enableStreamDetection,omitempty"`
	InstanceGroupingUpdate               struct {
		AutoFill bool `json:"autoFill,omitempty"`
	} `json:"instanceGroupingUpdate,omitempty"`
	EmailSetting struct {
		OnlySendWithRCA                    bool   `json:"onlySendWithRCA,omitempty"`
		EnableNotificationAW               bool   `json:"enableNotificationAW,omitempty"`
		EnableIncidentPredictionEmailAlert bool   `json:"enableIncidentPredictionEmailAlert,omitempty"`
		EnableIncidentDetectionEmailAlert  bool   `json:"enableIncidentDetectionEmailAlert,omitempty"`
		EnableAlertsEmail                  bool   `json:"enableAlertsEmail,omitempty"`
		EnableRootCauseEmailAlert          bool   `json:"enableRootCauseEmailAlert,omitempty"`
		EmailDampeningPeriod               int    `json:"emailDampeningPeriod"`
		AlertsEmailDampeningPeriod         int    `json:"alertsEmailDampeningPeriod"`
		PredictionEmailDampeningPeriod     int    `json:"predictionEmailDampeningPeriod"`
		AwSeverityLevel                    string `json:"awSeverityLevel,omitempty"`
	} `json:"emailSetting,omitempty"`
}

// ProjectResponse represents the API response for project operations
type ProjectResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetProject retrieves a project's configuration
func (c *Client) GetProject(projectName, username string) (*ProjectConfig, error) {
	params := url.Values{}
	params.Add("projectList", fmt.Sprintf(`[{"customerName":"%s","projectName":"%s"}]`, username, projectName))

	path := fmt.Sprintf("/api/external/v1/watch-tower-setting?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// 404 = endpoint not found, 204 = no content (project deleted)
	if statusCode == 404 || statusCode == 204 {
		return nil, nil // Project doesn't exist
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get project: HTTP %d", statusCode)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	settingList, ok := response["settingList"].(map[string]interface{})
	if !ok || settingList[projectName] == nil {
		return nil, nil // Project doesn't exist
	}

	// Parse the project settings
	var project ProjectConfig
	project.ProjectName = projectName

	// The actual settings are in a JSON string within the response
	if settingsStr, ok := settingList[projectName].(string); ok {
		var settingsWrapper map[string]interface{}
		if err := json.Unmarshal([]byte(settingsStr), &settingsWrapper); err != nil {
			return nil, fmt.Errorf("failed to parse project settings: %w", err)
		}

		// The actual settings are in the DATA field
		var settings map[string]interface{}
		if dataField, ok := settingsWrapper["DATA"].(map[string]interface{}); ok {
			settings = dataField
		} else {
			// Fallback to using the whole wrapper if DATA field doesn't exist
			settings = settingsWrapper
		}
		project.Settings = settings

		// Extract common fields from settings
		if displayName, ok := settings["projectDisplayName"].(string); ok {
			project.ProjectDisplayName = displayName
		}
		if cValue, ok := settings["cValue"].(float64); ok {
			project.CValue = int(cValue)
		}
		if pValue, ok := settings["pValue"].(float64); ok {
			project.PValue = pValue
		}
	}

	return &project, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(project *ProjectConfig) error {
	formData := url.Values{}
	formData.Set("operation", "create")
	formData.Set("projectName", project.ProjectName)
	formData.Set("systemName", project.SystemName)
	formData.Set("instanceType", project.InstanceType)
	formData.Set("dataType", project.DataType)
	formData.Set("projectCloudType", project.ProjectCloudType)

	if project.ProjectDisplayName != "" {
		formData.Set("projectDisplayName", project.ProjectDisplayName)
	}
	if project.InsightAgentType != "" {
		formData.Set("insightAgentType", project.InsightAgentType)
	}
	if project.ProjectCreationType != "" {
		formData.Set("projectCreationType", project.ProjectCreationType)
	}

	body, statusCode, err := c.DoFormRequest("POST", "/api/v1/check-and-add-custom-project", formData)
	if err != nil {
		return err
	}

	// Handle the response
	if statusCode == 400 {
		// Check if it's a "project already exists" error
		var response ProjectResponse
		if err := json.Unmarshal(body, &response); err == nil {
			if !response.Success && (strings.Contains(response.Message, "already existed") || strings.Contains(response.Message, "already exists")) {
				// Project already exists, treat as success
				return nil
			}
		}
		return fmt.Errorf("failed to create project: HTTP %d - %s", statusCode, string(body))
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to create project: HTTP %d - %s", statusCode, string(body))
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("failed to create project: %s", response.Message)
	}

	return nil
}

// UpdateProject updates an existing project's configuration
func (c *Client) UpdateProject(project *ProjectConfig) error {
	// Build the settings JSON
	settings := project.Settings
	if settings == nil {
		settings = make(map[string]interface{})
	}

	// Convert settings map to ProjectSettings struct to ensure type safety
	// and only send fields that are present in the map
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	var projectSettings ProjectSettings
	if err := json.Unmarshal(settingsJSON, &projectSettings); err != nil {
		return fmt.Errorf("failed to unmarshal to ProjectSettings: %w", err)
	}

	// Marshal back to map[string]interface{} with omitempty fields excluded
	finalSettingsJSON, err := json.Marshal(projectSettings)
	if err != nil {
		return fmt.Errorf("failed to marshal ProjectSettings: %w", err)
	}

	var finalSettings map[string]interface{}
	if err := json.Unmarshal(finalSettingsJSON, &finalSettings); err != nil {
		return fmt.Errorf("failed to unmarshal final settings: %w", err)
	}

	path := fmt.Sprintf("/api/external/v1/watch-tower-setting?projectName=%s&customerName=%s",
		url.QueryEscape(project.ProjectName), url.QueryEscape(c.Username))

	body, statusCode, err := c.DoRequest("POST", path, finalSettings)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to update project: HTTP %d - %s", statusCode, string(body))
	}

	// The update endpoint might return empty body or simple success message
	if len(body) == 0 {
		return nil
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// If response is not JSON, just log and continue
		return nil
	}

	if !response.Success {
		return fmt.Errorf("failed to update project: %s", response.Message)
	}

	return nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectName string) error {
	formData := url.Values{}
	formData.Set("projectName", projectName)

	body, statusCode, err := c.DoFormRequest("POST", "/api/v1/delete-project", formData)
	if err != nil {
		return err
	}

	// Treat 404 and 405 as success (project doesn't exist or endpoint not implemented)
	if statusCode == 404 || statusCode == 405 {
		return nil // Already deleted or deletion not supported
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to delete project: HTTP %d - %s", statusCode, string(body))
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("failed to delete project: %s", response.Message)
	}

	return nil
}
