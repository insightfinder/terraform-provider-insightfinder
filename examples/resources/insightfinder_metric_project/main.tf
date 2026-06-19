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

# Basic metric project
resource "insightfinder_metric_project" "basic_metrics" {
  project_name = "application-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  project_display_name = "Application Metrics"
  project_time_zone    = "UTC"
  sampling_interval    = 60
  retention_time       = 90
  ubl_retention_time   = 8
}

# Metric project with full detection tuning
resource "insightfinder_metric_project" "tuned_metrics" {
  project_name = "tuned-infrastructure-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "OnPremise"
    project_cloud_type = "OnPremise"
    insight_agent_type = "Custom"
  }

  project_display_name = "Tuned Infrastructure Metrics"
  project_time_zone    = "US/Eastern"
  sampling_interval    = 60
  retention_time       = 90
  ubl_retention_time   = 8

  # Detection tuning
  c_value                                  = 3
  p_value                                  = 0.95
  high_ratio_c_value                       = 3
  maximum_hint                             = 20
  dynamic_baseline_detection_flag          = true
  baseline_duration                        = 14400000
  positive_baseline_violation_factor       = 2.0
  negative_baseline_violation_factor       = 2.0
  enable_period_anomaly_filter             = false
  enable_ubl_detect                        = true
  enable_cumulative_detect                 = true
  enable_component_level_detection         = true
  filter_by_anomaly_in_baseline_generation = false
  anomaly_dampening                        = 50400000
  anomaly_gap_tolerance_count              = 1
  instance_down_ratio_threshold            = 0.05

  # Gap filling
  enable_fill_gap                  = false
  enable_store_filled_gap          = false
  gap_filling_training_data_length = 0
  enable_metric_data_prediction    = true

  # KPI prediction
  prediction_training_data_length    = 0
  prediction_correlation_sensitivity = 0.75
  enable_kpi_prediction              = false

  # Incident prediction and RCA
  incident_prediction_window       = 12
  min_incident_prediction_window   = 5
  incident_relation_search_window  = 6
  incident_prediction_event_limit  = 50
  root_cause_count_threshold       = 1
  root_cause_probability_threshold = 0.8
  maximum_root_cause_result_size   = 5
  multi_hop_search_level           = 2
  multi_hop_search_limit           = "30"

  # Instance down
  instance_down_threshold     = 3600000
  instance_down_report_number = 50
  instance_down_enable        = false
  show_instance_down          = false

  # Alerting
  avg_per_incident_downtime_cost = 5000.0
  alert_hourly_cost              = 200.0

  # Email
  enable_new_alert_email = true
  email_setting = jsonencode({
    onlySendWithRCA                    = false
    enableIncidentDetectionEmailAlert  = true
    enableIncidentPredictionEmailAlert = true
    enableRootCauseEmailAlert          = true
    emailDampeningPeriod               = 3600000
    awSeverityLevel                    = "Major"
  })

  # Webhook
  webhook_url = "https://hooks.example.com/metrics-alerts"
  webhook_type_set_str = jsonencode([
    "metric",
    "detectedIncident",
    "predictedIncident"
  ])
  webhook_alert_dampening = 18000000
}

# Metric project with holiday settings
resource "insightfinder_metric_project" "metrics_with_holidays" {
  project_name = "metrics-with-holidays"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  project_display_name = "Metrics with Holidays"
  sampling_interval    = 60
  retention_time       = 90

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
    },
    {
      name       = "independence_day"
      start_date = "07-04"
      end_date   = "07-04"
    }
  ]
}

# Metric project with per-metric alert threshold and component configurations
resource "insightfinder_metric_project" "with_metric_config" {
  project_name = "configured-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  metric_configurations = {
    "cpu_usage" = {
      escalate_incident_components = ["web-server-01", "web-server-02"]
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
          detection_type                          = "positive"
          c_value_override                        = null
          high_c_value_override                   = null
          pattern_name_higher                     = "High CPU"
          pattern_name_lower                      = ""
          metric_type                             = "CPU Utilization"
          fill_zero                               = false
          rouge_value                             = null
          enable_baseline_near_constance          = false
          compute_difference                      = false
          anomaly_gap_tolerance_duration          = 60000
        }
      ]
    }
    "memory_usage" = {
      escalate_incident_components = []
      ignored_components           = []
      metric_alert_settings        = []
    }
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

# Metric project with incident priority by anomaly score
resource "insightfinder_metric_project" "priority_example" {
  project_name = "priority-configured-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  project_display_name = "Priority Configured Metrics"
  project_time_zone    = "UTC"
  sampling_interval    = 60
  retention_time       = 90

  # Map anomaly score ranges to incident priority levels 1 (highest) through 5 (lowest)
  incident_priority_by_anomaly_score_setting = jsonencode({
    enabled = true
    priorityScoreRangeMap = {
      "1" = "10001-"
      "2" = "5001-10000"
      "3" = "2001-5000"
      "4" = "1001-2000"
      "5" = "0-1000"
    }
  })
}
