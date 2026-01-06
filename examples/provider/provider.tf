# Provider Configuration Example

This example demonstrates how to configure the InsightFinder provider.

## Basic Configuration

```hcl
terraform {
  required_providers {
    insightfinder = {
      source  = "insightfinder/insightfinder"
      version = "~> 1.0"
    }
  }
}

provider "insightfinder" {
  base_url    = "https://app.insightfinder.com"
  username    = var.username
  license_key = var.license_key
}
```

## Using Environment Variables

```bash
export INSIGHTFINDER_USERNAME="your-username"
export INSIGHTFINDER_LICENSE_KEY="your-license-key"
export INSIGHTFINDER_BASE_URL="https://app.insightfinder.com"
```

```hcl
provider "insightfinder" {
  # Configuration will be read from environment variables
}
```

## Variables

```hcl
variable "username" {
  description = "InsightFinder username"
  type        = string
  sensitive   = true
}

variable "license_key" {
  description = "InsightFinder license key"
  type        = string
  sensitive   = true
}
```

## For Staging Environment

```hcl
provider "insightfinder" {
  base_url    = "https://stg.insightfinder.com"
  username    = var.username
  license_key = var.license_key
}
```
 required_providers {
    insightfinder = {
      source = "insightfinder/insightfinder"
    }
  }
}

provider "insightfinder" {
  base_url    = "https://app.insightfinder.com"
  username    = var.username
  license_key = var.license_key
}

# Create a project
resource "insightfinder_project" "example" {
  project_name         = "my-production-project"
  project_display_name = "Production Monitoring"
  system_name          = "production-system"
  
  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }
  
  c_value             = 3
  p_value             = 0.9
  project_time_zone   = "UTC"
  sampling_interval   = 600
}

# Output the project details
output "project_id" {
  value = insightfinder_project.example.id
}

output "project_name" {
  value = insightfinder_project.example.project_name
}
