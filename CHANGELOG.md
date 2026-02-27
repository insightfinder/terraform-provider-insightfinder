# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.6.0] - 2026-02-26

### Added
- **insightfinder_project**: Extended `json_key_settings` with `metafield_setting` field
  - Added `metafield_setting` Boolean field to json_key_settings for metafield statistics management
  - Each JSON key now requires `json_key`, `type`, `summary_setting`, and `metafield_setting` fields
  - `metafield_setting`: Boolean flag to include the field in metafield statistics for enhanced log analysis
  - Both summary and metafield settings are collected separately during Create/Update operations
  - Both setting types are sent in a single API call with separate arrays for efficiency
  - Read operation fetches both settings from API and properly merges with configuration
  - Updated documentation and examples to demonstrate metafield_setting usage
  - Enables flexible statistics collection: track fields in summary stats, metafield stats, both, or neither

## [1.5.0] - 2026-02-20

### Added
- **insightfinder_project**: Added `json_key_settings` attribute for managing custom JSON fields in logs
  - New optional nested list to define custom JSON fields extracted from logs
  - Each JSON key requires `json_key`, `type`, and `summary_setting` fields
  - `json_key`: The JSON field name to extract from logs
  - `type`: The data type (e.g., `string`, `number`, `JSONArray`)
  - `summary_setting`: Boolean flag to include the field in summary statistics
  - Supports full CRUD operations via two separate endpoints:
    - `GET/POST /api/external/v1/logjsontype` for JSON key type management
    - `GET/POST /api/external/v1/logsummarysettings` for summary setting management
  - API response order is not relied upon; preserves configuration order in state
  - JSON keys are tracked by name to maintain consistent ordering across API calls
  - Updated documentation and examples to demonstrate JSON key configuration
  - Helps extract and analyze custom fields from structured JSON logs

## [1.4.1] - 2026-02-09

### Fixed
- **insightfinder_project**: Fixed holiday_settings order inconsistency issue
  - Holiday settings now preserve the order defined in configuration
  - Resolved "Provider produced inconsistent result after apply" errors
  - Read operation now maintains existing holiday order from state when refreshing
  - Prevents unnecessary plan changes due to reordering

## [1.4.0] - 2026-02-09

### Added
- **insightfinder_project**: Added `holiday_settings` attribute for managing project holidays
  - New optional nested list to define holiday periods for anomaly detection
  - Each holiday requires `name`, `start_date` (MM-DD format), and `end_date` (MM-DD format)
  - Holidays are automatically created, updated, and deleted to match configuration
  - Supports full CRUD operations via `/api/external/v1/holiday` endpoint
  - Holidays are sorted by name to maintain consistent ordering
  - Updated documentation and examples to demonstrate holiday configuration
  - Helps improve anomaly detection accuracy during expected holiday periods

## [1.3.2] - 2026-02-06

### Added
- **insightfinder_project**: Added `servicenow_table` field to `project_creation_config` for ServiceNow projects
  - New optional field in project creation configuration
  - When `project_cloud_type` is "ServiceNow", the `servicenow_table` field can be specified to set the ServiceNow table name
  - The field is passed as `tableName` parameter to the `/api/v1/check-and-add-custom-project` API endpoint
  - Backward compatible: only included in API request when project type is ServiceNow and the field has a value
  - Updated documentation and examples to demonstrate ServiceNow project creation with table configuration

## [1.3.1] - 2026-02-04

### Fixed
- **insightfinder_project**: Fixed panic when creating projects with nil JSON array fields
  - Added proper nil checks for `cdf_setting`, `log_to_log_setting_list`, `webhook_header_list`, and `shared_usernames` fields
  - Replaced unsafe type assertions with safe type checking using comma-ok idiom
  - Prevents "interface conversion: interface {} is nil, not []interface {}" panic errors
  - Resolves issue where projects would fail to create when these fields parsed to nil values

## [1.3.0] - 2026-02-03

### Added
- **insightfinder_project**: Added `project_servicenow_settings` attribute for ServiceNow integration
  - New optional nested object to configure ServiceNow third-party settings
  - Only applies when `project_cloud_type` is "ServiceNow" (case-insensitive)
  - Supports fields: `host`, `sysparm_query`, `proxy`, `servicenow_user`, `servicenow_password`, `instance_field`, `instance_field_regex`, `timestamp_format`, `client_id`, `client_secret`, and `additional_fields`
  - Automatically creates/updates ServiceNow integration via `/api/external/v1/thirdpartysetting` endpoint
  - Sensitive fields (`servicenow_password`, `client_secret`) are properly marked and protected

### Fixed
- **insightfinder_project**: Fixed log label drift detection and management
  - Implemented order-independent comparison for `log_label_settings` to prevent false drift detection when labels are reordered
  - Added automatic deletion of labels that exist in API but not in configuration
  - Labels added via UI are now properly detected as drift and can be removed via `terraform apply`
  - Fixed `DeleteLogLabels` function to properly send requests without double-marshaling JSON
  - Preserves original order in state when no actual changes occur
  
- **insightfinder_project**: Fixed "unknown value" error for non-ServiceNow projects
  - Explicitly set `project_servicenow_settings` to null for non-ServiceNow projects
  - Resolves "Provider returned invalid result object after apply" error
  - Ensures all computed values are known after create/update operations

## [1.2.1] - 2026-01-09

### Removed
- **insightfinder_project**: Completely removed `project_creation_type` attribute
  - Removed from resource schema definition
  - Removed from internal structs and client code
  - Removed from all test configurations
  - This field was not needed for project creation and has been eliminated to simplify the API

## [1.2.0] - 2026-01-07

### Fixed
- **insightfinder_project**: Fixed handling of optional `log_label_settings` attribute
  - Changed internal representation from Go slice to `types.List` to properly handle unknown values during plan phase
  - Projects can now be created without specifying `log_label_settings` attribute
  - Resolves "Value Conversion Error" when `log_label_settings` is not provided in configuration
  - Added proper null/unknown value checking before processing log label settings
  - Improved compatibility with Terraform Framework's type system

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
