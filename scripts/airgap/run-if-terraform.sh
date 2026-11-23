#!/usr/bin/env bash
# First run of the bundled configuration on the airgapped machine, from inside the unpacked
# bundle folder. Afterwards use terraform plan / apply / destroy from the same folder.
set -euo pipefail

export IF_BASE_URL=${IF_BASE_URL:-"https://<internal InsightFinder URL>"}
export IF_USERNAME=${IF_USERNAME:-"<your user name>"}
export IF_LICENSE_KEY=${IF_LICENSE_KEY:-"<your license key>"}
export CHECKPOINT_DISABLE=1                # don't try to reach the internet for version checks
# export SSL_CERT_FILE=/path/to/ca.pem    # only if InsightFinder uses a private certificate authority

FIRST=${FIRST:-insightfinder_metric_project.l2m}  # created before the rest; empty to skip

cd "$(dirname "$0")"

# Installs the provider from ./terraform.d (no internet needed)
terraform init

# Everything that will be created or changed
terraform plan

# Create the metric project first, so the log project can refer to it
if [ -n "$FIRST" ]; then
  terraform apply -target="$FIRST"
fi

# Create or update everything else
terraform apply

# Should say "No changes.": InsightFinder now matches the configuration
terraform plan
