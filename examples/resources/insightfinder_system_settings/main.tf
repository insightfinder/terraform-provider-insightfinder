terraform {
  required_providers {
    insightfinder = {
      source  = "insightfinder/insightfinder"
      version = "~> 1.8"
    }
  }
}

provider "insightfinder" {
  base_url    = "https://app.insightfinder.com"
  username    = var.username
  license_key = var.license_key
}

# Knowledge base settings only
resource "insightfinder_system_settings" "kb_only" {
  system_name = "my-production-system"

  knowledgebase_settings = {
    enable_global_knowledge_base      = true
    composite_valid_threshold         = 900000
    timeline_top_k                    = 50
    enable_ignore_instance_prediction = true
    prediction_source                 = 0
    share_system_type                 = 1
    action_execution_time             = 15
    auto_fix_validation_window        = 1
    filter_self_to_self               = true
    rule_source_type                  = 0
    satellite_system_set              = jsonencode([])

    rule_active_threshold            = 0.8
    rule_inactive_threshold          = 0.1
    rule_active_condition            = 0
    false_positive_tolerance         = 1
    kb_training_length               = 172800000
    tolerance                        = 0.08
    enable_insensitive_rule_matching = true
  }
}

# Notifications with system down, daily/weekly reports, and instance down alerts
resource "insightfinder_system_settings" "with_notifications" {
  system_name = "my-production-system"

  notifications_settings = {
    prediction_email                       = "alerts@example.com"
    enable_incident_prediction_email_alert = true
    enable_incident_detection_email_alert  = true
    enable_root_cause_email_alert          = true
    root_cause_email                       = "rca@example.com"
    incident_detection_email               = "incidents@example.com"
    alert_email                            = "alerts@example.com"
    health_alert_email                     = "health@example.com"
    alert_health_score                     = 0.5
    aggregation_interval                   = 10
    email_dampening_period                 = 59280000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 59220000
    order                                  = 0
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_system_down_email_alert         = false
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    alert_frequency                        = 0
    incident_count_threshold               = jsonencode({})
    assignment_map                         = jsonencode({})

    # System down notifications (via dedicated /api/external/v2/systemdownsetting API)
    system_down_notification = {
      enable_system_down_email_alert = true
      email_dampening_period         = 3600000
      email_set                      = ["oncall@example.com", "ops@example.com"]
    }

    # Daily insights report emails
    daily_report_notification = {
      enable_insights_report = true
      email_set              = ["manager@example.com", "team@example.com"]
    }

    # Weekly insights report emails
    weekly_report_notification = {
      enable_insights_report = true
      email_set              = ["executive@example.com", "manager@example.com"]
    }

    # Per-project instance down alerts
    instance_down_notification = [
      {
        project_name                = "production-metrics"
        instance_down_enable        = true
        instance_down_dampening     = 1800000
        instance_down_threshold     = 300000
        instance_down_report_number = 2
        instance_down_emails        = ["oncall@example.com"]
      },
      {
        project_name                = "staging-metrics"
        instance_down_enable        = false
        instance_down_dampening     = 3600000
        instance_down_threshold     = 600000
        instance_down_report_number = 5
        instance_down_emails        = []
      }
    ]

    # Project-level dampening window overrides
    project_level_dampening_windows = [
      {
        source_project = "change-detection-project"
        target_project = "change-detection-project"
        duration       = 21600000
      },
      {
        source_project  = "llm-trace-project"
        target_project  = "change-detection-project"
        source_customer = "admin"
        target_customer = "admin"
        duration        = 28800000
      }
    ]
  }
}

# Full configuration with both KB and extended notifications
resource "insightfinder_system_settings" "full" {
  system_name = "my-production-system"

  knowledgebase_settings = {
    enable_global_knowledge_base      = true
    composite_valid_threshold         = 900000
    timeline_top_k                    = 50
    enable_ignore_instance_prediction = true
    prediction_source                 = 0
    share_system_type                 = 1
    action_execution_time             = 15
    auto_fix_validation_window        = 1
    filter_self_to_self               = true
    rule_source_type                  = 0
    satellite_system_set              = jsonencode([])

    rule_active_threshold            = 0.8
    rule_inactive_threshold          = 0.1
    rule_active_condition            = 0
    false_positive_tolerance         = 1
    kb_training_length               = 172800000
    tolerance                        = 0.08
    enable_insensitive_rule_matching = true
  }

  notifications_settings = {
    prediction_email                       = "alerts@example.com"
    enable_incident_prediction_email_alert = true
    enable_incident_detection_email_alert  = true
    alert_health_score                     = 0.5
    aggregation_interval                   = 10
    email_dampening_period                 = 59280000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 59220000
    order                                  = 0
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_system_down_email_alert         = false
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    alert_frequency                        = 0
    incident_count_threshold               = jsonencode({})
    assignment_map                         = jsonencode({})

    system_down_notification = {
      enable_system_down_email_alert = true
      email_dampening_period         = 3600000
      email_set                      = ["oncall@example.com"]
    }

    daily_report_notification = {
      enable_insights_report = true
      email_set              = ["reports@example.com"]
    }

    weekly_report_notification = {
      enable_insights_report = true
      email_set              = ["reports@example.com"]
    }

    instance_down_notification = [
      {
        project_name                = "production-metrics"
        instance_down_enable        = true
        instance_down_dampening     = 1800000
        instance_down_threshold     = 300000
        instance_down_report_number = 1
        instance_down_emails        = ["oncall@example.com"]
      }
    ]

    project_level_dampening_windows = [
      {
        source_project = "production-metrics"
        target_project = "production-metrics"
        duration       = 21600000
      }
    ]
  }

  miscellaneous_settings = {
    healthview_longterm                      = false
    should_auto_share                        = true
    rootcause_reverse_entry_filter_threshold = 99
    enable_composite_timeline                = true
  }
}

variable "username" {
  description = "InsightFinder username"
  type        = string
  sensitive   = true
}

variable "license_key" {
  description = "InsightFinder license key"
  type        = string
  sensitive   = true
}
