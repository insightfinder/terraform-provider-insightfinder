// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *client.Client
}

type projectDataSourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	ProjectName        types.String  `tfsdk:"project_name"`
	ProjectDisplayName types.String  `tfsdk:"project_display_name"`
	CValue             types.Int64   `tfsdk:"c_value"`
	PValue             types.Float64 `tfsdk:"p_value"`
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an InsightFinder project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_name": schema.StringAttribute{
				Description: "The name of the project to fetch.",
				Required:    true,
			},
			"project_display_name": schema.StringAttribute{
				Computed: true,
			},
			"c_value": schema.Int64Attribute{
				Computed: true,
			},
			"p_value": schema.Float64Attribute{
				Computed: true,
			},
		},
	}
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "")
		return
	}

	d.client = client
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data projectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	project, err := d.client.GetProject(data.ProjectName.ValueString(), d.client.Username)
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	if project == nil {
		resp.Diagnostics.AddError("Project not found", "")
		return
	}

	data.ID = types.StringValue(project.ProjectName)
	data.ProjectDisplayName = types.StringValue(project.ProjectDisplayName)
	data.CValue = types.Int64Value(int64(project.CValue))
	data.PValue = types.Float64Value(project.PValue)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
