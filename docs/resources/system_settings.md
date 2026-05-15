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

- `system_name` (String) The name of the InsightFinder system to configure.

### Optional

- `knowledgebase_settings` (Attributes) Knowledge base and incident prediction settings for the system. See [knowledgebase_settings](#nested-schema-for-knowledgebase_settings) below.
- `notifications_settings` (Attributes) Notification and alert email settings for the system. See [notifications_settings](#nested-schema-for-notifications_settings) below.
- `miscellaneous_settings` (Attributes) Miscellaneous system framework settings. See [miscellaneous_settings](#nested-schema-for-miscellaneous_settings) below.

### Read-Only

- `id` (String) Identifier for this resource (same as `system_name`).

---

### Nested Schema for `knowledgebase_settings`

All attributes are Optional and Computed (server defaults are used when omitted).

#### Global Knowledge Base

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_global_knowledge_base` | Boolean | Enable the global knowledge base for this system. |
| `composite_valid_threshold` | Number | Minimum validity threshold (milliseconds) for composite knowledge base entries. |
| `timeline_top_k` | Number | Number of top-K timeline entries to retain in the knowledge base. |
| `enable_ignore_instance_prediction` | Boolean | When enabled, instance-level predictions are excluded from the knowledge base. |
| `prediction_source` | Number | Source of predictions used for knowledge base training (`0` = default). |
| `share_system_type` | Number | Sharing type for the knowledge base across systems (`0` = disabled, `1` = shared). |
| `action_execution_time` | Number | Time window (minutes) for executing automated actions. |
| `auto_fix_validation_window` | Number | Validation window (hours) for auto-fix actions. |
| `filter_self_to_self` | Boolean | Filter out self-to-self causal relationships in the knowledge base. |
| `rule_source_type` | Number | Source type for knowledge base rules (`0` = default). |
| `satellite_system_set` | String | JSON-encoded array of satellite systems linked to this system's knowledge base. Each entry has a `systemPartitionKey` object (`userName`, `systemName`, `envName`) and a `replay` boolean. Use `jsonencode([])` for no satellite systems. Example: `jsonencode([{systemPartitionKey={userName="admin",systemName="<id>",envName="All"},replay=false}])` |

#### Incident Prediction

| Attribute | Type | Description |
|-----------|------|-------------|
| `rule_active_threshold` | Number | Score threshold (0.0–1.0) above which a prediction rule is considered active. |
| `rule_inactive_threshold` | Number | Score threshold (0.0–1.0) below which a prediction rule is deactivated. |
| `rule_active_condition` | Number | Condition type for rule activation (`0` = default). |
| `false_positive_tolerance` | Number | Number of allowed false positives before a rule is deactivated. |
| `kb_training_length` | Number | Training window length for KB rules in milliseconds (e.g., `172800000` = 2 days). |
| `tolerance` | Number | Tolerance value for incident prediction scoring (0.0–1.0). |
| `enable_insensitive_rule_matching` | Boolean | Enable case-insensitive matching when applying KB rules. |

---

### Nested Schema for `notifications_settings`

All attributes are Optional and Computed (server defaults are used when omitted).

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
| `ticket_open_time` | Number | Time window in milliseconds for keeping a ticket open. |

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
| `enabled_consolidation_algorithms` | List of String | Consolidation algorithms to apply. Supported values: `"derivedIncidents"`, `"rcaChain"`, `"contentBased"`, `"metricInstanceTimestamp"`. Example: `["derivedIncidents", "rcaChain", "contentBased", "metricInstanceTimestamp"]`. |

#### System Down Notification

Configures system-down alerts via a dedicated API (`/api/external/v2/systemdownsetting`). All attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_system_down_email_alert` | Boolean | Enable email alert when the system is detected as down. |
| `email_dampening_period` | Number | Minimum interval in milliseconds between repeated system-down email alerts. |
| `email_set` | List of String | Email addresses to notify when the system goes down. |

#### Daily Report Notification

Configures daily insights report emails via `/api/external/v1/insightsreportsetting`. All attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_insights_report` | Boolean | Enable the daily insights summary email for this system. |
| `email_set` | List of String | Email addresses to receive the daily report. |

#### Weekly Report Notification

Configures weekly insights report emails (same API as daily, `isDaily=false`). All attributes are Optional and Computed.

| Attribute | Type | Description |
|-----------|------|-------------|
| `enable_insights_report` | Boolean | Enable the weekly insights summary email for this system. |
| `email_set` | List of String | Email addresses to receive the weekly report. |

#### Instance Down Notification

A list of per-project instance-down alert configurations via `/api/external/v1/projects/update`. Each entry configures one project.

| Attribute | Type | Description |
|-----------|------|-------------|
| `project_name` | String (Required) | The project to configure instance-down alerts for. |
| `instance_down_enable` | Boolean | Enable instance-down detection for this project. |
| `instance_down_dampening` | Number | Dampening window in milliseconds between repeated instance-down alerts. |
| `instance_down_threshold` | Number | Duration in milliseconds before an instance is considered down. |
| `instance_down_report_number` | Number | Number of instances that must be down before an alert is sent. |
| `instance_down_emails` | List of String | Email addresses to notify when instances go down. |

#### Project Level Dampening Windows

A set of project-pair dampening window rules stored in the health view setting. Each rule overrides the system-level `incident_dampening_window` for a specific source→target project relationship. Order does not matter — Terraform compares entries by value regardless of the order returned by the API.

| Attribute | Type | Description |
|-----------|------|-------------|
| `source_project` | String (Required) | The source project name (`ps`). |
| `target_project` | String (Required) | The target project name (`pt`). |
| `source_customer` | String | Customer (username) of the source project (`cs`). Defaults to the provider username when omitted. |
| `target_customer` | String | Customer (username) of the target project (`ct`). Defaults to the provider username when omitted. |
| `duration` | Number (Required) | Dampening duration in milliseconds (`d`). |

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
- **`project_level_dampening_windows`**: Optional and Computed list. When omitted, the server value is preserved in state. When set (even to `[]`), the declared list replaces any existing rules on the server. `source_customer` and `target_customer` default to the provider username when not specified.
- **`enabled_consolidation_algorithms`**: Optional and Computed list of strings. Supported algorithm names are `"derivedIncidents"`, `"rcaChain"`, `"contentBased"`, and `"metricInstanceTimestamp"`. When omitted, the server value is preserved in state.
- **API endpoints**: `knowledgebase_settings` maps to two separate API calls — `SetGlobalKBSetting` and `SetIncidentPredictionSetting`. `notifications_settings` maps to `SetHealthViewSetting`. `miscellaneous_settings` maps to two calls on `/api/external/v1/systemframework` — `operation=hideOrOrderOrLongTerm` for `healthview_longterm`, and `operation=systemFrameworkSetting` for the remaining three fields. All four fields are read via a single `GET /api/external/v1/systemframework` call.
