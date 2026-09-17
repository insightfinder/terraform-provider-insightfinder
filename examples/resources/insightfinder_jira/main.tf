# Copyright (c) InsightFinder Inc.
# SPDX-License-Identifier: MPL-2.0

# service_host must end with a trailing slash: the backend appends "rest/api/..." to it
# via plain string concatenation, so a missing slash produces a broken host
# (e.g. "company.atlassian.netrest").
resource "insightfinder_jira" "basic" {
  account          = "user@company.com"
  service_host     = "https://company.atlassian.net/"
  api_token        = var.jira_api_token
  system_names     = ["Production"]
  jira_project_key = "II"
  jira_reporter_id = "5f9b2c1e4b0f8e001a2b3c4d"
}

# Full configuration with all optional fields
resource "insightfinder_jira" "full" {
  account      = "user@company.com"
  service_host = "https://company.atlassian.net/"
  api_token    = var.jira_api_token
  system_names = [
    "Production-US-East",
    "Production-US-West",
  ]
  jira_project_key = "II"
  jira_reporter_id = "5f9b2c1e4b0f8e001a2b3c4d"

  options        = ["Detected Incident", "Predicted Incident", "Root Cause"]
  content_option = ["SUMMARY", "RECOMMENDATION"]

  # Default assignee for created tickets
  jira_assignee_id = "5f9b2c1e4b0f8e001a2b3c4e"

  # Transition tickets to this status when the InsightFinder incident resolves
  jira_close_status_id = "31"

  # Extra fields to set at ticket-creation time
  jira_issue_fields = {
    "customfield_10010" = "some-value"
  }

  # Custom description template + placeholder-to-field mapping
  description_template = "Incident: {{title}}\nRoot cause: {{rootCause}}"
  template_field_mapping = {
    "title"     = "summary"
    "rootCause" = "customfield_10020"
  }

  # Post-creation field patches, gated by root-cause project/metric/component match
  field_update_rules = [
    {
      projects = ["AccessParks"]
      field_id = "customfield_10050"
      value    = "High"
      enabled  = true
    },
  ]
}

variable "jira_api_token" {
  description = "Jira API token generated for the account"
  type        = string
  sensitive   = true
}
