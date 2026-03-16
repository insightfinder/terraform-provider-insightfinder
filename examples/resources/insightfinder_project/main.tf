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
      json_key                = "user"
      type                    = "JSONArray"
      summary_setting         = false
      metafield_setting       = false
      dampening_field_setting = false
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
