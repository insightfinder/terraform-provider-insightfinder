// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const backendDefaultSkipRegex = `(?i)\btest (inc|alert|incident|ticket)\b|^test$`

func serviceNowListServer(entry map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(map[string]interface{}{"success": true, "extServiceAllInfo": []interface{}{entry}})
		_, _ = w.Write(body)
	}))
}

func TestGetServiceNowConfig_ResolutionNotesFromRawConfigs(t *testing.T) {
	configs, _ := json.Marshal(map[string]interface{}{
		"systemIds":                        []string{},
		"enableResolutionNotesLearning":    true,
		"resolutionNotesPatternStripRegex": "Expired",
	})
	server := serviceNowListServer(map[string]interface{}{
		"account":                          "admin",
		"service_host":                     "https://sn.example.com",
		"enableResolutionNotesLearning":    true,
		"resolutionNotesSkipRegex":         backendDefaultSkipRegex,
		"resolutionNotesPatternStripRegex": "Expired",
		"configs":                          string(configs),
	})
	defer server.Close()

	c, _ := NewClient(server.URL, "test_user", "test_key")
	got, err := c.GetServiceNowConfig("admin", "https://sn.example.com", "test_user")
	if err != nil || got == nil {
		t.Fatalf("got %v, err %v", got, err)
	}
	if !got.EnableResolutionNotesLearning {
		t.Error("EnableResolutionNotesLearning = false, want true")
	}
	if got.ResolutionNotesSkipRegex != "" {
		t.Errorf("ResolutionNotesSkipRegex = %q, want empty: the root value is the backend default, not a stored one", got.ResolutionNotesSkipRegex)
	}
	if got.ResolutionNotesPatternStripRegex != "Expired" {
		t.Errorf("ResolutionNotesPatternStripRegex = %q, want Expired", got.ResolutionNotesPatternStripRegex)
	}
}

func TestCreateOrUpdateServiceNowConfig_ResolutionNotesForm(t *testing.T) {
	tests := []struct {
		name        string
		verify      bool
		config      ServiceNowConfig
		wantPresent map[string]string
		wantAbsent  []string
	}{
		{
			name:        "learning off, regexes unset",
			config:      ServiceNowConfig{Account: "a", ServiceHost: "h"},
			wantPresent: map[string]string{"enableResolutionNotesLearning": "false"},
			wantAbsent:  []string{"resolutionNotesSkipRegex", "resolutionNotesPatternStripRegex"},
		},
		{
			name:   "learning on, both regexes",
			config: ServiceNowConfig{Account: "a", ServiceHost: "h", EnableResolutionNotesLearning: true, ResolutionNotesSkipRegex: `(?i)^n/a$`, ResolutionNotesPatternStripRegex: `\d+`},
			wantPresent: map[string]string{
				"enableResolutionNotesLearning":    "true",
				"resolutionNotesSkipRegex":         `(?i)^n/a$`,
				"resolutionNotesPatternStripRegex": `\d+`,
			},
		},
		{
			name:       "verify call carries none of them",
			verify:     true,
			config:     ServiceNowConfig{Account: "a", ServiceHost: "h", EnableResolutionNotesLearning: true, ResolutionNotesSkipRegex: "x"},
			wantAbsent: []string{"enableResolutionNotesLearning", "resolutionNotesSkipRegex", "resolutionNotesPatternStripRegex"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Fatalf("Failed to parse form: %v", err)
				}
				for k, v := range tt.wantPresent {
					if got, ok := r.PostForm[k]; !ok || got[0] != v {
						t.Errorf("form %s = %v, want %q", k, got, v)
					}
				}
				for _, k := range tt.wantAbsent {
					if _, ok := r.PostForm[k]; ok {
						t.Errorf("form %s present, want absent", k)
					}
				}
				_, _ = w.Write([]byte(`{"success":true}`))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			cfg := tt.config
			if err := c.CreateOrUpdateServiceNowConfig(&cfg, "test_user", tt.verify); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
