// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var l2mSettingObjType = types.ObjectType{AttrTypes: l2mSettingAttrTypes}
var l2mJsonParserObjType = types.ObjectType{AttrTypes: l2mJsonParserAttrTypes}

// l2mTestJsonParser builds a json parser with the configured fields set and
// every other field set to unset (null in state, unknown in plan).
func l2mTestJsonParser(unset func(attr.Type) attr.Value) attr.Value {
	attrs := map[string]attr.Value{
		"additional_metric_name": types.StringValue(""),
		"aggregation_mode":       types.Int64Value(0),
		"aggregation_period":     types.Int64Value(0),
		"data_filter":            types.StringValue("alert->bic->qualified=true"),
		"grouping_by_component":  types.BoolValue(false),
		"instance_name_key":      types.StringValue("apikey->name"),
		"metric_name":            types.StringValue("Alert Volume"),
		"operation":              types.Int64Value(1),
	}
	for k, t := range l2mJsonParserAttrTypes {
		if _, ok := attrs[k]; !ok {
			attrs[k] = unset(t)
		}
	}
	return types.ObjectValueMust(l2mJsonParserAttrTypes, attrs)
}

func l2mTestSetting(unset func(attr.Type) attr.Value) types.Set {
	parsers := types.ListValueMust(l2mJsonParserObjType, []attr.Value{l2mTestJsonParser(unset)})
	elem := types.ObjectValueMust(l2mSettingAttrTypes, map[string]attr.Value{
		"metric_project_name": types.StringValue("atom-L2M"),
		"json_flag":           types.BoolValue(true),
		"enable_mapping":      types.BoolValue(false),
		"regexs":              unset(l2mSettingAttrTypes["regexs"]),
		"json_parsers":        parsers,
	})
	return types.SetValueMust(l2mSettingObjType, []attr.Value{elem})
}

func l2mNullOf(t attr.Type) attr.Value {
	switch tt := t.(type) {
	case types.ObjectType:
		return types.ObjectNull(tt.AttrTypes)
	case types.ListType:
		return types.ListNull(tt.ElemType)
	}
	switch t {
	case types.StringType:
		return types.StringNull()
	case types.Int64Type:
		return types.Int64Null()
	case types.BoolType:
		return types.BoolNull()
	}
	panic("unexpected type")
}

func l2mUnknownOf(t attr.Type) attr.Value {
	switch tt := t.(type) {
	case types.ObjectType:
		return types.ObjectUnknown(tt.AttrTypes)
	case types.ListType:
		return types.ListUnknown(tt.ElemType)
	}
	switch t {
	case types.StringType:
		return types.StringUnknown()
	case types.Int64Type:
		return types.Int64Unknown()
	case types.BoolType:
		return types.BoolUnknown()
	}
	panic("unexpected type")
}

func TestL2MSettingsPlanModifier_UnchangedEntryHasNoDiff(t *testing.T) {
	ctx := context.Background()
	state := l2mTestSetting(l2mNullOf)
	plan := l2mTestSetting(l2mUnknownOf)

	req := planmodifier.SetRequest{StateValue: state, PlanValue: plan, ConfigValue: state}
	resp := &planmodifier.SetResponse{PlanValue: plan}
	l2mSettingsPlanModifier{}.PlanModifySet(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.PlanValue.Equal(state) {
		t.Fatalf("expected plan to equal state\nplan:  %s\nstate: %s", resp.PlanValue, state)
	}
}

func TestL2MSettingsPlanModifier_NewEntryKeepsUnknowns(t *testing.T) {
	ctx := context.Background()
	plan := l2mTestSetting(l2mUnknownOf)
	other := types.SetValueMust(l2mSettingObjType, []attr.Value{
		types.ObjectValueMust(l2mSettingAttrTypes, map[string]attr.Value{
			"metric_project_name": types.StringValue("other-L2M"),
			"json_flag":           types.BoolNull(),
			"enable_mapping":      types.BoolNull(),
			"regexs":              types.ListNull(types.ObjectType{AttrTypes: l2mRegexParserAttrTypes}),
			"json_parsers":        types.ListNull(l2mJsonParserObjType),
		}),
	})

	req := planmodifier.SetRequest{StateValue: other, PlanValue: plan, ConfigValue: plan}
	resp := &planmodifier.SetResponse{PlanValue: plan}
	l2mSettingsPlanModifier{}.PlanModifySet(ctx, req, resp)

	if !resp.PlanValue.Equal(plan) {
		t.Fatalf("expected unmatched entry to be left as planned")
	}
}

func TestL2MUnknownsToNull(t *testing.T) {
	ctx := context.Background()
	got := l2mUnknownsToNull(ctx, l2mTestSetting(l2mUnknownOf))
	if !got.Equal(l2mTestSetting(l2mNullOf)) {
		t.Fatalf("expected unknowns to become nulls, got %s", got)
	}
}
