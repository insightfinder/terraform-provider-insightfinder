// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"insightfinder": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add pre-checks for environment variables, network access, etc.
	if v := os.Getenv("IF_USERNAME"); v == "" {
		t.Skip("IF_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("IF_LICENSE_KEY"); v == "" {
		t.Skip("IF_LICENSE_KEY must be set for acceptance tests")
	}
}

func TestProvider(t *testing.T) {
	// Test that the provider can be instantiated
	p := New("test")()
	if p == nil {
		t.Fatal("Expected provider instance, got nil")
	}

	// Verify provider implements the interface
	var _ provider.Provider = p
}

func TestProviderMetadata(t *testing.T) {
	p := New("1.0.0")()

	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(context.Background(), req, resp)

	if resp.TypeName != "insightfinder" {
		t.Errorf("Expected type name 'insightfinder', got '%s'", resp.TypeName)
	}

	if resp.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", resp.Version)
	}
}

func TestProviderSchema(t *testing.T) {
	p := New("test")()

	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() produced errors: %v", resp.Diagnostics.Errors())
	}

	// Verify required attributes exist
	if _, ok := resp.Schema.Attributes["base_url"]; !ok {
		t.Error("Expected 'base_url' attribute in schema")
	}
	if _, ok := resp.Schema.Attributes["username"]; !ok {
		t.Error("Expected 'username' attribute in schema")
	}
	if _, ok := resp.Schema.Attributes["license_key"]; !ok {
		t.Error("Expected 'license_key' attribute in schema")
	}
}

func TestProviderConfigure(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
	}{
		{
			name: "configure with environment variables",
			setupEnv: func() {
				os.Setenv("IF_BASE_URL", "https://test.insightfinder.com")
				os.Setenv("IF_USERNAME", "test_user")
				os.Setenv("IF_LICENSE_KEY", "test_key")
			},
			cleanupEnv: func() {
				os.Unsetenv("IF_BASE_URL")
				os.Unsetenv("IF_USERNAME")
				os.Unsetenv("IF_LICENSE_KEY")
			},
			expectError: false,
		},
		{
			name: "missing username",
			setupEnv: func() {
				os.Setenv("IF_BASE_URL", "https://test.insightfinder.com")
				os.Setenv("IF_LICENSE_KEY", "test_key")
				os.Unsetenv("IF_USERNAME")
			},
			cleanupEnv: func() {
				os.Unsetenv("IF_BASE_URL")
				os.Unsetenv("IF_LICENSE_KEY")
			},
			expectError: true,
		},
		{
			name: "missing license key",
			setupEnv: func() {
				os.Setenv("IF_BASE_URL", "https://test.insightfinder.com")
				os.Setenv("IF_USERNAME", "test_user")
				os.Unsetenv("IF_LICENSE_KEY")
			},
			cleanupEnv: func() {
				os.Unsetenv("IF_BASE_URL")
				os.Unsetenv("IF_USERNAME")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			// Note: Full provider configuration testing would require
			// implementing the provider.ConfigureRequest properly.
			// This is a basic structure for unit testing.
		})
	}
}

func TestProviderResources(t *testing.T) {
	p := New("test")().(interface {
		Resources(context.Context) []func() resource.Resource
	})

	resources := p.Resources(context.Background())

	expectedResources := map[string]bool{
		"insightfinder_project":    false,
		"insightfinder_servicenow": false,
		"insightfinder_jwt_config": false,
		"insightfinder_log_labels": false,
	}

	if len(resources) != len(expectedResources) {
		t.Errorf("Expected %d resources, got %d", len(expectedResources), len(resources))
	}

	// Check that all expected resources are registered
	for _, resourceFunc := range resources {
		r := resourceFunc()
		// Resources don't have a direct way to get their type name without metadata
		// This is a limitation of the test - in real scenarios, use acceptance tests
		t.Logf("Found resource: %T", r)
	}
}

func TestProviderDataSources(t *testing.T) {
	p := New("test")().(interface {
		DataSources(context.Context) []func() datasource.DataSource
	})

	dataSources := p.DataSources(context.Background())

	expectedCount := 2 // insightfinder_project, insightfinder_systems

	if len(dataSources) != expectedCount {
		t.Errorf("Expected %d data sources, got %d", expectedCount, len(dataSources))
	}

	// Check that data sources can be instantiated
	for _, dsFunc := range dataSources {
		ds := dsFunc()
		if ds == nil {
			t.Error("Data source function returned nil")
		}
		t.Logf("Found data source: %T", ds)
	}
}
