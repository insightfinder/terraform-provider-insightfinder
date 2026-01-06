variable "username" {
  description = "InsightFinder username"
  type        = string
  sensitive   = true
}

variable "license_key" {
  description = "InsightFinder license key (API key)"
  type        = string
  sensitive   = true
}
