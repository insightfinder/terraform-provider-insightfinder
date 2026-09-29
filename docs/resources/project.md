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

  # Arbitrary customer-specific settings
  custom_setting = jsonencode({
    nbc = {
      global_commerce = {
        flow = "enabled"
      }
    }
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
    },
    {
      json_key                                      = "service"
      type                                          = "string"
      summary_setting                               = true
      metafield_setting                             = false
      dampening_field_setting                       = false
      notification_setting                          = true
      notification_setting_display_name             = "Service"
      service_now_notification_setting              = true
      service_now_notification_setting_display_name = "Service (ServiceNow)"
    }
  ]

  # ServiceNow notification format strings
  service_now_short_description_format = "Alert: {service} anomaly detected"
  service_now_description_format       = "Detailed description: {service} reported anomaly at {timestamp}"

  # Slack notification template
  slack_block_template = "${suggestedPriority}, {service} anomaly detected"
}
```
### Project with Log-to-Metric Settings

```terraform
resource "insightfinder_project" "l2m_example" {
  project_name = "my-log-project"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "My Log Project"
  sampling_interval    = 600

  # Regex-based log-to-metric conversion
  l2m_settings = [
    {
      metric_project_name = "my-metric-project"
      json_flag           = false
      enable_mapping      = false
      regexs = [
        {
          metric_name_regex    = "metric_name=(\\w+)"
          metric_value_regex   = "value=([\\d.]+)"
          instance_name_regex  = "host=([\\w.-]+)"
          timestamp_regex      = "ts=([\\d]+)"
          timestamp_format     = "epoch"
          metric_name          = "response_time"
          operation            = 0
          aggregation_mode     = 2
          aggregation_period   = 60
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
          metric_value_key      = ""
          data_filter           = "alert->error->message=.*(?i)EFS.*"
          operation             = 1
          metric_name           = "EFS Error"
          aggregation_mode      = 0
          aggregation_period    = 0
          grouping_by_component = false
        },
        {
          metric_value_key     = "alert->core->value"
          instance_name_key    = "alert->asset->asset_id"
          timestamp_key        = "alert->timestamp"
          timestamp_format     = "epoch"
          operation            = 1
          aggregation_mode     = 3
          aggregation_period   = 60
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

  mode = 1
}
```

## Schema

-> **Note on defaults:** The provider does not define schema defaults. Unless noted otherwise, every optional attribute is also *Computed*: when it is omitted from the configuration, its value is read back from the InsightFinder API. The "Default" values listed below are the InsightFinder server-side defaults.

### Required

- `project_name` (String) Unique project identifier. Changing this forces a new resource to be created.
- `system_name` (String) Name of the system this project belongs to.
- `project_creation_config` (Attributes) Project creation configuration. See [below for nested schema](#nestedatt--project_creation_config).

### Optional

#### General Settings

- `project_display_name` (String) Display name for the project.
- `project_time_zone` (String) Time zone for the project (e.g., `UTC`, `America/New_York`). Default: `UTC`
- `sampling_interval` (Number) Data sampling interval in seconds. Default: `600`
- `c_value` (Number) The C value for anomaly detection sensitivity (typically 2-5).
- `p_value` (Number) The P value for anomaly detection probability (0.0-1.0).
- `retention_time` (Number) Data retention period in days. Default: `90`
- `ubl_retention_time` (Number) Retention time for UBL data in days. Default: `90`
- `mode` (Number) Process mode for the project. Controls which RabbitMQ processing queue the project's data is routed to. Set and read via the `/api/v1/logdedicatedmode` API. Maps to the `processMode` field in the API response. Default: `0` (LIVE).

  | Value | Name | Queue Suffix | Description |
  |-------|------|--------------|-------------|
  | `0` | `LIVE` | *(default queue)* | Default. Standard real-time processing queue |
  | `1` | `HISTORICAL` | `-historical` | Historical/backfill data processing — use when ingesting past data to avoid competing with live traffic |
  | `2` | `UPDATE` | `-update` | Re-processing of existing data |
  | `3` | `AW` | `-aw` | AI Watchtower dedicated queue |
  | `4` | `DEDICATED` | `-dedicated` | Isolated worker queue — use for log-heavy projects to prevent starving other projects on the shared queue |

- `shared_usernames` (String) JSON-encoded array of usernames to share the project with (e.g., `jsonencode(["user1", "user2"])`). Must be a JSON array; other values are ignored.
- `custom_setting` (String) Arbitrary customer-specific settings (JSON object). Accepts any nested structure of string values (e.g., `jsonencode({ nbc = { global_commerce = { flow = "enabled" } } })`).

#### Anomaly Detection & Alerts

- `anomaly_detection_mode` (Number) Enables or disables anomaly detection for log data. `0`: enabled (default); `-1`: disabled.
- `anomaly_sampling_interval` (Number) The time window (in seconds) used for log anomaly detection. Default: `60`
- `enable_anomaly_score_escalation` (Boolean) *(Metric projects only)* Automatically escalate incidents whose anomaly score is at or above `escalation_anomaly_score_threshold`. Default: `false`
- `escalation_anomaly_score_threshold` (String) *(Metric projects only)* If the calculated anomaly score is greater than or equal to this threshold, the anomaly is flagged for escalation. Numeric value passed as a string, typically between `0.0` and `1.0` (e.g., `"1.0"` escalates only the most extreme anomalies, `"0.1"` escalates almost any deviation).
- `ignore_anomaly_score_threshold` (String) *(Metric projects only)* Anomalies with a score below this value are ignored. Numeric value passed as a string, typically between `0.0` and `1.0` (e.g., `"0.3"` discards any anomaly scoring less than 0.3). Default: `0.0`
- `enable_stream_detection` (Boolean) *(Metric projects only)* Switch the metric detection engine from batch-processing mode to streaming detection mode. Default: `false`
- `enable_hot_event` (Boolean) *(Log projects only)* Enable hot event detection (frequency-based spikes of known log patterns). Default: `true`
- `hot_event_threshold` (Number) *(Log projects only)* Maximum allowable count of a specific log pattern within a sampling interval before it is classified as an anomaly. Must be `>= 0`; a very high value effectively suppresses hot event alerts. Default: `50`
- `hot_event_calm_down_period` (Number) *(Log projects only)* After a pattern triggers a hot event, subsequent spikes of the same pattern within this period are recorded but do not generate new alerts or incidents. Unit: sampling intervals (multiples of `anomaly_sampling_interval`). Must be `> 0`; values `<= 0` fall back to the default. Default: `3`
- `hot_event_detection_mode` (Number) *(Log projects only)* How the system decides whether a log volume is "hot". `0`: standard detection logic (default); `1`: advanced statistical approach relative to moving averages.
- `hot_number_limit` (Number) *(Log projects only)* Maximum number of unique log patterns that can be classified as "hot" within a single processing cycle. Default: `20`
- `cold_event_threshold` (Number) *(Log projects only)* Sensitivity threshold for triggering cold event alerts. Must be `>= 0`; `0` effectively disables cold event detection. Default: `10`
- `cold_number_limit` (Number) *(Log projects only)* Maximum number of cold events detected per day. Default: `0`
- `rare_anomaly_type` (Number) *(Log projects only)* How rare (infrequent or new) log patterns are categorized. `0`: detect both new patterns and known patterns with very low frequency (default); `1`: alert only when a log template is seen for the first time; `2`: alert only when an existing template contains a new/rare value.
- `rare_event_alert_thresholds` (Number) *(Log projects only)* Cluster size / frequency a rare event must reach before triggering an alert. Must be `>= 0`; higher values are less sensitive, `1` alerts on the very first occurrence. Default: `1`
- `rare_number_limit` (Number) *(Log projects only)* Maximum number of unique log patterns that can be classified as "rare" at once. Default: `20`
- `rare_event_auto_incident_flag` (Boolean) Automatically create an incident for detected rare events.
- `collect_all_rare_events_flag` (Boolean) *(Log projects only)* When `true`, collect and analyze every rare event; when `false`, use standard (potentially sampled) collection. Default: `false`
- `new_alert_flag` (Boolean) *(Log projects only)* When `true`, the system tracks and flags the "newness" of log alerts (new vs. recurring anomalies). Default: `false`
- `enable_new_alert_email` (Boolean) Send email alerts for newly detected anomalies and incidents. Default: `false`
- `alert_average_time` (Number) *(Metric projects only)* Duration over which average metric values are calculated when evaluating alert conditions (typically minutes). Default: `0`
- `alert_hourly_cost` (Number) Monetary value (e.g., USD per hour) of downtime or degraded performance for this project. Must be `>= 0.0`. Default: `0.0`
- `avg_per_incident_downtime_cost` (Number) Flat monetary value (e.g., USD) assigned to every incident in the project, used to quantify business impact. Default: `0.0`
- `incident_priority_by_anomaly_score_setting` (String) Configures how incidents are assigned a priority based on their anomaly score. Accepts a JSON-encoded object with two fields:
  - `enabled` (Boolean) — whether priority assignment is active.
  - `priorityScoreRangeMap` (Object) — maps priority levels (`"1"` through `"5"`) to score range strings. The format is `"<lower>-<upper>"` where the upper bound may be omitted for open-ended ranges (e.g. `"10001-"`).

  Example:
  ```hcl
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
  ```
- `incident_priority_cap_setting` (String) Incident priority cap settings (JSON). Contains `ticketCreationPriorityCap` and `suggestedPriorityCap` string values.

  Example:
  ```hcl
  incident_priority_cap_setting = jsonencode({
    ticketCreationPriorityCap = "2"
    suggestedPriorityCap      = "5"
  })
  ```

#### Log Settings

- `log_detection_min_count` (Number) Minimum count for triggering log detection in a batch. Detection is triggered immediately once at least this many log entries are available in a batch; otherwise it waits until 3 minutes later. Default: `10000`
- `log_detection_size` (Number) Maximum count for triggering log detection in a batch. Default: `30000`
- `log_pattern_limit_level` (Number) Limit on the number of distinct log patterns. Patterns generated beyond this number are all assigned to the MISC (`-2`) miscellaneous pattern. Default: `1024`
- `max_log_model_size` (Number) Maximum number of training samples (logs) per model. Default: `10000`
- `keyword_feature_number` (Number) Number of keyword features (feature vector length) generated for training the model. Default: `200`
- `keyword_setting` (Number) Keyword setting used when collecting and querying keywords. `-1`: disabled (default); `0`: letters only; `1`: letters and numbers.
- `model_keyword_setting` (Number) Keyword selection during model training. `0`: letters only (default); `1`: letters and numbers.
- `model_keyword_segment_k` (Number) Number of keyword segments (K) used for the log model.
- `model_match_threshold` (Number) Model match threshold for log pattern matching.
- `disable_model_keyword_stats_collection` (Boolean) Disable collection of keyword frequencies for model training. Default: `false`
- `disable_log_compress_event` (Boolean) Disable saving (compressing) the log event data. Default: `false`
- `disable_log_processing_flag` (Boolean) Disable log processing.
- `log_anomaly_event_base_score` (String) Base score weights for the different log anomaly event types, as a JSON array string in the order rare, hot, cold, detection alert, new pattern, critical. Default: `"[5,0.01,0.0075,0.01,1,100]"`
- `multi_line_flag` (Boolean) Enable multi-line (regex multiline) processing. Default: `false`
- `nlp_flag` (Boolean) Enable NLP. Deprecated. Default: `false`
- `pretty_json_convertor_flag` (Boolean) When enabled, the system reformats invalid JSON log data into valid JSON. Default: `false`
- `similarity_sensitivity` (String) *(Log projects only)* Strictness of the log clustering engine, i.e. how similar two log messages must be to be grouped under the same pattern. `high`: very strict; `medium`: balanced (default); `low`: loose.
- `feature_outlier_sensitivity` (String) *(Log projects only)* How aggressively numerical features extracted from log messages are flagged as outliers. `high`: flags minor deviations; `medium`: balanced (default); `low`: only extreme outliers.
- `feature_outlier_threshold` (Number) *(Log projects only)* Fixed cutoff for numerical outlier detection in log-extracted data. `0.0` relies entirely on `feature_outlier_sensitivity`; a positive value acts as a hard limit. Default: `0.0`
- `new_pattern_number_limit` (Number) *(Log projects only)* Cap on the number of new (previously unseen) log patterns identified and tracked during a single processing interval. Unset means no limit.
- `new_pattern_range` (Number) *(Log projects only)* Suppression ("calm down") window for alerts triggered by new log patterns. Unit: sampling intervals. Must be `>= 0`. Default: typically `3`
- `whitelist_number_limit` (Number) Maximum number of whitelisted ("known safe") log patterns the project can track. Default: `100`
- `project_model_flag` (Boolean) *(Log projects only)* When `true`, log patterns are shared across all instances in the project (project model); when `false`, pattern learning is isolated to each instance (instance model). Default: typically `true`
- `zone_name_key` (String) Name of the log field (key) that contains the zone information (e.g., cloud region or data center), as it appears in the logs.
- `component_name_auto_overwrite` (Boolean) Enable automatic overwrite of component names.

#### Incident Prediction & Root Cause Analysis

- `incident_prediction_window` (Number) "Look-ahead" time for the predictive analytics engine, in minutes. Must be `>= 0`. Default: `0`
- `min_incident_prediction_window` (Number) Lower bound on the predictive alerting time, in minutes. Default: `0`
- `incident_prediction_event_limit` (Number) Maximum number of predicted incidents tracked, displayed, or alerted on within a single processing window. `0` suppresses all predicted incident displays. Default: `5`
- `incident_relation_search_window` (Number) Time window (in minutes) used to link predicted incidents to actual detected events. Default: typically `60`
- `root_cause_count_threshold` (Number) Maximum number of root cause candidates identified and presented for a single anomaly or incident. `0` hides all root cause suggestions. Default: typically `10`
- `root_cause_probability_threshold` (Number) Probability threshold (0.0-1.0) used to filter root cause candidates. Default: `0.8`
- `root_cause_log_message_search_range` (Number) Time window (in minutes) around an anomaly or incident in which related log messages are searched during RCA. Default: typically `60`
- `root_cause_rank_setting` (Number) Ranking algorithm used to prioritize identified root causes. `0`: default balanced ranking; `1+`: alternative ranking modes. Default: `0`
- `maximum_root_cause_result_size` (Number) Hard limit on the number of root cause entries processed and presented. `0` uses the built-in global default. Default: typically `20`
- `multi_hop_search_level` (Number) Depth (number of hops) traversed in the causal dependency graph when searching for a root cause. `0` disables deep causal searching. Default: typically `1`
- `multi_hop_search_limit` (String) Breadth of the causal search: maximum number of neighbor nodes / candidates explored at each hop. Numeric value passed as a string. Default: typically `"10"`
- `causal_prediction_setting` (Number) Scope of causal analysis used for incident prediction. `0`: within the project and across other accessible projects (default); `1`: within the current project only; `2`: only relationships where the cause originates from a different project.
- `causal_min_delay` (String) Minimum time difference (in minutes) required between a cause event and an effect event to be considered a valid causal pair. Numeric value passed as a string. Default: typically `"0"` (no minimum)
- `normal_event_causal_flag` (Boolean) *(Log projects only)* When `true`, normal (non-anomalous) log events are also included as candidates in causal analysis / RCA. Default: `false`
- `training_filter` (Boolean) When `true`, suppress incidents that occur outside of a known training window. Default: typically `false`
- `daily_model_span` (Number) Number of days of historical data used to build the daily behavioral model. Must be `>= 1`. Default: `1`
- `min_valid_model_span` (Number) Minimum duration (in milliseconds) of data a model must contain to be considered valid. Default: typically `21600000` (6 hours)
- `maximum_detection_wait_time` (Number) *(Log projects only)* Maximum time (in minutes) log anomaly detection waits for late-arriving logs before analyzing a time window. Default: typically `15`
- `maximum_threads` (Number) *(Log projects only)* Degree of parallelism for heavy background tasks (log training, detection, replay). Minimum `1` (sequential). Default: `1`
- `large_project` (Boolean) Optimize processing for projects with a very large number of instances or very high data throughput. Default: `false`

#### Prediction Rules

- `prediction_count_threshold` (Number) Minimum amount of evidence required to trigger a predicted incident alert. `1` means any single piece of evidence is enough. Default: `1`
- `prediction_probability_threshold` (Number) Confidence filter (0.0-1.0) for predicted incidents; e.g., `0.9` only alerts on incidents with at least 90% predicted likelihood. Default: `0.8`
- `prediction_rule_active_condition` (Number) Maturity required of a learned causal rule before it is used for prediction. `0`: use all discovered rules immediately; `1`: only rules that have successfully predicted an event at least once; `2+`: require multiple historical verifications. Default: typically `1`
- `prediction_rule_active_threshold` (Number) Minimum probability (0.0-1.0) required for a causal rule to be promoted to "active". Must be greater than or equal to `prediction_rule_inactive_threshold`. Default: typically `0.7`
- `prediction_rule_inactive_threshold` (Number) Confidence score (0.0-1.0) below which an active causal rule is demoted. Must be less than or equal to `prediction_rule_active_threshold`. Default: typically `0.5`
- `prediction_rule_false_positive_threshold` (Number) Maximum allowable false-positive rate for a causal rule before it is disabled for proactive alerting. **Note:** the provider models this attribute as an integer.

#### Instance Settings

- `instance_convert_flag` (Boolean) *(Log projects only)* When `true`, apply conversion/normalization to instance names; when `false`, keep instance names exactly as they appear in the raw log stream. Default: `false`
- `instance_down_enable` (Boolean) Enable "instance down" detection (alert when an instance stops reporting data). Default: `false`
- `show_instance_down` (Boolean) Whether to show instance down events/incidents in dashboards, timelines, and incident reports. Default: `true`
- `is_grouping_by_instance` (Boolean) *(Log projects only)* When `true`, analyze each instance independently; when `false`, aggregate analysis across the entire project. Default: typically `true`
- `ignore_instance_for_kb` (Boolean) When `true`, ignore the instance name when matching against the Knowledge Base (match by pattern only). Default: `false`
- `is_edge_brain` (Boolean) When `true`, the project is treated as an "Edge" (resource-constrained) instance instead of a standard "Brain/Cloud" instance. Default: `false`
- `is_trace_prompt` (Boolean) *(Log projects only)* When `true`, the project is treated as an LLM trace/prompt evaluation project, enabling LLM-specific evaluation and observability features. Default: `false`
- `instance_grouping_update` (String) JSON-encoded instance grouping settings, e.g. `jsonencode({ autoFill = true })` to enable auto-fill for instance grouping.

#### Webhook & Email Notifications

- `webhook_url` (String) Webhook URL for notifications.
- `webhook_type_set_str` (String) JSON array string of webhook event types (e.g., `jsonencode(["log", "detectedIncident", "predictedIncident"])`).
- `webhook_black_list_set_str` (String) Blacklist set string for webhooks.
- `webhook_critical_keyword_set_str` (String) Critical keyword set string for webhooks.
- `webhook_alert_dampening` (Number) Alert dampening for webhooks.
- `max_web_hook_request_size` (Number) Maximum webhook request size.
- `webhook_header_list` (String) JSON-encoded array of header objects to include in webhook requests. Must be a JSON array; other values are ignored.

  ```hcl
  webhook_header_list = jsonencode([
    { headerName = "Authorization", headerValue = "Bearer token" }
  ])
  ```
- `email_setting` (String) JSON-encoded email notification configuration. Supported keys:
  - `enableIncidentDetectionEmailAlert` (Boolean) Enable email alerts for incident detection
  - `enableIncidentPredictionEmailAlert` (Boolean) Enable email alerts for incident prediction
  - `enableRootCauseEmailAlert` (Boolean) Enable email alerts including root cause analysis
  - `enableAlertsEmail` (Boolean) Enable general alerts email
  - `enableNotificationAW` (Boolean) Enable notification for AI Watchtower
  - `onlySendWithRCA` (Boolean) Only send alerts if RCA is available
  - `emailDampeningPeriod` (Number) Dampening period in milliseconds
  - `alertsEmailDampeningPeriod` (Number) Dampening period for alerts in milliseconds
  - `predictionEmailDampeningPeriod` (Number) Dampening period for prediction alerts in milliseconds
  - `awSeverityLevel` (String) Severity level for AI Watchtower notifications
- `service_now_short_description_format` (String, Optional) Short description format string for ServiceNow notifications. Maps to `serviceNowNotificationAdditionalSetting.shortDescriptionFormat`. Only applied when `json_key_settings` is configured.
- `service_now_description_format` (String, Optional) Description format string for ServiceNow notifications (non-key-value notification content). Maps to `serviceNowNotificationAdditionalSetting.descriptionFormat`. Only applied when `json_key_settings` is configured.
- `slack_block_template` (String, Optional) Slack block template string used for Slack notifications. Maps to `slackNotificationAdditionalSetting.slackBlockTemplate`. Only applied when `json_key_settings` is configured.

#### Other JSON Configuration Strings

These attributes accept JSON-encoded strings (use `jsonencode(...)`).

- `proxy` (String) Proxy server used when communicating with external systems or running actions (e.g., `http://proxy.internal:8080`). Default: empty.
- `base_value_setting` (String) JSON-encoded base value and metric mapping configuration. Supported keys: `isSourceProject` (Boolean), `mappingKeys` (List), `baseValueKeys` (List), `metricProjects` (List of String), `additionalMetricNames` (List of String).
- `cdf_setting` (String) JSON-encoded array of CDF setting objects. Must be a JSON array; other values are ignored.
- `llm_evaluation_setting` (String) JSON-encoded LLM evaluation settings. Supported Boolean keys: `isHallucinationEvaluation`, `isAnswerRelevantEvaluation`, `isLogicConsistencyEvaluation`, `isFactualInaccuracyEvaluation`, `isMaliciousPromptEvaluation`, `isToxicityEvaluation`, `isPiiPhiLeakageEvaluation`, `isTopicGuardrailsEvaluation`, `isToneDetectionEvaluation`, `isAnomalousOutliersEvaluation`, `showSafetyTemplate`, and the bias evaluations `isGenderBiasEvaluation`, `isRacialBiasEvaluation`, `isSocioeconomicBiasEvaluation`, `isCulturalBiasEvaluation`, `isReligiousBiasEvaluation`, `isPoliticalBiasEvaluation`, `isDisabilityBiasEvaluation`, `isAgeBiasEvaluation`.
- `log_to_log_setting_list` (String) JSON-encoded array of log-to-log setting objects. Must be a JSON array; other values are ignored.

#### Nested Settings

- `log_label_settings` (Set of Objects) Set of log label settings for the project. Each setting is applied individually via API. See [below for nested schema](#nestedatt--log_label_settings).
- `project_servicenow_settings` (Attributes) ServiceNow third-party settings. **Note:** Only applies when `project_creation_config.project_cloud_type` is `ServiceNow`; settings are ignored otherwise. See [below for nested schema](#nestedatt--project_servicenow_settings).
- `holiday_settings` (Set of Objects) Set of holiday settings for the project. Each holiday defines a period that should be treated as a holiday for anomaly detection purposes. See [below for nested schema](#nestedatt--holiday_settings).
- `json_key_settings` (Set of Objects) Set of JSON key settings for extracting custom fields from JSON-structured logs. Manages which JSON keys are available for analysis and which should be included in summary, metafield, dampening field, and notification statistics. See [below for nested schema](#nestedatt--json_key_settings).
- `l2m_settings` (Set of Objects) Log-to-metric settings. Each entry defines how log data from this project is parsed and converted into metrics for a target metric project. Stored as a set so order does not matter. See [below for nested schema](#nestedatt--l2m_settings).

### Read-Only

- `id` (String) Project identifier (same as `project_name`).

<a id="nestedatt--project_creation_config"></a>
### Nested Schema for `project_creation_config`

Required:

- `data_type` (String) Type of data (e.g., `Log`, `Metric`, `Trace`).
- `instance_type` (String) Instance type (e.g., `PrivateCloud`, `AWS`, `Azure`, `ServiceNow`).
- `project_cloud_type` (String) Cloud type for the project (usually the same as `instance_type`, e.g. `PrivateCloud`).

Optional:

- `insight_agent_type` (String) InsightFinder agent type (e.g., `Custom`, `LogStreaming`, `MetricFile`, `Historical`).
- `servicenow_table` (String) ServiceNow table name. Required when `project_cloud_type` is `ServiceNow`.

<a id="nestedatt--log_label_settings"></a>
### Nested Schema for `log_label_settings`

Required:

- `label_type` (String) Type of log label (e.g., `whitelist`, `blacklist`, `patternName`, `logSeverity`, `logEventID`, `logSession`, `logComponent`, `logTransactionID`, `logCustomParameter`).
- `log_label_string` (String) JSON-encoded log label value/pattern (e.g., `jsonencode(["ERROR", "WARN"])`, or `key=["ERROR","WARN"]` for JSON-structured logs).

<a id="nestedatt--project_servicenow_settings"></a>
### Nested Schema for `project_servicenow_settings`

Required:

- `host` (String) ServiceNow instance host URL (e.g., `https://dev123456.service-now.com/`).
- `servicenow_user` (String) ServiceNow username for authentication.
- `servicenow_password` (String, Sensitive) ServiceNow password for authentication.

Optional:

- `client_id` (String) OAuth client ID for ServiceNow.
- `client_secret` (String, Sensitive) OAuth client secret for ServiceNow.
- `instance_field` (String) Field in the ServiceNow record that contains the instance name (e.g., `short_description`).
- `instance_field_regex` (String) Regex applied to `instance_field` to extract the instance name.
- `timestamp_format` (String) Java SimpleDateFormat used to parse ServiceNow timestamps (e.g., `yyyy-MM-dd HH:mm:ss`).
- `sysparm_query` (String) ServiceNow filter query used to limit the fetched records. Default: empty.
- `proxy` (String) Proxy URL for the ServiceNow connection. Default: empty.
- `additional_fields` (List of String) Additional fields to fetch from ServiceNow records.
- `component_name_rule` (String) Rule for determining the component name from ServiceNow data.
- `service_now_import_flag` (Boolean) Whether to enable importing data from ServiceNow.

<a id="nestedatt--holiday_settings"></a>
### Nested Schema for `holiday_settings`

Required:

- `name` (String) Name of the holiday (must be unique within the project).
- `start_date` (String) Start date of the holiday in MM-DD format (e.g., `12-25`).
- `end_date` (String) End date of the holiday in MM-DD format (e.g., `12-26`).

<a id="nestedatt--json_key_settings"></a>
### Nested Schema for `json_key_settings`

Required:

- `json_key` (String) The JSON key name to extract from logs.
- `type` (String) The data type of the JSON value (e.g., `string`, `number`, `JSONArray`).
- `summary_setting` (Boolean) Whether to include this key in the summary statistics. When `true`, the key's values will be aggregated in summary reports.
- `metafield_setting` (Boolean) Whether to include this key in the metafield statistics. When `true`, the key's values will be tracked as metafield data for enhanced log analysis.
- `dampening_field_setting` (Boolean) Whether to include this key in the dampening field list. When `true`, the key is used to control alert dampening logic — alerts with the same value for this field will be grouped and suppressed during the dampening window.

Optional (not computed):

- `notification_setting` (Boolean) Whether to include this key in the notification settings. When `true`, the key is sent in the `notificationSetting` map with `selected: true`.
- `notification_setting_display_name` (String) Display name for this key in notification settings. Defaults to the `json_key` value when not specified.
- `service_now_notification_setting` (Boolean) Whether to include this key in the ServiceNow notification settings. When `true`, the key is sent in the `serviceNowNotificationSetting` map with `selected: true`.
- `service_now_notification_setting_display_name` (String) Display name for this key in ServiceNow notification settings. Defaults to the `json_key` value when not specified.

<a id="nestedatt--l2m_settings"></a>
### Nested Schema for `l2m_settings`

Required:

- `metric_project_name` (String) Name of the target metric project that receives the converted metrics.

Optional:

- `json_flag` (Boolean) When `true`, use `json_parsers`; when `false`, use `regexs`.
- `enable_mapping` (Boolean) Whether to enable mapping.
- `regexs` (List of Objects) Regex-based parser entries. Used when `json_flag` is `false`. See [below for nested schema](#nestedatt--l2m_settings--regexs).
- `json_parsers` (List of Objects) JSON-based parser entries. Used when `json_flag` is `true`. See [below for nested schema](#nestedatt--l2m_settings--json_parsers).

<a id="nestedatt--l2m_settings--regexs"></a>
### Nested Schema for `l2m_settings.regexs`

Optional:

- `metric_name_regex` (String) Regex to extract the metric name from log lines.
- `metric_value_regex` (String) Regex to extract the metric value from log lines.
- `base_value_key` (String) Base value key for derived calculations.
- `instance_name_regex` (String) Regex to extract the instance name.
- `container_name_regex` (String) Regex to extract the container name.
- `timestamp_regex` (String) Regex to extract the timestamp.
- `timestamp_format` (String) Format string for the extracted timestamp.
- `data_filter` (String) Filter expression for log line selection.
- `metric_name` (String) Static metric name.
- `operation` (Number) Parser operation type.
- `aggregation_mode` (Number) Aggregation mode.
- `aggregation_period` (Number) Aggregation period.
- `container_type` (Number) Container type identifier.
- `grouping_by_component` (Boolean) Whether to group metrics by component.

<a id="nestedatt--l2m_settings--json_parsers"></a>
### Nested Schema for `l2m_settings.json_parsers`

Optional:

- `metric_value_key` (String) JSON path key for the metric value.
- `base_value_key` (String) JSON path key for the base value.
- `instance_name_key` (String) JSON path key for the instance name.
- `container_name_key` (String) JSON path key for the container name.
- `timestamp_key` (String) JSON path key for the timestamp.
- `timestamp_format` (String) Format string for the timestamp value.
- `data_filter` (String) Filter expression applied to log data before parsing (e.g., `alert->error->message=.*(?i)EFS.*`).
- `metric_name` (String) Static name for the derived metric.
- `additional_metric_name` (String) JSON path key for an additional metric name.
- `operation` (Number) Parser operation type.
- `aggregation_mode` (Number) Aggregation mode.
- `aggregation_period` (Number) Aggregation period.
- `container_type` (Number) Container type identifier.
- `grouping_by_component` (Boolean) Whether to group metrics by component.
- `derived_value_model` (Attributes) Derived value transformation model. See [below for nested schema](#nestedatt--l2m_settings--json_parsers--derived_value_model).

<a id="nestedatt--l2m_settings--json_parsers--derived_value_model"></a>
### Nested Schema for `l2m_settings.json_parsers.derived_value_model`

Optional:

- `base_value` (String) Base value expression (e.g., `alert->field=value`).
- `actual_value` (String) Actual value expression.
- `operation` (Number) Derived value operation type.
- `mapping_id_list` (List of String) List of JSON path keys used for mapping.

## Import

Projects can be imported using the project name:

```shell
terraform import insightfinder_project.example my-project-name
```
