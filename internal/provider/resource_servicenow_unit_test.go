// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/insightfinder/terraform-provider-insightfinder/internal/provider/client"
)

// TestProjectConfigsToTF_AllFieldsRoundTrip guards against the class of bug fixed in
// commit 3ae7c50 ("Fix missing enable_metric_value_update in project_configs read path"):
// project_configs is the only attribute in this provider built from two independently
// hand-copied attr.Type maps (projectConfigAttrTypes and the local attrTypes literal
// inside projectConfigsToTF) describing the same object shape, rather than a single
// shared helper. If a new field is added to the ServiceNowProjectConfig struct and the
// object-value map in projectConfigsToTF without adding it to both attr.Type maps,
// types.ObjectValue returns an error diagnostic and every Read produces an "Extra Object
// Attribute Value" error on plan/refresh for every project — this test fails immediately
// instead of only surfacing against a real InsightFinder API + terraform plan.
func TestProjectConfigsToTF_AllFieldsRoundTrip(t *testing.T) {
	ctx := context.Background()
	configs := map[string]client.ServiceNowProjectConfig{
		"my-project": {
			EnableTicketCreation:                  true,
			EnableTicketUpdate:                    true,
			EnableIncidentConsolidationInfoUpdate: true,
			EnableIncidentResolveUpdate:           true,
			EnableIncidentFieldSync:               true,
			EnableMetricValueUpdate:               true,
			ConfigurationItem:                     "some-ci",
		},
	}

	result, diags := projectConfigsToTF(ctx, configs, types.MapNull(projectConfigAttrTypes()))
	if diags.HasError() {
		t.Fatalf("projectConfigsToTF returned unexpected error diagnostics: %v", diags)
	}

	roundTripped, err := projectConfigsFromTF(ctx, result)
	if err != nil {
		t.Fatalf("projectConfigsFromTF returned unexpected error: %v", err)
	}

	got, ok := roundTripped["my-project"]
	if !ok {
		t.Fatalf("expected project %q in round-tripped map, got %v", "my-project", roundTripped)
	}
	want := configs["my-project"]
	if got != want {
		t.Fatalf("round-tripped project config = %+v, want %+v", got, want)
	}
}
