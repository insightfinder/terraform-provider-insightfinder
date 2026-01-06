// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	ID              types.String `tfsdk:"id"`
	Account         types.String `tfsdk:"account"`
	ServiceHost     types.String `tfsdk:"service_host"`
	Password        types.String `tfsdk:"password"`
	Proxy           types.String `tfsdk:"proxy"`
	DampeningPeriod types.Int64  `tfsdk:"dampening_period"`
	AppID           types.String `tfsdk:"app_id"`
	AppKey          types.String `tfsdk:"app_key"`
	AuthType        types.String `tfsdk:"auth_type"`
	SystemNames     types.List   `tfsdk:"system_names"`
	SystemIDs       types.List   `tfsdk:"system_ids"`
	Options         types.List   `tfsdk:"options"`
	ContentOption   types.List   `tfsdk:"content_option"`
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
			"system_names": schema.ListAttribute{
				Description: "List of system names to integrate (will be resolved to system IDs).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"system_ids": schema.ListAttribute{
				Description: "List of system IDs to integrate (computed from system_names if not provided).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"options": schema.ListAttribute{
				Description: "ServiceNow integration options.",
				Required:    true,
				ElementType: types.StringType,
			},
			"content_option": schema.ListAttribute{
				Description: "ServiceNow content options.",
				Required:    true,
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

	// Resolve system names to system IDs if system_names is provided
	var systemIDs []string
	resolvedNames := make([]string, 0)
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
		resolvedNames = systemNames
	} else if !plan.SystemIDs.IsNull() && !plan.SystemIDs.IsUnknown() {
		// Use provided system IDs directly
		diags = plan.SystemIDs.ElementsAs(ctx, &systemIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Resolve IDs to names
		names, err := r.client.ResolveSystemIDsToNames(systemIDs, r.client.Username)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Resolving System IDs",
				fmt.Sprintf("Could not resolve system IDs to names: %s", err.Error()),
			)
			return
		}
		resolvedNames = names
	}

	systemIDs, resolvedNames = alignSystemMappings(systemIDs, resolvedNames)

	// Get options and content options
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

	// Create ServiceNow config
	config := &client.ServiceNowConfig{
		Account:         plan.Account.ValueString(),
		ServiceHost:     plan.ServiceHost.ValueString(),
		Password:        plan.Password.ValueString(),
		Proxy:           plan.Proxy.ValueString(),
		DampeningPeriod: int(plan.DampeningPeriod.ValueInt64()),
		AppID:           plan.AppID.ValueString(),
		AppKey:          plan.AppKey.ValueString(),
		AuthType:        authType,
		SystemIDs:       systemIDs,
		SystemNames:     resolvedNames,
		Options:         options,
		ContentOption:   contentOption,
	}

	// First call with verify=true
	err := r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating ServiceNow Config (Verification)",
			"Could not create ServiceNow config: "+err.Error(),
		)
		return
	}

	// Second call without verify flag
	err = r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating ServiceNow Config",
			"Could not create ServiceNow config: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(fmt.Sprintf("%s@%s", plan.Account.ValueString(), plan.ServiceHost.ValueString()))

	// Convert system IDs and names to Terraform values
	systemIDsList, diags := types.ListValueFrom(ctx, types.StringType, systemIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SystemIDs = systemIDsList

	systemNamesList, diags := types.ListValueFrom(ctx, types.StringType, resolvedNames)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SystemNames = systemNamesList
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

	// Get current ServiceNow configuration
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

	// If config doesn't exist, remove from state
	if config == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state with current values
	// Note: Don't update sensitive fields (password, app_id, app_key) if they come back empty
	// as the API doesn't return the actual values for security reasons.
	// We keep the existing values in state.
	if strings.TrimSpace(config.Proxy) == "" {
		state.Proxy = types.StringNull()
	} else {
		state.Proxy = types.StringValue(config.Proxy)
	}
	state.DampeningPeriod = types.Int64Value(int64(config.DampeningPeriod))

	// Update system IDs and names - preserve state order where possible
	var stateSystemNames []string
	if !state.SystemNames.IsNull() && !state.SystemNames.IsUnknown() {
		diags := state.SystemNames.ElementsAs(ctx, &stateSystemNames, false)
		resp.Diagnostics.Append(diags...)
		// Continue even if there's an error; just use API order as fallback
	}

	var stateSystemIDs []string
	if !state.SystemIDs.IsNull() && !state.SystemIDs.IsUnknown() {
		diags := state.SystemIDs.ElementsAs(ctx, &stateSystemIDs, false)
		resp.Diagnostics.Append(diags...)
	}

	var resolvedNames []string
	if len(config.SystemIDs) > 0 {
		// Try to preserve state order
		if len(stateSystemIDs) == len(config.SystemIDs) {
			// Build map from API
			apiIDSet := make(map[string]struct{}, len(config.SystemIDs))
			for _, id := range config.SystemIDs {
				apiIDSet[strings.TrimSpace(id)] = struct{}{}
			}

			// Check if state IDs are still valid
			allValid := true
			for _, stateID := range stateSystemIDs {
				if _, exists := apiIDSet[strings.TrimSpace(stateID)]; !exists {
					allValid = false
					break
				}
			}

			if allValid {
				// Preserve state order
				config.SystemIDs = stateSystemIDs
				resolvedNames = stateSystemNames
			}
		}

		// If we didn't preserve state order, resolve fresh
		if len(resolvedNames) == 0 {
			if names, err := r.client.ResolveSystemIDsToNames(config.SystemIDs, r.client.Username); err == nil {
				config.SystemIDs, resolvedNames = alignSystemMappings(config.SystemIDs, names)
			} else if len(config.SystemNames) > 0 {
				config.SystemIDs, resolvedNames = alignSystemMappings(config.SystemIDs, config.SystemNames)
			} else {
				config.SystemIDs, resolvedNames = alignSystemMappings(config.SystemIDs, nil)
			}
		}
	}

	systemIDsList, diags := types.ListValueFrom(ctx, types.StringType, config.SystemIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.SystemIDs = systemIDsList

	if resolvedNames != nil {
		systemNamesList, diags := types.ListValueFrom(ctx, types.StringType, resolvedNames)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.SystemNames = systemNamesList
	} else {
		state.SystemNames = types.ListNull(types.StringType)
	}

	// Update options
	optionsList, diags := types.ListValueFrom(ctx, types.StringType, config.Options)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Options = optionsList

	// Update content options
	contentOptionList, diags := types.ListValueFrom(ctx, types.StringType, config.ContentOption)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.ContentOption = contentOptionList

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

	effectiveSystemNames := plan.SystemNames
	if plan.SystemNames.IsUnknown() && !priorState.SystemNames.IsNull() && !priorState.SystemNames.IsUnknown() {
		effectiveSystemNames = priorState.SystemNames
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

	// Resolve system names to system IDs if system_names is provided
	var systemIDs []string
	resolvedNames := make([]string, 0)
	if !effectiveSystemNames.IsNull() && !effectiveSystemNames.IsUnknown() {
		var systemNames []string
		diags = effectiveSystemNames.ElementsAs(ctx, &systemNames, false)
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
		resolvedNames = systemNames
	} else if !plan.SystemIDs.IsNull() && !plan.SystemIDs.IsUnknown() {
		// Use provided system IDs directly
		diags = plan.SystemIDs.ElementsAs(ctx, &systemIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Resolve IDs to names
		names, err := r.client.ResolveSystemIDsToNames(systemIDs, r.client.Username)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Resolving System IDs",
				fmt.Sprintf("Could not resolve system IDs to names: %s", err.Error()),
			)
			return
		}
		resolvedNames = names
	}

	systemIDs, resolvedNames = alignSystemMappings(systemIDs, resolvedNames)

	// Get options and content options
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

	// Update ServiceNow config
	config := &client.ServiceNowConfig{
		Account:         plan.Account.ValueString(),
		ServiceHost:     plan.ServiceHost.ValueString(),
		Password:        plan.Password.ValueString(),
		Proxy:           plan.Proxy.ValueString(),
		DampeningPeriod: int(plan.DampeningPeriod.ValueInt64()),
		AppID:           appIDValue,
		AppKey:          appKeyValue,
		AuthType:        authType,
		SystemIDs:       systemIDs,
		SystemNames:     resolvedNames,
		Options:         options,
		ContentOption:   contentOption,
	}

	// First call with verify=true
	err := r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ServiceNow Config (Verification)",
			"Could not update ServiceNow config: "+err.Error(),
		)
		return
	}

	// Second call without verify flag
	err = r.client.CreateOrUpdateServiceNowConfig(config, r.client.Username, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ServiceNow Config",
			"Could not update ServiceNow config: "+err.Error(),
		)
		return
	}

	// Update system IDs and names in state
	systemIDsList, diags := types.ListValueFrom(ctx, types.StringType, systemIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SystemIDs = systemIDsList

	systemNamesList, diags := types.ListValueFrom(ctx, types.StringType, resolvedNames)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SystemNames = systemNamesList
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
	// Import using format: account@service_host
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

func alignSystemMappings(systemIDs []string, systemNames []string) ([]string, []string) {
	if len(systemIDs) == 0 {
		return []string{}, systemNames
	}

	// Build map of ID -> Name for quick lookup
	idToName := make(map[string]string, len(systemIDs))
	for i, id := range systemIDs {
		trimmedID := strings.TrimSpace(id)
		if trimmedID == "" {
			continue
		}
		var name string
		if i < len(systemNames) {
			name = strings.TrimSpace(systemNames[i])
		}
		if name == "" {
			name = trimmedID
		}
		idToName[trimmedID] = name
	}

	// Preserve order from systemIDs, remove duplicates
	seen := make(map[string]struct{}, len(systemIDs))
	alignedIDs := make([]string, 0, len(systemIDs))
	alignedNames := make([]string, 0, len(systemIDs))

	for _, rawID := range systemIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		alignedIDs = append(alignedIDs, id)
		alignedNames = append(alignedNames, idToName[id])
	}

	return alignedIDs, alignedNames
}
