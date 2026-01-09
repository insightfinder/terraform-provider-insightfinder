// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig("test-project-1", "Test Project 1", "test-system-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_name", "test-project-1"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_display_name", "Test Project 1"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "system_name", "test-system-1"),
					resource.TestCheckResourceAttrSet("insightfinder_project.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "insightfinder_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig("test-project-1", "Test Project 1 Updated", "test-system-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_display_name", "Test Project 1 Updated"),
				),
			},
			// Delete testing automatically occurs at the end
		},
	})
}

func TestAccProjectResourceLogType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfigLogType("test-log-project", "Test Log Project", "test-system-log"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_name", "test-log-project"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_creation_config.data_type", "Log"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_creation_config.instance_type", "PrivateCloud"),
				),
			},
		},
	})
}

func TestAccProjectResourceWithAlerting(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfigWithAlerting("test-alert-project", "test-system-alert"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_name", "test-alert-project"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "email_config.0.email_recipients.0", "admin@example.com"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "email_config.0.enable_alert_email", "true"),
				),
			},
		},
	})
}

func TestAccProjectResourceWithLLMEvaluation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfigWithLLM("test-llm-project", "test-system-llm"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_project.test", "project_name", "test-llm-project"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "llm_evaluation_setting.0.is_hallucination_evaluation", "true"),
					resource.TestCheckResourceAttr("insightfinder_project.test", "llm_evaluation_setting.0.is_toxicity_evaluation", "true"),
				),
			},
		},
	})
}

func testAccProjectResourceConfig(projectName, displayName, systemName string) string {
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
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }
}
`, projectName, displayName, systemName)
}

func testAccProjectResourceConfigLogType(projectName, displayName, systemName string) string {
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
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }

  sampling_interval = 10
  project_time_zone = "UTC"
}
`, projectName, displayName, systemName)
}

func testAccProjectResourceConfigWithAlerting(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_project" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  project_creation_config = {
    data_type               = "Log"
    instance_type           = "PrivateCloud"
    project_cloud_type      = "PrivateCloud"
    insight_agent_type      = "LogStreaming"
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }

  email_config = {
    email_recipients    = ["admin@example.com", "ops@example.com"]
    enable_alert_email  = true
    enable_hourly_email = false
  }
}
`, projectName, systemName)
}

func testAccProjectResourceConfigWithLLM(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_project" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  project_creation_config = {
    data_type               = "Log"
    instance_type           = "PrivateCloud"
    project_cloud_type      = "PrivateCloud"
    insight_agent_type      = "LogStreaming"
    sampling_interval       = 10
    sampling_interval_in_seconds = 600
  }

  llm_evaluation_setting = {
    is_hallucination_evaluation      = true
    is_answer_relevant_evaluation    = false
    is_logic_consistency_evaluation  = false
    is_factual_inaccuracy_evaluation = false
    is_malicious_prompt_evaluation   = false
    is_toxicity_evaluation           = true
    is_pii_phi_leakage_evaluation    = false
    is_topic_guardrails_evaluation   = false
    is_tone_detection_evaluation     = false
  }
}
`, projectName, systemName)
}
