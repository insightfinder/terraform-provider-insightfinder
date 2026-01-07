# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-01-07

### Fixed
- **insightfinder_project**: Fixed state drift issue with log_label_settings
  - Added support for missing label types: `featurelist`, `incidentlist`, `triagelist`, `anomalyFeature`, `dataFilter`, `instanceName`, `dataQualityCheck`, and `extractionBlacklist`
  - Implemented order-preserving logic to maintain user-specified order of log_label_settings from configuration
  - Prevents false change detection when API returns labels in different order than specified in Terraform configuration
  - Resolves issue where `terraform plan` would continuously show changes even after `terraform apply`
  - Removed unused wrapper function to fix linter errors

## [1.0.1] - 2026-01-07

### Fixed
- **insightfinder_project**: Fixed state drift issue with log_label_settings
  - Added support for missing label types: `featurelist`, `incidentlist`, `triagelist`, `anomalyFeature`, `dataFilter`, `instanceName`, `dataQualityCheck`, and `extractionBlacklist`
  - Implemented order-preserving logic to maintain user-specified order of log_label_settings from configuration
  - Prevents false change detection when API returns labels in different order than specified in Terraform configuration
  - Resolves issue where `terraform plan` would continuously show changes even after `terraform apply`

## [1.0.0] - 2026-01-05

### Added

#### Resources
- **insightfinder_project**: Comprehensive project management with full configuration support
  - Project creation with data type, instance type, and cloud type configuration
  - Watch tower settings for anomaly detection and alerting
  - Email notifications and webhook configurations
  - Log processing settings and anomaly detection parameters
  - LLM evaluation settings for AI/ML projects
  - Incident prediction and root cause analysis configuration
  
- **insightfinder_servicenow**: ServiceNow integration management
  - Support for both OAuth and Basic authentication
  - System-level integration with multiple systems support
  - Configurable dampening periods and alert options
  - Content options for incident details
  - Automatic system name to ID resolution
  
- **insightfinder_jwt_config**: JWT authentication configuration
  - System-level JWT token management
  - Secure secret storage (marked as sensitive)
  - Automatic system name resolution
  - Deletion support via empty string configuration
  
- **insightfinder_log_labels**: Log filtering and labeling
  - Whitelist/blacklist configuration
  - Training whitelist for model optimization
  - Pattern naming rules
  - Field-based and regex-based filtering

#### Data Sources
- **insightfinder_project**: Query existing project configurations
  - Full project details retrieval
  - Watch tower settings access
  - Log label settings access
  
- **insightfinder_systems**: List and query available systems
  - System framework information
  - System display names and IDs
  - Owner information

#### Features
- **Provider Configuration**
  - Base URL configuration for different environments
  - Username and license key authentication
  - Environment variable support
  - Secure credential handling

- **Client Library**
  - Comprehensive API client with retry logic
  - System name to ID resolution helpers
  - JWT configuration management
  - ServiceNow integration API
  - Project management API
  - Log labels API

- **Developer Experience**
  - Development override support
  - Detailed debug logging with tflog
  - Comprehensive error messages
  - Import state support for all resources

### Technical Details

#### Dependencies
- Terraform Plugin Framework v1.4+
- Go 1.21+
- Terraform >= 1.0

#### Architecture
- Built with Terraform Plugin Framework (modern approach)
- Structured client library for API interactions
- Resource-specific validation and error handling
- State preservation for stable plan/apply cycles

#### Known Limitations
- ServiceNow integration requires system to exist before configuration
- JWT secrets must be at least 6 characters
- Log labels must be associated with existing projects
- Some project settings are computed and cannot be modified directly

### Security
- All sensitive fields (passwords, secrets, keys) marked as sensitive
- No credentials logged in debug output
- Secure HTTPS communication with InsightFinder API

## [Unreleased]

### Planned
- Terraform acceptance tests
- Additional data sources for metrics and logs
- Enhanced error messages with remediation hints
- Support for bulk operations
- Rate limiting and retry strategies

---

## Version History

- **1.0.0** - Initial release with core functionality

[1.0.0]: https://github.com/insightfinder/terraform-provider-insightfinder/releases/tag/v1.0.0
[Unreleased]: https://github.com/insightfinder/terraform-provider-insightfinder/compare/v1.0.0...HEAD
