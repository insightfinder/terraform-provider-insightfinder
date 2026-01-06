// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLogLabelsResource_Whitelist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLogLabelsResourceConfigWhitelist("test-log-project", "test-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "project_name", "test-log-project"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "system_name", "test-system"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "whitelist.#", "2"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "whitelist.0.label_name", "important"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "whitelist.0.label_rule", "ERROR"),
					resource.TestCheckResourceAttrSet("insightfinder_log_labels.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "insightfinder_log_labels.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccLogLabelsResourceConfigWhitelistUpdated("test-log-project", "test-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "whitelist.#", "3"),
				),
			},
		},
	})
}

func TestAccLogLabelsResource_PatternNaming(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLogLabelsResourceConfigPatternNaming("pattern-project", "pattern-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "project_name", "pattern-project"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "pattern_naming.#", "2"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "pattern_naming.0.label_name", "database_error"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "pattern_naming.0.label_rule", "database.*connection.*failed"),
				),
			},
		},
	})
}

func TestAccLogLabelsResource_TrainingWhitelist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLogLabelsResourceConfigTrainingWhitelist("training-project", "training-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "project_name", "training-project"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "training_whitelist.#", "1"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "training_whitelist.0.label_name", "normal_operation"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "training_whitelist.0.label_rule", "INFO:.*started successfully"),
				),
			},
		},
	})
}

func TestAccLogLabelsResource_Combined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLogLabelsResourceConfigCombined("combined-project", "combined-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "whitelist.#", "1"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "pattern_naming.#", "1"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "training_whitelist.#", "1"),
				),
			},
		},
	})
}

func TestAccLogLabelsResource_MultipleRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLogLabelsResourceConfigMultipleRules("multi-rule-project", "multi-rule-system"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "whitelist.#", "5"),
					resource.TestCheckResourceAttr("insightfinder_log_labels.test", "pattern_naming.#", "3"),
				),
			},
		},
	})
}

func testAccLogLabelsResourceConfigWhitelist(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_log_labels" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  whitelist = [
    {
      label_name = "important"
      label_rule = "ERROR"
    },
    {
      label_name = "critical"
      label_rule = "CRITICAL|FATAL"
    }
  ]
}
`, projectName, systemName)
}

func testAccLogLabelsResourceConfigWhitelistUpdated(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_log_labels" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  whitelist = [
    {
      label_name = "important"
      label_rule = "ERROR"
    },
    {
      label_name = "critical"
      label_rule = "CRITICAL|FATAL"
    },
    {
      label_name = "warning"
      label_rule = "WARN"
    }
  ]
}
`, projectName, systemName)
}

func testAccLogLabelsResourceConfigPatternNaming(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_log_labels" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  pattern_naming = [
    {
      label_name = "database_error"
      label_rule = "database.*connection.*failed"
    },
    {
      label_name = "network_timeout"
      label_rule = "network.*timeout|connection.*timed out"
    }
  ]
}
`, projectName, systemName)
}

func testAccLogLabelsResourceConfigTrainingWhitelist(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_log_labels" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  training_whitelist = [
    {
      label_name = "normal_operation"
      label_rule = "INFO:.*started successfully"
    }
  ]
}
`, projectName, systemName)
}

func testAccLogLabelsResourceConfigCombined(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_log_labels" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  whitelist = [
    {
      label_name = "errors"
      label_rule = "ERROR|EXCEPTION"
    }
  ]

  pattern_naming = [
    {
      label_name = "auth_failure"
      label_rule = "authentication.*failed|login.*denied"
    }
  ]

  training_whitelist = [
    {
      label_name = "routine"
      label_rule = "INFO:.*routine check"
    }
  ]
}
`, projectName, systemName)
}

func testAccLogLabelsResourceConfigMultipleRules(projectName, systemName string) string {
	return fmt.Sprintf(`
resource "insightfinder_log_labels" "test" {
  project_name = %[1]q
  system_name  = %[2]q

  whitelist = [
    {
      label_name = "error"
      label_rule = "ERROR"
    },
    {
      label_name = "fatal"
      label_rule = "FATAL"
    },
    {
      label_name = "exception"
      label_rule = "Exception|Traceback"
    },
    {
      label_name = "warning"
      label_rule = "WARN"
    },
    {
      label_name = "critical"
      label_rule = "CRITICAL"
    }
  ]

  pattern_naming = [
    {
      label_name = "out_of_memory"
      label_rule = "OutOfMemory|OOM|memory.*exceeded"
    },
    {
      label_name = "disk_full"
      label_rule = "disk.*full|no space left"
    },
    {
      label_name = "connection_refused"
      label_rule = "connection refused|failed to connect"
    }
  ]
}
`, projectName, systemName)
}
