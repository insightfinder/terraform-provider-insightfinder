---
page_title: "insightfinder_project Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages an InsightFinder project with comprehensive configuration options for log and metric analysis.
---

# insightfinder_project (Resource)

Manages an InsightFinder project. Projects are the primary containers for log or metric data with configurable anomaly detection, alerting, and analysis settings.

## Example Usage

### Basic Log Project

```terraform
resource "insightfinder_project" "app_logs" {
  project_name = "application-logs"
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

### Advanced Project with Alerting

```terraform
resource "insightfinder_project" "advanced" {
  project_name = "critical-services"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "AWS"
    project_cloud_type = "AWS"
    insight_agent_type = "Historical"
  }

  project_display_name      = "Critical Services"
  project_time_zone         = "America/New_York"
  sampling_interval         = 600
  retention_time            = 180
  
  # Anomaly detection
  anomaly_detection_mode    = 1
  anomaly_sampling_interval = 600
  enable_hot_event          = true
  hot_event_threshold       = 10
  
  # Email alerts
  enable_new_alert_email = true
  email_setting = jsonencode({
    enableIncidentDetectionEmailAlert  = true
    enableIncidentPredictionEmailAlert = true
    enableRootCauseEmailAlert          = true
    emailDampeningPeriod               = 3600000
  })
  
  # Webhook
  webhook_url = "https://hooks.example.com/incidents"
  webhook_type_set_str = jsonencode([
    "log",
    "detectedIncident",
    "predictedIncident"
  ])
}
```

### ServiceNow Project with Third-Party Settings

```terraform
resource "insightfinder_project" "servicenow_project" {
  project_name = "servicenow-incidents"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "ServiceNow"
    project_cloud_type = "ServiceNow"
    insight_agent_type = "Custom"
  }

  project_display_name = "ServiceNow Incidents"
  project_time_zone    = "UTC"
  sampling_interval    = 600
  retention_time       = 90

  # ServiceNow third-party settings (only applies when project_cloud_type is ServiceNow)
  project_servicenow_settings = {
    host                 = "https://dev123456.service-now.com/"
    servicenow_user      = "admin"
    servicenow_password  = "your-password"
    client_id            = "your-oauth-client-id"
    client_secret        = "your-oauth-client-secret"
    instance_field       = "short_description"
    instance_field_regex = "1"
    timestamp_format     = "yyyy-MM-dd HH:mm:ss"
    sysparm_query        = ""
    proxy                = ""
    additional_fields    = ["work_end", "priority"]
  }
}
```

## Schema

### Required

- `project_name` (String) Unique project identifier
- `system_name` (String) Name of the system this project belongs to
- `project_creation_config` (Object) Project creation configuration
  - `data_type` (String) Type of data: `Log`, `Metric`, or `Alert`
  - `instance_type` (String) Instance type: `AWS`, `Azure`, `GCP`, `PrivateCloud`, `OnPremise`
  - `project_cloud_type` (String) Cloud type (usually same as instance_type)
  - `insight_agent_type` (String) Agent type: `LogStreaming`, `MetricFile`, `Historical`

### Optional

- `project_display_name` (String) Display name for the project
- `project_time_zone` (String) Time zone (e.g., `UTC`, `America/New_York`)
- `sampling_interval` (Number) Data sampling interval in seconds. Default: `600`
- `retention_time` (Number) Data retention period in days. Default: `90`
- `anomaly_detection_mode` (Number) Anomaly detection mode. Default: `0`
- `enable_hot_event` (Boolean) Enable hot event detection. Default: `true`
- `enable_new_alert_email` (Boolean) Enable email alerts. Default: `false`
- `email_setting` (String) JSON-encoded email configuration
- `webhook_url` (String) Webhook URL for notifications
- `webhook_type_set_str` (String) JSON array of webhook event types
- `log_label_settings` (List of Objects) List of log label settings for the project
  - `label_type` (String) Type of log label (whitelist, blacklist, patternName, etc.)
  - `log_label_string` (String) JSON-encoded log label configuration
- `project_servicenow_settings` (Object) ServiceNow third-party settings. **Note:** Only applies when `project_cloud_type` is "ServiceNow" (case insensitive). Settings will be ignored if project_cloud_type is not ServiceNow.
  - `host` (String, Required) ServiceNow instance host URL (e.g., `https://dev123456.service-now.com/`)
  - `servicenow_user` (String, Required) ServiceNow username for authentication
  - `servicenow_password` (String, Required, Sensitive) ServiceNow password for authentication
  - `client_id` (String, Optional) OAuth client ID for ServiceNow
  - `client_secret` (String, Optional, Sensitive) OAuth client secret for ServiceNow
  - `instance_field` (String, Optional) Field to use for instance identification (default: `short_description`)
  - `instance_field_regex` (String, Optional) Regex pattern for instance field
  - `timestamp_format` (String, Optional) Timestamp format for ServiceNow data (default: `yyyy-MM-dd HH:mm:ss`)
  - `sysparm_query` (String, Optional) ServiceNow query parameter
  - `proxy` (String, Optional) Proxy URL for ServiceNow connection
  - `additional_fields` (List of String, Optional) Additional fields to fetch from ServiceNow

See full schema in the [complete example](https://github.com/insightfinder/terraform-provider-insightfinder/tree/main/examples/resources/insightfinder_project).

### Read-Only

- `id` (String) Project identifier (same as project_name)

## Import

Projects can be imported using the project name:

```shell
terraform import insightfinder_project.example my-project-name
```
