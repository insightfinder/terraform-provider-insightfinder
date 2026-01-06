---
page_title: "InsightFinder Provider"
subcategory: ""
description: |-
  The InsightFinder provider allows you to manage InsightFinder resources using Terraform.
---

# InsightFinder Provider

The InsightFinder provider allows you to manage InsightFinder resources using Terraform. Use it to create and configure projects, integrate with ServiceNow, manage JWT authentication, and configure log processing rules.

## Example Usage

```terraform
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

## Authentication

The provider supports authentication via:

1. **Provider configuration** (shown above)
2. **Environment variables**:
   - `INSIGHTFINDER_USERNAME`
   - `INSIGHTFINDER_LICENSE_KEY`
   - `INSIGHTFINDER_BASE_URL`

## Schema

### Required

- `username` (String, Sensitive) InsightFinder username
- `license_key` (String, Sensitive) InsightFinder license key

### Optional

- `base_url` (String) InsightFinder API base URL. Defaults to `https://app.insightfinder.com`
