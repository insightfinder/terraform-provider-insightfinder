# Copyright (c) InsightFinder Inc.
# SPDX-License-Identifier: MPL-2.0

# Basic Authentication
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

# OAuth Authentication
resource "insightfinder_servicenow" "oauth" {
  account          = "admin"
  service_host     = "https://dev12345.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 7200000
  app_id           = var.servicenow_app_id
  app_key          = var.servicenow_app_key
  system_names     = ["Production", "Staging"]
  options          = ["Root Cause", "Prediction"]
  content_option   = ["SUMMARY", "DESCRIPTION"]
  auth_type        = "oauth"
}

# Full configuration with all optional fields
resource "insightfinder_servicenow" "full" {
  account          = "serviceaccount"
  service_host     = "https://company.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 3600000
  system_names = [
    "Production-US-East",
    "Production-US-West",
    "Production-EU"
  ]
  options        = ["Root Cause", "Prediction"]
  content_option = ["SUMMARY", "DESCRIPTION", "IMPACT"]
  auth_type      = "basic"

  # Optional proxy
  proxy = "http://proxy.company.com:8080"

  # Write content to a custom ServiceNow field
  service_now_field = "u_probable_cause"

  # Use comments instead of work notes
  content_source = "comments"

  # Only correlate events within a 7-day window
  trigger_window_in_mills = 604800000

  # Enable feedback collection from ServiceNow
  enable_feedback_collect = false

  # Enable ticket creation in ServiceNow
  enable_ticket_creation = false

  # Enable ticket update in ServiceNow
  enable_ticket_update = false

  # Filter ticket creation by a specific field value
  ticket_created_by_source_key   = "activity_due"
  ticket_created_by_source_value = "xyzzz"

  # Associate created tickets with a CMDB configuration item
  configuration_item = "My-Server-CI"

  # Map InsightFinder projects to ServiceNow tables
  table_mapping = {
    "my-project"      = "incident"
    "another-project" = "problem"
  }
}

variable "servicenow_password" {
  description = "ServiceNow account password"
  type        = string
  sensitive   = true
}

variable "servicenow_app_id" {
  description = "ServiceNow OAuth application ID"
  type        = string
  sensitive   = true
}

variable "servicenow_app_key" {
  description = "ServiceNow OAuth application key"
  type        = string
  sensitive   = true
}
