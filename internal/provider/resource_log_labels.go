// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

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
	_ resource.Resource                = &logLabelsResource{}
	_ resource.ResourceWithConfigure   = &logLabelsResource{}
	_ resource.ResourceWithImportState = &logLabelsResource{}
)

// NewLogLabelsResource is a helper function to simplify the provider implementation.
func NewLogLabelsResource() resource.Resource {
	return &logLabelsResource{}
}

// logLabelsResource is the resource implementation.
type logLabelsResource struct {
	client *client.Client
}

// logLabelsResourceModel maps the resource schema data.
type logLabelsResourceModel struct {
	ID            types.String           `tfsdk:"id"`
	ProjectName   types.String           `tfsdk:"project_name"`
	LabelSettings []logLabelSettingModel `tfsdk:"label_settings"`
}

// logLabelSettingModel represents a single log label setting
type logLabelSettingModel struct {
	LabelType      types.String `tfsdk:"label_type"`
	LogLabelString types.String `tfsdk:"log_label_string"`
}

// Metadata returns the resource type name.
func (r *logLabelsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_labels"
}

// Schema defines the schema for the resource.
func (r *logLabelsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages InsightFinder log label settings for a project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the log labels configuration (project_name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_name": schema.StringAttribute{
				Description: "The name of the project to configure log labels for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label_settings": schema.ListNestedAttribute{
				Description: "List of log label settings.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label_type": schema.StringAttribute{
							Description: "Type of log label (e.g., logSeverity, logEventID, logSession, logComponent, logTransactionID, logCustomParameter).",
							Required:    true,
						},
						"log_label_string": schema.StringAttribute{
							Description: "JSON array string of log labels (e.g., '[\"ERROR\",\"WARN\"]').",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *logLabelsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *logLabelsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan logLabelsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating log labels", map[string]interface{}{
		"project_name": plan.ProjectName.ValueString(),
	})

	// Validate and convert label settings
	settings, err := r.validateAndConvertSettings(plan.LabelSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Log Label Settings",
			err.Error(),
		)
		return
	}

	// Create log labels
	err = r.client.CreateOrUpdateLogLabels(
		plan.ProjectName.ValueString(),
		r.client.Username,
		settings,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Log Labels",
			"Could not create log labels: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = plan.ProjectName

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *logLabelsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state logLabelsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading log labels", map[string]interface{}{
		"project_name": state.ProjectName.ValueString(),
	})

	// Get current log labels
	currentLabels, err := r.client.GetLogLabels(
		state.ProjectName.ValueString(),
		r.client.Username,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Log Labels",
			"Could not read log labels: "+err.Error(),
		)
		return
	}

	// If no labels exist, remove from state
	if len(currentLabels) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state with current values
	// Map API fields back to label types
	updatedSettings := make([]logLabelSettingModel, 0)

	// For each setting in the plan, check if it exists in current state
	for _, setting := range state.LabelSettings {
		labelType := setting.LabelType.ValueString()
		apiField := client.MapLabelTypeToAPIField(labelType)

		if labels, ok := currentLabels[apiField]; ok && len(labels) > 0 {
			// Convert labels array to JSON string
			labelsJSON, err := json.Marshal(labels)
			if err != nil {
				continue
			}

			updatedSettings = append(updatedSettings, logLabelSettingModel{
				LabelType:      types.StringValue(labelType),
				LogLabelString: types.StringValue(string(labelsJSON)),
			})
		}
	}

	// If no settings remain, remove from state
	if len(updatedSettings) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	state.LabelSettings = updatedSettings

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *logLabelsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan logLabelsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating log labels", map[string]interface{}{
		"project_name": plan.ProjectName.ValueString(),
	})

	// Validate and convert label settings
	settings, err := r.validateAndConvertSettings(plan.LabelSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Log Label Settings",
			err.Error(),
		)
		return
	}

	// Update log labels
	err = r.client.CreateOrUpdateLogLabels(
		plan.ProjectName.ValueString(),
		r.client.Username,
		settings,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Log Labels",
			"Could not update log labels: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *logLabelsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state logLabelsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting log labels", map[string]interface{}{
		"project_name": state.ProjectName.ValueString(),
	})

	// Collect all label types to delete
	labelTypes := make([]string, 0, len(state.LabelSettings))
	for _, setting := range state.LabelSettings {
		labelTypes = append(labelTypes, setting.LabelType.ValueString())
	}

	err := r.client.DeleteLogLabels(
		state.ProjectName.ValueString(),
		r.client.Username,
		labelTypes,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Log Labels",
			"Could not delete log labels: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *logLabelsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using project_name as the ID
	resource.ImportStatePassthroughID(ctx, path.Root("project_name"), req, resp)
}

// validateAndConvertSettings validates log label strings are valid JSON arrays
func (r *logLabelsResource) validateAndConvertSettings(settings []logLabelSettingModel) ([]*client.LogLabelSetting, error) {
	result := make([]*client.LogLabelSetting, 0, len(settings))

	for i, setting := range settings {
		labelType := setting.LabelType.ValueString()
		logLabelString := setting.LogLabelString.ValueString()

		// Validate that logLabelString is a valid JSON array
		var labels []string
		if err := json.Unmarshal([]byte(logLabelString), &labels); err != nil {
			return nil, fmt.Errorf("label_settings[%d].log_label_string is not a valid JSON array: %s", i, err.Error())
		}

		result = append(result, &client.LogLabelSetting{
			LabelType:      labelType,
			LogLabelString: logLabelString,
			Labels:         labels,
		})
	}

	return result, nil
}
