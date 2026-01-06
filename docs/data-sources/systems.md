---
page_title: "insightfinder_systems Data Source - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Retrieves the list of available InsightFinder systems.
---

# insightfinder_systems (Data Source)

Retrieves the list of all available InsightFinder systems, including both owned and shared systems.

## Example Usage

### List All Systems

```terraform
data "insightfinder_systems" "all" {}

output "system_names" {
  value = [for s in data.insightfinder_systems.all.systems : s.system_name]
}

output "system_count" {
  value = length(data.insightfinder_systems.all.systems)
}
```

### Find Specific System

```terraform
data "insightfinder_systems" "all" {}

locals {
  production_system = [
    for s in data.insightfinder_systems.all.systems :
    s if s.system_name == "Production"
  ][0]
}

output "production_system_id" {
  value = local.production_system.system_id
}

output "production_display_name" {
  value = local.production_system.system_display_name
}
```

### Conditional Resource Creation

```terraform
data "insightfinder_systems" "all" {}

locals {
  has_staging = contains(
    [for s in data.insightfinder_systems.all.systems : s.system_name],
    "Staging"
  )
}

resource "insightfinder_project" "staging_logs" {
  count = local.has_staging ? 1 : 0
  
  project_name = "staging-logs"
  system_name  = "Staging"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Staging Logs"
  project_time_zone    = "UTC"
  sampling_interval    = 600
}
```

### Dynamic ServiceNow Configuration

```terraform
data "insightfinder_systems" "all" {}

locals {
  production_systems = [
    for s in data.insightfinder_systems.all.systems :
    s.system_name if can(regex("^Production-", s.system_name))
  ]
}

resource "insightfinder_servicenow" "prod" {
  account          = "admin"
  service_host     = "https://company.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 3600000
  system_names     = local.production_systems
  options          = ["Root Cause"]
  content_option   = ["SUMMARY"]
  auth_type        = "basic"
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier
- `systems` (List of Object) List of available systems
  - `system_id` (String) System identifier
  - `system_name` (String) System name
  - `system_display_name` (String) System display name
  - `owner` (String) System owner username

## Notes

- This data source returns both owned systems and systems shared with the user
- System IDs are used internally by InsightFinder
- System names are used in Terraform resource configurations
- The list is refreshed on each Terraform run
