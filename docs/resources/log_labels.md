---
page_title: "insightfinder_log_labels Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages log filtering and labeling rules for InsightFinder projects.
---

# insightfinder_log_labels (Resource)

Manages log filtering and labeling configuration for projects. Each entry in `label_settings` sets one label type (whitelist, blacklist, pattern naming, severity, etc.) on the project.

~> **Note:** This standalone resource is different from the `log_label_settings` attribute of `insightfinder_project`. Here the attribute is `label_settings`, and each `log_label_string` must be a JSON array of **strings**, not objects.

## Example Usage

### Basic Configuration

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

### JSON Logs (key=value rules)

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

### With Project Dependency

```terraform
resource "insightfinder_project" "app" {
  project_name = "my-app-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }
}

resource "insightfinder_log_labels" "app_labels" {
  project_name = insightfinder_project.app.project_name

  label_settings = [
    {
      label_type       = "whitelist"
      log_label_string = jsonencode(["ERROR", "WARN"])
    }
  ]
}
```

## Schema

### Required

- `project_name` (String) The name of the project to configure log labels for. Changing this forces a new resource.
- `label_settings` (Attributes List) List of log label settings. Each entry is applied separately. (see [below for nested schema](#nestedatt--label_settings))

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
- `log_label_string` (String) JSON array of strings. The provider rejects anything that isn't a JSON array of strings. Use `jsonencode()`.
  - Plain-text logs: values or regexes, e.g. `["ERROR","WARN"]` or `["^\\d+$"]`.
  - JSON logs: `key=value` rules, e.g. `["level=ERROR|WARN"]` or `["code=^\\d+$"]`.

## Import

Log labels can be imported using the project name:

```shell
terraform import insightfinder_log_labels.example my-project-name
```

## Notes

- The project must exist before configuring log labels.
- Use `jsonencode()` so `log_label_string` is a valid JSON array.
- Values are case-sensitive, and regular expressions are supported.
- Several label types can be set in one resource. Each `label_type` should appear only once.
- On destroy, the provider sets each managed `label_type` to an empty array (`[]`). Label types this resource doesn't manage are left alone.
- If the project has none of the configured label types during refresh, the resource is removed from state.
- For richer, object-based label rules (`isCritical`, `isHotEventOnly`, etc.) use the `log_label_settings` attribute of `insightfinder_project`.
