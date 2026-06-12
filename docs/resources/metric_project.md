---
page_title: "insightfinder_metric_project Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages an InsightFinder metric project with full configuration for anomaly detection, prediction, and alerting.
---

# insightfinder_metric_project (Resource)

Manages an InsightFinder metric project. Metric projects ingest time-series metric data and provide anomaly detection, incident prediction, and root cause analysis.

## Example Usage

### Basic Metric Project

```terraform
resource "insightfinder_metric_project" "app_metrics" {
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
```

### Metric Project with Process Mode (L2M)

```terraform
resource "insightfinder_metric_project" "conviva_alerts_l2m" {
  project_name = "Conviva-Alerts-Stage-L2M"
  system_name  = "my-system"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  mode = 4
}
```

### Metric Project with Detection Tuning

```terraform
resource "insightfinder_metric_project" "tuned_metrics" {
  project_name = "tuned-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "OnPremise"
    project_cloud_type = "OnPremise"
    insight_agent_type = "Custom"
  }

  project_display_name = "Tuned Metrics"
  project_time_zone    = "US/Eastern"
  sampling_interval    = 60
  retention_time       = 90
  ubl_retention_time   = 8

  # Detection tuning
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
  instance_down_ratio_threshold        = 0.05

  # Gap filling
  enable_fill_gap                  = false
  enable_store_filled_gap          = false
  gap_filling_training_data_length = 0
  enable_metric_data_prediction    = true

  # KPI prediction
  prediction_training_data_length     = 0
  prediction_correlation_sensitivity  = 0.75
  enable_kpi_prediction               = false

  # Incident prediction and RCA
  incident_prediction_window      = 12
  min_incident_prediction_window  = 5
  incident_relation_search_window = 6
  incident_prediction_event_limit = 50
  root_cause_count_threshold      = 1
  root_cause_probability_threshold = 0.8
  maximum_root_cause_result_size  = 5
  multi_hop_search_level          = 2
  multi_hop_search_limit          = "30"

  # Instance down
  instance_down_threshold      = 3600000
  instance_down_report_number  = 50
  instance_down_enable         = false
  show_instance_down           = false
}
```

### Metric Project with Holiday Settings

```terraform
resource "insightfinder_metric_project" "metrics_with_holidays" {
  project_name = "metrics-with-holidays"
  system_name  = "Production"

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
```

### Metric Project with Metric Configurations

```terraform
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
```

### Metric Project with Email and Webhook Alerts

```terraform
resource "insightfinder_metric_project" "alerted_metrics" {
  project_name = "alerted-metrics"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "AWS"
    project_cloud_type = "AWS"
    insight_agent_type = "Custom"
  }

  sampling_interval = 60
  retention_time    = 90

  enable_new_alert_email = true
  email_setting = jsonencode({
    onlySendWithRCA                    = false
    enableIncidentDetectionEmailAlert  = true
    enableIncidentPredictionEmailAlert = true
    enableRootCauseEmailAlert          = true
    emailDampeningPeriod               = 3600000
    awSeverityLevel                    = "Major"
  })

  webhook_url = "https://hooks.example.com/metrics-alerts"
  webhook_type_set_str = jsonencode([
    "metric",
    "detectedIncident",
    "predictedIncident"
  ])
  webhook_alert_dampening = 18000000
}
```

## Schema

### Required

- `project_name` (String) Unique project identifier. Forces replacement if changed.
- `system_name` (String) Name of the system this project belongs to.
- `project_creation_config` (Object) Project creation configuration.
  - `data_type` (String) Type of data — must be `Metric`.
  - `instance_type` (String) Instance type: `AWS`, `Azure`, `GCP`, `PrivateCloud`, `OnPremise`.
  - `project_cloud_type` (String) Cloud type (usually same as `instance_type`).
  - `insight_agent_type` (String, Optional, Computed) Agent type: `Custom`, `MetricFile`, etc.

### Optional — Common Settings

- `mode` (Number, Computed) Process mode for the project (set via `logdedicatedmode` API). Common values: `0` = normal, `4` = L2M (log-to-metric). Corresponds to `processMode` in the API payload.
- `project_display_name` (String, Computed) Display name for the project.
- `project_time_zone` (String, Computed) Time zone (e.g., `UTC`, `US/Eastern`).
- `sampling_interval` (Number, Computed) Data sampling interval in seconds.
- `retention_time` (Number, Computed) Data retention in days.
- `ubl_retention_time` (Number, Computed) UBL data retention in days.
- `training_filter` (Boolean, Computed) Enable training filter.
- `enable_new_alert_email` (Boolean, Computed) Enable email alerts.
- `large_project` (Boolean, Computed) Mark as a large project.
- `new_pattern_range` (Number, Computed) Range for new pattern detection.
- `proxy` (String, Computed) Proxy URL.
- `ignore_instance_for_kb` (Boolean, Computed) Ignore instance for knowledge base.
- `show_instance_down` (Boolean, Computed) Show instance-down incidents.
- `avg_per_incident_downtime_cost` (Number, Computed) Average downtime cost per incident.
- `alert_hourly_cost` (Number, Computed) Alert hourly cost.
- `alert_average_time` (Number, Computed) Alert average resolution time.

### Optional — Detection Tuning

- `c_value` (Number, Computed) C value for anomaly sensitivity (typically 2–5).
- `p_value` (Number, Computed) P value for anomaly probability (0.0–1.0).
- `high_ratio_c_value` (Number, Computed) C value for high-ratio anomalies.
- `maximum_hint` (Number, Computed) Maximum hint value.
- `dynamic_baseline_detection_flag` (Boolean, Computed) Enable dynamic baseline detection.
- `baseline_duration` (Number, Computed) Baseline duration in milliseconds.
- `positive_baseline_violation_factor` (Number, Computed) Positive violation factor.
- `negative_baseline_violation_factor` (Number, Computed) Negative violation factor.
- `enable_period_anomaly_filter` (Boolean, Computed) Enable period anomaly filter.
- `enable_ubl_detect` (Boolean, Computed) Enable UBL detection.
- `enable_cumulative_detect` (Boolean, Computed) Enable cumulative detection.
- `enable_baseline_detection_double_verify` (Boolean, Computed) Enable double-verify for baseline detection.
- `filter_by_anomaly_in_baseline_generation` (Boolean, Computed) Filter anomalies during baseline generation.
- `anomaly_dampening` (Number, Computed) Anomaly dampening window in milliseconds.
- `anomaly_gap_tolerance_count` (Number, Computed) Number of gaps to tolerate before flagging anomaly.
- `instance_down_ratio_threshold` (Number, Computed) Ratio of instances down to trigger alert.
- `model_span` (Number, Computed) Model span.
- `pattern_id_generation_rule` (Number, Computed) Rule for pattern ID generation.
- `component_name_auto_overwrite` (Boolean, Computed) Auto-overwrite component names.
- `enable_stream_detection` (Boolean, Computed) Enable stream detection.
- `enable_anomaly_score_escalation` (Boolean, Computed) Enable anomaly score escalation.
- `escalation_anomaly_score_threshold` (String, Computed) Escalation threshold.
- `ignore_anomaly_score_threshold` (String, Computed) Ignore threshold.

### Optional — Gap Filling and Prediction

- `enable_fill_gap` (Boolean, Computed) Enable gap filling.
- `enable_store_filled_gap` (Boolean, Computed) Store filled gap data.
- `gap_filling_training_data_length` (Number, Computed) Training data length for gap filling.
- `enable_metric_data_prediction` (Boolean, Computed) Enable metric data prediction.
- `prediction_training_data_length` (Number, Computed) Training data length for prediction.
- `prediction_correlation_sensitivity` (Number, Computed) Correlation sensitivity for prediction.
- `enable_kpi_prediction` (Boolean, Computed) Enable KPI prediction.

### Optional — Incident Prediction and RCA

- `incident_prediction_window` (Number, Computed) Incident prediction look-ahead window (hours).
- `min_incident_prediction_window` (Number, Computed) Minimum prediction window (hours).
- `incident_relation_search_window` (Number, Computed) Window for relation search (hours).
- `incident_prediction_event_limit` (Number, Computed) Max events for incident prediction.
- `root_cause_count_threshold` (Number, Computed) Root cause count threshold.
- `root_cause_probability_threshold` (Number, Computed) Root cause probability threshold.
- `composite_rca_limit` (Number, Computed) Composite RCA limit.
- `root_cause_log_message_search_range` (Number, Computed) Log message search range for RCA (minutes).
- `causal_prediction_setting` (Number, Computed) Causal prediction mode.
- `root_cause_rank_setting` (Number, Computed) Root cause ranking setting.
- `maximum_root_cause_result_size` (Number, Computed) Max root cause results.
- `multi_hop_search_level` (Number, Computed) Multi-hop causal search depth.
- `multi_hop_search_limit` (String, Computed) Multi-hop search limit.
- `prediction_count_threshold` (Number, Computed) Prediction count threshold.
- `prediction_probability_threshold` (Number, Computed) Prediction probability threshold.
- `prediction_rule_active_condition` (Number, Computed) Prediction rule activation condition.
- `prediction_rule_false_positive_threshold` (Number, Computed) False positive threshold for prediction rules.
- `prediction_rule_active_threshold` (Number, Computed) Active threshold for prediction rules.
- `prediction_rule_inactive_threshold` (Number, Computed) Inactive threshold for prediction rules.
- `min_valid_model_span` (Number, Computed) Minimum valid model span in milliseconds.

### Optional — Instance Down Detection

- `instance_down_threshold` (Number, Computed) Duration in ms before instance is considered down.
- `instance_down_report_number` (Number, Computed) Number of instances down before reporting.
- `instance_down_enable` (Boolean, Computed) Enable instance down detection.

### Optional — Webhook

- `webhook_url` (String, Computed) Webhook URL for notifications.
- `webhook_type_set_str` (String, Computed) JSON array of webhook event types.
- `webhook_black_list_set_str` (String, Computed) JSON array of webhook blacklist patterns.
- `webhook_critical_keyword_set_str` (String, Computed) JSON array of critical keywords.
- `webhook_alert_dampening` (Number, Computed) Webhook alert dampening window in milliseconds.
- `max_web_hook_request_size` (Number, Computed) Maximum webhook request size in MB.
- `webhook_header_list` (String, Computed) JSON array of custom webhook headers.

### Optional — Email

- `email_setting` (String, Computed) JSON-encoded email configuration object. Supported fields:
  - `onlySendWithRCA` (Boolean)
  - `enableIncidentDetectionEmailAlert` (Boolean)
  - `enableIncidentPredictionEmailAlert` (Boolean)
  - `enableRootCauseEmailAlert` (Boolean)
  - `enableAlertsEmail` (Boolean)
  - `emailDampeningPeriod` (Number, ms)
  - `alertsEmailDampeningPeriod` (Number, ms)
  - `predictionEmailDampeningPeriod` (Number, ms)
  - `awSeverityLevel` (String)

### Optional — Complex / Array Fields

- `linked_log_projects` (String, Computed) JSON array of linked log project names.
- `component_metric_setting_overall_model_list` (String, Computed) JSON array of component metric model settings.
- `shared_usernames` (String, Computed) JSON array of usernames to share the project with.
- `instance_grouping_update` (String, Computed) JSON object for instance grouping settings (e.g., `{"autoFill": false}`).

### Optional — Holidays

- `holiday_settings` (List of Objects) Holiday periods that suppress anomaly detection. Each holiday requires:
  - `name` (String, Required) Unique holiday name within the project.
  - `start_date` (String, Required) Start date in `MM-DD` format (e.g., `12-25`).
  - `end_date` (String, Required) End date in `MM-DD` format (e.g., `12-26`).

### Optional — Metric Configurations

- `metric_configurations` (Map of Objects) Per-metric alert threshold settings and component operation rules, **keyed by metric name**. Each map key is the exact metric name (e.g., `"cpu_usage"`); the value object contains:
  - `escalate_incident_components` (List of String, Optional) Component names that escalate incidents for this metric.
  - `ignored_components` (List of String, Optional) Component names excluded from anomaly detection for this metric.
  - `metric_alert_settings` (List of Objects, Optional) Per-component (or global) alert threshold rows. Each row has:
    - `component_name` (String) Component name, or empty string for the global (project-level) setting.
    - `threshold_alert_lower_bound` (String) Lower threshold for anomaly alert.
    - `threshold_alert_upper_bound` (String) Upper threshold for anomaly alert.
    - `threshold_alert_lower_bound_negative` (String) Negative direction lower alert threshold.
    - `threshold_alert_upper_bound_negative` (String) Negative direction upper alert threshold.
    - `threshold_no_alert_lower_bound` (String) Lower threshold below which no alert is raised.
    - `threshold_no_alert_upper_bound` (String) Upper threshold above which no alert is raised.
    - `threshold_no_alert_lower_bound_negative` (String) Negative direction lower no-alert threshold.
    - `threshold_no_alert_upper_bound_negative` (String) Negative direction upper no-alert threshold.
    - `incident_alert_lower_bound` (String) Lower threshold for incident alert.
    - `incident_alert_upper_bound` (String) Upper threshold for incident alert.
    - `incident_alert_lower_bound_negative` (String) Negative direction lower incident alert threshold.
    - `incident_alert_upper_bound_negative` (String) Negative direction upper incident alert threshold.
    - `incident_no_alert_lower_bound` (String) Lower threshold below which no incident alert is raised.
    - `incident_no_alert_upper_bound` (String) Upper threshold above which no incident alert is raised.
    - `incident_no_alert_lower_bound_negative` (String) Negative direction lower no-incident-alert threshold.
    - `incident_no_alert_upper_bound_negative` (String) Negative direction upper no-incident-alert threshold.
    - `is_kpi` (Boolean, Optional) Mark this metric as a KPI metric.
    - `is_flapping_result_only` (Boolean, Optional) Only report flapping anomalies.
    - `incident_duration_threshold` (Number, Optional) Minimum incident duration in milliseconds before alerting.
    - `detection_type` (String, Optional) Detection direction: `positive`, `negative`, or `both`.
    - `c_value_override` (Number, Optional, Computed) Per-metric override for the C value anomaly sensitivity. Null means use the project default.
    - `high_c_value_override` (Number, Optional, Computed) Per-metric override for the high-ratio C value anomaly sensitivity. Null means use the project default.
    - `pattern_name_higher` (String, Optional) Pattern name for values above threshold.
    - `pattern_name_lower` (String, Optional) Pattern name for values below threshold.
    - `metric_type` (String, Optional) Metric data type (e.g., `Unknown`, `CPU Utilization`, `Network Utilization`).
    - `fill_zero` (Boolean, Optional, Computed) Fill missing data points with zero.
    - `rouge_value` (String, Optional) Raw rouge value string from the API (e.g., `{"l":NaN,"s":NaN}`). Set to empty string to clear.
    - `enable_baseline_near_constance` (Boolean, Optional) Enable near-constance baseline detection.
    - `compute_difference` (Boolean, Optional) Compute derivative (difference) of this metric before detection.
    - `anomaly_gap_tolerance_duration` (Number, Optional, Computed) Anomaly gap tolerance in milliseconds. Internally converted to a count using the project `sampling_interval` before being sent to the API.

### Read-Only

- `id` (String) Project identifier (same as `project_name`).

## Import

Metric projects can be imported using the project name:

```shell
terraform import insightfinder_metric_project.example my-metric-project
```
