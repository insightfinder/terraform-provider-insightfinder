---
page_title: "insightfinder_log_labels Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages log filtering and labeling rules for InsightFinder projects.
---

# insightfinder_log_labels (Resource)

Manages log filtering and labeling configuration for projects. Configure whitelists, blacklists, pattern naming, and training filters to optimize log analysis.

## Example Usage

### Whitelist Configuration

```terraform
resource "insightfinder_log_labels" "errors" {
  project_name = "application-logs"
  
  log_label_settings = [
    {
      label_type       = "whitelist"
      log_label_string = jsonencode([
        {
          type           = "fieldName"
          keyword        = "severity=error|critical|fatal"
          isCritical     = true
          isHotEventOnly = false
        }
      ])
    }
  ]
}
```

### Pattern Naming

```terraform
resource "insightfinder_log_labels" "patterns" {
  project_name = "application-logs"
  
  log_label_settings = [
    {
      label_type       = "patternName"
      log_label_string = jsonencode([
        {
          type           = "fieldName"
          keyword        = "message"
          patternNameKey = "message"
        }
      ])
    }
  ]
}
```

### Complete Configuration

```terraform
resource "insightfinder_log_labels" "complete" {
  project_name = "application-logs"
  
  log_label_settings = [
    # Whitelist critical errors
    {
      label_type       = "whitelist"
      log_label_string = jsonencode([
        {
          type           = "fieldName"
          keyword        = "severity=error|critical"
          isCritical     = true
          isHotEventOnly = false
        }
      ])
    },
    
    # Blacklist noise
    {
      label_type       = "blacklist"
      log_label_string = jsonencode([
        {
          type    = "fieldName"
          keyword = "healthcheck|ping"
        }
      ])
    },
    
    # Training whitelist
    {
      label_type       = "trainingWhitelist"
      log_label_string = jsonencode([
        {
          type    = "fieldName"
          keyword = "service_name"
        }
      ])
    },
    
    # Pattern naming
    {
      label_type       = "patternName"
      log_label_string = jsonencode([
        {
          type           = "fieldName"
          keyword        = "message"
          patternNameKey = "message"
        }
      ])
    }
  ]
}
```

### With Project Dependency

```terraform
resource "insightfinder_project" "app" {
  project_name = "my-app-logs"
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
}

resource "insightfinder_log_labels" "app_labels" {
  project_name = insightfinder_project.app.project_name
  
  log_label_settings = [
    {
      label_type       = "whitelist"
      log_label_string = jsonencode([{
        type           = "fieldName"
        keyword        = "level=ERROR|FATAL"
        isCritical     = true
        isHotEventOnly = false
      }])
    }
  ]
}
```

## Schema

### Required

- `project_name` (String) Name of the project to configure labels for
- `log_label_settings` (List of Object) List of label configurations
  - `label_type` (String) Type of label: `whitelist`, `blacklist`, `trainingWhitelist`, `patternName`
  - `log_label_string` (String) JSON-encoded array of label rules

### Label Rule Schema

For `whitelist` and `blacklist`:
```json
{
  "type": "fieldName",
  "keyword": "field=regex|pattern",
  "isCritical": true,
  "isHotEventOnly": false
}
```

For `trainingWhitelist`:
```json
{
  "type": "fieldName",
  "keyword": "field_name"
}
```

For `patternName`:
```json
{
  "type": "fieldName",
  "keyword": "field_name",
  "patternNameKey": "field_name"
}
```

### Read-Only

- `id` (String) Log labels identifier (same as project_name)

## Import

Log labels can be imported using the project name:

```shell
terraform import insightfinder_log_labels.example my-project-name
```

## Notes

- The project must exist before configuring log labels
- Use `jsonencode()` to properly format label strings
- Field names are case-sensitive
- Regular expressions are supported in keyword fields
- Multiple label types can be configured simultaneously
- Empty `log_label_settings` will remove all labels from the project
