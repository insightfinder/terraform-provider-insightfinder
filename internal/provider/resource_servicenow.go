// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &servicenowResource{}
	_ resource.ResourceWithConfigure   = &servicenowResource{}
	_ resource.ResourceWithImportState = &servicenowResource{}
)

// NewServiceNowResource is a helper function to simplify the provider implementation.
func NewServiceNowResource() resource.Resource {
	return &servicenowResource{}
}

// servicenowResource is the resource implementation.
type servicenowResource struct {
	client *client.Client
}

// projectConfigModel maps per-project ServiceNow ticket settings.
type projectConfigModel struct {
	EnableTicketCreation                  types.Bool   `tfsdk:"enable_ticket_creation"`
	EnableTicketUpdate                    types.Bool   `tfsdk:"enable_ticket_update"`
	EnableIncidentConsolidationInfoUpdate types.Bool   `tfsdk:"enable_incident_consolidation_info_update"`
	EnableIncidentResolveUpdate           types.Bool   `tfsdk:"enable_incident_resolve_update"`
	EnableIncidentFieldSync               types.Bool   `tfsdk:"enable_incident_field_sync"`
	EnableMetricValueUpdate               types.Bool   `tfsdk:"enable_metric_value_update"`
	ConfigurationItem                     types.String `tfsdk:"configuration_item"`
}

// resolutionCodeRuleModel maps a single pattern-based resolution code classification rule.
type resolutionCodeRuleModel struct {
	Pattern types.String `tfsdk:"pattern"`
	Outcome types.String `tfsdk:"outcome"`
}

// servicenowResourceModel maps the resource schema data.
type servicenowResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Account                    types.String `tfsdk:"account"`
	ServiceHost                types.String `tfsdk:"service_host"`
	Password                   types.String `tfsdk:"password"`
	Proxy                      types.String `tfsdk:"proxy"`
	AppID                      types.String `tfsdk:"app_id"`
	AppKey                     types.String `tfsdk:"app_key"`
	AuthType                   types.String `tfsdk:"auth_type"`
	SystemNames                types.Set    `tfsdk:"system_names"`
	Options                    types.Set    `tfsdk:"options"`
	ContentOption              types.Set    `tfsdk:"content_option"`
	ServiceNowField            types.String `tfsdk:"service_now_field"`
	ContentSource              types.String `tfsdk:"content_source"`
	TriggerWindowInMills       types.Int64  `tfsdk:"trigger_window_in_mills"`
	EnableFeedbackCollect      types.Bool   `tfsdk:"enable_feedback_collect"`
	TicketCreatedBySourceKey   types.String `tfsdk:"ticket_created_by_source_key"`
	TicketCreatedBySourceValue types.String `tfsdk:"ticket_created_by_source_value"`
	ConfigurationItem          types.String `tfsdk:"configuration_item"`
	DepartmentID               types.String `tfsdk:"department_id"`
	ProjectConfigs             types.Map    `tfsdk:"project_configs"`
	TableMapping               types.Map    `tfsdk:"table_mapping"`
	ResolutionCodeRules        types.List   `tfsdk:"resolution_code_rules"`
}

// Metadata returns the resource type name.
func (r *servicenowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicenow"
}

// Schema defines the schema for the resource.
func (r *servicenowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages InsightFinder ServiceNow integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the ServiceNow configuration (account@service_host).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account": schema.StringAttribute{
				Description: "ServiceNow account username.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_host": schema.StringAttribute{
				Description: "ServiceNow service host URL.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "ServiceNow account password.",
				Required:    true,
				Sensitive:   true,
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy server URL (optional).",
				Optional:    true,
				Computed:    true,
			},
			"app_id": schema.StringAttribute{
				Description: "ServiceNow application ID (optional).",
				Optional:    true,
			},
			"app_key": schema.StringAttribute{
				Description: "ServiceNow application key (optional).",
				Optional:    true,
				Sensitive:   true,
			},
			"auth_type": schema.StringAttribute{
				Description: "Authentication type to use when connecting to ServiceNow. Must be 'basic' or 'oauth'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("basic"),
			},
			"system_names": schema.SetAttribute{
				Description: "Set of system names to integrate.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"options": schema.SetAttribute{
				Description: "ServiceNow integration options.",
				Required:    true,
				ElementType: types.StringType,
			},
			"content_option": schema.SetAttribute{
				Description: "ServiceNow content options.",
				Required:    true,
				ElementType: types.StringType,
			},
			"service_now_field": schema.StringAttribute{
				Description: "ServiceNow field to write integration content to (e.g., 'u_probable_cause').",
				Optional:    true,
			},
			"content_source": schema.StringAttribute{
				Description: "ServiceNow content source field (e.g., 'work_notes'). Defaults to 'work_notes'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("work_notes"),
			},
			"trigger_window_in_mills": schema.Int64Attribute{
				Description: "Trigger window in milliseconds for ServiceNow integration (e.g., 604800000 for 7 days).",
				Optional:    true,
			},
			"enable_feedback_collect": schema.BoolAttribute{
				Description: "Whether to enable ServiceNow feedback collection.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ticket_created_by_source_key": schema.StringAttribute{
				Description: "ServiceNow field key used to filter when a ticket is created (e.g., 'activity_due').",
				Optional:    true,
			},
			"ticket_created_by_source_value": schema.StringAttribute{
				Description: "Value for the ticket_created_by_source_key field filter.",
				Optional:    true,
			},
			"configuration_item": schema.StringAttribute{
				Description: "ServiceNow configuration item (CMDB CI) to associate with created tickets.",
				Optional:    true,
			},
			"department_id": schema.StringAttribute{
				Description: "ServiceNow department ID to associate with created tickets.",
				Optional:    true,
			},
			"project_configs": schema.MapNestedAttribute{
				Description: "Per-project ServiceNow ticket configuration.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable_ticket_creation": schema.BoolAttribute{
							Description: "Whether to enable ServiceNow ticket creation for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"enable_ticket_update": schema.BoolAttribute{
							Description: "Whether to enable ServiceNow ticket updates for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"enable_incident_consolidation_info_update": schema.BoolAttribute{
							Description: "Whether to enable incident consolidation info updates for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"enable_incident_resolve_update": schema.BoolAttribute{
							Description: "Whether to enable incident resolve updates for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"enable_incident_field_sync": schema.BoolAttribute{
							Description: "Whether to enable syncing incident field updates for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"enable_metric_value_update": schema.BoolAttribute{
							Description: "Whether to enable syncing metric value updates for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"configuration_item": schema.StringAttribute{
							Description: "ServiceNow configuration item (CMDB CI) for this specific project, overrides the top-level configuration_item.",
							Optional:    true,
						},
					},
				},
			},
			"table_mapping": schema.MapAttribute{
				Description: "Mapping of InsightFinder project names to ServiceNow table names (e.g., {\"my-project\" = \"incident\"}).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"resolution_code_rules": schema.ListNestedAttribute{
				Description: "Ordered list of pattern-based rules used to classify ServiceNow resolution/close codes as positive (\"like\") or negative (\"disLike\") feedback.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"pattern": schema.StringAttribute{
							Description: "Regular expression matched against the ServiceNow resolution/close code (e.g., '^Solved').",
							Required:    true,
						},
						"outcome": schema.StringAttribute{
							Description: "Feedback outcome when the pattern matches. Must be 'like' or 'disLike'.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *servicenowResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *servicenowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan servicenowResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating ServiceNow config", map[string]interface{}{
		"account":      plan.Account.ValueString(),
		"service_host": plan.ServiceHost.ValueString(),
	})

	authType := "basic"
	if !plan.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
		authType = strings.ToLower(strings.TrimSpace(plan.AuthType.ValueString()))
	}
	if authType == "" {
		authType = "basic"
	}
	if authType != "basic" && authType != "oauth" {
		resp.Diagnostics.AddError(
			"Invalid Authentication Type",
			fmt.Sprintf("auth_type must be either 'basic' or 'oauth', got '%s'", authType),
		)
		return
	}

	if authType == "oauth" {
		if plan.AppID.IsNull() || plan.AppID.IsUnknown() || strings.TrimSpace(plan.AppID.ValueString()) == "" {
			resp.Diagnostics.AddError(
				"Missing app_id for OAuth",
				"auth_type is set to 'oauth' but app_id is not provided.",
			)
			return
		}
		if plan.AppKey.IsNull() || plan.AppKey.IsUnknown() || strings.TrimSpace(plan.AppKey.ValueString()) == "" {
			resp.Diagnostics.AddError(
				"Missing app_key for OAuth",
				"auth_type is set to 'oauth' but app_key is not provided.",
			)
			return
		}
	}

	// Resolve system names to system IDs for the API call.
	var systemIDs []string
	if !plan.SystemNames.IsNull() && !plan.SystemNames.IsUnknown() {
		var systemNames []string
		diags = plan.SystemNames.ElementsAs(ctx, &systemNames, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		resolvedIDs, err := r.client.ResolveSystemNameToIDs(systemNames, r.client.Username)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Resolving System Names",
				fmt.Sprintf("Could not resolve system names to IDs: %s", err.Error()),
			)
			return
		}
		systemIDs = resolvedIDs
	}

	var options []string
	diags = plan.Options.ElementsAs(ctx, &options, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var contentOption []string
	diags = plan.ContentOption.ElementsAs(ctx, &contentOption, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectConfigs, err := projectConfigsFromTF(ctx, plan.ProjectConfigs)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading project_configs", err.Error())
		return
	}

	resolutionCodeRules, err := resolutionCodeRulesFromTF(ctx, plan.ResolutionCodeRules)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading resolution_code_rules", err.Error())
		return
	}

	config := &client.ServiceNowConfig{
		Account:                    plan.Account.ValueString(),
		ServiceHost:                plan.ServiceHost.ValueString(),
		Password:                   plan.Password.ValueString(),
		Proxy:                      plan.Proxy.ValueString(),
		AppID:                      plan.AppID.ValueString(),
		AppKey:                     plan.AppKey.ValueString(),
		AuthType:                   authType,
		SystemIDs:                  systemIDs,
		Options:                    options,
		ContentOption:              contentOption,
		ServiceNowField:            plan.ServiceNowField.ValueString(),
		ContentSource:              plan.ContentSource.ValueString(),
		TriggerWindowInMills:       plan.TriggerWindowInMills.ValueInt64(),
		EnableFeedbackCollect:      plan.EnableFeedbackCollect.ValueBool(),
		TicketCreatedBySourceKey:   plan.TicketCreatedBySourceKey.ValueString(),
		TicketCreatedBySourceValue: plan.TicketCreatedBySourceValue.ValueString(),
		ConfigurationItem:          plan.ConfigurationItem.ValueString(),
		DepartmentID:               plan.DepartmentID.ValueString(),
		ProjectConfigs:             projectConfigs,
		ResolutionCodeRules:        resolutionCodeRules,
	}

	err = r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating ServiceNow Config (Verification)",
			"Could not create ServiceNow config: "+err.Error(),
		)
		return
	}

	err = r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating ServiceNow Config",
			"Could not create ServiceNow config: "+err.Error(),
		)
		return
	}

	if !plan.TableMapping.IsNull() && !plan.TableMapping.IsUnknown() {
		tableMapping, err := tableMappingFromTF(ctx, plan.TableMapping)
		if err != nil {
			resp.Diagnostics.AddError("Error Reading table_mapping", err.Error())
			return
		}
		if len(tableMapping) > 0 {
			if err := r.client.UpdateServiceNowTableMapping(
				plan.Account.ValueString(),
				plan.ServiceHost.ValueString(),
				r.client.Username,
				tableMapping,
			); err != nil {
				resp.Diagnostics.AddError(
					"Error Creating ServiceNow Table Mapping",
					"Could not set table mapping: "+err.Error(),
				)
				return
			}
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s@%s", plan.Account.ValueString(), plan.ServiceHost.ValueString()))
	plan.AuthType = types.StringValue(authType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *servicenowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state servicenowResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading ServiceNow config", map[string]interface{}{
		"account":      state.Account.ValueString(),
		"service_host": state.ServiceHost.ValueString(),
	})

	config, err := r.client.GetServiceNowConfig(
		state.Account.ValueString(),
		state.ServiceHost.ValueString(),
		r.client.Username,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ServiceNow Config",
			"Could not read ServiceNow config: "+err.Error(),
		)
		return
	}

	if config == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Proxy = types.StringValue(config.Proxy)

	// Resolve system IDs from API to names and store only the names.
	if len(config.SystemIDs) > 0 {
		names, err := r.client.ResolveSystemIDsToNames(config.SystemIDs, r.client.Username)
		if err == nil && len(names) > 0 {
			systemNamesSet, diags := types.SetValueFrom(ctx, types.StringType, names)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			state.SystemNames = systemNamesSet
		}
		// If resolution fails, keep state.SystemNames as-is.
	}
	// If API returned no system IDs, keep state.SystemNames as-is.

	optionsSet, diags := types.SetValueFrom(ctx, types.StringType, config.Options)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Options = optionsSet

	contentOptionSlice := config.ContentOption
	if contentOptionSlice == nil {
		contentOptionSlice = []string{}
	}
	contentOptionSet, diags := types.SetValueFrom(ctx, types.StringType, contentOptionSlice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.ContentOption = contentOptionSet

	authType := strings.ToLower(strings.TrimSpace(config.AuthType))
	if authType == "" {
		if !state.AuthType.IsNull() && !state.AuthType.IsUnknown() {
			authType = strings.ToLower(strings.TrimSpace(state.AuthType.ValueString()))
		}
		if authType == "" {
			authType = "basic"
		}
	}
	state.AuthType = types.StringValue(authType)

	if config.ServiceNowField != "" {
		state.ServiceNowField = types.StringValue(config.ServiceNowField)
	} else {
		state.ServiceNowField = types.StringNull()
	}

	if config.ContentSource != "" {
		state.ContentSource = types.StringValue(config.ContentSource)
	} else {
		state.ContentSource = types.StringValue("work_notes")
	}

	if config.TriggerWindowInMills > 0 {
		state.TriggerWindowInMills = types.Int64Value(config.TriggerWindowInMills)
	} else {
		state.TriggerWindowInMills = types.Int64Null()
	}

	state.EnableFeedbackCollect = types.BoolValue(config.EnableFeedbackCollect)

	if config.TicketCreatedBySourceKey != "" {
		state.TicketCreatedBySourceKey = types.StringValue(config.TicketCreatedBySourceKey)
	} else {
		state.TicketCreatedBySourceKey = types.StringNull()
	}
	if config.TicketCreatedBySourceValue != "" {
		state.TicketCreatedBySourceValue = types.StringValue(config.TicketCreatedBySourceValue)
	} else {
		state.TicketCreatedBySourceValue = types.StringNull()
	}
	if config.ConfigurationItem != "" {
		state.ConfigurationItem = types.StringValue(config.ConfigurationItem)
	} else if !state.ConfigurationItem.IsNull() && !state.ConfigurationItem.IsUnknown() && state.ConfigurationItem.ValueString() == "" {
		// Prior state was ""; API treats "" and null identically — preserve "" to avoid perpetual diff.
	} else {
		state.ConfigurationItem = types.StringNull()
	}

	if config.DepartmentID != "" {
		state.DepartmentID = types.StringValue(config.DepartmentID)
	} else {
		state.DepartmentID = types.StringNull()
	}

	if len(config.ProjectConfigs) > 0 {
		projectConfigsMap, d := projectConfigsToTF(ctx, config.ProjectConfigs, state.ProjectConfigs)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.ProjectConfigs = projectConfigsMap
	} else {
		state.ProjectConfigs = types.MapNull(projectConfigAttrTypes())
	}

	if len(config.ResolutionCodeRules) > 0 {
		resolutionCodeRulesList, d := resolutionCodeRulesToTF(ctx, config.ResolutionCodeRules)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.ResolutionCodeRules = resolutionCodeRulesList
	} else {
		state.ResolutionCodeRules = types.ListNull(resolutionCodeRuleObjectType())
	}

	if len(config.TableMapping) > 0 {
		tableMappingValues := make(map[string]attr.Value, len(config.TableMapping))
		for _, row := range config.TableMapping {
			if len(row) == 2 {
				tableMappingValues[row[0]] = types.StringValue(row[1])
			}
		}
		tableMappingMap, d := types.MapValue(types.StringType, tableMappingValues)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.TableMapping = tableMappingMap
	} else {
		state.TableMapping = types.MapNull(types.StringType)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *servicenowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan servicenowResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState servicenowResourceModel
	diags = req.State.Get(ctx, &priorState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	authType := "basic"
	if !plan.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
		authType = strings.ToLower(strings.TrimSpace(plan.AuthType.ValueString()))
	} else if !priorState.AuthType.IsNull() && !priorState.AuthType.IsUnknown() {
		authType = strings.ToLower(strings.TrimSpace(priorState.AuthType.ValueString()))
	}
	if authType == "" {
		authType = "basic"
	}
	if authType != "basic" && authType != "oauth" {
		resp.Diagnostics.AddError(
			"Invalid Authentication Type",
			fmt.Sprintf("auth_type must be either 'basic' or 'oauth', got '%s'", authType),
		)
		return
	}

	appIDValue := strings.TrimSpace(plan.AppID.ValueString())
	if plan.AppID.IsNull() || plan.AppID.IsUnknown() {
		appIDValue = strings.TrimSpace(priorState.AppID.ValueString())
	}
	appKeyValue := strings.TrimSpace(plan.AppKey.ValueString())
	if plan.AppKey.IsNull() || plan.AppKey.IsUnknown() {
		appKeyValue = strings.TrimSpace(priorState.AppKey.ValueString())
	}

	if authType == "oauth" {
		if appIDValue == "" {
			resp.Diagnostics.AddError(
				"Missing app_id for OAuth",
				"auth_type is set to 'oauth' but app_id is not provided.",
			)
			return
		}
		if appKeyValue == "" {
			resp.Diagnostics.AddError(
				"Missing app_key for OAuth",
				"auth_type is set to 'oauth' but app_key is not provided.",
			)
			return
		}
	}

	tflog.Debug(ctx, "Updating ServiceNow config", map[string]interface{}{
		"account":      plan.Account.ValueString(),
		"service_host": plan.ServiceHost.ValueString(),
	})

	// Resolve system names to system IDs for the API call.
	var systemIDs []string
	if !plan.SystemNames.IsNull() && !plan.SystemNames.IsUnknown() {
		var systemNames []string
		diags = plan.SystemNames.ElementsAs(ctx, &systemNames, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		resolvedIDs, err := r.client.ResolveSystemNameToIDs(systemNames, r.client.Username)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Resolving System Names",
				fmt.Sprintf("Could not resolve system names to IDs: %s", err.Error()),
			)
			return
		}
		systemIDs = resolvedIDs
	}

	var options []string
	diags = plan.Options.ElementsAs(ctx, &options, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var contentOption []string
	diags = plan.ContentOption.ElementsAs(ctx, &contentOption, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectConfigs, err := projectConfigsFromTF(ctx, plan.ProjectConfigs)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading project_configs", err.Error())
		return
	}

	resolutionCodeRules, err := resolutionCodeRulesFromTF(ctx, plan.ResolutionCodeRules)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading resolution_code_rules", err.Error())
		return
	}

	config := &client.ServiceNowConfig{
		Account:                    plan.Account.ValueString(),
		ServiceHost:                plan.ServiceHost.ValueString(),
		Password:                   plan.Password.ValueString(),
		Proxy:                      plan.Proxy.ValueString(),
		AppID:                      appIDValue,
		AppKey:                     appKeyValue,
		AuthType:                   authType,
		SystemIDs:                  systemIDs,
		Options:                    options,
		ContentOption:              contentOption,
		ServiceNowField:            plan.ServiceNowField.ValueString(),
		ContentSource:              plan.ContentSource.ValueString(),
		TriggerWindowInMills:       plan.TriggerWindowInMills.ValueInt64(),
		EnableFeedbackCollect:      plan.EnableFeedbackCollect.ValueBool(),
		TicketCreatedBySourceKey:   plan.TicketCreatedBySourceKey.ValueString(),
		TicketCreatedBySourceValue: plan.TicketCreatedBySourceValue.ValueString(),
		ConfigurationItem:          plan.ConfigurationItem.ValueString(),
		DepartmentID:               plan.DepartmentID.ValueString(),
		ProjectConfigs:             projectConfigs,
		ResolutionCodeRules:        resolutionCodeRules,
	}

	err = r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ServiceNow Config (Verification)",
			"Could not update ServiceNow config: "+err.Error(),
		)
		return
	}

	err = r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ServiceNow Config",
			"Could not update ServiceNow config: "+err.Error(),
		)
		return
	}

	if !plan.TableMapping.IsNull() && !plan.TableMapping.IsUnknown() {
		tableMapping, err := tableMappingFromTF(ctx, plan.TableMapping)
		if err != nil {
			resp.Diagnostics.AddError("Error Reading table_mapping", err.Error())
			return
		}
		if len(tableMapping) > 0 {
			if err := r.client.UpdateServiceNowTableMapping(
				plan.Account.ValueString(),
				plan.ServiceHost.ValueString(),
				r.client.Username,
				tableMapping,
			); err != nil {
				resp.Diagnostics.AddError(
					"Error Updating ServiceNow Table Mapping",
					"Could not set table mapping: "+err.Error(),
				)
				return
			}
		}
	}

	plan.AuthType = types.StringValue(authType)
	if appIDValue == "" {
		plan.AppID = types.StringNull()
	} else {
		plan.AppID = types.StringValue(appIDValue)
	}
	if appKeyValue == "" {
		plan.AppKey = types.StringNull()
	} else {
		plan.AppKey = types.StringValue(appKeyValue)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *servicenowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state servicenowResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting ServiceNow config", map[string]interface{}{
		"account":      state.Account.ValueString(),
		"service_host": state.ServiceHost.ValueString(),
	})

	err := r.client.DeleteServiceNowConfig(
		state.Account.ValueString(),
		state.ServiceHost.ValueString(),
		r.client.Username,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ServiceNow Config",
			"Could not delete ServiceNow config: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *servicenowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "@")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: account@service_host",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_host"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// projectConfigAttrTypes returns the attr.Type map for a projectConfigModel object.
func projectConfigAttrTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enable_ticket_creation":                    types.BoolType,
			"enable_ticket_update":                      types.BoolType,
			"enable_incident_consolidation_info_update": types.BoolType,
			"enable_incident_resolve_update":            types.BoolType,
			"enable_incident_field_sync":                types.BoolType,
			"enable_metric_value_update":                types.BoolType,
			"configuration_item":                        types.StringType,
		},
	}
}

// projectConfigsFromTF converts a Terraform types.Map into a client.ProjectConfig map.
// Returns nil (not an error) when the map is null or unknown.
func projectConfigsFromTF(ctx context.Context, m types.Map) (map[string]client.ServiceNowProjectConfig, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}
	var models map[string]projectConfigModel
	diags := m.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read project_configs: %s", diags[0].Detail())
	}
	result := make(map[string]client.ServiceNowProjectConfig, len(models))
	for projectName, pc := range models {
		result[projectName] = client.ServiceNowProjectConfig{
			EnableTicketCreation:                  pc.EnableTicketCreation.ValueBool(),
			EnableTicketUpdate:                    pc.EnableTicketUpdate.ValueBool(),
			EnableIncidentConsolidationInfoUpdate: pc.EnableIncidentConsolidationInfoUpdate.ValueBool(),
			EnableIncidentResolveUpdate:           pc.EnableIncidentResolveUpdate.ValueBool(),
			EnableIncidentFieldSync:               pc.EnableIncidentFieldSync.ValueBool(),
			EnableMetricValueUpdate:               pc.EnableMetricValueUpdate.ValueBool(),
			ConfigurationItem:                     pc.ConfigurationItem.ValueString(),
		}
	}
	return result, nil
}

// projectConfigsToTF converts a ServiceNowProjectConfig map into a Terraform types.Map.
// priorConfigs is the prior state map used to preserve "" vs null for configuration_item
// when the API returns an empty string (API treats "" and null identically).
func projectConfigsToTF(ctx context.Context, configs map[string]client.ServiceNowProjectConfig, priorConfigs types.Map) (types.Map, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"enable_ticket_creation":                    types.BoolType,
		"enable_ticket_update":                      types.BoolType,
		"enable_incident_consolidation_info_update": types.BoolType,
		"enable_incident_resolve_update":            types.BoolType,
		"enable_incident_field_sync":                types.BoolType,
		"configuration_item":                        types.StringType,
	}

	var priorModels map[string]projectConfigModel
	if !priorConfigs.IsNull() && !priorConfigs.IsUnknown() {
		_ = priorConfigs.ElementsAs(ctx, &priorModels, false)
	}

	elements := make(map[string]attr.Value, len(configs))
	for projectName, pc := range configs {
		configItem := types.StringNull()
		if pc.ConfigurationItem != "" {
			configItem = types.StringValue(pc.ConfigurationItem)
		} else if prior, ok := priorModels[projectName]; ok &&
			!prior.ConfigurationItem.IsNull() && !prior.ConfigurationItem.IsUnknown() &&
			prior.ConfigurationItem.ValueString() == "" {
			// Prior state was ""; preserve it — API treats "" same as null.
			configItem = types.StringValue("")
		}
		obj, d := types.ObjectValue(attrTypes, map[string]attr.Value{
			"enable_ticket_creation":                    types.BoolValue(pc.EnableTicketCreation),
			"enable_ticket_update":                      types.BoolValue(pc.EnableTicketUpdate),
			"enable_incident_consolidation_info_update": types.BoolValue(pc.EnableIncidentConsolidationInfoUpdate),
			"enable_incident_resolve_update":            types.BoolValue(pc.EnableIncidentResolveUpdate),
			"enable_incident_field_sync":                types.BoolValue(pc.EnableIncidentFieldSync),
			"enable_metric_value_update":                types.BoolValue(pc.EnableMetricValueUpdate),
			"configuration_item":                        configItem,
		})
		if d.HasError() {
			return types.MapNull(types.ObjectType{AttrTypes: attrTypes}), d
		}
		elements[projectName] = obj
	}
	return types.MapValue(types.ObjectType{AttrTypes: attrTypes}, elements)
}

// resolutionCodeRuleObjectType returns the attr.Type for a resolutionCodeRuleModel object.
func resolutionCodeRuleObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"pattern": types.StringType,
			"outcome": types.StringType,
		},
	}
}

// resolutionCodeRulesFromTF converts a Terraform types.List into a client.ServiceNowResolutionCodeRule slice.
// Returns nil (not an error) when the list is null or unknown.
func resolutionCodeRulesFromTF(ctx context.Context, l types.List) ([]client.ServiceNowResolutionCodeRule, error) {
	if l.IsNull() || l.IsUnknown() {
		return nil, nil
	}
	var models []resolutionCodeRuleModel
	diags := l.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read resolution_code_rules: %s", diags[0].Detail())
	}
	result := make([]client.ServiceNowResolutionCodeRule, 0, len(models))
	for _, m := range models {
		result = append(result, client.ServiceNowResolutionCodeRule{
			Pattern: m.Pattern.ValueString(),
			Outcome: m.Outcome.ValueString(),
		})
	}
	return result, nil
}

// resolutionCodeRulesToTF converts a ServiceNowResolutionCodeRule slice into a Terraform types.List.
func resolutionCodeRulesToTF(ctx context.Context, rules []client.ServiceNowResolutionCodeRule) (types.List, diag.Diagnostics) {
	objType := resolutionCodeRuleObjectType()
	elements := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		obj, d := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
			"pattern": types.StringValue(rule.Pattern),
			"outcome": types.StringValue(rule.Outcome),
		})
		if d.HasError() {
			return types.ListNull(objType), d
		}
		elements = append(elements, obj)
	}
	return types.ListValue(objType, elements)
}

// tableMappingFromTF converts a Terraform types.Map to [][]string for the API
func tableMappingFromTF(ctx context.Context, m types.Map) ([][]string, error) {
	var tableMap map[string]string
	diags := m.ElementsAs(ctx, &tableMap, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read table_mapping: %s", diags[0].Detail())
	}
	result := make([][]string, 0, len(tableMap))
	for project, table := range tableMap {
		result = append(result, []string{project, table})
	}
	return result, nil
}
