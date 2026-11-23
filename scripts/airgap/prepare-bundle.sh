#!/usr/bin/env bash
# Builds a bundle that runs the provider on an airgapped machine. Run from the repository root:
#
#   scripts/airgap/prepare-bundle.sh              # bundles examples/airgap-data-quality
#   ARCH=arm64 scripts/airgap/prepare-bundle.sh   # target says aarch64 for uname -m
#   CONFIG_DIR=examples/airgap-<name> scripts/airgap/prepare-bundle.sh
#
# Produces dist/airgap/if-<config folder name>-bundle.tgz, e.g. dist/airgap/if-airgap-data-quality-bundle.tgz.
set -euo pipefail

ARCH=${ARCH:-amd64}                            # CPU of the airgapped machine: amd64 or arm64
VERSION=${VERSION:-$(cat VERSION)}
CONFIG_DIR=${CONFIG_DIR:-examples/airgap-data-quality} # the .tf (and optional .json) files to bundle
OUT_DIR=${OUT_DIR:-dist/airgap}

NAME=if-$(basename "$CONFIG_DIR")-bundle
BUNDLE=$OUT_DIR/$NAME
PLUGIN_DIR=$BUNDLE/terraform.d/plugins/registry.terraform.io/insightfinder/insightfinder/$VERSION/linux_$ARCH

rm -rf "$BUNDLE" "$BUNDLE.tgz"
mkdir -p "$PLUGIN_DIR"

CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build -trimpath \
  -ldflags "-s -w -X main.version=$VERSION" \
  -o "$PLUGIN_DIR/terraform-provider-insightfinder_v$VERSION" .

shopt -s nullglob
cp "$CONFIG_DIR"/*.tf "$CONFIG_DIR"/*.json scripts/airgap/run-if-terraform.sh "$BUNDLE/"

tar czf "$BUNDLE.tgz" -C "$OUT_DIR" "$NAME"

echo "Send this file to the airgapped machine: $BUNDLE.tgz"
