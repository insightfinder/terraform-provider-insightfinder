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
    awSeverityLevel                    = "Major"
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
    servicenow_table   = "incident"  # Required for ServiceNow projects
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
    sysparm_query           = ""
    proxy                   = ""
    additional_fields       = ["work_end", "priority"]
    component_name_rule     = ""
    service_now_import_flag = true
  }
}
```

### Project with Holiday Settings

```terraform
resource "insightfinder_project" "holidays_example" {
  project_name = "project-with-holidays"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
  }

  project_display_name = "Project with Holiday Settings"
  sampling_interval    = 600

  # Define holidays that affect anomaly detection
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
```

### Project with JSON Key Settings

```terraform
resource "insightfinder_project" "json_logs_example" {
  project_name = "json-structured-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Structured JSON Logs"
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
```
### Project with Mode

```terraform
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

  mode = 4
}
```

## Schema

### Required

- `project_name` (String) Unique project identifier
- `system_name` (String) Name of the system this project belongs to
- `project_creation_config` (Object) Project creation configuration
  - `data_type` (String) Type of data: `Log`, `Metric`, or `Alert`
  - `instance_type` (String) Instance type: `AWS`, `Azure`, `GCP`, `PrivateCloud`, `OnPremise`, `ServiceNow`
  - `project_cloud_type` (String) Cloud type (usually same as instance_type)
  - `insight_agent_type` (String) Agent type: `LogStreaming`, `MetricFile`, `Historical`, `Custom`
  - `servicenow_table` (String) ServiceNow table name (required when `project_cloud_type` is `ServiceNow`)

### Optional

- `project_display_name` (String) Display name for the project
- `project_time_zone` (String) Time zone (e.g., `UTC`, `America/New_York`)
- `sampling_interval` (Number) Data sampling interval in seconds. Default: `600`
- `retention_time` (Number) Data retention period in days. Default: `90`
- `anomaly_detection_mode` (Number) Anomaly detection mode. Default: `0`
- `component_name_auto_overwrite` (Boolean) Enable automatic overwrite of component names
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
  - `component_name_rule` (String, Optional) Rule for determining the component name from ServiceNow data
  - `service_now_import_flag` (Boolean, Optional, Computed) Whether to enable importing data from ServiceNow
- `holiday_settings` (List of Objects) List of holiday settings for the project. Each holiday defines a period that should be treated as a holiday for anomaly detection purposes.
  - `name` (String, Required) Name of the holiday (must be unique within the project)
  - `start_date` (String, Required) Start date of the holiday in MM-DD format (e.g., `12-25`)
  - `end_date` (String, Required) End date of the holiday in MM-DD format (e.g., `12-26`)
- `json_key_settings` (List of Objects) List of JSON key settings for extracting custom fields from JSON-structured logs. Manages which JSON keys are available for analysis and which should be included in summary, metafield, and dampening field statistics.
  - `json_key` (String, Required) The JSON key name to extract from logs
  - `type` (String, Required) The data type of the JSON value (e.g., `string`, `number`, `JSONArray`)
  - `summary_setting` (Boolean, Required) Whether to include this key in the summary statistics. When `true`, the key's values will be aggregated in summary reports.
  - `metafield_setting` (Boolean, Required) Whether to include this key in the metafield statistics. When `true`, the key's values will be tracked as metafield data for enhanced log analysis.
  - `dampening_field_setting` (Boolean, Required) Whether to include this key in the dampening field list. When `true`, the key is used to control alert dampening logic — alerts with the same value for this field will be grouped and suppressed during the dampening window.
- `mode` (Number, Optional, Computed) Process mode for the project. Controls which RabbitMQ processing queue the project's data is routed to. Set and read via the `/api/v1/logdedicatedmode` API. Maps to the `processMode` field in the API response. Default: `0` (LIVE).

  | Value | Name | Queue Suffix | Description |
  |-------|------|--------------|-------------|
  | `0` | `LIVE` | *(default queue)* | Default. Standard real-time processing queue |
  | `1` | `HISTORICAL` | `-historical` | Historical/backfill data processing — use when ingesting past data to avoid competing with live traffic |
  | `2` | `UPDATE` | `-update` | Re-processing of existing data |
  | `3` | `AW` | `-aw` | AI Watchtower dedicated queue |
  | `4` | `DEDICATED` | `-dedicated` | Isolated worker queue — use for log-heavy projects to prevent starving other projects on the shared queue |

See full schema in the [complete example](https://github.com/insightfinder/terraform-provider-insightfinder/tree/main/examples/resources/insightfinder_project).

### Read-Only

- `id` (String) Project identifier (same as project_name)

## Import

Projects can be imported using the project name:

```shell
terraform import insightfinder_project.example my-project-name
```
