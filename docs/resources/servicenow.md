---
page_title: "insightfinder_servicenow Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages ServiceNow integration for InsightFinder incident management.
---

# insightfinder_servicenow (Resource)

Manages ServiceNow integration configuration. Allows InsightFinder to create and update ServiceNow incidents based on detected anomalies and predictions.

## Example Usage

### Basic Authentication

```terraform
resource "insightfinder_servicenow" "basic" {
  account        = "admin"
  service_host   = "https://dev12345.service-now.com/"
  password       = var.servicenow_password
  system_names   = ["Production"]
  options        = ["Root Cause"]
  content_option = ["SUMMARY"]
  auth_type      = "basic"
}
```

### OAuth Authentication

```terraform
resource "insightfinder_servicenow" "oauth" {
  account        = "admin"
  service_host   = "https://company.service-now.com/"
  password       = var.servicenow_password
  app_id         = var.servicenow_app_id
  app_key        = var.servicenow_app_key
  system_names   = ["Production", "Staging"]
  options        = ["Root Cause", "Prediction"]
  content_option = ["SUMMARY", "DESCRIPTION"]
  auth_type      = "oauth"
}
```

### Full Configuration with Per-Project Settings and Table Mapping

```terraform
resource "insightfinder_servicenow" "full" {
  account      = "serviceaccount"
  service_host = "https://company.service-now.com/"
  password     = var.servicenow_password
  system_names = [
    "Production-US",
    "Production-EU",
    "Production-APAC"
  ]
  options        = ["Root Cause"]
  content_option = ["SUMMARY"]
  auth_type      = "basic"
  proxy          = "http://proxy.company.com:8080"

  # Write content to a custom ServiceNow field
  service_now_field = "u_probable_cause"

  # Use comments instead of work notes
  content_source = "comments"

  # Only correlate events within a 7-day window
  trigger_window_in_mills = 604800000

  # Enable feedback collection from ServiceNow
  enable_feedback_collect = true

  # Filter ticket creation by a specific field value
  ticket_created_by_source_key   = "activity_due"
  ticket_created_by_source_value = "xyzzz"

  # Associate created tickets with a CMDB configuration item
  configuration_item = "My-Server-CI"

  # Associate tickets with a ServiceNow department
  department_id = "9c777c281bd335906f2ca797b04bcb9e"

  # Per-project ticket creation, update, and consolidation settings
  project_configs = {
    "my-project" = {
      enable_ticket_creation                    = true
      enable_ticket_update                      = true
      enable_incident_consolidation_info_update = false
      enable_incident_resolve_update            = true
      configuration_item                        = "My-Server-CI"
    }
    "another-project" = {
      enable_ticket_creation                    = false
      enable_ticket_update                      = false
      enable_incident_consolidation_info_update = true
      enable_incident_resolve_update            = false
      configuration_item                        = "Another-Server-CI"
    }
  }

  # Map InsightFinder projects to ServiceNow tables
  table_mapping = {
    "my-project"      = "incident"
    "another-project" = "problem"
  }

  # Classify resolution codes as positive/negative feedback
  resolution_code_rules = [
    {
      pattern = "^(Could Not Replicate|No Fault Found|No Trouble Found)"
      outcome = "disLike"
    },
    {
      pattern = "^Solved"
      outcome = "like"
    },
  ]
}
```

## Schema

### Required

- `account` (String) ServiceNow account username. Forces replacement on change.
- `service_host` (String) ServiceNow instance URL (e.g., `https://dev12345.service-now.com/`). Forces replacement on change.
- `password` (String, Sensitive) ServiceNow account password.
- `options` (Set of String) Integration options (e.g., `Root Cause`, `Prediction`).
- `content_option` (Set of String) Incident content fields (e.g., `SUMMARY`, `DESCRIPTION`, `IMPACT`).

### Optional

- `auth_type` (String) Authentication type: `basic` or `oauth`. Default: `basic`.
- `app_id` (String) ServiceNow OAuth application ID (required when `auth_type = "oauth"`).
- `app_key` (String, Sensitive) ServiceNow OAuth application key (required when `auth_type = "oauth"`).
- `proxy` (String, Computed) Proxy server URL if required.
- `system_names` (List of String) List of InsightFinder system names to integrate.
- `service_now_field` (String) ServiceNow field to write integration content to (e.g., `u_probable_cause`).
- `content_source` (String, Computed) ServiceNow field to write incident notes to (e.g., `work_notes`, `comments`). Defaults to `work_notes`.
- `trigger_window_in_mills` (Number) Time window in milliseconds within which events are correlated into a single incident (e.g., `604800000` for 7 days).
- `enable_feedback_collect` (Boolean, Computed) Whether to enable ServiceNow feedback collection. Defaults to `false`.
- `ticket_created_by_source_key` (String) ServiceNow field key used to filter when a ticket is created (e.g., `activity_source`).
- `ticket_created_by_source_value` (String) Value matched against `ticket_created_by_source_key` to determine whether to create a ticket.
- `configuration_item` (String) ServiceNow CMDB configuration item to associate with created tickets.
- `department_id` (String) ServiceNow department ID to associate with created tickets.
- `project_configs` (Map of Object) Per-project ServiceNow ticket configuration. Each key is an InsightFinder project name, and the value is an object with:
  - `enable_ticket_creation` (Boolean, Computed) Whether to enable ticket creation for this project. Defaults to `false`.
  - `enable_ticket_update` (Boolean, Computed) Whether to enable ticket updates for this project. Defaults to `false`.
  - `enable_incident_consolidation_info_update` (Boolean, Computed) Whether to enable incident consolidation info updates for this project. Defaults to `false`.
  - `enable_incident_resolve_update` (Boolean, Computed) Whether to enable incident resolve updates for this project. Defaults to `false`.
  - `configuration_item` (String) ServiceNow CMDB configuration item for this specific project. Overrides the top-level `configuration_item` when set.
- `table_mapping` (Map of String) Mapping of InsightFinder project names to ServiceNow table names (e.g., `{ "my-project" = "incident" }`).
- `resolution_code_rules` (List of Object) Ordered list of pattern-based rules used to classify ServiceNow resolution/close codes as positive or negative feedback. Each object has:
  - `pattern` (String, Required) Regular expression matched against the ServiceNow resolution/close code (e.g., `^Solved`).
  - `outcome` (String, Required) Feedback outcome when the pattern matches. Must be `like` or `disLike`.

### Read-Only

- `id` (String) Integration identifier (`account@service_host`).

## Import

ServiceNow integrations can be imported using the format `account@service_host`:

```shell
terraform import insightfinder_servicenow.example admin@https://dev12345.service-now.com/
```

## Notes

- When using OAuth authentication, both `app_id` and `app_key` are required.
- `options` and `content_option` are sets — order does not matter and will not cause plan drift.
- System names are automatically resolved to system IDs internally.
- `table_mapping` sends each `[projectName, tableName]` pair to a dedicated PUT endpoint.
