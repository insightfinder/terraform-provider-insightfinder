// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemSettingsResource_KnowledgebaseOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSystemSettingsResourceConfigKBOnly("test-system-settings"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "system_name", "test-system-settings"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.enable_global_knowledge_base", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.composite_valid_threshold", "86400000"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.timeline_top_k", "5"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.rule_active_threshold", "0.8"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.rule_inactive_threshold", "0.2"),
					resource.TestCheckResourceAttrSet("insightfinder_system_settings.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "insightfinder_system_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccSystemSettingsResourceConfigKBUpdated("test-system-settings"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.enable_global_knowledge_base", "false"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.timeline_top_k", "10"),
				),
			},
		},
	})
}

func TestAccSystemSettingsResource_NotificationsOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemSettingsResourceConfigNotificationsOnly("test-system-notif"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "system_name", "test-system-notif"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.prediction_email", "alert@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.enable_system_down_email_alert", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.alert_health_score", "0.5"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.aggregation_interval", "5"),
					resource.TestCheckResourceAttrSet("insightfinder_system_settings.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "insightfinder_system_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSystemSettingsResource_BothSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemSettingsResourceConfigBoth("test-system-full"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "system_name", "test-system-full"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "knowledgebase_settings.enable_global_knowledge_base", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.prediction_email", "full@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.enable_alerts_email", "true"),
					resource.TestCheckResourceAttrSet("insightfinder_system_settings.test", "id"),
				),
			},
			// Update notifications
			{
				Config: testAccSystemSettingsResourceConfigBothUpdated("test-system-full"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.prediction_email", "updated@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.enable_system_down_email_alert", "true"),
				),
			},
		},
	})
}

func TestAccSystemSettingsResource_EmailSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemSettingsResourceConfigEmailSettings("test-system-email"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.alert_email", "alerts@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.health_alert_email", "health@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.incident_detection_email", "incidents@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.root_cause_email", "rca@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.enable_root_cause_email_alert", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.enable_incident_detection_email_alert", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.enable_incident_prediction_email_alert", "true"),
				),
			},
		},
	})
}

func testAccSystemSettingsResourceConfigKBOnly(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  knowledgebase_settings = {
    enable_global_knowledge_base      = true
    composite_valid_threshold         = 86400000
    timeline_top_k                    = 5
    enable_ignore_instance_prediction = false
    prediction_source                 = 0
    share_system_type                 = 0
    action_execution_time             = 30
    auto_fix_validation_window        = 10
    filter_self_to_self               = false
    rule_source_type                  = 0

    rule_active_threshold           = 0.8
    rule_inactive_threshold         = 0.2
    rule_active_condition           = 0
    false_positive_tolerance        = 3
    kb_training_length              = 604800000
    tolerance                       = 0.1
    enable_insensitive_rule_matching = false
  }
}
`
}

func testAccSystemSettingsResourceConfigKBUpdated(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  knowledgebase_settings = {
    enable_global_knowledge_base      = false
    composite_valid_threshold         = 86400000
    timeline_top_k                    = 10
    enable_ignore_instance_prediction = false
    prediction_source                 = 0
    share_system_type                 = 0
    action_execution_time             = 30
    auto_fix_validation_window        = 10
    filter_self_to_self               = false
    rule_source_type                  = 0

    rule_active_threshold           = 0.8
    rule_inactive_threshold         = 0.2
    rule_active_condition           = 0
    false_positive_tolerance        = 3
    kb_training_length              = 604800000
    tolerance                       = 0.1
    enable_insensitive_rule_matching = false
  }
}
`
}

func testAccSystemSettingsResourceConfigNotificationsOnly(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  notifications_settings = {
    prediction_email                    = "alert@example.com"
    enable_system_down_email_alert      = true
    alert_health_score                  = 0.5
    aggregation_interval                = 5
    alert_frequency                     = 1
    email_dampening_period              = 3600000
    alerts_email_dampening_period       = 3600000
    prediction_email_dampening_period   = 3600000
    incident_dampening_window           = 3600000
    hide_flag                           = false
    enable_splunk_export                = false
    only_send_with_rca                  = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_alerts_email                 = false
    enable_health_email_alert           = false
    enable_root_cause_email_alert       = false
    order                               = 0
  }
}
`
}

func testAccSystemSettingsResourceConfigBoth(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  knowledgebase_settings = {
    enable_global_knowledge_base      = true
    composite_valid_threshold         = 86400000
    timeline_top_k                    = 5
    enable_ignore_instance_prediction = false
    prediction_source                 = 0
    share_system_type                 = 0
    action_execution_time             = 30
    auto_fix_validation_window        = 10
    filter_self_to_self               = false
    rule_source_type                  = 0

    rule_active_threshold           = 0.8
    rule_inactive_threshold         = 0.2
    rule_active_condition           = 0
    false_positive_tolerance        = 3
    kb_training_length              = 604800000
    tolerance                       = 0.1
    enable_insensitive_rule_matching = false
  }

  notifications_settings = {
    prediction_email                    = "full@example.com"
    enable_alerts_email                 = true
    enable_system_down_email_alert      = false
    alert_health_score                  = 0.5
    aggregation_interval                = 5
    alert_frequency                     = 1
    email_dampening_period              = 3600000
    alerts_email_dampening_period       = 3600000
    prediction_email_dampening_period   = 3600000
    incident_dampening_window           = 3600000
    hide_flag                           = false
    enable_splunk_export                = false
    only_send_with_rca                  = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_health_email_alert           = false
    enable_root_cause_email_alert       = false
    order                               = 0
  }
}
`
}

func testAccSystemSettingsResourceConfigBothUpdated(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  knowledgebase_settings = {
    enable_global_knowledge_base      = true
    composite_valid_threshold         = 86400000
    timeline_top_k                    = 5
    enable_ignore_instance_prediction = false
    prediction_source                 = 0
    share_system_type                 = 0
    action_execution_time             = 30
    auto_fix_validation_window        = 10
    filter_self_to_self               = false
    rule_source_type                  = 0

    rule_active_threshold           = 0.8
    rule_inactive_threshold         = 0.2
    rule_active_condition           = 0
    false_positive_tolerance        = 3
    kb_training_length              = 604800000
    tolerance                       = 0.1
    enable_insensitive_rule_matching = false
  }

  notifications_settings = {
    prediction_email                    = "updated@example.com"
    enable_system_down_email_alert      = true
    enable_alerts_email                 = true
    alert_health_score                  = 0.6
    aggregation_interval                = 5
    alert_frequency                     = 1
    email_dampening_period              = 3600000
    alerts_email_dampening_period       = 3600000
    prediction_email_dampening_period   = 3600000
    incident_dampening_window           = 3600000
    hide_flag                           = false
    enable_splunk_export                = false
    only_send_with_rca                  = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_health_email_alert           = false
    enable_root_cause_email_alert       = false
    order                               = 0
  }
}
`
}

func testAccSystemSettingsResourceConfigEmailSettings(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  notifications_settings = {
    alert_email                         = "alerts@example.com"
    health_alert_email                  = "health@example.com"
    incident_detection_email            = "incidents@example.com"
    root_cause_email                    = "rca@example.com"
    enable_root_cause_email_alert       = true
    enable_incident_detection_email_alert = true
    enable_incident_prediction_email_alert = true
    prediction_email                    = ""
    enable_system_down_email_alert      = false
    alert_health_score                  = 0.5
    aggregation_interval                = 5
    alert_frequency                     = 1
    email_dampening_period              = 3600000
    alerts_email_dampening_period       = 3600000
    prediction_email_dampening_period   = 3600000
    incident_dampening_window           = 3600000
    hide_flag                           = false
    enable_splunk_export                = false
    only_send_with_rca                  = false
    enable_alerts_email                 = false
    enable_health_email_alert           = false
    order                               = 0
  }
}
`
}

func TestAccSystemSettingsResource_SystemDownNotification(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemSettingsResourceConfigSystemDown("test-system-sysdown"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "system_name", "test-system-sysdown"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.system_down_notification.enable_system_down_email_alert", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.system_down_notification.email_dampening_period", "3600000"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.system_down_notification.email_set.0", "ops@example.com"),
					resource.TestCheckResourceAttrSet("insightfinder_system_settings.test", "id"),
				),
			},
			{
				ResourceName:      "insightfinder_system_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSystemSettingsResource_DailyWeeklyReport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemSettingsResourceConfigReports("test-system-reports"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "system_name", "test-system-reports"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.daily_report_notification.enable_insights_report", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.daily_report_notification.email_set.0", "reports@example.com"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.weekly_report_notification.enable_insights_report", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.weekly_report_notification.email_set.0", "exec@example.com"),
					resource.TestCheckResourceAttrSet("insightfinder_system_settings.test", "id"),
				),
			},
		},
	})
}

func TestAccSystemSettingsResource_InstanceDownNotification(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemSettingsResourceConfigInstanceDown("test-system-instdown"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "system_name", "test-system-instdown"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.instance_down_notification.0.project_name", "my-project"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.instance_down_notification.0.instance_down_enable", "true"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.instance_down_notification.0.instance_down_threshold", "300000"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.instance_down_notification.0.instance_down_emails.0", "oncall@example.com"),
					resource.TestCheckResourceAttrSet("insightfinder_system_settings.test", "id"),
				),
			},
			{
				Config: testAccSystemSettingsResourceConfigInstanceDownUpdated("test-system-instdown"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.instance_down_notification.0.instance_down_enable", "false"),
					resource.TestCheckResourceAttr("insightfinder_system_settings.test", "notifications_settings.instance_down_notification.0.instance_down_threshold", "600000"),
				),
			},
		},
	})
}

func testAccSystemSettingsResourceConfigSystemDown(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  notifications_settings = {
    prediction_email                       = ""
    enable_system_down_email_alert         = false
    alert_health_score                     = 0.5
    aggregation_interval                   = 5
    alert_frequency                        = 0
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 3600000
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    order                                  = 0

    system_down_notification = {
      enable_system_down_email_alert = true
      email_dampening_period         = 3600000
      email_set                      = ["ops@example.com"]
    }
  }
}
`
}

func testAccSystemSettingsResourceConfigReports(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  notifications_settings = {
    prediction_email                       = ""
    enable_system_down_email_alert         = false
    alert_health_score                     = 0.5
    aggregation_interval                   = 5
    alert_frequency                        = 0
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 3600000
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    order                                  = 0

    daily_report_notification = {
      enable_insights_report = true
      email_set              = ["reports@example.com"]
    }

    weekly_report_notification = {
      enable_insights_report = true
      email_set              = ["exec@example.com"]
    }
  }
}
`
}

func testAccSystemSettingsResourceConfigInstanceDown(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  notifications_settings = {
    prediction_email                       = ""
    enable_system_down_email_alert         = false
    alert_health_score                     = 0.5
    aggregation_interval                   = 5
    alert_frequency                        = 0
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 3600000
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    order                                  = 0

    instance_down_notification = [
      {
        project_name                = "my-project"
        instance_down_enable        = true
        instance_down_dampening     = 1800000
        instance_down_threshold     = 300000
        instance_down_report_number = 2
        instance_down_emails        = ["oncall@example.com"]
      }
    ]
  }
}
`
}

func testAccSystemSettingsResourceConfigInstanceDownUpdated(systemName string) string {
	return `
resource "insightfinder_system_settings" "test" {
  system_name = "` + systemName + `"

  notifications_settings = {
    prediction_email                       = ""
    enable_system_down_email_alert         = false
    alert_health_score                     = 0.5
    aggregation_interval                   = 5
    alert_frequency                        = 0
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 3600000
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_incident_prediction_email_alert = false
    enable_incident_detection_email_alert  = false
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    order                                  = 0

    instance_down_notification = [
      {
        project_name                = "my-project"
        instance_down_enable        = false
        instance_down_dampening     = 3600000
        instance_down_threshold     = 600000
        instance_down_report_number = 5
        instance_down_emails        = []
      }
    ]
  }
}
`
}
