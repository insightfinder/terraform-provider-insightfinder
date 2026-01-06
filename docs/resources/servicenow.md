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
  account          = "admin"
  service_host     = "https://dev12345.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 3600000
  system_names     = ["Production"]
  options          = ["Root Cause"]
  content_option   = ["SUMMARY"]
  auth_type        = "basic"
}
```

### OAuth Authentication

```terraform
resource "insightfinder_servicenow" "oauth" {
  account          = "admin"
  service_host     = "https://company.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 7200000
  app_id           = var.servicenow_app_id
  app_key          = var.servicenow_app_key
  system_names     = ["Production", "Staging"]
  options          = ["Root Cause", "Prediction"]
  content_option   = ["SUMMARY", "DESCRIPTION"]
  auth_type        = "oauth"
}
```

### Multiple Systems

```terraform
resource "insightfinder_servicenow" "multi" {
  account          = "serviceaccount"
  service_host     = "https://company.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 3600000
  system_names     = [
    "Production-US",
    "Production-EU",
    "Production-APAC"
  ]
  options        = ["Root Cause"]
  content_option = ["SUMMARY"]
  auth_type      = "basic"
  proxy          = "http://proxy.company.com:8080"
}
```

## Schema

### Required

- `account` (String) ServiceNow account username
- `service_host` (String) ServiceNow instance URL (e.g., `https://dev12345.service-now.com/`)
- `password` (String, Sensitive) ServiceNow account password
- `dampening_period` (Number) Dampening period in milliseconds (e.g., `3600000` for 1 hour)
- `system_names` (List of String) List of InsightFinder system names to integrate
- `options` (List of String) Integration options: `Root Cause`, `Prediction`
- `content_option` (List of String) Incident content fields: `SUMMARY`, `DESCRIPTION`, `IMPACT`

### Optional

- `auth_type` (String) Authentication type: `basic` or `oauth`. Default: `basic`
- `app_id` (String) ServiceNow OAuth application ID (required when `auth_type = "oauth"`)
- `app_key` (String, Sensitive) ServiceNow OAuth application key (required when `auth_type = "oauth"`)
- `proxy` (String) Proxy server URL if required
- `system_ids` (List of String, Computed) Resolved system IDs (computed from system_names)

### Read-Only

- `id` (String) Integration identifier (`account@service_host`)

## Import

ServiceNow integrations can be imported using the format `account@service_host`:

```shell
terraform import insightfinder_servicenow.example admin@https://dev12345.service-now.com/
```

## Notes

- The `system_names` list order is preserved in the configuration
- When using OAuth authentication, both `app_id` and `app_key` are required
- System names are automatically resolved to system IDs
- The dampening period prevents duplicate incidents within the specified time window
