# Terraform Provider for InsightFinder

[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)

The Terraform InsightFinder provider allows you to manage InsightFinder resources using Infrastructure as Code (IaC).

## Features

- **Project Management**: Create and manage InsightFinder projects with comprehensive configuration options
- **ServiceNow Integration**: Configure ServiceNow integrations with OAuth/Basic authentication
- **JWT Configuration**: Manage system-level JWT authentication tokens
- **Log Labels**: Configure log filtering, whitelisting, and pattern naming rules
- **Data Sources**: Query existing projects and systems

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)

## Installation

### Terraform Registry (Recommended)

Add the provider to your Terraform configuration:

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
  base_url    = "https://app.insightfinder.com"  # or your instance URL
  username    = var.username
  license_key = var.license_key
}
```

Then run:
```bash
terraform init
```

### Local Development

For local development, build and install the provider:

```bash
make install
```

This will compile the provider and install it to `~/.terraform.d/plugins/`.

## Usage

### Provider Configuration

```hcl
provider "insightfinder" {
  base_url    = "https://app.insightfinder.com"
  username    = "your-username"
  license_key = "your-license-key"
}
```

#### Authentication

The provider supports the following authentication methods:

1. **Static credentials** (shown above)
2. **Environment variables**:
   ```bash
   export INSIGHTFINDER_USERNAME="your-username"
   export INSIGHTFINDER_LICENSE_KEY="your-license-key"
   export INSIGHTFINDER_BASE_URL="https://app.insightfinder.com"
   ```

### Example: Creating a Project

```hcl
resource "insightfinder_project" "example" {
  project_name = "my-application-logs"
  system_name  = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "LogStreaming"
  }

  project_display_name = "Application Logs"
  project_time_zone    = "America/New_York"
  sampling_interval    = 600
}
```

### Example: ServiceNow Integration

```hcl
resource "insightfinder_servicenow" "incident_management" {
  account          = "admin"
  service_host     = "https://dev12345.service-now.com/"
  password         = var.servicenow_password
  dampening_period = 3600000
  system_names     = ["Production", "Staging"]
  options          = ["Root Cause"]
  content_option   = ["SUMMARY"]
  auth_type        = "basic"
}
```

### Example: JWT Configuration

```hcl
resource "insightfinder_jwt_config" "production" {
  system_name = "Production"
  jwt_secret  = var.jwt_secret
}
```

### Example: Log Labels

```hcl
resource "insightfinder_log_labels" "filters" {
  project_name = insightfinder_project.example.project_name
  
  log_label_settings = [
    {
      label_type       = "whitelist"
      log_label_string = jsonencode([{
        type           = "fieldName"
        keyword        = "severity=error|critical"
        isCritical     = true
        isHotEventOnly = false
      }])
    },
    {
      label_type       = "patternName"
      log_label_string = jsonencode([{
        type           = "fieldName"
        keyword        = "message"
        patternNameKey = "message"
      }])
    }
  ]
}
```

### Data Sources

Query existing resources:

```hcl
data "insightfinder_project" "existing" {
  project_name = "my-existing-project"
}

data "insightfinder_systems" "all" {}

output "system_names" {
  value = data.insightfinder_systems.all.systems[*].system_name
}
```

## Resources

- [`insightfinder_project`](docs/resources/project.md) - Manage InsightFinder projects
- [`insightfinder_servicenow`](docs/resources/servicenow.md) - Configure ServiceNow integrations
- [`insightfinder_jwt_config`](docs/resources/jwt_config.md) - Manage JWT authentication
- [`insightfinder_log_labels`](docs/resources/log_labels.md) - Configure log filtering and labeling

## Data Sources

- [`insightfinder_project`](docs/data-sources/project.md) - Query project information
- [`insightfinder_systems`](docs/data-sources/systems.md) - List available systems

## Development

### Building the Provider

```bash
go build -o terraform-provider-insightfinder
```

### Running Tests

```bash
go test ./...
```

### Installing Locally

```bash
make install
```

This installs the provider to:
```
~/.terraform.d/plugins/registry.terraform.io/insightfinder/insightfinder/1.0.0/
```

### Using Development Override

Create a `.terraformrc` file in your project:

```hcl
provider_installation {
  dev_overrides {
    "insightfinder/insightfinder" = "/path/to/terraform-provider-insightfinder"
  }
  direct {}
}
```

Then set:
```bash
export TF_CLI_CONFIG_FILE=.terraformrc
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linters
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

For detailed development and testing instructions, see:
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [TESTING.md](TESTING.md) - Testing guide
- [TEST_SUMMARY.md](TEST_SUMMARY.md) - Test suite overview

### Testing

The provider includes a comprehensive test suite:

```bash
# Run unit tests
make test-unit

# Run acceptance tests (requires credentials)
export IF_USERNAME="your-username"
export IF_LICENSE_KEY="your-license-key"
make test-acc

# Run specific resource tests
make test-project
make test-servicenow
make test-jwt

# Generate coverage report
make test-coverage
```

**Test Coverage:**
- 10 unit tests (provider, client)
- 27 acceptance tests (resources, data sources)
- Full CRUD operation testing
- Import state verification
- Error handling validation

See [TESTING.md](TESTING.md) for complete testing documentation.

## Support

For issues and questions:
- Open an [issue](https://github.com/insightfinder/terraform-provider-insightfinder/issues)
- Contact InsightFinder support at [support@insightfinder.com](mailto:support@insightfinder.com)

## License

This provider is released under the Mozilla Public License 2.0. See [LICENSE](LICENSE) for details.

## Authors

Maintained by [InsightFinder](https://www.insightfinder.com/)

## Acknowledgments

Built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
