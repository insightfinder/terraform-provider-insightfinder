# ServiceNow Integration Examples

## Basic Authentication

```hcl
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
```

## OAuth Authentication

```hcl
resource "insightfinder_servicenow" "oauth" {
  account          = "admin"
  service_host     = "https://dev12345.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 7200000
  app_id           = var.servicenow_app_id
  app_key          = var.servicenow_app_key
  system_names     = ["Production", "Staging"]
  options          = ["Root Cause"]
  content_option   = ["SUMMARY", "DESCRIPTION"]
  auth_type        = "oauth"
}
```

## Multiple Systems

```hcl
resource "insightfinder_servicenow" "multi_system" {
  account          = "serviceaccount"
  service_host     = "https://company.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 3600000
  system_names     = [
    "Production-US-East",
    "Production-US-West",
    "Production-EU"
  ]
  options        = ["Root Cause", "Prediction"]
  content_option = ["SUMMARY", "DESCRIPTION", "IMPACT"]
  auth_type      = "basic"
  
  # Optional proxy
  proxy = "http://proxy.company.com:8080"
}
```

## Variables

```hcl
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
```
