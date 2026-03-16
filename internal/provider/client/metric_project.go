// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// MetricProjectSettings represents the settings specific to a metric project
type MetricProjectSettings struct {
	// Identification
	ProjectName        string `json:"projectName,omitempty"`
	ProjectDisplayName string `json:"projectDisplayName,omitempty"`

	// Detection tuning
	CValue                              int     `json:"cValue,omitempty"`
	PValue                              float64 `json:"pValue,omitempty"`
	HighRatioCValue                     int     `json:"highRatioCValue,omitempty"`
	MaximumHint                         int     `json:"maximumHint,omitempty"`
	DynamicBaselineDetectionFlag        bool    `json:"dynamicBaselineDetectionFlag,omitempty"`
	PositiveBaselineViolationFactor     float64 `json:"positiveBaselineViolationFactor,omitempty"`
	NegativeBaselineViolationFactor     float64 `json:"negativeBaselineViolationFactor,omitempty"`
	EnablePeriodAnomalyFilter           bool    `json:"enablePeriodAnomalyFilter,omitempty"`
	EnableUBLDetect                     bool    `json:"enableUBLDetect,omitempty"`
	EnableCumulativeDetect              bool    `json:"enableCumulativeDetect,omitempty"`
	ModelSpan                           int     `json:"modelSpan,omitempty"`
	EnableMetricDataPrediction          bool    `json:"enableMetricDataPrediction,omitempty"`
	EnableBaselineDetectionDoubleVerify bool    `json:"enableBaselineDetectionDoubleVerify,omitempty"`
	EnableFillGap                       bool    `json:"enableFillGap,omitempty"`
	EnableStoreFilledGap                bool    `json:"enableStoreFilledGap,omitempty"`
	GapFillingTrainingDataLength        int     `json:"gapFillingTrainingDataLength,omitempty"`
	PatternIdGenerationRule             int     `json:"patternIdGenerationRule,omitempty"`
	AnomalyGapToleranceCount            int     `json:"anomalyGapToleranceCount,omitempty"`
	FilterByAnomalyInBaselineGeneration bool    `json:"filterByAnomalyInBaselineGeneration,omitempty"`
	BaselineDuration                    int     `json:"baselineDuration,omitempty"`
	AnomalyDampening                    int     `json:"anomalyDampening,omitempty"`
	InstanceDownRatioThreshold          float64 `json:"instanceDownRatioThreshold,omitempty"`
	ComponentNameAutoOverwrite          bool    `json:"componentNameAutoOverwrite,omitempty"`

	// Prediction
	PredictionTrainingDataLength     int     `json:"predictionTrainingDataLength,omitempty"`
	PredictionCorrelationSensitivity float64 `json:"predictionCorrelationSensitivity,omitempty"`
	EnableKPIPrediction              bool    `json:"enableKPIPrediction,omitempty"`

	// Instance down
	InstanceDownThreshold    int  `json:"instanceDownThreshold,omitempty"`
	InstanceDownReportNumber int  `json:"instanceDownReportNumber,omitempty"`
	InstanceDownEnable       bool `json:"instanceDownEnable,omitempty"`

	// Common project settings
	ProjectTimeZone                 string  `json:"projectTimeZone,omitempty"`
	SamplingInterval                int     `json:"samplingInterval,omitempty"`
	RetentionTime                   int     `json:"retentionTime,omitempty"`
	UBLRetentionTime                int     `json:"UBLRetentionTime,omitempty"`
	TrainingFilter                  bool    `json:"trainingFilter,omitempty"`
	EnableNewAlertEmail             bool    `json:"enableNewAlertEmail,omitempty"`
	EnableAnomalyScoreEscalation    bool    `json:"enableAnomalyScoreEscalation,omitempty"`
	EscalationAnomalyScoreThreshold string  `json:"escalationAnomalyScoreThreshold,omitempty"`
	IgnoreAnomalyScoreThreshold     string  `json:"ignoreAnomalyScoreThreshold,omitempty"`
	EnableStreamDetection           bool    `json:"enableStreamDetection,omitempty"`
	LargeProject                    bool    `json:"largeProject,omitempty"`
	NewPatternRange                 int     `json:"newPatternRange,omitempty"`
	Proxy                           string  `json:"proxy,omitempty"`
	IgnoreInstanceForKB             bool    `json:"ignoreInstanceForKB,omitempty"`
	ShowInstanceDown                bool    `json:"showInstanceDown,omitempty"`
	AlertHourlyCost                 float64 `json:"alertHourlyCost,omitempty"`
	AlertAverageTime                int     `json:"alertAverageTime,omitempty"`
	AvgPerIncidentDowntimeCost      float64 `json:"avgPerIncidentDowntimeCost,omitempty"`

	// Incident prediction and RCA
	IncidentPredictionWindow             int     `json:"incidentPredictionWindow,omitempty"`
	MinIncidentPredictionWindow          int     `json:"minIncidentPredictionWindow,omitempty"`
	IncidentRelationSearchWindow         int     `json:"incidentRelationSearchWindow,omitempty"`
	IncidentPredictionEventLimit         int     `json:"incidentPredictionEventLimit,omitempty"`
	RootCauseCountThreshold              int     `json:"rootCauseCountThreshold,omitempty"`
	RootCauseProbabilityThreshold        float64 `json:"rootCauseProbabilityThreshold,omitempty"`
	CompositeRCALimit                    int     `json:"compositeRCALimit,omitempty"`
	RootCauseLogMessageSearchRange       int     `json:"rootCauseLogMessageSearchRange,omitempty"`
	CausalPredictionSetting              int     `json:"causalPredictionSetting,omitempty"`
	RootCauseRankSetting                 int     `json:"rootCauseRankSetting,omitempty"`
	MaximumRootCauseResultSize           int     `json:"maximumRootCauseResultSize,omitempty"`
	MultiHopSearchLevel                  int     `json:"multiHopSearchLevel,omitempty"`
	MultiHopSearchLimit                  string  `json:"multiHopSearchLimit,omitempty"`
	PredictionCountThreshold             int     `json:"predictionCountThreshold,omitempty"`
	PredictionProbabilityThreshold       float64 `json:"predictionProbabilityThreshold,omitempty"`
	PredictionRuleActiveCondition        int     `json:"predictionRuleActiveCondition,omitempty"`
	PredictionRuleFalsePositiveThreshold int     `json:"predictionRuleFalsePositiveThreshold,omitempty"`
	PredictionRuleActiveThreshold        float64 `json:"predictionRuleActiveThreshold,omitempty"`
	PredictionRuleInactiveThreshold      float64 `json:"predictionRuleInactiveThreshold,omitempty"`
	MinValidModelSpan                    int     `json:"minValidModelSpan,omitempty"`

	// Webhook
	MaxWebHookRequestSize        int    `json:"maxWebHookRequestSize,omitempty"`
	WebhookURL                   string `json:"webhookUrl,omitempty"`
	WebhookTypeSetStr            string `json:"webhookTypeSetStr,omitempty"`
	WebhookBlackListSetStr       string `json:"webhookBlackListSetStr,omitempty"`
	WebhookCriticalKeywordSetStr string `json:"webhookCriticalKeywordSetStr,omitempty"`
	WebhookAlertDampening        int    `json:"webhookAlertDampening,omitempty"`

	// Complex / array fields
	LinkedLogProjects                      []interface{} `json:"linkedLogProjects,omitempty"`
	ComponentMetricSettingOverallModelList []interface{} `json:"componentMetricSettingOverallModelList,omitempty"`
	SharedUsernames                        []interface{} `json:"sharedUsernames,omitempty"`
	WebhookHeaderList                      []interface{} `json:"webhookHeaderList,omitempty"`

	InstanceGroupingUpdate struct {
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

// UpdateMetricProject updates an existing metric project's configuration.
func (c *Client) UpdateMetricProject(project *ProjectConfig) error {
	settings := project.Settings
	if settings == nil {
		settings = make(map[string]interface{})
	}

	path := fmt.Sprintf("/api/external/v1/watch-tower-setting?projectName=%s&customerName=%s",
		url.QueryEscape(project.ProjectName), url.QueryEscape(c.Username))

	body, statusCode, err := c.DoRequest("POST", path, settings)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to update metric project: HTTP %d - %s", statusCode, string(body))
	}

	if len(body) == 0 {
		return nil
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	if !response.Success {
		return fmt.Errorf("failed to update metric project: %s", response.Message)
	}

	return nil
}
