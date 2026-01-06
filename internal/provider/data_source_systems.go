// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &systemsDataSource{}
	_ datasource.DataSourceWithConfigure = &systemsDataSource{}
)

// NewSystemsDataSource is a helper function to simplify the provider implementation.
func NewSystemsDataSource() datasource.DataSource {
	return &systemsDataSource{}
}

// systemsDataSource is the data source implementation.
type systemsDataSource struct {
	client *client.Client
}

// systemsDataSourceModel maps the data source schema data.
type systemsDataSourceModel struct {
	ID      types.String  `tfsdk:"id"`
	Systems []systemModel `tfsdk:"systems"`
}

// systemModel represents a single system
type systemModel struct {
	SystemID   types.String `tfsdk:"system_id"`
	SystemName types.String `tfsdk:"system_name"`
}

// Metadata returns the data source type name.
func (d *systemsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_systems"
}

// Schema defines the schema for the data source.
func (d *systemsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of systems from InsightFinder.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier for the data source.",
				Computed:    true,
			},
			"systems": schema.ListNestedAttribute{
				Description: "List of systems.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"system_id": schema.StringAttribute{
							Description: "The unique identifier for the system.",
							Computed:    true,
						},
						"system_name": schema.StringAttribute{
							Description: "The name of the system.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *systemsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *systemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state systemsDataSourceModel

	tflog.Debug(ctx, "Reading systems list")

	// Get system framework
	systemFramework, err := d.client.GetSystemFramework(d.client.Username, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Systems",
			"Could not read systems: "+err.Error(),
		)
		return
	}

	if systemFramework == nil || systemFramework.OwnSystemArr == nil {
		// No systems found
		state.ID = types.StringValue("systems")
		state.Systems = []systemModel{}
		diags := resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	// Parse systems
	systems := make([]systemModel, 0)
	for _, systemStr := range systemFramework.OwnSystemArr {
		var system client.SystemFramework
		if err := json.Unmarshal([]byte(systemStr), &system); err != nil {
			tflog.Warn(ctx, "Failed to parse system", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		systems = append(systems, systemModel{
			SystemID:   types.StringValue(system.SystemID),
			SystemName: types.StringValue(system.SystemName),
		})
	}

	// Set state
	state.ID = types.StringValue("systems")
	state.Systems = systems

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
