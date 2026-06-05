// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

var (
	_ resource.Resource                = &metricProjectResource{}
	_ resource.ResourceWithConfigure   = &metricProjectResource{}
	_ resource.ResourceWithImportState = &metricProjectResource{}
)

// useStateIfConfigNullModifier suppresses drift for a string attribute: if the
// config does not set the field (null), the prior state value is kept so the
// plan shows no change even when the API returns a different value.
type useStateIfConfigNullModifier struct{ desc string }

func (m useStateIfConfigNullModifier) Description(_ context.Context) string         { return m.desc }
func (m useStateIfConfigNullModifier) MarkdownDescription(_ context.Context) string { return m.desc }
func (m useStateIfConfigNullModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}
}

func ignoreStringDrift() planmodifier.String {
	return useStateIfConfigNullModifier{desc: "Preserves state value when config is null, suppressing API drift."}
}

type useStateIfConfigNullListModifier struct{ desc string }

func (m useStateIfConfigNullListModifier) Description(_ context.Context) string { return m.desc }
func (m useStateIfConfigNullListModifier) MarkdownDescription(_ context.Context) string {
	return m.desc
}
func (m useStateIfConfigNullListModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}
}

func ignoreListDrift() planmodifier.List {
	return useStateIfConfigNullListModifier{desc: "Preserves state value when config is null, suppressing API drift."}
}

func NewMetricProjectResource() resource.Resource {
	return &metricProjectResource{}
}

type metricProjectResource struct {
	client *client.Client
}

type metricProjectResourceModel struct {
	ID                    types.String                      `tfsdk:"id"`
	ProjectName           types.String                      `tfsdk:"project_name"`
	ProjectDisplayName    types.String                      `tfsdk:"project_display_name"`
	SystemName            types.String                      `tfsdk:"system_name"`
	ProjectCreationConfig *metricProjectCreationConfigModel `tfsdk:"project_creation_config"`

	// Common settings
	CValue           types.Int64   `tfsdk:"c_value"`
	PValue           types.Float64 `tfsdk:"p_value"`
	ProjectTimeZone  types.String  `tfsdk:"project_time_zone"`
	SamplingInterval types.Int64   `tfsdk:"sampling_interval"`

	RetentionTime       types.Int64  `tfsdk:"retention_time"`
	UBLRetentionTime    types.Int64  `tfsdk:"ubl_retention_time"`
	TrainingFilter      types.Bool   `tfsdk:"training_filter"`
	EnableNewAlertEmail types.Bool   `tfsdk:"enable_new_alert_email"`
	LargeProject        types.Bool   `tfsdk:"large_project"`
	NewPatternRange     types.Int64  `tfsdk:"new_pattern_range"`
	Proxy               types.String `tfsdk:"proxy"`

	EnableAnomalyScoreEscalation    types.Bool    `tfsdk:"enable_anomaly_score_escalation"`
	EscalationAnomalyScoreThreshold types.String  `tfsdk:"escalation_anomaly_score_threshold"`
	IgnoreAnomalyScoreThreshold     types.String  `tfsdk:"ignore_anomaly_score_threshold"`
	EnableStreamDetection           types.Bool    `tfsdk:"enable_stream_detection"`
	IgnoreInstanceForKB             types.Bool    `tfsdk:"ignore_instance_for_kb"`
	ShowInstanceDown                types.Bool    `tfsdk:"show_instance_down"`
	InstanceDownEnable              types.Bool    `tfsdk:"instance_down_enable"`
	AlertHourlyCost                 types.Float64 `tfsdk:"alert_hourly_cost"`
	AlertAverageTime                types.Int64   `tfsdk:"alert_average_time"`
	AvgPerIncidentDowntimeCost      types.Float64 `tfsdk:"avg_per_incident_downtime_cost"`

	// Incident prediction and RCA
	IncidentPredictionWindow             types.Int64   `tfsdk:"incident_prediction_window"`
	MinIncidentPredictionWindow          types.Int64   `tfsdk:"min_incident_prediction_window"`
	IncidentRelationSearchWindow         types.Int64   `tfsdk:"incident_relation_search_window"`
	IncidentPredictionEventLimit         types.Int64   `tfsdk:"incident_prediction_event_limit"`
	RootCauseCountThreshold              types.Int64   `tfsdk:"root_cause_count_threshold"`
	RootCauseProbabilityThreshold        types.Float64 `tfsdk:"root_cause_probability_threshold"`
	CompositeRCALimit                    types.Int64   `tfsdk:"composite_rca_limit"`
	RootCauseLogMessageSearchRange       types.Int64   `tfsdk:"root_cause_log_message_search_range"`
	CausalPredictionSetting              types.Int64   `tfsdk:"causal_prediction_setting"`
	RootCauseRankSetting                 types.Int64   `tfsdk:"root_cause_rank_setting"`
	MaximumRootCauseResultSize           types.Int64   `tfsdk:"maximum_root_cause_result_size"`
	MultiHopSearchLevel                  types.Int64   `tfsdk:"multi_hop_search_level"`
	MultiHopSearchLimit                  types.String  `tfsdk:"multi_hop_search_limit"`
	PredictionCountThreshold             types.Int64   `tfsdk:"prediction_count_threshold"`
	PredictionProbabilityThreshold       types.Float64 `tfsdk:"prediction_probability_threshold"`
	PredictionRuleActiveCondition        types.Int64   `tfsdk:"prediction_rule_active_condition"`
	PredictionRuleActiveThreshold        types.Float64 `tfsdk:"prediction_rule_active_threshold"`
	PredictionRuleFalsePositiveThreshold types.Int64   `tfsdk:"prediction_rule_false_positive_threshold"`
	PredictionRuleInactiveThreshold      types.Float64 `tfsdk:"prediction_rule_inactive_threshold"`
	MinValidModelSpan                    types.Int64   `tfsdk:"min_valid_model_span"`

	// Webhook
	MaxWebHookRequestSize        types.Int64  `tfsdk:"max_web_hook_request_size"`
	WebhookAlertDampening        types.Int64  `tfsdk:"webhook_alert_dampening"`
	WebhookBlackListSetStr       types.String `tfsdk:"webhook_black_list_set_str"`
	WebhookCriticalKeywordSetStr types.String `tfsdk:"webhook_critical_keyword_set_str"`
	WebhookTypeSetStr            types.String `tfsdk:"webhook_type_set_str"`
	WebhookUrl                   types.String `tfsdk:"webhook_url"`

	// Metric-specific fields
	HighRatioCValue                     types.Int64   `tfsdk:"high_ratio_c_value"`
	MaximumHint                         types.Int64   `tfsdk:"maximum_hint"`
	DynamicBaselineDetectionFlag        types.Bool    `tfsdk:"dynamic_baseline_detection_flag"`
	PositiveBaselineViolationFactor     types.Float64 `tfsdk:"positive_baseline_violation_factor"`
	NegativeBaselineViolationFactor     types.Float64 `tfsdk:"negative_baseline_violation_factor"`
	EnablePeriodAnomalyFilter           types.Bool    `tfsdk:"enable_period_anomaly_filter"`
	EnableUBLDetect                     types.Bool    `tfsdk:"enable_ubl_detect"`
	EnableCumulativeDetect              types.Bool    `tfsdk:"enable_cumulative_detect"`
	EnableComponentLevelDetection       types.Bool    `tfsdk:"enable_component_level_detection"`
	PredictionTrainingDataLength        types.Int64   `tfsdk:"prediction_training_data_length"`
	PredictionCorrelationSensitivity    types.Float64 `tfsdk:"prediction_correlation_sensitivity"`
	EnableKPIPrediction                 types.Bool    `tfsdk:"enable_kpi_prediction"`
	InstanceDownThreshold               types.Int64   `tfsdk:"instance_down_threshold"`
	InstanceDownReportNumber            types.Int64   `tfsdk:"instance_down_report_number"`
	ModelSpan                           types.Int64   `tfsdk:"model_span"`
	EnableMetricDataPrediction          types.Bool    `tfsdk:"enable_metric_data_prediction"`
	EnableBaselineDetectionDoubleVerify types.Bool    `tfsdk:"enable_baseline_detection_double_verify"`
	EnableFillGap                       types.Bool    `tfsdk:"enable_fill_gap"`
	EnableStoreFilledGap                types.Bool    `tfsdk:"enable_store_filled_gap"`
	GapFillingTrainingDataLength        types.Int64   `tfsdk:"gap_filling_training_data_length"`
	PatternIdGenerationRule             types.Int64   `tfsdk:"pattern_id_generation_rule"`
	AnomalyGapToleranceCount            types.Int64   `tfsdk:"anomaly_gap_tolerance_count"`
	FilterByAnomalyInBaselineGeneration types.Bool    `tfsdk:"filter_by_anomaly_in_baseline_generation"`
	BaselineDuration                    types.Int64   `tfsdk:"baseline_duration"`
	AnomalyDampening                    types.Int64   `tfsdk:"anomaly_dampening"`
	InstanceDownRatioThreshold          types.Float64 `tfsdk:"instance_down_ratio_threshold"`
	ComponentNameAutoOverwrite          types.Bool    `tfsdk:"component_name_auto_overwrite"`

	// Complex JSON fields
	LinkedLogProjects                      types.String `tfsdk:"linked_log_projects"`
	ComponentMetricSettingOverallModelList types.String `tfsdk:"component_metric_setting_overall_model_list"`
	EmailSetting                           types.String `tfsdk:"email_setting"`
	InstanceGroupingUpdate                 types.String `tfsdk:"instance_grouping_update"`
	SharedUsernames                        types.String `tfsdk:"shared_usernames"`
	WebhookHeaderList                      types.String `tfsdk:"webhook_header_list"`

	// Holiday settings
	HolidaySettings types.List `tfsdk:"holiday_settings"`

	// Metric configurations (per-metric alert thresholds + component operations)
	MetricConfigurations types.Map `tfsdk:"metric_configurations"`

	Mode types.Int64 `tfsdk:"mode"`
}

// metricAlertSettingModel maps one entry in metric_alert_settings.
type metricAlertSettingModel struct {
	ComponentName                      types.String `tfsdk:"component_name"`
	ThresholdAlertLowerBound           types.String `tfsdk:"threshold_alert_lower_bound"`
	ThresholdAlertUpperBound           types.String `tfsdk:"threshold_alert_upper_bound"`
	ThresholdAlertLowerBoundNegative   types.String `tfsdk:"threshold_alert_lower_bound_negative"`
	ThresholdAlertUpperBoundNegative   types.String `tfsdk:"threshold_alert_upper_bound_negative"`
	ThresholdNoAlertLowerBound         types.String `tfsdk:"threshold_no_alert_lower_bound"`
	ThresholdNoAlertUpperBound         types.String `tfsdk:"threshold_no_alert_upper_bound"`
	ThresholdNoAlertLowerBoundNegative types.String `tfsdk:"threshold_no_alert_lower_bound_negative"`
	ThresholdNoAlertUpperBoundNegative types.String `tfsdk:"threshold_no_alert_upper_bound_negative"`
	IncidentAlertLowerBound            types.String `tfsdk:"incident_alert_lower_bound"`
	IncidentAlertUpperBound            types.String `tfsdk:"incident_alert_upper_bound"`
	IncidentAlertLowerBoundNegative    types.String `tfsdk:"incident_alert_lower_bound_negative"`
	IncidentAlertUpperBoundNegative    types.String `tfsdk:"incident_alert_upper_bound_negative"`
	IncidentNoAlertLowerBound          types.String `tfsdk:"incident_no_alert_lower_bound"`
	IncidentNoAlertUpperBound          types.String `tfsdk:"incident_no_alert_upper_bound"`
	IncidentNoAlertLowerBoundNegative  types.String `tfsdk:"incident_no_alert_lower_bound_negative"`
	IncidentNoAlertUpperBoundNegative  types.String `tfsdk:"incident_no_alert_upper_bound_negative"`
	IsKPI                              types.Bool   `tfsdk:"is_kpi"`
	IsFlappingResultOnly               types.Bool   `tfsdk:"is_flapping_result_only"`
	IncidentDurationThreshold          types.Int64  `tfsdk:"incident_duration_threshold"`
	DetectionType                      types.String `tfsdk:"detection_type"`
	CValueOverride                     types.Int64  `tfsdk:"c_value_override"`
	HighCValueOverride                 types.Int64  `tfsdk:"high_c_value_override"`
	PatternNameHigher                  types.String `tfsdk:"pattern_name_higher"`
	PatternNameLower                   types.String `tfsdk:"pattern_name_lower"`
	MetricType                         types.String `tfsdk:"metric_type"`
	FillZero                           types.Bool   `tfsdk:"fill_zero"`
	RougeValue                         types.String `tfsdk:"rouge_value"` // null or raw string from API e.g. `{"l":NaN,"s":NaN}`
	EnableBaselineNearConstance        types.Bool   `tfsdk:"enable_baseline_near_constance"`
	ComputeDifference                  types.Bool   `tfsdk:"compute_difference"`
	AnomalyGapToleranceDuration        types.Int64  `tfsdk:"anomaly_gap_tolerance_duration"`
}

// metricConfigurationModel is the internal representation of one metric_configurations entry.
// MetricName holds the map key; used internally and not exposed to tfsdk directly.
type metricConfigurationModel struct {
	MetricName                 types.String
	EscalateIncidentComponents types.List // list of strings
	IgnoredComponents          types.List // list of strings
	MetricAlertSettings        types.List // list of metricAlertSettingModel
}

// metricConfigurationValueModel is the tfsdk-facing value type for metric_configurations map entries.
// The map key is the metric_name.
type metricConfigurationValueModel struct {
	EscalateIncidentComponents types.List `tfsdk:"escalate_incident_components"`
	IgnoredComponents          types.List `tfsdk:"ignored_components"`
	MetricAlertSettings        types.List `tfsdk:"metric_alert_settings"`
}

type metricProjectCreationConfigModel struct {
	DataType         types.String `tfsdk:"data_type"`
	InstanceType     types.String `tfsdk:"instance_type"`
	ProjectCloudType types.String `tfsdk:"project_cloud_type"`
	InsightAgentType types.String `tfsdk:"insight_agent_type"`
}

// metricAlertSettingAttrTypes returns the attr.Type map for a metricAlertSettingModel object.
func metricAlertSettingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"component_name":                          types.StringType,
		"threshold_alert_lower_bound":             types.StringType,
		"threshold_alert_upper_bound":             types.StringType,
		"threshold_alert_lower_bound_negative":    types.StringType,
		"threshold_alert_upper_bound_negative":    types.StringType,
		"threshold_no_alert_lower_bound":          types.StringType,
		"threshold_no_alert_upper_bound":          types.StringType,
		"threshold_no_alert_lower_bound_negative": types.StringType,
		"threshold_no_alert_upper_bound_negative": types.StringType,
		"incident_alert_lower_bound":              types.StringType,
		"incident_alert_upper_bound":              types.StringType,
		"incident_alert_lower_bound_negative":     types.StringType,
		"incident_alert_upper_bound_negative":     types.StringType,
		"incident_no_alert_lower_bound":           types.StringType,
		"incident_no_alert_upper_bound":           types.StringType,
		"incident_no_alert_lower_bound_negative":  types.StringType,
		"incident_no_alert_upper_bound_negative":  types.StringType,
		"is_kpi":                                  types.BoolType,
		"is_flapping_result_only":                 types.BoolType,
		"incident_duration_threshold":             types.Int64Type,
		"detection_type":                          types.StringType,
		"c_value_override":                        types.Int64Type,
		"high_c_value_override":                   types.Int64Type,
		"pattern_name_higher":                     types.StringType,
		"pattern_name_lower":                      types.StringType,
		"metric_type":                             types.StringType,
		"fill_zero":                               types.BoolType,
		"rouge_value":                             types.StringType,
		"enable_baseline_near_constance":          types.BoolType,
		"compute_difference":                      types.BoolType,
		"anomaly_gap_tolerance_duration":          types.Int64Type,
	}
}

// metricConfigurationValueAttrTypes returns the attr.Type map for a metricConfigurationValueModel
// (the map value; metric_name is the map key and not repeated here).
func metricConfigurationValueAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"escalate_incident_components": types.ListType{ElemType: types.StringType},
		"ignored_components":           types.ListType{ElemType: types.StringType},
		"metric_alert_settings":        types.ListType{ElemType: types.ObjectType{AttrTypes: metricAlertSettingAttrTypes()}},
	}
}

// metricConfigsFromMap extracts the metric_configurations map from state into the internal slice.
func metricConfigsFromMap(ctx context.Context, m types.Map) ([]metricConfigurationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	rawMap := make(map[string]metricConfigurationValueModel)
	diags.Append(m.ElementsAs(ctx, &rawMap, false)...)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]metricConfigurationModel, 0, len(rawMap))
	for k, v := range rawMap {
		result = append(result, metricConfigurationModel{
			MetricName:                 types.StringValue(k),
			EscalateIncidentComponents: v.EscalateIncidentComponents,
			IgnoredComponents:          v.IgnoredComponents,
			MetricAlertSettings:        v.MetricAlertSettings,
		})
	}
	return result, diags
}

// metricConfigsToMap converts the internal slice back to a types.Map for state storage.
func metricConfigsToMap(ctx context.Context, configs []metricConfigurationModel) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: metricConfigurationValueAttrTypes()}
	elements := make(map[string]attr.Value, len(configs))
	for _, c := range configs {
		objVal, d := types.ObjectValueFrom(ctx, metricConfigurationValueAttrTypes(), metricConfigurationValueModel{
			EscalateIncidentComponents: c.EscalateIncidentComponents,
			IgnoredComponents:          c.IgnoredComponents,
			MetricAlertSettings:        c.MetricAlertSettings,
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(elemType), diags
		}
		elements[c.MetricName.ValueString()] = objVal
	}
	mapVal, d := types.MapValue(elemType, elements)
	diags.Append(d...)
	return mapVal, diags
}

func (r *metricProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric_project"
}

func (r *metricProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an InsightFinder metric project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the project (same as project_name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_name": schema.StringAttribute{
				Description: "The name of the project (must be unique).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_display_name": schema.StringAttribute{
				Description: "The display name for the project.",
				Optional:    true,
				Computed:    true,
			},
			"system_name": schema.StringAttribute{
				Description: "The system name this project belongs to.",
				Required:    true,
			},
			"project_creation_config": schema.SingleNestedAttribute{
				Description: "Configuration for creating the project.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"data_type": schema.StringAttribute{
						Description: "The type of data (e.g., Metric).",
						Required:    true,
					},
					"instance_type": schema.StringAttribute{
						Description: "The instance type (e.g., PrivateCloud, AWS, Azure).",
						Required:    true,
					},
					"project_cloud_type": schema.StringAttribute{
						Description: "The cloud type for the project.",
						Required:    true,
					},
					"insight_agent_type": schema.StringAttribute{
						Description: "The InsightFinder agent type.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"c_value": schema.Int64Attribute{
				Description: "The C value for anomaly detection sensitivity (typically 2-5).",
				Optional:    true,
				Computed:    true,
			},
			"p_value": schema.Float64Attribute{
				Description: "The P value for anomaly detection probability (0.0-1.0).",
				Optional:    true,
				Computed:    true,
			},
			"project_time_zone": schema.StringAttribute{
				Description: "The timezone for the project (default: UTC).",
				Optional:    true,
				Computed:    true,
			},
			"sampling_interval": schema.Int64Attribute{
				Description: "The sampling interval in seconds.",
				Optional:    true,
				Computed:    true,
			},
			"retention_time": schema.Int64Attribute{
				Description: "The retention time in days.",
				Optional:    true,
				Computed:    true,
			},
			"ubl_retention_time": schema.Int64Attribute{
				Description: "Retention time for UBL data in days.",
				Optional:    true,
				Computed:    true,
			},
			"training_filter": schema.BoolAttribute{
				Description: "Training filter flag.",
				Optional:    true,
				Computed:    true,
			},
			"enable_new_alert_email": schema.BoolAttribute{
				Description: "Enable new alert email notifications.",
				Optional:    true,
				Computed:    true,
			},
			"large_project": schema.BoolAttribute{
				Description: "Is this a large project.",
				Optional:    true,
				Computed:    true,
			},
			"new_pattern_range": schema.Int64Attribute{
				Description: "Range for new patterns.",
				Optional:    true,
				Computed:    true,
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy configuration.",
				Optional:    true,
				Computed:    true,
			},
			"enable_anomaly_score_escalation": schema.BoolAttribute{
				Description: "Enable anomaly score escalation.",
				Optional:    true,
				Computed:    true,
			},
			"escalation_anomaly_score_threshold": schema.StringAttribute{
				Description: "Threshold for anomaly score escalation.",
				Optional:    true,
				Computed:    true,
			},
			"ignore_anomaly_score_threshold": schema.StringAttribute{
				Description: "Threshold to ignore anomaly scores.",
				Optional:    true,
				Computed:    true,
			},
			"enable_stream_detection": schema.BoolAttribute{
				Description: "Enable stream detection.",
				Optional:    true,
				Computed:    true,
			},
			"ignore_instance_for_kb": schema.BoolAttribute{
				Description: "Ignore instance for knowledge base.",
				Optional:    true,
				Computed:    true,
			},
			"show_instance_down": schema.BoolAttribute{
				Description: "Whether to show instance down incidents for this project.",
				Optional:    true,
				Computed:    true,
			},
			"instance_down_enable": schema.BoolAttribute{
				Description: "Enable instance down detection.",
				Optional:    true,
				Computed:    true,
			},
			"alert_hourly_cost": schema.Float64Attribute{
				Description: "Hourly cost for alerts.",
				Optional:    true,
				Computed:    true,
			},
			"alert_average_time": schema.Int64Attribute{
				Description: "Average time for alerts.",
				Optional:    true,
				Computed:    true,
			},
			"avg_per_incident_downtime_cost": schema.Float64Attribute{
				Description: "Average cost per incident downtime.",
				Optional:    true,
				Computed:    true,
			},
			"incident_prediction_window": schema.Int64Attribute{
				Description: "Window for incident prediction.",
				Optional:    true,
				Computed:    true,
			},
			"min_incident_prediction_window": schema.Int64Attribute{
				Description: "Minimum incident prediction window.",
				Optional:    true,
				Computed:    true,
			},
			"incident_relation_search_window": schema.Int64Attribute{
				Description: "Window for incident relation search.",
				Optional:    true,
				Computed:    true,
			},
			"incident_prediction_event_limit": schema.Int64Attribute{
				Description: "Event limit for incident prediction.",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_count_threshold": schema.Int64Attribute{
				Description: "Threshold for root cause count.",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_probability_threshold": schema.Float64Attribute{
				Description: "Threshold for root cause probability.",
				Optional:    true,
				Computed:    true,
			},
			"composite_rca_limit": schema.Int64Attribute{
				Description: "Limit for composite root cause analysis.",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_log_message_search_range": schema.Int64Attribute{
				Description: "Search range for root cause log messages.",
				Optional:    true,
				Computed:    true,
			},
			"causal_prediction_setting": schema.Int64Attribute{
				Description: "Causal prediction setting.",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_rank_setting": schema.Int64Attribute{
				Description: "Rank setting for root cause.",
				Optional:    true,
				Computed:    true,
			},
			"maximum_root_cause_result_size": schema.Int64Attribute{
				Description: "Maximum root cause result size.",
				Optional:    true,
				Computed:    true,
			},
			"multi_hop_search_level": schema.Int64Attribute{
				Description: "Multi-hop search level.",
				Optional:    true,
				Computed:    true,
			},
			"multi_hop_search_limit": schema.StringAttribute{
				Description: "Multi-hop search limit.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_count_threshold": schema.Int64Attribute{
				Description: "Threshold for prediction count.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_probability_threshold": schema.Float64Attribute{
				Description: "Threshold for prediction probability.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_active_condition": schema.Int64Attribute{
				Description: "Active condition for prediction rules.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_active_threshold": schema.Float64Attribute{
				Description: "Active threshold for prediction rules.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_false_positive_threshold": schema.Int64Attribute{
				Description: "False positive threshold for prediction rules.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_inactive_threshold": schema.Float64Attribute{
				Description: "Inactive threshold for prediction rules.",
				Optional:    true,
				Computed:    true,
			},
			"min_valid_model_span": schema.Int64Attribute{
				Description: "Minimum valid model span.",
				Optional:    true,
				Computed:    true,
			},
			"max_web_hook_request_size": schema.Int64Attribute{
				Description: "Maximum webhook request size.",
				Optional:    true,
				Computed:    true,
			},
			"webhook_alert_dampening": schema.Int64Attribute{
				Description: "Alert dampening for webhooks.",
				Optional:    true,
				Computed:    true,
			},
			"webhook_black_list_set_str": schema.StringAttribute{
				Description: "Blacklist set string for webhooks (JSON).",
				Optional:    true,
				Computed:    true,
			},
			"webhook_critical_keyword_set_str": schema.StringAttribute{
				Description: "Critical keyword set string for webhooks (JSON).",
				Optional:    true,
				Computed:    true,
			},
			"webhook_type_set_str": schema.StringAttribute{
				Description: "Type set string for webhooks (JSON).",
				Optional:    true,
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "Webhook URL.",
				Optional:    true,
				Computed:    true,
			},
			// Metric-specific
			"high_ratio_c_value": schema.Int64Attribute{
				Description: "High ratio C value for anomaly detection.",
				Optional:    true,
				Computed:    true,
			},
			"maximum_hint": schema.Int64Attribute{
				Description: "Maximum hint value.",
				Optional:    true,
				Computed:    true,
			},
			"dynamic_baseline_detection_flag": schema.BoolAttribute{
				Description: "Enable dynamic baseline detection.",
				Optional:    true,
				Computed:    true,
			},
			"positive_baseline_violation_factor": schema.Float64Attribute{
				Description: "Positive baseline violation factor.",
				Optional:    true,
				Computed:    true,
			},
			"negative_baseline_violation_factor": schema.Float64Attribute{
				Description: "Negative baseline violation factor.",
				Optional:    true,
				Computed:    true,
			},
			"enable_period_anomaly_filter": schema.BoolAttribute{
				Description: "Enable period anomaly filter.",
				Optional:    true,
				Computed:    true,
			},
			"enable_ubl_detect": schema.BoolAttribute{
				Description: "Enable UBL detection.",
				Optional:    true,
				Computed:    true,
			},
			"enable_cumulative_detect": schema.BoolAttribute{
				Description: "Enable cumulative detection.",
				Optional:    true,
				Computed:    true,
			},
			"enable_component_level_detection": schema.BoolAttribute{
				Description: "Enable component level detection.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_training_data_length": schema.Int64Attribute{
				Description: "Length of training data used for prediction.",
				Optional:    true,
				Computed:    true,
			},
			"prediction_correlation_sensitivity": schema.Float64Attribute{
				Description: "Sensitivity for prediction correlation.",
				Optional:    true,
				Computed:    true,
			},
			"enable_kpi_prediction": schema.BoolAttribute{
				Description: "Enable KPI prediction.",
				Optional:    true,
				Computed:    true,
			},
			"instance_down_threshold": schema.Int64Attribute{
				Description: "Threshold (ms) before an instance is considered down.",
				Optional:    true,
				Computed:    true,
			},
			"instance_down_report_number": schema.Int64Attribute{
				Description: "Number of instances down before reporting.",
				Optional:    true,
				Computed:    true,
			},
			"model_span": schema.Int64Attribute{
				Description: "Model span setting.",
				Optional:    true,
				Computed:    true,
			},
			"enable_metric_data_prediction": schema.BoolAttribute{
				Description: "Enable metric data prediction.",
				Optional:    true,
				Computed:    true,
			},
			"enable_baseline_detection_double_verify": schema.BoolAttribute{
				Description: "Enable double verification for baseline detection.",
				Optional:    true,
				Computed:    true,
			},
			"enable_fill_gap": schema.BoolAttribute{
				Description: "Enable gap filling for metric data.",
				Optional:    true,
				Computed:    true,
			},
			"enable_store_filled_gap": schema.BoolAttribute{
				Description: "Enable storing filled gap data.",
				Optional:    true,
				Computed:    true,
			},
			"gap_filling_training_data_length": schema.Int64Attribute{
				Description: "Length of training data for gap filling.",
				Optional:    true,
				Computed:    true,
			},
			"pattern_id_generation_rule": schema.Int64Attribute{
				Description: "Rule for pattern ID generation.",
				Optional:    true,
				Computed:    true,
			},
			"anomaly_gap_tolerance_count": schema.Int64Attribute{
				Description: "Number of gaps tolerated before reporting anomaly.",
				Optional:    true,
				Computed:    true,
			},
			"filter_by_anomaly_in_baseline_generation": schema.BoolAttribute{
				Description: "Filter anomalies when generating baseline.",
				Optional:    true,
				Computed:    true,
			},
			"baseline_duration": schema.Int64Attribute{
				Description: "Duration (ms) for baseline calculation.",
				Optional:    true,
				Computed:    true,
			},
			"anomaly_dampening": schema.Int64Attribute{
				Description: "Dampening period (ms) between anomaly alerts.",
				Optional:    true,
				Computed:    true,
			},
			"instance_down_ratio_threshold": schema.Float64Attribute{
				Description: "Ratio threshold for instance down detection.",
				Optional:    true,
				Computed:    true,
			},
			"component_name_auto_overwrite": schema.BoolAttribute{
				Description: "Auto-overwrite component names.",
				Optional:    true,
				Computed:    true,
			},
			// Complex JSON fields
			"linked_log_projects": schema.StringAttribute{
				Description: "List of linked log projects (JSON array).",
				Optional:    true,
				Computed:    true,
			},
			"component_metric_setting_overall_model_list": schema.StringAttribute{
				Description: "Component metric setting overall model list (JSON array).",
				Optional:    true,
				Computed:    true,
			},
			"email_setting": schema.StringAttribute{
				Description: "Email notification settings (JSON).",
				Optional:    true,
				Computed:    true,
			},
			"instance_grouping_update": schema.StringAttribute{
				Description: "Instance grouping update settings (JSON).",
				Optional:    true,
				Computed:    true,
			},
			"shared_usernames": schema.StringAttribute{
				Description: "List of shared usernames (JSON array).",
				Optional:    true,
				Computed:    true,
			},
			"webhook_header_list": schema.StringAttribute{
				Description: "List of webhook headers (JSON array).",
				Optional:    true,
				Computed:    true,
			},
			"holiday_settings": schema.ListNestedAttribute{
				Description: "List of holiday settings for the project. Each holiday has a name, start date, and end date (MM-DD format).",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the holiday.",
							Required:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Start date of the holiday in MM-DD format (e.g., '12-25').",
							Required:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "End date of the holiday in MM-DD format (e.g., '12-26').",
							Required:    true,
						},
					},
				},
			},
			"mode": schema.Int64Attribute{
				Description: "The process mode for the project (set via logdedicatedmode API).",
				Optional:    true,
				Computed:    true,
			},
			"metric_configurations": schema.MapNestedAttribute{
				Description: "Per-metric alert threshold and component configurations, keyed by metric name.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"escalate_incident_components": schema.ListAttribute{
							Description:   "Components for which incidents are escalated. Use ['Global_<hash>'] to select all.",
							ElementType:   types.StringType,
							Optional:      true,
							Computed:      true,
							PlanModifiers: []planmodifier.List{ignoreListDrift()},
						},
						"ignored_components": schema.ListAttribute{
							Description:   "Components that are ignored for this metric. Use ['Global_<hash>'] to select all.",
							ElementType:   types.StringType,
							Optional:      true,
							Computed:      true,
							PlanModifiers: []planmodifier.List{ignoreListDrift()},
						},
						"metric_alert_settings": schema.ListNestedAttribute{
							Description: "Alert threshold settings per component for this metric.",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"component_name":                          schema.StringAttribute{Description: "Component name (use Global_<hash> for the global setting).", Required: true},
									"threshold_alert_lower_bound":             schema.StringAttribute{Description: "Alert lower bound threshold.", Optional: true, Computed: true},
									"threshold_alert_upper_bound":             schema.StringAttribute{Description: "Alert upper bound threshold.", Optional: true, Computed: true},
									"threshold_alert_lower_bound_negative":    schema.StringAttribute{Description: "Alert lower bound negative threshold.", Optional: true, Computed: true},
									"threshold_alert_upper_bound_negative":    schema.StringAttribute{Description: "Alert upper bound negative threshold.", Optional: true, Computed: true},
									"threshold_no_alert_lower_bound":          schema.StringAttribute{Description: "No-alert lower bound threshold.", Optional: true, Computed: true},
									"threshold_no_alert_upper_bound":          schema.StringAttribute{Description: "No-alert upper bound threshold.", Optional: true, Computed: true},
									"threshold_no_alert_lower_bound_negative": schema.StringAttribute{Description: "No-alert lower bound negative threshold.", Optional: true, Computed: true},
									"threshold_no_alert_upper_bound_negative": schema.StringAttribute{Description: "No-alert upper bound negative threshold.", Optional: true, Computed: true},
									"incident_alert_lower_bound":              schema.StringAttribute{Description: "Incident alert lower bound.", Optional: true, Computed: true},
									"incident_alert_upper_bound":              schema.StringAttribute{Description: "Incident alert upper bound.", Optional: true, Computed: true},
									"incident_alert_lower_bound_negative":     schema.StringAttribute{Description: "Incident alert lower bound negative.", Optional: true, Computed: true},
									"incident_alert_upper_bound_negative":     schema.StringAttribute{Description: "Incident alert upper bound negative.", Optional: true, Computed: true},
									"incident_no_alert_lower_bound":           schema.StringAttribute{Description: "Incident no-alert lower bound.", Optional: true, Computed: true},
									"incident_no_alert_upper_bound":           schema.StringAttribute{Description: "Incident no-alert upper bound.", Optional: true, Computed: true},
									"incident_no_alert_lower_bound_negative":  schema.StringAttribute{Description: "Incident no-alert lower bound negative.", Optional: true, Computed: true},
									"incident_no_alert_upper_bound_negative":  schema.StringAttribute{Description: "Incident no-alert upper bound negative.", Optional: true, Computed: true},
									"is_kpi":                                  schema.BoolAttribute{Description: "Whether this metric is a KPI.", Optional: true, Computed: true},
									"is_flapping_result_only":                 schema.BoolAttribute{Description: "Whether to report flapping results only.", Optional: true},
									"incident_duration_threshold":             schema.Int64Attribute{Description: "Minimum incident duration (ms) to trigger.", Optional: true},
									"detection_type":                          schema.StringAttribute{Description: "Detection direction: 'positive', 'negative', or 'both'.", Optional: true},
									"c_value_override": schema.Int64Attribute{
										Description: "Override for the C value anomaly sensitivity. Null means use project default.",
										Optional:    true,
									},
									"high_c_value_override": schema.Int64Attribute{
										Description: "Override for the high-ratio C value anomaly sensitivity. Null means use project default.",
										Optional:    true,
									},
									"pattern_name_higher":            schema.StringAttribute{Description: "Pattern name for higher anomalies.", Optional: true},
									"pattern_name_lower":             schema.StringAttribute{Description: "Pattern name for lower anomalies.", Optional: true},
									"metric_type":                    schema.StringAttribute{Description: "Metric type classification (e.g., 'Unknown', 'CPU Utilization').", Optional: true},
									"fill_zero":                      schema.BoolAttribute{Description: "Fill missing data with zero.", Optional: true, Computed: true},
									"enable_baseline_near_constance": schema.BoolAttribute{Description: "Enable baseline near constance detection.", Optional: true},
									"compute_difference":             schema.BoolAttribute{Description: "Compute difference for this metric.", Optional: true},
									"anomaly_gap_tolerance_duration": schema.Int64Attribute{Description: "Anomaly gap tolerance duration in milliseconds.", Optional: true, Computed: true},
									"rouge_value": schema.StringAttribute{
										Description: "Rouge value as raw JSON string from API (e.g., '{\"l\":NaN,\"s\":NaN}'). Null to disable.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
											ignoreStringDrift(),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *metricProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

// populateMetricSettings converts the Terraform plan/state into a settings map for API calls.
func populateMetricSettings(plan *metricProjectResourceModel) map[string]interface{} {
	parseJSONField := func(jsonStr string) interface{} {
		if jsonStr == "" {
			return nil
		}
		var result interface{}
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			return jsonStr
		}
		return result
	}

	// Build the map directly to preserve false/zero values that omitempty would strip.
	s := make(map[string]interface{})

	if !plan.ProjectName.IsNull() {
		s["projectName"] = plan.ProjectName.ValueString()
	}
	if !plan.ProjectDisplayName.IsNull() {
		s["projectDisplayName"] = plan.ProjectDisplayName.ValueString()
	}
	if !plan.CValue.IsNull() {
		s["cValue"] = int(plan.CValue.ValueInt64())
	}
	if !plan.PValue.IsNull() {
		s["pValue"] = plan.PValue.ValueFloat64()
	}
	if !plan.ProjectTimeZone.IsNull() {
		s["projectTimeZone"] = plan.ProjectTimeZone.ValueString()
	}
	if !plan.SamplingInterval.IsNull() {
		s["samplingInterval"] = int(plan.SamplingInterval.ValueInt64())
	}
	if !plan.RetentionTime.IsNull() {
		s["retentionTime"] = int(plan.RetentionTime.ValueInt64())
	}
	if !plan.UBLRetentionTime.IsNull() {
		s["UBLRetentionTime"] = int(plan.UBLRetentionTime.ValueInt64())
	}
	if !plan.TrainingFilter.IsNull() {
		s["trainingFilter"] = plan.TrainingFilter.ValueBool()
	}
	if !plan.EnableNewAlertEmail.IsNull() {
		s["enableNewAlertEmail"] = plan.EnableNewAlertEmail.ValueBool()
	}
	if !plan.LargeProject.IsNull() {
		s["largeProject"] = plan.LargeProject.ValueBool()
	}
	if !plan.NewPatternRange.IsNull() {
		s["newPatternRange"] = int(plan.NewPatternRange.ValueInt64())
	}
	if !plan.Proxy.IsNull() {
		s["proxy"] = plan.Proxy.ValueString()
	}
	if !plan.EnableAnomalyScoreEscalation.IsNull() {
		s["enableAnomalyScoreEscalation"] = plan.EnableAnomalyScoreEscalation.ValueBool()
	}
	if !plan.EscalationAnomalyScoreThreshold.IsNull() {
		s["escalationAnomalyScoreThreshold"] = plan.EscalationAnomalyScoreThreshold.ValueString()
	}
	if !plan.IgnoreAnomalyScoreThreshold.IsNull() {
		s["ignoreAnomalyScoreThreshold"] = plan.IgnoreAnomalyScoreThreshold.ValueString()
	}
	if !plan.EnableStreamDetection.IsNull() {
		s["enableStreamDetection"] = plan.EnableStreamDetection.ValueBool()
	}
	if !plan.IgnoreInstanceForKB.IsNull() {
		s["ignoreInstanceForKB"] = plan.IgnoreInstanceForKB.ValueBool()
	}
	if !plan.ShowInstanceDown.IsNull() {
		s["showInstanceDown"] = plan.ShowInstanceDown.ValueBool()
	}
	if !plan.InstanceDownEnable.IsNull() {
		s["instanceDownEnable"] = plan.InstanceDownEnable.ValueBool()
	}
	if !plan.AlertHourlyCost.IsNull() {
		s["alertHourlyCost"] = plan.AlertHourlyCost.ValueFloat64()
	}
	if !plan.AlertAverageTime.IsNull() {
		s["alertAverageTime"] = int(plan.AlertAverageTime.ValueInt64())
	}
	if !plan.AvgPerIncidentDowntimeCost.IsNull() {
		s["avgPerIncidentDowntimeCost"] = plan.AvgPerIncidentDowntimeCost.ValueFloat64()
	}

	// Incident prediction and RCA
	if !plan.IncidentPredictionWindow.IsNull() {
		s["incidentPredictionWindow"] = int(plan.IncidentPredictionWindow.ValueInt64())
	}
	if !plan.MinIncidentPredictionWindow.IsNull() {
		s["minIncidentPredictionWindow"] = int(plan.MinIncidentPredictionWindow.ValueInt64())
	}
	if !plan.IncidentRelationSearchWindow.IsNull() {
		s["incidentRelationSearchWindow"] = int(plan.IncidentRelationSearchWindow.ValueInt64())
	}
	if !plan.IncidentPredictionEventLimit.IsNull() {
		s["incidentPredictionEventLimit"] = int(plan.IncidentPredictionEventLimit.ValueInt64())
	}
	if !plan.RootCauseCountThreshold.IsNull() {
		s["rootCauseCountThreshold"] = int(plan.RootCauseCountThreshold.ValueInt64())
	}
	if !plan.RootCauseProbabilityThreshold.IsNull() {
		s["rootCauseProbabilityThreshold"] = plan.RootCauseProbabilityThreshold.ValueFloat64()
	}
	if !plan.CompositeRCALimit.IsNull() {
		s["compositeRCALimit"] = int(plan.CompositeRCALimit.ValueInt64())
	}
	if !plan.RootCauseLogMessageSearchRange.IsNull() {
		s["rootCauseLogMessageSearchRange"] = int(plan.RootCauseLogMessageSearchRange.ValueInt64())
	}
	if !plan.CausalPredictionSetting.IsNull() {
		s["causalPredictionSetting"] = int(plan.CausalPredictionSetting.ValueInt64())
	}
	if !plan.RootCauseRankSetting.IsNull() {
		s["rootCauseRankSetting"] = int(plan.RootCauseRankSetting.ValueInt64())
	}
	if !plan.MaximumRootCauseResultSize.IsNull() {
		s["maximumRootCauseResultSize"] = int(plan.MaximumRootCauseResultSize.ValueInt64())
	}
	if !plan.MultiHopSearchLevel.IsNull() {
		s["multiHopSearchLevel"] = int(plan.MultiHopSearchLevel.ValueInt64())
	}
	if !plan.MultiHopSearchLimit.IsNull() {
		s["multiHopSearchLimit"] = plan.MultiHopSearchLimit.ValueString()
	}
	if !plan.PredictionCountThreshold.IsNull() {
		s["predictionCountThreshold"] = int(plan.PredictionCountThreshold.ValueInt64())
	}
	if !plan.PredictionProbabilityThreshold.IsNull() {
		s["predictionProbabilityThreshold"] = plan.PredictionProbabilityThreshold.ValueFloat64()
	}
	if !plan.PredictionRuleActiveCondition.IsNull() {
		s["predictionRuleActiveCondition"] = int(plan.PredictionRuleActiveCondition.ValueInt64())
	}
	if !plan.PredictionRuleActiveThreshold.IsNull() {
		s["predictionRuleActiveThreshold"] = plan.PredictionRuleActiveThreshold.ValueFloat64()
	}
	if !plan.PredictionRuleFalsePositiveThreshold.IsNull() {
		s["predictionRuleFalsePositiveThreshold"] = int(plan.PredictionRuleFalsePositiveThreshold.ValueInt64())
	}
	if !plan.PredictionRuleInactiveThreshold.IsNull() {
		s["predictionRuleInactiveThreshold"] = plan.PredictionRuleInactiveThreshold.ValueFloat64()
	}
	if !plan.MinValidModelSpan.IsNull() {
		s["minValidModelSpan"] = int(plan.MinValidModelSpan.ValueInt64())
	}

	// Webhook
	if !plan.MaxWebHookRequestSize.IsNull() {
		s["maxWebHookRequestSize"] = int(plan.MaxWebHookRequestSize.ValueInt64())
	}
	if !plan.WebhookAlertDampening.IsNull() {
		s["webhookAlertDampening"] = int(plan.WebhookAlertDampening.ValueInt64())
	}
	if !plan.WebhookBlackListSetStr.IsNull() {
		s["webhookBlackListSetStr"] = plan.WebhookBlackListSetStr.ValueString()
	}
	if !plan.WebhookCriticalKeywordSetStr.IsNull() {
		s["webhookCriticalKeywordSetStr"] = plan.WebhookCriticalKeywordSetStr.ValueString()
	}
	if !plan.WebhookTypeSetStr.IsNull() {
		s["webhookTypeSetStr"] = plan.WebhookTypeSetStr.ValueString()
	}
	if !plan.WebhookUrl.IsNull() {
		s["webhookUrl"] = plan.WebhookUrl.ValueString()
	}

	// Metric-specific
	if !plan.HighRatioCValue.IsNull() {
		s["highRatioCValue"] = int(plan.HighRatioCValue.ValueInt64())
	}
	if !plan.MaximumHint.IsNull() {
		s["maximumHint"] = int(plan.MaximumHint.ValueInt64())
	}
	if !plan.DynamicBaselineDetectionFlag.IsNull() {
		s["dynamicBaselineDetectionFlag"] = plan.DynamicBaselineDetectionFlag.ValueBool()
	}
	if !plan.PositiveBaselineViolationFactor.IsNull() {
		s["positiveBaselineViolationFactor"] = plan.PositiveBaselineViolationFactor.ValueFloat64()
	}
	if !plan.NegativeBaselineViolationFactor.IsNull() {
		s["negativeBaselineViolationFactor"] = plan.NegativeBaselineViolationFactor.ValueFloat64()
	}
	if !plan.EnablePeriodAnomalyFilter.IsNull() {
		s["enablePeriodAnomalyFilter"] = plan.EnablePeriodAnomalyFilter.ValueBool()
	}
	if !plan.EnableUBLDetect.IsNull() {
		s["enableUBLDetect"] = plan.EnableUBLDetect.ValueBool()
	}
	if !plan.EnableCumulativeDetect.IsNull() {
		s["enableCumulativeDetect"] = plan.EnableCumulativeDetect.ValueBool()
	}
	if !plan.EnableComponentLevelDetection.IsNull() {
		s["enableComponentLevelDetection"] = plan.EnableComponentLevelDetection.ValueBool()
	}
	if !plan.PredictionTrainingDataLength.IsNull() {
		s["predictionTrainingDataLength"] = int(plan.PredictionTrainingDataLength.ValueInt64())
	}
	if !plan.PredictionCorrelationSensitivity.IsNull() {
		s["predictionCorrelationSensitivity"] = plan.PredictionCorrelationSensitivity.ValueFloat64()
	}
	if !plan.EnableKPIPrediction.IsNull() {
		s["enableKPIPrediction"] = plan.EnableKPIPrediction.ValueBool()
	}
	if !plan.InstanceDownThreshold.IsNull() {
		s["instanceDownThreshold"] = int(plan.InstanceDownThreshold.ValueInt64())
	}
	if !plan.InstanceDownReportNumber.IsNull() {
		s["instanceDownReportNumber"] = int(plan.InstanceDownReportNumber.ValueInt64())
	}
	if !plan.ModelSpan.IsNull() {
		s["modelSpan"] = int(plan.ModelSpan.ValueInt64())
	}
	if !plan.EnableMetricDataPrediction.IsNull() {
		s["enableMetricDataPrediction"] = plan.EnableMetricDataPrediction.ValueBool()
	}
	if !plan.EnableBaselineDetectionDoubleVerify.IsNull() {
		s["enableBaselineDetectionDoubleVerify"] = plan.EnableBaselineDetectionDoubleVerify.ValueBool()
	}
	if !plan.EnableFillGap.IsNull() {
		s["enableFillGap"] = plan.EnableFillGap.ValueBool()
	}
	if !plan.EnableStoreFilledGap.IsNull() {
		s["enableStoreFilledGap"] = plan.EnableStoreFilledGap.ValueBool()
	}
	if !plan.GapFillingTrainingDataLength.IsNull() {
		s["gapFillingTrainingDataLength"] = int(plan.GapFillingTrainingDataLength.ValueInt64())
	}
	if !plan.PatternIdGenerationRule.IsNull() {
		s["patternIdGenerationRule"] = int(plan.PatternIdGenerationRule.ValueInt64())
	}
	if !plan.AnomalyGapToleranceCount.IsNull() {
		s["anomalyGapToleranceCount"] = int(plan.AnomalyGapToleranceCount.ValueInt64())
	}
	if !plan.FilterByAnomalyInBaselineGeneration.IsNull() {
		s["filterByAnomalyInBaselineGeneration"] = plan.FilterByAnomalyInBaselineGeneration.ValueBool()
	}
	if !plan.BaselineDuration.IsNull() {
		s["baselineDuration"] = int(plan.BaselineDuration.ValueInt64())
	}
	if !plan.AnomalyDampening.IsNull() {
		s["anomalyDampening"] = int(plan.AnomalyDampening.ValueInt64())
	}
	if !plan.InstanceDownRatioThreshold.IsNull() {
		s["instanceDownRatioThreshold"] = plan.InstanceDownRatioThreshold.ValueFloat64()
	}
	if !plan.ComponentNameAutoOverwrite.IsNull() {
		s["componentNameAutoOverwrite"] = plan.ComponentNameAutoOverwrite.ValueBool()
	}
	if !plan.Mode.IsNull() {
		s["processMode"] = int(plan.Mode.ValueInt64())
	}

	// Complex JSON fields
	if !plan.LinkedLogProjects.IsNull() {
		if parsed := parseJSONField(plan.LinkedLogProjects.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["linkedLogProjects"] = list
			}
		}
	}
	if !plan.ComponentMetricSettingOverallModelList.IsNull() {
		if parsed := parseJSONField(plan.ComponentMetricSettingOverallModelList.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["componentMetricSettingOverallModelList"] = list
			}
		}
	}
	if !plan.EmailSetting.IsNull() {
		if parsed := parseJSONField(plan.EmailSetting.ValueString()); parsed != nil {
			s["emailSetting"] = parsed
		}
	}
	if !plan.InstanceGroupingUpdate.IsNull() {
		if parsed := parseJSONField(plan.InstanceGroupingUpdate.ValueString()); parsed != nil {
			s["instanceGroupingUpdate"] = parsed
		}
	}
	if !plan.SharedUsernames.IsNull() {
		if parsed := parseJSONField(plan.SharedUsernames.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["sharedUsernames"] = list
			}
		}
	}
	if !plan.WebhookHeaderList.IsNull() {
		if parsed := parseJSONField(plan.WebhookHeaderList.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["webhookHeaderList"] = list
			}
		}
	}

	return s
}

// populateMetricStateFromSettings fills a metricProjectResourceModel from an API settings map.
func populateMetricStateFromSettings(m *metricProjectResourceModel, settings map[string]interface{}) {
	getInt64 := func(key string) types.Int64 {
		if val, ok := settings[key]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				return types.Int64Value(int64(v))
			case int64:
				return types.Int64Value(v)
			case int:
				return types.Int64Value(int64(v))
			}
		}
		return types.Int64Null()
	}

	getFloat64 := func(key string) types.Float64 {
		if val, ok := settings[key]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				return types.Float64Value(v)
			case int64:
				return types.Float64Value(float64(v))
			case int:
				return types.Float64Value(float64(v))
			}
		}
		return types.Float64Null()
	}

	getString := func(key string) types.String {
		if val, ok := settings[key]; ok && val != nil {
			if str, ok := val.(string); ok {
				return types.StringValue(str)
			}
		}
		return types.StringNull()
	}

	getBool := func(key string) types.Bool {
		if val, ok := settings[key]; ok && val != nil {
			if b, ok := val.(bool); ok {
				return types.BoolValue(b)
			}
		}
		return types.BoolNull()
	}

	getJSONString := func(key string) types.String {
		if val, ok := settings[key]; ok && val != nil {
			if str, ok := val.(string); ok {
				return types.StringValue(str)
			}
			if jsonBytes, err := json.Marshal(val); err == nil {
				return types.StringValue(string(jsonBytes))
			}
		}
		return types.StringNull()
	}

	m.ProjectDisplayName = getString("projectDisplayName")
	m.CValue = getInt64("cValue")
	m.PValue = getFloat64("pValue")
	m.ProjectTimeZone = getString("projectTimeZone")
	m.SamplingInterval = getInt64("samplingInterval")
	m.RetentionTime = getInt64("retentionTime")
	m.UBLRetentionTime = getInt64("UBLRetentionTime")
	m.TrainingFilter = getBool("trainingFilter")
	m.EnableNewAlertEmail = getBool("enableNewAlertEmail")
	m.LargeProject = getBool("largeProject")
	m.NewPatternRange = getInt64("newPatternRange")
	m.Proxy = getString("proxy")
	m.EnableAnomalyScoreEscalation = getBool("enableAnomalyScoreEscalation")
	m.EscalationAnomalyScoreThreshold = getString("escalationAnomalyScoreThreshold")
	m.IgnoreAnomalyScoreThreshold = getString("ignoreAnomalyScoreThreshold")
	m.EnableStreamDetection = getBool("enableStreamDetection")
	m.IgnoreInstanceForKB = getBool("ignoreInstanceForKB")
	m.ShowInstanceDown = getBool("showInstanceDown")
	m.InstanceDownEnable = getBool("instanceDownEnable")
	m.AlertHourlyCost = getFloat64("alertHourlyCost")
	m.AlertAverageTime = getInt64("alertAverageTime")
	m.AvgPerIncidentDowntimeCost = getFloat64("avgPerIncidentDowntimeCost")

	m.IncidentPredictionWindow = getInt64("incidentPredictionWindow")
	m.MinIncidentPredictionWindow = getInt64("minIncidentPredictionWindow")
	m.IncidentRelationSearchWindow = getInt64("incidentRelationSearchWindow")
	m.IncidentPredictionEventLimit = getInt64("incidentPredictionEventLimit")
	m.RootCauseCountThreshold = getInt64("rootCauseCountThreshold")
	m.RootCauseProbabilityThreshold = getFloat64("rootCauseProbabilityThreshold")
	m.CompositeRCALimit = getInt64("compositeRCALimit")
	m.RootCauseLogMessageSearchRange = getInt64("rootCauseLogMessageSearchRange")
	m.CausalPredictionSetting = getInt64("causalPredictionSetting")
	m.RootCauseRankSetting = getInt64("rootCauseRankSetting")
	m.MaximumRootCauseResultSize = getInt64("maximumRootCauseResultSize")
	m.MultiHopSearchLevel = getInt64("multiHopSearchLevel")
	m.MultiHopSearchLimit = getString("multiHopSearchLimit")
	m.PredictionCountThreshold = getInt64("predictionCountThreshold")
	m.PredictionProbabilityThreshold = getFloat64("predictionProbabilityThreshold")
	m.PredictionRuleActiveCondition = getInt64("predictionRuleActiveCondition")
	m.PredictionRuleActiveThreshold = getFloat64("predictionRuleActiveThreshold")
	m.PredictionRuleFalsePositiveThreshold = getInt64("predictionRuleFalsePositiveThreshold")
	m.PredictionRuleInactiveThreshold = getFloat64("predictionRuleInactiveThreshold")
	m.MinValidModelSpan = getInt64("minValidModelSpan")

	m.MaxWebHookRequestSize = getInt64("maxWebHookRequestSize")
	m.WebhookAlertDampening = getInt64("webhookAlertDampening")
	m.WebhookBlackListSetStr = getJSONString("webhookBlackListSetStr")
	m.WebhookCriticalKeywordSetStr = getJSONString("webhookCriticalKeywordSetStr")
	m.WebhookTypeSetStr = getJSONString("webhookTypeSetStr")
	m.WebhookUrl = getString("webhookUrl")

	// Metric-specific
	m.HighRatioCValue = getInt64("highRatioCValue")
	m.MaximumHint = getInt64("maximumHint")
	m.DynamicBaselineDetectionFlag = getBool("dynamicBaselineDetectionFlag")
	m.PositiveBaselineViolationFactor = getFloat64("positiveBaselineViolationFactor")
	m.NegativeBaselineViolationFactor = getFloat64("negativeBaselineViolationFactor")
	m.EnablePeriodAnomalyFilter = getBool("enablePeriodAnomalyFilter")
	m.EnableUBLDetect = getBool("enableUBLDetect")
	m.EnableCumulativeDetect = getBool("enableCumulativeDetect")
	m.EnableComponentLevelDetection = getBool("enableComponentLevelDetection")
	m.PredictionTrainingDataLength = getInt64("predictionTrainingDataLength")
	m.PredictionCorrelationSensitivity = getFloat64("predictionCorrelationSensitivity")
	m.EnableKPIPrediction = getBool("enableKPIPrediction")
	m.InstanceDownThreshold = getInt64("instanceDownThreshold")
	m.InstanceDownReportNumber = getInt64("instanceDownReportNumber")
	m.ModelSpan = getInt64("modelSpan")
	m.EnableMetricDataPrediction = getBool("enableMetricDataPrediction")
	m.EnableBaselineDetectionDoubleVerify = getBool("enableBaselineDetectionDoubleVerify")
	m.EnableFillGap = getBool("enableFillGap")
	m.EnableStoreFilledGap = getBool("enableStoreFilledGap")
	m.GapFillingTrainingDataLength = getInt64("gapFillingTrainingDataLength")
	m.PatternIdGenerationRule = getInt64("patternIdGenerationRule")
	m.AnomalyGapToleranceCount = getInt64("anomalyGapToleranceCount")
	m.FilterByAnomalyInBaselineGeneration = getBool("filterByAnomalyInBaselineGeneration")
	m.BaselineDuration = getInt64("baselineDuration")
	m.AnomalyDampening = getInt64("anomalyDampening")
	m.InstanceDownRatioThreshold = getFloat64("instanceDownRatioThreshold")
	m.ComponentNameAutoOverwrite = getBool("componentNameAutoOverwrite")
	m.Mode = getInt64("processMode")

	// Complex JSON fields
	m.LinkedLogProjects = getJSONString("linkedLogProjects")
	m.ComponentMetricSettingOverallModelList = getJSONString("componentMetricSettingOverallModelList")
	m.InstanceGroupingUpdate = getJSONString("instanceGroupingUpdate")
	m.SharedUsernames = getJSONString("sharedUsernames")
	m.WebhookHeaderList = getJSONString("webhookHeaderList")
}

// preserveMetricConfigValues overwrites state fields with explicit config values to prevent drift.
func preserveMetricConfigValues(plan *metricProjectResourceModel, config *metricProjectResourceModel) {
	preserve := func(planVal, configVal *types.String) {
		if !configVal.IsNull() {
			*planVal = *configVal
		}
	}
	preserveInt := func(planVal, configVal *types.Int64) {
		if !configVal.IsNull() {
			*planVal = *configVal
		}
	}
	preserveFloat := func(planVal, configVal *types.Float64) {
		if !configVal.IsNull() {
			*planVal = *configVal
		}
	}
	preserveBool := func(planVal, configVal *types.Bool) {
		if !configVal.IsNull() {
			*planVal = *configVal
		}
	}

	preserve(&plan.ProjectDisplayName, &config.ProjectDisplayName)
	preserveInt(&plan.CValue, &config.CValue)
	preserveFloat(&plan.PValue, &config.PValue)
	preserve(&plan.ProjectTimeZone, &config.ProjectTimeZone)
	preserveInt(&plan.SamplingInterval, &config.SamplingInterval)
	preserveInt(&plan.RetentionTime, &config.RetentionTime)
	preserveInt(&plan.UBLRetentionTime, &config.UBLRetentionTime)
	preserveBool(&plan.TrainingFilter, &config.TrainingFilter)
	preserveBool(&plan.EnableNewAlertEmail, &config.EnableNewAlertEmail)
	preserveBool(&plan.LargeProject, &config.LargeProject)
	preserveInt(&plan.NewPatternRange, &config.NewPatternRange)
	preserve(&plan.Proxy, &config.Proxy)
	preserveBool(&plan.EnableAnomalyScoreEscalation, &config.EnableAnomalyScoreEscalation)
	preserve(&plan.EscalationAnomalyScoreThreshold, &config.EscalationAnomalyScoreThreshold)
	preserve(&plan.IgnoreAnomalyScoreThreshold, &config.IgnoreAnomalyScoreThreshold)
	preserveBool(&plan.EnableStreamDetection, &config.EnableStreamDetection)
	preserveBool(&plan.IgnoreInstanceForKB, &config.IgnoreInstanceForKB)
	preserveBool(&plan.ShowInstanceDown, &config.ShowInstanceDown)
	preserveBool(&plan.InstanceDownEnable, &config.InstanceDownEnable)
	preserveFloat(&plan.AlertHourlyCost, &config.AlertHourlyCost)
	preserveInt(&plan.AlertAverageTime, &config.AlertAverageTime)
	preserveFloat(&plan.AvgPerIncidentDowntimeCost, &config.AvgPerIncidentDowntimeCost)
	preserveInt(&plan.IncidentPredictionWindow, &config.IncidentPredictionWindow)
	preserveInt(&plan.MinIncidentPredictionWindow, &config.MinIncidentPredictionWindow)
	preserveInt(&plan.IncidentRelationSearchWindow, &config.IncidentRelationSearchWindow)
	preserveInt(&plan.IncidentPredictionEventLimit, &config.IncidentPredictionEventLimit)
	preserveInt(&plan.RootCauseCountThreshold, &config.RootCauseCountThreshold)
	preserveFloat(&plan.RootCauseProbabilityThreshold, &config.RootCauseProbabilityThreshold)
	preserveInt(&plan.CompositeRCALimit, &config.CompositeRCALimit)
	preserveInt(&plan.RootCauseLogMessageSearchRange, &config.RootCauseLogMessageSearchRange)
	preserveInt(&plan.CausalPredictionSetting, &config.CausalPredictionSetting)
	preserveInt(&plan.RootCauseRankSetting, &config.RootCauseRankSetting)
	preserveInt(&plan.MaximumRootCauseResultSize, &config.MaximumRootCauseResultSize)
	preserveInt(&plan.MultiHopSearchLevel, &config.MultiHopSearchLevel)
	preserve(&plan.MultiHopSearchLimit, &config.MultiHopSearchLimit)
	preserveInt(&plan.PredictionCountThreshold, &config.PredictionCountThreshold)
	preserveFloat(&plan.PredictionProbabilityThreshold, &config.PredictionProbabilityThreshold)
	if !config.PredictionRuleActiveCondition.IsNull() {
		plan.PredictionRuleActiveCondition = config.PredictionRuleActiveCondition
	}
	if !config.PredictionRuleActiveThreshold.IsNull() {
		plan.PredictionRuleActiveThreshold = config.PredictionRuleActiveThreshold
	}
	if !config.PredictionRuleFalsePositiveThreshold.IsNull() {
		plan.PredictionRuleFalsePositiveThreshold = config.PredictionRuleFalsePositiveThreshold
	}
	if !config.PredictionRuleInactiveThreshold.IsNull() {
		plan.PredictionRuleInactiveThreshold = config.PredictionRuleInactiveThreshold
	}
	preserveInt(&plan.MinValidModelSpan, &config.MinValidModelSpan)
	preserveInt(&plan.MaxWebHookRequestSize, &config.MaxWebHookRequestSize)
	preserveInt(&plan.WebhookAlertDampening, &config.WebhookAlertDampening)
	preserve(&plan.WebhookBlackListSetStr, &config.WebhookBlackListSetStr)
	preserve(&plan.WebhookCriticalKeywordSetStr, &config.WebhookCriticalKeywordSetStr)
	preserve(&plan.WebhookTypeSetStr, &config.WebhookTypeSetStr)
	preserve(&plan.WebhookUrl, &config.WebhookUrl)
	preserveInt(&plan.HighRatioCValue, &config.HighRatioCValue)
	preserveInt(&plan.MaximumHint, &config.MaximumHint)
	preserveBool(&plan.DynamicBaselineDetectionFlag, &config.DynamicBaselineDetectionFlag)
	preserveFloat(&plan.PositiveBaselineViolationFactor, &config.PositiveBaselineViolationFactor)
	preserveFloat(&plan.NegativeBaselineViolationFactor, &config.NegativeBaselineViolationFactor)
	preserveBool(&plan.EnablePeriodAnomalyFilter, &config.EnablePeriodAnomalyFilter)
	preserveBool(&plan.EnableUBLDetect, &config.EnableUBLDetect)
	preserveBool(&plan.EnableCumulativeDetect, &config.EnableCumulativeDetect)
	preserveBool(&plan.EnableComponentLevelDetection, &config.EnableComponentLevelDetection)
	preserveInt(&plan.PredictionTrainingDataLength, &config.PredictionTrainingDataLength)
	preserveFloat(&plan.PredictionCorrelationSensitivity, &config.PredictionCorrelationSensitivity)
	preserveBool(&plan.EnableKPIPrediction, &config.EnableKPIPrediction)
	preserveInt(&plan.InstanceDownThreshold, &config.InstanceDownThreshold)
	preserveInt(&plan.InstanceDownReportNumber, &config.InstanceDownReportNumber)
	preserveInt(&plan.ModelSpan, &config.ModelSpan)
	preserveBool(&plan.EnableMetricDataPrediction, &config.EnableMetricDataPrediction)
	preserveBool(&plan.EnableBaselineDetectionDoubleVerify, &config.EnableBaselineDetectionDoubleVerify)
	preserveBool(&plan.EnableFillGap, &config.EnableFillGap)
	preserveBool(&plan.EnableStoreFilledGap, &config.EnableStoreFilledGap)
	preserveInt(&plan.GapFillingTrainingDataLength, &config.GapFillingTrainingDataLength)
	preserveInt(&plan.PatternIdGenerationRule, &config.PatternIdGenerationRule)
	preserveInt(&plan.AnomalyGapToleranceCount, &config.AnomalyGapToleranceCount)
	preserveBool(&plan.FilterByAnomalyInBaselineGeneration, &config.FilterByAnomalyInBaselineGeneration)
	preserveInt(&plan.BaselineDuration, &config.BaselineDuration)
	preserveInt(&plan.AnomalyDampening, &config.AnomalyDampening)
	preserveFloat(&plan.InstanceDownRatioThreshold, &config.InstanceDownRatioThreshold)
	preserveBool(&plan.ComponentNameAutoOverwrite, &config.ComponentNameAutoOverwrite)
	preserve(&plan.LinkedLogProjects, &config.LinkedLogProjects)
	preserve(&plan.ComponentMetricSettingOverallModelList, &config.ComponentMetricSettingOverallModelList)
	// Preserve email_setting from config to avoid drift from API-enriched values
	plan.EmailSetting = config.EmailSetting
	preserve(&plan.InstanceGroupingUpdate, &config.InstanceGroupingUpdate)
	preserve(&plan.SharedUsernames, &config.SharedUsernames)
	preserve(&plan.WebhookHeaderList, &config.WebhookHeaderList)
	if !config.MetricConfigurations.IsNull() {
		plan.MetricConfigurations = config.MetricConfigurations
	}
	preserveInt(&plan.Mode, &config.Mode)
}

func (r *metricProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan metricProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating metric project", map[string]any{"project_name": plan.ProjectName.ValueString()})

	projectConfig := &client.ProjectConfig{
		ProjectName:        plan.ProjectName.ValueString(),
		ProjectDisplayName: plan.ProjectDisplayName.ValueString(),
		SystemName:         plan.SystemName.ValueString(),
		DataType:           plan.ProjectCreationConfig.DataType.ValueString(),
		InstanceType:       plan.ProjectCreationConfig.InstanceType.ValueString(),
		ProjectCloudType:   plan.ProjectCreationConfig.ProjectCloudType.ValueString(),
		InsightAgentType:   plan.ProjectCreationConfig.InsightAgentType.ValueString(),
		CValue:             int(plan.CValue.ValueInt64()),
		PValue:             plan.PValue.ValueFloat64(),
	}

	if err := r.client.CreateProject(projectConfig); err != nil {
		resp.Diagnostics.AddError("Error creating metric project",
			"Could not create metric project, unexpected error: "+err.Error())
		return
	}

	plan.ID = plan.ProjectName

	settings := populateMetricSettings(&plan)
	if len(settings) > 0 {
		updateConfig := &client.ProjectConfig{
			ProjectName: plan.ProjectName.ValueString(),
			Settings:    settings,
		}
		if err := r.client.UpdateMetricProject(updateConfig); err != nil {
			tflog.Warn(ctx, "Could not apply all settings on creation", map[string]any{
				"error": err.Error(),
				"note":  "Settings can be applied on next terraform apply",
			})
		}
	}

	var config metricProjectResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(plan.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		tflog.Warn(ctx, "Could not read metric project after creation", map[string]any{"error": err.Error()})
		config.ID = plan.ProjectName
		resp.State.Set(ctx, config)
		return
	}

	if project != nil {
		plan = config
		plan.ID = plan.ProjectName

		apiSettings := project.Settings
		if apiSettings == nil {
			apiSettings = make(map[string]interface{})
		}
		populateMetricStateFromSettings(&plan, apiSettings)
	}

	preserveMetricConfigValues(&plan, &config)

	// Process holiday_settings
	if !config.HolidaySettings.IsNull() && !config.HolidaySettings.IsUnknown() {
		var holidays []holidaySettingModel
		diags = config.HolidaySettings.ElementsAs(ctx, &holidays, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, holiday := range holidays {
			clientHoliday := &client.Holiday{
				Name:      holiday.Name.ValueString(),
				StartDate: holiday.StartDate.ValueString(),
				EndDate:   holiday.EndDate.ValueString(),
			}
			if err := r.client.CreateHoliday(plan.ProjectName.ValueString(), clientHoliday); err != nil {
				resp.Diagnostics.AddError("Error creating holiday",
					fmt.Sprintf("Could not create holiday '%s': %s", clientHoliday.Name, err.Error()))
				return
			}
		}

		plan.HolidaySettings = config.HolidaySettings
	} else {
		plan.HolidaySettings = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":       types.StringType,
				"start_date": types.StringType,
				"end_date":   types.StringType,
			},
		})
	}

	// Process metric_configurations
	if !config.MetricConfigurations.IsNull() && !config.MetricConfigurations.IsUnknown() {
		metricConfigs, mcDiags := metricConfigsFromMap(ctx, config.MetricConfigurations)
		resp.Diagnostics.Append(mcDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		patternIdRule := int(plan.PatternIdGenerationRule.ValueInt64())
		applyDiags := applyMetricConfigurations(ctx, r.client, plan.ProjectName.ValueString(), metricConfigs, patternIdRule, plan.SamplingInterval.ValueInt64())
		resp.Diagnostics.Append(applyDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Normalize null component lists → [] so state is consistent with what Read returns,
		// preventing a perpetual plan diff.
		normalizedConfigs, normDiags := normalizeMetricConfigsForState(ctx, metricConfigs)
		resp.Diagnostics.Append(normDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		setVal, d := metricConfigsToMap(ctx, normalizedConfigs)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.MetricConfigurations = setVal
	} else {
		plan.MetricConfigurations = types.MapNull(types.ObjectType{AttrTypes: metricConfigurationValueAttrTypes()})
	}

	plan.SystemName = config.SystemName
	plan.ProjectCreationConfig = config.ProjectCreationConfig

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Metric project created successfully", map[string]any{"project_name": plan.ProjectName.ValueString()})
}

func (r *metricProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state metricProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading metric project", map[string]any{"project_name": state.ProjectName.ValueString()})

	project, err := r.client.GetProject(state.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError("Error reading metric project",
			"Could not read metric project, unexpected error: "+err.Error())
		return
	}

	if project == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	settings := project.Settings
	if settings == nil {
		settings = make(map[string]interface{})
	}
	populateMetricStateFromSettings(&state, settings)

	// Read holiday settings from API
	holidays, err := r.client.GetHolidays(state.ProjectName.ValueString())
	if err != nil {
		tflog.Warn(ctx, "Could not read holidays", map[string]any{"error": err.Error()})
	} else if holidays != nil {
		var existingHolidays []holidaySettingModel
		if !state.HolidaySettings.IsNull() && !state.HolidaySettings.IsUnknown() {
			state.HolidaySettings.ElementsAs(ctx, &existingHolidays, false)
		}

		apiHolidayMap := make(map[string]holidaySettingModel)
		for name, dates := range holidays {
			var startDate, endDate string
			if dates != "" {
				parts := splitByComma(dates)
				if len(parts) >= 2 {
					startDate = parts[0]
					endDate = parts[1]
				} else if len(parts) == 1 {
					startDate = parts[0]
					endDate = parts[0]
				}
			}
			apiHolidayMap[name] = holidaySettingModel{
				Name:      types.StringValue(name),
				StartDate: types.StringValue(startDate),
				EndDate:   types.StringValue(endDate),
			}
		}

		var holidaySettings []holidaySettingModel
		seenNames := make(map[string]bool)
		for _, existing := range existingHolidays {
			name := existing.Name.ValueString()
			if apiHoliday, exists := apiHolidayMap[name]; exists {
				holidaySettings = append(holidaySettings, apiHoliday)
				seenNames[name] = true
			}
		}
		var newHolidays []holidaySettingModel
		for name, holiday := range apiHolidayMap {
			if !seenNames[name] {
				newHolidays = append(newHolidays, holiday)
			}
		}
		sort.Slice(newHolidays, func(i, j int) bool {
			return newHolidays[i].Name.ValueString() < newHolidays[j].Name.ValueString()
		})
		holidaySettings = append(holidaySettings, newHolidays...)

		if len(holidaySettings) > 0 {
			listValue, diags := types.ListValueFrom(ctx, types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":       types.StringType,
					"start_date": types.StringType,
					"end_date":   types.StringType,
				},
			}, holidaySettings)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				state.HolidaySettings = listValue
			}
		} else {
			state.HolidaySettings = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":       types.StringType,
					"start_date": types.StringType,
					"end_date":   types.StringType,
				},
			})
		}
	}

	// Read metric_configurations from API (only for metrics already tracked in state)
	if !state.MetricConfigurations.IsNull() && !state.MetricConfigurations.IsUnknown() {
		existingConfigs, mcDiags := metricConfigsFromMap(ctx, state.MetricConfigurations)
		resp.Diagnostics.Append(mcDiags...)
		if !resp.Diagnostics.HasError() {
			updatedConfigs, readDiags := readMetricConfigurationsFromAPI(ctx, r.client, state.ProjectName.ValueString(), existingConfigs)
			resp.Diagnostics.Append(readDiags...)
			if !resp.Diagnostics.HasError() && len(updatedConfigs) > 0 {
				setVal, d := metricConfigsToMap(ctx, updatedConfigs)
				resp.Diagnostics.Append(d...)
				if !resp.Diagnostics.HasError() {
					state.MetricConfigurations = setVal
				}
			}
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Metric project read successfully", map[string]any{"project_name": state.ProjectName.ValueString()})
}

func (r *metricProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan metricProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config metricProjectResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating metric project", map[string]any{"project_name": config.ProjectName.ValueString()})

	projectConfig := &client.ProjectConfig{
		ProjectName:        config.ProjectName.ValueString(),
		ProjectDisplayName: config.ProjectDisplayName.ValueString(),
		SystemName:         config.SystemName.ValueString(),
		CValue:             int(config.CValue.ValueInt64()),
		PValue:             config.PValue.ValueFloat64(),
		Settings:           populateMetricSettings(&config),
	}

	if err := r.client.UpdateMetricProject(projectConfig); err != nil {
		resp.Diagnostics.AddError("Error updating metric project",
			"Could not update metric project, unexpected error: "+err.Error())
		return
	}

	project, err := r.client.GetProject(plan.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		resp.Diagnostics.AddWarning("Error reading metric project after update",
			"Could not read metric project after update: "+err.Error()+". State may be out of sync.")
	} else if project != nil {
		apiSettings := project.Settings
		if apiSettings == nil {
			apiSettings = make(map[string]interface{})
		}
		populateMetricStateFromSettings(&plan, apiSettings)
	}

	preserveMetricConfigValues(&plan, &config)

	// Process holiday_settings
	if !config.HolidaySettings.IsNull() && !config.HolidaySettings.IsUnknown() {
		var configHolidays []holidaySettingModel
		diags = config.HolidaySettings.ElementsAs(ctx, &configHolidays, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		currentHolidays, err := r.client.GetHolidays(plan.ProjectName.ValueString())
		if err != nil {
			tflog.Warn(ctx, "Could not get current holidays", map[string]any{"error": err.Error()})
			currentHolidays = make(map[string]string)
		}

		configHolidayNames := make(map[string]bool)
		for _, h := range configHolidays {
			configHolidayNames[h.Name.ValueString()] = true
		}

		var holidaysToDelete []string
		for name := range currentHolidays {
			if !configHolidayNames[name] {
				holidaysToDelete = append(holidaysToDelete, name)
			}
		}
		if len(holidaysToDelete) > 0 {
			if err := r.client.DeleteHolidays(plan.ProjectName.ValueString(), holidaysToDelete); err != nil {
				resp.Diagnostics.AddError("Error deleting holidays",
					fmt.Sprintf("Could not delete holidays: %s", err.Error()))
				return
			}
		}

		for _, holiday := range configHolidays {
			clientHoliday := &client.Holiday{
				Name:      holiday.Name.ValueString(),
				StartDate: holiday.StartDate.ValueString(),
				EndDate:   holiday.EndDate.ValueString(),
			}
			if existingDates, exists := currentHolidays[clientHoliday.Name]; exists {
				parts := splitByComma(existingDates)
				var existingStart, existingEnd string
				if len(parts) >= 2 {
					existingStart = parts[0]
					existingEnd = parts[1]
				} else if len(parts) == 1 {
					existingStart = parts[0]
					existingEnd = parts[0]
				}
				if existingStart == clientHoliday.StartDate && existingEnd == clientHoliday.EndDate {
					continue
				}
				if err := r.client.DeleteHolidays(plan.ProjectName.ValueString(), []string{clientHoliday.Name}); err != nil {
					resp.Diagnostics.AddError("Error updating holiday",
						fmt.Sprintf("Could not delete existing holiday '%s': %s", clientHoliday.Name, err.Error()))
					return
				}
			}
			if err := r.client.CreateHoliday(plan.ProjectName.ValueString(), clientHoliday); err != nil {
				resp.Diagnostics.AddError("Error creating holiday",
					fmt.Sprintf("Could not create holiday '%s': %s", clientHoliday.Name, err.Error()))
				return
			}
		}

		plan.HolidaySettings = config.HolidaySettings
	} else {
		currentHolidays, err := r.client.GetHolidays(plan.ProjectName.ValueString())
		if err != nil {
			tflog.Warn(ctx, "Could not get current holidays for deletion", map[string]any{"error": err.Error()})
		} else if len(currentHolidays) > 0 {
			var holidaysToDelete []string
			for name := range currentHolidays {
				holidaysToDelete = append(holidaysToDelete, name)
			}
			if err := r.client.DeleteHolidays(plan.ProjectName.ValueString(), holidaysToDelete); err != nil {
				resp.Diagnostics.AddError("Error deleting holidays",
					fmt.Sprintf("Could not delete holidays: %s", err.Error()))
				return
			}
		}
		plan.HolidaySettings = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":       types.StringType,
				"start_date": types.StringType,
				"end_date":   types.StringType,
			},
		})
	}

	// Process metric_configurations
	if !config.MetricConfigurations.IsNull() && !config.MetricConfigurations.IsUnknown() {
		metricConfigs, mcDiags := metricConfigsFromMap(ctx, config.MetricConfigurations)
		resp.Diagnostics.Append(mcDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		patternIdRule := int(plan.PatternIdGenerationRule.ValueInt64())
		applyDiags := applyMetricConfigurations(ctx, r.client, plan.ProjectName.ValueString(), metricConfigs, patternIdRule, plan.SamplingInterval.ValueInt64())
		resp.Diagnostics.Append(applyDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Normalize null component lists → [] so state is consistent with what Read returns,
		// preventing a perpetual plan diff.
		normalizedConfigs, normDiags := normalizeMetricConfigsForState(ctx, metricConfigs)
		resp.Diagnostics.Append(normDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		setVal, d := metricConfigsToMap(ctx, normalizedConfigs)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.MetricConfigurations = setVal
	} else {
		plan.MetricConfigurations = types.MapNull(types.ObjectType{AttrTypes: metricConfigurationValueAttrTypes()})
	}

	plan.SystemName = config.SystemName
	plan.ProjectCreationConfig = config.ProjectCreationConfig

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Metric project updated successfully", map[string]any{"project_name": plan.ProjectName.ValueString()})
}

func (r *metricProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state metricProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting metric project", map[string]any{"project_name": state.ProjectName.ValueString()})

	if err := r.client.DeleteProject(state.ProjectName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting metric project",
			"Could not delete metric project, unexpected error: "+err.Error())
		return
	}

	tflog.Info(ctx, "Metric project deleted successfully", map[string]any{"project_name": state.ProjectName.ValueString()})
}

func (r *metricProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildMetricAlertSettingFromAPI converts a client.MetricAlertSetting to a metricAlertSettingModel.
//
//nolint:unused
func buildMetricAlertSettingFromAPI(s client.MetricAlertSetting) metricAlertSettingModel {
	m := metricAlertSettingModel{
		ComponentName:                      types.StringValue(s.ComponentName),
		ThresholdAlertLowerBound:           types.StringValue(s.ThresholdAlertLowerBound),
		ThresholdAlertUpperBound:           types.StringValue(s.ThresholdAlertUpperBound),
		ThresholdAlertLowerBoundNegative:   types.StringValue(s.ThresholdAlertLowerBoundNegative),
		ThresholdAlertUpperBoundNegative:   types.StringValue(s.ThresholdAlertUpperBoundNegative),
		ThresholdNoAlertLowerBound:         types.StringValue(s.ThresholdNoAlertLowerBound),
		ThresholdNoAlertUpperBound:         types.StringValue(s.ThresholdNoAlertUpperBound),
		ThresholdNoAlertLowerBoundNegative: types.StringValue(s.ThresholdNoAlertLowerBoundNegative),
		ThresholdNoAlertUpperBoundNegative: types.StringValue(s.ThresholdNoAlertUpperBoundNegative),
		IncidentAlertLowerBound:            types.StringValue(s.IncidentAlertLowerBound),
		IncidentAlertUpperBound:            types.StringValue(s.IncidentAlertUpperBound),
		IncidentAlertLowerBoundNegative:    types.StringValue(s.IncidentAlertLowerBoundNegative),
		IncidentAlertUpperBoundNegative:    types.StringValue(s.IncidentAlertUpperBoundNegative),
		IncidentNoAlertLowerBound:          types.StringValue(s.IncidentNoAlertLowerBound),
		IncidentNoAlertUpperBound:          types.StringValue(s.IncidentNoAlertUpperBound),
		IncidentNoAlertLowerBoundNegative:  types.StringValue(s.IncidentNoAlertLowerBoundNegative),
		IncidentNoAlertUpperBoundNegative:  types.StringValue(s.IncidentNoAlertUpperBoundNegative),
		IsKPI:                              types.BoolValue(s.IsKPI),
		IsFlappingResultOnly:               types.BoolValue(s.IsFlappingResultOnly),
		IncidentDurationThreshold:          types.Int64Value(s.IncidentDurationThreshold),
		DetectionType:                      types.StringValue(s.DetectionType),
		PatternNameHigher:                  types.StringValue(s.PatternNameHigher),
		PatternNameLower:                   types.StringValue(s.PatternNameLower),
		MetricType:                         types.StringValue(s.MetricType),
		FillZero:                           types.BoolValue(s.FillZero),
		EnableBaselineNearConstance:        types.BoolValue(s.EnableBaselineNearConstance),
		ComputeDifference:                  types.BoolValue(s.ComputeDifference),
		AnomalyGapToleranceDuration:        types.Int64Value(s.AnomalyGapToleranceDuration),
	}
	if s.CValueOverride != nil {
		m.CValueOverride = types.Int64Value(*s.CValueOverride)
	} else {
		m.CValueOverride = types.Int64Null()
	}
	if s.HighCValueOverride != nil {
		m.HighCValueOverride = types.Int64Value(*s.HighCValueOverride)
	} else {
		m.HighCValueOverride = types.Int64Null()
	}
	// Treat both nil pointer and the literal string "null" (returned by the API) as Terraform null.
	if s.RougeValue == nil || *s.RougeValue == "null" {
		m.RougeValue = types.StringNull()
	} else {
		m.RougeValue = types.StringValue(*s.RougeValue)
	}
	return m
}

// buildSingleMetricAlertSettingPost converts a metricAlertSettingModel to a client.MetricAlertSettingPost.
// metricName is used for the SMetric field (not stored in the model itself).
// samplingIntervalS is the project sampling interval in seconds, used to convert the duration to a count.
func buildSingleMetricAlertSettingPost(metricName string, s metricAlertSettingModel, samplingIntervalS int64) client.MetricAlertSettingPost {
	detectionType := s.DetectionType.ValueString()
	if detectionType == "" {
		detectionType = "positive"
	}

	var anomalyCount int64 = 1
	if samplingIntervalS > 0 {
		durationMS := s.AnomalyGapToleranceDuration.ValueInt64()
		samplingMS := samplingIntervalS * 1000
		if samplingMS > 0 {
			anomalyCount = durationMS / samplingMS
		}
	}
	if anomalyCount < 1 {
		anomalyCount = 1
	}

	return client.MetricAlertSettingPost{
		SMetric:                            metricName,
		ComponentName:                      s.ComponentName.ValueString(),
		ThresholdAlertLowerBound:           s.ThresholdAlertLowerBound.ValueString(),
		ThresholdAlertUpperBound:           s.ThresholdAlertUpperBound.ValueString(),
		ThresholdAlertLowerBoundNegative:   s.ThresholdAlertLowerBoundNegative.ValueString(),
		ThresholdAlertUpperBoundNegative:   s.ThresholdAlertUpperBoundNegative.ValueString(),
		ThresholdNoAlertLowerBound:         s.ThresholdNoAlertLowerBound.ValueString(),
		ThresholdNoAlertUpperBound:         s.ThresholdNoAlertUpperBound.ValueString(),
		ThresholdNoAlertLowerBoundNegative: s.ThresholdNoAlertLowerBoundNegative.ValueString(),
		ThresholdNoAlertUpperBoundNegative: s.ThresholdNoAlertUpperBoundNegative.ValueString(),
		IncidentAlertLowerBound:            s.IncidentAlertLowerBound.ValueString(),
		IncidentAlertUpperBound:            s.IncidentAlertUpperBound.ValueString(),
		IncidentAlertLowerBoundNegative:    s.IncidentAlertLowerBoundNegative.ValueString(),
		IncidentAlertUpperBoundNegative:    s.IncidentAlertUpperBoundNegative.ValueString(),
		IncidentNoAlertLowerBound:          s.IncidentNoAlertLowerBound.ValueString(),
		IncidentNoAlertUpperBound:          s.IncidentNoAlertUpperBound.ValueString(),
		IncidentNoAlertLowerBoundNegative:  s.IncidentNoAlertLowerBoundNegative.ValueString(),
		IncidentNoAlertUpperBoundNegative:  s.IncidentNoAlertUpperBoundNegative.ValueString(),
		IsKPI:                              s.IsKPI.ValueBool(),
		IsFlappingResultOnly:               s.IsFlappingResultOnly.ValueBool(),
		IncidentDurationThreshold:          s.IncidentDurationThreshold.ValueInt64(),
		DetectionType:                      detectionType,
		PatternNameHigher:                  s.PatternNameHigher.ValueString(),
		PatternNameLower:                   s.PatternNameLower.ValueString(),
		MetricType:                         s.MetricType.ValueString(),
		FillZero:                           s.FillZero.ValueBool(),
		RougeValue:                         client.ConvertRougeValueForPost(s.RougeValue.ValueString()),
		EnableBaselineNearConstance:        s.EnableBaselineNearConstance.ValueBool(),
		ComputeDifference:                  s.ComputeDifference.ValueBool(),
		AnomalyGapToleranceDurationCount:   anomalyCount,
		CValueOverride:                     nullableInt64Div60(s.CValueOverride),
		HighCValueOverride:                 nullableInt64Div60(s.HighCValueOverride),
	}
}

// readMetricConfigurationsFromAPI reads current metric configuration state from the API
// for all metrics listed in existingConfigs.
func readMetricConfigurationsFromAPI(ctx context.Context, c *client.Client, projectName string, existingConfigs []metricConfigurationModel) ([]metricConfigurationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(existingConfigs) == 0 {
		return nil, diags
	}

	escalateMap, err := c.GetMetricComponents(projectName, "escalateIncident")
	if err != nil {
		diags.AddWarning("Could not read escalateIncident components", err.Error())
		escalateMap = make(map[string][]string)
	}
	ignoredMap, err := c.GetMetricComponents(projectName, "ignored")
	if err != nil {
		diags.AddWarning("Could not read ignored components", err.Error())
		ignoredMap = make(map[string][]string)
	}

	var result []metricConfigurationModel
	for _, existingCfg := range existingConfigs {
		metricName := existingCfg.MetricName.ValueString()

		// Build a lookup of existing alert settings by component name so we can
		// preserve null for c_value_override / high_c_value_override when the user
		// didn't specify them in config. Without this, a UI change that sets these
		// on the API side would cause every metric in the set to appear changed
		// (because the set element hash changes null → value for all entries).
		existingSettingByComponent := make(map[string]metricAlertSettingModel)
		if !existingCfg.MetricAlertSettings.IsNull() && !existingCfg.MetricAlertSettings.IsUnknown() {
			var existingSettings []metricAlertSettingModel
			if d := existingCfg.MetricAlertSettings.ElementsAs(ctx, &existingSettings, false); !d.HasError() {
				for _, s := range existingSettings {
					existingSettingByComponent[s.ComponentName.ValueString()] = s
				}
			}
		}

		// Use an explicit empty slice when the metric has no entries so the framework
		// stores [] (empty list) in state rather than null. This prevents a perpetual
		// diff where Terraform plans [] for a null Computed list attribute.
		escalateComps := escalateMap[metricName]
		if escalateComps == nil {
			escalateComps = []string{}
		}
		escalateListVal, d := types.ListValueFrom(ctx, types.StringType, escalateComps)
		diags.Append(d...)

		ignoredComps := ignoredMap[metricName]
		if ignoredComps == nil {
			ignoredComps = []string{}
		}
		ignoredListVal, d := types.ListValueFrom(ctx, types.StringType, ignoredComps)
		diags.Append(d...)

		settingEntry, err := c.GetMetricSettings(projectName, metricName)
		var alertSettingModels []metricAlertSettingModel
		if err != nil {
			diags.AddWarning(fmt.Sprintf("Could not read metric settings for %s", metricName), err.Error())
		} else if settingEntry != nil {
			global := buildMetricAlertSettingFromAPI(settingEntry.GlobalSetting)
			preserveConfiguredNulls(&global, existingSettingByComponent)
			alertSettingModels = append(alertSettingModels, global)
			for _, compSetting := range settingEntry.ComponentLevelSettingList {
				cm := buildMetricAlertSettingFromAPI(compSetting)
				preserveConfiguredNulls(&cm, existingSettingByComponent)
				alertSettingModels = append(alertSettingModels, cm)
			}
		}

		alertSettingsListVal, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: metricAlertSettingAttrTypes()}, alertSettingModels)
		diags.Append(d...)

		result = append(result, metricConfigurationModel{
			MetricName:                 types.StringValue(metricName),
			EscalateIncidentComponents: escalateListVal,
			IgnoredComponents:          ignoredListVal,
			MetricAlertSettings:        alertSettingsListVal,
		})
	}
	return result, diags
}

// preserveConfiguredNulls checks the existing state for a metric alert setting and, for
// fields where the prior state was null (user didn't configure them), resets the
// freshly-read model back to null. This prevents perpetual diffs when the API returns
// server-side defaults (e.g. detection_type="positive", rouge_value=NaN JSON) that were
// never explicitly set by the user.
func preserveConfiguredNulls(m *metricAlertSettingModel, existingByComponent map[string]metricAlertSettingModel) {
	old, ok := existingByComponent[m.ComponentName.ValueString()]
	if !ok {
		return
	}
	if old.CValueOverride.IsNull() {
		m.CValueOverride = types.Int64Null()
	}
	if old.HighCValueOverride.IsNull() {
		m.HighCValueOverride = types.Int64Null()
	}
	if old.DetectionType.IsNull() || old.DetectionType.ValueString() == "" {
		m.DetectionType = old.DetectionType
	}
	if old.RougeValue.IsNull() {
		m.RougeValue = types.StringNull()
	}
	// For fields that are Optional-only (not Computed), restore null when old state was null
	// so they don't produce spurious diffs from API-side defaults.
	if old.IsFlappingResultOnly.IsNull() {
		m.IsFlappingResultOnly = types.BoolNull()
	}
	if old.IncidentDurationThreshold.IsNull() {
		m.IncidentDurationThreshold = types.Int64Null()
	}
	if old.PatternNameHigher.IsNull() || old.PatternNameHigher.ValueString() == "" {
		m.PatternNameHigher = old.PatternNameHigher
	}
	if old.PatternNameLower.IsNull() || old.PatternNameLower.ValueString() == "" {
		m.PatternNameLower = old.PatternNameLower
	}
	if old.MetricType.IsNull() || old.MetricType.ValueString() == "" {
		m.MetricType = old.MetricType
	}
	if old.EnableBaselineNearConstance.IsNull() {
		m.EnableBaselineNearConstance = types.BoolNull()
	}
	if old.ComputeDifference.IsNull() {
		m.ComputeDifference = types.BoolNull()
	}
}

// applyMetricConfigurations pushes metric_configurations to the API.
// All metric alert settings are batched into a single POST request.
// Component escalation/ignored settings are still sent per-metric (separate endpoint).
func applyMetricConfigurations(ctx context.Context, c *client.Client, projectName string, planConfigs []metricConfigurationModel, patternIdRule int, samplingIntervalS int64) diag.Diagnostics {
	var diags diag.Diagnostics
	var allPostData []client.MetricAlertSettingPost

	for _, cfg := range planConfigs {
		metricName := cfg.MetricName.ValueString()

		if !cfg.EscalateIncidentComponents.IsNull() && !cfg.EscalateIncidentComponents.IsUnknown() {
			var escalateComps []string
			d := cfg.EscalateIncidentComponents.ElementsAs(ctx, &escalateComps, false)
			diags.Append(d...)
			if !diags.HasError() {
				if err := c.SetMetricComponents(projectName, metricName, "escalateIncident", escalateComps); err != nil {
					diags.AddError(fmt.Sprintf("Error setting escalateIncident components for %s", metricName), err.Error())
				}
			}
		}

		if !cfg.IgnoredComponents.IsNull() && !cfg.IgnoredComponents.IsUnknown() {
			var ignoredComps []string
			d := cfg.IgnoredComponents.ElementsAs(ctx, &ignoredComps, false)
			diags.Append(d...)
			if !diags.HasError() {
				if err := c.SetMetricComponents(projectName, metricName, "ignored", ignoredComps); err != nil {
					diags.AddError(fmt.Sprintf("Error setting ignored components for %s", metricName), err.Error())
				}
			}
		}

		if !cfg.MetricAlertSettings.IsNull() && !cfg.MetricAlertSettings.IsUnknown() {
			var alertSettings []metricAlertSettingModel
			d := cfg.MetricAlertSettings.ElementsAs(ctx, &alertSettings, false)
			diags.Append(d...)
			if !diags.HasError() {
				for _, s := range alertSettings {
					allPostData = append(allPostData, buildSingleMetricAlertSettingPost(metricName, s, samplingIntervalS))
				}
			}
		}
	}

	if len(allPostData) > 0 {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(allPostData); err != nil {
			diags.AddError("Error marshaling metric alert settings", err.Error())
			return diags
		}
		jsonBytes := buf.Bytes()
		if err := c.SetMetricSettings(projectName, patternIdRule, jsonBytes); err != nil {
			diags.AddError("Error setting metric alert settings", err.Error())
		}
	}

	return diags
}

// normalizeMetricConfigsForState converts null component lists to empty lists so that
// Terraform's state always holds [] (not null) for escalate_incident_components and
// ignored_components. Without this, Optional+Computed null attributes cause a perpetual
// plan diff: Terraform plans [] for null Computed list attributes, apply stores null
// from config, and the cycle repeats.
func normalizeMetricConfigsForState(ctx context.Context, configs []metricConfigurationModel) ([]metricConfigurationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	emptyList, d := types.ListValueFrom(ctx, types.StringType, []string{})
	diags.Append(d...)
	if diags.HasError() {
		return configs, diags
	}
	for i, mc := range configs {
		if mc.EscalateIncidentComponents.IsNull() {
			configs[i].EscalateIncidentComponents = emptyList
		}
		if mc.IgnoredComponents.IsNull() {
			configs[i].IgnoredComponents = emptyList
		}
	}
	return configs, diags
}

// splitByComma splits a string by comma. Extracted to avoid repetition in date parsing.
func nullableInt64Div60(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	val := v.ValueInt64() / 60
	return &val
}

func splitByComma(s string) []string {
	var parts []string
	current := ""
	for _, c := range s {
		if c == ',' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
