// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
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
	EnableComponentLevelDetection       bool    `json:"enableComponentLevelDetection,omitempty"`
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
	IncidentPriorityByAnomalyScoreSetting struct {
		Enabled               bool              `json:"enabled,omitempty"`
		PriorityScoreRangeMap map[string]string `json:"priorityScoreRangeMap,omitempty"`
	} `json:"incidentPriorityByAnomalyScoreSetting,omitempty"`
}

// ─── Metric Configuration Types ───────────────────────────────────────────────

// MetricAlertSetting is one row in the globalSetting or componentLevelSettingList
// returned by GET /api/external/v1/componentmetricupdate.
// RougeValue is *string because the API returns it as a JSON string or null.
type MetricAlertSetting struct {
	SMetric                            string  `json:"smetric"`
	ComponentName                      string  `json:"componentName"`
	ThresholdAlertLowerBound           string  `json:"thresholdAlertLowerBound"`
	ThresholdAlertUpperBound           string  `json:"thresholdAlertUpperBound"`
	ThresholdAlertLowerBoundNegative   string  `json:"thresholdAlertLowerBoundNegative"`
	ThresholdAlertUpperBoundNegative   string  `json:"thresholdAlertUpperBoundNegative"`
	ThresholdNoAlertLowerBound         string  `json:"thresholdNoAlertLowerBound"`
	ThresholdNoAlertUpperBound         string  `json:"thresholdNoAlertUpperBound"`
	ThresholdNoAlertLowerBoundNegative string  `json:"thresholdNoAlertLowerBoundNegative"`
	ThresholdNoAlertUpperBoundNegative string  `json:"thresholdNoAlertUpperBoundNegative"`
	IncidentAlertLowerBound            string  `json:"incidentAlertLowerBound"`
	IncidentAlertUpperBound            string  `json:"incidentAlertUpperBound"`
	IncidentAlertLowerBoundNegative    string  `json:"incidentAlertLowerBoundNegative"`
	IncidentAlertUpperBoundNegative    string  `json:"incidentAlertUpperBoundNegative"`
	IncidentNoAlertLowerBound          string  `json:"incidentNoAlertLowerBound"`
	IncidentNoAlertUpperBound          string  `json:"incidentNoAlertUpperBound"`
	IncidentNoAlertLowerBoundNegative  string  `json:"incidentNoAlertLowerBoundNegative"`
	IncidentNoAlertUpperBoundNegative  string  `json:"incidentNoAlertUpperBoundNegative"`
	IsKPI                              bool    `json:"isKPI"`
	IsFlappingResultOnly               bool    `json:"isFlappingResultOnly"`
	IncidentDurationThreshold          int64   `json:"incidentDurationThreshold"`
	DetectionType                      string  `json:"detectionType"`
	CValueOverride                     *int64  `json:"cValueOverride"`
	HighCValueOverride                 *int64  `json:"highCValueOverride"`
	PatternNameHigher                  string  `json:"patternNameHigher"`
	PatternNameLower                   string  `json:"patternNameLower"`
	MetricType                         string  `json:"metricType"`
	FillZero                           bool    `json:"fillZero"`
	RougeValue                         *string `json:"rougeValue"` // null or JSON string e.g. `{"l":NaN,"s":NaN}`
	EnableBaselineNearConstance        bool    `json:"enableBaselineNearConstance"`
	ComputeDifference                  bool    `json:"computeDifference"`
	AnomalyGapToleranceDuration        int64   `json:"anomalyGapToleranceDuration"`
}

// MetricAlertSettingPost is used when writing to POST /api/external/v1/componentmetricupdate.
// RougeValue uses json.RawMessage so we can inject the non-standard `{l: "NaN",s: "NaN"}` bytes verbatim.
type MetricAlertSettingPost struct {
	SMetric                            string          `json:"smetric"`
	ComponentName                      string          `json:"componentName"`
	ThresholdAlertLowerBound           string          `json:"thresholdAlertLowerBound"`
	ThresholdAlertUpperBound           string          `json:"thresholdAlertUpperBound"`
	ThresholdAlertLowerBoundNegative   string          `json:"thresholdAlertLowerBoundNegative"`
	ThresholdAlertUpperBoundNegative   string          `json:"thresholdAlertUpperBoundNegative"`
	ThresholdNoAlertLowerBound         string          `json:"thresholdNoAlertLowerBound"`
	ThresholdNoAlertUpperBound         string          `json:"thresholdNoAlertUpperBound"`
	ThresholdNoAlertLowerBoundNegative string          `json:"thresholdNoAlertLowerBoundNegative"`
	ThresholdNoAlertUpperBoundNegative string          `json:"thresholdNoAlertUpperBoundNegative"`
	IncidentAlertLowerBound            string          `json:"incidentAlertLowerBound"`
	IncidentAlertUpperBound            string          `json:"incidentAlertUpperBound"`
	IncidentAlertLowerBoundNegative    string          `json:"incidentAlertLowerBoundNegative"`
	IncidentAlertUpperBoundNegative    string          `json:"incidentAlertUpperBoundNegative"`
	IncidentNoAlertLowerBound          string          `json:"incidentNoAlertLowerBound"`
	IncidentNoAlertUpperBound          string          `json:"incidentNoAlertUpperBound"`
	IncidentNoAlertLowerBoundNegative  string          `json:"incidentNoAlertLowerBoundNegative"`
	IncidentNoAlertUpperBoundNegative  string          `json:"incidentNoAlertUpperBoundNegative"`
	IsKPI                              bool            `json:"isKPI"`
	IsFlappingResultOnly               bool            `json:"isFlappingResultOnly"`
	IncidentDurationThreshold          int64           `json:"incidentDurationThreshold"`
	DetectionType                      string          `json:"detectionType"`
	CValueOverride                     *int64          `json:"cValueOverride"`
	HighCValueOverride                 *int64          `json:"highCValueOverride"`
	PatternNameHigher                  string          `json:"patternNameHigher"`
	PatternNameLower                   string          `json:"patternNameLower"`
	MetricType                         string          `json:"metricType"`
	FillZero                           bool            `json:"fillZero"`
	RougeValue                         json.RawMessage `json:"rougeValue"` // injected verbatim: null or {l: "NaN",s: "NaN"}
	EnableBaselineNearConstance        bool            `json:"enableBaselineNearConstance"`
	ComputeDifference                  bool            `json:"computeDifference"`
	AnomalyGapToleranceDurationCount   int64           `json:"anomalyGapToleranceDurationCount"`
}

// MetricSettingEntry is the per-metric object inside the metricSetting array from the GET response.
type MetricSettingEntry struct {
	GlobalSetting             MetricAlertSetting   `json:"globalSetting"`
	ComponentLevelSettingList []MetricAlertSetting `json:"componentLevelSettingList"`
}

// GetMetricSettingsResponse is the top-level GET /api/external/v1/componentmetricupdate response.
type GetMetricSettingsResponse struct {
	MetricSettingArrCount   int                  `json:"metricSettingArrCount"`
	ReachEnd                bool                 `json:"reachEnd"`
	MetricSetting           []MetricSettingEntry `json:"metricSetting"`
	PatternIdGenerationRule int                  `json:"patternIdGenerationRule"`
}

// MetricComponentEntry is one element decoded from the componentIgnored / componentEscalateIncident JSON string.
type MetricComponentEntry struct {
	MetricLevelPrimaryKey struct {
		ProjectLevelPartitionKey struct {
			ProjectName string `json:"projectName"`
			UserName    string `json:"userName"`
		} `json:"projectLevelPartitionKey"`
		MetricName string `json:"metricName"`
	} `json:"metricLevelPrimaryKey"`
	ComponentNameSet []string `json:"componentNameSet"`
}

// MetricComponentOperationResponse is the outer wrapper from GET /api/external/v1/metriccomponent.
type MetricComponentOperationResponse struct {
	ComponentEscalateIncident string `json:"componentEscalateIncident,omitempty"`
	ComponentIgnored          string `json:"componentIgnored,omitempty"`
}

// ─── Metric Configuration API Functions ───────────────────────────────────────

// GetMetricSettings fetches alert settings for a single metric in a project.
// projectName should NOT include the @username suffix here; it is appended internally.
// metricFilter is a fuzzy search, so we paginate until reachEnd=true and return
// only the entry whose GlobalSetting.SMetric exactly matches metricName.
func (c *Client) GetMetricSettings(projectName, metricName string) (*MetricSettingEntry, error) {
	start := 0
	limit := 500

	for {
		params := url.Values{}
		params.Set("onlyIsKpi", "false")
		params.Set("onlyComputeDifference", "false")
		params.Set("projectName", projectName+"@"+c.Username)
		params.Set("start", fmt.Sprintf("%d", start))
		params.Set("limit", fmt.Sprintf("%d", limit))
		params.Set("metricFilter", metricName)
		params.Set("customerName", c.Username)
		params.Set("tzOffset", "-14400000")

		apiPath := "/api/external/v1/componentmetricupdate?" + params.Encode()
		body, statusCode, err := c.DoRequest("GET", apiPath, nil)
		if err != nil {
			return nil, err
		}
		if statusCode == 404 || statusCode == 204 {
			return nil, nil
		}
		if statusCode != 200 {
			return nil, fmt.Errorf("failed to get metric settings: HTTP %d - %s", statusCode, string(body))
		}

		var resp GetMetricSettingsResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse metric settings response: %w", err)
		}

		for i, entry := range resp.MetricSetting {
			if entry.GlobalSetting.SMetric == metricName {
				return &resp.MetricSetting[i], nil
			}
		}

		if resp.ReachEnd {
			break
		}
		start += limit
	}

	return nil, nil
}

// GetAllMetricSettings fetches all metric alert settings for a project using paginated
// requests without a metric filter, returning a map keyed by metric name.
// This is far more efficient than calling GetMetricSettings once per metric when
// reading configurations for many metrics.
func (c *Client) GetAllMetricSettings(projectName string) (map[string]*MetricSettingEntry, error) {
	result := make(map[string]*MetricSettingEntry)
	start := 0
	limit := 500

	for {
		params := url.Values{}
		params.Set("onlyIsKpi", "false")
		params.Set("onlyComputeDifference", "false")
		params.Set("projectName", projectName+"@"+c.Username)
		params.Set("start", fmt.Sprintf("%d", start))
		params.Set("limit", fmt.Sprintf("%d", limit))
		params.Set("customerName", c.Username)
		params.Set("tzOffset", "-14400000")

		apiPath := "/api/external/v1/componentmetricupdate?" + params.Encode()
		body, statusCode, err := c.DoRequest("GET", apiPath, nil)
		if err != nil {
			return nil, err
		}
		if statusCode == 404 || statusCode == 204 {
			break
		}
		if statusCode != 200 {
			return nil, fmt.Errorf("failed to get all metric settings: HTTP %d - %s", statusCode, string(body))
		}

		var resp GetMetricSettingsResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse metric settings response: %w", err)
		}

		for i, entry := range resp.MetricSetting {
			if name := entry.GlobalSetting.SMetric; name != "" {
				result[name] = &resp.MetricSetting[i]
			}
		}

		if resp.ReachEnd {
			break
		}
		start += limit
	}

	return result, nil
}

// SetMetricSettings sends a pre-serialized JSON array of MetricAlertSettingPost entries for a project.
// data is the raw JSON bytes (built by the resource layer to handle rougeValue properly).
func (c *Client) SetMetricSettings(projectName string, patternIdGenerationRule int, data []byte) error {
	fields := map[string]string{
		"projectName":             projectName + "@" + c.Username,
		"patternIdGenerationRule": fmt.Sprintf("%d", patternIdGenerationRule),
		"customerName":            c.Username,
	}
	fileParts := map[string][]byte{
		"data": data,
	}

	_ = os.WriteFile("/tmp/tf_metric_payload.json", data, 0644)

	// Retry up to 5 times: the API returns 204 or 502 when the project was just created
	// and the backend hasn't fully provisioned it yet.
	var body []byte
	var statusCode int
	var err error
	for attempt := 1; attempt <= 5; attempt++ {
		body, statusCode, err = c.DoMultipartFormRequest("POST", "/api/external/v1/componentmetricupdate", fields, fileParts)
		if err != nil {
			return err
		}
		if statusCode != 204 && statusCode != 502 {
			break
		}
		time.Sleep(time.Duration(attempt) * 3 * time.Second)
	}

	// If still 204/502 after all retries, fall back to sending each entry individually.
	if statusCode == 502 {
		return c.setMetricSettingsOneByOne(projectName, patternIdGenerationRule, data)
	}
	// 204 after all retries means the backend is not ready for this API — ignore silently.
	if statusCode == 204 {
		return nil
	}

	_ = os.WriteFile("/tmp/tf_metric_response.json", body, 0644)
	if statusCode != 200 {
		return fmt.Errorf("failed to set metric settings: HTTP %d - %s", statusCode, string(body))
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil // treat parse error as success when status was 200
	}
	if !response.Success {
		return fmt.Errorf("failed to set metric settings: %s", response.Message)
	}
	return nil
}

// setMetricSettingsOneByOne sends each entry in data as a separate request.
// Used as a fallback when the batch request keeps returning 204.
func (c *Client) setMetricSettingsOneByOne(projectName string, patternIdGenerationRule int, data []byte) error {
	var entries []json.RawMessage
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse metric settings for one-by-one fallback: %w", err)
	}

	fields := map[string]string{
		"projectName":             projectName + "@" + c.Username,
		"patternIdGenerationRule": fmt.Sprintf("%d", patternIdGenerationRule),
		"customerName":            c.Username,
	}

	for i, entry := range entries {
		entryBytes, _ := json.Marshal([]json.RawMessage{entry})
		var body []byte
		var statusCode int
		var err error
		for attempt := 1; attempt <= 5; attempt++ {
			body, statusCode, err = c.DoMultipartFormRequest("POST", "/api/external/v1/componentmetricupdate", fields, map[string][]byte{"data": entryBytes})
			if err != nil {
				return fmt.Errorf("failed to set metric setting [%d]: %w", i, err)
			}
			if statusCode != 204 && statusCode != 502 {
				break
			}
			time.Sleep(time.Duration(attempt) * 3 * time.Second)
		}
		if statusCode == 204 {
			continue
		}
		if statusCode != 200 {
			return fmt.Errorf("failed to set metric setting [%d]: HTTP %d - %s", i, statusCode, string(body))
		}
		var response ProjectResponse
		if err := json.Unmarshal(body, &response); err != nil {
			continue
		}
		if !response.Success {
			return fmt.Errorf("failed to set metric setting [%d]: %s", i, response.Message)
		}
	}
	return nil
}

// GetMetricComponents returns a map[metricName][]componentName for the given operation
// (operation = "escalateIncident" or "ignored").
func (c *Client) GetMetricComponents(projectName, operation string) (map[string][]string, error) {
	params := url.Values{}
	params.Set("projectName", projectName)
	params.Set("customerName", c.Username)
	params.Set("operation", operation)
	params.Set("tzOffset", "-14400000")

	apiPath := "/api/external/v1/metriccomponent?" + params.Encode()
	body, statusCode, err := c.DoRequest("GET", apiPath, nil)
	if err != nil {
		return nil, err
	}
	if statusCode == 404 || statusCode == 204 {
		return make(map[string][]string), nil
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get metric components (%s): HTTP %d - %s", operation, statusCode, string(body))
	}

	var outer MetricComponentOperationResponse
	if err := json.Unmarshal(body, &outer); err != nil {
		return nil, fmt.Errorf("failed to parse metric component response: %w", err)
	}

	var encodedStr string
	if operation == "escalateIncident" {
		encodedStr = outer.ComponentEscalateIncident
	} else {
		encodedStr = outer.ComponentIgnored
	}

	if encodedStr == "" {
		return make(map[string][]string), nil
	}

	var entries []MetricComponentEntry
	if err := json.Unmarshal([]byte(encodedStr), &entries); err != nil {
		return nil, fmt.Errorf("failed to parse %s entries: %w", operation, err)
	}

	result := make(map[string][]string)
	for _, e := range entries {
		result[e.MetricLevelPrimaryKey.MetricName] = e.ComponentNameSet
	}
	return result, nil
}

// SetMetricComponents sets the component list for a metric+operation pair.
// If desiredComponents contains a "Global_*" prefix component, selectAll=true is used.
// Otherwise: reset with selectAll=true, then add desired and remove the Global_* entry.
func (c *Client) SetMetricComponents(projectName, metricName, operation string, desiredComponents []string) error {
	hasGlobal := false
	for _, comp := range desiredComponents {
		if strings.HasPrefix(comp, "Global_") {
			hasGlobal = true
			break
		}
	}

	// Step 1: always reset via selectAll=true
	if err := c.postMetricComponent(projectName, metricName, operation, true, nil, nil); err != nil {
		return err
	}

	// If desired is "all" (has Global_*), we're done
	if hasGlobal {
		return nil
	}

	// If desired is empty, remove the Global_* that selectAll just added
	if len(desiredComponents) == 0 {
		// GET current to find the actual Global_* component name
		current, err := c.GetMetricComponents(projectName, operation)
		if err != nil {
			return fmt.Errorf("failed to get current components after selectAll: %w", err)
		}
		var globalComps []string
		for _, comp := range current[metricName] {
			if strings.HasPrefix(comp, "Global_") {
				globalComps = append(globalComps, comp)
			}
		}
		return c.postMetricComponent(projectName, metricName, operation, false, nil, globalComps)
	}

	// Step 2: add desired components, remove the Global_* entry
	current, err := c.GetMetricComponents(projectName, operation)
	if err != nil {
		return fmt.Errorf("failed to get current components after selectAll: %w", err)
	}
	var globalComps []string
	for _, comp := range current[metricName] {
		if strings.HasPrefix(comp, "Global_") {
			globalComps = append(globalComps, comp)
		}
	}
	return c.postMetricComponent(projectName, metricName, operation, false, desiredComponents, globalComps)
}

// postMetricComponent is the low-level POST to /api/external/v1/metriccomponent.
func (c *Client) postMetricComponent(projectName, metricName, operation string, selectAll bool, addSet, removeSet []string) error {
	formData := url.Values{}
	formData.Set("projectName", projectName)
	formData.Set("customerName", c.Username)
	formData.Set("metricName", metricName)
	formData.Set("operation", operation)
	if selectAll {
		formData.Set("selectAll", "true")
	} else {
		formData.Set("selectAll", "false")
		if len(addSet) > 0 {
			addJSON, _ := json.Marshal(addSet)
			formData.Set("addComponentSet", string(addJSON))
		}
		if len(removeSet) > 0 {
			removeJSON, _ := json.Marshal(removeSet)
			formData.Set("removeComponentSet", string(removeJSON))
		}
	}

	var body []byte
	var statusCode int
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		body, statusCode, err = c.DoFormRequest("POST", "/api/external/v1/metriccomponent", formData)
		if err != nil {
			return err
		}
		if statusCode != 502 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if statusCode != 200 {
		return fmt.Errorf("failed to set metric component (%s/%s): HTTP %d - %s", metricName, operation, statusCode, string(body))
	}

	var response ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}
	if !response.Success {
		return fmt.Errorf("failed to set metric component (%s/%s): %s", metricName, operation, response.Message)
	}
	return nil
}

// ConvertRougeValueForPost converts a Terraform-stored rougeValue string to the raw bytes
// for the POST API, or null.
// GET response delivers rougeValue as a JSON string like `{"l":NaN,"s":NaN}` (unquoted NaN).
// This function normalizes that into valid JSON `{"l":"NaN","s":"NaN"}` for the POST body.
var defaultRougeValue = json.RawMessage(`{"l":"NaN","s":"NaN"}`)

func ConvertRougeValueForPost(storedVal string) json.RawMessage {
	if storedVal == "" || storedVal == "null" {
		return defaultRougeValue
	}
	// Normalize: replace unquoted NaN → "NaN" to make it valid JSON for parsing
	normalized := strings.ReplaceAll(storedVal, ":NaN", `:"NaN"`)
	normalized = strings.ReplaceAll(normalized, ": NaN", `:"NaN"`)
	// Use interface{} so numeric values (e.g. 3.0) and string "NaN" are both handled correctly.
	var rv map[string]interface{}
	if err := json.Unmarshal([]byte(normalized), &rv); err == nil {
		out, _ := json.Marshal(rv)
		return json.RawMessage(out)
	}
	return defaultRougeValue
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
