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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jwtConfigResource{}
	_ resource.ResourceWithConfigure   = &jwtConfigResource{}
	_ resource.ResourceWithImportState = &jwtConfigResource{}
)

// NewJWTConfigResource is a helper function to simplify the provider implementation.
func NewJWTConfigResource() resource.Resource {
	return &jwtConfigResource{}
}

// jwtConfigResource is the resource implementation.
type jwtConfigResource struct {
	client *client.Client
}

// jwtConfigResourceModel maps the resource schema data.
type jwtConfigResourceModel struct {
	ID         types.String `tfsdk:"id"`
	SystemName types.String `tfsdk:"system_name"`
	JWTSecret  types.String `tfsdk:"jwt_secret"`
	JWTType    types.Int64  `tfsdk:"jwt_type"`
}

// Metadata returns the resource type name.
func (r *jwtConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwt_config"
}

// Schema defines the schema for the resource.
func (r *jwtConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages InsightFinder JWT configuration for a system.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the JWT configuration (system_name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system_name": schema.StringAttribute{
				Description: "The name of the system to configure JWT for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"jwt_secret": schema.StringAttribute{
				Description: "The JWT secret token (minimum 6 characters).",
				Required:    true,
				Sensitive:   true,
			},
			"jwt_type": schema.Int64Attribute{
				Description: "The JWT type (1 for system-level JWT).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *jwtConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *jwtConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan jwtConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating JWT config", map[string]interface{}{
		"system_name": plan.SystemName.ValueString(),
	})

	// Validate JWT secret length
	jwtSecret := plan.JWTSecret.ValueString()
	if len(jwtSecret) < 6 {
		resp.Diagnostics.AddError(
			"Invalid JWT Secret",
			"JWT secret must be at least 6 characters long.",
		)
		return
	}

	// Resolve the system name to a system ID using the shared client helper
	systemID, err := r.resolveSystemID(plan.SystemName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Resolving System Name",
			fmt.Sprintf("Could not find system '%s': %s", plan.SystemName.ValueString(), err.Error()),
		)
		return
	}

	// Set default JWT type if not specified
	jwtType := int64(1)
	if !plan.JWTType.IsNull() && !plan.JWTType.IsUnknown() {
		jwtType = plan.JWTType.ValueInt64()
	}

	// Create JWT config
	jwtConfig := &client.JWTConfig{
		SystemName: plan.SystemName.ValueString(),
		SystemID:   systemID,
		JWTSecret:  jwtSecret,
		JWTType:    int(jwtType),
	}

	err = r.client.CreateOrUpdateJWTConfig(jwtConfig, r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating JWT Config",
			"Could not create JWT config: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = plan.SystemName
	plan.JWTType = types.Int64Value(jwtType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *jwtConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state jwtConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading JWT config", map[string]interface{}{
		"system_name": state.SystemName.ValueString(),
	})

	// Get current JWT configuration
	jwtConfig, err := r.client.GetJWTConfig(state.SystemName.ValueString(), r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading JWT Config",
			"Could not read JWT config: "+err.Error(),
		)
		return
	}

	// If JWT config doesn't exist, remove from state
	if jwtConfig == nil || jwtConfig.JWTSecret == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state with current values
	state.JWTSecret = types.StringValue(jwtConfig.JWTSecret)
	state.JWTType = types.Int64Value(int64(jwtConfig.JWTType))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *jwtConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan jwtConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating JWT config", map[string]interface{}{
		"system_name": plan.SystemName.ValueString(),
	})

	// Validate JWT secret length
	jwtSecret := plan.JWTSecret.ValueString()
	if len(jwtSecret) < 6 {
		resp.Diagnostics.AddError(
			"Invalid JWT Secret",
			"JWT secret must be at least 6 characters long.",
		)
		return
	}

	// Resolve the system name to a system ID using the shared client helper
	systemID, err := r.resolveSystemID(plan.SystemName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Resolving System Name",
			fmt.Sprintf("Could not find system '%s': %s", plan.SystemName.ValueString(), err.Error()),
		)
		return
	}

	// Set default JWT type if not specified
	jwtType := int64(1)
	if !plan.JWTType.IsNull() && !plan.JWTType.IsUnknown() {
		jwtType = plan.JWTType.ValueInt64()
	}

	// Update JWT config
	jwtConfig := &client.JWTConfig{
		SystemName: plan.SystemName.ValueString(),
		SystemID:   systemID,
		JWTSecret:  jwtSecret,
		JWTType:    int(jwtType),
	}

	err = r.client.CreateOrUpdateJWTConfig(jwtConfig, r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating JWT Config",
			"Could not update JWT config: "+err.Error(),
		)
		return
	}

	// Update state
	plan.JWTType = types.Int64Value(jwtType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jwtConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state jwtConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting JWT config", map[string]interface{}{
		"system_name": state.SystemName.ValueString(),
	})

	// Resolve the system name to a system ID
	systemID, err := r.resolveSystemID(state.SystemName.ValueString())
	if err != nil {
		// If system not found, it's already deleted
		tflog.Debug(ctx, "System not found, considering JWT config as already deleted")
		return
	}

	// Delete JWT config by setting empty secret
	jwtConfig := &client.JWTConfig{
		SystemName: state.SystemName.ValueString(),
		SystemID:   systemID,
		JWTSecret:  "",
		JWTType:    0,
	}

	err = r.client.DeleteJWTConfig(jwtConfig, r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting JWT Config",
			"Could not delete JWT config: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *jwtConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using system_name as the ID
	resource.ImportStatePassthroughID(ctx, path.Root("system_name"), req, resp)
}

// resolveSystemID finds the system ID for a given system name
func (r *jwtConfigResource) resolveSystemID(systemName string) (string, error) {
	trimmedName := strings.TrimSpace(systemName)
	if trimmedName == "" {
		return "", fmt.Errorf("system name cannot be empty")
	}

	resolvedIDs, err := r.client.ResolveSystemNameToIDs([]string{trimmedName}, r.client.Username)
	if err != nil {
		return "", err
	}

	if len(resolvedIDs) == 0 {
		return "", fmt.Errorf("system '%s' not found", trimmedName)
	}

	id := strings.TrimSpace(resolvedIDs[0])
	if id == "" {
		return "", fmt.Errorf("system '%s' returned empty identifier", trimmedName)
	}

	return id, nil
}
