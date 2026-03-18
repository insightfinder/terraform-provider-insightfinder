# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.8.1] - 2026-03-17

### Fixed
- **insightfinder_project**: Changed `json_key_settings`, `holiday_settings`, and `log_label_settings` from lists to sets so that order changes in the API response no longer produce a diff in the Terraform plan

## [1.8.0] - 2026-03-16

### Added
- **insightfinder_metric_project**: New resource for managing InsightFinder metric projects
  - Full CRUD support for metric projects via `/api/external/v1/watch-tower-setting`
  - Project creation via `/api/v1/check-and-add-custom-project` with `data_type`, `instance_type`, `project_cloud_type`, and `insight_agent_type`
  - All metric-specific detection tuning fields: `c_value`, `p_value`, `high_ratio_c_value`, `maximum_hint`, `dynamic_baseline_detection_flag`, `baseline_duration`, `positive_baseline_violation_factor`, `negative_baseline_violation_factor`, `enable_period_anomaly_filter`, `enable_ubl_detect`, `enable_cumulative_detect`, `enable_baseline_detection_double_verify`, `filter_by_anomaly_in_baseline_generation`, `anomaly_dampening`, `anomaly_gap_tolerance_count`
  - Gap filling and prediction fields: `enable_fill_gap`, `enable_store_filled_gap`, `gap_filling_training_data_length`, `enable_metric_data_prediction`, `prediction_training_data_length`, `prediction_correlation_sensitivity`, `enable_kpi_prediction`
  - Incident prediction and RCA fields: `incident_prediction_window`, `min_incident_prediction_window`, `incident_relation_search_window`, `incident_prediction_event_limit`, `root_cause_count_threshold`, `root_cause_probability_threshold`, `causal_prediction_setting`, `root_cause_rank_setting`, `maximum_root_cause_result_size`, `multi_hop_search_level`, `multi_hop_search_limit`
  - Instance down detection: `instance_down_threshold`, `instance_down_report_number`, `instance_down_enable`, `instance_down_ratio_threshold`, `show_instance_down`
  - Holiday settings management (create/update/delete via `/api/external/v1/holiday`)
  - Complex JSON fields: `linked_log_projects`, `component_metric_setting_overall_model_list`, `email_setting`, `instance_grouping_update`, `shared_usernames`, `webhook_header_list`
  - Full webhook configuration support
  - Import support via project name
  - **`metric_configurations`** block: Per-metric alert threshold and component operation settings
    - Each entry targets a named metric (`metric_name`) and supports:
    - `escalate_incident_components` (List of String): component names that escalate incidents for this metric
    - `ignored_components` (List of String): component names excluded from detection for this metric
    - `metric_alert_settings` (List of Objects): per-component (or global) alert threshold rows with full threshold bands (`threshold_alert_lower_bound`, `threshold_alert_upper_bound`, `incident_alert_lower_bound`, `incident_alert_upper_bound`, and their negative variants), detection flags (`is_kpi`, `is_flapping_result_only`, `fill_zero`, `compute_difference`, `enable_baseline_near_constance`), and display fields (`detection_type`, `pattern_name_higher`, `pattern_name_lower`, `metric_type`, `rouge_value`, `incident_duration_threshold`, `anomaly_gap_tolerance_duration`)
    - API endpoints: GET/POST `/api/external/v1/componentmetricupdate` for alert settings; GET/POST `/api/external/v1/metriccomponent` for escalate/ignore component operations

- **insightfinder_system_settings**: Extended `notifications_settings` with dedicated notification sub-blocks
  - **`system_down_notification`** block: Configures system-down email alerts via `/api/external/v2/systemdownsetting`
    - `enable_system_down_email_alert` (Boolean): Enable email when the system goes down
    - `email_dampening_period` (Number, ms): Minimum interval between repeated system-down emails
    - `email_set` (List of String): Recipient addresses for system-down notifications
  - **`daily_report_notification`** block: Configures daily insights report emails via `/api/external/v1/insightsreportsetting`
    - `enable_insights_report` (Boolean): Enable the daily summary email
    - `email_set` (List of String): Recipient addresses for the daily report
  - **`weekly_report_notification`** block: Configures weekly insights report emails (same API, `isDaily=false`)
    - `enable_insights_report` (Boolean): Enable the weekly summary email
    - `email_set` (List of String): Recipient addresses for the weekly report
  - **`instance_down_notification`** block (List): Per-project instance-down alert settings via `/api/external/v1/projects/update`
    - `project_name` (String, Required): The project to configure
    - `instance_down_enable` (Boolean): Enable instance-down detection for this project
    - `instance_down_dampening` (Number, ms): Dampening window between repeated instance-down alerts
    - `instance_down_threshold` (Number, ms): Duration before an instance is considered down
    - `instance_down_report_number` (Number): Number of down instances before an alert is sent
    - `instance_down_emails` (List of String): Recipient addresses for instance-down notifications

### Fixed
- **insightfinder_metric_project**: Fixed `false` boolean and `0` integer values being silently dropped from API requests
  - `populateMetricSettings` now builds `map[string]interface{}` directly instead of marshaling through a struct with `omitempty` tags
  - `UpdateMetricProject` now sends the settings map directly without double-marshaling through `MetricProjectSettings`
  - Affected fields included `enable_kpi_prediction`, `filter_by_anomaly_in_baseline_generation`, `show_instance_down`, `gap_filling_training_data_length`, `prediction_training_data_length`, and all other boolean/integer fields set to their zero values

- **insightfinder_project**: Fixed same `omitempty` bug where `false` booleans and `0` integers were dropped from API update requests
  - `populateSettings` now builds `map[string]interface{}` directly
  - `UpdateProject` now sends the settings map directly without the unnecessary double-marshal through `ProjectSettings`
  - Resolves issues where settings like `training_filter = false`, `enable_hot_event = false`, or any integer field set to `0` would not be applied

## [1.7.1] - 2026-03-12

### Changed
- **insightfinder_servicenow**: Extended resource with new fields and updated API integration
  - Added `service_now_field` (String, Optional): ServiceNow field to write integration content to (e.g., `u_probable_cause`)
  - Added `content_source` (String, Optional, Computed): ServiceNow content source field (e.g., `work_notes`). Defaults to `work_notes`
  - Added `trigger_window_in_mills` (Number, Optional): Trigger window in milliseconds (e.g., `604800000` for 7 days)
  - Added `table_mapping` (Map of String, Optional): Mapping of InsightFinder project names to ServiceNow table names
  - Changed `options` and `content_option` from `List` to `Set` type to avoid order-sensitivity drift
  - Removed `system_ids` computed attribute — system resolution is now handled internally
  - Updated `GetServiceNowConfig` to use the new `/api/external/v1/system/externalServlies/list` endpoint and match entries by `account` + `service_host`
  - Added `UpdateServiceNowTableMapping` client method using a dedicated PUT endpoint
  - `contentSource` and `serviceNowField` are now sent in the Create/Update API call
  - Auth type is inferred from `app_id`/`app_key` presence during Read

## [1.7.0] - 2026-03-11

### Added
- **insightfinder_project**: Extended `json_key_settings` with `dampening_field_setting` field
  - Added `dampening_field_setting` Boolean field to `json_key_settings` for dampening field management
  - Each JSON key now requires `json_key`, `type`, `summary_setting`, `metafield_setting`, and `dampening_field_setting` fields
  - `dampening_field_setting`: Boolean flag to include the field in the dampening field list, controlling which fields are used for alert dampening logic
  - Dampening field keys are collected alongside summary and metafield keys during Create/Update operations
  - All three setting types (`summarySetting`, `metaFieldSetting`, `dampeningFieldSetting`) are sent in a single POST to `/api/external/v1/logsummarysettings`
  - Read operation fetches `dampeningFieldSetting` array from API response and maps it back to state
  - Updated documentation and examples to demonstrate `dampening_field_setting` usage

- **insightfinder_system_settings**: Added documentation for `knowledgebase_settings` and `notifications_settings`
  - Added `docs/resources/system_settings.md` with full schema reference for both nested blocks
  - `knowledgebase_settings`: covers global KB fields (`enable_global_knowledge_base`, `satellite_system_set`, `composite_valid_threshold`, `timeline_top_k`, etc.) and incident prediction fields (`rule_active_threshold`, `rule_inactive_threshold`, `kb_training_length`, `tolerance`, etc.)
  - `notifications_settings`: covers health view display, alert thresholds, email dampening periods, per-event-type email toggles, recipient addresses, and the JSON-encoded map fields (`incident_count_threshold`, `assignment_map`)
  - Includes example usage blocks for full configuration, KB-only, and satellite system linking
  - Documents import behavior (`terraform import insightfinder_system_settings.example <system_name>`) and delete-only semantics

## [1.6.1] - 2026-02-27

### Fixed
- **insightfinder_project**: Fixed missing log label type mappings
  - Added support for `incidentFieldVerification` label type
  - Added support for `incidentPriority` label type
  - Updated label type to API field name mappings to include all label types
  - Fixes issue where incident field verification and incident priority labels were not properly handled

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

- **1.8.0** - Added `insightfinder_metric_project` resource; fixed `omitempty` zero-value bug in both metric and log project updates
- **1.7.0** - Added `dampening_field_setting` to `json_key_settings`; fixed system settings and email_setting generation in CLI tool
- **1.0.0** - Initial release with core functionality

[1.8.0]: https://github.com/insightfinder/terraform-provider-insightfinder/releases/tag/v1.8.0
[1.7.0]: https://github.com/insightfinder/terraform-provider-insightfinder/releases/tag/v1.7.0
[1.0.0]: https://github.com/insightfinder/terraform-provider-insightfinder/releases/tag/v1.0.0
[Unreleased]: https://github.com/insightfinder/terraform-provider-insightfinder/compare/v1.8.0...HEAD
