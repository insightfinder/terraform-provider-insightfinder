# Data Source Examples

## Query Project Information

```hcl
data "insightfinder_project" "existing" {
  project_name = "my-existing-project"
}

output "project_system" {
  value = data.insightfinder_project.existing.system_name
}

output "project_timezone" {
  value = data.insightfinder_project.existing.project_time_zone
}
```

## List All Systems

```hcl
data "insightfinder_systems" "all" {}

output "system_names" {
  value = [for s in data.insightfinder_systems.all.systems : s.system_name]
}

output "system_count" {
  value = length(data.insightfinder_systems.all.systems)
}
```

## Find Specific System

```hcl
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
```

## Use Data Source with Resources

```hcl
# Query existing system
data "insightfinder_systems" "all" {}

locals {
  available_systems = [
    for s in data.insightfinder_systems.all.systems : s.system_name
  ]
}

# Create project using first available system
resource "insightfinder_project" "dynamic" {
  project_name = "dynamic-project"
  system_name  = local.available_systems[0]

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Dynamic Project"
  project_time_zone    = "UTC"
  sampling_interval    = 600
}
```

## Conditional Resource Creation

```hcl
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
