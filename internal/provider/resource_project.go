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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	ComponentNameAutoOverwrite types.Bool `tfsdk:"component_name_auto_overwrite"`
	InstanceConvertFlag        types.Bool `tfsdk:"instance_convert_flag"`
	InstanceDownEnable         types.Bool `tfsdk:"instance_down_enable"`
	IsEdgeBrain                types.Bool `tfsdk:"is_edge_brain"`
	IsGroupingByInstance       types.Bool `tfsdk:"is_grouping_by_instance"`
	IsTracePrompt              types.Bool `tfsdk:"is_trace_prompt"`
	ShowInstanceDown           types.Bool `tfsdk:"show_instance_down"`

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
	BaseValueSetting          types.String `tfsdk:"base_value_setting"`
	CdfSetting                types.String `tfsdk:"cdf_setting"`
	EmailSetting              types.String `tfsdk:"email_setting"`
	InstanceGroupingUpdate    types.String `tfsdk:"instance_grouping_update"`
	LlmEvaluationSetting      types.String `tfsdk:"llm_evaluation_setting"`
	LogToLogSettingList       types.String `tfsdk:"log_to_log_setting_list"`
	WebhookHeaderList         types.String `tfsdk:"webhook_header_list"`
	SharedUsernames           types.String `tfsdk:"shared_usernames"`
	LogLabelSettings          types.Set    `tfsdk:"log_label_settings"`
	JsonKeySettings           types.Set    `tfsdk:"json_key_settings"`
	ProjectServiceNowSettings types.Object `tfsdk:"project_servicenow_settings"`
	HolidaySettings           types.Set    `tfsdk:"holiday_settings"`
}

type holidaySettingModel struct {
	Name      types.String `tfsdk:"name"`
	StartDate types.String `tfsdk:"start_date"`
	EndDate   types.String `tfsdk:"end_date"`
}

type jsonKeySettingModel struct {
	JsonKey               types.String `tfsdk:"json_key"`
	Type                  types.String `tfsdk:"type"`
	SummarySetting        types.Bool   `tfsdk:"summary_setting"`
	MetafieldSetting      types.Bool   `tfsdk:"metafield_setting"`
	DampeningfieldSetting types.Bool   `tfsdk:"dampening_field_setting"`
}

type projectCreationConfigModel struct {
	DataType         types.String `tfsdk:"data_type"`
	InstanceType     types.String `tfsdk:"instance_type"`
	ProjectCloudType types.String `tfsdk:"project_cloud_type"`
	InsightAgentType types.String `tfsdk:"insight_agent_type"`
	ServiceNowTable  types.String `tfsdk:"servicenow_table"`
}

type projectServiceNowSettingsModel struct {
	Host               types.String `tfsdk:"host"`
	SysparmQuery       types.String `tfsdk:"sysparm_query"`
	Proxy              types.String `tfsdk:"proxy"`
	ServiceNowUser     types.String `tfsdk:"servicenow_user"`
	ServiceNowPassword types.String `tfsdk:"servicenow_password"`
	InstanceField      types.String `tfsdk:"instance_field"`
	InstanceFieldRegex types.String `tfsdk:"instance_field_regex"`
	TimestampFormat    types.String `tfsdk:"timestamp_format"`
	ClientID           types.String `tfsdk:"client_id"`
	ClientSecret       types.String `tfsdk:"client_secret"`
	AdditionalFields   types.List   `tfsdk:"additional_fields"`
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
					"servicenow_table": schema.StringAttribute{
						Description: "The ServiceNow table name. Required when project_cloud_type is ServiceNow.",
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
			"component_name_auto_overwrite": schema.BoolAttribute{
				Description: "Enable automatic overwrite of component names",
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
			"log_label_settings": schema.SetNestedAttribute{
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
			"project_servicenow_settings": schema.SingleNestedAttribute{
				Description: "ServiceNow third-party settings for the project. Only applies when project_cloud_type is 'ServiceNow' (case insensitive).",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Description: "ServiceNow instance host URL",
						Required:    true,
					},
					"sysparm_query": schema.StringAttribute{
						Description: "ServiceNow query parameter",
						Optional:    true,
						Computed:    true,
					},
					"proxy": schema.StringAttribute{
						Description: "Proxy URL for ServiceNow connection",
						Optional:    true,
						Computed:    true,
					},
					"servicenow_user": schema.StringAttribute{
						Description: "ServiceNow username for authentication",
						Required:    true,
					},
					"servicenow_password": schema.StringAttribute{
						Description: "ServiceNow password for authentication",
						Required:    true,
						Sensitive:   true,
					},
					"instance_field": schema.StringAttribute{
						Description: "Field to use for instance identification",
						Optional:    true,
						Computed:    true,
					},
					"instance_field_regex": schema.StringAttribute{
						Description: "Regex pattern for instance field",
						Optional:    true,
						Computed:    true,
					},
					"timestamp_format": schema.StringAttribute{
						Description: "Timestamp format for ServiceNow data",
						Optional:    true,
						Computed:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "OAuth client ID for ServiceNow",
						Optional:    true,
						Computed:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "OAuth client secret for ServiceNow",
						Optional:    true,
						Computed:    true,
						Sensitive:   true,
					},
					"additional_fields": schema.ListAttribute{
						Description: "Additional fields to fetch from ServiceNow",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"holiday_settings": schema.SetNestedAttribute{
				Description: "List of holiday settings for the project. Each holiday has a name, start date, and end date (MM-DD format).",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the holiday",
							Required:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Start date of the holiday in MM-DD format (e.g., '12-25')",
							Required:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "End date of the holiday in MM-DD format (e.g., '12-26')",
							Required:    true,
						},
					},
				},
			},
			"json_key_settings": schema.SetNestedAttribute{
				Description: "Set of JSON key settings for the project. Manages custom JSON fields extracted from logs.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"json_key": schema.StringAttribute{
							Description: "The JSON key name to extract from logs",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "The data type of the JSON value (e.g., 'string', 'number', 'JSONArray')",
							Required:    true,
						},
						"summary_setting": schema.BoolAttribute{
							Description: "Whether to include this key in the summary statistics",
							Required:    true,
						},
						"metafield_setting": schema.BoolAttribute{
							Description: "Whether to include this key in the metafield settings",
							Required:    true,
						},
						"dampening_field_setting": schema.BoolAttribute{
							Description: "Whether to include this key in the dampening field settings",
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
	if !plan.MinValidModelSpan.IsNull() {
		s["minValidModelSpan"] = int(plan.MinValidModelSpan.ValueInt64())
	}
	if !plan.MaxWebHookRequestSize.IsNull() {
		s["maxWebHookRequestSize"] = int(plan.MaxWebHookRequestSize.ValueInt64())
	}
	if !plan.WebhookUrl.IsNull() {
		s["webhookUrl"] = plan.WebhookUrl.ValueString()
	}
	if !plan.WebhookTypeSetStr.IsNull() {
		s["webhookTypeSetStr"] = plan.WebhookTypeSetStr.ValueString()
	}
	if !plan.WebhookBlackListSetStr.IsNull() {
		s["webhookBlackListSetStr"] = plan.WebhookBlackListSetStr.ValueString()
	}
	if !plan.WebhookCriticalKeywordSetStr.IsNull() {
		s["webhookCriticalKeywordSetStr"] = plan.WebhookCriticalKeywordSetStr.ValueString()
	}
	if !plan.WebhookAlertDampening.IsNull() {
		s["webhookAlertDampening"] = int(plan.WebhookAlertDampening.ValueInt64())
	}
	if !plan.Proxy.IsNull() {
		s["proxy"] = plan.Proxy.ValueString()
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
	if !plan.MultiHopSearchLimit.IsNull() {
		s["multiHopSearchLimit"] = plan.MultiHopSearchLimit.ValueString()
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

	// Log-specific fields
	if !plan.DailyModelSpan.IsNull() {
		s["dailyModelSpan"] = int(plan.DailyModelSpan.ValueInt64())
	}
	if !plan.KeywordFeatureNumber.IsNull() {
		s["keywordFeatureNumber"] = int(plan.KeywordFeatureNumber.ValueInt64())
	}
	if !plan.MaxLogModelSize.IsNull() {
		s["maxLogModelSize"] = int(plan.MaxLogModelSize.ValueInt64())
	}
	if !plan.ModelKeywordSetting.IsNull() {
		s["modelKeywordSetting"] = int(plan.ModelKeywordSetting.ValueInt64())
	}
	if !plan.NlpFlag.IsNull() {
		s["nlpFlag"] = plan.NlpFlag.ValueBool()
	}
	if !plan.ProjectModelFlag.IsNull() {
		s["projectModelFlag"] = plan.ProjectModelFlag.ValueBool()
	}
	if !plan.MaximumThreads.IsNull() {
		s["maximumThreads"] = int(plan.MaximumThreads.ValueInt64())
	}
	if !plan.LogDetectionMinCount.IsNull() {
		s["logDetectionMinCount"] = int(plan.LogDetectionMinCount.ValueInt64())
	}
	if !plan.LogDetectionSize.IsNull() {
		s["logDetectionSize"] = int(plan.LogDetectionSize.ValueInt64())
	}
	if !plan.MaximumDetectionWaitTime.IsNull() {
		s["maximumDetectionWaitTime"] = int(plan.MaximumDetectionWaitTime.ValueInt64())
	}
	if !plan.KeywordSetting.IsNull() {
		s["keywordSetting"] = int(plan.KeywordSetting.ValueInt64())
	}
	if !plan.LogPatternLimitLevel.IsNull() {
		s["logPatternLimitLevel"] = int(plan.LogPatternLimitLevel.ValueInt64())
	}
	if !plan.NormalEventCausalFlag.IsNull() {
		s["normalEventCausalFlag"] = plan.NormalEventCausalFlag.ValueBool()
	}
	if !plan.SimilaritySensitivity.IsNull() {
		s["similaritySensitivity"] = plan.SimilaritySensitivity.ValueString()
	}
	if !plan.CollectAllRareEventsFlag.IsNull() {
		s["collectAllRareEventsFlag"] = plan.CollectAllRareEventsFlag.ValueBool()
	}
	if !plan.RareEventAlertThresholds.IsNull() {
		s["rareEventAlertThresholds"] = int(plan.RareEventAlertThresholds.ValueInt64())
	}
	if !plan.LogAnomalyEventBaseScore.IsNull() {
		s["logAnomalyEventBaseScore"] = plan.LogAnomalyEventBaseScore.ValueString()
	}
	if !plan.RareNumberLimit.IsNull() {
		s["rareNumberLimit"] = int(plan.RareNumberLimit.ValueInt64())
	}
	if !plan.WhitelistNumberLimit.IsNull() {
		s["whitelistNumberLimit"] = int(plan.WhitelistNumberLimit.ValueInt64())
	}
	if !plan.NewPatternNumberLimit.IsNull() {
		s["newPatternNumberLimit"] = int(plan.NewPatternNumberLimit.ValueInt64())
	}
	if !plan.HotNumberLimit.IsNull() {
		s["hotNumberLimit"] = int(plan.HotNumberLimit.ValueInt64())
	}
	if !plan.ColdNumberLimit.IsNull() {
		s["coldNumberLimit"] = int(plan.ColdNumberLimit.ValueInt64())
	}
	if !plan.RareAnomalyType.IsNull() {
		s["rareAnomalyType"] = int(plan.RareAnomalyType.ValueInt64())
	}
	if !plan.HotEventThreshold.IsNull() {
		s["hotEventThreshold"] = int(plan.HotEventThreshold.ValueInt64())
	}
	if !plan.ColdEventThreshold.IsNull() {
		s["coldEventThreshold"] = int(plan.ColdEventThreshold.ValueInt64())
	}
	if !plan.DisableLogCompressEvent.IsNull() {
		s["disableLogCompressEvent"] = plan.DisableLogCompressEvent.ValueBool()
	}
	if !plan.EnableHotEvent.IsNull() {
		s["enableHotEvent"] = plan.EnableHotEvent.ValueBool()
	}
	if !plan.HotEventCalmDownPeriod.IsNull() {
		s["hotEventCalmDownPeriod"] = int(plan.HotEventCalmDownPeriod.ValueInt64())
	}
	if !plan.InstanceDownEnable.IsNull() {
		s["instanceDownEnable"] = plan.InstanceDownEnable.ValueBool()
	}
	if !plan.AnomalySamplingInterval.IsNull() {
		s["anomalySamplingInterval"] = int(plan.AnomalySamplingInterval.ValueInt64())
	}
	if !plan.HotEventDetectionMode.IsNull() {
		s["hotEventDetectionMode"] = int(plan.HotEventDetectionMode.ValueInt64())
	}
	if !plan.AnomalyDetectionMode.IsNull() {
		s["anomalyDetectionMode"] = int(plan.AnomalyDetectionMode.ValueInt64())
	}
	if !plan.PrettyJsonConvertorFlag.IsNull() {
		s["prettyJsonConvertorFlag"] = plan.PrettyJsonConvertorFlag.ValueBool()
	}
	if !plan.ZoneNameKey.IsNull() {
		s["zoneNameKey"] = plan.ZoneNameKey.ValueString()
	}
	if !plan.MultiLineFlag.IsNull() {
		s["multiLineFlag"] = plan.MultiLineFlag.ValueBool()
	}
	if !plan.FeatureOutlierSensitivity.IsNull() {
		s["featureOutlierSensitivity"] = plan.FeatureOutlierSensitivity.ValueString()
	}
	if !plan.DisableModelKeywordStatsCollection.IsNull() {
		s["disableModelKeywordStatsCollection"] = plan.DisableModelKeywordStatsCollection.ValueBool()
	}
	if !plan.ComponentNameAutoOverwrite.IsNull() {
		s["componentNameAutoOverwrite"] = plan.ComponentNameAutoOverwrite.ValueBool()
	}
	if !plan.InstanceConvertFlag.IsNull() {
		s["instanceConvertFlag"] = plan.InstanceConvertFlag.ValueBool()
	}
	if !plan.NewAlertFlag.IsNull() {
		s["newAlertFlag"] = plan.NewAlertFlag.ValueBool()
	}
	if !plan.IsGroupingByInstance.IsNull() {
		s["isGroupingByInstance"] = plan.IsGroupingByInstance.ValueBool()
	}
	if !plan.FeatureOutlierThreshold.IsNull() {
		s["featureOutlierThreshold"] = plan.FeatureOutlierThreshold.ValueFloat64()
	}
	if !plan.IsTracePrompt.IsNull() {
		s["isTracePrompt"] = plan.IsTracePrompt.ValueBool()
	}
	if !plan.IsEdgeBrain.IsNull() {
		s["isEdgeBrain"] = plan.IsEdgeBrain.ValueBool()
	}

	// Incident prediction and RCA fields
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
	if !plan.RootCauseLogMessageSearchRange.IsNull() {
		s["rootCauseLogMessageSearchRange"] = int(plan.RootCauseLogMessageSearchRange.ValueInt64())
	}
	if !plan.CausalPredictionSetting.IsNull() {
		s["causalPredictionSetting"] = int(plan.CausalPredictionSetting.ValueInt64())
	}
	if !plan.CausalMinDelay.IsNull() {
		s["causalMinDelay"] = plan.CausalMinDelay.ValueString()
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
	if !plan.AvgPerIncidentDowntimeCost.IsNull() {
		s["avgPerIncidentDowntimeCost"] = plan.AvgPerIncidentDowntimeCost.ValueFloat64()
	}
	if !plan.PredictionRuleActiveCondition.IsNull() {
		s["predictionRuleActiveCondition"] = int(plan.PredictionRuleActiveCondition.ValueInt64())
	}
	if !plan.PredictionRuleFalsePositiveThreshold.IsNull() {
		s["predictionRuleFalsePositiveThreshold"] = int(plan.PredictionRuleFalsePositiveThreshold.ValueInt64())
	}
	if !plan.PredictionRuleActiveThreshold.IsNull() {
		s["predictionRuleActiveThreshold"] = plan.PredictionRuleActiveThreshold.ValueFloat64()
	}
	if !plan.PredictionRuleInactiveThreshold.IsNull() {
		s["predictionRuleInactiveThreshold"] = plan.PredictionRuleInactiveThreshold.ValueFloat64()
	}
	if !plan.PredictionProbabilityThreshold.IsNull() {
		s["predictionProbabilityThreshold"] = plan.PredictionProbabilityThreshold.ValueFloat64()
	}
	if !plan.AlertHourlyCost.IsNull() {
		s["alertHourlyCost"] = plan.AlertHourlyCost.ValueFloat64()
	}
	if !plan.AlertAverageTime.IsNull() {
		s["alertAverageTime"] = int(plan.AlertAverageTime.ValueInt64())
	}
	if !plan.IgnoreInstanceForKB.IsNull() {
		s["ignoreInstanceForKB"] = plan.IgnoreInstanceForKB.ValueBool()
	}
	if !plan.ShowInstanceDown.IsNull() {
		s["showInstanceDown"] = plan.ShowInstanceDown.ValueBool()
	}
	if !plan.PredictionCountThreshold.IsNull() {
		s["predictionCountThreshold"] = int(plan.PredictionCountThreshold.ValueInt64())
	}

	// Complex JSON fields
	if !plan.BaseValueSetting.IsNull() {
		if parsed := parseJSONField(plan.BaseValueSetting.ValueString()); parsed != nil {
			s["baseValueSetting"] = parsed
		}
	}
	if !plan.CdfSetting.IsNull() {
		if parsed := parseJSONField(plan.CdfSetting.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["cdfSetting"] = list
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
	if !plan.LlmEvaluationSetting.IsNull() {
		if parsed := parseJSONField(plan.LlmEvaluationSetting.ValueString()); parsed != nil {
			s["llmEvaluationSetting"] = parsed
		}
	}
	if !plan.LogToLogSettingList.IsNull() {
		if parsed := parseJSONField(plan.LogToLogSettingList.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["logToLogSettingList"] = list
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
	if !plan.SharedUsernames.IsNull() {
		if parsed := parseJSONField(plan.SharedUsernames.ValueString()); parsed != nil {
			if list, ok := parsed.([]interface{}); ok {
				s["sharedUsernames"] = list
			}
		}
	}

	return s
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
		ProjectName:        plan.ProjectName.ValueString(),
		ProjectDisplayName: plan.ProjectDisplayName.ValueString(),
		SystemName:         plan.SystemName.ValueString(),
		DataType:           plan.ProjectCreationConfig.DataType.ValueString(),
		InstanceType:       plan.ProjectCreationConfig.InstanceType.ValueString(),
		ProjectCloudType:   plan.ProjectCreationConfig.ProjectCloudType.ValueString(),
		InsightAgentType:   plan.ProjectCreationConfig.InsightAgentType.ValueString(),
		ServiceNowTable:    plan.ProjectCreationConfig.ServiceNowTable.ValueString(),
		CValue:             int(plan.CValue.ValueInt64()),
		PValue:             plan.PValue.ValueFloat64(),
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
		// If we can't read from API, use the config as state but with ID set
		config.ID = plan.ProjectName
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
		plan.ComponentNameAutoOverwrite = getBool("componentNameAutoOverwrite")
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
		// Preserve user's email_setting from config to avoid drift from API-enriched values
		plan.EmailSetting = config.EmailSetting
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

			// Convert back to types.Set
			setValue, diags := types.SetValueFrom(ctx, config.LogLabelSettings.ElementType(ctx), normalizedSettings)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			plan.LogLabelSettings = setValue
		}
	}

	// If LogLabelSettings is null or empty, set it to an empty set
	if config.LogLabelSettings.IsNull() {
		emptySet, diags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"label_type":       types.StringType,
				"log_label_string": types.StringType,
			},
		}, []logLabelSettingModel{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.LogLabelSettings = emptySet
		}
	}

	// Process project_servicenow_settings if project_cloud_type is ServiceNow
	if !config.ProjectServiceNowSettings.IsNull() && !config.ProjectServiceNowSettings.IsUnknown() {
		// Check if project_cloud_type is ServiceNow (case insensitive)
		cloudType := config.ProjectCreationConfig.ProjectCloudType.ValueString()
		if cloudType != "" && (cloudType == "ServiceNow" || cloudType == "servicenow" || cloudType == "SERVICENOW") {
			var serviceNowSettings projectServiceNowSettingsModel
			diags = config.ProjectServiceNowSettings.As(ctx, &serviceNowSettings, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			tflog.Info(ctx, "Processing ServiceNow third-party settings")

			// Convert to client model
			clientSettings := &client.ServiceNowThirdPartySettings{
				Host:               serviceNowSettings.Host.ValueString(),
				SysparmQuery:       serviceNowSettings.SysparmQuery.ValueString(),
				Proxy:              serviceNowSettings.Proxy.ValueString(),
				ServiceNowUser:     serviceNowSettings.ServiceNowUser.ValueString(),
				ServiceNowPassword: serviceNowSettings.ServiceNowPassword.ValueString(),
				InstanceField:      serviceNowSettings.InstanceField.ValueString(),
				InstanceFieldRegex: serviceNowSettings.InstanceFieldRegex.ValueString(),
				TimestampFormat:    serviceNowSettings.TimestampFormat.ValueString(),
				ClientID:           serviceNowSettings.ClientID.ValueString(),
				ClientSecret:       serviceNowSettings.ClientSecret.ValueString(),
			}

			// Convert additional fields list
			if !serviceNowSettings.AdditionalFields.IsNull() && !serviceNowSettings.AdditionalFields.IsUnknown() {
				var additionalFields []string
				diags = serviceNowSettings.AdditionalFields.ElementsAs(ctx, &additionalFields, false)
				resp.Diagnostics.Append(diags...)
				if !resp.Diagnostics.HasError() {
					clientSettings.AdditionalFields = additionalFields
				}
			}

			// Apply ServiceNow settings
			err := r.client.CreateOrUpdateServiceNowThirdPartySettings(plan.ProjectName.ValueString(), clientSettings)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error applying ServiceNow third-party settings",
					fmt.Sprintf("Could not apply ServiceNow settings: %s", err.Error()),
				)
				return
			}

			// Store the config in plan
			plan.ProjectServiceNowSettings = config.ProjectServiceNowSettings
		} else {
			tflog.Warn(ctx, "project_servicenow_settings provided but project_cloud_type is not ServiceNow, ignoring settings", map[string]any{
				"project_cloud_type": cloudType,
			})
			// Set to null for non-ServiceNow projects
			plan.ProjectServiceNowSettings = types.ObjectNull(map[string]attr.Type{
				"host":                 types.StringType,
				"sysparm_query":        types.StringType,
				"proxy":                types.StringType,
				"servicenow_user":      types.StringType,
				"servicenow_password":  types.StringType,
				"instance_field":       types.StringType,
				"instance_field_regex": types.StringType,
				"timestamp_format":     types.StringType,
				"client_id":            types.StringType,
				"client_secret":        types.StringType,
				"additional_fields":    types.ListType{ElemType: types.StringType},
			})
		}
	} else {
		// No ServiceNow settings in config - explicitly set to null
		plan.ProjectServiceNowSettings = types.ObjectNull(map[string]attr.Type{
			"host":                 types.StringType,
			"sysparm_query":        types.StringType,
			"proxy":                types.StringType,
			"servicenow_user":      types.StringType,
			"servicenow_password":  types.StringType,
			"instance_field":       types.StringType,
			"instance_field_regex": types.StringType,
			"timestamp_format":     types.StringType,
			"client_id":            types.StringType,
			"client_secret":        types.StringType,
			"additional_fields":    types.ListType{ElemType: types.StringType},
		})
	}

	// Process holiday_settings
	if !config.HolidaySettings.IsNull() && !config.HolidaySettings.IsUnknown() {
		var holidays []holidaySettingModel
		diags = config.HolidaySettings.ElementsAs(ctx, &holidays, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Info(ctx, "Processing holiday settings", map[string]any{"count": len(holidays)})

		// Create each holiday (preserving config order)
		for _, holiday := range holidays {
			clientHoliday := &client.Holiday{
				Name:      holiday.Name.ValueString(),
				StartDate: holiday.StartDate.ValueString(),
				EndDate:   holiday.EndDate.ValueString(),
			}

			err := r.client.CreateHoliday(plan.ProjectName.ValueString(), clientHoliday)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating holiday",
					fmt.Sprintf("Could not create holiday '%s': %s", clientHoliday.Name, err.Error()),
				)
				return
			}
		}

		plan.HolidaySettings = config.HolidaySettings
	} else {
		// No holiday settings in config - set to null
		plan.HolidaySettings = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":       types.StringType,
				"start_date": types.StringType,
				"end_date":   types.StringType,
			},
		})
	}

	// Process json_key_settings
	if !config.JsonKeySettings.IsNull() && !config.JsonKeySettings.IsUnknown() {
		var configJsonKeys []jsonKeySettingModel
		diags = config.JsonKeySettings.ElementsAs(ctx, &configJsonKeys, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Info(ctx, "Processing JSON key settings", map[string]any{"count": len(configJsonKeys)})

		// Convert to client format
		var jsonKeysToUpdate []client.JsonKeyType
		var summaryKeys []string
		var metafieldKeys []string
		var dampeningFieldKeys []string

		for _, jsonKeySetting := range configJsonKeys {
			jsonKeyType := client.JsonKeyType{
				JsonKey:        jsonKeySetting.JsonKey.ValueString(),
				Type:           jsonKeySetting.Type.ValueString(),
				SummaryCheck:   jsonKeySetting.SummarySetting.ValueBool(),
				MetaFieldCheck: jsonKeySetting.MetafieldSetting.ValueBool(),
			}
			jsonKeysToUpdate = append(jsonKeysToUpdate, jsonKeyType)

			// Track which keys have summary settings enabled
			if jsonKeySetting.SummarySetting.ValueBool() {
				summaryKeys = append(summaryKeys, jsonKeySetting.JsonKey.ValueString())
			}

			// Track which keys have metafield settings enabled
			if jsonKeySetting.MetafieldSetting.ValueBool() {
				metafieldKeys = append(metafieldKeys, jsonKeySetting.JsonKey.ValueString())
			}

			// Track which keys have dampening field settings enabled
			if jsonKeySetting.DampeningfieldSetting.ValueBool() {
				dampeningFieldKeys = append(dampeningFieldKeys, jsonKeySetting.JsonKey.ValueString())
			}
		}

		// Update JSON key types
		if len(jsonKeysToUpdate) > 0 {
			err := r.client.UpdateJsonKeyTypes(plan.ProjectName.ValueString(), jsonKeysToUpdate)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating JSON key types",
					fmt.Sprintf("Could not update JSON key types: %s", err.Error()),
				)
				return
			}
		}

		// Update summary, metafield, and dampening field settings
		err := r.client.UpdateJsonKeySummarySettings(plan.ProjectName.ValueString(), summaryKeys, metafieldKeys, dampeningFieldKeys)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating JSON key summary and metafield settings",
				fmt.Sprintf("Could not update settings: %s", err.Error()),
			)
			return
		}

		// Preserve the config order in plan
		plan.JsonKeySettings = config.JsonKeySettings
	} else {
		// No JSON key settings in config - set to null
		plan.JsonKeySettings = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"json_key":                types.StringType,
				"type":                    types.StringType,
				"summary_setting":         types.BoolType,
				"metafield_setting":       types.BoolType,
				"dampening_field_setting": types.BoolType,
			},
		})
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
	state.ComponentNameAutoOverwrite = getBool("componentNameAutoOverwrite")
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

		// Convert back to types.Set
		setValue, diags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"label_type":       types.StringType,
				"log_label_string": types.StringType,
			},
		}, convertedSettings)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.LogLabelSettings = setValue
		}
	}

	// Read ServiceNow third-party settings if project_cloud_type is ServiceNow
	// Check if we have a ProjectCreationConfig in state first
	if state.ProjectCreationConfig != nil {
		cloudType := state.ProjectCreationConfig.ProjectCloudType.ValueString()
		if cloudType != "" && (cloudType == "ServiceNow" || cloudType == "servicenow" || cloudType == "SERVICENOW") {
			serviceNowSettings, err := r.client.GetServiceNowThirdPartySettings(state.ProjectName.ValueString())
			if err != nil {
				tflog.Warn(ctx, "Could not read ServiceNow third-party settings", map[string]any{"error": err.Error()})
				// Keep existing state if we can't read from API
			} else if serviceNowSettings != nil {
				// Convert to terraform model
				serviceNowModel := projectServiceNowSettingsModel{
					Host:               types.StringValue(serviceNowSettings.Host),
					SysparmQuery:       types.StringValue(serviceNowSettings.SysparmQuery),
					Proxy:              types.StringValue(serviceNowSettings.Proxy),
					ServiceNowUser:     types.StringValue(serviceNowSettings.ServiceNowUser),
					ServiceNowPassword: types.StringValue(serviceNowSettings.ServiceNowPassword),
					InstanceField:      types.StringValue(serviceNowSettings.InstanceField),
					InstanceFieldRegex: types.StringValue(serviceNowSettings.InstanceFieldRegex),
					TimestampFormat:    types.StringValue(serviceNowSettings.TimestampFormat),
					ClientID:           types.StringValue(serviceNowSettings.ClientID),
					ClientSecret:       types.StringValue(serviceNowSettings.ClientSecret),
				}

				// Convert additional fields
				if len(serviceNowSettings.AdditionalFields) > 0 {
					additionalFieldsList, diags := types.ListValueFrom(ctx, types.StringType, serviceNowSettings.AdditionalFields)
					resp.Diagnostics.Append(diags...)
					if !resp.Diagnostics.HasError() {
						serviceNowModel.AdditionalFields = additionalFieldsList
					}
				} else {
					serviceNowModel.AdditionalFields = types.ListNull(types.StringType)
				}

				// Convert to types.Object
				serviceNowObject, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
					"host":                 types.StringType,
					"sysparm_query":        types.StringType,
					"proxy":                types.StringType,
					"servicenow_user":      types.StringType,
					"servicenow_password":  types.StringType,
					"instance_field":       types.StringType,
					"instance_field_regex": types.StringType,
					"timestamp_format":     types.StringType,
					"client_id":            types.StringType,
					"client_secret":        types.StringType,
					"additional_fields":    types.ListType{ElemType: types.StringType},
				}, serviceNowModel)
				resp.Diagnostics.Append(diags...)
				if !resp.Diagnostics.HasError() {
					state.ProjectServiceNowSettings = serviceNowObject
				}
			}
		}
	}

	// Read holiday settings from API
	holidays, err := r.client.GetHolidays(state.ProjectName.ValueString())
	if err != nil {
		tflog.Warn(ctx, "Could not read holidays", map[string]any{"error": err.Error()})
		// Keep existing state if we can't read from API
	} else if holidays != nil {
		// Convert API response to state model
		// Response format: {"holidayName": "startDate,endDate", ...}

		// Build holiday settings from API response
		var holidaySettings []holidaySettingModel
		for name, dates := range holidays {
			// Parse the dates string "MM-DD,MM-DD"
			var startDate, endDate string
			if dates != "" {
				// Simple split by comma
				parts := []string{}
				current := ""
				for _, c := range dates {
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

				if len(parts) >= 2 {
					startDate = parts[0]
					endDate = parts[1]
				} else if len(parts) == 1 {
					startDate = parts[0]
					endDate = parts[0]
				}
			}

			holidaySettings = append(holidaySettings, holidaySettingModel{
				Name:      types.StringValue(name),
				StartDate: types.StringValue(startDate),
				EndDate:   types.StringValue(endDate),
			})
		}

		// Convert to types.Set
		if len(holidaySettings) > 0 {
			setValue, diags := types.SetValueFrom(ctx, types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":       types.StringType,
					"start_date": types.StringType,
					"end_date":   types.StringType,
				},
			}, holidaySettings)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				state.HolidaySettings = setValue
			}
		} else {
			state.HolidaySettings = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":       types.StringType,
					"start_date": types.StringType,
					"end_date":   types.StringType,
				},
			})
		}
	}

	// Read JSON key settings from API
	jsonKeyTypes, err := r.client.GetJsonKeyTypes(state.ProjectName.ValueString())
	if err != nil {
		tflog.Warn(ctx, "Could not read JSON key types", map[string]any{"error": err.Error()})
		// Keep existing state if we can't read from API
	} else if len(jsonKeyTypes) > 0 {
		// Get summary, metafield, and dampening field settings
		summarySettingsResp, err := r.client.GetJsonKeySummarySettings(state.ProjectName.ValueString())
		if err != nil {
			tflog.Warn(ctx, "Could not read JSON key summary settings", map[string]any{"error": err.Error()})
			summarySettingsResp = &client.JsonKeySummarySettings{}
		}

		// Create sets for quick lookup
		summarySet := make(map[string]bool)
		for _, key := range summarySettingsResp.SummarySetting {
			summarySet[key] = true
		}

		metafieldSet := make(map[string]bool)
		for _, key := range summarySettingsResp.MetaFieldSetting {
			metafieldSet[key] = true
		}

		dampeningFieldSet := make(map[string]bool)
		for _, key := range summarySettingsResp.DampeningFieldSetting {
			dampeningFieldSet[key] = true
		}

		// Build the set of JSON key settings from API response
		var jsonKeySettings []jsonKeySettingModel
		for _, jsonKey := range jsonKeyTypes {
			jsonKeySettings = append(jsonKeySettings, jsonKeySettingModel{
				JsonKey:               types.StringValue(jsonKey.JsonKey),
				Type:                  types.StringValue(jsonKey.Type),
				SummarySetting:        types.BoolValue(summarySet[jsonKey.JsonKey]),
				MetafieldSetting:      types.BoolValue(metafieldSet[jsonKey.JsonKey]),
				DampeningfieldSetting: types.BoolValue(dampeningFieldSet[jsonKey.JsonKey]),
			})
		}

		// Convert to types.Set
		if len(jsonKeySettings) > 0 {
			setValue, diags := types.SetValueFrom(ctx, types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"json_key":                types.StringType,
					"type":                    types.StringType,
					"summary_setting":         types.BoolType,
					"metafield_setting":       types.BoolType,
					"dampening_field_setting": types.BoolType,
				},
			}, jsonKeySettings)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				state.JsonKeySettings = setValue
			}
		} else {
			state.JsonKeySettings = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"json_key":                types.StringType,
					"type":                    types.StringType,
					"summary_setting":         types.BoolType,
					"metafield_setting":       types.BoolType,
					"dampening_field_setting": types.BoolType,
				},
			})
		}
	} else {
		state.JsonKeySettings = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"json_key":                types.StringType,
				"type":                    types.StringType,
				"summary_setting":         types.BoolType,
				"metafield_setting":       types.BoolType,
				"dampening_field_setting": types.BoolType,
			},
		})
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
		plan.ComponentNameAutoOverwrite = getBool("componentNameAutoOverwrite")
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
		// PredictionRule* fields are managed by the incident prediction API, not the
		// watch-tower-setting GET API, so it returns stale values. Preserve config values.
		plan.PredictionCountThreshold = getInt64("predictionCountThreshold")
		plan.PredictionProbabilityThreshold = getFloat64("predictionProbabilityThreshold")
		if !config.PredictionRuleActiveCondition.IsNull() {
			plan.PredictionRuleActiveCondition = config.PredictionRuleActiveCondition
		} else {
			plan.PredictionRuleActiveCondition = getInt64("predictionRuleActiveCondition")
		}
		if !config.PredictionRuleActiveThreshold.IsNull() {
			plan.PredictionRuleActiveThreshold = config.PredictionRuleActiveThreshold
		} else {
			plan.PredictionRuleActiveThreshold = getFloat64("predictionRuleActiveThreshold")
		}
		if !config.PredictionRuleFalsePositiveThreshold.IsNull() {
			plan.PredictionRuleFalsePositiveThreshold = config.PredictionRuleFalsePositiveThreshold
		} else {
			plan.PredictionRuleFalsePositiveThreshold = getInt64("predictionRuleFalsePositiveThreshold")
		}
		if !config.PredictionRuleInactiveThreshold.IsNull() {
			plan.PredictionRuleInactiveThreshold = config.PredictionRuleInactiveThreshold
		} else {
			plan.PredictionRuleInactiveThreshold = getFloat64("predictionRuleInactiveThreshold")
		}
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
		// Preserve user's email_setting from config to avoid drift from API-enriched values
		plan.EmailSetting = config.EmailSetting
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

		// Get current API labels to detect what needs to be deleted
		apiLabels, err := r.client.GetLogLabels(plan.ProjectName.ValueString(), r.client.Username)
		if err != nil {
			tflog.Warn(ctx, "Failed to get current log labels for comparison", map[string]any{"error": err.Error()})
			apiLabels = make(map[string]string) // Continue with empty map
		}

		// Build a map of label types in config
		configLabelTypes := make(map[string]bool)
		for _, setting := range configSettings {
			configLabelTypes[setting.LabelType.ValueString()] = true
		}

		// Map API field names to label types for comparison
		apiFieldToLabelType := map[string]string{
			"whitelist":                       "whitelist",
			"trainingWhitelist":               "trainingWhitelist",
			"trainingBlacklistLabels":         "blacklist",
			"featurelist":                     "featurelist",
			"incidentlist":                    "incidentlist",
			"triagelist":                      "triagelist",
			"patternNameLabels":               "patternName",
			"patternSignatureLabels":          "patternSignature",
			"patternMatchRegexLabels":         "patternMatchRegex",
			"patternIgnoreRegexLabels":        "patternIgnoreRegex",
			"customActionLabels":              "customAction",
			"logEventIDLabels":                "logEventID",
			"logSeverityLabels":               "logSeverity",
			"logStatusCodeLabels":             "logStatusCode",
			"alertEventTypeLabels":            "alertEventType",
			"anomalyFeatureLabels":            "anomalyFeature",
			"dataFilterLabels":                "dataFilter",
			"instanceNameLabels":              "instanceName",
			"dataQualityCheckLabels":          "dataQualityCheck",
			"incidentFieldVerificationLabels": "incidentFieldVerification",
			"incidentPriorityLabels":          "incidentPriority",
			"extractionBlacklist":             "extractionBlacklist",
		}

		// Find labels in API that are not in config (need to be deleted)
		labelsToDelete := []string{}
		for apiField, jsonString := range apiLabels {
			if labelType, ok := apiFieldToLabelType[apiField]; ok {
				// Only delete if it has non-empty content and is not in config
				if jsonString != "" && jsonString != "[]" && !configLabelTypes[labelType] {
					labelsToDelete = append(labelsToDelete, labelType)
					tflog.Info(ctx, "Label will be deleted (in API but not in config)", map[string]any{"label_type": labelType})
				}
			}
		}

		// Delete labels that are in API but not in config
		if len(labelsToDelete) > 0 {
			tflog.Info(ctx, "Deleting log labels not in config", map[string]any{"count": len(labelsToDelete)})
			err := r.client.DeleteLogLabels(
				plan.ProjectName.ValueString(),
				r.client.Username,
				labelsToDelete,
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting log label settings",
					fmt.Sprintf("Could not delete log label settings: %s", err.Error()),
				)
				return
			}
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
		// Explicitly set empty set if none configured
		emptySet, diags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"label_type":       types.StringType,
				"log_label_string": types.StringType,
			},
		}, []logLabelSettingModel{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.LogLabelSettings = emptySet
		}
	}

	// Process project_servicenow_settings if provided and project_cloud_type is ServiceNow
	if !config.ProjectServiceNowSettings.IsNull() && !config.ProjectServiceNowSettings.IsUnknown() {
		// Check if project_cloud_type is ServiceNow (case insensitive)
		cloudType := config.ProjectCreationConfig.ProjectCloudType.ValueString()
		if cloudType != "" && (cloudType == "ServiceNow" || cloudType == "servicenow" || cloudType == "SERVICENOW") {
			var serviceNowSettings projectServiceNowSettingsModel
			diags = config.ProjectServiceNowSettings.As(ctx, &serviceNowSettings, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			tflog.Info(ctx, "Updating ServiceNow third-party settings")

			// Convert to client model
			clientSettings := &client.ServiceNowThirdPartySettings{
				Host:               serviceNowSettings.Host.ValueString(),
				SysparmQuery:       serviceNowSettings.SysparmQuery.ValueString(),
				Proxy:              serviceNowSettings.Proxy.ValueString(),
				ServiceNowUser:     serviceNowSettings.ServiceNowUser.ValueString(),
				ServiceNowPassword: serviceNowSettings.ServiceNowPassword.ValueString(),
				InstanceField:      serviceNowSettings.InstanceField.ValueString(),
				InstanceFieldRegex: serviceNowSettings.InstanceFieldRegex.ValueString(),
				TimestampFormat:    serviceNowSettings.TimestampFormat.ValueString(),
				ClientID:           serviceNowSettings.ClientID.ValueString(),
				ClientSecret:       serviceNowSettings.ClientSecret.ValueString(),
			}

			// Convert additional fields list
			if !serviceNowSettings.AdditionalFields.IsNull() && !serviceNowSettings.AdditionalFields.IsUnknown() {
				var additionalFields []string
				diags = serviceNowSettings.AdditionalFields.ElementsAs(ctx, &additionalFields, false)
				resp.Diagnostics.Append(diags...)
				if !resp.Diagnostics.HasError() {
					clientSettings.AdditionalFields = additionalFields
				}
			}

			// Apply ServiceNow settings
			err := r.client.CreateOrUpdateServiceNowThirdPartySettings(plan.ProjectName.ValueString(), clientSettings)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating ServiceNow third-party settings",
					fmt.Sprintf("Could not update ServiceNow settings: %s", err.Error()),
				)
				return
			}

			// Store the config in plan
			plan.ProjectServiceNowSettings = config.ProjectServiceNowSettings
		} else {
			tflog.Warn(ctx, "project_servicenow_settings provided but project_cloud_type is not ServiceNow, ignoring settings", map[string]any{
				"project_cloud_type": cloudType,
			})
			// Set to null for non-ServiceNow projects
			plan.ProjectServiceNowSettings = types.ObjectNull(map[string]attr.Type{
				"host":                 types.StringType,
				"sysparm_query":        types.StringType,
				"proxy":                types.StringType,
				"servicenow_user":      types.StringType,
				"servicenow_password":  types.StringType,
				"instance_field":       types.StringType,
				"instance_field_regex": types.StringType,
				"timestamp_format":     types.StringType,
				"client_id":            types.StringType,
				"client_secret":        types.StringType,
				"additional_fields":    types.ListType{ElemType: types.StringType},
			})
		}
	} else {
		// No ServiceNow settings in config - explicitly set to null
		plan.ProjectServiceNowSettings = types.ObjectNull(map[string]attr.Type{
			"host":                 types.StringType,
			"sysparm_query":        types.StringType,
			"proxy":                types.StringType,
			"servicenow_user":      types.StringType,
			"servicenow_password":  types.StringType,
			"instance_field":       types.StringType,
			"instance_field_regex": types.StringType,
			"timestamp_format":     types.StringType,
			"client_id":            types.StringType,
			"client_secret":        types.StringType,
			"additional_fields":    types.ListType{ElemType: types.StringType},
		})
	}

	// Process holiday_settings
	if !config.HolidaySettings.IsNull() && !config.HolidaySettings.IsUnknown() {
		var configHolidays []holidaySettingModel
		diags = config.HolidaySettings.ElementsAs(ctx, &configHolidays, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Info(ctx, "Processing holiday settings update", map[string]any{"count": len(configHolidays)})

		// Get current holidays from API
		currentHolidays, err := r.client.GetHolidays(plan.ProjectName.ValueString())
		if err != nil {
			tflog.Warn(ctx, "Could not get current holidays", map[string]any{"error": err.Error()})
			currentHolidays = make(map[string]string)
		}

		// Create a map of config holiday names
		configHolidayNames := make(map[string]bool)
		for _, h := range configHolidays {
			configHolidayNames[h.Name.ValueString()] = true
		}

		// Find holidays to delete (in API but not in config)
		var holidaysToDelete []string
		for name := range currentHolidays {
			if !configHolidayNames[name] {
				holidaysToDelete = append(holidaysToDelete, name)
			}
		}

		// Delete holidays not in config
		if len(holidaysToDelete) > 0 {
			tflog.Info(ctx, "Deleting holidays not in config", map[string]any{"count": len(holidaysToDelete)})
			err := r.client.DeleteHolidays(plan.ProjectName.ValueString(), holidaysToDelete)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting holidays",
					fmt.Sprintf("Could not delete holidays: %s", err.Error()),
				)
				return
			}
		}

		// Create or update holidays from config
		for _, holiday := range configHolidays {
			clientHoliday := &client.Holiday{
				Name:      holiday.Name.ValueString(),
				StartDate: holiday.StartDate.ValueString(),
				EndDate:   holiday.EndDate.ValueString(),
			}

			// Check if holiday needs to be updated
			if existingDates, exists := currentHolidays[clientHoliday.Name]; exists {
				// Parse existing dates
				parts := []string{}
				current := ""
				for _, c := range existingDates {
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

				var existingStart, existingEnd string
				if len(parts) >= 2 {
					existingStart = parts[0]
					existingEnd = parts[1]
				} else if len(parts) == 1 {
					existingStart = parts[0]
					existingEnd = parts[0]
				}

				// Skip if dates are the same
				if existingStart == clientHoliday.StartDate && existingEnd == clientHoliday.EndDate {
					tflog.Debug(ctx, "Holiday unchanged, skipping", map[string]any{"name": clientHoliday.Name})
					continue
				}

				// Delete and recreate to update
				tflog.Info(ctx, "Updating holiday", map[string]any{"name": clientHoliday.Name})
				err := r.client.DeleteHolidays(plan.ProjectName.ValueString(), []string{clientHoliday.Name})
				if err != nil {
					resp.Diagnostics.AddError(
						"Error updating holiday",
						fmt.Sprintf("Could not delete existing holiday '%s': %s", clientHoliday.Name, err.Error()),
					)
					return
				}
			}

			// Create the holiday
			err := r.client.CreateHoliday(plan.ProjectName.ValueString(), clientHoliday)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating holiday",
					fmt.Sprintf("Could not create holiday '%s': %s", clientHoliday.Name, err.Error()),
				)
				return
			}
		}

		plan.HolidaySettings = config.HolidaySettings
	} else {
		// No holiday settings in config - delete all existing holidays
		tflog.Info(ctx, "Holiday settings removed from config, deleting all holidays")

		// Get current holidays from API
		currentHolidays, err := r.client.GetHolidays(plan.ProjectName.ValueString())
		if err != nil {
			tflog.Warn(ctx, "Could not get current holidays for deletion", map[string]any{"error": err.Error()})
		} else if len(currentHolidays) > 0 {
			// Delete all holidays
			var holidaysToDelete []string
			for name := range currentHolidays {
				holidaysToDelete = append(holidaysToDelete, name)
			}

			tflog.Info(ctx, "Deleting all holidays", map[string]any{"count": len(holidaysToDelete)})
			err := r.client.DeleteHolidays(plan.ProjectName.ValueString(), holidaysToDelete)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting holidays",
					fmt.Sprintf("Could not delete holidays: %s", err.Error()),
				)
				return
			}
		}

		// Set to null
		plan.HolidaySettings = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":       types.StringType,
				"start_date": types.StringType,
				"end_date":   types.StringType,
			},
		})
	}

	// Process json_key_settings
	if !config.JsonKeySettings.IsNull() && !config.JsonKeySettings.IsUnknown() {
		var configJsonKeys []jsonKeySettingModel
		diags = config.JsonKeySettings.ElementsAs(ctx, &configJsonKeys, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Info(ctx, "Processing JSON key settings", map[string]any{"count": len(configJsonKeys)})

		// Convert to client format
		var jsonKeysToUpdate []client.JsonKeyType
		var summaryKeys []string
		var metafieldKeys []string
		var dampeningFieldKeys []string

		for _, jsonKeySetting := range configJsonKeys {
			jsonKeyType := client.JsonKeyType{
				JsonKey:        jsonKeySetting.JsonKey.ValueString(),
				Type:           jsonKeySetting.Type.ValueString(),
				SummaryCheck:   jsonKeySetting.SummarySetting.ValueBool(),
				MetaFieldCheck: jsonKeySetting.MetafieldSetting.ValueBool(),
			}
			jsonKeysToUpdate = append(jsonKeysToUpdate, jsonKeyType)

			// Track which keys have summary settings enabled
			if jsonKeySetting.SummarySetting.ValueBool() {
				summaryKeys = append(summaryKeys, jsonKeySetting.JsonKey.ValueString())
			}

			// Track which keys have metafield settings enabled
			if jsonKeySetting.MetafieldSetting.ValueBool() {
				metafieldKeys = append(metafieldKeys, jsonKeySetting.JsonKey.ValueString())
			}

			// Track which keys have dampening field settings enabled
			if jsonKeySetting.DampeningfieldSetting.ValueBool() {
				dampeningFieldKeys = append(dampeningFieldKeys, jsonKeySetting.JsonKey.ValueString())
			}
		}

		// Update JSON key types
		if len(jsonKeysToUpdate) > 0 {
			err := r.client.UpdateJsonKeyTypes(plan.ProjectName.ValueString(), jsonKeysToUpdate)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating JSON key types",
					fmt.Sprintf("Could not update JSON key types: %s", err.Error()),
				)
				return
			}
		}

		// Update summary, metafield, and dampening field settings
		err := r.client.UpdateJsonKeySummarySettings(plan.ProjectName.ValueString(), summaryKeys, metafieldKeys, dampeningFieldKeys)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating JSON key summary and metafield settings",
				fmt.Sprintf("Could not update settings: %s", err.Error()),
			)
			return
		}

		plan.JsonKeySettings = config.JsonKeySettings
	} else {
		// No JSON key settings in config - set to null
		plan.JsonKeySettings = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"json_key":                types.StringType,
				"type":                    types.StringType,
				"summary_setting":         types.BoolType,
				"metafield_setting":       types.BoolType,
				"dampening_field_setting": types.BoolType,
			},
		})
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
// Compares labels in sorted order (order-independent) but preserves original order in state
func convertLogLabelsToState(apiLabels map[string]string, existingState []logLabelSettingModel) []logLabelSettingModel {
	var result []logLabelSettingModel

	// Map from API field names to label types
	apiFieldToLabelType := map[string]string{
		"whitelist":                       "whitelist",
		"trainingWhitelist":               "trainingWhitelist",
		"trainingBlacklistLabels":         "blacklist",
		"featurelist":                     "featurelist",
		"incidentlist":                    "incidentlist",
		"triagelist":                      "triagelist",
		"patternNameLabels":               "patternName",
		"patternSignatureLabels":          "patternSignature",
		"patternMatchRegexLabels":         "patternMatchRegex",
		"patternIgnoreRegexLabels":        "patternIgnoreRegex",
		"customActionLabels":              "customAction",
		"logEventIDLabels":                "logEventID",
		"logSeverityLabels":               "logSeverity",
		"logStatusCodeLabels":             "logStatusCode",
		"alertEventTypeLabels":            "alertEventType",
		"anomalyFeatureLabels":            "anomalyFeature",
		"dataFilterLabels":                "dataFilter",
		"instanceNameLabels":              "instanceName",
		"dataQualityCheckLabels":          "dataQualityCheck",
		"incidentFieldVerificationLabels": "incidentFieldVerification",
		"incidentPriorityLabels":          "incidentPriority",
		"extractionBlacklist":             "extractionBlacklist",
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

	// Create a map of existing state for comparison (order-independent)
	existingStateMap := make(map[string]string)
	for _, existing := range existingState {
		labelType := existing.LabelType.ValueString()
		existingStateMap[labelType] = normalizeJSON(existing.LogLabelString.ValueString())
	}

	// Compare the two maps (order-independent comparison)
	stateMatchesAPI := len(existingStateMap) == len(apiDataMap)
	if stateMatchesAPI {
		for labelType, existingValue := range existingStateMap {
			apiValue, exists := apiDataMap[labelType]
			if !exists || apiValue != existingValue {
				stateMatchesAPI = false
				break
			}
		}
	}

	// If state matches API (ignoring order), preserve the existing order
	if stateMatchesAPI && len(existingState) > 0 {
		return existingState
	}

	// State doesn't match API - rebuild from API data
	// If we have existing state, preserve its order for labels that still exist
	if len(existingState) > 0 {
		processedTypes := make(map[string]bool)

		// First pass: preserve order from existing state where labels still exist in API
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
		// This allows detection of labels added in the UI
		defaultOrder := []string{
			"trainingWhitelist",
			"whitelist",
			"blacklist",
			"featurelist",
			"incidentlist",
			"triagelist",
			"patternName",
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
			"incidentFieldVerification",
			"incidentPriority",
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
			"incidentFieldVerification",
			"incidentPriority",
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
