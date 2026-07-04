// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &slackResource{}
	_ resource.ResourceWithConfigure   = &slackResource{}
	_ resource.ResourceWithImportState = &slackResource{}
)

// NewSlackResource is a helper function to simplify the provider implementation.
func NewSlackResource() resource.Resource {
	return &slackResource{}
}

// slackResource is the resource implementation.
type slackResource struct {
	client *client.Client
}

// slackProjectRuleModel maps an optional per-project Slack match rule.
type slackProjectRuleModel struct {
	Type    types.String `tfsdk:"type"`
	Keyword types.String `tfsdk:"keyword"`
}

// slackProjectConfigModel maps a per-project Slack notification override.
type slackProjectConfigModel struct {
	ProjectName                   types.String `tfsdk:"project_name"`
	Channel                       types.String `tfsdk:"channel"`
	Webhook                       types.String `tfsdk:"webhook"`
	Options                       types.Set    `tfsdk:"options"`
	EnableConsolidationInfoUpdate types.Bool   `tfsdk:"enable_consolidation_info_update"`
	PriorityLevels                types.List   `tfsdk:"priority_levels"`
	Rule                          types.Object `tfsdk:"rule"`
}

// slackResourceModel maps the resource schema data.
type slackResourceModel struct {
	ID             types.String `tfsdk:"id"`
	SystemName     types.String `tfsdk:"system_name"`
	Webhook        types.String `tfsdk:"webhook"`
	ChannelName    types.String `tfsdk:"channel_name"`
	Options        types.Set    `tfsdk:"options"`
	ProjectConfigs types.List   `tfsdk:"project_configs"`
}

// Metadata returns the resource type name.
func (r *slackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slack"
}

// Schema defines the schema for the resource.
func (r *slackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an InsightFinder Slack integration for a system.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "API-generated identifier (account ID) for this Slack integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system_name": schema.StringAttribute{
				Description: "InsightFinder system name to integrate with Slack.",
				Required:    true,
			},
			"webhook": schema.StringAttribute{
				Description: "Slack incoming webhook URL.",
				Required:    true,
				Sensitive:   true,
			},
			"channel_name": schema.StringAttribute{
				Description: "Slack channel to send notifications to (e.g., '#my-channel').",
				Required:    true,
			},
			"options": schema.SetAttribute{
				Description: "Notification types to send to Slack (e.g., 'Detected Incident', 'Predicted Incident', 'New Pattern Alert', 'Missing Monitoring Data').",
				Required:    true,
				ElementType: types.StringType,
			},
			"project_configs": schema.ListNestedAttribute{
				Description: "Per-project Slack notification overrides. A project may appear more than once to notify multiple channels.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"project_name": schema.StringAttribute{
							Description: "InsightFinder project name.",
							Required:    true,
						},
						"channel": schema.StringAttribute{
							Description: "Slack channel override for this project. Defaults to the top-level channel_name when empty.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"webhook": schema.StringAttribute{
							Description: "Slack webhook override for this project. Defaults to the top-level webhook when empty.",
							Optional:    true,
							Computed:    true,
							Sensitive:   true,
							Default:     stringdefault.StaticString(""),
						},
						"options": schema.SetAttribute{
							Description: "Notification types to send for this project. Defaults to the top-level options when empty.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"enable_consolidation_info_update": schema.BoolAttribute{
							Description: "Whether to send incident consolidation info updates for this project.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"priority_levels": schema.ListAttribute{
							Description: "Priority levels that trigger notifications for this project (e.g., [1, 2, 3]).",
							Optional:    true,
							ElementType: types.Int64Type,
						},
						"rule": schema.SingleNestedAttribute{
							Description: "Optional match rule that filters which alerts are sent for this project.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "Rule type (e.g., 'fieldName' or 'content').",
									Optional:    true,
								},
								"keyword": schema.StringAttribute{
									Description: "Rule match expression.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *slackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *slackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan slackResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Slack integration", map[string]interface{}{
		"system_name":  plan.SystemName.ValueString(),
		"channel_name": plan.ChannelName.ValueString(),
	})

	systemID, err := r.resolveSystemID(plan.SystemName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Resolving System Name", err.Error())
		return
	}

	var options []string
	diags = plan.Options.ElementsAs(ctx, &options, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectConfigs, err := slackProjectConfigsFromTF(ctx, plan.ProjectConfigs)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading project_configs", err.Error())
		return
	}

	config := &client.SlackConfig{
		SystemID:       systemID,
		WebHook:        plan.Webhook.ValueString(),
		ChannelName:    plan.ChannelName.ValueString(),
		Options:        options,
		ProjectConfigs: projectConfigs,
	}

	account, err := r.client.CreateSlackConfig(config, r.client.Username)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Slack Integration",
			"Could not create Slack integration: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(account)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *slackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state slackResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Slack integration", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	config, err := r.client.GetSlackConfig(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Slack Integration",
			"Could not read Slack integration: "+err.Error(),
		)
		return
	}

	if config == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Webhook = types.StringValue(config.WebHook)
	state.ChannelName = types.StringValue(config.ChannelName)

	optionsSet, diags := types.SetValueFrom(ctx, types.StringType, config.Options)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Options = optionsSet

	if config.SystemID != "" {
		names, err := r.client.ResolveSystemIDsToNames([]string{config.SystemID}, r.client.Username)
		if err == nil && len(names) > 0 {
			state.SystemName = types.StringValue(names[0])
		}
		// If resolution fails, keep state.SystemName as-is.
	}

	if len(config.ProjectConfigs) > 0 {
		projectConfigsList, d := slackProjectConfigsToTF(ctx, config.ProjectConfigs)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.ProjectConfigs = projectConfigsList
	} else {
		state.ProjectConfigs = types.ListNull(types.ObjectType{AttrTypes: slackProjectConfigAttrTypes()})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *slackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan slackResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState slackResourceModel
	diags = req.State.Get(ctx, &priorState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Slack integration", map[string]interface{}{
		"id": priorState.ID.ValueString(),
	})

	systemID, err := r.resolveSystemID(plan.SystemName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Resolving System Name", err.Error())
		return
	}

	oldSystemID, err := r.resolveSystemID(priorState.SystemName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Resolving System Name", err.Error())
		return
	}

	var options []string
	diags = plan.Options.ElementsAs(ctx, &options, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectConfigs, err := slackProjectConfigsFromTF(ctx, plan.ProjectConfigs)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading project_configs", err.Error())
		return
	}

	newConfig := &client.SlackConfig{
		SystemID:       systemID,
		WebHook:        plan.Webhook.ValueString(),
		ChannelName:    plan.ChannelName.ValueString(),
		Options:        options,
		ProjectConfigs: projectConfigs,
	}

	// The integration to update is located by matching the prior system, webhook, and
	// channel name against the list of existing integrations, rather than trusting the
	// account ID already stored in state.
	account, err := r.client.UpdateSlackConfig(
		oldSystemID,
		priorState.Webhook.ValueString(),
		priorState.ChannelName.ValueString(),
		newConfig,
		r.client.Username,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Slack Integration",
			"Could not update Slack integration: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(account)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *slackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state slackResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Slack integration", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	if err := r.client.DeleteSlackConfig(state.ID.ValueString(), r.client.Username); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Slack Integration",
			"Could not delete Slack integration: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state using the account ID as the import ID.
func (r *slackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// resolveSystemID resolves a single system name to its system ID.
func (r *slackResource) resolveSystemID(systemName string) (string, error) {
	systemIDs, err := r.client.ResolveSystemNameToIDs([]string{systemName}, r.client.Username)
	if err != nil {
		return "", fmt.Errorf("could not resolve system name to ID: %w", err)
	}
	if len(systemIDs) == 0 {
		return "", fmt.Errorf("could not resolve system name %q to an ID", systemName)
	}
	return systemIDs[0], nil
}

// slackProjectRuleAttrTypes returns the attr.Type map for a slackProjectRuleModel object.
func slackProjectRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":    types.StringType,
		"keyword": types.StringType,
	}
}

// slackProjectConfigAttrTypes returns the attr.Type map for a slackProjectConfigModel object.
func slackProjectConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_name":                     types.StringType,
		"channel":                          types.StringType,
		"webhook":                          types.StringType,
		"options":                          types.SetType{ElemType: types.StringType},
		"enable_consolidation_info_update": types.BoolType,
		"priority_levels":                  types.ListType{ElemType: types.Int64Type},
		"rule":                             types.ObjectType{AttrTypes: slackProjectRuleAttrTypes()},
	}
}

// slackProjectConfigsFromTF converts a Terraform types.List into a client.SlackProjectConfig slice.
// Returns nil (not an error) when the list is null or unknown.
func slackProjectConfigsFromTF(ctx context.Context, l types.List) ([]client.SlackProjectConfig, error) {
	if l.IsNull() || l.IsUnknown() {
		return nil, nil
	}

	var models []slackProjectConfigModel
	diags := l.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read project_configs: %s", diags[0].Detail())
	}

	result := make([]client.SlackProjectConfig, 0, len(models))
	for _, m := range models {
		pc := client.SlackProjectConfig{
			ProjectName:                   m.ProjectName.ValueString(),
			Channel:                       m.Channel.ValueString(),
			Webhook:                       m.Webhook.ValueString(),
			EnableConsolidationInfoUpdate: m.EnableConsolidationInfoUpdate.ValueBool(),
		}

		if !m.Options.IsNull() && !m.Options.IsUnknown() {
			var opts []string
			d := m.Options.ElementsAs(ctx, &opts, false)
			if d.HasError() {
				return nil, fmt.Errorf("failed to read project_configs options: %s", d[0].Detail())
			}
			pc.Options = opts
		}

		if !m.PriorityLevels.IsNull() && !m.PriorityLevels.IsUnknown() {
			var levels []int64
			d := m.PriorityLevels.ElementsAs(ctx, &levels, false)
			if d.HasError() {
				return nil, fmt.Errorf("failed to read project_configs priority_levels: %s", d[0].Detail())
			}
			pc.PriorityLevels = levels
		}

		if !m.Rule.IsNull() && !m.Rule.IsUnknown() {
			var ruleModel slackProjectRuleModel
			d := m.Rule.As(ctx, &ruleModel, basetypes.ObjectAsOptions{})
			if d.HasError() {
				return nil, fmt.Errorf("failed to read project_configs rule: %s", d[0].Detail())
			}
			pc.Rule = &client.SlackProjectRule{
				Type:    ruleModel.Type.ValueString(),
				Keyword: ruleModel.Keyword.ValueString(),
			}
		}

		result = append(result, pc)
	}
	return result, nil
}

// slackProjectConfigsToTF converts a client.SlackProjectConfig slice into a Terraform types.List.
func slackProjectConfigsToTF(ctx context.Context, configs []client.SlackProjectConfig) (types.List, diag.Diagnostics) {
	elemType := types.ObjectType{AttrTypes: slackProjectConfigAttrTypes()}

	models := make([]slackProjectConfigModel, 0, len(configs))
	for _, pc := range configs {
		optionsSet, d := types.SetValueFrom(ctx, types.StringType, pc.Options)
		if d.HasError() {
			return types.ListNull(elemType), d
		}

		// pc.PriorityLevels is nil whenever the API returned no priority levels for this
		// project (whether omitted or sent as an empty array) — pass it through as-is so
		// ListValueFrom produces null, matching the schema's un-defaulted Optional attribute.
		levelsList, d := types.ListValueFrom(ctx, types.Int64Type, pc.PriorityLevels)
		if d.HasError() {
			return types.ListNull(elemType), d
		}

		ruleObj := types.ObjectNull(slackProjectRuleAttrTypes())
		if pc.Rule != nil {
			obj, d := types.ObjectValueFrom(ctx, slackProjectRuleAttrTypes(), slackProjectRuleModel{
				Type:    types.StringValue(pc.Rule.Type),
				Keyword: types.StringValue(pc.Rule.Keyword),
			})
			if d.HasError() {
				return types.ListNull(elemType), d
			}
			ruleObj = obj
		}

		models = append(models, slackProjectConfigModel{
			ProjectName:                   types.StringValue(pc.ProjectName),
			Channel:                       types.StringValue(pc.Channel),
			Webhook:                       types.StringValue(pc.Webhook),
			Options:                       optionsSet,
			EnableConsolidationInfoUpdate: types.BoolValue(pc.EnableConsolidationInfoUpdate),
			PriorityLevels:                levelsList,
			Rule:                          ruleObj,
		})
	}

	return types.ListValueFrom(ctx, elemType, models)
}
