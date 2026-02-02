// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "id"),
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.#"),
				),
			},
		},
	})
}

func TestAccSystemsDataSource_WithMultipleSystems(t *testing.T) {
	// This test assumes there are existing systems in the account
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.#"),
					// Verify systems have required attributes
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.0.system_name"),
				),
			},
		},
	})
}

func TestAccSystemsDataSource_AfterCreatingProject(t *testing.T) {
	// Create a project which creates a system, then query systems
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemsDataSourceConfigWithProject(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check project was created
					resource.TestCheckResourceAttr("insightfinder_project.test", "system_name", "systems-test-system"),
					// Check data source returns systems
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.#"),
				),
			},
		},
	})
}

func TestAccSystemsDataSource_VerifyStructure(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemsDataSourceConfigWithProject(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify system structure
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.#"),
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.0.system_name"),
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.test", "systems.0.system_id"),
				),
			},
		},
	})
}

func TestAccSystemsDataSource_DynamicUsage(t *testing.T) {
	// Test using systems data source to dynamically create resources
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemsDataSourceConfigDynamic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.all", "systems.#"),
					resource.TestCheckResourceAttrSet("data.insightfinder_systems.all", "id"),
				),
			},
		},
	})
}

func testAccSystemsDataSourceConfig() string {
	return `
data "insightfinder_systems" "test" {}
`
}

func testAccSystemsDataSourceConfigWithProject() string {
	return `
resource "insightfinder_project" "test" {
  project_name = "systems-datasource-test-project"
  system_name  = "systems-test-system"

  project_creation_config = {
    data_type               = "Log"
    instance_type           = "PrivateCloud"
    project_cloud_type      = "PrivateCloud"
    insight_agent_type      = "LogStreaming"
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }
}

data "insightfinder_systems" "test" {
  depends_on = [insightfinder_project.test]
}
`
}

func testAccSystemsDataSourceConfigDynamic() string {
	return `
# Query all systems
data "insightfinder_systems" "all" {}

# Example of how systems could be used dynamically
# (This is a structural test, not a full dynamic resource creation)
output "system_count" {
  value = length(data.insightfinder_systems.all.systems)
}

output "system_names" {
  value = [for system in data.insightfinder_systems.all.systems : system.system_name]
}
`
}
