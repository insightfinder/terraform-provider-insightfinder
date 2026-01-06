// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &insightfinderProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &insightfinderProvider{
			version: version,
		}
	}
}

// insightfinderProvider is the provider implementation.
type insightfinderProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// insightfinderProviderModel maps provider schema data to a Go type.
type insightfinderProviderModel struct {
	BaseURL    types.String `tfsdk:"base_url"`
	Username   types.String `tfsdk:"username"`
	LicenseKey types.String `tfsdk:"license_key"`
}

// Metadata returns the provider type name.
func (p *insightfinderProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "insightfinder"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *insightfinderProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with InsightFinder API to manage projects, configurations, and integrations.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "The base URL for the InsightFinder API. May also be provided via IF_BASE_URL environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username for InsightFinder authentication. May also be provided via IF_USERNAME environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"license_key": schema.StringAttribute{
				Description: "The license key (API key) for InsightFinder authentication. May also be provided via IF_LICENSE_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a InsightFinder API client for data sources and resources.
func (p *insightfinderProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring InsightFinder client")

	// Retrieve provider data from configuration
	var config insightfinderProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown InsightFinder API Base URL",
			"The provider cannot create the InsightFinder API client as there is an unknown configuration value for the base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IF_BASE_URL environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown InsightFinder API Username",
			"The provider cannot create the InsightFinder API client as there is an unknown configuration value for the username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IF_USERNAME environment variable.",
		)
	}

	if config.LicenseKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("license_key"),
			"Unknown InsightFinder API License Key",
			"The provider cannot create the InsightFinder API client as there is an unknown configuration value for the license key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IF_LICENSE_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	baseURL := os.Getenv("IF_BASE_URL")
	username := os.Getenv("IF_USERNAME")
	licenseKey := os.Getenv("IF_LICENSE_KEY")

	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.LicenseKey.IsNull() {
		licenseKey = config.LicenseKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if baseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing InsightFinder API Base URL",
			"The provider cannot create the InsightFinder API client as there is a missing or empty value for the base URL. "+
				"Set the base_url value in the configuration or use the IF_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing InsightFinder API Username",
			"The provider cannot create the InsightFinder API client as there is a missing or empty value for the username. "+
				"Set the username value in the configuration or use the IF_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if licenseKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("license_key"),
			"Missing InsightFinder API License Key",
			"The provider cannot create the InsightFinder API client as there is a missing or empty value for the license key. "+
				"Set the license_key value in the configuration or use the IF_LICENSE_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "insightfinder_base_url", baseURL)
	ctx = tflog.SetField(ctx, "insightfinder_username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "insightfinder_license_key")

	tflog.Debug(ctx, "Creating InsightFinder client")

	// Create a new InsightFinder client using the configuration values
	c, err := client.NewClient(baseURL, username, licenseKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create InsightFinder API Client",
			"An unexpected error occurred when creating the InsightFinder API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"InsightFinder Client Error: "+err.Error(),
		)
		return
	}

	// Make the InsightFinder client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured InsightFinder client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *insightfinderProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewSystemsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *insightfinderProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewLogLabelsResource,
		NewJWTConfigResource,
		NewServiceNowResource,
	}
}
