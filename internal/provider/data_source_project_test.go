// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create a project
			{
				Config: testAccProjectDataSourceConfigWithResource("datasource-test-project", "DataSource Test Project", "datasource-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that the resource was created
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_name", "datasource-test-project"),
					// Check that the data source can read it
					resource.TestCheckResourceAttr("data.insightfinder_project.test", "project_name", "datasource-test-project"),
					resource.TestCheckResourceAttr("data.insightfinder_project.test", "project_display_name", "DataSource Test Project"),
					resource.TestCheckResourceAttr("data.insightfinder_project.test", "system_name", "datasource-system"),
					resource.TestCheckResourceAttrSet("data.insightfinder_project.test", "id"),
					resource.TestCheckResourceAttrSet("data.insightfinder_project.test", "data_type"),
				),
			},
		},
	})
}

func TestAccProjectDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceConfigByName("query-by-name-project", "Query By Name Test", "query-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.insightfinder_project.test", "project_name", "query-by-name-project"),
					resource.TestCheckResourceAttr("data.insightfinder_project.test", "project_display_name", "Query By Name Test"),
				),
			},
		},
	})
}

func TestAccProjectDataSource_WithSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceConfigWithSettings("settings-project", "settings-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.insightfinder_project.test", "project_name", "settings-project"),
					resource.TestCheckResourceAttrSet("data.insightfinder_project.test", "sampling_interval"),
					resource.TestCheckResourceAttrSet("data.insightfinder_project.test", "project_time_zone"),
				),
			},
		},
	})
}

func TestAccProjectDataSource_NonExistent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectDataSourceConfigNonExistent(),
				ExpectError: nil, // Will return error from API
			},
		},
	})
}

func testAccProjectDataSourceConfigWithResource(projectName, displayName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_project" "test" {
  project_name         = %[1]q
  project_display_name = %[2]q
  system_name          = %[3]q

  project_creation_config = {
    data_type               = "Log"
    instance_type           = "PrivateCloud"
    project_cloud_type      = "PrivateCloud"
    insight_agent_type      = "LogStreaming"
    project_creation_type   = "Kafka"
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }
}

data "insightfinder_project" "test" {
  project_name = insightfinder_project.test.project_name
  
  depends_on = [insightfinder_project.test]
}
`, projectName, displayName, systemName)
}

func testAccProjectDataSourceConfigByName(projectName, displayName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_project" "test" {
  project_name         = %[1]q
  project_display_name = %[2]q
  system_name          = %[3]q

  project_creation_config = {
    data_type               = "Log"
    instance_type           = "PrivateCloud"
    project_cloud_type      = "PrivateCloud"
    insight_agent_type      = "LogStreaming"
    project_creation_type   = "Kafka"
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }
}

data "insightfinder_project" "test" {
  project_name = %[1]q
  
  depends_on = [insightfinder_project.test]
}
`, projectName, displayName, systemName)
}

func testAccProjectDataSourceConfigWithSettings(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_project" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  project_creation_config = {
    data_type               = "Log"
    instance_type           = "PrivateCloud"
    project_cloud_type      = "PrivateCloud"
    insight_agent_type      = "LogStreaming"
    project_creation_type   = "Kafka"
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }

  sampling_interval = 10
  project_time_zone = "America/New_York"
}

data "insightfinder_project" "test" {
  project_name = insightfinder_project.test.project_name
  
  depends_on = [insightfinder_project.test]
}
`, projectName, systemName)
}

func testAccProjectDataSourceConfigNonExistent() string {
	return `
data "insightfinder_project" "test" {
  project_name = "this-project-should-not-exist-12345"
}
`
}
