// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetricProjectResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMetricProjectResourceConfigBasic("test-metric-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "project_name", "test-metric-basic"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "project_creation_config.data_type", "Metric"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "project_creation_config.instance_type", "PrivateCloud"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "sampling_interval", "60"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "retention_time", "90"),
					resource.TestCheckResourceAttrSet("insightfinder_metric_project.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "insightfinder_metric_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccMetricProjectResourceConfigBasicUpdated("test-metric-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "sampling_interval", "300"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "retention_time", "180"),
				),
			},
		},
	})
}

func TestAccMetricProjectResource_DetectionTuning(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricProjectResourceConfigDetectionTuning("test-metric-tuned"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "project_name", "test-metric-tuned"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "c_value", "3"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "p_value", "0.95"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "dynamic_baseline_detection_flag", "true"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "enable_ubl_detect", "true"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "anomaly_dampening", "50400000"),
					resource.TestCheckResourceAttrSet("insightfinder_metric_project.test", "id"),
				),
			},
		},
	})
}

func TestAccMetricProjectResource_WithMetricConfigurations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with metric_configurations
			{
				Config: testAccMetricProjectResourceConfigWithMetricConfigurations("test-metric-configured"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "project_name", "test-metric-configured"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.metric_name", "cpu_usage"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.escalate_incident_components.0", "web-server-01"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.ignored_components.0", "test-instance"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.metric_alert_settings.0.threshold_alert_upper_bound", "95"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.metric_alert_settings.0.is_kpi", "true"),
					resource.TestCheckResourceAttrSet("insightfinder_metric_project.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "insightfinder_metric_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update metric_configurations
			{
				Config: testAccMetricProjectResourceConfigWithMetricConfigurationsUpdated("test-metric-configured"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.metric_alert_settings.0.threshold_alert_upper_bound", "90"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "metric_configurations.0.metric_alert_settings.0.is_kpi", "false"),
				),
			},
		},
	})
}

func TestAccMetricProjectResource_WithHolidays(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricProjectResourceConfigWithHolidays("test-metric-holidays"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "project_name", "test-metric-holidays"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "holiday_settings.0.name", "christmas"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "holiday_settings.0.start_date", "12-25"),
					resource.TestCheckResourceAttr("insightfinder_metric_project.test", "holiday_settings.0.end_date", "12-26"),
					resource.TestCheckResourceAttrSet("insightfinder_metric_project.test", "id"),
				),
			},
		},
	})
}

func testAccMetricProjectResourceConfigBasic(projectName string) string {
	return `
resource "insightfinder_metric_project" "test" {
  project_name = "` + projectName + `"
  system_name  = "test-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90
}
`
}

func testAccMetricProjectResourceConfigBasicUpdated(projectName string) string {
	return `
resource "insightfinder_metric_project" "test" {
  project_name = "` + projectName + `"
  system_name  = "test-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 300
  retention_time    = 180
}
`
}

func testAccMetricProjectResourceConfigDetectionTuning(projectName string) string {
	return `
resource "insightfinder_metric_project" "test" {
  project_name = "` + projectName + `"
  system_name  = "test-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "OnPremise"
    project_cloud_type = "OnPremise"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  c_value                              = 3
  p_value                              = 0.95
  high_ratio_c_value                   = 3
  maximum_hint                         = 20
  dynamic_baseline_detection_flag      = true
  baseline_duration                    = 14400000
  positive_baseline_violation_factor   = 2.0
  negative_baseline_violation_factor   = 2.0
  enable_period_anomaly_filter         = false
  enable_ubl_detect                    = true
  enable_cumulative_detect             = true
  filter_by_anomaly_in_baseline_generation = false
  anomaly_dampening                    = 50400000
  anomaly_gap_tolerance_count          = 1

  enable_fill_gap                  = false
  enable_store_filled_gap          = false
  gap_filling_training_data_length = 0
  enable_metric_data_prediction    = true

  prediction_training_data_length    = 0
  prediction_correlation_sensitivity = 0.75
  enable_kpi_prediction              = false
}
`
}

func testAccMetricProjectResourceConfigWithMetricConfigurations(projectName string) string {
	return `
resource "insightfinder_metric_project" "test" {
  project_name = "` + projectName + `"
  system_name  = "test-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  metric_configurations = [
    {
      metric_name                  = "cpu_usage"
      escalate_incident_components = ["web-server-01"]
      ignored_components           = ["test-instance"]
      metric_alert_settings = [
        {
          component_name                          = ""
          threshold_alert_lower_bound             = ""
          threshold_alert_upper_bound             = "95"
          threshold_alert_lower_bound_negative    = ""
          threshold_alert_upper_bound_negative    = ""
          threshold_no_alert_lower_bound          = ""
          threshold_no_alert_upper_bound          = "80"
          threshold_no_alert_lower_bound_negative = ""
          threshold_no_alert_upper_bound_negative = ""
          incident_alert_lower_bound              = ""
          incident_alert_upper_bound              = "99"
          incident_alert_lower_bound_negative     = ""
          incident_alert_upper_bound_negative     = ""
          incident_no_alert_lower_bound           = ""
          incident_no_alert_upper_bound           = "90"
          incident_no_alert_lower_bound_negative  = ""
          incident_no_alert_upper_bound_negative  = ""
          is_kpi                                  = true
          is_flapping_result_only                 = false
          incident_duration_threshold             = 300000
          detection_type                          = "Threshold"
          pattern_name_higher                     = "High CPU"
          pattern_name_lower                      = ""
          metric_type                             = "Gauge"
          fill_zero                               = false
          rouge_value                             = ""
          enable_baseline_near_constance          = false
          compute_difference                      = false
          anomaly_gap_tolerance_duration          = 0
        }
      ]
    }
  ]
}
`
}

func testAccMetricProjectResourceConfigWithMetricConfigurationsUpdated(projectName string) string {
	return `
resource "insightfinder_metric_project" "test" {
  project_name = "` + projectName + `"
  system_name  = "test-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  metric_configurations = [
    {
      metric_name                  = "cpu_usage"
      escalate_incident_components = ["web-server-01"]
      ignored_components           = []
      metric_alert_settings = [
        {
          component_name                          = ""
          threshold_alert_lower_bound             = ""
          threshold_alert_upper_bound             = "90"
          threshold_alert_lower_bound_negative    = ""
          threshold_alert_upper_bound_negative    = ""
          threshold_no_alert_lower_bound          = ""
          threshold_no_alert_upper_bound          = "75"
          threshold_no_alert_lower_bound_negative = ""
          threshold_no_alert_upper_bound_negative = ""
          incident_alert_lower_bound              = ""
          incident_alert_upper_bound              = "95"
          incident_alert_lower_bound_negative     = ""
          incident_alert_upper_bound_negative     = ""
          incident_no_alert_lower_bound           = ""
          incident_no_alert_upper_bound           = "85"
          incident_no_alert_lower_bound_negative  = ""
          incident_no_alert_upper_bound_negative  = ""
          is_kpi                                  = false
          is_flapping_result_only                 = false
          incident_duration_threshold             = 600000
          detection_type                          = "Threshold"
          pattern_name_higher                     = "High CPU"
          pattern_name_lower                      = ""
          metric_type                             = "Gauge"
          fill_zero                               = false
          rouge_value                             = ""
          enable_baseline_near_constance          = false
          compute_difference                      = false
          anomaly_gap_tolerance_duration          = 0
        }
      ]
    }
  ]
}
`
}

func testAccMetricProjectResourceConfigWithHolidays(projectName string) string {
	return `
resource "insightfinder_metric_project" "test" {
  project_name = "` + projectName + `"
  system_name  = "test-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  holiday_settings = [
    {
      name       = "christmas"
      start_date = "12-25"
      end_date   = "12-26"
    },
    {
      name       = "new_year"
      start_date = "01-01"
      end_date   = "01-01"
    }
  ]
}
`
}
