---
page_title: "insightfinder_project Data Source - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Retrieves information about an existing InsightFinder project.
---

# insightfinder_project (Data Source)

Retrieves information about an existing InsightFinder project, including its configuration, watch tower settings, and log labels.

## Example Usage

### Basic Usage

```terraform
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

### Use in Resource

```terraform
data "insightfinder_project" "template" {
  project_name = "template-project"
}

resource "insightfinder_project" "new" {
  project_name = "new-project"
  system_name  = data.insightfinder_project.template.system_name

  project_creation_config = {
    data_type          = data.insightfinder_project.template.project_creation_config.data_type
    instance_type      = data.insightfinder_project.template.project_creation_config.instance_type
    project_cloud_type = data.insightfinder_project.template.project_creation_config.project_cloud_type
    insight_agent_type = data.insightfinder_project.template.project_creation_config.insight_agent_type
  }

  project_display_name = "New Project"
  project_time_zone    = data.insightfinder_project.template.project_time_zone
  sampling_interval    = data.insightfinder_project.template.sampling_interval
}
```

## Schema

### Required

- `project_name` (String) Name of the project to retrieve

### Read-Only

All project attributes are available as read-only outputs:

- `id` (String) Project identifier
- `system_name` (String) System name
- `project_display_name` (String) Display name
- `project_time_zone` (String) Time zone
- `sampling_interval` (Number) Sampling interval in seconds
- `retention_time` (Number) Retention period in days
- `project_creation_config` (Object) Creation configuration
- `anomaly_detection_mode` (Number) Anomaly detection mode
- `email_setting` (String) Email configuration (JSON)
- `webhook_url` (String) Webhook URL
- `log_label_settings` (List) Log label configurations

And many more... See the [resource documentation](../resources/project.md) for the complete list of available attributes.
