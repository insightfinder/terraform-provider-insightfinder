# Copyright (c) InsightFinder Inc.
# SPDX-License-Identifier: MPL-2.0

# Basic Slack integration
resource "insightfinder_slack" "basic" {
  system_name  = "Production"
  webhook      = var.slack_webhook
  channel_name = "#aiops-incidents"

  options = [
    "Detected Incident"
  ]
}

# Full configuration with per-project overrides
resource "insightfinder_slack" "full" {
  system_name  = "Production"
  webhook      = var.slack_webhook
  channel_name = "#aiops-incidents"

  options = [
    "Detected Incident",
    "Predicted Incident",
    "New Pattern Alert",
    "Missing Monitoring Data"
  ]

  # Notify a dedicated channel when an incident's priority is upgraded.
  priority_upgrade_channel = "#aiops-priority-upgrades"
  priority_upgrade_webhook = var.slack_webhook

  # Suppress Slack notifications for incidents that did not originate from InsightFinder.
  disable_slack_for_non_insightfinder_incidents = true

  # Send a subset of alert types for specific projects to a different channel,
  # optionally filtered by a match rule and priority levels.
  project_configs = [
    {
      project_name = "my-project"
      channel      = "my-project-incidents"
      options = [
        "Detected Incident",
        "Predicted Incident"
      ]
      enable_consolidation_info_update = true
      priority_levels                  = [1, 2, 3]
      rule = {
        type    = "fieldName"
        keyword = "alert->core->monitored_item=.*"
      }
    },
    {
      project_name = "another-project"
      channel      = "another-project-incidents"
      options = [
        "Detected Incident"
      ]
    }
  ]
}

variable "slack_webhook" {
  description = "Slack incoming webhook URL"
  type        = string
  sensitive   = true
}
