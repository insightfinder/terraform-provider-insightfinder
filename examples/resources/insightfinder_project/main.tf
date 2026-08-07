terraform {
  required_providers {
    insightfinder = {
      source  = "insightfinder/insightfinder"
      version = "~> 1.9"
    }
  }
}

provider "insightfinder" {
  base_url    = "https://app.insightfinder.com"
  username    = var.username
  license_key = var.license_key
}

# Basic log project
resource "insightfinder_project" "basic_logs" {
  project_name = "basic-application-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Application Logs"
  project_time_zone    = "UTC"
  sampling_interval    = 600
  retention_time       = 90
}

# Advanced project with anomaly detection
resource "insightfinder_project" "advanced_logs" {
  project_name = "advanced-system-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "AWS"
    project_cloud_type = "AWS"
    insight_agent_type = "Historical"
  }

  project_display_name = "System Logs with ML"
  project_time_zone    = "America/New_York"
  sampling_interval    = 600
  retention_time       = 180

  # Anomaly detection
  anomaly_detection_mode    = 1
  anomaly_sampling_interval = 600
  enable_hot_event          = true
  hot_event_threshold       = 10
  hot_number_limit          = 20

  # Rare event handling
  rare_event_auto_incident_flag = true

  # Email alerts
  enable_new_alert_email = true
  email_setting = jsonencode({
    enableIncidentDetectionEmailAlert  = true
    enableIncidentPredictionEmailAlert = true
    enableRootCauseEmailAlert          = true
    emailDampeningPeriod               = 3600000
    onlySendWithRCA                    = false
    awSeverityLevel                    = "Major"
  })

  # Webhook configuration
  webhook_url = "https://hooks.example.com/incident"
  webhook_type_set_str = jsonencode([
    "log",
    "detectedIncident",
    "predictedIncident",
    "detectedIncidentWithRC"
  ])

  # Root cause analysis
  root_cause_probability_threshold = 0.8
  root_cause_count_threshold       = 1
  maximum_root_cause_result_size   = 5
}

# ServiceNow project with table name
resource "insightfinder_project" "servicenow_incidents" {
  project_name = "servicenow-incidents"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "ServiceNow"
    project_cloud_type = "ServiceNow"
    servicenow_table   = "incident" # Required for ServiceNow projects
  }

  project_display_name = "ServiceNow Incidents"
  project_time_zone    = "UTC"
  sampling_interval    = 600
  retention_time       = 90

  # ServiceNow third-party settings
  project_servicenow_settings = {
    host                 = "https://dev123456.service-now.com/"
    servicenow_user      = "admin"
    servicenow_password  = var.servicenow_password
    instance_field       = "configuration item"
    instance_field_regex = ""
    timestamp_format     = "yyyy-MM-dd HH:mm:ss"
    sysparm_query        = ""
    proxy                = ""
    additional_fields    = ["work_end", "priority"]
    component_name_rule  = ""
  }
}

# Project with custom JSON key settings
resource "insightfinder_project" "json_logs" {
  project_name = "json-structured-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Structured JSON Logs"
  project_time_zone    = "UTC"
  sampling_interval    = 600
  retention_time       = 90

  # Define JSON key settings for extracting custom fields from logs
  json_key_settings = [
    {
      json_key                = "api"
      type                    = "string"
      summary_setting         = false
      metafield_setting       = false
      dampening_field_setting = false
    },
    {
      json_key                = "api2"
      type                    = "string"
      summary_setting         = true
      metafield_setting       = false
      dampening_field_setting = false
    },
    {
      json_key                = "state"
      type                    = "number"
      summary_setting         = true
      metafield_setting       = true
      dampening_field_setting = false
    },
    {
      json_key                = "status"
      type                    = "string"
      summary_setting         = false
      metafield_setting       = true
      dampening_field_setting = true
    },
    {
      json_key                          = "user"
      type                              = "JSONArray"
      summary_setting                   = false
      metafield_setting                 = false
      dampening_field_setting           = false
      notification_setting              = true
      notification_setting_display_name = "User Field"
    }
  ]
}
# Project with process mode
resource "insightfinder_project" "loki_logs" {
  project_name = "loki-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Loki Logs"
  sampling_interval    = 600
  retention_time       = 90

  mode = 1
}

# Project with log-to-metric (L2M) settings
resource "insightfinder_project" "l2m_log_project" {
  project_name = "l2m-source-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "L2M Source Logs"
  sampling_interval    = 600
  retention_time       = 90

  l2m_settings = [
    {
      metric_project_name = "my-metric-project"
      json_flag           = false
      enable_mapping      = false
      regexs = [
        {
          metric_name_regex     = "metric=(\\w+)"
          metric_value_regex    = "value=([\\d.]+)"
          instance_name_regex   = "host=([\\w.-]+)"
          timestamp_regex       = "ts=([\\d]+)"
          timestamp_format      = "epoch"
          metric_name           = "response_time"
          operation             = 0
          aggregation_mode      = 2
          aggregation_period    = 60
          grouping_by_component = false
        }
      ]
    },
    {
      metric_project_name = "my-json-metric-project"
      json_flag           = true
      enable_mapping      = false
      json_parsers = [
        {
          metric_value_key      = "alert->core->summary"
          instance_name_key     = "alert->asset->asset_id"
          timestamp_key         = "alert->asset->asset_video_format"
          timestamp_format      = "epoch"
          operation             = 1
          aggregation_mode      = 3
          aggregation_period    = 60
          grouping_by_component = true
          derived_value_model = {
            base_value      = "alert->asset->asset_id=v"
            actual_value    = "alert->cloud->availability_zone=b"
            operation       = 2
            mapping_id_list = ["alert->asset->stream_type", "alert->asset->pid"]
          }
        }
      ]
    }
  ]
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

variable "servicenow_password" {
  description = "ServiceNow account password"
  type        = string
  sensitive   = true
}

# Project with incident priority by anomaly score
resource "insightfinder_project" "priority_example" {
  project_name = "priority-configured-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Priority Configured Logs"
  project_time_zone    = "UTC"
  sampling_interval    = 600
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
