# Copyright (c) InsightFinder Inc.
# SPDX-License-Identifier: Apache-2.0

# The configuration scripts/airgap/prepare-bundle.sh packs for an airgapped machine: a data
# quality log project and the metric project its log-to-metric (L2M) data goes to. Adjust the
# names and system below before building the bundle, or after unpacking it.

terraform {
  required_providers {
    insightfinder = { source = "insightfinder/insightfinder" }
  }
}

provider "insightfinder" {} # credentials come from the IF_* environment variables

locals {
  system_name = "<system name>"
}

# Metric project that receives the L2M data quality metrics
resource "insightfinder_metric_project" "l2m" {
  project_name = "<dq-l2m-metric-project>"
  system_name  = local.system_name

  project_creation_config = {
    data_type          = "Metric"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  # Disable baseline anomaly filter ("Enable anomaly filter for baseline" in the UI)
  filter_by_anomaly_in_baseline_generation = false

  # Negative baseline detection is set per metric, so it can only be added once L2M has created
  # the metrics. Then add one entry per metric and apply again:
  # metric_configurations = {
  #   "<metric name>" = { detection_type = "both" }
  # }
}

# Data quality log project
resource "insightfinder_project" "dq_log" {
  project_name = "<dq-log-project>"
  system_name  = local.system_name

  project_creation_config = {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
    insight_agent_type = "Custom"
  }

  anomaly_detection_mode             = -1   # disable log anomaly detection
  log_detection_min_count            = 1    # minimum batch size
  log_to_metric_instance_aggregation = true # log to metric instance aggregation

  # Automatic historical log processing (automatic reruns). Needs an InsightFinder backend with
  # II-25564; remove this line on an older backend.
  automatic_historical_log_processing = true

  # Training allowlist on the "group" field. Confirm the keyword against an existing project's UI.
  log_label_settings = [
    {
      label_type       = "trainingWhitelist"
      log_label_string = jsonencode([{ type = "fieldName", keyword = "group" }])
    }
  ]

  depends_on = [insightfinder_metric_project.l2m]
}
