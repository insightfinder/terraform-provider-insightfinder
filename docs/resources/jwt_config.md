---
page_title: "insightfinder_jwt_config Resource - terraform-provider-insightfinder"
subcategory: ""
description: |-
  Manages JWT authentication configuration for InsightFinder systems.
---

# insightfinder_jwt_config (Resource)

Manages system-level JWT authentication configuration. JWT tokens are used for secure API authentication at the system level.

## Example Usage

### Basic Configuration

```terraform
resource "insightfinder_jwt_config" "production" {
  system_name = "Production"
  jwt_secret  = var.jwt_secret
}
```

### Multiple Systems

```terraform
resource "insightfinder_jwt_config" "prod" {
  system_name = "Production"
  jwt_secret  = var.prod_jwt_secret
}

resource "insightfinder_jwt_config" "staging" {
  system_name = "Staging"
  jwt_secret  = var.staging_jwt_secret
}
```

### With Validation

```terraform
variable "jwt_secret" {
  description = "JWT secret token"
  type        = string
  sensitive   = true
  
  validation {
    condition     = length(var.jwt_secret) >= 6
    error_message = "JWT secret must be at least 6 characters long."
  }
}

resource "insightfinder_jwt_config" "system" {
  system_name = "Production"
  jwt_secret  = var.jwt_secret
  jwt_type    = 1
}
```

## Schema

### Required

- `system_name` (String) Name of the system to configure JWT for
- `jwt_secret` (String, Sensitive) JWT secret token (minimum 6 characters)

### Optional

- `jwt_type` (Number) JWT type. Default: `1` (system-level JWT)

### Read-Only

- `id` (String) JWT configuration identifier (same as system_name)

## Import

JWT configurations can be imported using the system name:

```shell
terraform import insightfinder_jwt_config.example Production
```

## Notes

- JWT secrets must be at least 6 characters long
- The `jwt_secret` field is marked as sensitive and will not appear in logs or console output
- To delete JWT configuration, remove the resource from your configuration and run `terraform apply`
- The provider sends an empty string to the API to delete the JWT configuration
- System names are automatically resolved to system IDs
