// Copyright (c) InsightFinder Inc.package provider

// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

// NewProjectResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// projectResource is the resource implementation.
type projectResource struct {
	client *client.Client
}

// projectResourceModel maps the resource schema data.
type projectResourceModel struct {
	ID                    types.String                `tfsdk:"id"`
	ProjectName           types.String                `tfsdk:"project_name"`
	ProjectDisplayName    types.String                `tfsdk:"project_display_name"`
	SystemName            types.String                `tfsdk:"system_name"`
	ProjectCreationConfig *projectCreationConfigModel `tfsdk:"project_creation_config"`
	CValue                types.Int64                 `tfsdk:"c_value"`
	PValue                types.Float64               `tfsdk:"p_value"`
	ProjectTimeZone       types.String                `tfsdk:"project_time_zone"`
	SamplingInterval      types.Int64                 `tfsdk:"sampling_interval"`

	// Basic Configuration
	UBLRetentionTime                   types.Int64   `tfsdk:"ubl_retention_time"`
	AlertAverageTime                   types.Int64   `tfsdk:"alert_average_time"`
	AlertHourlyCost                    types.Float64 `tfsdk:"alert_hourly_cost"`
	AnomalyDetectionMode               types.Int64   `tfsdk:"anomaly_detection_mode"`
	AnomalySamplingInterval            types.Int64   `tfsdk:"anomaly_sampling_interval"`
	AvgPerIncidentDowntimeCost         types.Float64 `tfsdk:"avg_per_incident_downtime_cost"`
	CausalMinDelay                     types.String  `tfsdk:"causal_min_delay"`
	CausalPredictionSetting            types.Int64   `tfsdk:"causal_prediction_setting"`
	ColdEventThreshold                 types.Int64   `tfsdk:"cold_event_threshold"`
	ColdNumberLimit                    types.Int64   `tfsdk:"cold_number_limit"`
	CollectAllRareEventsFlag           types.Bool    `tfsdk:"collect_all_rare_events_flag"`
	DailyModelSpan                     types.Int64   `tfsdk:"daily_model_span"`
	DisableLogCompressEvent            types.Bool    `tfsdk:"disable_log_compress_event"`
	DisableModelKeywordStatsCollection types.Bool    `tfsdk:"disable_model_keyword_stats_collection"`

	// Anomaly and Detection Settings
	EnableAnomalyScoreEscalation    types.Bool    `tfsdk:"enable_anomaly_score_escalation"`
	EnableHotEvent                  types.Bool    `tfsdk:"enable_hot_event"`
	EnableNewAlertEmail             types.Bool    `tfsdk:"enable_new_alert_email"`
	EnableStreamDetection           types.Bool    `tfsdk:"enable_stream_detection"`
	EscalationAnomalyScoreThreshold types.String  `tfsdk:"escalation_anomaly_score_threshold"`
	FeatureOutlierSensitivity       types.String  `tfsdk:"feature_outlier_sensitivity"`
	FeatureOutlierThreshold         types.Float64 `tfsdk:"feature_outlier_threshold"`
	HotEventCalmDownPeriod          types.Int64   `tfsdk:"hot_event_calm_down_period"`
	HotEventDetectionMode           types.Int64   `tfsdk:"hot_event_detection_mode"`
	HotEventThreshold               types.Int64   `tfsdk:"hot_event_threshold"`
	HotNumberLimit                  types.Int64   `tfsdk:"hot_number_limit"`
	IgnoreAnomalyScoreThreshold     types.String  `tfsdk:"ignore_anomaly_score_threshold"`
	IgnoreInstanceForKB             types.Bool    `tfsdk:"ignore_instance_for_kb"`

	// Incident Settings
	IncidentPredictionEventLimit types.Int64 `tfsdk:"incident_prediction_event_limit"`
	IncidentPredictionWindow     types.Int64 `tfsdk:"incident_prediction_window"`
	IncidentRelationSearchWindow types.Int64 `tfsdk:"incident_relation_search_window"`

	// Instance Settings
	InstanceConvertFlag  types.Bool `tfsdk:"instance_convert_flag"`
	InstanceDownEnable   types.Bool `tfsdk:"instance_down_enable"`
	IsEdgeBrain          types.Bool `tfsdk:"is_edge_brain"`
	IsGroupingByInstance types.Bool `tfsdk:"is_grouping_by_instance"`
	IsTracePrompt        types.Bool `tfsdk:"is_trace_prompt"`
	ShowInstanceDown     types.Bool `tfsdk:"show_instance_down"`

	// Log Settings
	KeywordFeatureNumber     types.Int64  `tfsdk:"keyword_feature_number"`
	KeywordSetting           types.Int64  `tfsdk:"keyword_setting"`
	LargeProject             types.Bool   `tfsdk:"large_project"`
	LogAnomalyEventBaseScore types.String `tfsdk:"log_anomaly_event_base_score"`
	LogDetectionMinCount     types.Int64  `tfsdk:"log_detection_min_count"`
	LogDetectionSize         types.Int64  `tfsdk:"log_detection_size"`
	LogPatternLimitLevel     types.Int64  `tfsdk:"log_pattern_limit_level"`
	MaxLogModelSize          types.Int64  `tfsdk:"max_log_model_size"`
	ModelKeywordSetting      types.Int64  `tfsdk:"model_keyword_setting"`
	MultiLineFlag            types.Bool   `tfsdk:"multi_line_flag"`
	NlpFlag                  types.Bool   `tfsdk:"nlp_flag"`
	PrettyJsonConvertorFlag  types.Bool   `tfsdk:"pretty_json_convertor_flag"`

	// Prediction and Root Cause Settings
	MaximumDetectionWaitTime             types.Int64   `tfsdk:"maximum_detection_wait_time"`
	MaximumRootCauseResultSize           types.Int64   `tfsdk:"maximum_root_cause_result_size"`
	MaximumThreads                       types.Int64   `tfsdk:"maximum_threads"`
	MinIncidentPredictionWindow          types.Int64   `tfsdk:"min_incident_prediction_window"`
	MinValidModelSpan                    types.Int64   `tfsdk:"min_valid_model_span"`
	MultiHopSearchLevel                  types.Int64   `tfsdk:"multi_hop_search_level"`
	MultiHopSearchLimit                  types.String  `tfsdk:"multi_hop_search_limit"`
	NewAlertFlag                         types.Bool    `tfsdk:"new_alert_flag"`
	NewPatternNumberLimit                types.Int64   `tfsdk:"new_pattern_number_limit"`
	NewPatternRange                      types.Int64   `tfsdk:"new_pattern_range"`
	NormalEventCausalFlag                types.Bool    `tfsdk:"normal_event_causal_flag"`
	PredictionCountThreshold             types.Int64   `tfsdk:"prediction_count_threshold"`
	PredictionProbabilityThreshold       types.Float64 `tfsdk:"prediction_probability_threshold"`
	PredictionRuleActiveCondition        types.Int64   `tfsdk:"prediction_rule_active_condition"`
	PredictionRuleActiveThreshold        types.Float64 `tfsdk:"prediction_rule_active_threshold"`
	PredictionRuleFalsePositiveThreshold types.Int64   `tfsdk:"prediction_rule_false_positive_threshold"`
	PredictionRuleInactiveThreshold      types.Float64 `tfsdk:"prediction_rule_inactive_threshold"`
	RootCauseCountThreshold              types.Int64   `tfsdk:"root_cause_count_threshold"`
	RootCauseLogMessageSearchRange       types.Int64   `tfsdk:"root_cause_log_message_search_range"`
	RootCauseProbabilityThreshold        types.Float64 `tfsdk:"root_cause_probability_threshold"`
	RootCauseRankSetting                 types.Int64   `tfsdk:"root_cause_rank_setting"`

	// Pattern and Rare Event Settings
	ProjectModelFlag         types.Bool   `tfsdk:"project_model_flag"`
	Proxy                    types.String `tfsdk:"proxy"`
	RareAnomalyType          types.Int64  `tfsdk:"rare_anomaly_type"`
	RareEventAlertThresholds types.Int64  `tfsdk:"rare_event_alert_thresholds"`
	RareNumberLimit          types.Int64  `tfsdk:"rare_number_limit"`
	RetentionTime            types.Int64  `tfsdk:"retention_time"`
	SimilaritySensitivity    types.String `tfsdk:"similarity_sensitivity"`
	TrainingFilter           types.Bool   `tfsdk:"training_filter"`

	// Webhook Settings
	MaxWebHookRequestSize        types.Int64  `tfsdk:"max_web_hook_request_size"`
	WebhookAlertDampening        types.Int64  `tfsdk:"webhook_alert_dampening"`
	WebhookBlackListSetStr       types.String `tfsdk:"webhook_black_list_set_str"`
	WebhookCriticalKeywordSetStr types.String `tfsdk:"webhook_critical_keyword_set_str"`
	WebhookTypeSetStr            types.String `tfsdk:"webhook_type_set_str"`
	WebhookUrl                   types.String `tfsdk:"webhook_url"`
	WhitelistNumberLimit         types.Int64  `tfsdk:"whitelist_number_limit"`
	ZoneNameKey                  types.String `tfsdk:"zone_name_key"`

	// Complex object fields (will use types.String for JSON encoding)
	BaseValueSetting       types.String `tfsdk:"base_value_setting"`
	CdfSetting             types.String `tfsdk:"cdf_setting"`
	EmailSetting           types.String `tfsdk:"email_setting"`
	InstanceGroupingUpdate types.String `tfsdk:"instance_grouping_update"`
	LlmEvaluationSetting   types.String `tfsdk:"llm_evaluation_setting"`
	LogToLogSettingList    types.String `tfsdk:"log_to_log_setting_list"`
	WebhookHeaderList      types.String `tfsdk:"webhook_header_list"`
	SharedUsernames        types.String `tfsdk:"shared_usernames"`
	LogLabelSettings       types.List   `tfsdk:"log_label_settings"`
}

type projectCreationConfigModel struct {
	DataType         types.String `tfsdk:"data_type"`
	InstanceType     types.String `tfsdk:"instance_type"`
	ProjectCloudType types.String `tfsdk:"project_cloud_type"`
	InsightAgentType types.String `tfsdk:"insight_agent_type"`
}

// Metadata returns the resource type name.
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an InsightFinder project.",
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
			"project_creation_config": schema.SingleNestedAttribute{
				Description: "Configuration for creating the project.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"data_type": schema.StringAttribute{
						Description: "The type of data (e.g., Log, Metric, Trace).",
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
			// Additional configuration fields
			"ubl_retention_time": schema.Int64Attribute{
				Description: "Retention time for UBL data in days",
				Optional:    true,
				Computed:    true,
			},
			"alert_average_time": schema.Int64Attribute{
				Description: "Average time for alerts",
				Optional:    true,
				Computed:    true,
			},
			"alert_hourly_cost": schema.Float64Attribute{
				Description: "Hourly cost for alerts",
				Optional:    true,
				Computed:    true,
			},
			"anomaly_detection_mode": schema.Int64Attribute{
				Description: "Anomaly detection mode",
				Optional:    true,
				Computed:    true,
			},
			"anomaly_sampling_interval": schema.Int64Attribute{
				Description: "Sampling interval for anomaly detection",
				Optional:    true,
				Computed:    true,
			},
			"avg_per_incident_downtime_cost": schema.Float64Attribute{
				Description: "Average cost per incident downtime",
				Optional:    true,
				Computed:    true,
			},
			"causal_min_delay": schema.StringAttribute{
				Description: "Minimum delay for causal analysis",
				Optional:    true,
				Computed:    true,
			},
			"causal_prediction_setting": schema.Int64Attribute{
				Description: "Causal prediction setting",
				Optional:    true,
				Computed:    true,
			},
			"cold_event_threshold": schema.Int64Attribute{
				Description: "Threshold for cold events",
				Optional:    true,
				Computed:    true,
			},
			"cold_number_limit": schema.Int64Attribute{
				Description: "Limit for cold numbers",
				Optional:    true,
				Computed:    true,
			},
			"collect_all_rare_events_flag": schema.BoolAttribute{
				Description: "Flag to collect all rare events",
				Optional:    true,
				Computed:    true,
			},
			"daily_model_span": schema.Int64Attribute{
				Description: "Daily model span setting",
				Optional:    true,
				Computed:    true,
			},
			"disable_log_compress_event": schema.BoolAttribute{
				Description: "Disable log compress event",
				Optional:    true,
				Computed:    true,
			},
			"disable_model_keyword_stats_collection": schema.BoolAttribute{
				Description: "Disable model keyword stats collection",
				Optional:    true,
				Computed:    true,
			},
			"enable_anomaly_score_escalation": schema.BoolAttribute{
				Description: "Enable anomaly score escalation",
				Optional:    true,
				Computed:    true,
			},
			"enable_hot_event": schema.BoolAttribute{
				Description: "Enable hot event detection",
				Optional:    true,
				Computed:    true,
			},
			"enable_new_alert_email": schema.BoolAttribute{
				Description: "Enable new alert email notifications",
				Optional:    true,
				Computed:    true,
			},
			"enable_stream_detection": schema.BoolAttribute{
				Description: "Enable stream detection",
				Optional:    true,
				Computed:    true,
			},
			"escalation_anomaly_score_threshold": schema.StringAttribute{
				Description: "Threshold for anomaly score escalation",
				Optional:    true,
				Computed:    true,
			},
			"feature_outlier_sensitivity": schema.StringAttribute{
				Description: "Sensitivity for feature outlier detection",
				Optional:    true,
				Computed:    true,
			},
			"feature_outlier_threshold": schema.Float64Attribute{
				Description: "Threshold for feature outlier detection",
				Optional:    true,
				Computed:    true,
			},
			"hot_event_calm_down_period": schema.Int64Attribute{
				Description: "Calm down period for hot events",
				Optional:    true,
				Computed:    true,
			},
			"hot_event_detection_mode": schema.Int64Attribute{
				Description: "Detection mode for hot events",
				Optional:    true,
				Computed:    true,
			},
			"hot_event_threshold": schema.Int64Attribute{
				Description: "Threshold for hot event detection",
				Optional:    true,
				Computed:    true,
			},
			"hot_number_limit": schema.Int64Attribute{
				Description: "Limit for hot numbers",
				Optional:    true,
				Computed:    true,
			},
			"ignore_anomaly_score_threshold": schema.StringAttribute{
				Description: "Threshold to ignore anomaly scores",
				Optional:    true,
				Computed:    true,
			},
			"ignore_instance_for_kb": schema.BoolAttribute{
				Description: "Ignore instance for knowledge base",
				Optional:    true,
				Computed:    true,
			},
			"incident_prediction_event_limit": schema.Int64Attribute{
				Description: "Event limit for incident prediction",
				Optional:    true,
				Computed:    true,
			},
			"incident_prediction_window": schema.Int64Attribute{
				Description: "Window for incident prediction",
				Optional:    true,
				Computed:    true,
			},
			"incident_relation_search_window": schema.Int64Attribute{
				Description: "Window for incident relation search",
				Optional:    true,
				Computed:    true,
			},
			"instance_convert_flag": schema.BoolAttribute{
				Description: "Flag for instance conversion",
				Optional:    true,
				Computed:    true,
			},
			"instance_down_enable": schema.BoolAttribute{
				Description: "Enable instance down report",
				Optional:    true,
				Computed:    true,
			},
			"is_edge_brain": schema.BoolAttribute{
				Description: "Is edge brain enabled",
				Optional:    true,
				Computed:    true,
			},
			"is_grouping_by_instance": schema.BoolAttribute{
				Description: "Is grouping by instance enabled",
				Optional:    true,
				Computed:    true,
			},
			"is_trace_prompt": schema.BoolAttribute{
				Description: "Is trace prompt enabled",
				Optional:    true,
				Computed:    true,
			},
			"show_instance_down": schema.BoolAttribute{
				Description: "Whether to show instance down incidents for this project",
				Optional:    true,
				Computed:    true,
			},
			"keyword_feature_number": schema.Int64Attribute{
				Description: "Number of keyword features",
				Optional:    true,
				Computed:    true,
			},
			"keyword_setting": schema.Int64Attribute{
				Description: "Keyword setting configuration",
				Optional:    true,
				Computed:    true,
			},
			"large_project": schema.BoolAttribute{
				Description: "Is this a large project",
				Optional:    true,
				Computed:    true,
			},
			"log_anomaly_event_base_score": schema.StringAttribute{
				Description: "Base score for log anomaly events",
				Optional:    true,
				Computed:    true,
			},
			"log_detection_min_count": schema.Int64Attribute{
				Description: "Minimum count for log detection",
				Optional:    true,
				Computed:    true,
			},
			"log_detection_size": schema.Int64Attribute{
				Description: "Size for log detection",
				Optional:    true,
				Computed:    true,
			},
			"log_pattern_limit_level": schema.Int64Attribute{
				Description: "Limit level for log patterns",
				Optional:    true,
				Computed:    true,
			},
			"max_log_model_size": schema.Int64Attribute{
				Description: "Maximum log model size",
				Optional:    true,
				Computed:    true,
			},
			"model_keyword_setting": schema.Int64Attribute{
				Description: "Model keyword setting",
				Optional:    true,
				Computed:    true,
			},
			"multi_line_flag": schema.BoolAttribute{
				Description: "Multi-line flag",
				Optional:    true,
				Computed:    true,
			},
			"nlp_flag": schema.BoolAttribute{
				Description: "NLP flag",
				Optional:    true,
				Computed:    true,
			},
			"pretty_json_convertor_flag": schema.BoolAttribute{
				Description: "Pretty JSON convertor flag",
				Optional:    true,
				Computed:    true,
			},
			"maximum_detection_wait_time": schema.Int64Attribute{
				Description: "Maximum detection wait time",
				Optional:    true,
				Computed:    true,
			},
			"maximum_root_cause_result_size": schema.Int64Attribute{
				Description: "Maximum root cause result size",
				Optional:    true,
				Computed:    true,
			},
			"maximum_threads": schema.Int64Attribute{
				Description: "Maximum number of threads",
				Optional:    true,
				Computed:    true,
			},
			"min_incident_prediction_window": schema.Int64Attribute{
				Description: "Minimum incident prediction window",
				Optional:    true,
				Computed:    true,
			},
			"min_valid_model_span": schema.Int64Attribute{
				Description: "Minimum valid model span",
				Optional:    true,
				Computed:    true,
			},
			"multi_hop_search_level": schema.Int64Attribute{
				Description: "Multi-hop search level",
				Optional:    true,
				Computed:    true,
			},
			"multi_hop_search_limit": schema.StringAttribute{
				Description: "Multi-hop search limit",
				Optional:    true,
				Computed:    true,
			},
			"new_alert_flag": schema.BoolAttribute{
				Description: "New alert flag",
				Optional:    true,
				Computed:    true,
			},
			"new_pattern_number_limit": schema.Int64Attribute{
				Description: "Limit for new pattern numbers",
				Optional:    true,
				Computed:    true,
			},
			"new_pattern_range": schema.Int64Attribute{
				Description: "Range for new patterns",
				Optional:    true,
				Computed:    true,
			},
			"normal_event_causal_flag": schema.BoolAttribute{
				Description: "Normal event causal flag",
				Optional:    true,
				Computed:    true,
			},
			"prediction_count_threshold": schema.Int64Attribute{
				Description: "Threshold for prediction count",
				Optional:    true,
				Computed:    true,
			},
			"prediction_probability_threshold": schema.Float64Attribute{
				Description: "Threshold for prediction probability",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_active_condition": schema.Int64Attribute{
				Description: "Active condition for prediction rules",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_active_threshold": schema.Float64Attribute{
				Description: "Active threshold for prediction rules",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_false_positive_threshold": schema.Int64Attribute{
				Description: "False positive threshold for prediction rules",
				Optional:    true,
				Computed:    true,
			},
			"prediction_rule_inactive_threshold": schema.Float64Attribute{
				Description: "Inactive threshold for prediction rules",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_count_threshold": schema.Int64Attribute{
				Description: "Threshold for root cause count",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_log_message_search_range": schema.Int64Attribute{
				Description: "Search range for root cause log messages",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_probability_threshold": schema.Float64Attribute{
				Description: "Threshold for root cause probability",
				Optional:    true,
				Computed:    true,
			},
			"root_cause_rank_setting": schema.Int64Attribute{
				Description: "Rank setting for root cause",
				Optional:    true,
				Computed:    true,
			},
			"project_model_flag": schema.BoolAttribute{
				Description: "Project model flag",
				Optional:    true,
				Computed:    true,
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy configuration",
				Optional:    true,
				Computed:    true,
			},
			"rare_anomaly_type": schema.Int64Attribute{
				Description: "Type of rare anomaly",
				Optional:    true,
				Computed:    true,
			},
			"rare_event_alert_thresholds": schema.Int64Attribute{
				Description: "Alert thresholds for rare events",
				Optional:    true,
				Computed:    true,
			},
			"rare_number_limit": schema.Int64Attribute{
				Description: "Limit for rare numbers",
				Optional:    true,
				Computed:    true,
			},
			"retention_time": schema.Int64Attribute{
				Description: "The retention time in days",
				Optional:    true,
				Computed:    true,
			},
			"similarity_sensitivity": schema.StringAttribute{
				Description: "Sensitivity for similarity detection",
				Optional:    true,
				Computed:    true,
			},
			"training_filter": schema.BoolAttribute{
				Description: "Training filter flag",
				Optional:    true,
				Computed:    true,
			},
			"max_web_hook_request_size": schema.Int64Attribute{
				Description: "Maximum webhook request size",
				Optional:    true,
				Computed:    true,
			},
			"webhook_alert_dampening": schema.Int64Attribute{
				Description: "Alert dampening for webhooks",
				Optional:    true,
				Computed:    true,
			},
			"webhook_black_list_set_str": schema.StringAttribute{
				Description: "Blacklist set string for webhooks",
				Optional:    true,
				Computed:    true,
			},
			"webhook_critical_keyword_set_str": schema.StringAttribute{
				Description: "Critical keyword set string for webhooks",
				Optional:    true,
				Computed:    true,
			},
			"webhook_type_set_str": schema.StringAttribute{
				Description: "Type set string for webhooks",
				Optional:    true,
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "Webhook URL",
				Optional:    true,
				Computed:    true,
			},
			"whitelist_number_limit": schema.Int64Attribute{
				Description: "Limit for whitelist numbers",
				Optional:    true,
				Computed:    true,
			},
			"zone_name_key": schema.StringAttribute{
				Description: "Zone name key",
				Optional:    true,
				Computed:    true,
			},
			"base_value_setting": schema.StringAttribute{
				Description: "Base value setting configuration (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"cdf_setting": schema.StringAttribute{
				Description: "CDF setting configuration (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"email_setting": schema.StringAttribute{
				Description: "Email notification settings (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"instance_grouping_update": schema.StringAttribute{
				Description: "Instance grouping update settings (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"llm_evaluation_setting": schema.StringAttribute{
				Description: "LLM evaluation settings (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"log_to_log_setting_list": schema.StringAttribute{
				Description: "List of log to log settings (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"webhook_header_list": schema.StringAttribute{
				Description: "List of webhook headers (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"shared_usernames": schema.StringAttribute{
				Description: "List of shared usernames (JSON)",
				Optional:    true,
				Computed:    true,
			},
			"log_label_settings": schema.ListNestedAttribute{
				Description: "List of log label settings for the project. Each setting is applied individually via API.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label_type": schema.StringAttribute{
							Description: "Type of log label (whitelist, blacklist, patternName, etc.)",
							Required:    true,
						},
						"log_label_string": schema.StringAttribute{
							Description: "The log label value/pattern",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// populateSettings converts the Terraform plan/state into a settings map for API calls
func populateSettings(plan *projectResourceModel) map[string]interface{} {
	// Helper function to parse JSON fields
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

	// Build ProjectSettings struct from the plan
	projectSettings := client.ProjectSettings{}

	// Populate only non-null fields to match the struct
	if !plan.ProjectName.IsNull() {
		projectSettings.ProjectName = plan.ProjectName.ValueString()
	}
	if !plan.ProjectDisplayName.IsNull() {
		projectSettings.ProjectDisplayName = plan.ProjectDisplayName.ValueString()
	}
	if !plan.CValue.IsNull() {
		projectSettings.CValue = int(plan.CValue.ValueInt64())
	}
	if !plan.PValue.IsNull() {
		projectSettings.PValue = plan.PValue.ValueFloat64()
	}
	if !plan.ProjectTimeZone.IsNull() {
		projectSettings.ProjectTimeZone = plan.ProjectTimeZone.ValueString()
	}
	if !plan.SamplingInterval.IsNull() {
		projectSettings.SamplingInterval = int(plan.SamplingInterval.ValueInt64())
	}
	if !plan.MinValidModelSpan.IsNull() {
		projectSettings.MinValidModelSpan = int(plan.MinValidModelSpan.ValueInt64())
	}
	if !plan.MaxWebHookRequestSize.IsNull() {
		projectSettings.MaxWebHookRequestSize = int(plan.MaxWebHookRequestSize.ValueInt64())
	}
	if !plan.WebhookUrl.IsNull() {
		projectSettings.WebhookURL = plan.WebhookUrl.ValueString()
	}
	if !plan.WebhookTypeSetStr.IsNull() {
		projectSettings.WebhookTypeSetStr = plan.WebhookTypeSetStr.ValueString()
	}
	if !plan.WebhookBlackListSetStr.IsNull() {
		projectSettings.WebhookBlackListSetStr = plan.WebhookBlackListSetStr.ValueString()
	}
	if !plan.WebhookCriticalKeywordSetStr.IsNull() {
		projectSettings.WebhookCriticalKeywordSetStr = plan.WebhookCriticalKeywordSetStr.ValueString()
	}
	if !plan.WebhookAlertDampening.IsNull() {
		projectSettings.WebhookAlertDampening = int(plan.WebhookAlertDampening.ValueInt64())
	}
	if !plan.Proxy.IsNull() {
		projectSettings.Proxy = plan.Proxy.ValueString()
	}
	if !plan.RetentionTime.IsNull() {
		projectSettings.RetentionTime = int(plan.RetentionTime.ValueInt64())
	}
	if !plan.UBLRetentionTime.IsNull() {
		projectSettings.UBLRetentionTime = int(plan.UBLRetentionTime.ValueInt64())
	}
	if !plan.TrainingFilter.IsNull() {
		projectSettings.TrainingFilter = plan.TrainingFilter.ValueBool()
	}
	if !plan.MultiHopSearchLimit.IsNull() {
		projectSettings.MultiHopSearchLimit = plan.MultiHopSearchLimit.ValueString()
	}
	if !plan.EnableNewAlertEmail.IsNull() {
		projectSettings.EnableNewAlertEmail = plan.EnableNewAlertEmail.ValueBool()
	}
	if !plan.LargeProject.IsNull() {
		projectSettings.LargeProject = plan.LargeProject.ValueBool()
	}
	if !plan.NewPatternRange.IsNull() {
		projectSettings.NewPatternRange = int(plan.NewPatternRange.ValueInt64())
	}
	if !plan.EnableAnomalyScoreEscalation.IsNull() {
		projectSettings.EnableAnomalyScoreEscalation = plan.EnableAnomalyScoreEscalation.ValueBool()
	}
	if !plan.EscalationAnomalyScoreThreshold.IsNull() {
		projectSettings.EscalationAnomalyScoreThreshold = plan.EscalationAnomalyScoreThreshold.ValueString()
	}
	if !plan.IgnoreAnomalyScoreThreshold.IsNull() {
		projectSettings.IgnoreAnomalyScoreThreshold = plan.IgnoreAnomalyScoreThreshold.ValueString()
	}
	if !plan.EnableStreamDetection.IsNull() {
		projectSettings.EnableStreamDetection = plan.EnableStreamDetection.ValueBool()
	}

	// Log-specific fields
	if !plan.DailyModelSpan.IsNull() {
		projectSettings.DailyModelSpan = int(plan.DailyModelSpan.ValueInt64())
	}
	if !plan.KeywordFeatureNumber.IsNull() {
		projectSettings.KeywordFeatureNumber = int(plan.KeywordFeatureNumber.ValueInt64())
	}
	if !plan.MaxLogModelSize.IsNull() {
		projectSettings.MaxLogModelSize = int(plan.MaxLogModelSize.ValueInt64())
	}
	if !plan.ModelKeywordSetting.IsNull() {
		projectSettings.ModelKeywordSetting = int(plan.ModelKeywordSetting.ValueInt64())
	}
	if !plan.NlpFlag.IsNull() {
		projectSettings.NlpFlag = plan.NlpFlag.ValueBool()
	}
	if !plan.ProjectModelFlag.IsNull() {
		projectSettings.ProjectModelFlag = plan.ProjectModelFlag.ValueBool()
	}
	if !plan.MaximumThreads.IsNull() {
		projectSettings.MaximumThreads = int(plan.MaximumThreads.ValueInt64())
	}
	if !plan.LogDetectionMinCount.IsNull() {
		projectSettings.LogDetectionMinCount = int(plan.LogDetectionMinCount.ValueInt64())
	}
	if !plan.LogDetectionSize.IsNull() {
		projectSettings.LogDetectionSize = int(plan.LogDetectionSize.ValueInt64())
	}
	if !plan.MaximumDetectionWaitTime.IsNull() {
		projectSettings.MaximumDetectionWaitTime = int(plan.MaximumDetectionWaitTime.ValueInt64())
	}
	if !plan.KeywordSetting.IsNull() {
		projectSettings.KeywordSetting = int(plan.KeywordSetting.ValueInt64())
	}
	if !plan.LogPatternLimitLevel.IsNull() {
		projectSettings.LogPatternLimitLevel = int(plan.LogPatternLimitLevel.ValueInt64())
	}
	if !plan.NormalEventCausalFlag.IsNull() {
		projectSettings.NormalEventCausalFlag = plan.NormalEventCausalFlag.ValueBool()
	}
	if !plan.SimilaritySensitivity.IsNull() {
		projectSettings.SimilaritySensitivity = plan.SimilaritySensitivity.ValueString()
	}
	if !plan.CollectAllRareEventsFlag.IsNull() {
		projectSettings.CollectAllRareEventsFlag = plan.CollectAllRareEventsFlag.ValueBool()
	}
	if !plan.RareEventAlertThresholds.IsNull() {
		projectSettings.RareEventAlertThresholds = int(plan.RareEventAlertThresholds.ValueInt64())
	}
	if !plan.LogAnomalyEventBaseScore.IsNull() {
		projectSettings.LogAnomalyEventBaseScore = plan.LogAnomalyEventBaseScore.ValueString()
	}
	if !plan.RareNumberLimit.IsNull() {
		projectSettings.RareNumberLimit = int(plan.RareNumberLimit.ValueInt64())
	}
	if !plan.WhitelistNumberLimit.IsNull() {
		projectSettings.WhitelistNumberLimit = int(plan.WhitelistNumberLimit.ValueInt64())
	}
	if !plan.NewPatternNumberLimit.IsNull() {
		projectSettings.NewPatternNumberLimit = int(plan.NewPatternNumberLimit.ValueInt64())
	}
	if !plan.HotNumberLimit.IsNull() {
		projectSettings.HotNumberLimit = int(plan.HotNumberLimit.ValueInt64())
	}
	if !plan.ColdNumberLimit.IsNull() {
		projectSettings.ColdNumberLimit = int(plan.ColdNumberLimit.ValueInt64())
	}
	if !plan.RareAnomalyType.IsNull() {
		projectSettings.RareAnomalyType = int(plan.RareAnomalyType.ValueInt64())
	}
	if !plan.HotEventThreshold.IsNull() {
		projectSettings.HotEventThreshold = int(plan.HotEventThreshold.ValueInt64())
	}
	if !plan.ColdEventThreshold.IsNull() {
		projectSettings.ColdEventThreshold = int(plan.ColdEventThreshold.ValueInt64())
	}
	if !plan.DisableLogCompressEvent.IsNull() {
		projectSettings.DisableLogCompressEvent = plan.DisableLogCompressEvent.ValueBool()
	}
	if !plan.EnableHotEvent.IsNull() {
		projectSettings.EnableHotEvent = plan.EnableHotEvent.ValueBool()
	}
	if !plan.HotEventCalmDownPeriod.IsNull() {
		projectSettings.HotEventCalmDownPeriod = int(plan.HotEventCalmDownPeriod.ValueInt64())
	}
	if !plan.InstanceDownEnable.IsNull() {
		projectSettings.InstanceDownEnable = plan.InstanceDownEnable.ValueBool()
	}
	if !plan.AnomalySamplingInterval.IsNull() {
		projectSettings.AnomalySamplingInterval = int(plan.AnomalySamplingInterval.ValueInt64())
	}
	if !plan.HotEventDetectionMode.IsNull() {
		projectSettings.HotEventDetectionMode = int(plan.HotEventDetectionMode.ValueInt64())
	}
	if !plan.AnomalyDetectionMode.IsNull() {
		projectSettings.AnomalyDetectionMode = int(plan.AnomalyDetectionMode.ValueInt64())
	}
	if !plan.PrettyJsonConvertorFlag.IsNull() {
		projectSettings.PrettyJSONConvertorFlag = plan.PrettyJsonConvertorFlag.ValueBool()
	}
	if !plan.ZoneNameKey.IsNull() {
		projectSettings.ZoneNameKey = plan.ZoneNameKey.ValueString()
	}
	if !plan.MultiLineFlag.IsNull() {
		projectSettings.MultiLineFlag = plan.MultiLineFlag.ValueBool()
	}
	if !plan.FeatureOutlierSensitivity.IsNull() {
		projectSettings.FeatureOutlierSensitivity = plan.FeatureOutlierSensitivity.ValueString()
	}
	if !plan.DisableModelKeywordStatsCollection.IsNull() {
		projectSettings.DisableModelKeywordStatsCollection = plan.DisableModelKeywordStatsCollection.ValueBool()
	}
	if !plan.InstanceConvertFlag.IsNull() {
		projectSettings.InstanceConvertFlag = plan.InstanceConvertFlag.ValueBool()
	}
	if !plan.NewAlertFlag.IsNull() {
		projectSettings.NewAlertFlag = plan.NewAlertFlag.ValueBool()
	}
	if !plan.IsGroupingByInstance.IsNull() {
		projectSettings.IsGroupingByInstance = plan.IsGroupingByInstance.ValueBool()
	}
	if !plan.FeatureOutlierThreshold.IsNull() {
		projectSettings.FeatureOutlierThreshold = plan.FeatureOutlierThreshold.ValueFloat64()
	}
	if !plan.IsTracePrompt.IsNull() {
		projectSettings.IsTracePrompt = plan.IsTracePrompt.ValueBool()
	}
	if !plan.IsEdgeBrain.IsNull() {
		projectSettings.IsEdgeBrain = plan.IsEdgeBrain.ValueBool()
	}

	// Incident prediction and RCA fields
	if !plan.IncidentPredictionWindow.IsNull() {
		projectSettings.IncidentPredictionWindow = int(plan.IncidentPredictionWindow.ValueInt64())
	}
	if !plan.MinIncidentPredictionWindow.IsNull() {
		projectSettings.MinIncidentPredictionWindow = int(plan.MinIncidentPredictionWindow.ValueInt64())
	}
	if !plan.IncidentRelationSearchWindow.IsNull() {
		projectSettings.IncidentRelationSearchWindow = int(plan.IncidentRelationSearchWindow.ValueInt64())
	}
	if !plan.IncidentPredictionEventLimit.IsNull() {
		projectSettings.IncidentPredictionEventLimit = int(plan.IncidentPredictionEventLimit.ValueInt64())
	}
	if !plan.RootCauseCountThreshold.IsNull() {
		projectSettings.RootCauseCountThreshold = int(plan.RootCauseCountThreshold.ValueInt64())
	}
	if !plan.RootCauseProbabilityThreshold.IsNull() {
		projectSettings.RootCauseProbabilityThreshold = plan.RootCauseProbabilityThreshold.ValueFloat64()
	}
	if !plan.RootCauseLogMessageSearchRange.IsNull() {
		projectSettings.RootCauseLogMessageSearchRange = int(plan.RootCauseLogMessageSearchRange.ValueInt64())
	}
	if !plan.CausalPredictionSetting.IsNull() {
		projectSettings.CausalPredictionSetting = int(plan.CausalPredictionSetting.ValueInt64())
	}
	if !plan.CausalMinDelay.IsNull() {
		projectSettings.CausalMinDelay = plan.CausalMinDelay.ValueString()
	}
	if !plan.RootCauseRankSetting.IsNull() {
		projectSettings.RootCauseRankSetting = int(plan.RootCauseRankSetting.ValueInt64())
	}
	if !plan.MaximumRootCauseResultSize.IsNull() {
		projectSettings.MaximumRootCauseResultSize = int(plan.MaximumRootCauseResultSize.ValueInt64())
	}
	if !plan.MultiHopSearchLevel.IsNull() {
		projectSettings.MultiHopSearchLevel = int(plan.MultiHopSearchLevel.ValueInt64())
	}
	if !plan.AvgPerIncidentDowntimeCost.IsNull() {
		projectSettings.AvgPerIncidentDowntimeCost = plan.AvgPerIncidentDowntimeCost.ValueFloat64()
	}
	if !plan.PredictionRuleActiveCondition.IsNull() {
		projectSettings.PredictionRuleActiveCondition = int(plan.PredictionRuleActiveCondition.ValueInt64())
	}
	if !plan.PredictionRuleFalsePositiveThreshold.IsNull() {
		projectSettings.PredictionRuleFalsePositiveThreshold = int(plan.PredictionRuleFalsePositiveThreshold.ValueInt64())
	}
	if !plan.PredictionRuleActiveThreshold.IsNull() {
		projectSettings.PredictionRuleActiveThreshold = plan.PredictionRuleActiveThreshold.ValueFloat64()
	}
	if !plan.PredictionRuleInactiveThreshold.IsNull() {
		projectSettings.PredictionRuleInactiveThreshold = plan.PredictionRuleInactiveThreshold.ValueFloat64()
	}
	if !plan.PredictionProbabilityThreshold.IsNull() {
		projectSettings.PredictionProbabilityThreshold = plan.PredictionProbabilityThreshold.ValueFloat64()
	}
	if !plan.AlertHourlyCost.IsNull() {
		projectSettings.AlertHourlyCost = plan.AlertHourlyCost.ValueFloat64()
	}
	if !plan.AlertAverageTime.IsNull() {
		projectSettings.AlertAverageTime = int(plan.AlertAverageTime.ValueInt64())
	}
	if !plan.IgnoreInstanceForKB.IsNull() {
		projectSettings.IgnoreInstanceForKB = plan.IgnoreInstanceForKB.ValueBool()
	}
	if !plan.ShowInstanceDown.IsNull() {
		projectSettings.ShowInstanceDown = plan.ShowInstanceDown.ValueBool()
	}
	if !plan.PredictionCountThreshold.IsNull() {
		projectSettings.PredictionCountThreshold = int(plan.PredictionCountThreshold.ValueInt64())
	}

	// Complex JSON fields - parse and assign to struct
	if !plan.BaseValueSetting.IsNull() {
		if parsed := parseJSONField(plan.BaseValueSetting.ValueString()); parsed != nil {
			// Marshal to JSON and unmarshal to the struct field
			if jsonBytes, err := json.Marshal(parsed); err == nil {
				_ = json.Unmarshal(jsonBytes, &projectSettings.BaseValueSetting)
			}
		}
	}
	if !plan.CdfSetting.IsNull() {
		projectSettings.CdfSetting = parseJSONField(plan.CdfSetting.ValueString()).([]interface{})
	}
	if !plan.EmailSetting.IsNull() {
		if parsed := parseJSONField(plan.EmailSetting.ValueString()); parsed != nil {
			if jsonBytes, err := json.Marshal(parsed); err == nil {
				_ = json.Unmarshal(jsonBytes, &projectSettings.EmailSetting)
			}
		}
	}
	if !plan.InstanceGroupingUpdate.IsNull() {
		if parsed := parseJSONField(plan.InstanceGroupingUpdate.ValueString()); parsed != nil {
			if jsonBytes, err := json.Marshal(parsed); err == nil {
				_ = json.Unmarshal(jsonBytes, &projectSettings.InstanceGroupingUpdate)
			}
		}
	}
	if !plan.LlmEvaluationSetting.IsNull() {
		if parsed := parseJSONField(plan.LlmEvaluationSetting.ValueString()); parsed != nil {
			if jsonBytes, err := json.Marshal(parsed); err == nil {
				_ = json.Unmarshal(jsonBytes, &projectSettings.LlmEvaluationSetting)
			}
		}
	}
	if !plan.LogToLogSettingList.IsNull() {
		projectSettings.LogToLogSettingList = parseJSONField(plan.LogToLogSettingList.ValueString()).([]interface{})
	}
	if !plan.WebhookHeaderList.IsNull() {
		projectSettings.WebhookHeaderList = parseJSONField(plan.WebhookHeaderList.ValueString()).([]interface{})
	}
	if !plan.SharedUsernames.IsNull() {
		projectSettings.SharedUsernames = parseJSONField(plan.SharedUsernames.ValueString()).([]interface{})
	}

	// Convert struct to map[string]interface{} using JSON marshal/unmarshal
	// This automatically excludes omitempty fields that are zero values
	jsonBytes, err := json.Marshal(projectSettings)
	if err != nil {
		// Fallback to empty map if marshaling fails
		return make(map[string]interface{})
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &settings); err != nil {
		// Fallback to empty map if unmarshaling fails
		return make(map[string]interface{})
	}

	return settings
}

// Create creates the resource and sets the initial Terraform state.
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating project", map[string]any{"project_name": plan.ProjectName.ValueString()})

	// Create the project via API
	projectConfig := &client.ProjectConfig{
		ProjectName:         plan.ProjectName.ValueString(),
		ProjectDisplayName:  plan.ProjectDisplayName.ValueString(),
		SystemName:          plan.SystemName.ValueString(),
		DataType:            plan.ProjectCreationConfig.DataType.ValueString(),
		InstanceType:        plan.ProjectCreationConfig.InstanceType.ValueString(),
		ProjectCloudType:    plan.ProjectCreationConfig.ProjectCloudType.ValueString(),
		InsightAgentType:    plan.ProjectCreationConfig.InsightAgentType.ValueString(),
		CValue:              int(plan.CValue.ValueInt64()),
		PValue:              plan.PValue.ValueFloat64(),
	}

	err := r.client.CreateProject(projectConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the ID
	plan.ID = plan.ProjectName

	// Apply additional settings if any are provided
	settings := populateSettings(&plan)
	// Only update if we have settings beyond just the project name
	if len(settings) > 0 {
		tflog.Debug(ctx, "Applying additional project settings", map[string]any{"settings_count": len(settings)})
		updateConfig := &client.ProjectConfig{
			ProjectName: plan.ProjectName.ValueString(),
			Settings:    settings,
		}
		err = r.client.UpdateProject(updateConfig)
		if err != nil {
			// Log the error but don't fail - project is created
			tflog.Warn(ctx, "Could not apply all settings on creation", map[string]any{
				"error": err.Error(),
				"note":  "Settings can be applied on next terraform apply",
			})
		}
	}

	// Read back the project configuration after creation to populate computed fields
	// We need to merge config values (from req.Config) with API values
	var config projectResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading project configuration after creation")
	project, err := r.client.GetProject(plan.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		tflog.Warn(ctx, "Could not read project after creation", map[string]any{
			"error": err.Error(),
			"note":  "State may not reflect all API values",
		})
		// If we can't read from API, use the config as state
		resp.State.Set(ctx, config)
		return
	}

	if project != nil {
		// Start with config values
		plan = config
		plan.ID = plan.ProjectName

		// Populate all fields from the API response (same logic as Read/Update)
		settings := project.Settings
		if settings == nil {
			settings = make(map[string]interface{})
		}

		// Helper functions
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
				// If it's already a string, return it
				if str, ok := val.(string); ok {
					return types.StringValue(str)
				}
				// Otherwise, marshal it to JSON
				if jsonBytes, err := json.Marshal(val); err == nil {
					return types.StringValue(string(jsonBytes))
				}
			}
			return types.StringNull()
		}

		// Populate all fields from API response
		plan.ProjectDisplayName = getString("projectDisplayName")
		plan.CValue = getInt64("cValue")
		plan.PValue = getFloat64("pValue")
		plan.ProjectTimeZone = getString("projectTimeZone")
		plan.SamplingInterval = getInt64("samplingInterval")

		// Basic Configuration
		plan.UBLRetentionTime = getInt64("UBLRetentionTime")
		plan.AlertAverageTime = getInt64("alertAverageTime")
		plan.AlertHourlyCost = getFloat64("alertHourlyCost")
		plan.AnomalyDetectionMode = getInt64("anomalyDetectionMode")
		plan.AnomalySamplingInterval = getInt64("anomalySamplingInterval")
		plan.AvgPerIncidentDowntimeCost = getFloat64("avgPerIncidentDowntimeCost")
		plan.CausalPredictionSetting = getInt64("causalPredictionSetting")
		plan.CausalMinDelay = getString("causalMinDelay")
		plan.ColdEventThreshold = getInt64("coldEventThreshold")
		plan.ColdNumberLimit = getInt64("coldNumberLimit")
		plan.CollectAllRareEventsFlag = getBool("collectAllRareEventsFlag")
		plan.DailyModelSpan = getInt64("dailyModelSpan")
		plan.DisableLogCompressEvent = getBool("disableLogCompressEvent")
		plan.DisableModelKeywordStatsCollection = getBool("disableModelKeywordStatsCollection")

		// Anomaly and Detection Settings
		plan.EnableAnomalyScoreEscalation = getBool("enableAnomalyScoreEscalation")
		plan.EnableHotEvent = getBool("enableHotEvent")
		plan.EnableNewAlertEmail = getBool("enableNewAlertEmail")
		plan.EnableStreamDetection = getBool("enableStreamDetection")
		plan.EscalationAnomalyScoreThreshold = getString("escalationAnomalyScoreThreshold")
		plan.FeatureOutlierSensitivity = getString("featureOutlierSensitivity")
		plan.FeatureOutlierThreshold = getFloat64("featureOutlierThreshold")
		plan.HotEventCalmDownPeriod = getInt64("hotEventCalmDownPeriod")
		plan.HotEventDetectionMode = getInt64("hotEventDetectionMode")
		plan.HotEventThreshold = getInt64("hotEventThreshold")
		plan.HotNumberLimit = getInt64("hotNumberLimit")
		plan.IgnoreAnomalyScoreThreshold = getString("ignoreAnomalyScoreThreshold")
		plan.IgnoreInstanceForKB = getBool("ignoreInstanceForKB")

		// Incident Settings
		plan.IncidentPredictionEventLimit = getInt64("incidentPredictionEventLimit")
		plan.IncidentPredictionWindow = getInt64("incidentPredictionWindow")
		plan.IncidentRelationSearchWindow = getInt64("incidentRelationSearchWindow")

		// Instance Settings
		plan.InstanceConvertFlag = getBool("instanceConvertFlag")
		plan.InstanceDownEnable = getBool("instanceDownEnable")
		plan.IsEdgeBrain = getBool("isEdgeBrain")
		plan.IsGroupingByInstance = getBool("isGroupingByInstance")
		plan.IsTracePrompt = getBool("isTracePrompt")
		plan.ShowInstanceDown = getBool("showInstanceDown")

		// Log Settings
		plan.KeywordFeatureNumber = getInt64("keywordFeatureNumber")
		plan.KeywordSetting = getInt64("keywordSetting")
		plan.LargeProject = getBool("largeProject")
		plan.LogAnomalyEventBaseScore = getString("logAnomalyEventBaseScore")
		plan.LogDetectionMinCount = getInt64("logDetectionMinCount")
		plan.LogDetectionSize = getInt64("logDetectionSize")
		plan.LogPatternLimitLevel = getInt64("logPatternLimitLevel")
		plan.MaxLogModelSize = getInt64("maxLogModelSize")
		plan.MaximumDetectionWaitTime = getInt64("maximumDetectionWaitTime")
		plan.MaximumThreads = getInt64("maximumThreads")
		plan.ModelKeywordSetting = getInt64("modelKeywordSetting")
		plan.MultiLineFlag = getBool("multiLineFlag")
		plan.NlpFlag = getBool("nlpFlag")
		plan.PrettyJsonConvertorFlag = getBool("prettyJsonConvertorFlag")

		// Model Settings
		plan.MaximumRootCauseResultSize = getInt64("maximumRootCauseResultSize")
		plan.MinIncidentPredictionWindow = getInt64("minIncidentPredictionWindow")
		plan.MinValidModelSpan = getInt64("minValidModelSpan")
		plan.MultiHopSearchLevel = getInt64("multiHopSearchLevel")
		plan.MultiHopSearchLimit = getString("multiHopSearchLimit")

		// Pattern and Event Settings
		plan.NewAlertFlag = getBool("newAlertFlag")
		plan.NewPatternNumberLimit = getInt64("newPatternNumberLimit")
		plan.NewPatternRange = getInt64("newPatternRange")
		plan.NormalEventCausalFlag = getBool("normalEventCausalFlag")

		// Prediction Settings
		plan.PredictionCountThreshold = getInt64("predictionCountThreshold")
		plan.PredictionProbabilityThreshold = getFloat64("predictionProbabilityThreshold")
		plan.PredictionRuleActiveCondition = getInt64("predictionRuleActiveCondition")
		plan.PredictionRuleActiveThreshold = getFloat64("predictionRuleActiveThreshold")
		plan.PredictionRuleFalsePositiveThreshold = getInt64("predictionRuleFalsePositiveThreshold")
		plan.PredictionRuleInactiveThreshold = getFloat64("predictionRuleInactiveThreshold")
		plan.ProjectModelFlag = getBool("projectModelFlag")
		plan.Proxy = getString("proxy")

		// Rare Event Settings
		plan.RareAnomalyType = getInt64("rareAnomalyType")
		plan.RareEventAlertThresholds = getInt64("rareEventAlertThresholds")
		plan.RareNumberLimit = getInt64("rareNumberLimit")
		plan.RetentionTime = getInt64("retentionTime")

		// Root Cause Settings
		plan.RootCauseCountThreshold = getInt64("rootCauseCountThreshold")
		plan.RootCauseLogMessageSearchRange = getInt64("rootCauseLogMessageSearchRange")
		plan.RootCauseProbabilityThreshold = getFloat64("rootCauseProbabilityThreshold")
		plan.RootCauseRankSetting = getInt64("rootCauseRankSetting")

		// Similarity and Training
		plan.SimilaritySensitivity = getString("similaritySensitivity")
		plan.TrainingFilter = getBool("trainingFilter")

		// Webhook Settings
		plan.MaxWebHookRequestSize = getInt64("maxWebHookRequestSize")
		plan.WebhookAlertDampening = getInt64("webhookAlertDampening")
		plan.WebhookBlackListSetStr = getString("webhookBlackListSetStr")
		plan.WebhookCriticalKeywordSetStr = getString("webhookCriticalKeywordSetStr")
		plan.WebhookTypeSetStr = getString("webhookTypeSetStr")
		plan.WebhookUrl = getString("webhookUrl")
		plan.WhitelistNumberLimit = getInt64("whitelistNumberLimit")
		plan.ZoneNameKey = getString("zoneNameKey")

		// Metric Project Fields

		// JSON String Fields
		plan.BaseValueSetting = getJSONString("baseValueSetting")
		plan.CdfSetting = getJSONString("cdfSetting")
		plan.EmailSetting = getJSONString("emailSetting")
		plan.InstanceGroupingUpdate = getJSONString("instanceGroupingUpdate")
		plan.LlmEvaluationSetting = getJSONString("llmEvaluationSetting")
		plan.LogToLogSettingList = getJSONString("logToLogSettingList")
		plan.WebhookHeaderList = getJSONString("webhookHeaderList")
		plan.SharedUsernames = getJSONString("sharedUsernames")
	}

	// Always preserve config values over API values for fields explicitly set by user
	// This ensures user-specified values take precedence
	if !config.ProjectDisplayName.IsNull() {
		plan.ProjectDisplayName = config.ProjectDisplayName
	}
	if !config.CValue.IsNull() {
		plan.CValue = config.CValue
	}
	if !config.PValue.IsNull() {
		plan.PValue = config.PValue
	}
	if !config.ProjectTimeZone.IsNull() {
		plan.ProjectTimeZone = config.ProjectTimeZone
	}
	if !config.SamplingInterval.IsNull() {
		plan.SamplingInterval = config.SamplingInterval
	}
	if !config.UBLRetentionTime.IsNull() {
		plan.UBLRetentionTime = config.UBLRetentionTime
	}
	if !config.AlertAverageTime.IsNull() {
		plan.AlertAverageTime = config.AlertAverageTime
	}
	if !config.AlertHourlyCost.IsNull() {
		plan.AlertHourlyCost = config.AlertHourlyCost
	}
	if !config.AnomalyDetectionMode.IsNull() {
		plan.AnomalyDetectionMode = config.AnomalyDetectionMode
	}
	if !config.AnomalySamplingInterval.IsNull() {
		plan.AnomalySamplingInterval = config.AnomalySamplingInterval
	}
	if !config.AvgPerIncidentDowntimeCost.IsNull() {
		plan.AvgPerIncidentDowntimeCost = config.AvgPerIncidentDowntimeCost
	}
	if !config.CausalPredictionSetting.IsNull() {
		plan.CausalPredictionSetting = config.CausalPredictionSetting
	}
	if !config.CausalMinDelay.IsNull() {
		plan.CausalMinDelay = config.CausalMinDelay
	}
	if !config.ColdEventThreshold.IsNull() {
		plan.ColdEventThreshold = config.ColdEventThreshold
	}
	if !config.ColdNumberLimit.IsNull() {
		plan.ColdNumberLimit = config.ColdNumberLimit
	}
	if !config.CollectAllRareEventsFlag.IsNull() {
		plan.CollectAllRareEventsFlag = config.CollectAllRareEventsFlag
	}
	if !config.DailyModelSpan.IsNull() {
		plan.DailyModelSpan = config.DailyModelSpan
	}
	if !config.DisableLogCompressEvent.IsNull() {
		plan.DisableLogCompressEvent = config.DisableLogCompressEvent
	}
	if !config.DisableModelKeywordStatsCollection.IsNull() {
		plan.DisableModelKeywordStatsCollection = config.DisableModelKeywordStatsCollection
	}
	if !config.EnableAnomalyScoreEscalation.IsNull() {
		plan.EnableAnomalyScoreEscalation = config.EnableAnomalyScoreEscalation
	}
	if !config.EnableHotEvent.IsNull() {
		plan.EnableHotEvent = config.EnableHotEvent
	}
	if !config.EnableNewAlertEmail.IsNull() {
		plan.EnableNewAlertEmail = config.EnableNewAlertEmail
	}
	if !config.EnableStreamDetection.IsNull() {
		plan.EnableStreamDetection = config.EnableStreamDetection
	}
	if !config.EscalationAnomalyScoreThreshold.IsNull() {
		plan.EscalationAnomalyScoreThreshold = config.EscalationAnomalyScoreThreshold
	}
	if !config.FeatureOutlierSensitivity.IsNull() {
		plan.FeatureOutlierSensitivity = config.FeatureOutlierSensitivity
	}
	if !config.FeatureOutlierThreshold.IsNull() {
		plan.FeatureOutlierThreshold = config.FeatureOutlierThreshold
	}
	if !config.HotEventCalmDownPeriod.IsNull() {
		plan.HotEventCalmDownPeriod = config.HotEventCalmDownPeriod
	}
	if !config.HotEventDetectionMode.IsNull() {
		plan.HotEventDetectionMode = config.HotEventDetectionMode
	}
	if !config.HotEventThreshold.IsNull() {
		plan.HotEventThreshold = config.HotEventThreshold
	}
	if !config.HotNumberLimit.IsNull() {
		plan.HotNumberLimit = config.HotNumberLimit
	}
	if !config.IgnoreAnomalyScoreThreshold.IsNull() {
		plan.IgnoreAnomalyScoreThreshold = config.IgnoreAnomalyScoreThreshold
	}
	if !config.IgnoreInstanceForKB.IsNull() {
		plan.IgnoreInstanceForKB = config.IgnoreInstanceForKB
	}
	if !config.IncidentPredictionEventLimit.IsNull() {
		plan.IncidentPredictionEventLimit = config.IncidentPredictionEventLimit
	}
	if !config.IncidentPredictionWindow.IsNull() {
		plan.IncidentPredictionWindow = config.IncidentPredictionWindow
	}
	if !config.IncidentRelationSearchWindow.IsNull() {
		plan.IncidentRelationSearchWindow = config.IncidentRelationSearchWindow
	}
	if !config.InstanceConvertFlag.IsNull() {
		plan.InstanceConvertFlag = config.InstanceConvertFlag
	}
	if !config.InstanceDownEnable.IsNull() {
		plan.InstanceDownEnable = config.InstanceDownEnable
	}
	if !config.IsEdgeBrain.IsNull() {
		plan.IsEdgeBrain = config.IsEdgeBrain
	}
	if !config.IsGroupingByInstance.IsNull() {
		plan.IsGroupingByInstance = config.IsGroupingByInstance
	}
	if !config.IsTracePrompt.IsNull() {
		plan.IsTracePrompt = config.IsTracePrompt
	}
	if !config.ShowInstanceDown.IsNull() {
		plan.ShowInstanceDown = config.ShowInstanceDown
	}
	if !config.KeywordFeatureNumber.IsNull() {
		plan.KeywordFeatureNumber = config.KeywordFeatureNumber
	}
	if !config.KeywordSetting.IsNull() {
		plan.KeywordSetting = config.KeywordSetting
	}
	if !config.LargeProject.IsNull() {
		plan.LargeProject = config.LargeProject
	}
	if !config.LogAnomalyEventBaseScore.IsNull() {
		plan.LogAnomalyEventBaseScore = config.LogAnomalyEventBaseScore
	}
	if !config.LogDetectionMinCount.IsNull() {
		plan.LogDetectionMinCount = config.LogDetectionMinCount
	}
	if !config.LogDetectionSize.IsNull() {
		plan.LogDetectionSize = config.LogDetectionSize
	}
	if !config.LogPatternLimitLevel.IsNull() {
		plan.LogPatternLimitLevel = config.LogPatternLimitLevel
	}
	if !config.MaxLogModelSize.IsNull() {
		plan.MaxLogModelSize = config.MaxLogModelSize
	}
	if !config.MaximumDetectionWaitTime.IsNull() {
		plan.MaximumDetectionWaitTime = config.MaximumDetectionWaitTime
	}
	if !config.MaximumThreads.IsNull() {
		plan.MaximumThreads = config.MaximumThreads
	}
	if !config.ModelKeywordSetting.IsNull() {
		plan.ModelKeywordSetting = config.ModelKeywordSetting
	}
	if !config.MultiLineFlag.IsNull() {
		plan.MultiLineFlag = config.MultiLineFlag
	}
	if !config.NlpFlag.IsNull() {
		plan.NlpFlag = config.NlpFlag
	}
	if !config.PrettyJsonConvertorFlag.IsNull() {
		plan.PrettyJsonConvertorFlag = config.PrettyJsonConvertorFlag
	}
	if !config.MaximumRootCauseResultSize.IsNull() {
		plan.MaximumRootCauseResultSize = config.MaximumRootCauseResultSize
	}
	if !config.MinIncidentPredictionWindow.IsNull() {
		plan.MinIncidentPredictionWindow = config.MinIncidentPredictionWindow
	}
	if !config.MinValidModelSpan.IsNull() {
		plan.MinValidModelSpan = config.MinValidModelSpan
	}
	if !config.MultiHopSearchLevel.IsNull() {
		plan.MultiHopSearchLevel = config.MultiHopSearchLevel
	}
	if !config.MultiHopSearchLimit.IsNull() {
		plan.MultiHopSearchLimit = config.MultiHopSearchLimit
	}
	if !config.NewAlertFlag.IsNull() {
		plan.NewAlertFlag = config.NewAlertFlag
	}
	if !config.NewPatternNumberLimit.IsNull() {
		plan.NewPatternNumberLimit = config.NewPatternNumberLimit
	}
	if !config.NewPatternRange.IsNull() {
		plan.NewPatternRange = config.NewPatternRange
	}
	if !config.NormalEventCausalFlag.IsNull() {
		plan.NormalEventCausalFlag = config.NormalEventCausalFlag
	}
	if !config.PredictionCountThreshold.IsNull() {
		plan.PredictionCountThreshold = config.PredictionCountThreshold
	}
	if !config.PredictionProbabilityThreshold.IsNull() {
		plan.PredictionProbabilityThreshold = config.PredictionProbabilityThreshold
	}
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
	if !config.ProjectModelFlag.IsNull() {
		plan.ProjectModelFlag = config.ProjectModelFlag
	}
	if !config.Proxy.IsNull() {
		plan.Proxy = config.Proxy
	}
	if !config.RareAnomalyType.IsNull() {
		plan.RareAnomalyType = config.RareAnomalyType
	}
	if !config.RareEventAlertThresholds.IsNull() {
		plan.RareEventAlertThresholds = config.RareEventAlertThresholds
	}
	if !config.RareNumberLimit.IsNull() {
		plan.RareNumberLimit = config.RareNumberLimit
	}
	if !config.RetentionTime.IsNull() {
		plan.RetentionTime = config.RetentionTime
	}
	if !config.RootCauseCountThreshold.IsNull() {
		plan.RootCauseCountThreshold = config.RootCauseCountThreshold
	}
	if !config.RootCauseLogMessageSearchRange.IsNull() {
		plan.RootCauseLogMessageSearchRange = config.RootCauseLogMessageSearchRange
	}
	if !config.RootCauseProbabilityThreshold.IsNull() {
		plan.RootCauseProbabilityThreshold = config.RootCauseProbabilityThreshold
	}
	if !config.RootCauseRankSetting.IsNull() {
		plan.RootCauseRankSetting = config.RootCauseRankSetting
	}
	if !config.SimilaritySensitivity.IsNull() {
		plan.SimilaritySensitivity = config.SimilaritySensitivity
	}
	if !config.TrainingFilter.IsNull() {
		plan.TrainingFilter = config.TrainingFilter
	}
	if !config.MaxWebHookRequestSize.IsNull() {
		plan.MaxWebHookRequestSize = config.MaxWebHookRequestSize
	}
	if !config.WebhookAlertDampening.IsNull() {
		plan.WebhookAlertDampening = config.WebhookAlertDampening
	}
	if !config.WebhookBlackListSetStr.IsNull() {
		plan.WebhookBlackListSetStr = config.WebhookBlackListSetStr
	}
	if !config.WebhookCriticalKeywordSetStr.IsNull() {
		plan.WebhookCriticalKeywordSetStr = config.WebhookCriticalKeywordSetStr
	}
	if !config.WebhookTypeSetStr.IsNull() {
		plan.WebhookTypeSetStr = config.WebhookTypeSetStr
	}
	if !config.WebhookUrl.IsNull() {
		plan.WebhookUrl = config.WebhookUrl
	}
	if !config.WhitelistNumberLimit.IsNull() {
		plan.WhitelistNumberLimit = config.WhitelistNumberLimit
	}
	if !config.ZoneNameKey.IsNull() {
		plan.ZoneNameKey = config.ZoneNameKey
	}
	if !config.BaseValueSetting.IsNull() {
		plan.BaseValueSetting = config.BaseValueSetting
	}
	if !config.CdfSetting.IsNull() {
		plan.CdfSetting = config.CdfSetting
	}
	if !config.EmailSetting.IsNull() {
		plan.EmailSetting = config.EmailSetting
	}
	if !config.InstanceGroupingUpdate.IsNull() {
		plan.InstanceGroupingUpdate = config.InstanceGroupingUpdate
	}
	if !config.LlmEvaluationSetting.IsNull() {
		plan.LlmEvaluationSetting = config.LlmEvaluationSetting
	}
	if !config.LogToLogSettingList.IsNull() {
		plan.LogToLogSettingList = config.LogToLogSettingList
	}
	if !config.WebhookHeaderList.IsNull() {
		plan.WebhookHeaderList = config.WebhookHeaderList
	}
	if !config.SharedUsernames.IsNull() {
		plan.SharedUsernames = config.SharedUsernames
	}

	// Process log_label_settings if provided - each setting must be applied individually
	if !config.LogLabelSettings.IsNull() && !config.LogLabelSettings.IsUnknown() {
		var configSettings []logLabelSettingModel
		diags = config.LogLabelSettings.ElementsAs(ctx, &configSettings, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(configSettings) > 0 {
			tflog.Info(ctx, "Processing log label settings", map[string]any{"count": len(configSettings)})

			// Convert terraform model to client model
			settings := make([]*client.LogLabelSetting, 0, len(configSettings))
			for _, setting := range configSettings {
				settings = append(settings, &client.LogLabelSetting{
					LabelType:      setting.LabelType.ValueString(),
					LogLabelString: setting.LogLabelString.ValueString(),
				})
			}

			// Apply all settings (function will iterate and call API for each)
			err := r.client.CreateOrUpdateLogLabels(
				plan.ProjectName.ValueString(),
				r.client.Username,
				settings,
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error applying log label settings",
					fmt.Sprintf("Could not apply log label settings: %s", err.Error()),
				)
				return
			}

			// Use config values in state (normalized for consistency)
			// Note: GetLogLabels uses a different endpoint that may not immediately reflect changes
			normalizedSettings := make([]logLabelSettingModel, 0, len(configSettings))
			for _, setting := range configSettings {
				normalizedSettings = append(normalizedSettings, logLabelSettingModel{
					LabelType:      setting.LabelType,
					LogLabelString: types.StringValue(normalizeJSON(setting.LogLabelString.ValueString())),
				})
			}

			// Convert back to types.List
			listValue, diags := types.ListValueFrom(ctx, config.LogLabelSettings.ElementType(ctx), normalizedSettings)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			plan.LogLabelSettings = listValue
		}
	}

	// If LogLabelSettings is null or empty, set it to an empty list
	if config.LogLabelSettings.IsNull() {
		emptyList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"label_type":       types.StringType,
				"log_label_string": types.StringType,
			},
		}, []logLabelSettingModel{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.LogLabelSettings = emptyList
		}
	}

	// SystemName and ProjectCreationConfig are config-only (not returned by API)
	plan.SystemName = config.SystemName
	plan.ProjectCreationConfig = config.ProjectCreationConfig

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Project created successfully", map[string]any{"project_name": plan.ProjectName.ValueString()})
}

// Read refreshes the Terraform state with the latest data.
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading project", map[string]any{"project_name": state.ProjectName.ValueString()})

	// Get the project from the API
	project, err := r.client.GetProject(state.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			"Could not read project, unexpected error: "+err.Error(),
		)
		return
	}

	// If project doesn't exist, remove from state
	if project == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state with API data
	// The project.Settings map contains all the configuration from the API
	settings := project.Settings
	if settings == nil {
		settings = make(map[string]interface{})
	}

	// Helper function to safely get values from settings map
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
			// If it's already a string, return it
			if str, ok := val.(string); ok {
				return types.StringValue(str)
			}
			// Otherwise, marshal it to JSON
			if jsonBytes, err := json.Marshal(val); err == nil {
				return types.StringValue(string(jsonBytes))
			}
		}
		return types.StringNull()
	}

	// Populate all fields from API response
	state.ProjectDisplayName = getString("projectDisplayName")
	state.CValue = getInt64("cValue")
	state.PValue = getFloat64("pValue")
	state.ProjectTimeZone = getString("projectTimeZone")
	state.SamplingInterval = getInt64("samplingInterval")

	// Basic Configuration
	state.UBLRetentionTime = getInt64("UBLRetentionTime")
	state.AlertAverageTime = getInt64("alertAverageTime")
	state.AlertHourlyCost = getFloat64("alertHourlyCost")
	state.AnomalyDetectionMode = getInt64("anomalyDetectionMode")
	state.AnomalySamplingInterval = getInt64("anomalySamplingInterval")
	state.AvgPerIncidentDowntimeCost = getFloat64("avgPerIncidentDowntimeCost")
	state.CausalPredictionSetting = getInt64("causalPredictionSetting")
	state.CausalMinDelay = getString("causalMinDelay")
	state.ColdEventThreshold = getInt64("coldEventThreshold")
	state.ColdNumberLimit = getInt64("coldNumberLimit")
	state.CollectAllRareEventsFlag = getBool("collectAllRareEventsFlag")
	state.DailyModelSpan = getInt64("dailyModelSpan")
	state.DisableLogCompressEvent = getBool("disableLogCompressEvent")
	state.DisableModelKeywordStatsCollection = getBool("disableModelKeywordStatsCollection")

	// Anomaly and Detection Settings
	state.EnableAnomalyScoreEscalation = getBool("enableAnomalyScoreEscalation")
	state.EnableHotEvent = getBool("enableHotEvent")
	state.EnableNewAlertEmail = getBool("enableNewAlertEmail")
	state.EnableStreamDetection = getBool("enableStreamDetection")
	state.EscalationAnomalyScoreThreshold = getString("escalationAnomalyScoreThreshold")
	state.FeatureOutlierSensitivity = getString("featureOutlierSensitivity")
	state.FeatureOutlierThreshold = getFloat64("featureOutlierThreshold")
	state.HotEventCalmDownPeriod = getInt64("hotEventCalmDownPeriod")
	state.HotEventDetectionMode = getInt64("hotEventDetectionMode")
	state.HotEventThreshold = getInt64("hotEventThreshold")
	state.HotNumberLimit = getInt64("hotNumberLimit")
	state.IgnoreAnomalyScoreThreshold = getString("ignoreAnomalyScoreThreshold")
	state.IgnoreInstanceForKB = getBool("ignoreInstanceForKB")

	// Incident Settings
	state.IncidentPredictionEventLimit = getInt64("incidentPredictionEventLimit")
	state.IncidentPredictionWindow = getInt64("incidentPredictionWindow")
	state.IncidentRelationSearchWindow = getInt64("incidentRelationSearchWindow")

	// Instance Settings
	state.InstanceConvertFlag = getBool("instanceConvertFlag")
	state.InstanceDownEnable = getBool("instanceDownEnable")
	state.IsEdgeBrain = getBool("isEdgeBrain")
	state.IsGroupingByInstance = getBool("isGroupingByInstance")
	state.IsTracePrompt = getBool("isTracePrompt")
	state.ShowInstanceDown = getBool("showInstanceDown")

	// Log Settings
	state.KeywordFeatureNumber = getInt64("keywordFeatureNumber")
	state.KeywordSetting = getInt64("keywordSetting")
	state.LargeProject = getBool("largeProject")
	state.LogAnomalyEventBaseScore = getString("logAnomalyEventBaseScore")
	state.LogDetectionMinCount = getInt64("logDetectionMinCount")
	state.LogDetectionSize = getInt64("logDetectionSize")
	state.LogPatternLimitLevel = getInt64("logPatternLimitLevel")
	state.MaxLogModelSize = getInt64("maxLogModelSize")
	state.MaximumDetectionWaitTime = getInt64("maximumDetectionWaitTime")
	state.MaximumThreads = getInt64("maximumThreads")
	state.ModelKeywordSetting = getInt64("modelKeywordSetting")
	state.MultiLineFlag = getBool("multiLineFlag")
	state.NlpFlag = getBool("nlpFlag")
	state.PrettyJsonConvertorFlag = getBool("prettyJsonConvertorFlag")

	// Model Settings
	state.MaximumRootCauseResultSize = getInt64("maximumRootCauseResultSize")
	state.MinIncidentPredictionWindow = getInt64("minIncidentPredictionWindow")
	state.MinValidModelSpan = getInt64("minValidModelSpan")
	state.MultiHopSearchLevel = getInt64("multiHopSearchLevel")
	state.MultiHopSearchLimit = getString("multiHopSearchLimit")

	// Pattern and Event Settings
	state.NewAlertFlag = getBool("newAlertFlag")
	state.NewPatternNumberLimit = getInt64("newPatternNumberLimit")
	state.NewPatternRange = getInt64("newPatternRange")
	state.NormalEventCausalFlag = getBool("normalEventCausalFlag")

	// Prediction Settings
	state.PredictionCountThreshold = getInt64("predictionCountThreshold")
	state.PredictionProbabilityThreshold = getFloat64("predictionProbabilityThreshold")
	state.PredictionRuleActiveCondition = getInt64("predictionRuleActiveCondition")
	state.PredictionRuleActiveThreshold = getFloat64("predictionRuleActiveThreshold")
	state.PredictionRuleFalsePositiveThreshold = getInt64("predictionRuleFalsePositiveThreshold")
	state.PredictionRuleInactiveThreshold = getFloat64("predictionRuleInactiveThreshold")
	state.ProjectModelFlag = getBool("projectModelFlag")
	state.Proxy = getString("proxy")

	// Rare Event Settings
	state.RareAnomalyType = getInt64("rareAnomalyType")
	state.RareEventAlertThresholds = getInt64("rareEventAlertThresholds")
	state.RareNumberLimit = getInt64("rareNumberLimit")
	state.RetentionTime = getInt64("retentionTime")

	// Root Cause Settings
	state.RootCauseCountThreshold = getInt64("rootCauseCountThreshold")
	state.RootCauseLogMessageSearchRange = getInt64("rootCauseLogMessageSearchRange")
	state.RootCauseProbabilityThreshold = getFloat64("rootCauseProbabilityThreshold")
	state.RootCauseRankSetting = getInt64("rootCauseRankSetting")

	// Similarity and Training
	state.SimilaritySensitivity = getString("similaritySensitivity")
	state.TrainingFilter = getBool("trainingFilter")

	// Webhook Settings
	state.MaxWebHookRequestSize = getInt64("maxWebHookRequestSize")
	state.WebhookAlertDampening = getInt64("webhookAlertDampening")
	state.WebhookBlackListSetStr = getString("webhookBlackListSetStr")
	state.WebhookCriticalKeywordSetStr = getString("webhookCriticalKeywordSetStr")
	state.WebhookTypeSetStr = getString("webhookTypeSetStr")
	state.WebhookUrl = getString("webhookUrl")
	state.WhitelistNumberLimit = getInt64("whitelistNumberLimit")
	state.ZoneNameKey = getString("zoneNameKey")

	// Metric Project Fields

	// JSON String Fields
	state.BaseValueSetting = getJSONString("baseValueSetting")
	state.CdfSetting = getJSONString("cdfSetting")
	state.EmailSetting = getJSONString("emailSetting")
	state.InstanceGroupingUpdate = getJSONString("instanceGroupingUpdate")
	state.LlmEvaluationSetting = getJSONString("llmEvaluationSetting")
	state.LogToLogSettingList = getJSONString("logToLogSettingList")
	state.WebhookHeaderList = getJSONString("webhookHeaderList")
	state.SharedUsernames = getJSONString("sharedUsernames")

	// Read log label settings from API
	logLabels, err := r.client.GetLogLabels(state.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		tflog.Warn(ctx, "Could not read log labels", map[string]any{"error": err.Error()})
		// Keep existing state if we can't read from API
	} else if logLabels != nil {
		// Extract existing state for comparison
		var existingSettings []logLabelSettingModel
		if !state.LogLabelSettings.IsNull() && !state.LogLabelSettings.IsUnknown() {
			diags = state.LogLabelSettings.ElementsAs(ctx, &existingSettings, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		// Convert API response to state model, preserving the order from existing state
		convertedSettings := convertLogLabelsToState(logLabels, existingSettings)

		// Convert back to types.List
		listValue, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"label_type":       types.StringType,
				"log_label_string": types.StringType,
			},
		}, convertedSettings)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.LogLabelSettings = listValue
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Project read successfully", map[string]any{"project_name": state.ProjectName.ValueString()})
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get the plan (desired state after update)
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the config (only user-specified values, not computed ones)
	var config projectResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating project", map[string]any{"project_name": config.ProjectName.ValueString()})

	// Use config (not plan) to populate settings - this ensures we only send user-specified values
	projectConfig := &client.ProjectConfig{
		ProjectName:        config.ProjectName.ValueString(),
		ProjectDisplayName: config.ProjectDisplayName.ValueString(),
		SystemName:         config.SystemName.ValueString(),
		CValue:             int(config.CValue.ValueInt64()),
		PValue:             config.PValue.ValueFloat64(),
		Settings:           populateSettings(&config),
	}

	err := r.client.UpdateProject(projectConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			"Could not update project, unexpected error: "+err.Error(),
		)
		return
	}

	// After successful update, read back the actual state from API
	project, err := r.client.GetProject(plan.ProjectName.ValueString(), r.client.Username)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading project after update",
			"Could not read project after update: "+err.Error()+". State may be out of sync.",
		)
	} else if project != nil {
		// Populate state with actual API values (same logic as Read method)
		settings := project.Settings
		if settings == nil {
			settings = make(map[string]interface{})
		}

		// Helper functions (same as in Read)
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
				// If it's already a string, return it
				if str, ok := val.(string); ok {
					return types.StringValue(str)
				}
				// Otherwise, marshal it to JSON
				if jsonBytes, err := json.Marshal(val); err == nil {
					return types.StringValue(string(jsonBytes))
				}
			}
			return types.StringNull()
		}

		// Populate all fields from API response
		plan.ProjectDisplayName = getString("projectDisplayName")
		plan.CValue = getInt64("cValue")
		plan.PValue = getFloat64("pValue")
		plan.ProjectTimeZone = getString("projectTimeZone")
		plan.SamplingInterval = getInt64("samplingInterval")

		// Basic Configuration
		plan.UBLRetentionTime = getInt64("UBLRetentionTime")
		plan.AlertAverageTime = getInt64("alertAverageTime")
		plan.AlertHourlyCost = getFloat64("alertHourlyCost")
		plan.AnomalyDetectionMode = getInt64("anomalyDetectionMode")
		plan.AnomalySamplingInterval = getInt64("anomalySamplingInterval")
		plan.AvgPerIncidentDowntimeCost = getFloat64("avgPerIncidentDowntimeCost")
		plan.CausalPredictionSetting = getInt64("causalPredictionSetting")
		plan.CausalMinDelay = getString("causalMinDelay")
		plan.ColdEventThreshold = getInt64("coldEventThreshold")
		plan.ColdNumberLimit = getInt64("coldNumberLimit")
		plan.CollectAllRareEventsFlag = getBool("collectAllRareEventsFlag")
		plan.DailyModelSpan = getInt64("dailyModelSpan")
		plan.DisableLogCompressEvent = getBool("disableLogCompressEvent")
		plan.DisableModelKeywordStatsCollection = getBool("disableModelKeywordStatsCollection")

		// Anomaly and Detection Settings
		plan.EnableAnomalyScoreEscalation = getBool("enableAnomalyScoreEscalation")
		plan.EnableHotEvent = getBool("enableHotEvent")
		plan.EnableNewAlertEmail = getBool("enableNewAlertEmail")
		plan.EnableStreamDetection = getBool("enableStreamDetection")
		plan.EscalationAnomalyScoreThreshold = getString("escalationAnomalyScoreThreshold")
		plan.FeatureOutlierSensitivity = getString("featureOutlierSensitivity")
		plan.FeatureOutlierThreshold = getFloat64("featureOutlierThreshold")
		plan.HotEventCalmDownPeriod = getInt64("hotEventCalmDownPeriod")
		plan.HotEventDetectionMode = getInt64("hotEventDetectionMode")
		plan.HotEventThreshold = getInt64("hotEventThreshold")
		plan.HotNumberLimit = getInt64("hotNumberLimit")
		plan.IgnoreAnomalyScoreThreshold = getString("ignoreAnomalyScoreThreshold")
		plan.IgnoreInstanceForKB = getBool("ignoreInstanceForKB")

		// Incident Settings
		plan.IncidentPredictionEventLimit = getInt64("incidentPredictionEventLimit")
		plan.IncidentPredictionWindow = getInt64("incidentPredictionWindow")
		plan.IncidentRelationSearchWindow = getInt64("incidentRelationSearchWindow")

		// Instance Settings
		plan.InstanceConvertFlag = getBool("instanceConvertFlag")
		plan.InstanceDownEnable = getBool("instanceDownEnable")
		plan.IsEdgeBrain = getBool("isEdgeBrain")
		plan.IsGroupingByInstance = getBool("isGroupingByInstance")
		plan.IsTracePrompt = getBool("isTracePrompt")
		plan.ShowInstanceDown = getBool("showInstanceDown")

		// Log Settings
		plan.KeywordFeatureNumber = getInt64("keywordFeatureNumber")
		plan.KeywordSetting = getInt64("keywordSetting")
		plan.LargeProject = getBool("largeProject")
		plan.LogAnomalyEventBaseScore = getString("logAnomalyEventBaseScore")
		plan.LogDetectionMinCount = getInt64("logDetectionMinCount")
		plan.LogDetectionSize = getInt64("logDetectionSize")
		plan.LogPatternLimitLevel = getInt64("logPatternLimitLevel")
		plan.MaxLogModelSize = getInt64("maxLogModelSize")
		plan.MaximumDetectionWaitTime = getInt64("maximumDetectionWaitTime")
		plan.MaximumThreads = getInt64("maximumThreads")
		plan.ModelKeywordSetting = getInt64("modelKeywordSetting")
		plan.MultiLineFlag = getBool("multiLineFlag")
		plan.NlpFlag = getBool("nlpFlag")
		plan.PrettyJsonConvertorFlag = getBool("prettyJsonConvertorFlag")

		// Model Settings
		plan.MaximumRootCauseResultSize = getInt64("maximumRootCauseResultSize")
		plan.MinIncidentPredictionWindow = getInt64("minIncidentPredictionWindow")
		plan.MinValidModelSpan = getInt64("minValidModelSpan")
		plan.MultiHopSearchLevel = getInt64("multiHopSearchLevel")
		plan.MultiHopSearchLimit = getString("multiHopSearchLimit")

		// Pattern and Event Settings
		plan.NewAlertFlag = getBool("newAlertFlag")
		plan.NewPatternNumberLimit = getInt64("newPatternNumberLimit")
		plan.NewPatternRange = getInt64("newPatternRange")
		plan.NormalEventCausalFlag = getBool("normalEventCausalFlag")

		// Prediction Settings
		plan.PredictionCountThreshold = getInt64("predictionCountThreshold")
		plan.PredictionProbabilityThreshold = getFloat64("predictionProbabilityThreshold")
		plan.PredictionRuleActiveCondition = getInt64("predictionRuleActiveCondition")
		plan.PredictionRuleActiveThreshold = getFloat64("predictionRuleActiveThreshold")
		plan.PredictionRuleFalsePositiveThreshold = getInt64("predictionRuleFalsePositiveThreshold")
		plan.PredictionRuleInactiveThreshold = getFloat64("predictionRuleInactiveThreshold")
		plan.ProjectModelFlag = getBool("projectModelFlag")
		plan.Proxy = getString("proxy")

		// Rare Event Settings
		plan.RareAnomalyType = getInt64("rareAnomalyType")
		plan.RareEventAlertThresholds = getInt64("rareEventAlertThresholds")
		plan.RareNumberLimit = getInt64("rareNumberLimit")
		plan.RetentionTime = getInt64("retentionTime")

		// Root Cause Settings
		plan.RootCauseCountThreshold = getInt64("rootCauseCountThreshold")
		plan.RootCauseLogMessageSearchRange = getInt64("rootCauseLogMessageSearchRange")
		plan.RootCauseProbabilityThreshold = getFloat64("rootCauseProbabilityThreshold")
		plan.RootCauseRankSetting = getInt64("rootCauseRankSetting")

		// Similarity and Training
		plan.SimilaritySensitivity = getString("similaritySensitivity")
		plan.TrainingFilter = getBool("trainingFilter")

		// Webhook Settings
		plan.MaxWebHookRequestSize = getInt64("maxWebHookRequestSize")
		plan.WebhookAlertDampening = getInt64("webhookAlertDampening")
		plan.WebhookBlackListSetStr = getString("webhookBlackListSetStr")
		plan.WebhookCriticalKeywordSetStr = getString("webhookCriticalKeywordSetStr")
		plan.WebhookTypeSetStr = getString("webhookTypeSetStr")
		plan.WebhookUrl = getString("webhookUrl")
		plan.WhitelistNumberLimit = getInt64("whitelistNumberLimit")
		plan.ZoneNameKey = getString("zoneNameKey")

		// Metric Project Fields

		// JSON String Fields
		plan.BaseValueSetting = getJSONString("baseValueSetting")
		plan.CdfSetting = getJSONString("cdfSetting")
		plan.EmailSetting = getJSONString("emailSetting")
		plan.InstanceGroupingUpdate = getJSONString("instanceGroupingUpdate")
		plan.LlmEvaluationSetting = getJSONString("llmEvaluationSetting")
		plan.LogToLogSettingList = getJSONString("logToLogSettingList")
		plan.WebhookHeaderList = getJSONString("webhookHeaderList")
		plan.SharedUsernames = getJSONString("sharedUsernames")
	}

	// Process log_label_settings if provided - each setting must be applied individually
	if !config.LogLabelSettings.IsNull() && !config.LogLabelSettings.IsUnknown() {
		var configSettings []logLabelSettingModel
		diags = config.LogLabelSettings.ElementsAs(ctx, &configSettings, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(configSettings) > 0 {
			tflog.Info(ctx, "Processing log label settings", map[string]any{"count": len(configSettings)})

			// Convert terraform model to client model
			settings := make([]*client.LogLabelSetting, 0, len(configSettings))
			for _, setting := range configSettings {
				settings = append(settings, &client.LogLabelSetting{
					LabelType:      setting.LabelType.ValueString(),
					LogLabelString: setting.LogLabelString.ValueString(),
				})
			}

			// Apply all settings (function will iterate and call API for each)
			err := r.client.CreateOrUpdateLogLabels(
				plan.ProjectName.ValueString(),
				r.client.Username,
				settings,
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error applying log label settings",
					fmt.Sprintf("Could not apply log label settings: %s", err.Error()),
				)
				return
			}

			// Log labels applied successfully - use config values
			plan.LogLabelSettings = config.LogLabelSettings
		}
	} else if config.LogLabelSettings.IsNull() {
		// Explicitly set empty list if none configured
		emptyList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"label_type":       types.StringType,
				"log_label_string": types.StringType,
			},
		}, []logLabelSettingModel{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.LogLabelSettings = emptyList
		}
	}

	// Preserve config-only fields that don't come from API
	plan.SystemName = config.SystemName
	plan.ProjectCreationConfig = config.ProjectCreationConfig

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Project updated successfully", map[string]any{"project_name": plan.ProjectName.ValueString()})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting project", map[string]any{"project_name": state.ProjectName.ValueString()})

	err := r.client.DeleteProject(state.ProjectName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			"Could not delete project, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Project deleted successfully", map[string]any{"project_name": state.ProjectName.ValueString()})
}

// ImportState imports the resource into Terraform state.
func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID (project name) as the identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// normalizeJSON parses and re-marshals JSON to normalize formatting
// This ensures that semantically equivalent JSON strings are byte-for-byte identical
// Uses the same format as Terraform's jsonencode(): compact with HTML escaping
func normalizeJSON(jsonStr string) string {
	if jsonStr == "" || jsonStr == "[]" {
		return jsonStr
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// If parsing fails, return as-is
		return jsonStr
	}

	// Re-marshal in compact form with HTML escaping enabled (matches Terraform's jsonencode)
	normalized, err := json.Marshal(data)
	if err != nil {
		return jsonStr
	}

	return string(normalized)
}

// convertLogLabelsToState converts API log labels response to Terraform state model
// while preserving the order from existing state when possible
func convertLogLabelsToState(apiLabels map[string]string, existingState []logLabelSettingModel) []logLabelSettingModel {
	var result []logLabelSettingModel

	// Map from API field names to label types
	apiFieldToLabelType := map[string]string{
		"whitelist":                "whitelist",
		"trainingWhitelist":        "trainingWhitelist",
		"trainingBlacklistLabels":  "blacklist",
		"featurelist":              "featurelist",
		"incidentlist":             "incidentlist",
		"triagelist":               "triagelist",
		"patternNameLabels":        "patternName",
		"patternSignatureLabels":   "patternSignature",
		"patternMatchRegexLabels":  "patternMatchRegex",
		"patternIgnoreRegexLabels": "patternIgnoreRegex",
		"customActionLabels":       "customAction",
		"logEventIDLabels":         "logEventID",
		"logSeverityLabels":        "logSeverity",
		"logStatusCodeLabels":      "logStatusCode",
		"alertEventTypeLabels":     "alertEventType",
		"anomalyFeatureLabels":     "anomalyFeature",
		"dataFilterLabels":         "dataFilter",
		"instanceNameLabels":       "instanceName",
		"dataQualityCheckLabels":   "dataQualityCheck",
		"extractionBlacklist":      "extractionBlacklist",
	}

	// Reverse map for looking up API fields from label types
	labelTypeToAPIField := make(map[string]string)
	for apiField, labelType := range apiFieldToLabelType {
		labelTypeToAPIField[labelType] = apiField
	}

	// Create a map of API data for quick lookup
	apiDataMap := make(map[string]string)
	for apiField, jsonString := range apiLabels {
		if labelType, ok := apiFieldToLabelType[apiField]; ok {
			if jsonString != "" && jsonString != "[]" {
				apiDataMap[labelType] = normalizeJSON(jsonString)
			}
		}
	}

	// If we have existing state, preserve its order and only include items that exist in API
	if len(existingState) > 0 {
		processedTypes := make(map[string]bool)

		// First pass: preserve order from existing state
		for _, existing := range existingState {
			labelType := existing.LabelType.ValueString()
			if normalizedJSON, ok := apiDataMap[labelType]; ok {
				result = append(result, logLabelSettingModel{
					LabelType:      types.StringValue(labelType),
					LogLabelString: types.StringValue(normalizedJSON),
				})
				processedTypes[labelType] = true
			}
		}

		// Second pass: add any new types from API that weren't in existing state
		// Use a consistent order for new types
		defaultOrder := []string{
			"trainingWhitelist",
			"featurelist",
			"incidentlist",
			"triagelist",
			"patternName",
			"whitelist",
			"blacklist",
			"patternSignature",
			"patternMatchRegex",
			"patternIgnoreRegex",
			"customAction",
			"logEventID",
			"logSeverity",
			"logStatusCode",
			"alertEventType",
			"anomalyFeature",
			"dataFilter",
			"instanceName",
			"dataQualityCheck",
			"extractionBlacklist",
		}

		for _, labelType := range defaultOrder {
			if !processedTypes[labelType] {
				if normalizedJSON, ok := apiDataMap[labelType]; ok {
					result = append(result, logLabelSettingModel{
						LabelType:      types.StringValue(labelType),
						LogLabelString: types.StringValue(normalizedJSON),
					})
				}
			}
		}
	} else {
		// No existing state, use default order
		labelTypeOrder := []string{
			"trainingWhitelist",
			"featurelist",
			"incidentlist",
			"triagelist",
			"patternName",
			"whitelist",
			"blacklist",
			"patternSignature",
			"patternMatchRegex",
			"patternIgnoreRegex",
			"customAction",
			"logEventID",
			"logSeverity",
			"logStatusCode",
			"alertEventType",
			"anomalyFeature",
			"dataFilter",
			"instanceName",
			"dataQualityCheck",
			"extractionBlacklist",
		}

		for _, labelType := range labelTypeOrder {
			if normalizedJSON, ok := apiDataMap[labelType]; ok {
				result = append(result, logLabelSettingModel{
					LabelType:      types.StringValue(labelType),
					LogLabelString: types.StringValue(normalizedJSON),
				})
			}
		}
	}

	return result
}
