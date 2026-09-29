---
page_title: "insightfinder_log_labels Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages log filtering and labeling rules for InsightFinder projects.
---

# insightfinder_log_labels (Resource)

Manages log filtering and labeling configuration for projects. Configure whitelists, blacklists, pattern naming, and training filters to optimize log analysis. Each entry in `label_settings` sets one label type on the project.

~> **Note:** The attribute on this resource is `label_settings`. The `log_label_settings` attribute belongs to `insightfinder_project`.

## Example Usage

### Whitelist Configuration

```terraform
resource "insightfinder_log_labels" "errors" {
  project_name = "application-logs"
  
  label_settings = [
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
  
  label_settings = [
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
  
  label_settings = [
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
  
  label_settings = [
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

### Simple String Labels

```terraform
resource "insightfinder_log_labels" "errors" {
  project_name = "application-logs"

  label_settings = [
    {
      label_type       = "whitelist"
      log_label_string = jsonencode(["ERROR", "FATAL"])
    }
  ]
}
```

### JSON Logs (key=value string rules)

For JSON-structured logs, prefix each value with the JSON key path it applies to:

```terraform
resource "insightfinder_log_labels" "json_logs" {
  project_name = "application-logs"

  label_settings = [
    {
      label_type       = "whitelist"
      log_label_string = jsonencode(["level=ERROR|FATAL"])
    },
    {
      label_type       = "blacklist"
      log_label_string = jsonencode(["message=.*healthcheck.*"])
    },
    {
      label_type       = "logSeverity"
      log_label_string = jsonencode(["severity=^(ERROR|CRITICAL)$"])
    }
  ]
}
```

## Schema

### Required

- `project_name` (String) The name of the project to configure log labels for. Changing this forces a new resource.
- `label_settings` (Attributes List) List of log label settings. Each entry is applied separately. (see [below for nested schema](#nestedatt--label_settings))

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

- `id` (String) Identifier for the log labels configuration (same as `project_name`).

<a id="nestedatt--label_settings"></a>
### Nested Schema for `label_settings`

Required:

- `label_type` (String) Type of log label. Values the provider recognizes (the API field each maps to is in parentheses):
  - `whitelist` (`whitelist`)
  - `trainingWhitelist` (`trainingWhitelist`)
  - `blacklist` (`trainingBlacklistLabels`)
  - `featurelist` (`featurelist`)
  - `incidentlist` (`incidentlist`)
  - `triagelist` (`triagelist`)
  - `anomalyFeature` (`anomalyFeatureLabels`)
  - `dataFilter` (`dataFilterLabels`)
  - `patternName` (`patternNameLabels`)
  - `patternSignature` (`patternSignatureLabels`)
  - `patternMatchRegex` (`patternMatchRegexLabels`)
  - `patternIgnoreRegex` (`patternIgnoreRegexLabels`)
  - `customAction` (`customActionLabels`)
  - `logEventID` (`logEventIDLabels`)
  - `logSeverity` (`logSeverityLabels`)
  - `logStatusCode` (`logStatusCodeLabels`)
  - `alertEventType` (`alertEventTypeLabels`)
  - `instanceName` (`instanceNameLabels`)
  - `dataQualityCheck` (`dataQualityCheckLabels`)
  - `incidentFieldVerification` (`incidentFieldVerificationLabels`)
  - `incidentPriority` (`incidentPriorityLabels`)
  - `extractionBlacklist` (`extractionBlacklist`)
  - `rareEventEscalationExclusion` (`rareEventEscalationExclusionLabels`)

  Any other value is sent to the API unchanged.
- `log_label_string` (String) JSON-encoded array of label rules. Use `jsonencode()`. See [Label Rule Schema](#label-rule-schema) for the rule object formats. Plain string arrays are also accepted:
  - Plain-text logs: values or regexes, e.g. `["ERROR","WARN"]` or `["^\\d+$"]`.
  - JSON logs: `key=value` rules, e.g. `["level=ERROR|WARN"]` or `["code=^\\d+$"]`.

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
- Empty `label_settings` will remove all labels from the project
- Several label types can be set in one resource. Each `label_type` should appear only once.
- On destroy, the provider sets each managed `label_type` to an empty array (`[]`). Label types this resource doesn't manage are left alone.
- If the project has none of the configured label types during refresh, the resource is removed from state.
