terraform {
  required_providers {
    insightfinder = {
      source = "insightfinder/insightfinder"
    }
  }
}

provider "insightfinder" {
  # Configuration options
  license_key = var.license_key
  user_name   = var.user_name
  server_url  = var.server_url # e.g., "https://stg.insightfinder.com"
}

resource "insightfinder_project" "example_with_holidays" {
  project_name         = "example-project-with-holidays"
  project_display_name = "Example Project with Holidays"
  system_name          = "example-system"

  project_creation_config {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
  }

  # Basic anomaly detection settings
  c_value           = 3
  p_value           = 0.05
  sampling_interval = 300

  # Holiday settings - list of holidays with name, start date, and end date
  holiday_settings = [
    {
      name       = "christmas"
      start_date = "12-25" # December 25th
      end_date   = "12-26" # December 26th
    },
    {
      name       = "new_year"
      start_date = "01-01" # January 1st
      end_date   = "01-01" # January 1st
    },
    {
      name       = "independence_day"
      start_date = "07-04" # July 4th
      end_date   = "07-04" # July 4th
    },
    {
      name       = "thanksgiving"
      start_date = "11-28" # Example: November 28th
      end_date   = "11-29" # Through November 29th
    }
  ]
}

# Variables
variable "license_key" {
  description = "InsightFinder license key"
  type        = string
  sensitive   = true
}

variable "user_name" {
  description = "InsightFinder username"
  type        = string
}

variable "server_url" {
  description = "InsightFinder server URL"
  type        = string
  default     = "https://app.insightfinder.com"
}

# Outputs
output "project_id" {
  description = "The ID of the created project"
  value       = insightfinder_project.example_with_holidays.id
}

output "project_name" {
  description = "The name of the created project"
  value       = insightfinder_project.example_with_holidays.project_name
}

output "holiday_settings" {
  description = "The holiday settings configured for the project"
  value       = insightfinder_project.example_with_holidays.holiday_settings
}
