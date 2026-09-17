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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jiraResource{}
	_ resource.ResourceWithConfigure   = &jiraResource{}
	_ resource.ResourceWithImportState = &jiraResource{}
)

// NewJiraResource is a helper function to simplify the provider implementation.
func NewJiraResource() resource.Resource {
	return &jiraResource{}
}

// jiraResource is the resource implementation.
type jiraResource struct {
	client *client.Client
}

// jiraFieldUpdateRuleModel maps a single post-creation Jira field-update rule.
type jiraFieldUpdateRuleModel struct {
	Projects   types.List   `tfsdk:"projects"`
	Metrics    types.List   `tfsdk:"metrics"`
	Components types.List   `tfsdk:"components"`
	FieldID    types.String `tfsdk:"field_id"`
	FieldName  types.String `tfsdk:"field_name"`
	ValueType  types.String `tfsdk:"value_type"`
	Value      types.String `tfsdk:"value"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

// jiraResourceModel maps the resource schema data.
type jiraResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Account              types.String `tfsdk:"account"`
	ServiceHost          types.String `tfsdk:"service_host"`
	APIToken             types.String `tfsdk:"api_token"`
	SystemNames          types.Set    `tfsdk:"system_names"`
	Options              types.Set    `tfsdk:"options"`
	ContentOption        types.Set    `tfsdk:"content_option"`
	JiraProjectKey       types.String `tfsdk:"jira_project_key"`
	JiraReporterID       types.String `tfsdk:"jira_reporter_id"`
	JiraAssigneeID       types.String `tfsdk:"jira_assignee_id"`
	JiraIssueFields      types.Map    `tfsdk:"jira_issue_fields"`
	JiraCloseStatusID    types.String `tfsdk:"jira_close_status_id"`
	DescriptionTemplate  types.String `tfsdk:"description_template"`
	TemplateFieldMapping types.Map    `tfsdk:"template_field_mapping"`
	FieldUpdateRules     types.List   `tfsdk:"field_update_rules"`
}

// Metadata returns the resource type name.
func (r *jiraResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira"
}

// Schema defines the schema for the resource.
func (r *jiraResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages InsightFinder Jira integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the Jira configuration (account@service_host).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account": schema.StringAttribute{
				Description: "Jira account email/username.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_host": schema.StringAttribute{
				Description: "Jira instance base URL (e.g. https://yourcompany.atlassian.net/). A trailing slash is added automatically if omitted.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_token": schema.StringAttribute{
				Description: "Jira API token generated for the account (not the account's login password).",
				Required:    true,
				Sensitive:   true,
			},
			"system_names": schema.SetAttribute{
				Description: "InsightFinder system names to associate with this Jira integration. At least one is required.",
				Required:    true,
				ElementType: types.StringType,
			},
			"options": schema.SetAttribute{
				Description: "Notification types that trigger Jira ticket creation (e.g., 'Detected Incident', 'Predicted Incident', 'Root Cause').",
				Optional:    true,
				ElementType: types.StringType,
			},
			"content_option": schema.SetAttribute{
				Description: "Content sections to include in the Jira ticket. Must be 'SUMMARY' and/or 'RECOMMENDATION'.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"jira_project_key": schema.StringAttribute{
				Description: "Jira project key that created tickets belong to (e.g., 'II').",
				Required:    true,
			},
			"jira_reporter_id": schema.StringAttribute{
				Description: "Jira account ID to set as the reporter on created tickets.",
				Required:    true,
			},
			"jira_assignee_id": schema.StringAttribute{
				Description: "Jira account ID to assign created tickets to. Leave unset to not set an assignee.",
				Optional:    true,
			},
			"jira_issue_fields": schema.MapAttribute{
				Description: "Additional Jira issue fields to set at ticket-creation time (field ID/name to value).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"jira_close_status_id": schema.StringAttribute{
				Description: "Jira status ID to transition a ticket to when the InsightFinder incident is resolved.",
				Optional:    true,
			},
			"description_template": schema.StringAttribute{
				Description: "Template used to render the Jira ticket description.",
				Optional:    true,
			},
			"template_field_mapping": schema.MapAttribute{
				Description: "Mapping of template placeholder names to Jira field IDs, used with description_template.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"field_update_rules": schema.ListNestedAttribute{
				Description: "Ordered post-creation field-update rules: when a ticket's project/metric/component match, field_id is set to value after creation (reaches fields that are not on the issue type's create screen).",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"projects": schema.ListAttribute{
							Description: "Root-cause project names this rule matches. Empty/omitted matches all.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"metrics": schema.ListAttribute{
							Description: "Root-cause metric names this rule matches. Empty/omitted matches all.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"components": schema.ListAttribute{
							Description: "Root-cause component names this rule matches. Empty/omitted matches all.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"field_id": schema.StringAttribute{
							Description: "Jira field ID to update (e.g., 'customfield_10050').",
							Required:    true,
						},
						"field_name": schema.StringAttribute{
							Description: "Human-readable name of the field, for documentation purposes.",
							Optional:    true,
						},
						"value_type": schema.StringAttribute{
							Description: "Type of the value being set (e.g., 'string', 'option'), used to shape the update payload.",
							Optional:    true,
						},
						"value": schema.StringAttribute{
							Description: "Value to set on field_id.",
							Required:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether this rule is active.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(true),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *jiraResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *jiraResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan jiraResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Jira config", map[string]interface{}{
		"account":      plan.Account.ValueString(),
		"service_host": plan.ServiceHost.ValueString(),
	})

	config, diags := r.buildJiraConfig(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateOrUpdateJiraConfig(config, r.client.Username, true); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Jira Config (Verification)",
			"Could not create Jira config: "+err.Error(),
		)
		return
	}

	if err := r.client.CreateOrUpdateJiraConfig(config, r.client.Username, false); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Jira Config",
			"Could not create Jira config: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s@%s", plan.Account.ValueString(), plan.ServiceHost.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *jiraResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state jiraResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Jira config", map[string]interface{}{
		"account":      state.Account.ValueString(),
		"service_host": state.ServiceHost.ValueString(),
	})

	config, err := r.client.GetJiraConfig(
		state.Account.ValueString(),
		state.ServiceHost.ValueString(),
		r.client.Username,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jira Config",
			"Could not read Jira config: "+err.Error(),
		)
		return
	}

	if config == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// api_token is intentionally left untouched: the backend never returns the raw token back
	// (only a derived encodedToken), so — same as password on the ServiceNow resource — state
	// keeps whatever value it already has instead of trying to round-trip it.

	if len(config.SystemIDs) > 0 {
		names, err := r.client.ResolveSystemIDsToNames(config.SystemIDs, r.client.Username)
		if err == nil && len(names) > 0 {
			systemNamesSet, d := types.SetValueFrom(ctx, types.StringType, names)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			state.SystemNames = systemNamesSet
		}
		// If resolution fails, keep state.SystemNames as-is.
	}

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

	if config.JiraProjectKey != "" {
		state.JiraProjectKey = types.StringValue(config.JiraProjectKey)
	}
	if config.JiraReporterID != "" {
		state.JiraReporterID = types.StringValue(config.JiraReporterID)
	}

	if config.JiraAssigneeID != "" {
		state.JiraAssigneeID = types.StringValue(config.JiraAssigneeID)
	} else {
		state.JiraAssigneeID = types.StringNull()
	}
	if config.JiraCloseStatusID != "" {
		state.JiraCloseStatusID = types.StringValue(config.JiraCloseStatusID)
	} else {
		state.JiraCloseStatusID = types.StringNull()
	}
	if config.DescriptionTemplate != "" {
		state.DescriptionTemplate = types.StringValue(config.DescriptionTemplate)
	} else {
		state.DescriptionTemplate = types.StringNull()
	}

	if len(config.JiraIssueFields) > 0 {
		m, d := types.MapValueFrom(ctx, types.StringType, config.JiraIssueFields)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.JiraIssueFields = m
	} else {
		state.JiraIssueFields = types.MapNull(types.StringType)
	}

	if len(config.TemplateFieldMapping) > 0 {
		m, d := types.MapValueFrom(ctx, types.StringType, config.TemplateFieldMapping)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.TemplateFieldMapping = m
	} else {
		state.TemplateFieldMapping = types.MapNull(types.StringType)
	}

	if len(config.FieldUpdateRules) > 0 {
		l, d := jiraFieldUpdateRulesToTF(ctx, config.FieldUpdateRules)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.FieldUpdateRules = l
	} else {
		state.FieldUpdateRules = types.ListNull(jiraFieldUpdateRuleObjectType())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *jiraResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan jiraResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Jira config", map[string]interface{}{
		"account":      plan.Account.ValueString(),
		"service_host": plan.ServiceHost.ValueString(),
	})

	config, diags := r.buildJiraConfig(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateOrUpdateJiraConfig(config, r.client.Username, true); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Jira Config (Verification)",
			"Could not update Jira config: "+err.Error(),
		)
		return
	}

	if err := r.client.CreateOrUpdateJiraConfig(config, r.client.Username, false); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Jira Config",
			"Could not update Jira config: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jiraResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state jiraResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Jira config", map[string]interface{}{
		"account":      state.Account.ValueString(),
		"service_host": state.ServiceHost.ValueString(),
	})

	err := r.client.DeleteJiraConfig(
		state.Account.ValueString(),
		state.ServiceHost.ValueString(),
		r.client.Username,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Jira Config",
			"Could not delete Jira config: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *jiraResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// A Jira account is an email address, so the ID carries an "@" of its own
	// ("user@company.com@https://company.atlassian.net/"). Split on the LAST "@" — the one
	// separating account from service_host — because splitting on every "@" yields three
	// parts and would reject every well-formed Jira import ID.
	sep := strings.LastIndex(req.ID, "@")
	if sep <= 0 || sep == len(req.ID)-1 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: account@service_host, "+
				"for example user@company.com@https://company.atlassian.net/",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account"), req.ID[:sep])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_host"), req.ID[sep+1:])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// buildJiraConfig converts the plan (shared shape between Create and Update) into a client.JiraConfig.
func (r *jiraResource) buildJiraConfig(ctx context.Context, plan jiraResourceModel) (*client.JiraConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	var systemIDs []string
	if !plan.SystemNames.IsNull() && !plan.SystemNames.IsUnknown() {
		var systemNames []string
		d := plan.SystemNames.ElementsAs(ctx, &systemNames, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		resolvedIDs, err := r.client.ResolveSystemNameToIDs(systemNames, r.client.Username)
		if err != nil {
			diags.AddError("Error Resolving System Names", fmt.Sprintf("Could not resolve system names to IDs: %s", err.Error()))
			return nil, diags
		}
		systemIDs = resolvedIDs
	}

	var options []string
	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		d := plan.Options.ElementsAs(ctx, &options, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	var contentOption []string
	if !plan.ContentOption.IsNull() && !plan.ContentOption.IsUnknown() {
		d := plan.ContentOption.ElementsAs(ctx, &contentOption, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	jiraIssueFields, err := stringMapFromTF(ctx, plan.JiraIssueFields)
	if err != nil {
		diags.AddError("Error Reading jira_issue_fields", err.Error())
		return nil, diags
	}

	templateFieldMapping, err := stringMapFromTF(ctx, plan.TemplateFieldMapping)
	if err != nil {
		diags.AddError("Error Reading template_field_mapping", err.Error())
		return nil, diags
	}

	fieldUpdateRules, err := jiraFieldUpdateRulesFromTF(ctx, plan.FieldUpdateRules)
	if err != nil {
		diags.AddError("Error Reading field_update_rules", err.Error())
		return nil, diags
	}

	return &client.JiraConfig{
		Account:              plan.Account.ValueString(),
		HostURL:              plan.ServiceHost.ValueString(),
		APIToken:             plan.APIToken.ValueString(),
		SystemIDs:            systemIDs,
		Options:              options,
		ContentOption:        contentOption,
		JiraProjectKey:       plan.JiraProjectKey.ValueString(),
		JiraReporterID:       plan.JiraReporterID.ValueString(),
		JiraAssigneeID:       plan.JiraAssigneeID.ValueString(),
		JiraIssueFields:      jiraIssueFields,
		JiraCloseStatusID:    plan.JiraCloseStatusID.ValueString(),
		DescriptionTemplate:  plan.DescriptionTemplate.ValueString(),
		TemplateFieldMapping: templateFieldMapping,
		FieldUpdateRules:     fieldUpdateRules,
	}, diags
}

// jiraFieldUpdateRuleObjectType returns the attr.Type for a jiraFieldUpdateRuleModel object.
func jiraFieldUpdateRuleObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"projects":   types.ListType{ElemType: types.StringType},
			"metrics":    types.ListType{ElemType: types.StringType},
			"components": types.ListType{ElemType: types.StringType},
			"field_id":   types.StringType,
			"field_name": types.StringType,
			"value_type": types.StringType,
			"value":      types.StringType,
			"enabled":    types.BoolType,
		},
	}
}

// jiraFieldUpdateRulesFromTF converts a Terraform types.List into a client.JiraFieldUpdateRule slice.
// Returns nil (not an error) when the list is null or unknown.
func jiraFieldUpdateRulesFromTF(ctx context.Context, l types.List) ([]client.JiraFieldUpdateRule, error) {
	if l.IsNull() || l.IsUnknown() {
		return nil, nil
	}
	var models []jiraFieldUpdateRuleModel
	diags := l.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read field_update_rules: %s", diags[0].Detail())
	}
	result := make([]client.JiraFieldUpdateRule, 0, len(models))
	for _, m := range models {
		rule := client.JiraFieldUpdateRule{
			FieldID:   m.FieldID.ValueString(),
			FieldName: m.FieldName.ValueString(),
			ValueType: m.ValueType.ValueString(),
			Value:     m.Value.ValueString(),
			Enabled:   m.Enabled.ValueBool(),
		}
		if !m.Projects.IsNull() && !m.Projects.IsUnknown() {
			d := m.Projects.ElementsAs(ctx, &rule.Projects, false)
			if d.HasError() {
				return nil, fmt.Errorf("failed to read field_update_rules projects: %s", d[0].Detail())
			}
		}
		if !m.Metrics.IsNull() && !m.Metrics.IsUnknown() {
			d := m.Metrics.ElementsAs(ctx, &rule.Metrics, false)
			if d.HasError() {
				return nil, fmt.Errorf("failed to read field_update_rules metrics: %s", d[0].Detail())
			}
		}
		if !m.Components.IsNull() && !m.Components.IsUnknown() {
			d := m.Components.ElementsAs(ctx, &rule.Components, false)
			if d.HasError() {
				return nil, fmt.Errorf("failed to read field_update_rules components: %s", d[0].Detail())
			}
		}
		result = append(result, rule)
	}
	return result, nil
}

// jiraFieldUpdateRulesToTF converts a client.JiraFieldUpdateRule slice into a Terraform types.List.
func jiraFieldUpdateRulesToTF(ctx context.Context, rules []client.JiraFieldUpdateRule) (types.List, diag.Diagnostics) {
	objType := jiraFieldUpdateRuleObjectType()
	var diags diag.Diagnostics
	elements := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		projects, d := types.ListValueFrom(ctx, types.StringType, rule.Projects)
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objType), diags
		}
		metrics, d := types.ListValueFrom(ctx, types.StringType, rule.Metrics)
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objType), diags
		}
		components, d := types.ListValueFrom(ctx, types.StringType, rule.Components)
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objType), diags
		}
		obj, d := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
			"projects":   projects,
			"metrics":    metrics,
			"components": components,
			"field_id":   types.StringValue(rule.FieldID),
			"field_name": types.StringValue(rule.FieldName),
			"value_type": types.StringValue(rule.ValueType),
			"value":      types.StringValue(rule.Value),
			"enabled":    types.BoolValue(rule.Enabled),
		})
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objType), diags
		}
		elements = append(elements, obj)
	}
	l, d := types.ListValue(objType, elements)
	diags.Append(d...)
	return l, diags
}

// stringMapFromTF converts a Terraform types.Map into a plain map[string]string.
// Returns nil (not an error) when the map is null or unknown.
func stringMapFromTF(ctx context.Context, m types.Map) (map[string]string, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}
	var result map[string]string
	diags := m.ElementsAs(ctx, &result, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read map: %s", diags[0].Detail())
	}
	return result, nil
}
