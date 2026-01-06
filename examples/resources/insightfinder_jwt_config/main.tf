# JWT Configuration Examples

## Basic JWT Configuration

```hcl
resource "insightfinder_jwt_config" "production" {
  system_name = "Production"
  jwt_secret  = var.jwt_secret
}
```

## JWT for Multiple Systems

```hcl
resource "insightfinder_jwt_config" "prod" {
  system_name = "Production"
  jwt_secret  = var.prod_jwt_secret
}

resource "insightfinder_jwt_config" "staging" {
  system_name = "Staging"
  jwt_secret  = var.staging_jwt_secret
}

resource "insightfinder_jwt_config" "dev" {
  system_name = "Development"
  jwt_secret  = var.dev_jwt_secret
}
```

## JWT with Explicit Type

```hcl
resource "insightfinder_jwt_config" "system_level" {
  system_name = "Production"
  jwt_secret  = var.jwt_secret
  jwt_type    = 1  # System-level JWT
}
```

## Variables

```hcl
variable "jwt_secret" {
  description = "JWT secret token (minimum 6 characters)"
  type        = string
  sensitive   = true
  
  validation {
    condition     = length(var.jwt_secret) >= 6
    error_message = "JWT secret must be at least 6 characters long."
  }
}

variable "prod_jwt_secret" {
  description = "Production JWT secret"
  type        = string
  sensitive   = true
}

variable "staging_jwt_secret" {
  description = "Staging JWT secret"
  type        = string
  sensitive   = true
}

variable "dev_jwt_secret" {
  description = "Development JWT secret"
  type        = string
  sensitive   = true
}
```

## Outputs

```hcl
output "jwt_system_id" {
  description = "The system ID for the JWT configuration"
  value       = insightfinder_jwt_config.production.id
}
```

## Notes

- JWT secrets must be at least 6 characters long
- Secrets are stored securely and marked as sensitive
- To delete JWT configuration, use `terraform destroy`
- The provider sends an empty string to the API to delete JWT configuration
