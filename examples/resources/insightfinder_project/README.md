# Project Resource Examples

This directory contains examples for using the `insightfinder_project` resource.

## Basic Log Project

```hcl
resource "insightfinder_project" "application_logs" {
  project_name = "my-app-logs"
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
```

## Metric Project with Anomaly Detection

```hcl
resource "insightfinder_project" "infrastructure_metrics" {
  project_name = "infra-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "AWS"
    project_cloud_type = "AWS"
    insight_agent_type = "MetricFile"
  }

  project_display_name = "Infrastructure Metrics"
  project_time_zone    = "America/New_York"
  sampling_interval    = 300
  
  # Anomaly detection settings
  anomaly_detection_mode    = 1
  anomaly_sampling_interval = 600
  enable_hot_event          = true
  hot_event_threshold       = 10
  
  # Alert settings
  enable_new_alert_email = true
  email_setting = jsonencode({
    enableIncidentDetectionEmailAlert = true
    enableIncidentPredictionEmailAlert = true
    enableRootCauseEmailAlert = true
    emailDampeningPeriod = 3600000
  })
}
```

## Complete Project with All Settings

See `complete-project.tf` for a comprehensive example.

## Usage

1. Copy the example to your Terraform configuration
2. Modify the values to match your requirements
3. Run:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```
