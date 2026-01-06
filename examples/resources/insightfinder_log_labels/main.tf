# Log Labels Examples

## Whitelist Configuration

```hcl
resource "insightfinder_log_labels" "error_whitelist" {
  project_name = "my-application-logs"
  
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

## Pattern Naming

```hcl
resource "insightfinder_log_labels" "pattern_names" {
  project_name = "my-application-logs"
  
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

## Training Whitelist

```hcl
resource "insightfinder_log_labels" "training" {
  project_name = "my-application-logs"
  
  log_label_settings = [
    {
      label_type       = "trainingWhitelist"
      log_label_string = jsonencode([
        {
          type    = "fieldName"
          keyword = "component"
        }
      ])
    }
  ]
}
```

## Combined Configuration

```hcl
resource "insightfinder_log_labels" "complete" {
  project_name = "my-application-logs"
  
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
        },
        {
          type           = "fieldName"
          keyword        = "level=ERROR|FATAL"
          isCritical     = true
          isHotEventOnly = true
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
        },
        {
          type    = "fieldName"
          keyword = "component"
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

## With Project Dependency

```hcl
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
        keyword        = "severity=.*error.*|.*critical.*"
        isCritical     = true
        isHotEventOnly = false
      }])
    }
  ]
}
```

## Notes

- Log labels must be associated with an existing project
- Multiple label types can be configured simultaneously
- Use `jsonencode()` to properly format the label strings
- Field names are case-sensitive
- Regular expressions are supported in keyword fields
