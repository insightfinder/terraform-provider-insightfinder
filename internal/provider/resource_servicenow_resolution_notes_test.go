// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestOptionalStringFromAPI(t *testing.T) {
	tests := []struct {
		name  string
		api   string
		prior types.String
		want  types.String
	}{
		{name: "stored value", api: "x", prior: types.StringNull(), want: types.StringValue("x")},
		{name: "unset stays null", api: "", prior: types.StringNull(), want: types.StringNull()},
		{name: "empty string in config is kept", api: "", prior: types.StringValue(""), want: types.StringValue("")},
		{name: "value cleared outside terraform", api: "", prior: types.StringValue("x"), want: types.StringNull()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optionalStringFromAPI(tt.api, tt.prior); !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
