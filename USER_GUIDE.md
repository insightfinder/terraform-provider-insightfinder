# InsightFinder Terraform Provider — User Guide

This document is a comprehensive reference for every configurable item exposed by the
InsightFinder Terraform Provider. It covers the provider block, all seven resources,
and both data sources.

**How to read the tables**

- **Type** is the Terraform type of the argument.
- **Required / Optional / Computed** follow Terraform semantics. Almost every optional
  argument on the project resources is `Optional + Computed`: if you omit it, Terraform
  does **not** send a value and the InsightFinder platform applies its own server-side
  default, which is then recorded in state. The "Default" values listed below are
  **platform defaults** unless the row explicitly says *provider default*, which means
  the provider itself sends that value when the argument is omitted.
- Arguments marked **Sensitive** are redacted from Terraform plan/apply output.

---

## Table of Contents

1. [Settings Map — by Function](#settings-map--by-function)
2. [Provider Configuration](#provider-configuration)
3. [Resources](#resources)
   - [insightfinder_project](#insightfinder_project)
   - [insightfinder_metric_project](#insightfinder_metric_project)
   - [insightfinder_log_labels](#insightfinder_log_labels)
   - [insightfinder_jwt_config](#insightfinder_jwt_config)
   - [insightfinder_servicenow](#insightfinder_servicenow)
   - [insightfinder_slack](#insightfinder_slack)
   - [insightfinder_system_settings](#insightfinder_system_settings)
3. [Data Sources](#data-sources)
   - [insightfinder_project (data source)](#insightfinder_project-data-source)
   - [insightfinder_systems](#insightfinder_systems)
5. [Appendix A — Import ID Reference](#appendix-a--import-id-reference)
6. [Appendix B — Force-Replacement Arguments](#appendix-b--force-replacement-arguments)

---

## Settings Map — by Function

The rest of this guide is organized by Terraform resource. This section organizes the same
settings the way you encounter them in InsightFinder itself — by what they *do* rather than
by which resource declares them. Each entry links to the reference section that documents it.

Where a functional area spans more than one resource, the owning resource is named in
parentheses.

### 1. System Settings

System-wide configuration that applies to a whole system rather than an individual project.
All of it lives in [`insightfinder_system_settings`](#insightfinder_system_settings), which
targets a system by display name.

1. **[Notification settings](#notifications_settings-block)** — everything that decides who
   gets told about an incident, through which channel, and how often.
   - **[Health view & general fields](#health-view--general-fields)** — dashboard placement
     and ordering, aggregation interval, Splunk export, per-project incident count
     thresholds, assignee maps, alert health score, and the system-level incident dampening
     window.
   - **[Email alert fields](#email-alert-fields)** — recipient addresses and enable/disable
     toggles for prediction, detection, health, general alert, and root-cause emails, plus
     their individual dampening periods.
   - **Consolidation settings** — how separate incidents get merged into one.
     - **[Toggles and algorithms](#health-view--general-fields)** —
       `component_level_incident_consolidation`, `component_level_dampening`,
       `enabled_consolidation_algorithms` (`derivedIncidents`, `rcaChain`, `contentBased`,
       `metricInstanceTimestamp`), and `metric_co_occurrence_buffer_ms`.
     - **[Custom consolidation rules](#custom_consolidation_rules-block)** — hand-written
       rules that group incidents across named projects by matching field values or content,
       with field correlations mapping equivalent fields between those projects.
     - **[Metric-to-log consolidation](#metric_log_consolidation_configs-block)** — pair a
       metric project with a log project and list the field keys used to correlate them.
   - **[Project-level dampening windows](#project_level_dampening_windows-block)** — override
     the system-level dampening window for a specific source → target project pair.
   - **[System-down notification](#system_down_notification-block)** — email alerting when
     the system as a whole stops reporting. Managed through its own API.
   - **[Instance-down notification](#instance_down_notification-block)** — per-project
     alerting when individual instances stop reporting, with their own thresholds,
     dampening, and recipient lists.
   - **Insights reports** — scheduled summary emails and their recipient lists, configured
     separately for the [daily](#daily_report_notification-block) and
     [weekly](#weekly_report_notification-block) report.
2. **[Knowledgebase settings](#knowledgebase_settings-block)** — the Knowledge Base and
   incident prediction engine.
   - **[Global KB fields](#global-kb-fields)** — enable the global knowledge base, composite
     thresholds, timeline retention, prediction source, auto-fix validation window, and the
     satellite systems linked to this system's KB.
   - **[Incident prediction fields](#incident-prediction-fields)** — the promotion and
     retirement thresholds that decide when a learned causal rule becomes an active alerting
     rule, false-positive tolerance, and the KB training window.
3. **[Miscellaneous settings](#miscellaneous_settings-block)** — long-term health view
   storage, automatic system sharing, the root cause reverse entry filter threshold, and the
   composite timeline view.

### 2. General Project Settings

Settings that apply to any project regardless of data type. Unless noted, these exist on both
[`insightfinder_project`](#insightfinder_project) and
[`insightfinder_metric_project`](#insightfinder_metric_project) and behave identically —
see [Common Optional Arguments](#common-optional-arguments-shared-with-insightfinder_project)
for the metric project's copy of the list.

- **[Identity and retention](#optional-arguments--general-settings)** — display name,
  timezone, sampling interval, data and UBL retention, and the C/P anomaly sensitivity pair.
- **[Anomaly detection and alerts](#optional-arguments--anomaly-detection--alerts)** — anomaly
  score escalation and ignore thresholds, alert cost and averaging window, new-alert email,
  and streaming detection. *(`anomaly_detection_mode` and `anomaly_sampling_interval` are on
  `insightfinder_project` only.)*
- **[Incident and root cause analysis](#optional-arguments--incident--root-cause-analysis)** —
  the prediction look-ahead window, RCA candidate counts and probability thresholds, RCA
  ranking, causal analysis scope, and per-incident downtime cost. *(`causal_min_delay` and
  `normal_event_causal_flag` are on `insightfinder_project` only; `composite_rca_limit` is on
  `insightfinder_metric_project` only.)*
- **[Prediction rules](#optional-arguments--prediction-rules)** — the evidence count and
  confidence thresholds that gate predicted-incident alerts, and the active/inactive
  thresholds that promote and retire causal rules.
- **[Instance settings](#optional-arguments--instance-settings)** — instance-down detection
  and visibility, grouping by instance, Knowledge Base instance matching, and the edge/brain
  and trace-prompt flags. *(`instance_convert_flag`, `is_grouping_by_instance`,
  `is_edge_brain`, and `is_trace_prompt` are on `insightfinder_project` only.)*
- **[Advanced settings](#optional-arguments--advanced-settings)** — proxy, model spans,
  detection wait time, thread count, multi-hop causal search depth and breadth, large-project
  optimization, and the training filter. *(Only `proxy`, `min_valid_model_span`,
  `multi_hop_search_level`, `multi_hop_search_limit`, `large_project`, `training_filter`, and
  `new_pattern_range` are shared with metric projects.)*
- **[Incident priority](#optional-arguments--incident-priority)** — map anomaly score ranges
  to priority levels, and cap the priority used for ticket creation and suggestions.
- **[Holiday settings](#holiday_settings-block)** — date ranges treated as holidays for
  anomaly detection.
- **[Email notification settings](#email_setting)** — the per-project `email_setting` JSON
  object controlling detection, prediction, root-cause, and AI Watchtower emails.
- **[Webhook settings](#optional-arguments--webhook-settings)** — webhook URL, request size,
  dampening, type/blacklist/critical-keyword sets, and
  [custom headers](#webhook_header_list).
- **[Project sharing](#shared_usernames)** — the list of usernames the project is shared with.
- **[Instance grouping](#instance_grouping_update)** — the auto-fill toggle for instance
  grouping.

### 3. Log Project Settings

Applies to [`insightfinder_project`](#insightfinder_project) when
`project_creation_config.data_type` is `Log`. These settings have no effect on metric
projects.

- **[Log parsing and model training](#optional-arguments--log-settings)** — detection batch
  sizes, pattern limits, model size, keyword feature extraction, multi-line handling, JSON
  repair, and the anomaly event base-score weights.
- **[Hot, cold, and rare event detection](#optional-arguments--anomaly-detection--alerts)** —
  the frequency-spike, low-frequency, and novel-pattern detectors, with their thresholds,
  calm-down periods, and per-cycle limits.
- **[Log clustering and outlier sensitivity](#optional-arguments--advanced-settings)** —
  `similarity_sensitivity` (how alike two messages must be to share a pattern),
  `feature_outlier_sensitivity` / `feature_outlier_threshold` for numeric values extracted
  from logs, pattern model scope, whitelist capacity, and `zone_name_key`.
- **[JSON key settings](#json_key_settings-block)** — which JSON key paths to extract from
  log data, and which downstream features (summary, metafield, dampening, notifications)
  each key feeds.
1. **Log Label Settings** — filtering and labeling rules that tell the engine which log lines
   matter and how to classify them.
   - **[Standalone resource](#insightfinder_log_labels)** —
     `insightfinder_log_labels`, managed separately from the project.
   - **[Inline block](#log_label_settings-block)** — the equivalent
     `log_label_settings` block on `insightfinder_project`. Use one or the other for a given
     project, not both.
   - **[Supported label types](#supported-label_type-values)** — whitelist, blacklist,
     pattern naming and matching, severity, event ID, incident priority, and the rest.
2. **[Log-to-metric settings](#l2m_settings-block)** — turn log data into metric data points
   in a target metric project.
   - **[JSON parsers](#l2m_settings--json_parsers)** — key-path based extraction of metric
     value, instance, container, and timestamp, with aggregation mode and period.
   - **[Derived value model](#l2m_settings--json_parsers--derived_value_model)** — compute a
     metric from a base and actual value expression.
   - **[Regex parsers](#l2m_settings--regexs)** — the regex-based alternative, used when
     `json_flag` is `false`.
- **[LLM / trace-prompt evaluation](#llm_evaluation_setting)** — hallucination, relevance,
  toxicity, PII leakage, and bias evaluations for projects with `is_trace_prompt` enabled.
- **[ServiceNow as a data source](#project_servicenow_settings-block)** — connection and
  parsing settings for a project whose `project_cloud_type` is `ServiceNow`. This is about
  *ingesting* ServiceNow records; for *writing tickets back*, see
  [Integrations](#5-integrations-and-access).
- **[ServiceNow notification templates](#optional-arguments--servicenow-notification-templates)** —
  short-description and description format rules.
- **Other JSON configuration blocks** —
  [`base_value_setting`](#base_value_setting) (base values and metric mappings),
  [`log_to_log_setting_list`](#log_to_log_setting_list) (log-to-log transformations), and
  [`cdf_setting`](#cdf_setting).

### 4. Metric Project Settings

Applies to [`insightfinder_metric_project`](#insightfinder_metric_project).

- **[Shared settings](#common-optional-arguments-shared-with-insightfinder_project)** — the
  subset of general project settings that metric projects also accept, listed in one place.
- **[Baseline and detection tuning](#metric-specific-optional-arguments)** — dynamic baseline
  detection, positive and negative violation factors, baseline duration, double verification,
  period anomaly filtering, UBL and cumulative detection, and component-level detection.
- **[Gap filling](#metric-specific-optional-arguments)** — automatic filling of missing data
  points, whether filled points are stored, the training data length, and gap tolerance.
- **[Prediction](#metric-specific-optional-arguments)** — KPI prediction, forward prediction
  of metric values, training data length, and correlation sensitivity.
- **[Instance-down detection](#metric-specific-optional-arguments)** — the silence threshold,
  how many instances must be down, and the ratio threshold that triggers an alert.
- **[Per-metric alert thresholds](#metric_configurations-block)** — a map keyed by metric
  name, letting you escalate or ignore specific components per metric.
  - **[Metric alert settings](#metric_configurations--metric_alert_settings)** — the full set
    of alert and incident bounds per component, KPI flag, detection direction, per-metric
    C-value overrides, and pattern naming.

### 5. Integrations and Access

1. **[ServiceNow ticketing](#insightfinder_servicenow)** — write incidents out to ServiceNow.
   Covers authentication, the target field and content source, ticket trigger window,
   feedback collection, CMDB configuration items, and project-to-table mapping.
   - **[Per-project ticket settings](#project_configs-block)** — enable creation, update,
     consolidation-info update, and resolve-update per project, with a per-project
     configuration item override.
2. **[Slack notifications](#insightfinder_slack)** — post incident, prediction, and
   pattern-alert notifications to a Slack channel via an incoming webhook. A system can have
   several of these, each notifying a different channel.
   - **[Per-project overrides](#project_configs-block-1)** — a different channel, webhook,
     or notification type set per project, plus priority level filtering.
   - **[Match rules](#project_configs--rule)** — filter which alerts reach Slack by field
     name or content.
   - **[Behavior notes](#behavior-notes)** — how existing integrations are adopted rather
     than duplicated, and how updates locate the integration to modify.
3. **[JWT configuration](#insightfinder_jwt_config)** — the system-level JWT secret and type.

### 6. Lookups

Read-only data sources for referencing objects that already exist in InsightFinder.

- **[Project lookup](#insightfinder_project-data-source)** — fetch an existing project's
  display name and C/P values by name.
- **[Systems lookup](#insightfinder_systems)** — list all systems with their IDs and names.

---

## Provider Configuration

The provider authenticates against the InsightFinder API. Each argument may be supplied
in the provider block or through an environment variable; the provider block takes
precedence when both are set.

| Argument | Type | Description | Environment Variable | Required |
|---|---|---|---|---|
| `base_url` | String | The base URL for the InsightFinder API (e.g. `https://app.insightfinder.com`). | `IF_BASE_URL` | Yes (block or env var) |
| `username` | String, Sensitive | The username for InsightFinder authentication. | `IF_USERNAME` | Yes (block or env var) |
| `license_key` | String, Sensitive | The license key (API key) for InsightFinder authentication. | `IF_LICENSE_KEY` | Yes (block or env var) |

All three values must resolve to a non-empty string. If any is missing from both the
configuration and the environment, the provider fails at configure time with an
attribute-specific error.

### Example

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

Using environment variables instead:

```bash
export IF_BASE_URL="https://app.insightfinder.com"
export IF_USERNAME="your-username"
export IF_LICENSE_KEY="your-license-key"
```

```hcl
provider "insightfinder" {
  # values are read from IF_BASE_URL / IF_USERNAME / IF_LICENSE_KEY
}
```

---

## Resources

### insightfinder_project

Manages an InsightFinder project. This resource has extensive configuration options for
fine-tuning project behavior, and is the right choice for Log, Trace, and other
non-metric data types. For metric data, use
[`insightfinder_metric_project`](#insightfinder_metric_project).

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `project_name` | String | The unique name of the project. Changing this forces a new resource. |
| `system_name` | String | The system name this project belongs to. |
| `project_creation_config` | Block | A block defining project creation parameters (see below). |

#### project_creation_config Block

| Argument | Type | Required | Description |
|---|---|---|---|
| `data_type` | String | Yes | The type of data (e.g., `Log`, `Metric`, `Trace`). |
| `instance_type` | String | Yes | The instance type (e.g., `PrivateCloud`, `AWS`, `Azure`). Default: `PrivateCloud`. |
| `project_cloud_type` | String | Yes | The cloud type for the project (e.g., `PrivateCloud`). |
| `insight_agent_type` | String | No | The InsightFinder agent type (e.g., `Custom`, `LogStreaming`, `Historical`). |
| `servicenow_table` | String | No | The ServiceNow table name. Required when `project_cloud_type` is `ServiceNow`. |

#### Optional Arguments — General Settings

| Argument | Type | Description |
|---|---|---|
| `project_display_name` | String | The display name for the project. |
| `project_time_zone` | String | The timezone for the project. Default: `UTC`. |
| `sampling_interval` | Integer | The sampling interval in seconds. Default: `600`. |
| `c_value` | Integer | The C value for anomaly detection sensitivity (typically 2–5). |
| `p_value` | Float | The P value for anomaly detection probability (0.0–1.0). |
| `retention_time` | Integer | Data retention time in days. Default: `90`. |
| `ubl_retention_time` | Integer | Retention time for UBL data in days. Default: `90`. |
| `mode` | Integer | The process mode for the project. Applied via the `logdedicatedmode` API. |

#### Optional Arguments — Anomaly Detection & Alerts

<table>
<tr><th>Argument</th><th>Description</th></tr>

<tr><td><code>anomaly_detection_mode</code></td><td>
Enable/disable anomaly detection for log data. <strong>Type:</strong> Integer.<br>
Options:
<ul>
<li><strong>0</strong>: Enable anomaly detection (default).</li>
<li><strong>-1</strong>: Disable anomaly detection.</li>
</ul>
</td></tr>

<tr><td><code>anomaly_sampling_interval</code></td><td>
The time window (in seconds) used for log anomaly detection. <strong>Type:</strong> Integer.
<strong>Default:</strong> 60.
</td></tr>

<tr><td><code>enable_anomaly_score_escalation</code></td><td>
<strong>(Metric projects only)</strong> Enable/disable the system automatically "escalating"
incidents based on <code>escalation_anomaly_score_threshold</code>. <strong>Type:</strong> Boolean.
<strong>Default:</strong> false.
</td></tr>

<tr><td><code>escalation_anomaly_score_threshold</code></td><td>
<strong>(Metric projects only)</strong> If the calculated anomaly score is <strong>greater than
or equal to</strong> this threshold, the system flags it for escalation.
<strong>Type:</strong> String (numeric).<br>
<strong>Range:</strong> typically 0.0–1.0.
<ul>
<li>1.0 means only the most extreme anomalies are escalated.</li>
<li>0.1 means almost any deviation is escalated.</li>
</ul>
</td></tr>

<tr><td><code>ignore_anomaly_score_threshold</code></td><td>
<strong>(Metric projects only)</strong> Ignore anomalies with a score below this value.
<strong>Type:</strong> String (numeric). <strong>Default:</strong> 0.0. <strong>Range:</strong> 0.0–1.0.<br>
A value of 0.3 means any anomaly scoring <strong>less than 0.3</strong> is discarded.
</td></tr>

<tr><td><code>enable_hot_event</code></td><td>
<strong>(Log projects only)</strong> Enable/disable hot (frequency-spike) event detection.
<strong>Type:</strong> Boolean. <strong>Default:</strong> true.
<ul>
<li><strong>true</strong>: monitor for and report frequency-based log spikes.</li>
<li><strong>false</strong>: ignore log frequency spikes for known patterns.</li>
</ul>
</td></tr>

<tr><td><code>hot_event_threshold</code></td><td>
<strong>(Log projects only)</strong> The <strong>maximum allowable count</strong> of a specific log
pattern within a sampling interval before it is classified as an anomaly.
<strong>Type:</strong> Integer. <strong>Default:</strong> 50. <strong>Range:</strong> &gt;= 0.<br>
A very high number effectively suppresses Hot Event alerts.
</td></tr>

<tr><td><code>hot_event_calm_down_period</code></td><td>
<strong>(Log projects only)</strong> After a pattern triggers a Hot Event, subsequent spikes of the
same pattern within this period are recorded but do <strong>not</strong> generate new alerts.
<strong>Type:</strong> Integer. <strong>Default:</strong> 3.
<strong>Unit:</strong> sampling intervals (multiples of <code>anomaly_sampling_interval</code>).
<strong>Range:</strong> &gt; 0 — 0 or less falls back to the default.
</td></tr>

<tr><td><code>hot_event_detection_mode</code></td><td>
<strong>(Log projects only)</strong> How the system decides whether a log volume is high enough to
be "Hot." <strong>Type:</strong> Integer. <strong>Default:</strong> 0.
<ul>
<li><strong>0</strong>: standard detection logic (static thresholds plus basic heuristics).</li>
<li><strong>1</strong>: advanced statistical approach detecting spikes relative to moving averages.</li>
</ul>
</td></tr>

<tr><td><code>hot_number_limit</code></td><td>
<strong>(Log projects only)</strong> Limits the total number of unique log patterns that can be
classified as "Hot" in a single processing cycle. <strong>Type:</strong> Integer.
<strong>Default:</strong> 20. Empty/unset is treated as infinity.
</td></tr>

<tr><td><code>cold_event_threshold</code></td><td>
<strong>(Log projects only)</strong> Sensitivity for triggering "Cold Event" alerts.
<strong>Type:</strong> Integer. <strong>Default:</strong> 10. <strong>Range:</strong> &gt;= 0 —
0 effectively disables cold event detection.
</td></tr>

<tr><td><code>cold_number_limit</code></td><td>
<strong>(Log projects only)</strong> The maximum number of cold events detected per day.
<strong>Type:</strong> Integer. <strong>Default:</strong> 0.
</td></tr>

<tr><td><code>rare_anomaly_type</code></td><td>
<strong>(Log projects only)</strong> How the system categorizes "rare" (infrequent or new) log
patterns. <strong>Type:</strong> Integer. <strong>Default:</strong> 0.
<ul>
<li><strong>0</strong>: all rare events — both new patterns and known patterns with very low frequency.</li>
<li><strong>1</strong>: alert only when a log template is seen for the first time.</li>
<li><strong>2</strong>: alert only when an existing template contains a new/rare value.</li>
</ul>
</td></tr>

<tr><td><code>rare_event_alert_thresholds</code></td><td>
<strong>(Log projects only)</strong> Limit on the cluster size or frequency of a rare event before
it triggers an alert. <strong>Type:</strong> Integer. <strong>Default:</strong> 1.
<strong>Range:</strong> &gt;= 0. Higher is less sensitive; 1 alerts on the first occurrence.
</td></tr>

<tr><td><code>rare_number_limit</code></td><td>
<strong>(Log projects only)</strong> Limits the total number of unique log patterns that can be
classified as "Rare" at once. <strong>Type:</strong> Integer. <strong>Default:</strong> 20.
</td></tr>

<tr><td><code>collect_all_rare_events_flag</code></td><td>
<strong>(Log projects only)</strong> Whether to capture every rare event or only a representative
subset. <strong>Type:</strong> Boolean. <strong>Default:</strong> false.
</td></tr>

<tr><td><code>enable_new_alert_email</code></td><td>
Controls whether the system sends email alerts for newly detected anomalies and incidents.
<strong>Type:</strong> Boolean. <strong>Default:</strong> false.
</td></tr>

<tr><td><code>new_alert_flag</code></td><td>
<strong>(Log projects only)</strong> Manages how the system distinguishes "new" from "recurring"
anomalies in the UI and notification pipeline. <strong>Type:</strong> Boolean.
<strong>Default:</strong> false.
</td></tr>

<tr><td><code>alert_average_time</code></td><td>
<strong>(Metric projects only)</strong> The duration over which the system averages metric values
when evaluating alert conditions. <strong>Type:</strong> Integer. <strong>Unit:</strong> minutes.
<strong>Default:</strong> 0. A value of 1 gives virtually no smoothing; 15 or 30 gives significant
smoothing for stable, long-term trends.
</td></tr>

<tr><td><code>alert_hourly_cost</code></td><td>
A monetary value assigned to downtime or degraded performance.
<strong>Type:</strong> Float. <strong>Unit:</strong> currency per hour (e.g. USD).
<strong>Default:</strong> 0.0. <strong>Range:</strong> &gt;= 0.0.
</td></tr>

<tr><td><code>enable_stream_detection</code></td><td>
<strong>(Metric projects only)</strong> Switches the metric detection engine from batch processing
to <strong>Streaming Detection</strong>. <strong>Type:</strong> Boolean. <strong>Default:</strong> false.
</td></tr>

</table>

#### Optional Arguments — Log Settings

<table>
<tr><th>Argument</th><th>Description</th></tr>

<tr><td><code>log_detection_min_count</code></td><td>
Minimum count for triggering log detection in a batch. <strong>Type:</strong> Integer.
<strong>Default:</strong> 10000.<br>
Log detection triggers immediately once at least this many log entries are available in a
batch; otherwise it waits up to 3 minutes.
</td></tr>

<tr><td><code>log_detection_size</code></td><td>
Maximum count for triggering log detection in a batch. <strong>Type:</strong> Integer.
<strong>Default:</strong> 30000.
</td></tr>

<tr><td><code>log_pattern_limit_level</code></td><td>
The limit on distinct log patterns. <strong>Type:</strong> Integer. <strong>Default:</strong> 1024.<br>
Once more patterns than this are generated, new patterns are all assigned to the
MISC (-2) miscellaneous pattern.
</td></tr>

<tr><td><code>max_log_model_size</code></td><td>
Maximum training data sample size per model — the number of logs used to train the model.
<strong>Type:</strong> Integer. <strong>Default:</strong> 10000.
</td></tr>

<tr><td><code>keyword_feature_number</code></td><td>
Number of keyword features (feature vector length) generated for training the model.
<strong>Type:</strong> Integer. <strong>Default:</strong> 200.
</td></tr>

<tr><td><code>keyword_setting</code></td><td>
The keyword setting used when collecting keywords and during keyword queries.
<strong>Type:</strong> Integer. <strong>Default:</strong> -1.
<ul>
<li><strong>-1</strong>: disabled (default).</li>
<li><strong>0</strong>: letters only.</li>
<li><strong>1</strong>: letters and numbers.</li>
</ul>
</td></tr>

<tr><td><code>model_keyword_setting</code></td><td>
Keyword selection during model training. <strong>Type:</strong> Integer. <strong>Default:</strong> 0.
<ul>
<li><strong>0</strong>: letters only — collects only purely alphabetic keywords.</li>
<li><strong>1</strong>: letters and numbers.</li>
</ul>
</td></tr>

<tr><td><code>model_keyword_segment_k</code></td><td>
Number of keyword segments (K) used for the log model. <strong>Type:</strong> Integer.
</td></tr>

<tr><td><code>model_match_threshold</code></td><td>
Model match threshold for log pattern matching. <strong>Type:</strong> Float.
</td></tr>

<tr><td><code>disable_model_keyword_stats_collection</code></td><td>
Disable collection of keyword frequency statistics for model training.
<strong>Type:</strong> Boolean. <strong>Default:</strong> false.
</td></tr>

<tr><td><code>disable_log_compress_event</code></td><td>
Disable saving the log data (log compress event). <strong>Type:</strong> Boolean.
<strong>Default:</strong> false.
</td></tr>

<tr><td><code>disable_log_processing_flag</code></td><td>
Disable log processing for the project. <strong>Type:</strong> Boolean.
</td></tr>

<tr><td><code>log_anomaly_event_base_score</code></td><td>
Adjusts the anomaly score weight for each event type, in the order:
rare, hot, cold, detection alert, new pattern, critical. <strong>Type:</strong> String (JSON array).
<strong>Default:</strong> <code>"[5,0.01,0.0075,0.01,1,100]"</code>.
</td></tr>

<tr><td><code>multi_line_flag</code></td><td>
Enable the regex multiline flag for multi-line log processing. <strong>Type:</strong> Boolean.
<strong>Default:</strong> false.
</td></tr>

<tr><td><code>nlp_flag</code></td><td>
Enable/disable NLP. <strong>Type:</strong> Boolean. <strong>Default:</strong> false.
<strong>Deprecated.</strong>
</td></tr>

<tr><td><code>pretty_json_convertor_flag</code></td><td>
When enabled, the system reforms invalid JSON data into valid JSON.
<strong>Type:</strong> Boolean. <strong>Default:</strong> false.
</td></tr>

</table>

#### Optional Arguments — Incident & Root Cause Analysis

<table>
<tr><th>Argument</th><th>Description</th></tr>

<tr><td><code>incident_prediction_window</code></td><td>
The "look-ahead" time for the predictive analytics engine, for <strong>Metric</strong> and
<strong>Log</strong> projects. <strong>Type:</strong> Integer. <strong>Unit:</strong> minutes.
<strong>Default:</strong> 0. <strong>Range:</strong> &gt;= 0.
</td></tr>

<tr><td><code>min_incident_prediction_window</code></td><td>
Lower bound on the predictive alerting time, for <strong>Metric</strong> and <strong>Log</strong>
projects. <strong>Type:</strong> Integer. <strong>Unit:</strong> minutes. <strong>Default:</strong> 0.
</td></tr>

<tr><td><code>incident_prediction_event_limit</code></td><td>
Caps the maximum number of predicted incidents the system tracks, displays, or alerts on in a
single processing window. <strong>Type:</strong> Integer. <strong>Default:</strong> 5.
<strong>Range:</strong> &gt;= 0 — 0 suppresses all predicted incident displays.
</td></tr>

<tr><td><code>incident_relation_search_window</code></td><td>
Governs how the system links predicted incidents to actual detected events.
<strong>Type:</strong> Integer. <strong>Unit:</strong> minutes (converted internally to
milliseconds). <strong>Default:</strong> typically 60 (1 hour), or matches
<code>incident_prediction_window</code> depending on project type.
</td></tr>

<tr><td><code>root_cause_count_threshold</code></td><td>
The <strong>maximum number of root cause candidates</strong> the system identifies and presents for
a single anomaly or incident. <strong>Type:</strong> Integer. <strong>Default:</strong> 10.
<strong>Range:</strong> &gt;= 0 — 0 hides all root cause suggestions.
</td></tr>

<tr><td><code>root_cause_probability_threshold</code></td><td>
Statistical filter for root cause identification. <strong>Type:</strong> Float (percentage).
<strong>Range:</strong> 0.0–1.0. <strong>Default:</strong> 0.8 (80%).
</td></tr>

<tr><td><code>root_cause_log_message_search_range</code></td><td>
The time window surrounding an anomaly or incident during which the RCA engine searches for
relevant log entries. <strong>Type:</strong> Integer. <strong>Unit:</strong> minutes.
<strong>Range:</strong> &gt;= 0. <strong>Default:</strong> typically 60 (1 hour).
</td></tr>

<tr><td><code>root_cause_rank_setting</code></td><td>
Controls the algorithm used to prioritize and sort identified root causes. Where other settings
filter by <em>probability</em> or <em>count</em>, this determines the <strong>ranking logic</strong>
itself. <strong>Type:</strong> Integer. <strong>Default:</strong> 0.
<ul>
<li><strong>0</strong>: default (balanced) ranking logic.</li>
<li><strong>1+</strong>: alternative ranking modes (e.g. higher weight on specific data types or
different statistical models).</li>
</ul>
</td></tr>

<tr><td><code>maximum_root_cause_result_size</code></td><td>
Hard limit on the absolute maximum number of root cause entries the system will process and
present. <strong>Type:</strong> Integer. <strong>Range:</strong> &gt;= 0 — 0 makes the system fall
back to a built-in global default. <strong>Default:</strong> typically 20.
</td></tr>

<tr><td><code>avg_per_incident_downtime_cost</code></td><td>
A flat monetary value assigned to every incident or anomaly in the project.
<strong>Type:</strong> Float. <strong>Unit:</strong> currency (e.g. USD). <strong>Default:</strong> 0.0.
</td></tr>

<tr><td><code>causal_prediction_setting</code></td><td>
Determines the scope of causal analysis used for incident prediction — whether the system looks
for relationships within a single project, across projects, or both.
<strong>Type:</strong> Integer.
<ul>
<li><strong>0</strong>: <strong>default</strong> — the most comprehensive mode; looks for causal links
both within the project and across all other accessible projects.</li>
<li><strong>1</strong>: restrict the causal search to the current project only. Faster, but may miss
external root causes.</li>
<li><strong>2</strong>: only look for relationships where the cause originates in a different project.</li>
</ul>
</td></tr>

<tr><td><code>causal_min_delay</code></td><td>
The <strong>minimum time difference</strong> required between a "cause" event and an "effect" event
for the pair to be considered causal. <strong>Type:</strong> String (numeric).
<strong>Unit:</strong> minutes (converted internally to milliseconds). <strong>Range:</strong> &gt;= 0.
<strong>Default:</strong> typically 0 (no minimum delay enforced).
</td></tr>

<tr><td><code>normal_event_causal_flag</code></td><td>
<strong>(Log projects only)</strong> Whether "normal" log events — those <em>not</em> flagged as
anomalies — are included as candidates during causal analysis and RCA.
<strong>Type:</strong> Boolean. <strong>Default:</strong> false.
<ul>
<li><strong>true</strong>: include all events (normal and anomalous) in causal analysis.</li>
<li><strong>false</strong>: use only anomalous events.</li>
</ul>
</td></tr>

</table>

#### Optional Arguments — Prediction Rules

<table>
<tr><th>Argument</th><th>Description</th></tr>

<tr><td><code>prediction_count_threshold</code></td><td>
Minimum volume of evidence needed to trigger a predicted incident alert — a "frequency gate" for
the predictive engine. <strong>Type:</strong> Integer. <strong>Range:</strong> &gt;= 0.
<strong>Default:</strong> 1 (any single piece of evidence triggers a prediction).
</td></tr>

<tr><td><code>prediction_probability_threshold</code></td><td>
Confidence filter for proactive incident forecasting. <strong>Type:</strong> Float (percentage).
<strong>Range:</strong> 0.0–1.0. <strong>Default:</strong> 0.8 (80%).<br>
0.9 means only incidents with a 90%-or-higher predicted likelihood trigger an alert.
</td></tr>

<tr><td><code>prediction_rule_active_condition</code></td><td>
Maturity filter defining the prerequisite status a causal rule must reach before it can generate
a predicted incident alert. <strong>Type:</strong> Integer.
<ul>
<li><strong>0</strong>: <strong>unfiltered</strong> — use all discovered rules immediately.</li>
<li><strong>1</strong>: <strong>verified only</strong> — use only rules that have successfully predicted
an event at least once.</li>
<li><strong>2+</strong>: <strong>high confidence</strong> — require multiple historical verifications or
specific statistical significance.</li>
</ul>
<strong>Default:</strong> typically 1 (verified only).
</td></tr>

<tr><td><code>prediction_rule_active_threshold</code></td><td>
Minimum probability required for a causal rule to be promoted to "Active" status.
<strong>Type:</strong> Float (percentage). <strong>Range:</strong> 0.0–1.0.
<strong>Default:</strong> typically 0.7 (70%). Must be &gt;= <code>prediction_rule_inactive_threshold</code>.
</td></tr>

<tr><td><code>prediction_rule_false_positive_threshold</code></td><td>
Quality control for unreliable predictive rules: the number of false positives a causal rule may
accumulate before the system stops using it for proactive alerting.
<strong>Type:</strong> Integer (a count, not a rate — the Terraform schema is an integer, so a
fractional rate such as <code>0.3</code> is not a valid value). <strong>Range:</strong> &gt;= 0.
This is the project-level counterpart to <code>false_positive_tolerance</code> in
<a href="#incident-prediction-fields">knowledgebase_settings</a>.
</td></tr>

<tr><td><code>prediction_rule_inactive_threshold</code></td><td>
Retirement/demotion criterion — the lower boundary for a causal rule's confidence score, and the
counterpart to <code>prediction_rule_active_threshold</code>. <strong>Type:</strong> Float
(percentage). <strong>Range:</strong> 0.0–1.0. Must be &lt;= <code>prediction_rule_active_threshold</code>.
<strong>Default:</strong> typically 0.5 (50%).
</td></tr>

</table>

#### Optional Arguments — Instance Settings

<table>
<tr><th>Argument</th><th>Description</th></tr>

<tr><td><code>instance_convert_flag</code></td><td>
<strong>(Log projects)</strong> Manages how instance identifiers (hostnames, IP addresses) are
processed and indexed. <strong>Type:</strong> Boolean. <strong>Default:</strong> false.
<ul>
<li><strong>true</strong>: apply instance name conversion and normalization.</li>
<li><strong>false</strong>: keep instance names exactly as they appear in the raw log stream.</li>
</ul>
</td></tr>

<tr><td><code>instance_down_enable</code></td><td>
Toggles the "Instance Down" detection engine, which detects when a data source stops sending
information. <strong>Type:</strong> Boolean. <strong>Default:</strong> false.
</td></tr>

<tr><td><code>show_instance_down</code></td><td>
Whether "Instance Down" events are displayed in dashboards, anomaly timelines, and incident
reports. <strong>Type:</strong> Boolean. <strong>Default:</strong> true.
</td></tr>

<tr><td><code>is_grouping_by_instance</code></td><td>
<strong>(Log projects)</strong> Whether the log engine treats each instance as an independent entity
or pools them together for analysis. <strong>Type:</strong> Boolean.
<strong>Default:</strong> typically true.
<ul>
<li><strong>true</strong>: isolate analysis by instance.</li>
<li><strong>false</strong>: aggregate analysis across the entire project.</li>
</ul>
</td></tr>

<tr><td><code>ignore_instance_for_kb</code></td><td>
Whether the instance name is a required match when looking up known issues in the Knowledge Base.
<strong>Type:</strong> Boolean. <strong>Default:</strong> false.
<ul>
<li><strong>true</strong>: ignore the instance name (match by pattern only).</li>
<li><strong>false</strong>: require the instance name to match.</li>
</ul>
</td></tr>

<tr><td><code>is_edge_brain</code></td><td>
Architectural flag determining whether the project runs in a resource-constrained "Edge"
environment or a full-scale "Cloud/Brain" environment. <strong>Type:</strong> Boolean.
<strong>Default:</strong> false.
</td></tr>

<tr><td><code>is_trace_prompt</code></td><td>
<strong>(Log projects)</strong> Identifies the project as an LLM <strong>Trace Prompt</strong>
monitoring system, activating LLM-specific evaluation and observability features.
<strong>Type:</strong> Boolean. <strong>Default:</strong> false. See
<a href="#llm_evaluation_setting"><code>llm_evaluation_setting</code></a>.
</td></tr>

<tr><td><code>component_name_auto_overwrite</code></td><td>
Automatically overwrite component names with values from the data source.
<strong>Type:</strong> Boolean.
</td></tr>

</table>

#### Optional Arguments — Advanced Settings

<table>
<tr><th>Argument</th><th>Description</th></tr>

<tr><td><code>proxy</code></td><td>
Address of a proxy server InsightFinder must use when communicating with external systems, or
when an Action (such as an automated remediation script) executes on a remote server.
<strong>Type:</strong> String. <strong>Default:</strong> null / empty string.
<strong>Format:</strong> URL or hostname with optional port (e.g. <code>http://proxy.internal:8080</code>).
</td></tr>

<tr><td><code>daily_model_span</code></td><td>
How many days of historical data the ML engine looks back on to build the baseline for "normal"
behavior. <strong>Type:</strong> Integer. <strong>Unit:</strong> days. <strong>Range:</strong> &gt;= 1.
<strong>Default:</strong> 1.<br>
1 uses only the previous 24 hours; 14 or 30 gives a much more robust statistical baseline.
</td></tr>

<tr><td><code>min_valid_model_span</code></td><td>
Minimum required duration of data present in a model before it is considered "valid" and ready
for production use. <strong>Type:</strong> Integer. <strong>Unit:</strong> milliseconds.
<strong>Range:</strong> &gt; 0. <strong>Default:</strong> typically 21,600,000 ms (6 hours).
</td></tr>

<tr><td><code>maximum_detection_wait_time</code></td><td>
<strong>(Log projects)</strong> The maximum time the log anomaly detection engine waits for
late-arriving logs before analyzing a time window. <strong>Type:</strong> Integer.
<strong>Unit:</strong> minutes. <strong>Range:</strong> typically 1–60. <strong>Default:</strong> 15.
</td></tr>

<tr><td><code>maximum_threads</code></td><td>
<strong>(Log projects)</strong> Degree of parallelism for heavy background tasks such as log
training, detection, and replay. <strong>Type:</strong> Integer. <strong>Range:</strong> typically
1–16 (minimum floor of 1 enforced by the API). <strong>Default:</strong> 1 (sequential processing).
</td></tr>

<tr><td><code>multi_hop_search_level</code></td><td>
Depth — the "number of hops" — the system traverses in the causal dependency graph during RCA.
<strong>Type:</strong> Integer. <strong>Range:</strong> &gt;= 0 — 0 disables deep causal searching
beyond the most obvious correlations. <strong>Default:</strong> typically 1.
</td></tr>

<tr><td><code>multi_hop_search_limit</code></td><td>
Breadth of the causal search — the maximum number of neighbor nodes or candidates explored at each
hop. Where <code>multi_hop_search_level</code> controls <em>depth</em>, this controls <em>breadth</em>.
<strong>Type:</strong> String (numeric). <strong>Default:</strong> typically 10.
</td></tr>

<tr><td><code>new_pattern_number_limit</code></td><td>
<strong>(Log projects)</strong> Numerical cap on the number of unique, previously unseen log patterns
the system may identify and track in a single processing interval. <strong>Type:</strong> Integer.
Empty or null is treated as infinity.
</td></tr>

<tr><td><code>new_pattern_range</code></td><td>
<strong>(Log projects)</strong> Suppression / calm-down window for alerts triggered by newly
discovered log templates. <strong>Type:</strong> Integer. <strong>Unit:</strong> sampling intervals
(multiples of <code>anomaly_sampling_interval</code>). <strong>Range:</strong> &gt;= 0.
<strong>Default:</strong> typically 3.
</td></tr>

<tr><td><code>project_model_flag</code></td><td>
<strong>(Log projects)</strong> Toggles between <strong>Local</strong> and <strong>Global</strong> pattern
modeling. <strong>Type:</strong> Boolean. <strong>Default:</strong> typically true.
<ul>
<li><strong>true</strong>: share log patterns across all instances in the project (Project Model).</li>
<li><strong>false</strong>: isolate log pattern learning to each individual instance (Instance Model).</li>
</ul>
</td></tr>

<tr><td><code>large_project</code></td><td>
Informs the processing engine that the project contains an exceptionally large number of instances
or very high data throughput, enabling several optimizations. <strong>Type:</strong> Boolean.
<strong>Default:</strong> false.
</td></tr>

<tr><td><code>similarity_sensitivity</code></td><td>
<strong>(Log projects)</strong> Strictness of the log clustering engine — how similar two log messages
must be to be grouped under the same Log Template (Pattern). <strong>Type:</strong> String.
<ul>
<li><strong>high</strong>: very strict; logs must be almost exactly the same.</li>
<li><strong>medium</strong>: default; balanced for most standard application logs.</li>
<li><strong>low</strong>: loose; groups logs even with significant differences.</li>
</ul>
</td></tr>

<tr><td><code>feature_outlier_sensitivity</code></td><td>
<strong>(Log projects)</strong> How aggressively the system flags numerical values extracted from
logs (latency, status codes, thread counts) as anomalous. <strong>Type:</strong> String.
<ul>
<li><strong>high</strong>: highly sensitive; flags minor statistical deviations.</li>
<li><strong>medium</strong>: default; flags significant deviations.</li>
<li><strong>low</strong>: least sensitive; flags only extreme statistical outliers.</li>
</ul>
</td></tr>

<tr><td><code>feature_outlier_threshold</code></td><td>
<strong>(Log projects)</strong> Manual override / hard cutoff for numerical outlier detection.
<strong>Type:</strong> Float. <strong>Default:</strong> 0.0.
<ul>
<li><strong>0.0</strong>: rely entirely on the statistical logic of <code>feature_outlier_sensitivity</code>.</li>
<li><strong>positive value</strong>: acts as a hard limit for detection.</li>
</ul>
</td></tr>

<tr><td><code>training_filter</code></td><td>
Historical filter for incident generation. When enabled, the system suppresses incidents occurring
outside a "known training window." <strong>Type:</strong> Boolean. <strong>Default:</strong> typically false.
</td></tr>

<tr><td><code>whitelist_number_limit</code></td><td>
How many "known safe" log patterns the project may track. If a log message matches a pattern on
this list, the system treats it as normal behavior and does not flag it.
<strong>Type:</strong> Integer. <strong>Default:</strong> 100.
</td></tr>

<tr><td><code>zone_name_key</code></td><td>
Tells the system which field ("key") in your logs contains the Zone — a cloud region
(<code>us-east-1</code>), data center (<code>DC-01</code>), or environment area.
<strong>Type:</strong> String. Enter the attribute name exactly as it appears in your log files
(often in JSON logs).
</td></tr>

</table>

#### Optional Arguments — Webhook Settings

| Argument | Type | Description |
|---|---|---|
| `webhook_url` | String | Webhook URL. |
| `max_web_hook_request_size` | Integer | Maximum webhook request size. |
| `webhook_alert_dampening` | Integer | Alert dampening period for webhooks (milliseconds). |
| `webhook_type_set_str` | String (JSON) | Type set string for webhooks. |
| `webhook_black_list_set_str` | String (JSON) | Blacklist set string for webhooks. |
| `webhook_critical_keyword_set_str` | String (JSON) | Critical keyword set string for webhooks. |
| `webhook_header_list` | String (JSON) | JSON array of header objects — see [webhook_header_list](#webhook_header_list). |

#### Optional Arguments — Incident Priority

| Argument | Type | Description |
|---|---|---|
| `incident_priority_by_anomaly_score_setting` | String (JSON) | Incident priority derived from anomaly score. Contains an `enabled` boolean and a `priorityScoreRangeMap` mapping priority levels (1–5) to score ranges. Example: `{"enabled":true,"priorityScoreRangeMap":{"1":"10001-","2":"5001-10000"}}`. |
| `incident_priority_cap_setting` | String (JSON) | Incident priority caps. Contains `ticketCreationPriorityCap` and `suggestedPriorityCap` string values. Example: `{"ticketCreationPriorityCap":"2","suggestedPriorityCap":"5"}`. |

#### Optional Arguments — ServiceNow Notification Templates

| Argument | Type | Description |
|---|---|---|
| `service_now_short_description_format` | String | Rules for short-description content in ServiceNow notifications. |
| `service_now_description_format` | String | Rules for non-key-value notification content in ServiceNow notifications. |

#### JSON Configuration Strings

These arguments accept JSON-formatted strings for complex configurations. Use `jsonencode()`
to build them from HCL where practical.

##### email_setting

A JSON object used for configuring email notifications.

```json
{
  "enableIncidentDetectionEmailAlert": true,
  "enableIncidentPredictionEmailAlert": true,
  "enableRootCauseEmailAlert": true,
  "enableAlertsEmail": true,
  "enableNotificationAW": false,
  "onlySendWithRCA": false,
  "emailDampeningPeriod": 3600000,
  "alertsEmailDampeningPeriod": 3600000,
  "predictionEmailDampeningPeriod": 3600000,
  "awSeverityLevel": "critical"
}
```

| Field | Type | Description |
|---|---|---|
| `enableIncidentDetectionEmailAlert` | Boolean | Enable email alerts for incident detection. |
| `enableIncidentPredictionEmailAlert` | Boolean | Enable email alerts for incident prediction. |
| `enableRootCauseEmailAlert` | Boolean | Enable email alerts including root cause analysis. |
| `enableAlertsEmail` | Boolean | Enable general alerts email. |
| `enableNotificationAW` | Boolean | Enable notification for AI Watchtower. |
| `onlySendWithRCA` | Boolean | Only send alerts if RCA is available. |
| `emailDampeningPeriod` | Integer | Dampening period in milliseconds. |
| `alertsEmailDampeningPeriod` | Integer | Dampening period for alerts, in milliseconds. |
| `predictionEmailDampeningPeriod` | Integer | Dampening period for prediction alerts, in milliseconds. |
| `awSeverityLevel` | String | Severity level for AI Watchtower notifications. |

##### llm_evaluation_setting

A JSON object for configuring LLM evaluation metrics. Relevant when
[`is_trace_prompt`](#optional-arguments--instance-settings) is enabled.

```json
{
  "isHallucinationEvaluation": true,
  "isAnswerRelevantEvaluation": true,
  "isLogicConsistencyEvaluation": true,
  "isFactualInaccuracyEvaluation": true,
  "isMaliciousPromptEvaluation": true,
  "isToxicityEvaluation": true,
  "isPiiPhiLeakageEvaluation": true,
  "isTopicGuardrailsEvaluation": false,
  "isToneDetectionEvaluation": false,
  "isAnomalousOutliersEvaluation": false,
  "showSafetyTemplate": false,
  "isGenderBiasEvaluation": false,
  "isRacialBiasEvaluation": false,
  "isSocioeconomicBiasEvaluation": false,
  "isCulturalBiasEvaluation": false,
  "isReligiousBiasEvaluation": false,
  "isPoliticalBiasEvaluation": false,
  "isDisabilityBiasEvaluation": false,
  "isAgeBiasEvaluation": false
}
```

All fields are booleans. The last eight are bias evaluations.

##### base_value_setting

A JSON object for configuring base values and metric mappings.

```json
{
  "isSourceProject": true,
  "mappingKeys": [],
  "baseValueKeys": [],
  "metricProjects": [],
  "additionalMetricNames": []
}
```

| Field | Type | Description |
|---|---|---|
| `isSourceProject` | Boolean | Whether this is a source project. |
| `mappingKeys` | Array of String/Object | Keys for mapping. |
| `baseValueKeys` | Array of String/Object | Keys for base values. |
| `metricProjects` | Array of String | List of metric projects. |
| `additionalMetricNames` | Array of String | Additional metric names. |

##### instance_grouping_update

A JSON object for instance grouping settings.

```json
{ "autoFill": true }
```

| Field | Type | Description |
|---|---|---|
| `autoFill` | Boolean | Enable auto-fill for instance grouping. |

##### shared_usernames

A JSON array of usernames to share the project with.

```json
["user1", "user2"]
```

##### webhook_header_list

A JSON array of header objects to be included in webhook requests.

```json
[
  { "headerName": "Authorization", "headerValue": "Bearer token" }
]
```

##### log_to_log_setting_list

A JSON array for configuring log-to-log transformations. Each element is a transformation
rules object.

##### cdf_setting

A JSON array for configuring Conditional Data Filtering (CDF) / Component Definition File
settings. Each element is a CDF object.

#### project_servicenow_settings Block

Used when `project_creation_config.project_cloud_type` is set to `ServiceNow`
(case-insensitive). Defines the parameters InsightFinder needs to connect to and retrieve
data from a ServiceNow instance.

| Argument | Type | Required | Description |
|---|---|---|---|
| `host` | String | Yes | The base URL of the ServiceNow instance (e.g. `https://[instance].service-now.com/`). |
| `servicenow_user` | String | Yes | The username for the ServiceNow account. |
| `servicenow_password` | String, Sensitive | Yes | The password for the ServiceNow account. |
| `client_id` | String | No | The OAuth client ID used for token-based authentication. |
| `client_secret` | String, Sensitive | No | The OAuth client secret used for token-based authentication. |
| `instance_field` | String | No | The field in the ServiceNow record (e.g. `short_description`) that contains the Instance Name to be monitored. |
| `instance_field_regex` | String | No | The regex applied to `instance_field` to extract the Instance Name. |
| `component_name_rule` | String | No | Rule for determining the component name from ServiceNow data. |
| `timestamp_format` | String | No | The Java `SimpleDateFormat` used to parse the timestamp field (e.g. `yyyy-MM-dd HH:mm:ss`). |
| `sysparm_query` | String | No | An optional ServiceNow filter query (encoded string) limiting the records fetched. Default: empty string. |
| `proxy` | String | No | Proxy server URL for InsightFinder to use when connecting to ServiceNow. Default: empty string. |
| `additional_fields` | List of String | No | Extra fields to retrieve from the ServiceNow record for inclusion in the InsightFinder data stream. |
| `service_now_import_flag` | Boolean | No | Whether to enable importing data from ServiceNow. |

#### holiday_settings Block

A set of holiday periods to be treated as holidays for anomaly detection purposes.

| Argument | Type | Required | Description |
|---|---|---|---|
| `name` | String | Yes | Name of the holiday. |
| `start_date` | String | Yes | Start date in `MM-DD` format (e.g. `12-25`). |
| `end_date` | String | Yes | End date in `MM-DD` format (e.g. `12-26`). |

```hcl
holiday_settings = [
  { name = "Christmas", start_date = "12-25", end_date = "12-26" },
  { name = "New Year",  start_date = "01-01", end_date = "01-01" },
]
```

#### log_label_settings Block

A set of log label settings applied to the project. Each setting is applied individually
via the API. This is the inline equivalent of the standalone
[`insightfinder_log_labels`](#insightfinder_log_labels) resource — use one or the other
for a given project, not both.

| Argument | Type | Required | Description |
|---|---|---|---|
| `label_type` | String | Yes | Type of log label (`whitelist`, `blacklist`, `patternName`, etc. — see the [label type reference](#supported-label_type-values)). |
| `log_label_string` | String | Yes | The log label value/pattern, as a JSON array string. |

#### json_key_settings Block

A set of custom JSON key extraction settings. Each entry defines one JSON key to extract
from log data and controls which downstream features it participates in.

| Argument | Type | Required | Description |
|---|---|---|---|
| `json_key` | String | Yes | The JSON key path to extract from logs (e.g. `alert->core->id`). |
| `type` | String | Yes | The data type of the JSON value (e.g. `string`, `number`, `JSONArray`). |
| `summary_setting` | Boolean | Yes | Include this key in summary statistics. |
| `metafield_setting` | Boolean | Yes | Include this key in metafield settings. |
| `dampening_field_setting` | Boolean | Yes | Include this key in dampening field settings. |
| `notification_setting` | Boolean | No | Include this key in notification settings. |
| `notification_setting_display_name` | String | No | Display name for this key in notification settings. |
| `service_now_notification_setting` | Boolean | No | Include this key in ServiceNow notification settings. |
| `service_now_notification_setting_display_name` | String | No | Display name for this key in ServiceNow notifications. |

#### l2m_settings Block

A set of log-to-metric (L2M) settings. Each entry maps this log project to a target metric
project and specifies how log data is parsed into metric data points.

| Argument | Type | Required | Description |
|---|---|---|---|
| `metric_project_name` | String | Yes | Name of the target metric project. |
| `json_flag` | Boolean | No | Use JSON parsers (`true`) or regex parsers (`false`). |
| `enable_mapping` | Boolean | No | Enable this L2M mapping. |
| `json_parsers` | List of Block | No | JSON parser objects. Used when `json_flag = true`. |
| `regexs` | List of Block | No | Regex parser objects. Used when `json_flag = false`. |

##### l2m_settings → json_parsers

| Argument | Type | Description |
|---|---|---|
| `metric_value_key` | String | JSON key path for the metric value. |
| `base_value_key` | String | JSON key path for the base value. |
| `instance_name_key` | String | JSON key path for the instance name. |
| `container_name_key` | String | JSON key path for the container name. |
| `timestamp_key` | String | JSON key path for the timestamp. |
| `timestamp_format` | String | Format string for the extracted timestamp. |
| `data_filter` | String | Filter expression applied to log data before parsing. |
| `operation` | Integer | Parser operation type. |
| `metric_name` | String | Name for the derived metric. |
| `additional_metric_name` | String | JSON key path for an additional metric name. |
| `aggregation_mode` | Integer | Aggregation mode for combining values. |
| `grouping_by_component` | Boolean | Group metric data by component. |
| `aggregation_period` | Integer | Aggregation period. |
| `container_type` | Integer | Container type. |
| `derived_value_model` | Block | Optional nested derived-value configuration (see below). |

##### l2m_settings → json_parsers → derived_value_model

| Argument | Type | Description |
|---|---|---|
| `base_value` | String | Base value expression. |
| `actual_value` | String | Actual value expression. |
| `operation` | Integer | Derived value operation type. |
| `mapping_id_list` | List of String | List of mapping IDs. |

##### l2m_settings → regexs

| Argument | Type | Description |
|---|---|---|
| `metric_name_regex` | String | Regex to extract the metric name. |
| `metric_value_regex` | String | Regex to extract the metric value. |
| `base_value_key` | String | Key for the base value. |
| `instance_name_regex` | String | Regex to extract the instance name. |
| `container_name_regex` | String | Regex to extract the container name. |
| `timestamp_regex` | String | Regex to extract the timestamp. |
| `timestamp_format` | String | Format string for the extracted timestamp. |
| `data_filter` | String | Data filter expression. |
| `operation` | Integer | Parser operation type. |
| `aggregation_mode` | Integer | Aggregation mode. |
| `metric_name` | String | Metric name. |
| `grouping_by_component` | Boolean | Group metric data by component. |
| `aggregation_period` | Integer | Aggregation period. |
| `container_type` | Integer | Container type. |

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | The unique identifier for the project (same as `project_name`). |

#### Example

```hcl
resource "insightfinder_project" "advanced_logs" {
  project_name         = "advanced-system-logs"
  project_display_name = "System Logs with ML"
  system_name          = "Production"

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "AWS"
    project_cloud_type = "AWS"
    insight_agent_type = "Historical"
  }

  project_time_zone = "America/New_York"
  sampling_interval = 600
  retention_time    = 180
  c_value           = 3
  p_value           = 0.95

  enable_hot_event      = true
  hot_event_threshold   = 50
  rare_anomaly_type     = 0
  similarity_sensitivity = "medium"

  email_setting = jsonencode({
    enableIncidentDetectionEmailAlert = true
    enableRootCauseEmailAlert         = true
    emailDampeningPeriod              = 3600000
  })

  holiday_settings = [
    { name = "Christmas", start_date = "12-25", end_date = "12-26" },
  ]
}
```

#### Import

```shell
terraform import insightfinder_project.example my-project-name
```

---

### insightfinder_metric_project

Manages an InsightFinder metric project. This resource is purpose-built for metric data and
provides metric-specific settings (baseline detection, gap filling, KPI prediction, per-metric
alert thresholds) in addition to the general settings shared with `insightfinder_project`.

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `project_name` | String | The unique name of the metric project. Changing this forces a new resource. |
| `system_name` | String | The system name this project belongs to. |
| `project_creation_config` | Block | A block defining project creation parameters (see below). |

#### project_creation_config Block

| Argument | Type | Required | Description |
|---|---|---|---|
| `data_type` | String | Yes | The type of data. Typically `Metric`. |
| `instance_type` | String | Yes | The instance type (e.g. `PrivateCloud`, `LogToMetric`, `OnPremise`). |
| `project_cloud_type` | String | Yes | The cloud type for the project (e.g. `PrivateCloud`, `LogToMetric`). |
| `insight_agent_type` | String | No | The InsightFinder agent type (e.g. `Custom`). |

#### Common Optional Arguments (shared with insightfinder_project)

These behave identically to the same-named arguments on `insightfinder_project`. Refer to
that section for full descriptions.

| Argument | Type | Description |
|---|---|---|
| `project_display_name` | String | Display name for the project. |
| `project_time_zone` | String | Timezone (default: `UTC`). |
| `sampling_interval` | Integer | Sampling interval in seconds. |
| `c_value` | Integer | C value for anomaly sensitivity (typically 2–5). |
| `p_value` | Float | P value for anomaly probability (0.0–1.0). |
| `retention_time` | Integer | Data retention time in days. |
| `ubl_retention_time` | Integer | UBL retention time in days. |
| `training_filter` | Boolean | Training filter flag. |
| `enable_new_alert_email` | Boolean | Enable new alert email notifications. |
| `large_project` | Boolean | Optimize processing for large-scale data. |
| `new_pattern_range` | Integer | Suppression window for new pattern alerts (sampling intervals). |
| `proxy` | String | Proxy server URL for external connections. |
| `mode` | Integer | The process mode for the project (set via the `logdedicatedmode` API). |
| `enable_anomaly_score_escalation` | Boolean | Enable anomaly score escalation. |
| `escalation_anomaly_score_threshold` | String | Threshold for anomaly score escalation. |
| `ignore_anomaly_score_threshold` | String | Ignore anomalies with a score below this threshold. |
| `enable_stream_detection` | Boolean | Enable the streaming detection pipeline. |
| `ignore_instance_for_kb` | Boolean | Ignore instance name when matching KB entries. |
| `show_instance_down` | Boolean | Show instance-down incidents in the UI. |
| `instance_down_enable` | Boolean | Enable instance-down detection. |
| `alert_hourly_cost` | Float | Hourly monetary cost for alerts (e.g. USD). |
| `alert_average_time` | Integer | Smoothing window for alert average time calculations. |
| `avg_per_incident_downtime_cost` | Float | Average monetary cost per incident downtime. |
| `incident_prediction_window` | Integer | Look-ahead window for incident prediction (minutes). |
| `min_incident_prediction_window` | Integer | Minimum incident prediction window (minutes). |
| `incident_relation_search_window` | Integer | Window for linking predicted to actual incidents (minutes). |
| `incident_prediction_event_limit` | Integer | Max predicted incidents tracked per processing window. |
| `root_cause_count_threshold` | Integer | Max root cause candidates returned per incident. |
| `root_cause_probability_threshold` | Float | Min probability for a root cause candidate (0.0–1.0). |
| `root_cause_log_message_search_range` | Integer | Search range for RCA log messages (minutes). |
| `causal_prediction_setting` | Integer | Causal analysis scope (0 = all, 1 = within project, 2 = cross project). |
| `root_cause_rank_setting` | Integer | Ranking algorithm for root causes. |
| `maximum_root_cause_result_size` | Integer | Hard limit on RCA entries returned. |
| `multi_hop_search_level` | Integer | Depth of causal graph traversal. |
| `multi_hop_search_limit` | String | Max neighbors explored at each causal hop (string-encoded integer). |
| `prediction_count_threshold` | Integer | Min evidence count to trigger a prediction alert. |
| `prediction_probability_threshold` | Float | Min confidence for prediction alerts (0.0–1.0). |
| `prediction_rule_active_condition` | Integer | Maturity filter for causal rules used in prediction. |
| `prediction_rule_active_threshold` | Float | Min probability to promote a rule to Active (0.0–1.0). |
| `prediction_rule_false_positive_threshold` | Integer | Max false-positive count before disabling a rule. |
| `prediction_rule_inactive_threshold` | Float | Probability below which a rule is demoted (0.0–1.0). |
| `min_valid_model_span` | Integer | Min data duration (ms) required before a model is used. |
| `component_name_auto_overwrite` | Boolean | Automatically overwrite component names from the data source. |
| `webhook_url` | String | Webhook URL. |
| `max_web_hook_request_size` | Integer | Maximum webhook request size. |
| `webhook_alert_dampening` | Integer | Alert dampening period for webhooks (ms). |
| `webhook_black_list_set_str` | String (JSON) | Blacklist pattern set for webhooks. |
| `webhook_critical_keyword_set_str` | String (JSON) | Critical keyword set for webhooks. |
| `webhook_type_set_str` | String (JSON) | Type set for webhooks. |
| `webhook_header_list` | String (JSON) | JSON array of webhook header objects (`{headerName, headerValue}`). |
| `email_setting` | String (JSON) | Email notification settings. Same structure as `insightfinder_project`. |
| `instance_grouping_update` | String (JSON) | Instance grouping update settings (e.g. `{"autoFill": false}`). |
| `shared_usernames` | String (JSON) | JSON array of usernames to share the project with. |
| `incident_priority_by_anomaly_score_setting` | String (JSON) | Incident priority derived from anomaly score. |
| `incident_priority_cap_setting` | String (JSON) | Incident priority caps (`ticketCreationPriorityCap`, `suggestedPriorityCap`). |
| `linked_log_projects` | String (JSON) | JSON array of log project names linked to this metric project for RCA. |
| `holiday_settings` | List of Block | Holiday settings — same `name` / `start_date` / `end_date` structure as `insightfinder_project`. |

#### Metric-Specific Optional Arguments

| Argument | Type | Description |
|---|---|---|
| `composite_rca_limit` | Integer | Limit for composite root cause analysis. |
| `high_ratio_c_value` | Integer | High-ratio C value for anomaly detection. Used for metrics that change dramatically. |
| `maximum_hint` | Integer | Maximum hint value for anomaly detection. |
| `dynamic_baseline_detection_flag` | Boolean | Enable dynamic baseline detection instead of static thresholds. |
| `positive_baseline_violation_factor` | Float | Multiplier for detecting positive (upward) baseline violations. |
| `negative_baseline_violation_factor` | Float | Multiplier for detecting negative (downward) baseline violations. |
| `enable_period_anomaly_filter` | Boolean | Filter out anomalies that follow a known periodic pattern. |
| `enable_ubl_detect` | Boolean | Enable UBL (Unsupervised Baseline Learning) detection. |
| `enable_cumulative_detect` | Boolean | Enable cumulative anomaly detection mode. |
| `enable_component_level_detection` | Boolean | Enable anomaly detection at the component level. |
| `prediction_training_data_length` | Integer | Historical data length (ms) for training the prediction model. |
| `prediction_correlation_sensitivity` | Float | Sensitivity for detecting metric correlations in prediction (0.0–1.0). |
| `enable_kpi_prediction` | Boolean | Enable KPI prediction. |
| `instance_down_threshold` | Integer | Silence duration (ms) before an instance is considered down. |
| `instance_down_report_number` | Integer | Number of instances that must be down before an alert is generated. |
| `instance_down_ratio_threshold` | Float | Fraction (0.0–1.0) of instances that must be down to trigger an alert. |
| `model_span` | Integer | Data span (ms) used by the detection model. |
| `enable_metric_data_prediction` | Boolean | Enable forward prediction of metric data values. |
| `enable_baseline_detection_double_verify` | Boolean | Require a second verification pass before flagging a baseline deviation. |
| `enable_fill_gap` | Boolean | Enable automatic gap-filling for missing metric data points. |
| `enable_store_filled_gap` | Boolean | Persist gap-filled data points to storage. |
| `gap_filling_training_data_length` | Integer | Historical data length (ms) used to train the gap-filling model. |
| `pattern_id_generation_rule` | Integer | Rule used to generate internal pattern IDs. |
| `anomaly_gap_tolerance_count` | Integer | Consecutive missing data points tolerated before being counted as an anomaly. |
| `filter_by_anomaly_in_baseline_generation` | Boolean | Exclude anomalous data when building the baseline model. |
| `baseline_duration` | Integer | Duration (ms) of the window used to calculate the baseline. |
| `anomaly_dampening` | Integer | Dampening period (ms) between consecutive anomaly alerts for the same metric. |
| `component_metric_setting_overall_model_list` | String (JSON) | JSON array specifying the overall model list for component-level metric settings. |

#### metric_configurations Block

An optional map of per-metric alert threshold and component configurations. **The map key is
the metric name.**

| Argument | Type | Description |
|---|---|---|
| `escalate_incident_components` | List of String | Components for which incidents are escalated. Use `["Global_<hash>"]` to select all. |
| `ignored_components` | List of String | Components ignored for this metric. Use `["Global_<hash>"]` to select all. |
| `metric_alert_settings` | List of Block | Alert threshold settings per component for this metric (see below). |

##### metric_configurations → metric_alert_settings

| Argument | Type | Required | Description |
|---|---|---|---|
| `component_name` | String | Yes | Component name. Use `Global_<hash>` for the global setting. |
| `threshold_alert_lower_bound` | String | No | Alert lower bound threshold. |
| `threshold_alert_upper_bound` | String | No | Alert upper bound threshold. |
| `threshold_alert_lower_bound_negative` | String | No | Alert lower bound negative threshold. |
| `threshold_alert_upper_bound_negative` | String | No | Alert upper bound negative threshold. |
| `threshold_no_alert_lower_bound` | String | No | No-alert lower bound threshold. |
| `threshold_no_alert_upper_bound` | String | No | No-alert upper bound threshold. |
| `threshold_no_alert_lower_bound_negative` | String | No | No-alert lower bound negative threshold. |
| `threshold_no_alert_upper_bound_negative` | String | No | No-alert upper bound negative threshold. |
| `incident_alert_lower_bound` | String | No | Incident alert lower bound. |
| `incident_alert_upper_bound` | String | No | Incident alert upper bound. |
| `incident_alert_lower_bound_negative` | String | No | Incident alert lower bound negative. |
| `incident_alert_upper_bound_negative` | String | No | Incident alert upper bound negative. |
| `incident_no_alert_lower_bound` | String | No | Incident no-alert lower bound. |
| `incident_no_alert_upper_bound` | String | No | Incident no-alert upper bound. |
| `incident_no_alert_lower_bound_negative` | String | No | Incident no-alert lower bound negative. |
| `incident_no_alert_upper_bound_negative` | String | No | Incident no-alert upper bound negative. |
| `is_kpi` | Boolean | No | Whether this metric is a KPI. |
| `is_flapping_result_only` | Boolean | No | Whether to report flapping results only. |
| `incident_duration_threshold` | Integer | No | Minimum incident duration (ms) required to trigger. |
| `detection_type` | String | No | Detection direction: `positive`, `negative`, or `both`. |
| `c_value_override` | Integer | No | Override for the C value anomaly sensitivity. Null uses the project default. |
| `high_c_value_override` | Integer | No | Override for the high-ratio C value anomaly sensitivity. Null uses the project default. |
| `pattern_name_higher` | String | No | Pattern name for higher anomalies. |
| `pattern_name_lower` | String | No | Pattern name for lower anomalies. |
| `metric_type` | String | No | Metric type classification (e.g. `Unknown`, `CPU Utilization`). |
| `fill_zero` | Boolean | No | Fill missing data with zero. |
| `enable_baseline_near_constance` | Boolean | No | Enable baseline near-constance detection. |
| `compute_difference` | Boolean | No | Compute difference for this metric. |
| `anomaly_gap_tolerance_duration` | Integer | No | Anomaly gap tolerance duration in milliseconds. |
| `detection_anomaly_type` | Integer | No | Anomaly detection type integer (0, 1, 2, …). |
| `rouge_value` | String (JSON) | No | Rouge value as a raw JSON string from the API (e.g. `{"l":NaN,"s":NaN}`). Null to disable. |

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | The unique identifier for the project (same as `project_name`). |

#### Example

```hcl
resource "insightfinder_metric_project" "tuned_metrics" {
  project_name         = "tuned-infrastructure-metrics"
  project_display_name = "Tuned Infrastructure Metrics"
  system_name          = "Production"

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "OnPremise"
    project_cloud_type = "OnPremise"
    insight_agent_type = "Custom"
  }

  project_time_zone = "US/Eastern"
  sampling_interval = 60
  retention_time    = 90

  c_value                            = 3
  p_value                            = 0.95
  high_ratio_c_value                 = 3
  dynamic_baseline_detection_flag    = true
  baseline_duration                  = 14400000
  positive_baseline_violation_factor = 2.0
  enable_fill_gap                    = true

  metric_configurations = {
    "cpu_usage" = {
      metric_alert_settings = [
        {
          component_name             = "Global_0"
          threshold_alert_upper_bound = "90"
          is_kpi                     = true
          detection_type             = "positive"
        }
      ]
    }
  }
}
```

#### Import

```shell
terraform import insightfinder_metric_project.example my-metric-project
```

---

### insightfinder_log_labels

Manages InsightFinder log label settings for a project as a standalone resource. Use this
*or* the inline [`log_label_settings`](#log_label_settings-block) block on
`insightfinder_project` — not both for the same project.

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `project_name` | String | The name of the project to configure log labels for. Changing this forces a new resource. |
| `label_settings` | List of Block | List of log label settings. |

#### label_settings Block

| Argument | Type | Required | Description |
|---|---|---|---|
| `label_type` | String | Yes | Type of log label — see the table below. |
| `log_label_string` | String | Yes | JSON array string of log labels. For plain-string logs, e.g. `'["ERROR","WARN"]'` or a regex such as `'^\d+$'`. For JSON-structured logs, e.g. `'key=["ERROR","WARN"]'` or `'key=^\d+$'`. |

##### Supported label_type values

| `label_type` | Underlying API field |
|---|---|
| `whitelist` | `whitelist` |
| `trainingWhitelist` | `trainingWhitelist` |
| `blacklist` | `trainingBlacklistLabels` |
| `extractionBlacklist` | `extractionBlacklist` |
| `featurelist` | `featurelist` |
| `incidentlist` | `incidentlist` |
| `triagelist` | `triagelist` |
| `anomalyFeature` | `anomalyFeatureLabels` |
| `dataFilter` | `dataFilterLabels` |
| `patternName` | `patternNameLabels` |
| `patternSignature` | `patternSignatureLabels` |
| `patternMatchRegex` | `patternMatchRegexLabels` |
| `patternIgnoreRegex` | `patternIgnoreRegexLabels` |
| `customAction` | `customActionLabels` |
| `logEventID` | `logEventIDLabels` |
| `logSeverity` | `logSeverityLabels` |
| `logStatusCode` | `logStatusCodeLabels` |
| `alertEventType` | `alertEventTypeLabels` |
| `instanceName` | `instanceNameLabels` |
| `dataQualityCheck` | `dataQualityCheckLabels` |
| `incidentFieldVerification` | `incidentFieldVerificationLabels` |
| `incidentPriority` | `incidentPriorityLabels` |

Any `label_type` not in this table is passed through to the API unchanged. This includes the
additional log-field label types `logSession`, `logComponent`, `logTransactionID`, and
`logCustomParameter`, which are sent to the API under exactly those names.

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | Identifier for the log labels configuration (same as `project_name`). |

#### Example

```hcl
resource "insightfinder_log_labels" "severity" {
  project_name = "application-logs"

  label_settings = [
    {
      label_type       = "logSeverity"
      log_label_string = jsonencode(["ERROR", "WARN", "FATAL"])
    },
    {
      label_type = "whitelist"
      log_label_string = jsonencode([
        {
          type           = "fieldName"
          keyword        = "severity=error|critical|fatal"
          isCritical     = true
          isHotEventOnly = false
        }
      ])
    }
  ]
}
```

#### Import

```shell
terraform import insightfinder_log_labels.example my-project-name
```

---

### insightfinder_jwt_config

Manages InsightFinder JWT configuration for a system.

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `system_name` | String | The name of the system to configure JWT for. Changing this forces a new resource. |
| `jwt_secret` | String, Sensitive | The JWT secret token (minimum 6 characters). |

#### Optional Arguments

| Argument | Type | Description |
|---|---|---|
| `jwt_type` | Integer | The JWT type. Default: `1` (system-level JWT). |

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | Identifier for the JWT configuration (same as `system_name`). |

#### Example

```hcl
resource "insightfinder_jwt_config" "prod" {
  system_name = "Production"
  jwt_secret  = var.jwt_secret
  jwt_type    = 1
}
```

#### Import

```shell
terraform import insightfinder_jwt_config.example Production
```

---

### insightfinder_servicenow

Manages an InsightFinder ServiceNow integration.

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `account` | String | ServiceNow account username. Changing this forces a new resource. |
| `service_host` | String | ServiceNow service host URL. Changing this forces a new resource. |
| `password` | String, Sensitive | ServiceNow account password. |
| `options` | Set of String | ServiceNow integration options (e.g. `Root Cause`). |
| `content_option` | Set of String | ServiceNow content options (e.g. `SUMMARY`). |

> **Note:** the Word-document version of this reference listed `dampening_period` as required
> and `options` / `content_option` as optional. In the current provider, `options` and
> `content_option` are required, there is no `dampening_period` argument, and the equivalent
> timing control is `trigger_window_in_mills`.

#### Optional Arguments

| Argument | Type | Description |
|---|---|---|
| `proxy` | String | Proxy server URL. |
| `app_id` | String | ServiceNow application ID. |
| `app_key` | String, Sensitive | ServiceNow application key. |
| `auth_type` | String | Authentication type — must be `basic` or `oauth`. *Provider default:* `basic`. |
| `system_names` | Set of String | Set of system names to integrate. System IDs are resolved internally from these names. |
| `service_now_field` | String | The ServiceNow ticket field where InsightFinder writes incident analysis data (e.g. `u_probable_cause`). |
| `content_source` | String | Source used to populate the ServiceNow ticket content (e.g. `work_notes`). *Provider default:* `work_notes`. |
| `trigger_window_in_mills` | Integer | Time window in milliseconds controlling when a new ServiceNow ticket is created for a recurring incident (e.g. `604800000` for 7 days). |
| `enable_feedback_collect` | Boolean | Collect feedback from resolved ServiceNow tickets back into InsightFinder. *Provider default:* `false`. |
| `ticket_created_by_source_key` | String | ServiceNow field name used to identify the ticket creator (e.g. `opened_by`). |
| `ticket_created_by_source_value` | String | Expected value of `ticket_created_by_source_key` identifying InsightFinder-created tickets (e.g. `Insight Finder Platform`). |
| `configuration_item` | String | Default ServiceNow CMDB configuration item applied to tickets when no project-level override is set. |
| `department_id` | String | ServiceNow department ID to associate with created tickets. |
| `table_mapping` | Map of String | Mapping of InsightFinder project names to ServiceNow table names (e.g. `{ "my-project" = "incident" }`). |
| `project_configs` | Map of Block | Per-project ServiceNow ticket settings (see below). |

#### project_configs Block

An optional map of project-specific ServiceNow ticket settings. **The map key is the
InsightFinder project name.** Each entry can override the default configuration item and
control which ticket operations are enabled for that project.

| Argument | Type | Description |
|---|---|---|
| `enable_ticket_creation` | Boolean | Enable automatic creation of ServiceNow tickets for incidents in this project. *Provider default:* `false`. |
| `enable_ticket_update` | Boolean | Enable updating existing ServiceNow tickets when new incident data arrives. *Provider default:* `false`. |
| `enable_incident_consolidation_info_update` | Boolean | Enable updating tickets with incident consolidation information. *Provider default:* `false`. |
| `enable_incident_resolve_update` | Boolean | Enable updating tickets when an incident is resolved. *Provider default:* `false`. |
| `configuration_item` | String | ServiceNow CMDB configuration item for this project. Overrides the top-level `configuration_item`. |

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | Identifier for the ServiceNow configuration, formatted `account@service_host`. |

#### Example

```hcl
resource "insightfinder_servicenow" "prod" {
  account      = "svc_insightfinder"
  service_host = "https://acme.service-now.com"
  password     = var.servicenow_password
  auth_type    = "basic"

  options        = ["Root Cause"]
  content_option = ["SUMMARY"]

  system_names            = ["Production"]
  service_now_field       = "u_probable_cause"
  content_source          = "work_notes"
  trigger_window_in_mills = 604800000
  configuration_item      = "AcmeApp"

  table_mapping = {
    "application-logs" = "incident"
  }

  project_configs = {
    "application-logs" = {
      enable_ticket_creation = true
      enable_ticket_update   = true
    }
  }
}
```

#### Import

The import ID is `account@service_host`:

```shell
terraform import insightfinder_servicenow.example svc_insightfinder@https://acme.service-now.com
```

---

### insightfinder_slack

Manages an InsightFinder Slack integration for a system. Allows InsightFinder to post
incident, prediction, and pattern-alert notifications to a Slack channel via an incoming
webhook.

A system may have more than one `insightfinder_slack` resource — each creates a distinct,
independently identified integration (its own webhook and channel), so you can notify
different channels for the same system.

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `system_name` | String | InsightFinder system name to integrate with Slack. |
| `webhook` | String, Sensitive | Slack incoming webhook URL. |
| `channel_name` | String | Slack channel to send notifications to (e.g. `#my-channel`). |
| `options` | Set of String | Notification types to send to Slack (e.g. `Detected Incident`, `Predicted Incident`, `New Pattern Alert`, `Missing Monitoring Data`). |

#### Optional Arguments

| Argument | Type | Description |
|---|---|---|
| `priority_upgrade_channel` | String | Slack channel to notify when an incident's priority is upgraded. *Provider default:* `""`. |
| `priority_upgrade_webhook` | String, Sensitive | Slack webhook URL to notify when an incident's priority is upgraded. *Provider default:* `""`. |
| `disable_slack_for_non_insightfinder_incidents` | Boolean | Suppress Slack notifications for incidents that did not originate from InsightFinder. *Provider default:* `false`. |
| `project_configs` | List of Block | Per-project Slack notification overrides (see below). |

#### project_configs Block

A list of per-project overrides. A project may appear more than once to notify multiple
channels.

| Argument | Type | Required | Description |
|---|---|---|---|
| `project_name` | String | Yes | InsightFinder project name. |
| `channel` | String | No | Slack channel override for this project. Defaults to the top-level `channel_name` when empty. *Provider default:* `""`. |
| `webhook` | String, Sensitive | No | Slack webhook override for this project. Defaults to the top-level `webhook` when empty. *Provider default:* `""`. |
| `options` | Set of String | No | Notification types to send for this project. Defaults to the top-level `options` when empty. |
| `enable_consolidation_info_update` | Boolean | No | Send incident consolidation info updates for this project. *Provider default:* `false`. |
| `priority_levels` | List of Number | No | Priority levels that trigger notifications for this project (e.g. `[1, 2, 3]`). |
| `rule` | Block | No | Optional match rule that filters which alerts are sent for this project (see below). |

##### project_configs → rule

| Argument | Type | Description |
|---|---|---|
| `type` | String | Rule type (e.g. `fieldName` or `content`). |
| `keyword` | String | Rule match expression. |

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | API-generated identifier (account ID) for this Slack integration. Cannot be set by the practitioner. |

#### Example

```hcl
resource "insightfinder_slack" "full" {
  system_name  = "Production"
  webhook      = var.slack_webhook
  channel_name = "#aiops-incidents"

  options = [
    "Detected Incident",
    "Predicted Incident",
    "New Pattern Alert",
    "Missing Monitoring Data",
  ]

  priority_upgrade_channel = "#aiops-priority-upgrades"
  priority_upgrade_webhook = var.slack_webhook

  disable_slack_for_non_insightfinder_incidents = true

  project_configs = [
    {
      project_name                     = "my-project"
      channel                          = "my-project-incidents"
      options                          = ["Detected Incident", "Predicted Incident"]
      enable_consolidation_info_update = true
      priority_levels                  = [1, 2, 3]
      rule = {
        type    = "fieldName"
        keyword = "alert->core->monitored_item=.*"
      }
    },
  ]
}
```

#### Behavior Notes

- `system_name` is automatically resolved to a system ID internally.
- `options` is a set — order does not matter and will not cause plan drift.
- If an integration already exists for the exact same system, webhook, and channel (for
  example it was created outside Terraform, or state was lost and `terraform apply` runs
  again), creating the resource **adopts** the existing integration and brings it in line
  with the given configuration rather than creating a duplicate.
- Updates locate the integration to modify by matching its system, webhook, and channel from
  *before* the change — not the stored `id` — against InsightFinder's list of integrations,
  then apply the new values to whatever account is found. This mirrors how the InsightFinder
  UI resolves which integration to update.

#### Import

```shell
terraform import insightfinder_slack.example ec2dc73e-9f7b-4d59-ab7b-07a7434d2040
```

---

### insightfinder_system_settings

Manages InsightFinder system-level settings: knowledge base configuration, notification and
alert settings, and miscellaneous system framework settings. This resource targets a
**system**, not a project.

> **Lifecycle note:** system settings cannot be deleted through the API. Removing the
> resource from Terraform removes it from state only; the settings are left unchanged on the
> server.

#### Required Arguments

| Argument | Type | Description |
|---|---|---|
| `system_name` | String | Display name of the system. Used to resolve the system ID via the system framework API. Changing this forces a new resource. |

#### knowledgebase_settings Block

Optional block. Controls the global Knowledge Base (KB) and Incident Prediction engine for
the system.

##### Global KB Fields

| Argument | Type | Description |
|---|---|---|
| `enable_global_knowledge_base` | Boolean | Enable the global knowledge base for the system. |
| `composite_valid_threshold` | Integer | Composite valid threshold in milliseconds. |
| `timeline_top_k` | Integer | Number of top timeline entries to retain. |
| `enable_ignore_instance_prediction` | Boolean | When true, the KB ignores instance-level prediction data. |
| `prediction_source` | Integer | Prediction source type (0 = default, 1 = custom). |
| `share_system_type` | Integer | Share system type for the KB. |
| `action_execution_time` | Integer | Action execution time in minutes. |
| `auto_fix_validation_window` | Integer | Validation window used by the auto-fix feature. |
| `filter_self_to_self` | Boolean | Filter out self-to-self KB entries. |
| `rule_source_type` | Integer | Rule source type (0 = default). |
| `satellite_system_set` | String (JSON) | JSON array of satellite systems linked to this system's knowledge base. Each entry requires a `systemPartitionKey` object (`userName`, `systemName`, `envName`) and a `replay` boolean. |

```hcl
satellite_system_set = jsonencode([
  {
    systemPartitionKey = { userName = "u", systemName = "<id>", envName = "All" }
    replay             = false
  }
])
```

##### Incident Prediction Fields

| Argument | Type | Description |
|---|---|---|
| `rule_active_threshold` | Float | Minimum probability to promote a causal rule to Active status (0.0–1.0). |
| `rule_inactive_threshold` | Float | Probability below which a rule is demoted (0.0–1.0). Must be `<= rule_active_threshold`. |
| `rule_active_condition` | Integer | Prerequisite a rule must meet before generating alerts (0 = unfiltered, 1 = verified only). |
| `false_positive_tolerance` | Integer | False positive count tolerated before a rule is deactivated. |
| `kb_training_length` | Integer | Length of the KB training window in milliseconds. |
| `tolerance` | Float | Tolerance value for incident prediction calculations. |
| `enable_insensitive_rule_matching` | Boolean | Enable case-insensitive rule matching in the KB. |

#### notifications_settings Block

Optional block. Controls all notification and alert settings for the system. Persisted via
`/api/external/v2/healthviewsetting` plus separate sub-APIs for system-down, instance-down,
and insights-report notifications.

##### Health View / General Fields

| Argument | Type | Description |
|---|---|---|
| `order` | Integer | Display order for the system in the health view dashboard. |
| `hide_flag` | Boolean | Hide this system from the health view. |
| `aggregation_interval` | Integer | Aggregation interval in minutes for health view metrics. |
| `enable_splunk_export` | Boolean | Enable exporting system data to Splunk. |
| `incident_count_threshold` | String (JSON) | JSON map of project names (format `"ProjectName@username"`) to incident count thresholds. Example: `jsonencode({"MyProject@admin": 5})`, or `jsonencode({})` to clear. |
| `assignment_map` | String (JSON) | JSON map of zone/component keys to assignee lists. Each value may contain `emailAssignees`, `jiraAssignees`, and `serviceNowAssignees` arrays. |
| `alert_health_score` | Float | Health score threshold (0.0–1.0) below which an alert is triggered. |
| `alert_frequency` | Integer | Alert frequency setting. |
| `incident_dampening_window` | Integer | Dampening window for incident notifications, in milliseconds. |
| `ticket_open_time` | Integer | Time window (ms) to keep a ticket open after an incident resolves. |
| `component_level_incident_consolidation` | Boolean | Enable component-level incident consolidation. |
| `component_level_dampening` | Boolean | Enable component-level dampening. |
| `enabled_consolidation_algorithms` | List of String | Consolidation algorithms to enable. Valid values: `derivedIncidents`, `rcaChain`, `contentBased`, `metricInstanceTimestamp`. |
| `max_notification_delay_tolerance` | Integer | Maximum notification delay tolerance in milliseconds. |
| `metric_co_occurrence_buffer_ms` | Integer | Metric co-occurrence buffer window in milliseconds. |

##### Email Alert Fields

| Argument | Type | Description |
|---|---|---|
| `prediction_email` | String | Email address for incident prediction notifications. |
| `alert_email` | String | Email address for general alert notifications. |
| `health_alert_email` | String | Email address for health alert notifications. |
| `incident_detection_email` | String | Email address for incident detection notifications. |
| `root_cause_email` | String | Email address for root cause analysis notifications. |
| `email_dampening_period` | Integer | Dampening period for health alert emails, in milliseconds. |
| `alerts_email_dampening_period` | Integer | Dampening period for general alert emails, in milliseconds. |
| `prediction_email_dampening_period` | Integer | Dampening period for prediction emails, in milliseconds. |
| `enable_system_down_email_alert` | Boolean | Enable email alert when the system is down. |
| `only_send_with_rca` | Boolean | Only send notifications when root cause analysis data is available. |
| `enable_incident_prediction_email_alert` | Boolean | Enable email alert for incident predictions. |
| `enable_incident_detection_email_alert` | Boolean | Enable email alert for incident detections. |
| `enable_alerts_email` | Boolean | Enable general alert emails. |
| `enable_health_email_alert` | Boolean | Enable health score email alerts. |
| `enable_root_cause_email_alert` | Boolean | Enable email alerts including root cause analysis results. |

##### system_down_notification Block

Optional nested block. Managed via a dedicated system-down API
(`/api/external/v2/systemdownsetting`).

| Argument | Type | Description |
|---|---|---|
| `enable_system_down_email_alert` | Boolean | Enable email alert when the system is down. |
| `email_dampening_period` | Integer | Dampening period for system-down emails, in milliseconds. |
| `email_set` | List of String | Email addresses to notify when the system is down. |

##### daily_report_notification Block

Optional nested block. Daily insights report notification settings.

| Argument | Type | Description |
|---|---|---|
| `enable_insights_report` | Boolean | Enable the daily insights report email. |
| `email_set` | List of String | Email addresses to receive the daily insights report. |

##### weekly_report_notification Block

Optional nested block. Weekly insights report notification settings.

| Argument | Type | Description |
|---|---|---|
| `enable_insights_report` | Boolean | Enable the weekly insights report email. |
| `email_set` | List of String | Email addresses to receive the weekly insights report. |

##### instance_down_notification Block

Optional list of nested blocks. Instance-down notification settings, one entry per project.

| Argument | Type | Required | Description |
|---|---|---|---|
| `project_name` | String | Yes | The project to configure instance-down notifications for. |
| `instance_down_enable` | Boolean | No | Enable instance-down email alerts for this project. |
| `instance_down_dampening` | Integer | No | Dampening period for instance-down alerts, in milliseconds. |
| `instance_down_threshold` | Integer | No | Threshold (milliseconds) before an instance is considered down. |
| `instance_down_report_number` | Integer | No | Number of instance-down events to include in the report. |
| `instance_down_emails` | List of String | No | Email addresses to notify when instances are down. |

##### project_level_dampening_windows Block

Optional set of nested blocks. Each block overrides the system-level incident dampening
window for a specific source → target project pair.

| Argument | Type | Required | Description |
|---|---|---|---|
| `source_project` | String | Yes | The source project name. |
| `target_project` | String | Yes | The target project name. |
| `duration` | Integer | Yes | Dampening duration in milliseconds. |
| `source_customer` | String | No | Username of the source project owner. Defaults to the provider username. |
| `target_customer` | String | No | Username of the target project owner. Defaults to the provider username. |
| `similarity_threshold` | Float | No | Similarity threshold (`st`) for this dampening window. |

##### custom_consolidation_rules Block

Optional list of nested blocks defining custom incident consolidation rules.

| Argument | Type | Description |
|---|---|---|
| `project_entries` | List of Block | Projects and their matching conditions for this consolidation rule. |
| `field_correlations` | List of Block | Field correlations mapping project fields across the consolidation rule. |

**`custom_consolidation_rules` → `project_entries`**

| Argument | Type | Required | Description |
|---|---|---|---|
| `project_name` | String | Yes | The project name. |
| `conditions` | List of Block | No | Matching conditions for this project entry. |

**`custom_consolidation_rules` → `project_entries` → `conditions`**

| Argument | Type | Required | Description |
|---|---|---|---|
| `type` | String | Yes | Condition type: `fieldName` or `content`. |
| `keyword` | String | Yes | The keyword or field expression to match. |

**`custom_consolidation_rules` → `field_correlations`**

| Argument | Type | Description |
|---|---|---|
| `project_field_keys` | List of Block | List of project field key mappings. |

**`custom_consolidation_rules` → `field_correlations` → `project_field_keys`**

| Argument | Type | Required | Description |
|---|---|---|---|
| `project_name` | String | Yes | The project name. |
| `type` | String | Yes | Field type: `fieldName` or `content`. |
| `field_key` | String | No | The field key path. Empty or omitted for the `content` type. |

##### metric_log_consolidation_configs Block

Optional list of nested blocks configuring metric-to-log project consolidation.

| Argument | Type | Required | Description |
|---|---|---|---|
| `metric_project_name` | String | Yes | The metric project name. |
| `log_project_name` | String | Yes | The log project name. |
| `field_keys` | List of String | No | Field keys used for consolidation. |

#### miscellaneous_settings Block

Optional block controlling miscellaneous system framework settings.

| Argument | Type | Description |
|---|---|---|
| `healthview_longterm` | Boolean | Enable long-term storage mode for the system health view. |
| `should_auto_share` | Boolean | Enable automatic sharing of system data. |
| `rootcause_reverse_entry_filter_threshold` | Integer | Threshold (0–100) for the root cause reverse entry filter. |
| `enable_composite_timeline` | Boolean | Enable the composite timeline view for the system. |

#### Read-Only Attributes

| Attribute | Type | Description |
|---|---|---|
| `id` | String | Identifier for the system settings (same as `system_name`). |

#### Example

```hcl
resource "insightfinder_system_settings" "prod" {
  system_name = "Production"

  knowledgebase_settings = {
    enable_global_knowledge_base = true
    timeline_top_k               = 10
    rule_active_threshold        = 0.7
    rule_inactive_threshold      = 0.5
    rule_active_condition        = 1
    kb_training_length           = 604800000
  }

  notifications_settings = {
    order                     = 1
    aggregation_interval      = 10
    alert_health_score        = 0.8
    incident_dampening_window = 3600000

    enable_incident_detection_email_alert = true
    incident_detection_email              = "oncall@example.com"

    enabled_consolidation_algorithms = ["derivedIncidents", "rcaChain"]

    system_down_notification = {
      enable_system_down_email_alert = true
      email_dampening_period         = 3600000
      email_set                      = ["oncall@example.com"]
    }

    daily_report_notification = {
      enable_insights_report = true
      email_set              = ["reports@example.com"]
    }

    project_level_dampening_windows = [
      {
        source_project = "app-logs"
        target_project = "app-metrics"
        duration       = 1800000
      }
    ]
  }

  miscellaneous_settings = {
    healthview_longterm       = true
    enable_composite_timeline = true
  }
}
```

#### Import

```shell
terraform import insightfinder_system_settings.example Production
```

---

## Data Sources

### insightfinder_project (data source)

Fetches an InsightFinder project.

| Argument / Attribute | Type | Mode | Description |
|---|---|---|---|
| `project_name` | String | Required | The name of the project to fetch. |
| `id` | String | Computed | Identifier for the project. |
| `project_display_name` | String | Computed | The display name for the project. |
| `c_value` | Integer | Computed | The C value. |
| `p_value` | Float | Computed | The P value. |

```hcl
data "insightfinder_project" "existing" {
  project_name = "application-logs"
}
```

### insightfinder_systems

Fetches the list of systems from InsightFinder.

| Attribute | Type | Description |
|---|---|---|
| `id` | String | Placeholder identifier for the data source. |
| `systems` | List of Object | List of systems. Each element has `system_id` and `system_name`. |

**`systems` element**

| Attribute | Type | Description |
|---|---|---|
| `system_id` | String | The unique identifier for the system. |
| `system_name` | String | The name of the system. |

```hcl
data "insightfinder_systems" "all" {}

output "system_names" {
  value = [for s in data.insightfinder_systems.all.systems : s.system_name]
}
```

---

## Appendix A — Import ID Reference

| Resource | Import ID | Example |
|---|---|---|
| `insightfinder_project` | Project name | `my-project-name` |
| `insightfinder_metric_project` | Project name | `my-metric-project` |
| `insightfinder_log_labels` | Project name | `my-project-name` |
| `insightfinder_jwt_config` | System name | `Production` |
| `insightfinder_system_settings` | System name | `Production` |
| `insightfinder_servicenow` | `account@service_host` | `svc_user@https://acme.service-now.com` |
| `insightfinder_slack` | API-generated account ID | `ec2dc73e-9f7b-4d59-ab7b-07a7434d2040` |

---

## Appendix B — Force-Replacement Arguments

Changing any of these arguments destroys and recreates the resource:

| Resource | Argument(s) |
|---|---|
| `insightfinder_project` | `project_name` |
| `insightfinder_metric_project` | `project_name` |
| `insightfinder_log_labels` | `project_name` |
| `insightfinder_jwt_config` | `system_name` |
| `insightfinder_system_settings` | `system_name` |
| `insightfinder_servicenow` | `account`, `service_host` |
