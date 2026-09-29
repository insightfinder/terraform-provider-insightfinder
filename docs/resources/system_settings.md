---
page_title: "insightfinder_system_settings Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages knowledge base, notification, and miscellaneous settings for an InsightFinder system.
---

# insightfinder_system_settings (Resource)

Manages knowledge base, notification/alert, and miscellaneous system framework settings for an InsightFinder system. All three blocks (`knowledgebase_settings`, `notifications_settings`, `miscellaneous_settings`) are optional — you can configure one or more independently.

Deleting this resource removes it from Terraform state only; the underlying settings on the InsightFinder server are left unchanged.

## Example Usage

### Full Configuration

```terraform
resource "insightfinder_system_settings" "example" {
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

    rule_active_threshold           = 0.8
    rule_inactive_threshold         = 0.1
    rule_active_condition           = 0
    false_positive_tolerance        = 1
    kb_training_length              = 172800000
    tolerance                       = 0.08
    enable_insensitive_rule_matching = true
  }

  notifications_settings = {
    order                = 0
    hide_flag            = false
    aggregation_interval = 10
    enable_splunk_export = false

    prediction_email   = ""
    alert_health_score = 0.0
    alert_frequency    = 0

    email_dampening_period             = 59280000
    alerts_email_dampening_period      = 3600000
    prediction_email_dampening_period  = 3600000
    incident_dampening_window          = 59220000
    ticket_open_time                   = 59940000

    enable_system_down_email_alert         = false
    only_send_with_rca                     = false
    enable_incident_prediction_email_alert = true
    enable_incident_detection_email_alert  = true
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false

    alert_email             = ""
    health_alert_email      = ""
    incident_detection_email = ""
    root_cause_email        = ""

    incident_count_threshold               = jsonencode({})
    assignment_map                         = jsonencode({})
    component_level_incident_consolidation = true
    enabled_consolidation_algorithms       = ["derivedIncidents", "rcaChain", "contentBased", "metricInstanceTimestamp"]

    system_down_notification = {
      enable_system_down_email_alert = true
      email_dampening_period         = 3600000
      email_set                      = ["ops@example.com"]
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
        project_name              = "my-metric-project"
        instance_down_enable      = true
        instance_down_dampening   = 3600000
        instance_down_threshold   = 3600000
        instance_down_report_number = 1
        instance_down_emails      = ["ops@example.com"]
      }
    ]
  }

  miscellaneous_settings = {
    healthview_longterm                       = false
    should_auto_share                         = true
    rootcause_reverse_entry_filter_threshold  = 99
    enable_composite_timeline                 = true
  }
}
```

### Notifications with Project-Level Dampening Windows

```terraform
resource "insightfinder_system_settings" "with_dampening_windows" {
  system_name = "my-production-system"

  notifications_settings = {
    aggregation_interval                   = 10
    order                                  = 0
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_system_down_email_alert         = false
    enable_incident_prediction_email_alert = true
    enable_incident_detection_email_alert  = true
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    prediction_email                       = ""
    alert_email                            = ""
    health_alert_email                     = ""
    incident_detection_email               = ""
    root_cause_email                       = ""
    alert_health_score                     = 0.0
    alert_frequency                        = 0
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 14400000
    incident_count_threshold               = jsonencode({ "my-llm-project@admin" = 3 })
    assignment_map                         = jsonencode({})
    component_level_incident_consolidation = false
    enabled_consolidation_algorithms       = ["derivedIncidents", "contentBased"]

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

    project_level_dampening_periods = [
      {
        project  = "change-detection-project"
        duration = 3600000
      }
    ]
  }
}
```

### Notifications with System Down, Reports, and Instance Down

```terraform
resource "insightfinder_system_settings" "notifications_extended" {
  system_name = "my-production-system"

  notifications_settings = {
    prediction_email                       = "alerts@example.com"
    enable_incident_prediction_email_alert = true
    enable_incident_detection_email_alert  = true
    alert_health_score                     = 0.5
    aggregation_interval                   = 10
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 3600000
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
      email_set                      = ["oncall@example.com", "ops@example.com"]
    }

    daily_report_notification = {
      enable_insights_report = true
      email_set              = ["manager@example.com"]
    }

    weekly_report_notification = {
      enable_insights_report = true
      email_set              = ["executive@example.com", "manager@example.com"]
    }

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
  }
}
```

### Notifications with Custom Consolidation Rules and Metric-Log Configs

```terraform
resource "insightfinder_system_settings" "with_consolidation" {
  system_name = "my-production-system"

  notifications_settings = {
    aggregation_interval                   = 10
    order                                  = 0
    hide_flag                              = false
    enable_splunk_export                   = false
    only_send_with_rca                     = false
    enable_system_down_email_alert         = false
    enable_incident_prediction_email_alert = true
    enable_incident_detection_email_alert  = true
    enable_alerts_email                    = false
    enable_health_email_alert              = false
    enable_root_cause_email_alert          = false
    prediction_email                       = ""
    alert_email                            = ""
    health_alert_email                     = ""
    incident_detection_email               = ""
    root_cause_email                       = ""
    alert_health_score                     = 0.0
    alert_frequency                        = 0
    email_dampening_period                 = 3600000
    alerts_email_dampening_period          = 3600000
    prediction_email_dampening_period      = 3600000
    incident_dampening_window              = 14400000
    ticket_open_time                       = 5400000
    incident_count_threshold               = jsonencode({})
    assignment_map                         = jsonencode({})
    component_level_incident_consolidation = false
    enabled_consolidation_algorithms       = ["contentBased", "metricInstanceTimestamp", "consolidationCustom"]

    max_notification_delay_tolerance = 10800000

    custom_consolidation_rules = [
      {
        project_entries = [
          {
            project_name = "frontend-logs"
            conditions = [
              {
                type    = "fieldName"
                keyword = "alert->options->severity=critical"
              },
              {
                type    = "content"
                keyword = "OutOfMemoryError"
              }
            ]
          },
          {
            project_name = "backend-metrics"
            conditions = [
              {
                type    = "fieldName"
                keyword = "region=us-east-1"
              }
            ]
          }
        ]
        field_correlations = [
          {
            project_field_keys = [
              {
                project_name = "frontend-logs"
                type         = "fieldName"
                field_key    = "alert->server->host"
              },
              {
                project_name = "backend-metrics"
                type         = "fieldName"
                field_key    = "hostname"
              }
            ]
          },
          {
            project_field_keys = [
              {
                project_name = "frontend-logs"
                type         = "content"
              },
              {
                project_name = "backend-metrics"
                type         = "content"
              }
            ]
          }
        ]
      }
    ]

    metric_log_consolidation_configs = [
      {
        metric_project_name = "backend-metrics"
        log_project_name    = "frontend-logs"
        field_keys          = ["alert->server->ip", "alert->asset->asset_id"]
      }
    ]
  }
}
```

### Miscellaneous Settings Only

```terraform
resource "insightfinder_system_settings" "misc_only" {
  system_name = "my-production-system"

  miscellaneous_settings = {
    healthview_longterm                      = false
    should_auto_share                        = true
    rootcause_reverse_entry_filter_threshold = 99
    enable_composite_timeline                = true
  }
}
```

### Knowledge Base Only

```terraform
resource "insightfinder_system_settings" "kb_only" {
  system_name = "my-system"

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
```

### Satellite System Linking

```terraform
resource "insightfinder_system_settings" "with_satellite" {
  system_name = "primary-system"

  knowledgebase_settings = {
    enable_global_knowledge_base      = true
    composite_valid_threshold         = 900000
    timeline_top_k                    = 50
    enable_ignore_instance_prediction = false
    prediction_source                 = 0
    share_system_type                 = 1
    action_execution_time             = 15
    auto_fix_validation_window        = 1
    filter_self_to_self               = true
    rule_source_type                  = 0

    satellite_system_set = jsonencode([
      {
        systemPartitionKey = {
          userName   = "admin"
          systemName = "satellite-system-id"
          envName    = "All"
        }
        replay = false
      }
    ])

    rule_active_threshold            = 0.8
    rule_inactive_threshold          = 0.1
    rule_active_condition            = 0
    false_positive_tolerance         = 1
    kb_training_length               = 172800000
    tolerance                        = 0.08
    enable_insensitive_rule_matching = true
  }
}
```

## Schema

### Required

- `system_name` (String) The display name of the InsightFinder system to configure. Used to resolve the system ID. Forces replacement if changed.

### Optional

- `knowledgebase_settings` (Attributes) Knowledge base and incident prediction settings for the system. See [knowledgebase_settings](#nested-schema-for-knowledgebase_settings) below.
- `notifications_settings` (Attributes) Notification and alert email settings for the system. See [notifications_settings](#nested-schema-for-notifications_settings) below.
- `miscellaneous_settings` (Attributes) Miscellaneous system framework settings. See [miscellaneous_settings](#nested-schema-for-miscellaneous_settings) below.

### Read-Only

- `id` (String) Identifier for this resource (same as `system_name`).

---

### Nested Schema for `knowledgebase_settings`

All attributes are Optional and Computed. The block is written to the API as a whole on every apply: an attribute left unset re-sends its prior state value on update (on the initial create, where there is no prior state, the type's zero value is sent), so set explicitly every value you want to control.

#### Global Knowledge Base

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_global_knowledge_base` | Boolean | Enable the global knowledge base for this system. |
| `composite_valid_threshold` | Number | Minimum validity threshold (milliseconds) for composite knowledge base entries. |
| `timeline_top_k` | Number | Number of top-K timeline entries to retain in the knowledge base. |
| `enable_ignore_instance_prediction` | Boolean | When enabled, instance-level predictions are excluded from the knowledge base. |
| `prediction_source` | Number | Prediction source type (`0` = default, `1` = custom). |
| `share_system_type` | Number | Sharing type for the knowledge base across systems (`0` = disabled, `1` = shared). |
| `action_execution_time` | Number | Time window (minutes) for executing automated actions. |
| `auto_fix_validation_window` | Number | Validation window (hours) for auto-fix actions. |
| `filter_self_to_self` | Boolean | Filter out self-to-self causal relationships in the knowledge base. |
| `rule_source_type` | Number | Source type for knowledge base rules (`0` = default). |
| `satellite_system_set` | String | JSON-encoded array of satellite systems linked to this system's knowledge base. Each entry has a `systemPartitionKey` object (`userName`, `systemName`, `envName`) and a `replay` boolean. Use `jsonencode([])` for no satellite systems. Example: `jsonencode([{systemPartitionKey={userName="admin",systemName="<id>",envName="All"},replay=false}])` |

#### Incident Prediction

| Attribute | Type | Description |
|-----------|------|-------------|
| `rule_active_threshold` | Number | Minimum probability (0.0–1.0) to promote a causal prediction rule to active status. |
| `rule_inactive_threshold` | Number | Probability (0.0–1.0) below which a prediction rule is demoted. Should be less than or equal to `rule_active_threshold` (not validated by the provider). |
| `rule_active_condition` | Number | Prerequisite a rule must meet before generating alerts (`0` = unfiltered, `1` = verified only). |
| `false_positive_tolerance` | Number | Number of allowed false positives before a rule is deactivated. |
| `kb_training_length` | Number | Training window length for KB rules in milliseconds (e.g., `172800000` = 2 days). |
| `tolerance` | Number | Tolerance value for incident prediction scoring (0.0–1.0). |
| `enable_insensitive_rule_matching` | Boolean | Enable case-insensitive matching when applying KB rules. |

---

### Nested Schema for `notifications_settings`

Unless stated otherwise below, attributes are Optional and Computed. The health view portion of this block is written to the API as a whole on every apply: a scalar attribute left unset re-sends its prior state value on update (on the initial create the type's zero value is sent), so set explicitly every value you want to control. Exceptions:

- `anomaly_score_notification_sensitivity` is Read-Only (Computed only).
- `system_down_notification`, `daily_report_notification`, `weekly_report_notification`, `instance_down_notification`, `project_level_dampening_windows`, and `project_level_dampening_periods` are Optional only (not Computed).

#### Health View Display

| Attribute | Type | Description |
|-----------|------|-------------|
| `order` | Number | Display order for this system in the health view (`0` = first). |
| `hide_flag` | Boolean | Hide this system from the health view dashboard. |
| `aggregation_interval` | Number | Aggregation interval in minutes for health view metrics. |
| `enable_splunk_export` | Boolean | Enable exporting health view data to Splunk. |

#### Alert Thresholds

| Attribute | Type | Description |
|-----------|------|-------------|
| `alert_health_score` | Number | Health score threshold (0.0–1.0) below which an alert is triggered. |
| `alert_frequency` | Number | Alert frequency limiter — maximum number of alerts per interval. |

#### Email Dampening

| Attribute | Type | Description |
|-----------|------|-------------|
| `email_dampening_period` | Number | Dampening period for health alert emails in milliseconds. |
| `alerts_email_dampening_period` | Number | Dampening period for alert emails in milliseconds. |
| `prediction_email_dampening_period` | Number | Dampening period for prediction emails in milliseconds. |
| `incident_dampening_window` | Number | Dampening window for incident notification emails in milliseconds. |
| `ticket_open_time` | Number | Time window in milliseconds to keep a ticket open after an incident resolves. |

#### Email Alert Toggles

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_system_down_email_alert` | Boolean | Send an email alert when the system is detected as down. |
| `only_send_with_rca` | Boolean | Only send incident notifications when a root cause analysis result is available. |
| `enable_incident_prediction_email_alert` | Boolean | Send email alerts for incident prediction events. |
| `enable_incident_detection_email_alert` | Boolean | Send email alerts for incident detection events. |
| `enable_alerts_email` | Boolean | Send alert emails (metric/log anomaly alerts). |
| `enable_health_email_alert` | Boolean | Send email alerts when the health score drops below `alert_health_score`. |
| `enable_root_cause_email_alert` | Boolean | Send email alerts when a root cause analysis completes. |

#### Email Recipients

| Attribute | Type | Description |
|-----------|------|-------------|
| `prediction_email` | String | Comma-separated email addresses for prediction notifications. |
| `alert_email` | String | Comma-separated email addresses for alert notifications. |
| `health_alert_email` | String | Comma-separated email addresses for health score alerts. |
| `incident_detection_email` | String | Comma-separated email addresses for incident detection notifications. |
| `root_cause_email` | String | Comma-separated email addresses for root cause analysis notifications. |

#### JSON-Encoded Map Fields

| Attribute | Type | Description |
|-----------|------|-------------|
| `incident_count_threshold` | String | JSON-encoded map of `"ProjectName@username"` keys to integer thresholds. An incident alert is suppressed until the count exceeds the threshold for that project. Example: `jsonencode({"MyProject@admin": 5})`. Use `jsonencode({})` for no thresholds. |
| `assignment_map` | String | JSON-encoded map of zone/component keys to assignee lists. Each value is an object with `jiraAssignees`, `emailAssignees`, and `serviceNowAssignees` arrays. Use `jsonencode({})` for no assignments. |

#### Incident Consolidation

| Attribute | Type | Description |
|-----------|------|-------------|
| `component_level_incident_consolidation` | Boolean | Enable component-level incident consolidation. When enabled, incidents from different components are consolidated before alerting. Maps to `componentLevelIncidentConsolidation` in the health view API. |
| `component_level_dampening` | Boolean | Enable component-level dampening. Maps to `componentLevelDampening` in the health view API. |
| `enabled_consolidation_algorithms` | List of String | Consolidation algorithms to apply. Supported values: `"derivedIncidents"`, `"rcaChain"`, `"contentBased"`, `"metricInstanceTimestamp"`, `"consolidationCustom"` (not validated by the provider). Example: `["derivedIncidents", "rcaChain", "contentBased", "metricInstanceTimestamp"]`. |

#### System Down Notification

`system_down_notification` (Attributes, Optional) — configures system-down alerts via a dedicated API (`/api/external/v2/systemdownsetting`). Only read and written when the block is configured. All nested attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_system_down_email_alert` | Boolean | Enable email alert when the system is detected as down. |
| `email_dampening_period` | Number | Minimum interval in milliseconds between repeated system-down email alerts. |
| `email_set` | List of String | Email addresses to notify when the system goes down. |

#### Daily Report Notification

`daily_report_notification` (Attributes, Optional) — configures daily insights report emails via `/api/external/v1/insightsreportsetting`. Only read and written when the block is configured. All nested attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_insights_report` | Boolean | Enable the daily insights summary email for this system. |
| `email_set` | List of String | Email addresses to receive the daily report. |

#### Weekly Report Notification

`weekly_report_notification` (Attributes, Optional) — configures weekly insights report emails (same API as daily, `isDaily=false`). Only read and written when the block is configured. All nested attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_insights_report` | Boolean | Enable the weekly insights summary email for this system. |
| `email_set` | List of String | Email addresses to receive the weekly report. |

#### Instance Down Notification

`instance_down_notification` (List of Attributes, Optional) — a list of per-project instance-down alert configurations via `/api/external/v1/projects/update` (one API call per entry). Each entry configures one project; only the listed projects are read back. All nested attributes other than `project_name` are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `project_name` | String (Required) | The project to configure instance-down alerts for. |
| `instance_down_enable` | Boolean | Enable instance-down detection for this project. |
| `instance_down_dampening` | Number | Dampening window in milliseconds between repeated instance-down alerts. |
| `instance_down_threshold` | Number | Duration in milliseconds before an instance is considered down. |
| `instance_down_report_number` | Number | Number of instance-down events to include in the report. |
| `instance_down_emails` | List of String | Email addresses to notify when instances go down. |

#### Project Level Dampening Windows

`project_level_dampening_windows` (Set of Attributes, Optional) — a set of project-pair dampening window rules stored in the health view setting. Each rule overrides the system-level `incident_dampening_window` for a specific source→target project relationship. Order does not matter — Terraform compares entries by value regardless of the order returned by the API.

| Attribute | Type | Description |
|-----------|------|-------------|
| `source_project` | String (Required) | The source project name (`ps`). |
| `target_project` | String (Required) | The target project name (`pt`). |
| `source_customer` | String (Optional, Computed) | Customer (username) of the source project (`cs`). Defaults to the provider username when omitted. |
| `target_customer` | String (Optional, Computed) | Customer (username) of the target project (`ct`). Defaults to the provider username when omitted. |
| `duration` | Number (Required) | Dampening duration in milliseconds (`d`). |
| `similarity_threshold` | Number (Optional, Computed) | Similarity threshold for this dampening window (`st`). |

#### Project Level Dampening Periods

`project_level_dampening_periods` (Set of Attributes, Optional) — a set of per-project dampening period rules stored in the health view setting, distinct from `project_level_dampening_windows`. Each rule overrides the system-level `incident_dampening_window` for a specific project, without a target project or similarity threshold.

| Attribute | Type | Description |
|-----------|------|-------------|
| `project` | String (Required) | The project name (`p`). |
| `customer` | String (Optional, Computed) | Customer (username) of the project (`c`). Defaults to the provider username when omitted. |
| `duration` | Number (Required) | Dampening duration in milliseconds (`d`). |

#### Max Notification Delay Tolerance

| Attribute | Type | Description |
|-----------|------|-------------|
| `max_notification_delay_tolerance` | Number | Maximum delay in milliseconds before a notification must fire regardless of dampening windows (e.g. `10800000` = 3 hours). Maps to `maxNotificationDelayTolerance`. |
| `metric_co_occurrence_buffer_ms` | Number | Metric co-occurrence buffer window in milliseconds. Maps to `metricCoOccurrenceBufferMs`. |

#### Anomaly Score Notifications

| Attribute | Type | Description |
|-----------|------|-------------|
| `anomaly_score_notification_min_delta` | Number | Minimum delta threshold for anomaly score notifications. Maps to `anomalyScoreNotificationMinDelta`. |
| `anomaly_score_notification_sensitivity` | String (Read-Only) | Anomaly score notification sensitivity, derived server-side from `local_kb_sensitivities`. Maps to `anomalyScoreNotificationSensitivity`. |

#### Local KB Sensitivities

| Attribute | Type | Description |
|-----------|------|-------------|
| `local_kb_sensitivities` | String | JSON-encoded array of per-project knowledge base sensitivity overrides, mirroring the API's internal format: `p` (project name), `u` (user), `s` (sensitivity), `m` (mode), `d` (delta), `sa` (bool), `cc` (map), `udc` (list), `crp` (list), `sc` (list). Maps to `localKbSensitivities`. Example: `jsonencode([{p = "MyProject", u = "user", s = 2, m = 0, d = 43, sa = false, cc = {}, udc = [], crp = [], sc = []}])`. Use `jsonencode([])` for no overrides. |

#### Notification Delay Config

| Attribute | Type | Description |
|-----------|------|-------------|
| `notification_delay_config` | String | JSON-encoded object configuring per-project notification delay overrides. Fields: `e` (enabled bool), `d` (delay in milliseconds), `u` (username), `p` (map of project name to `{d: delay in milliseconds}`). Maps to `notificationDelayConfig`. Example: `jsonencode({e = true, d = 3300000, u = "admin", p = {"MyProject" = {d = 3300000}}})`. Use `jsonencode({})` to clear. |

#### Dependency Consolidation Setting

| Attribute | Type | Description |
|-----------|------|-------------|
| `dependency_consolidation_setting` | String | JSON-encoded object configuring dependency-based incident consolidation. Fields: `sn` (enabled bool), `lw` (lookback window in milliseconds). Maps to `dependencyConsolidationSetting`. Example: `jsonencode({sn = true, lw = 1020000})`. Use `jsonencode({})` to clear. |

#### Custom Consolidation Rules

`custom_consolidation_rules` (List of Attributes, Optional, Computed) — a list of custom incident consolidation rules. When `consolidationCustom` is included in `enabled_consolidation_algorithms`, these rules control which incidents from different projects are consolidated into a single notification. Each rule has two sub-blocks:

**`project_entries`** (List of Attributes, Optional) — projects and their keyword/field matching conditions:

| Attribute | Type | Description |
|-----------|------|-------------|
| `project_name` | String (Required) | The project this entry applies to. |
| `conditions` | List of Object (Optional) | Matching conditions. Each has a `type` (String, Required: `"fieldName"` or `"content"`) and a `keyword` (String, Required: the keyword or field expression to match). |

**`field_correlations`** (List of Attributes, Optional) — cross-project field mappings that determine which field values must match for incidents to be consolidated:

| Attribute | Type | Description |
|-----------|------|-------------|
| `project_field_keys` | List of Object (Required) | One entry per project in the correlation. Each has `project_name` (String, Required), `type` (String, Required: `"fieldName"` or `"content"`), and `field_key` (String, Optional, Computed: the field path; omit or set to `null` for `"content"` type). |

#### Metric-Log Consolidation Configs

`metric_log_consolidation_configs` (List of Attributes, Optional, Computed) — a list of metric-to-log project consolidation mappings. Each entry links one metric project with one log project and specifies the field keys used to correlate their incidents.

| Attribute | Type | Description |
|-----------|------|-------------|
| `metric_project_name` | String (Required) | The metric project name. |
| `log_project_name` | String (Required) | The log project name. |
| `field_keys` | List of String (Optional, Computed) | Field key paths used to match incidents between the metric and log projects. |

---

### Nested Schema for `miscellaneous_settings`

Configures miscellaneous system framework settings via `/api/external/v1/systemframework`. All attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `healthview_longterm` | Boolean | Enable long-term storage mode for the system health view. Written via `operation=hideOrOrderOrLongTerm`; the current system order is read first to avoid overwriting it. |
| `should_auto_share` | Boolean | Enable automatic sharing of system data with linked systems. |
| `rootcause_reverse_entry_filter_threshold` | Number | Threshold (0–100) for root cause reverse entry filtering. |
| `enable_composite_timeline` | Boolean | Enable the composite timeline view for the system. |

---

## Import

`insightfinder_system_settings` resources can be imported using the system name:

```shell
terraform import insightfinder_system_settings.example my-production-system
```

After import, run `terraform plan` to review which computed fields will be populated from the API.

## Notes

- **Partial configuration**: All three blocks (`knowledgebase_settings`, `notifications_settings`, `miscellaneous_settings`) are independently optional. Omitting a block means those settings are not managed by Terraform.
- **Delete behavior**: Removing this resource from Terraform state does not change the settings on the InsightFinder server. The settings persist and must be manually reset if needed.
- **`satellite_system_set`**: Must be provided as a JSON-encoded string using `jsonencode(...)`. The value is semantically compared during plan/apply to avoid spurious diffs caused by JSON key ordering.
- **`incident_count_threshold` and `assignment_map`**: Must be provided as JSON-encoded strings using `jsonencode(...)`. These fields are stored as serialized JSON in Terraform state and compared semantically to avoid key-ordering diffs.
- **`project_level_dampening_windows`**: Optional set (not Computed). Whenever `notifications_settings` is applied, the declared set replaces any existing rules on the server — omitting the attribute (or setting it to `[]`) clears them. `source_customer` and `target_customer` default to the provider username when not specified.
- **`project_level_dampening_periods`**: Optional set (not Computed), separate from `project_level_dampening_windows` (the API tracks them independently). Whenever `notifications_settings` is applied, the declared set replaces any existing rules on the server — omitting the attribute (or setting it to `[]`) clears them. `customer` defaults to the provider username when not specified.
- **`enabled_consolidation_algorithms`**: Optional and Computed list of strings. Supported algorithm names are `"derivedIncidents"`, `"rcaChain"`, `"contentBased"`, `"metricInstanceTimestamp"`, and `"consolidationCustom"`. The value is always read back from the API on refresh; set it explicitly, because an unset (unknown) value is sent to the API as an empty list on apply.
- **`max_notification_delay_tolerance`**: Optional and Computed number (milliseconds). When omitted, the prior state value is re-sent on update.
- **`custom_consolidation_rules`**: Optional and Computed list of rule objects. Always read back from the API on refresh; when set, the declared list replaces the rules on the server. Include `"consolidationCustom"` in `enabled_consolidation_algorithms` to activate these rules. For `"content"` type `project_field_keys` entries, `field_key` can be omitted or set to `null`.
- **`metric_log_consolidation_configs`**: Optional and Computed list of metric-log mapping objects. Always read back from the API on refresh; when set, the declared list replaces the mappings on the server.
- **API endpoints**: `knowledgebase_settings` maps to two separate API calls — `SetGlobalKBSetting` and `SetIncidentPredictionSetting`. `notifications_settings` maps to `SetHealthViewSetting`, plus the dedicated system-down, insights-report (daily/weekly), and per-project instance-down APIs for `system_down_notification`, `daily_report_notification`/`weekly_report_notification`, and `instance_down_notification`. `miscellaneous_settings` maps to two calls on `/api/external/v1/systemframework` — `operation=hideOrOrderOrLongTerm` for `healthview_longterm`, and `operation=systemFrameworkSetting` for the remaining three fields. All four fields are read via a single `GET /api/external/v1/systemframework` call.
