// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// servicenowResourceModel maps the resource schema data.
type servicenowResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Account               types.String `tfsdk:"account"`
	ServiceHost           types.String `tfsdk:"service_host"`
	Password              types.String `tfsdk:"password"`
	Proxy                 types.String `tfsdk:"proxy"`
	DampeningPeriod       types.Int64  `tfsdk:"dampening_period"`
	AppID                 types.String `tfsdk:"app_id"`
	AppKey                types.String `tfsdk:"app_key"`
	AuthType              types.String `tfsdk:"auth_type"`
	SystemNames           types.Set    `tfsdk:"system_names"`
	Options               types.Set    `tfsdk:"options"`
	ContentOption         types.Set    `tfsdk:"content_option"`
	ServiceNowField       types.String `tfsdk:"service_now_field"`
	ContentSource         types.String `tfsdk:"content_source"`
	TriggerWindowInMills  types.Int64  `tfsdk:"trigger_window_in_mills"`
	EnableFeedbackCollect types.Bool   `tfsdk:"enable_feedback_collect"`
	EnableTicketCreation  types.Bool   `tfsdk:"enable_ticket_creation"`
	TableMapping          types.Map    `tfsdk:"table_mapping"`
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
			"dampening_period": schema.Int64Attribute{
				Description: "Dampening period in seconds.",
				Required:    true,
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
			"enable_ticket_creation": schema.BoolAttribute{
				Description: "Whether to enable ServiceNow ticket creation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"table_mapping": schema.MapAttribute{
				Description: "Mapping of InsightFinder project names to ServiceNow table names (e.g., {\"my-project\" = \"incident\"}).",
				Optional:    true,
				ElementType: types.StringType,
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

	config := &client.ServiceNowConfig{
		Account:               plan.Account.ValueString(),
		ServiceHost:           plan.ServiceHost.ValueString(),
		Password:              plan.Password.ValueString(),
		Proxy:                 plan.Proxy.ValueString(),
		DampeningPeriod:       int(plan.DampeningPeriod.ValueInt64()),
		AppID:                 plan.AppID.ValueString(),
		AppKey:                plan.AppKey.ValueString(),
		AuthType:              authType,
		SystemIDs:             systemIDs,
		Options:               options,
		ContentOption:         contentOption,
		ServiceNowField:       plan.ServiceNowField.ValueString(),
		ContentSource:         plan.ContentSource.ValueString(),
		TriggerWindowInMills:  plan.TriggerWindowInMills.ValueInt64(),
		EnableFeedbackCollect: plan.EnableFeedbackCollect.ValueBool(),
		EnableTicketCreation:  plan.EnableTicketCreation.ValueBool(),
	}

	err := r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, true)
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
	state.DampeningPeriod = types.Int64Value(int64(config.DampeningPeriod))

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
	state.EnableTicketCreation = types.BoolValue(config.EnableTicketCreation)

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

	config := &client.ServiceNowConfig{
		Account:               plan.Account.ValueString(),
		ServiceHost:           plan.ServiceHost.ValueString(),
		Password:              plan.Password.ValueString(),
		Proxy:                 plan.Proxy.ValueString(),
		DampeningPeriod:       int(plan.DampeningPeriod.ValueInt64()),
		AppID:                 appIDValue,
		AppKey:                appKeyValue,
		AuthType:              authType,
		SystemIDs:             systemIDs,
		Options:               options,
		ContentOption:         contentOption,
		ServiceNowField:       plan.ServiceNowField.ValueString(),
		ContentSource:         plan.ContentSource.ValueString(),
		TriggerWindowInMills:  plan.TriggerWindowInMills.ValueInt64(),
		EnableFeedbackCollect: plan.EnableFeedbackCollect.ValueBool(),
		EnableTicketCreation:  plan.EnableTicketCreation.ValueBool(),
	}

	err := r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, true)
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
