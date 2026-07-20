// Copyright (c) InsightFinder Inc.
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &systemSettingsResource{}
	_ resource.ResourceWithConfigure   = &systemSettingsResource{}
	_ resource.ResourceWithImportState = &systemSettingsResource{}
)

// NewSystemSettingsResource is a helper function to simplify the provider implementation.
func NewSystemSettingsResource() resource.Resource {
	return &systemSettingsResource{}
}

// systemSettingsResource is the resource implementation.
type systemSettingsResource struct {
	client *client.Client
}

// systemSettingsResourceModel maps the resource schema data.
type systemSettingsResourceModel struct {
	ID                    types.String                `tfsdk:"id"`
	SystemName            types.String                `tfsdk:"system_name"`
	KnowledgebaseSettings *knowledgebaseSettingsModel `tfsdk:"knowledgebase_settings"`
	NotificationsSettings *notificationsSettingsModel `tfsdk:"notifications_settings"`
	MiscellaneousSettings *miscellaneousSettingsModel `tfsdk:"miscellaneous_settings"`
}

// miscellaneousSettingsModel holds miscellaneous system framework settings
type miscellaneousSettingsModel struct {
	HealthviewLongterm                   types.Bool  `tfsdk:"healthview_longterm"`
	ShouldAutoShare                      types.Bool  `tfsdk:"should_auto_share"`
	RootcauseReverseEntryFilterThreshold types.Int64 `tfsdk:"rootcause_reverse_entry_filter_threshold"`
	EnableCompositeTimeline              types.Bool  `tfsdk:"enable_composite_timeline"`
}

// knowledgebaseSettingsModel holds both global KB and incident prediction settings
type knowledgebaseSettingsModel struct {
	// Global KB fields
	EnableGlobalKnowledgeBase      types.Bool   `tfsdk:"enable_global_knowledge_base"`
	SatelliteSystemSet             types.String `tfsdk:"satellite_system_set"`
	CompositeValidThreshold        types.Int64  `tfsdk:"composite_valid_threshold"`
	TimelineTopK                   types.Int64  `tfsdk:"timeline_top_k"`
	EnableIgnoreInstancePrediction types.Bool   `tfsdk:"enable_ignore_instance_prediction"`
	PredictionSource               types.Int64  `tfsdk:"prediction_source"`
	ShareSystemType                types.Int64  `tfsdk:"share_system_type"`
	ActionExecutionTime            types.Int64  `tfsdk:"action_execution_time"`
	AutoFixValidationWindow        types.Int64  `tfsdk:"auto_fix_validation_window"`
	FilterSelfToSelf               types.Bool   `tfsdk:"filter_self_to_self"`
	RuleSourceType                 types.Int64  `tfsdk:"rule_source_type"`
	// Incident prediction fields
	RuleActiveThreshold           types.Float64 `tfsdk:"rule_active_threshold"`
	RuleInactiveThreshold         types.Float64 `tfsdk:"rule_inactive_threshold"`
	RuleActiveCondition           types.Int64   `tfsdk:"rule_active_condition"`
	FalsePositiveTolerance        types.Int64   `tfsdk:"false_positive_tolerance"`
	KBTrainingLength              types.Int64   `tfsdk:"kb_training_length"`
	Tolerance                     types.Float64 `tfsdk:"tolerance"`
	EnableInsensitiveRuleMatching types.Bool    `tfsdk:"enable_insensitive_rule_matching"`
}

// notificationsSettingsModel holds notification/health view settings
type notificationsSettingsModel struct {
	Order                               types.Int64   `tfsdk:"order"`
	HideFlag                            types.Bool    `tfsdk:"hide_flag"`
	AggregationInterval                 types.Int64   `tfsdk:"aggregation_interval"`
	EnableSplunkExport                  types.Bool    `tfsdk:"enable_splunk_export"`
	IncidentCountThreshold              types.String  `tfsdk:"incident_count_threshold"`
	AssignmentMap                       types.String  `tfsdk:"assignment_map"`
	PredictionEmail                     types.String  `tfsdk:"prediction_email"`
	AlertHealthScore                    types.Float64 `tfsdk:"alert_health_score"`
	AlertFrequency                      types.Int64   `tfsdk:"alert_frequency"`
	EmailDampeningPeriod                types.Int64   `tfsdk:"email_dampening_period"`
	AlertsEmailDampeningPeriod          types.Int64   `tfsdk:"alerts_email_dampening_period"`
	PredictionEmailDampeningPeriod      types.Int64   `tfsdk:"prediction_email_dampening_period"`
	EnableSystemDownEmailAlert          types.Bool    `tfsdk:"enable_system_down_email_alert"`
	OnlySendWithRCA                     types.Bool    `tfsdk:"only_send_with_rca"`
	EnableIncidentPredictionEmailAlert  types.Bool    `tfsdk:"enable_incident_prediction_email_alert"`
	EnableIncidentDetectionEmailAlert   types.Bool    `tfsdk:"enable_incident_detection_email_alert"`
	EnableAlertsEmail                   types.Bool    `tfsdk:"enable_alerts_email"`
	EnableHealthEmailAlert              types.Bool    `tfsdk:"enable_health_email_alert"`
	AlertEmail                          types.String  `tfsdk:"alert_email"`
	HealthAlertEmail                    types.String  `tfsdk:"health_alert_email"`
	IncidentDetectionEmail              types.String  `tfsdk:"incident_detection_email"`
	EnableRootCauseEmailAlert           types.Bool    `tfsdk:"enable_root_cause_email_alert"`
	RootCauseEmail                      types.String  `tfsdk:"root_cause_email"`
	IncidentDampeningWindow             types.Int64   `tfsdk:"incident_dampening_window"`
	TicketOpenTime                      types.Int64   `tfsdk:"ticket_open_time"`
	ComponentLevelIncidentConsolidation types.Bool    `tfsdk:"component_level_incident_consolidation"`
	ComponentLevelDampening             types.Bool    `tfsdk:"component_level_dampening"`
	EnabledConsolidationAlgorithms      types.List    `tfsdk:"enabled_consolidation_algorithms"`
	// New notification settings
	SystemDownNotification        *systemDownNotificationModel        `tfsdk:"system_down_notification"`
	DailyReportNotification       *insightsReportNotificationModel    `tfsdk:"daily_report_notification"`
	WeeklyReportNotification      *insightsReportNotificationModel    `tfsdk:"weekly_report_notification"`
	InstanceDownNotification      []instanceDownNotificationModel     `tfsdk:"instance_down_notification"`
	ProjectLevelDampeningWindows  []projectLevelDampeningWindowModel  `tfsdk:"project_level_dampening_windows"`
	MaxNotificationDelayTolerance types.Int64                         `tfsdk:"max_notification_delay_tolerance"`
	CustomConsolidationRules      []customConsolidationRuleModel      `tfsdk:"custom_consolidation_rules"`
	MetricLogConsolidationConfigs []metricLogConsolidationConfigModel `tfsdk:"metric_log_consolidation_configs"`
	MetricCoOccurrenceBufferMs    types.Int64                         `tfsdk:"metric_co_occurrence_buffer_ms"`
}

// systemDownNotificationModel holds system down notification settings
type systemDownNotificationModel struct {
	EnableSystemDownEmailAlert types.Bool  `tfsdk:"enable_system_down_email_alert"`
	EmailDampeningPeriod       types.Int64 `tfsdk:"email_dampening_period"`
	EmailSet                   types.List  `tfsdk:"email_set"`
}

// insightsReportNotificationModel holds daily or weekly insights report notification settings
type insightsReportNotificationModel struct {
	EnableInsightsReport types.Bool `tfsdk:"enable_insights_report"`
	EmailSet             types.List `tfsdk:"email_set"`
}

// instanceDownNotificationModel holds instance down notification settings for one project
type instanceDownNotificationModel struct {
	ProjectName              types.String `tfsdk:"project_name"`
	InstanceDownEnable       types.Bool   `tfsdk:"instance_down_enable"`
	InstanceDownDampening    types.Int64  `tfsdk:"instance_down_dampening"`
	InstanceDownThreshold    types.Int64  `tfsdk:"instance_down_threshold"`
	InstanceDownReportNumber types.Int64  `tfsdk:"instance_down_report_number"`
	InstanceDownEmails       types.List   `tfsdk:"instance_down_emails"`
}

// projectLevelDampeningWindowModel holds a single project-level dampening window entry
type projectLevelDampeningWindowModel struct {
	SourceProject  types.String  `tfsdk:"source_project"`
	TargetProject  types.String  `tfsdk:"target_project"`
	SourceCustomer types.String  `tfsdk:"source_customer"`
	TargetCustomer types.String  `tfsdk:"target_customer"`
	Duration       types.Int64   `tfsdk:"duration"`
	ScoreThreshold types.Float64 `tfsdk:"score_threshold"`
}

// conditionModel holds a single matching condition in a custom consolidation rule project entry
type conditionModel struct {
	Type    types.String `tfsdk:"type"`
	Keyword types.String `tfsdk:"keyword"`
}

// projectEntryModel holds a project and its matching conditions in a custom consolidation rule
type projectEntryModel struct {
	ProjectName types.String     `tfsdk:"project_name"`
	Conditions  []conditionModel `tfsdk:"conditions"`
}

// projectFieldKeyModel holds one project's field key in a field correlation
type projectFieldKeyModel struct {
	ProjectName types.String `tfsdk:"project_name"`
	Type        types.String `tfsdk:"type"`
	FieldKey    types.String `tfsdk:"field_key"`
}

// fieldCorrelationModel holds a set of correlated field keys across projects
type fieldCorrelationModel struct {
	ProjectFieldKeys []projectFieldKeyModel `tfsdk:"project_field_keys"`
}

// customConsolidationRuleModel holds one custom incident consolidation rule
type customConsolidationRuleModel struct {
	ProjectEntries    []projectEntryModel     `tfsdk:"project_entries"`
	FieldCorrelations []fieldCorrelationModel `tfsdk:"field_correlations"`
}

// metricLogConsolidationConfigModel holds a metric-to-log project consolidation mapping
type metricLogConsolidationConfigModel struct {
	MetricProjectName types.String `tfsdk:"metric_project_name"`
	LogProjectName    types.String `tfsdk:"log_project_name"`
	FieldKeys         types.List   `tfsdk:"field_keys"`
}

// Metadata returns the resource type name.
func (r *systemSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_settings"
}

// Schema defines the schema for the resource.
func (r *systemSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages InsightFinder system-level settings including knowledge base and notifications configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the system settings (system_name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system_name": schema.StringAttribute{
				Description: "The display name of the system. Used to resolve the system ID via the system framework API.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"knowledgebase_settings": schema.SingleNestedAttribute{
				Description: "Knowledge base and incident prediction settings for the system.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					// Global KB fields
					"enable_global_knowledge_base": schema.BoolAttribute{
						Description: "Enable global knowledge base for the system.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"composite_valid_threshold": schema.Int64Attribute{
						Description: "Composite valid threshold in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"timeline_top_k": schema.Int64Attribute{
						Description: "Number of top timeline entries to keep.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"enable_ignore_instance_prediction": schema.BoolAttribute{
						Description: "Enable ignoring instance prediction in KB.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"prediction_source": schema.Int64Attribute{
						Description: "Prediction source type (0 = default).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"share_system_type": schema.Int64Attribute{
						Description: "Share system type for KB.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"action_execution_time": schema.Int64Attribute{
						Description: "Action execution time in minutes.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"auto_fix_validation_window": schema.Int64Attribute{
						Description: "Auto fix validation window.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"filter_self_to_self": schema.BoolAttribute{
						Description: "Filter self-to-self KB entries.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"rule_source_type": schema.Int64Attribute{
						Description: "Rule source type (0 = default).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"satellite_system_set": schema.StringAttribute{
						Description: "JSON array of satellite systems linked to this system's knowledge base. " +
							"Each entry has systemPartitionKey (userName, systemName, envName) and replay fields. " +
							"Example: jsonencode([{systemPartitionKey={userName=\"u\",systemName=\"<id>\",envName=\"All\"},replay=false}])",
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					// Incident prediction fields
					"rule_active_threshold": schema.Float64Attribute{
						Description: "Threshold to activate a prediction rule (0.0 - 1.0).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"rule_inactive_threshold": schema.Float64Attribute{
						Description: "Threshold to deactivate a prediction rule (0.0 - 1.0).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"rule_active_condition": schema.Int64Attribute{
						Description: "Condition for rule activation (0 = default).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"false_positive_tolerance": schema.Int64Attribute{
						Description: "Number of false positives tolerated before deactivating a rule.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"kb_training_length": schema.Int64Attribute{
						Description: "Length of KB training window in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"tolerance": schema.Float64Attribute{
						Description: "Tolerance value for incident prediction.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"enable_insensitive_rule_matching": schema.BoolAttribute{
						Description: "Enable insensitive rule matching for KB.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"notifications_settings": schema.SingleNestedAttribute{
				Description: "Notification and alert email settings for the system.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"order": schema.Int64Attribute{
						Description: "Display order for the system in the health view.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"hide_flag": schema.BoolAttribute{
						Description: "Hide this system from the health view.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"aggregation_interval": schema.Int64Attribute{
						Description: "Aggregation interval in minutes for health view metrics.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"enable_splunk_export": schema.BoolAttribute{
						Description: "Enable exporting data to Splunk.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"incident_count_threshold": schema.StringAttribute{
						Description: "JSON map of project names to incident count thresholds, e.g. {\"MyProject@user\": 5}.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"assignment_map": schema.StringAttribute{
						Description: "JSON map of zone/component keys to assignee lists (jiraAssignees, emailAssignees, serviceNowAssignees).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"prediction_email": schema.StringAttribute{
						Description: "Email address for prediction notifications.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"alert_health_score": schema.Float64Attribute{
						Description: "Health score threshold for triggering alerts.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"alert_frequency": schema.Int64Attribute{
						Description: "Alert frequency setting.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"email_dampening_period": schema.Int64Attribute{
						Description: "Dampening period for health alert emails in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"alerts_email_dampening_period": schema.Int64Attribute{
						Description: "Dampening period for alert emails in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"prediction_email_dampening_period": schema.Int64Attribute{
						Description: "Dampening period for prediction emails in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"enable_system_down_email_alert": schema.BoolAttribute{
						Description: "Enable email alert when the system is down.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"only_send_with_rca": schema.BoolAttribute{
						Description: "Only send notifications when root cause analysis is available.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_incident_prediction_email_alert": schema.BoolAttribute{
						Description: "Enable email alert for incident predictions.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_incident_detection_email_alert": schema.BoolAttribute{
						Description: "Enable email alert for incident detections.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_alerts_email": schema.BoolAttribute{
						Description: "Enable alert emails.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_health_email_alert": schema.BoolAttribute{
						Description: "Enable health score email alerts.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"alert_email": schema.StringAttribute{
						Description: "Email address for alert notifications.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"health_alert_email": schema.StringAttribute{
						Description: "Email address for health alert notifications.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"incident_detection_email": schema.StringAttribute{
						Description: "Email address for incident detection notifications.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_root_cause_email_alert": schema.BoolAttribute{
						Description: "Enable email alert for root cause analysis results.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"root_cause_email": schema.StringAttribute{
						Description: "Email address for root cause analysis notifications.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"incident_dampening_window": schema.Int64Attribute{
						Description: "Dampening window for incident notifications in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"ticket_open_time": schema.Int64Attribute{
						Description: "Time window in milliseconds for keeping a ticket open.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"component_level_incident_consolidation": schema.BoolAttribute{
						Description: "Enable component-level incident consolidation.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"component_level_dampening": schema.BoolAttribute{
						Description: "Enable component-level dampening.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enabled_consolidation_algorithms": schema.ListAttribute{
						Description: "List of consolidation algorithms to enable (e.g. [\"derivedIncidents\", \"rcaChain\", \"contentBased\", \"metricInstanceTimestamp\"]).",
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
					},
					"system_down_notification": schema.SingleNestedAttribute{
						Description: "System down notification settings (dedicated API at /api/external/v2/systemdownsetting).",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"enable_system_down_email_alert": schema.BoolAttribute{
								Description: "Enable email alert when the system is down.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"email_dampening_period": schema.Int64Attribute{
								Description: "Dampening period for system down email alerts in milliseconds.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"email_set": schema.ListAttribute{
								Description: "List of email addresses to notify when the system is down.",
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
						},
					},
					"daily_report_notification": schema.SingleNestedAttribute{
						Description: "Daily insights report notification settings.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"enable_insights_report": schema.BoolAttribute{
								Description: "Enable daily insights report email.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"email_set": schema.ListAttribute{
								Description: "List of email addresses to receive the daily insights report.",
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
						},
					},
					"weekly_report_notification": schema.SingleNestedAttribute{
						Description: "Weekly insights report notification settings.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"enable_insights_report": schema.BoolAttribute{
								Description: "Enable weekly insights report email.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"email_set": schema.ListAttribute{
								Description: "List of email addresses to receive the weekly insights report.",
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
						},
					},
					"instance_down_notification": schema.ListNestedAttribute{
						Description: "Instance down notification settings per project.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"project_name": schema.StringAttribute{
									Description: "The name of the project to configure instance down notifications for.",
									Required:    true,
								},
								"instance_down_enable": schema.BoolAttribute{
									Description: "Enable instance down email alerts for this project.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"instance_down_dampening": schema.Int64Attribute{
									Description: "Dampening period for instance down alerts in milliseconds.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"instance_down_threshold": schema.Int64Attribute{
									Description: "Threshold (in milliseconds) before an instance is considered down.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"instance_down_report_number": schema.Int64Attribute{
									Description: "Number of instance down events to include in the report.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"instance_down_emails": schema.ListAttribute{
									Description: "List of email addresses to notify when instances are down.",
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
								},
							},
						},
					},
					"project_level_dampening_windows": schema.SetNestedAttribute{
						Description: "Project-level dampening window rules. Each rule overrides the system-level incident dampening window for a specific source→target project pair.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source_project": schema.StringAttribute{
									Description: "The source project name.",
									Required:    true,
								},
								"target_project": schema.StringAttribute{
									Description: "The target project name.",
									Required:    true,
								},
								"source_customer": schema.StringAttribute{
									Description: "The customer (username) of the source project. Defaults to the provider username.",
									Optional:    true,
									Computed:    true,
								},
								"target_customer": schema.StringAttribute{
									Description: "The customer (username) of the target project. Defaults to the provider username.",
									Optional:    true,
									Computed:    true,
								},
								"duration": schema.Int64Attribute{
									Description: "Dampening duration in milliseconds.",
									Required:    true,
								},
								"score_threshold": schema.Float64Attribute{
									Description: "Score threshold (st) for this dampening window.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Float64{
										float64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"max_notification_delay_tolerance": schema.Int64Attribute{
						Description: "Maximum notification delay tolerance in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"metric_co_occurrence_buffer_ms": schema.Int64Attribute{
						Description: "Metric co-occurrence buffer window in milliseconds.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"custom_consolidation_rules": schema.ListNestedAttribute{
						Description: "Custom incident consolidation rules.",
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"project_entries": schema.ListNestedAttribute{
									Description: "Projects and their matching conditions for this consolidation rule.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"project_name": schema.StringAttribute{
												Description: "The project name.",
												Required:    true,
											},
											"conditions": schema.ListNestedAttribute{
												Description: "Matching conditions for this project entry.",
												Optional:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"type": schema.StringAttribute{
															Description: "Condition type: \"fieldName\" or \"content\".",
															Required:    true,
														},
														"keyword": schema.StringAttribute{
															Description: "The keyword or field expression to match.",
															Required:    true,
														},
													},
												},
											},
										},
									},
								},
								"field_correlations": schema.ListNestedAttribute{
									Description: "Field correlations mapping project fields across the consolidation rule.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"project_field_keys": schema.ListNestedAttribute{
												Description: "List of project field key mappings.",
												Required:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"project_name": schema.StringAttribute{
															Description: "The project name.",
															Required:    true,
														},
														"type": schema.StringAttribute{
															Description: "Field type: \"fieldName\" or \"content\".",
															Required:    true,
														},
														"field_key": schema.StringAttribute{
															Description: "The field key path. Empty or omitted for content type.",
															Optional:    true,
															Computed:    true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"metric_log_consolidation_configs": schema.ListNestedAttribute{
						Description: "Metric-to-log project consolidation configuration.",
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"metric_project_name": schema.StringAttribute{
									Description: "The metric project name.",
									Required:    true,
								},
								"log_project_name": schema.StringAttribute{
									Description: "The log project name.",
									Required:    true,
								},
								"field_keys": schema.ListAttribute{
									Description: "List of field keys used for consolidation.",
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"miscellaneous_settings": schema.SingleNestedAttribute{
				Description: "Miscellaneous system framework settings.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"healthview_longterm": schema.BoolAttribute{
						Description: "Enable long-term storage mode for the system health view.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"should_auto_share": schema.BoolAttribute{
						Description: "Enable automatic sharing of system data.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"rootcause_reverse_entry_filter_threshold": schema.Int64Attribute{
						Description: "Threshold for root cause reverse entry filter (0-100).",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"enable_composite_timeline": schema.BoolAttribute{
						Description: "Enable composite timeline view for the system.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *systemSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// resolveSystemID resolves a system name to its ID using the system framework API
func (r *systemSettingsResource) resolveSystemID(ctx context.Context, systemName string) (string, error) {
	ids, err := r.client.ResolveSystemNameToIDs([]string{systemName}, r.client.Username)
	if err != nil {
		return "", fmt.Errorf("failed to resolve system name '%s' to ID: %w", systemName, err)
	}
	if len(ids) == 0 {
		return "", fmt.Errorf("system '%s' not found", systemName)
	}
	return ids[0], nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan systemSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemName := plan.SystemName.ValueString()
	tflog.Debug(ctx, "Creating system settings", map[string]interface{}{"system_name": systemName})

	systemID, err := r.resolveSystemID(ctx, systemName)
	if err != nil {
		resp.Diagnostics.AddError("Error Resolving System ID", err.Error())
		return
	}

	if plan.KnowledgebaseSettings != nil {
		if err := r.applyKnowledgebaseSettings(ctx, systemID, plan.KnowledgebaseSettings); err != nil {
			resp.Diagnostics.AddError("Error Setting Knowledge Base Settings", err.Error())
			return
		}
	}

	if plan.NotificationsSettings != nil {
		if err := r.applyNotificationsSettings(ctx, systemID, plan.NotificationsSettings); err != nil {
			resp.Diagnostics.AddError("Error Setting Notifications Settings", err.Error())
			return
		}
	}

	if plan.MiscellaneousSettings != nil {
		if err := r.applyMiscellaneousSettings(ctx, systemID, plan.MiscellaneousSettings); err != nil {
			resp.Diagnostics.AddError("Error Setting Miscellaneous Settings", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(systemName)

	// Read back state to populate computed fields
	if err := r.readIntoModel(ctx, systemID, &plan); err != nil {
		resp.Diagnostics.AddError("Error Reading System Settings After Create", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *systemSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state systemSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemName := state.SystemName.ValueString()
	tflog.Debug(ctx, "Reading system settings", map[string]interface{}{"system_name": systemName})

	systemID, err := r.resolveSystemID(ctx, systemName)
	if err != nil {
		resp.Diagnostics.AddError("Error Resolving System ID", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, systemID, &state); err != nil {
		resp.Diagnostics.AddError("Error Reading System Settings", err.Error())
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan systemSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemName := plan.SystemName.ValueString()
	tflog.Debug(ctx, "Updating system settings", map[string]interface{}{"system_name": systemName})

	systemID, err := r.resolveSystemID(ctx, systemName)
	if err != nil {
		resp.Diagnostics.AddError("Error Resolving System ID", err.Error())
		return
	}

	if plan.KnowledgebaseSettings != nil {
		if err := r.applyKnowledgebaseSettings(ctx, systemID, plan.KnowledgebaseSettings); err != nil {
			resp.Diagnostics.AddError("Error Updating Knowledge Base Settings", err.Error())
			return
		}
	}

	if plan.NotificationsSettings != nil {
		if err := r.applyNotificationsSettings(ctx, systemID, plan.NotificationsSettings); err != nil {
			resp.Diagnostics.AddError("Error Updating Notifications Settings", err.Error())
			return
		}
	}

	if plan.MiscellaneousSettings != nil {
		if err := r.applyMiscellaneousSettings(ctx, systemID, plan.MiscellaneousSettings); err != nil {
			resp.Diagnostics.AddError("Error Updating Miscellaneous Settings", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(systemName)

	// Read back state to populate computed fields
	if err := r.readIntoModel(ctx, systemID, &plan); err != nil {
		resp.Diagnostics.AddError("Error Reading System Settings After Update", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete removes the resource from Terraform state (settings are left as-is on the server).
func (r *systemSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// System settings cannot be "deleted" from the API; removing from Terraform state only.
	tflog.Debug(ctx, "Deleting system settings resource (removing from state only)")
}

// ImportState imports the resource state using system_name as the ID.
func (r *systemSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("system_name"), req, resp)
}

// applyKnowledgebaseSettings writes knowledgebase settings to the API
func (r *systemSettingsResource) applyKnowledgebaseSettings(_ context.Context, systemID string, m *knowledgebaseSettingsModel) error {
	globalKB := &client.GlobalKBSetting{
		EnableGlobalKnowledgeBase:      m.EnableGlobalKnowledgeBase.ValueBool(),
		CompositeValidThreshold:        m.CompositeValidThreshold.ValueInt64(),
		TimelineTopK:                   m.TimelineTopK.ValueInt64(),
		EnableIgnoreInstancePrediction: m.EnableIgnoreInstancePrediction.ValueBool(),
		PredictionSource:               m.PredictionSource.ValueInt64(),
		ShareSystemType:                m.ShareSystemType.ValueInt64(),
		ActionExecutionTime:            m.ActionExecutionTime.ValueInt64(),
		AutoFixValidationWindow:        m.AutoFixValidationWindow.ValueInt64(),
		FilterSelfToSelf:               m.FilterSelfToSelf.ValueBool(),
		RuleSourceType:                 m.RuleSourceType.ValueInt64(),
	}

	// Parse satellite_system_set JSON string and convert to the SET-format list
	if v := m.SatelliteSystemSet.ValueString(); v != "" && v != "null" && v != "[]" {
		var setEntries []client.SatelliteSystemSetEntry
		if err := json.Unmarshal([]byte(v), &setEntries); err != nil {
			return fmt.Errorf("satellite_system_set is not valid JSON: %w", err)
		}
		globalKB.SatelliteSystemList = client.SatelliteSystemSetToList(setEntries)
	}

	if err := r.client.SetGlobalKBSetting(systemID, globalKB); err != nil {
		return fmt.Errorf("failed to set global KB setting: %w", err)
	}

	incidentPrediction := &client.IncidentPredictionSetting{
		RuleActiveThreshold:           m.RuleActiveThreshold.ValueFloat64(),
		RuleInactiveThreshold:         m.RuleInactiveThreshold.ValueFloat64(),
		RuleActiveCondition:           m.RuleActiveCondition.ValueInt64(),
		FalsePositiveTolerance:        m.FalsePositiveTolerance.ValueInt64(),
		KBTrainingLength:              m.KBTrainingLength.ValueInt64(),
		Tolerance:                     m.Tolerance.ValueFloat64(),
		EnableInsensitiveRuleMatching: m.EnableInsensitiveRuleMatching.ValueBool(),
	}

	if err := r.client.SetIncidentPredictionSetting(systemID, incidentPrediction); err != nil {
		return fmt.Errorf("failed to set incident prediction setting: %w", err)
	}

	return nil
}

// applyNotificationsSettings writes notifications settings to the API
func (r *systemSettingsResource) applyNotificationsSettings(_ context.Context, systemID string, m *notificationsSettingsModel) error {
	updates := &client.HealthViewSetting{
		Order:                               m.Order.ValueInt64(),
		HideFlag:                            m.HideFlag.ValueBool(),
		AggregationInterval:                 m.AggregationInterval.ValueInt64(),
		EnableSplunkExport:                  m.EnableSplunkExport.ValueBool(),
		PredictionEmail:                     m.PredictionEmail.ValueString(),
		AlertHealthScore:                    m.AlertHealthScore.ValueFloat64(),
		AlertFrequency:                      m.AlertFrequency.ValueInt64(),
		EmailDampeningPeriod:                m.EmailDampeningPeriod.ValueInt64(),
		AlertsEmailDampeningPeriod:          m.AlertsEmailDampeningPeriod.ValueInt64(),
		PredictionEmailDampeningPeriod:      m.PredictionEmailDampeningPeriod.ValueInt64(),
		EnableSystemDownEmailAlert:          m.EnableSystemDownEmailAlert.ValueBool(),
		OnlySendWithRCA:                     m.OnlySendWithRCA.ValueBool(),
		EnableIncidentPredictionEmailAlert:  m.EnableIncidentPredictionEmailAlert.ValueBool(),
		EnableIncidentDetectionEmailAlert:   m.EnableIncidentDetectionEmailAlert.ValueBool(),
		EnableAlertsEmail:                   m.EnableAlertsEmail.ValueBool(),
		EnableHealthEmailAlert:              m.EnableHealthEmailAlert.ValueBool(),
		AlertEmail:                          m.AlertEmail.ValueString(),
		HealthAlertEmail:                    m.HealthAlertEmail.ValueString(),
		IncidentDetectionEmail:              m.IncidentDetectionEmail.ValueString(),
		EnableRootCauseEmailAlert:           m.EnableRootCauseEmailAlert.ValueBool(),
		RootCauseEmail:                      m.RootCauseEmail.ValueString(),
		IncidentDampeningWindow:             m.IncidentDampeningWindow.ValueInt64(),
		TicketOpenTime:                      m.TicketOpenTime.ValueInt64(),
		ComponentLevelIncidentConsolidation: m.ComponentLevelIncidentConsolidation.ValueBool(),
		ComponentLevelDampening:             m.ComponentLevelDampening.ValueBool(),
		EnabledConsolidationAlgorithms:      typesListToStrings(m.EnabledConsolidationAlgorithms),
		MetricCoOccurrenceBufferMs:          m.MetricCoOccurrenceBufferMs.ValueInt64(),
	}

	if v := m.IncidentCountThreshold.ValueString(); v != "" && v != "null" {
		var ict map[string]int64
		if err := json.Unmarshal([]byte(v), &ict); err != nil {
			return fmt.Errorf("incident_count_threshold is not valid JSON: %w", err)
		}
		updates.IncidentCountThreshold = ict
	}

	if v := m.AssignmentMap.ValueString(); v != "" && v != "null" {
		var am map[string]any
		if err := json.Unmarshal([]byte(v), &am); err != nil {
			return fmt.Errorf("assignment_map is not valid JSON: %w", err)
		}
		updates.AssignmentMap = am
	}

	windows := make([]client.ProjectLevelDampeningWindow, 0, len(m.ProjectLevelDampeningWindows))
	for _, w := range m.ProjectLevelDampeningWindows {
		sc := w.SourceCustomer.ValueString()
		if sc == "" {
			sc = r.client.Username
		}
		tc := w.TargetCustomer.ValueString()
		if tc == "" {
			tc = r.client.Username
		}
		windows = append(windows, client.ProjectLevelDampeningWindow{
			SourceProject:  w.SourceProject.ValueString(),
			TargetProject:  w.TargetProject.ValueString(),
			SourceCustomer: sc,
			TargetCustomer: tc,
			Duration:       w.Duration.ValueInt64(),
			ScoreThreshold: w.ScoreThreshold.ValueFloat64(),
		})
	}
	updates.ProjectLevelDampeningWindows = windows

	updates.MaxNotificationDelayTolerance = m.MaxNotificationDelayTolerance.ValueInt64()

	rules := make([]client.CustomConsolidationRule, 0, len(m.CustomConsolidationRules))
	for _, rule := range m.CustomConsolidationRules {
		entries := make([]client.ProjectEntry, 0, len(rule.ProjectEntries))
		for _, pe := range rule.ProjectEntries {
			conds := make([]client.Condition, 0, len(pe.Conditions))
			for _, c := range pe.Conditions {
				conds = append(conds, client.Condition{
					Type:    c.Type.ValueString(),
					Keyword: c.Keyword.ValueString(),
				})
			}
			entries = append(entries, client.ProjectEntry{
				ProjectName: pe.ProjectName.ValueString(),
				Conditions:  conds,
			})
		}
		corrs := make([]client.FieldCorrelation, 0, len(rule.FieldCorrelations))
		for _, fc := range rule.FieldCorrelations {
			keys := make([]client.ProjectFieldKey, 0, len(fc.ProjectFieldKeys))
			for _, pk := range fc.ProjectFieldKeys {
				var fkPtr *string
				if !pk.FieldKey.IsNull() && !pk.FieldKey.IsUnknown() {
					fk := pk.FieldKey.ValueString()
					fkPtr = &fk
				}
				keys = append(keys, client.ProjectFieldKey{
					ProjectName: pk.ProjectName.ValueString(),
					Type:        pk.Type.ValueString(),
					FieldKey:    fkPtr,
				})
			}
			corrs = append(corrs, client.FieldCorrelation{
				ProjectFieldKeys: keys,
			})
		}
		rules = append(rules, client.CustomConsolidationRule{
			ProjectEntries:    entries,
			FieldCorrelations: corrs,
		})
	}
	updates.CustomConsolidationRules = rules

	mlConfigs := make([]client.MetricLogConsolidationConfig, 0, len(m.MetricLogConsolidationConfigs))
	for _, cfg := range m.MetricLogConsolidationConfigs {
		mlConfigs = append(mlConfigs, client.MetricLogConsolidationConfig{
			MetricProjectName: cfg.MetricProjectName.ValueString(),
			LogProjectName:    cfg.LogProjectName.ValueString(),
			FieldKeys:         typesListToStrings(cfg.FieldKeys),
		})
	}
	updates.MetricLogConsolidationConfigs = mlConfigs

	if err := r.client.SetHealthViewSetting(systemID, updates); err != nil {
		return fmt.Errorf("failed to set health view setting: %w", err)
	}

	// System down notification (dedicated API)
	if m.SystemDownNotification != nil {
		sdSetting := &client.SystemDownSetting{
			EnableSystemDownEmailAlert: m.SystemDownNotification.EnableSystemDownEmailAlert.ValueBool(),
			EmailDampeningPeriod:       m.SystemDownNotification.EmailDampeningPeriod.ValueInt64(),
			EmailSet:                   typesListToStrings(m.SystemDownNotification.EmailSet),
		}
		if err := r.client.SetSystemDownSetting(systemID, sdSetting); err != nil {
			return fmt.Errorf("failed to set system down setting: %w", err)
		}
	}

	// Daily report notification
	if m.DailyReportNotification != nil {
		if err := r.client.SetInsightsReportSetting(
			systemID,
			typesListToStrings(m.DailyReportNotification.EmailSet),
			m.DailyReportNotification.EnableInsightsReport.ValueBool(),
			true,
		); err != nil {
			return fmt.Errorf("failed to set daily insights report setting: %w", err)
		}
	}

	// Weekly report notification
	if m.WeeklyReportNotification != nil {
		if err := r.client.SetInsightsReportSetting(
			systemID,
			typesListToStrings(m.WeeklyReportNotification.EmailSet),
			m.WeeklyReportNotification.EnableInsightsReport.ValueBool(),
			false,
		); err != nil {
			return fmt.Errorf("failed to set weekly insights report setting: %w", err)
		}
	}

	// Instance down notifications (one API call per project)
	for _, item := range m.InstanceDownNotification {
		setting := &client.InstanceDownSetting{
			ProjectName:              item.ProjectName.ValueString(),
			InstanceDownEnable:       item.InstanceDownEnable.ValueBool(),
			InstanceDownDampening:    item.InstanceDownDampening.ValueInt64(),
			InstanceDownThreshold:    item.InstanceDownThreshold.ValueInt64(),
			InstanceDownReportNumber: item.InstanceDownReportNumber.ValueInt64(),
			InstanceDownEmails:       typesListToStrings(item.InstanceDownEmails),
		}
		if err := r.client.SetInstanceDownSetting(setting); err != nil {
			return fmt.Errorf("failed to set instance down setting for project %s: %w", setting.ProjectName, err)
		}
	}

	return nil
}

// applyMiscellaneousSettings writes miscellaneous system framework settings to the API
func (r *systemSettingsResource) applyMiscellaneousSettings(_ context.Context, systemID string, m *miscellaneousSettingsModel) error {
	if err := r.client.SetLongTermSetting(systemID, m.HealthviewLongterm.ValueBool()); err != nil {
		return fmt.Errorf("failed to set longTerm setting: %w", err)
	}

	fwSettings := &client.SystemFrameworkSetting{
		ShouldAutoShare:                      m.ShouldAutoShare.ValueBool(),
		RootCauseReverseEntryFilterThreshold: m.RootcauseReverseEntryFilterThreshold.ValueInt64(),
		EnableCompositeTimeline:              m.EnableCompositeTimeline.ValueBool(),
	}
	if err := r.client.SetSystemFrameworkSetting(systemID, fwSettings); err != nil {
		return fmt.Errorf("failed to set system framework setting: %w", err)
	}

	return nil
}

// readIntoModel reads the current state from the API and populates the model
func (r *systemSettingsResource) readIntoModel(_ context.Context, systemID string, m *systemSettingsResourceModel) error {
	if m.KnowledgebaseSettings != nil {
		kbSetting, err := r.client.GetGlobalKBSetting(systemID)
		if err != nil {
			return fmt.Errorf("failed to read global KB setting: %w", err)
		}
		ipSetting, err := r.client.GetIncidentPredictionSetting(systemID)
		if err != nil {
			return fmt.Errorf("failed to read incident prediction setting: %w", err)
		}

		if kbSetting != nil {
			m.KnowledgebaseSettings.EnableGlobalKnowledgeBase = types.BoolValue(kbSetting.EnableGlobalKnowledgeBase)
			m.KnowledgebaseSettings.CompositeValidThreshold = types.Int64Value(kbSetting.CompositeValidThreshold)
			m.KnowledgebaseSettings.TimelineTopK = types.Int64Value(kbSetting.TimelineTopK)
			m.KnowledgebaseSettings.EnableIgnoreInstancePrediction = types.BoolValue(kbSetting.EnableIgnoreInstancePrediction)
			m.KnowledgebaseSettings.PredictionSource = types.Int64Value(kbSetting.PredictionSource)
			m.KnowledgebaseSettings.ShareSystemType = types.Int64Value(kbSetting.ShareSystemType)
			m.KnowledgebaseSettings.ActionExecutionTime = types.Int64Value(kbSetting.ActionExecutionTime)
			m.KnowledgebaseSettings.AutoFixValidationWindow = types.Int64Value(kbSetting.AutoFixValidationWindow)
			m.KnowledgebaseSettings.FilterSelfToSelf = types.BoolValue(kbSetting.FilterSelfToSelf)
			m.KnowledgebaseSettings.RuleSourceType = types.Int64Value(kbSetting.RuleSourceType)

			// Serialize satellite_system_set, preserving the existing state string when
			// the data is semantically equal (avoids key-ordering diffs on every plan/apply).
			var apiSSJSON string
			if len(kbSetting.SatelliteSystemSet) > 0 {
				if b, err := json.Marshal(kbSetting.SatelliteSystemSet); err == nil {
					apiSSJSON = string(b)
				}
			} else {
				apiSSJSON = "[]"
			}
			existing := m.KnowledgebaseSettings.SatelliteSystemSet.ValueString()
			if normalizeJSONString(existing) == normalizeJSONString(apiSSJSON) && existing != "" {
				// Keep the existing state value (preserves user's key order)
			} else {
				m.KnowledgebaseSettings.SatelliteSystemSet = types.StringValue(apiSSJSON)
			}
		}

		if ipSetting != nil {
			m.KnowledgebaseSettings.RuleActiveThreshold = types.Float64Value(ipSetting.RuleActiveThreshold)
			m.KnowledgebaseSettings.RuleInactiveThreshold = types.Float64Value(ipSetting.RuleInactiveThreshold)
			m.KnowledgebaseSettings.RuleActiveCondition = types.Int64Value(ipSetting.RuleActiveCondition)
			m.KnowledgebaseSettings.FalsePositiveTolerance = types.Int64Value(ipSetting.FalsePositiveTolerance)
			m.KnowledgebaseSettings.KBTrainingLength = types.Int64Value(ipSetting.KBTrainingLength)
			m.KnowledgebaseSettings.Tolerance = types.Float64Value(ipSetting.Tolerance)
			m.KnowledgebaseSettings.EnableInsensitiveRuleMatching = types.BoolValue(ipSetting.EnableInsensitiveRuleMatching)
		}
	}

	if m.NotificationsSettings != nil {
		hvSetting, err := r.client.GetHealthViewSetting(systemID)
		if err != nil {
			return fmt.Errorf("failed to read health view setting: %w", err)
		}

		if hvSetting != nil {
			m.NotificationsSettings.Order = types.Int64Value(hvSetting.Order)
			m.NotificationsSettings.HideFlag = types.BoolValue(hvSetting.HideFlag)
			m.NotificationsSettings.AggregationInterval = types.Int64Value(hvSetting.AggregationInterval)
			m.NotificationsSettings.EnableSplunkExport = types.BoolValue(hvSetting.EnableSplunkExport)

			// For JSON map fields, marshal the API value then preserve the existing state
			// string when it is semantically equal (avoids key-ordering diffs).
			var apiICT string
			if hvSetting.IncidentCountThreshold != nil {
				if b, err := json.Marshal(hvSetting.IncidentCountThreshold); err == nil {
					apiICT = string(b)
				}
			} else {
				apiICT = "{}"
			}
			existingICT := m.NotificationsSettings.IncidentCountThreshold.ValueString()
			if normalizeJSONString(existingICT) == normalizeJSONString(apiICT) && existingICT != "" {
				// Keep existing state value
			} else {
				m.NotificationsSettings.IncidentCountThreshold = types.StringValue(apiICT)
			}

			var apiAM string
			if hvSetting.AssignmentMap != nil {
				if b, err := json.Marshal(hvSetting.AssignmentMap); err == nil {
					apiAM = string(b)
				}
			} else {
				apiAM = "{}"
			}
			existingAM := m.NotificationsSettings.AssignmentMap.ValueString()
			if normalizeJSONString(existingAM) == normalizeJSONString(apiAM) && existingAM != "" {
				// Keep existing state value
			} else {
				m.NotificationsSettings.AssignmentMap = types.StringValue(apiAM)
			}

			m.NotificationsSettings.PredictionEmail = types.StringValue(hvSetting.PredictionEmail)
			m.NotificationsSettings.AlertHealthScore = types.Float64Value(hvSetting.AlertHealthScore)
			m.NotificationsSettings.AlertFrequency = types.Int64Value(hvSetting.AlertFrequency)
			m.NotificationsSettings.EmailDampeningPeriod = types.Int64Value(hvSetting.EmailDampeningPeriod)
			m.NotificationsSettings.AlertsEmailDampeningPeriod = types.Int64Value(hvSetting.AlertsEmailDampeningPeriod)
			m.NotificationsSettings.PredictionEmailDampeningPeriod = types.Int64Value(hvSetting.PredictionEmailDampeningPeriod)
			m.NotificationsSettings.EnableSystemDownEmailAlert = types.BoolValue(hvSetting.EnableSystemDownEmailAlert)
			m.NotificationsSettings.OnlySendWithRCA = types.BoolValue(hvSetting.OnlySendWithRCA)
			m.NotificationsSettings.EnableIncidentPredictionEmailAlert = types.BoolValue(hvSetting.EnableIncidentPredictionEmailAlert)
			m.NotificationsSettings.EnableIncidentDetectionEmailAlert = types.BoolValue(hvSetting.EnableIncidentDetectionEmailAlert)
			m.NotificationsSettings.EnableAlertsEmail = types.BoolValue(hvSetting.EnableAlertsEmail)
			m.NotificationsSettings.EnableHealthEmailAlert = types.BoolValue(hvSetting.EnableHealthEmailAlert)
			m.NotificationsSettings.AlertEmail = types.StringValue(hvSetting.AlertEmail)
			m.NotificationsSettings.HealthAlertEmail = types.StringValue(hvSetting.HealthAlertEmail)
			m.NotificationsSettings.IncidentDetectionEmail = types.StringValue(hvSetting.IncidentDetectionEmail)
			m.NotificationsSettings.EnableRootCauseEmailAlert = types.BoolValue(hvSetting.EnableRootCauseEmailAlert)
			m.NotificationsSettings.RootCauseEmail = types.StringValue(hvSetting.RootCauseEmail)
			m.NotificationsSettings.IncidentDampeningWindow = types.Int64Value(hvSetting.IncidentDampeningWindow)
			m.NotificationsSettings.TicketOpenTime = types.Int64Value(hvSetting.TicketOpenTime)
			m.NotificationsSettings.ComponentLevelIncidentConsolidation = types.BoolValue(hvSetting.ComponentLevelIncidentConsolidation)
			m.NotificationsSettings.ComponentLevelDampening = types.BoolValue(hvSetting.ComponentLevelDampening)
			m.NotificationsSettings.EnabledConsolidationAlgorithms = stringsToTypesList(hvSetting.EnabledConsolidationAlgorithms)
			m.NotificationsSettings.MaxNotificationDelayTolerance = types.Int64Value(hvSetting.MaxNotificationDelayTolerance)
			m.NotificationsSettings.MetricCoOccurrenceBufferMs = types.Int64Value(hvSetting.MetricCoOccurrenceBufferMs)

			// Custom consolidation rules
			newRules := make([]customConsolidationRuleModel, 0, len(hvSetting.CustomConsolidationRules))
			for _, rule := range hvSetting.CustomConsolidationRules {
				entries := make([]projectEntryModel, 0, len(rule.ProjectEntries))
				for _, pe := range rule.ProjectEntries {
					// Use nil (not empty slice) for empty conditions so Terraform
					// state stays null when the user omits the field in config.
					var conds []conditionModel
					if len(pe.Conditions) > 0 {
						conds = make([]conditionModel, 0, len(pe.Conditions))
						for _, c := range pe.Conditions {
							conds = append(conds, conditionModel{
								Type:    types.StringValue(c.Type),
								Keyword: types.StringValue(c.Keyword),
							})
						}
					}
					entries = append(entries, projectEntryModel{
						ProjectName: types.StringValue(pe.ProjectName),
						Conditions:  conds,
					})
				}
				corrs := make([]fieldCorrelationModel, 0, len(rule.FieldCorrelations))
				for _, fc := range rule.FieldCorrelations {
					keys := make([]projectFieldKeyModel, 0, len(fc.ProjectFieldKeys))
					for _, pk := range fc.ProjectFieldKeys {
						fk := types.StringNull()
						if pk.FieldKey != nil {
							fk = types.StringValue(*pk.FieldKey)
						}
						keys = append(keys, projectFieldKeyModel{
							ProjectName: types.StringValue(pk.ProjectName),
							Type:        types.StringValue(pk.Type),
							FieldKey:    fk,
						})
					}
					corrs = append(corrs, fieldCorrelationModel{
						ProjectFieldKeys: keys,
					})
				}
				newRules = append(newRules, customConsolidationRuleModel{
					ProjectEntries:    entries,
					FieldCorrelations: corrs,
				})
			}
			m.NotificationsSettings.CustomConsolidationRules = newRules

			// Metric-log consolidation configs
			newMLConfigs := make([]metricLogConsolidationConfigModel, 0, len(hvSetting.MetricLogConsolidationConfigs))
			for _, cfg := range hvSetting.MetricLogConsolidationConfigs {
				newMLConfigs = append(newMLConfigs, metricLogConsolidationConfigModel{
					MetricProjectName: types.StringValue(cfg.MetricProjectName),
					LogProjectName:    types.StringValue(cfg.LogProjectName),
					FieldKeys:         stringsToTypesList(cfg.FieldKeys),
				})
			}
			m.NotificationsSettings.MetricLogConsolidationConfigs = newMLConfigs

			if m.NotificationsSettings.ProjectLevelDampeningWindows != nil {
				newWindows := make([]projectLevelDampeningWindowModel, 0, len(hvSetting.ProjectLevelDampeningWindows))
				for _, w := range hvSetting.ProjectLevelDampeningWindows {
					newWindows = append(newWindows, projectLevelDampeningWindowModel{
						SourceProject:  types.StringValue(w.SourceProject),
						TargetProject:  types.StringValue(w.TargetProject),
						SourceCustomer: types.StringValue(w.SourceCustomer),
						TargetCustomer: types.StringValue(w.TargetCustomer),
						Duration:       types.Int64Value(w.Duration),
						ScoreThreshold: types.Float64Value(w.ScoreThreshold),
					})
				}
				m.NotificationsSettings.ProjectLevelDampeningWindows = newWindows
			}
		}

		// System down notification
		if m.NotificationsSettings.SystemDownNotification != nil {
			sdSetting, err := r.client.GetSystemDownSetting(systemID)
			if err != nil {
				return fmt.Errorf("failed to read system down setting: %w", err)
			}
			if sdSetting != nil {
				m.NotificationsSettings.SystemDownNotification.EnableSystemDownEmailAlert = types.BoolValue(sdSetting.EnableSystemDownEmailAlert)
				m.NotificationsSettings.SystemDownNotification.EmailDampeningPeriod = types.Int64Value(sdSetting.EmailDampeningPeriod)
				m.NotificationsSettings.SystemDownNotification.EmailSet = stringsToTypesList(sdSetting.EmailSet)
			}
		}

		// Daily and weekly report notifications (single GET returns both)
		if m.NotificationsSettings.DailyReportNotification != nil || m.NotificationsSettings.WeeklyReportNotification != nil {
			irSetting, err := r.client.GetInsightsReportSetting(systemID)
			if err != nil {
				return fmt.Errorf("failed to read insights report setting: %w", err)
			}
			if irSetting != nil {
				if m.NotificationsSettings.DailyReportNotification != nil {
					m.NotificationsSettings.DailyReportNotification.EnableInsightsReport = types.BoolValue(irSetting.EnableDailyInsightsReport)
					m.NotificationsSettings.DailyReportNotification.EmailSet = stringsToTypesList(irSetting.EmailSet)
				}
				if m.NotificationsSettings.WeeklyReportNotification != nil {
					m.NotificationsSettings.WeeklyReportNotification.EnableInsightsReport = types.BoolValue(irSetting.EnableWeeklyInsightsReport)
					m.NotificationsSettings.WeeklyReportNotification.EmailSet = stringsToTypesList(irSetting.WeeklyEmailSet)
				}
			}
		}

		// Instance down notifications (one GET per project)
		if len(m.NotificationsSettings.InstanceDownNotification) > 0 {
			newList := make([]instanceDownNotificationModel, 0, len(m.NotificationsSettings.InstanceDownNotification))
			for _, item := range m.NotificationsSettings.InstanceDownNotification {
				projectName := item.ProjectName.ValueString()
				idSetting, err := r.client.GetInstanceDownSetting(projectName)
				if err != nil {
					return fmt.Errorf("failed to read instance down setting for project %s: %w", projectName, err)
				}
				if idSetting != nil {
					newList = append(newList, instanceDownNotificationModel{
						ProjectName:              types.StringValue(idSetting.ProjectName),
						InstanceDownEnable:       types.BoolValue(idSetting.InstanceDownEnable),
						InstanceDownDampening:    types.Int64Value(idSetting.InstanceDownDampening),
						InstanceDownThreshold:    types.Int64Value(idSetting.InstanceDownThreshold),
						InstanceDownReportNumber: types.Int64Value(idSetting.InstanceDownReportNumber),
						InstanceDownEmails:       stringsToTypesList(idSetting.InstanceDownEmails),
					})
				} else {
					newList = append(newList, item)
				}
			}
			m.NotificationsSettings.InstanceDownNotification = newList
		}
	}

	if m.MiscellaneousSettings != nil {
		miscSetting, err := r.client.GetMiscellaneousSettings(systemID)
		if err != nil {
			return fmt.Errorf("failed to read miscellaneous settings: %w", err)
		}
		if miscSetting != nil {
			m.MiscellaneousSettings.HealthviewLongterm = types.BoolValue(miscSetting.LongTerm)
			m.MiscellaneousSettings.ShouldAutoShare = types.BoolValue(miscSetting.ShouldAutoShare)
			m.MiscellaneousSettings.RootcauseReverseEntryFilterThreshold = types.Int64Value(miscSetting.RootCauseReverseEntryFilterThreshold)
			m.MiscellaneousSettings.EnableCompositeTimeline = types.BoolValue(miscSetting.EnableCompositeTimeline)
		}
	}

	return nil
}

// stringsToTypesList converts a Go []string to a types.List of string elements.
func stringsToTypesList(strs []string) types.List {
	if len(strs) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}
	vals := make([]attr.Value, len(strs))
	for i, s := range strs {
		vals[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, vals)
}

// typesListToStrings converts a types.List of string elements to a Go []string.
func typesListToStrings(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return []string{}
	}
	result := make([]string, 0, len(l.Elements()))
	for _, el := range l.Elements() {
		if s, ok := el.(types.String); ok {
			result = append(result, s.ValueString())
		}
	}
	return result
}

// normalizeJSONString parses and re-marshals a JSON string via interface{} so that
// semantically equivalent JSON with different key orderings produces the same bytes.
// This is the same approach used by resource_project.go's normalizeJSON helper.
func normalizeJSONString(s string) string {
	if s == "" {
		return s
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return s
	}
	return string(b)
}
