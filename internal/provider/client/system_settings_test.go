// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ---- GetGlobalKBSetting ----

func TestGetGlobalKBSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		statusCode   int
		responseBody string
		expectError  bool
		expectNil    bool
		expectedEGB  bool
	}{
		{
			name:       "successful get",
			systemID:   "sys-abc123",
			statusCode: 200,
			responseBody: `[{
				"systemId": "sys-abc123",
				"enableGlobalKnowledgeBase": true,
				"compositeValidThreshold": 86400000,
				"timelineTopK": 5,
				"enableIgnoreInstancePrediction": false,
				"predictionSource": 0,
				"shareSystemType": 0,
				"actionExecutionTime": 30,
				"autoFixValidationWindow": 10,
				"filterSelfToSelf": true,
				"ruleSourceType": 0
			}]`,
			expectError: false,
			expectNil:   false,
			expectedEGB: true,
		},
		{
			name:         "not found 404",
			systemID:     "sys-notfound",
			statusCode:   404,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "no content 204",
			systemID:     "sys-empty",
			statusCode:   204,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "empty array",
			systemID:     "sys-abc123",
			statusCode:   200,
			responseBody: `[]`,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
			expectNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				if r.Header.Get("X-User-Name") == "" {
					t.Error("Expected X-User-Name header")
				}
				if r.Header.Get("X-API-Key") == "" {
					t.Error("Expected X-API-Key header")
				}
				params := r.URL.Query()
				if params.Get("customerName") == "" {
					t.Error("Expected customerName query param")
				}
				if params.Get("systemIds") == "" {
					t.Error("Expected systemIds query param")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			setting, err := c.GetGlobalKBSetting(tt.systemID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && setting != nil {
				t.Error("Expected nil but got non-nil setting")
			}
			if !tt.expectNil && setting == nil && !tt.expectError {
				t.Error("Expected non-nil setting but got nil")
			}
			if setting != nil {
				if setting.EnableGlobalKnowledgeBase != tt.expectedEGB {
					t.Errorf("Expected EnableGlobalKnowledgeBase=%v, got %v", tt.expectedEGB, setting.EnableGlobalKnowledgeBase)
				}
				if setting.SystemID != tt.systemID {
					t.Errorf("Expected SystemID=%s, got %s", tt.systemID, setting.SystemID)
				}
			}
		})
	}
}

func TestSetGlobalKBSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		setting      *GlobalKBSetting
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name:     "successful set",
			systemID: "sys-abc123",
			setting: &GlobalKBSetting{
				EnableGlobalKnowledgeBase: true,
				CompositeValidThreshold:   86400000,
				TimelineTopK:              5,
			},
			statusCode:   200,
			responseBody: `{"success": true}`,
			expectError:  false,
		},
		{
			name:         "api returns success false",
			systemID:     "sys-abc123",
			setting:      &GlobalKBSetting{},
			statusCode:   200,
			responseBody: `{"success": false, "message": "not found"}`,
			expectError:  true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			setting:      &GlobalKBSetting{},
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("Failed to parse form: %v", err)
				}
				if r.FormValue("customerName") == "" {
					t.Error("Expected customerName in form data")
				}
				if r.FormValue("settingModels") == "" {
					t.Error("Expected settingModels in form data")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			err := c.SetGlobalKBSetting(tt.systemID, tt.setting)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// ---- GetIncidentPredictionSetting ----

func TestGetIncidentPredictionSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		statusCode   int
		responseBody string
		expectError  bool
		expectNil    bool
		expectedRAT  float64
	}{
		{
			name:       "successful get",
			systemID:   "sys-abc123",
			statusCode: 200,
			responseBody: `[{
				"systemId": "sys-abc123",
				"ruleActiveThreshold": 0.8,
				"ruleInactiveThreshold": 0.2,
				"ruleActiveCondition": 1,
				"falsePositiveTolerance": 3,
				"kbTrainingLength": 604800000,
				"tolerance": 0.1,
				"enableInsensitiveRuleMatching": false
			}]`,
			expectError: false,
			expectNil:   false,
			expectedRAT: 0.8,
		},
		{
			name:         "not found",
			systemID:     "sys-notfound",
			statusCode:   404,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "empty array",
			systemID:     "sys-abc123",
			statusCode:   200,
			responseBody: `[]`,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			setting, err := c.GetIncidentPredictionSetting(tt.systemID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && setting != nil {
				t.Error("Expected nil but got non-nil setting")
			}
			if !tt.expectNil && setting == nil && !tt.expectError {
				t.Error("Expected non-nil setting but got nil")
			}
			if setting != nil && setting.RuleActiveThreshold != tt.expectedRAT {
				t.Errorf("Expected RuleActiveThreshold=%v, got %v", tt.expectedRAT, setting.RuleActiveThreshold)
			}
		})
	}
}

func TestSetIncidentPredictionSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		setting      *IncidentPredictionSetting
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name:     "successful set",
			systemID: "sys-abc123",
			setting: &IncidentPredictionSetting{
				RuleActiveThreshold:   0.8,
				RuleInactiveThreshold: 0.2,
			},
			statusCode:   200,
			responseBody: `{"success": true}`,
			expectError:  false,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			setting:      &IncidentPredictionSetting{},
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
		{
			name:         "api returns success false",
			systemID:     "sys-abc123",
			setting:      &IncidentPredictionSetting{},
			statusCode:   200,
			responseBody: `{"success": false, "message": "invalid system"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("Failed to parse form: %v", err)
				}
				if r.FormValue("settingModels") == "" {
					t.Error("Expected settingModels in form data")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			err := c.SetIncidentPredictionSetting(tt.systemID, tt.setting)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// ---- GetHealthViewSetting ----

func TestGetHealthViewSetting(t *testing.T) {
	tests := []struct {
		name          string
		systemID      string
		statusCode    int
		responseBody  string
		expectError   bool
		expectNil     bool
		expectedEmail string
	}{
		{
			name:       "successful get - system found in map",
			systemID:   "sys-abc123",
			statusCode: 200,
			responseBody: `{
				"sys-abc123": {
					"key": {
						"systemPartitionKey": {
							"userName": "test_user",
							"systemName": "sys-abc123",
							"envName": "All"
						},
						"replay": false
					},
					"order": 1,
					"predictionEmail": "alert@example.com",
					"alertHealthScore": 0.5
				}
			}`,
			expectError:   false,
			expectNil:     false,
			expectedEmail: "alert@example.com",
		},
		{
			name:       "system not in map",
			systemID:   "sys-notfound",
			statusCode: 200,
			responseBody: `{
				"sys-other": {
					"key": {"systemPartitionKey": {"userName": "u","systemName": "sys-other","envName": "All"},"replay": false}
				}
			}`,
			expectError: false,
			expectNil:   true,
		},
		{
			name:         "not found 404",
			systemID:     "sys-abc123",
			statusCode:   404,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			setting, err := c.GetHealthViewSetting(tt.systemID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && setting != nil {
				t.Error("Expected nil but got non-nil setting")
			}
			if !tt.expectNil && setting == nil && !tt.expectError {
				t.Error("Expected non-nil setting but got nil")
			}
			if setting != nil {
				if setting.PredictionEmail != tt.expectedEmail {
					t.Errorf("Expected PredictionEmail=%s, got %s", tt.expectedEmail, setting.PredictionEmail)
				}
				if setting.SystemID != tt.systemID {
					t.Errorf("Expected SystemID=%s, got %s", tt.systemID, setting.SystemID)
				}
			}
		})
	}
}

func TestSetHealthViewSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		updates      *HealthViewSetting
		getResponse  string
		postResponse string
		statusCode   int
		expectError  bool
	}{
		{
			name:     "successful update - system exists",
			systemID: "sys-abc123",
			updates: &HealthViewSetting{
				PredictionEmail:  "new@example.com",
				AlertHealthScore: 0.7,
				AssignmentMap:    map[string]any{},
			},
			getResponse: `{
				"sys-abc123": {
					"key": {"systemPartitionKey": {"userName": "test_user","systemName": "sys-abc123","envName": "All"},"replay": false},
					"predictionEmail": "old@example.com"
				}
			}`,
			postResponse: `{"success": true}`,
			statusCode:   200,
			expectError:  false,
		},
		{
			name:     "successful update - system not in map (new entry)",
			systemID: "sys-new",
			updates: &HealthViewSetting{
				PredictionEmail: "new@example.com",
				AssignmentMap:   map[string]any{},
			},
			getResponse:  `{}`,
			postResponse: `{"success": true}`,
			statusCode:   200,
			expectError:  false,
		},
		{
			name:     "post returns success false",
			systemID: "sys-abc123",
			updates: &HealthViewSetting{
				AssignmentMap: map[string]any{},
			},
			getResponse:  `{"sys-abc123": {"key": {"systemPartitionKey": {"userName": "u","systemName": "sys-abc123","envName": "All"},"replay": false}}}`,
			postResponse: `{"success": false, "message": "update failed"}`,
			statusCode:   200,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				if callCount == 1 {
					// GET call
					if r.Method != "GET" {
						t.Errorf("Expected GET on first call, got %s", r.Method)
					}
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write([]byte(tt.getResponse))
				} else {
					// POST call
					if r.Method != "POST" {
						t.Errorf("Expected POST on second call, got %s", r.Method)
					}
					if err := r.ParseForm(); err != nil {
						t.Errorf("Failed to parse form: %v", err)
					}
					if r.FormValue("settings") == "" {
						t.Error("Expected settings in form data")
					}
					// Verify settings is valid JSON array
					var settingsArr []HealthViewSetting
					if err := json.Unmarshal([]byte(r.FormValue("settings")), &settingsArr); err != nil {
						t.Errorf("settings is not valid JSON array: %v", err)
					}
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write([]byte(tt.postResponse))
				}
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			err := c.SetHealthViewSetting(tt.systemID, tt.updates)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// ---- SatelliteSystemSetToList ----

func TestSatelliteSystemSetToList(t *testing.T) {
	input := []SatelliteSystemSetEntry{
		{
			SystemPartitionKey: SatelliteSystemPartitionKey{
				UserName:   "user1",
				SystemName: "sys-hash-1",
				EnvName:    "All",
			},
			Replay: false,
		},
		{
			SystemPartitionKey: SatelliteSystemPartitionKey{
				UserName:   "user2",
				SystemName: "sys-hash-2",
				EnvName:    "All",
			},
			Replay: true,
		},
	}

	result := SatelliteSystemSetToList(input)

	if len(result) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(result))
	}
	if result[0].SystemID != "sys-hash-1" {
		t.Errorf("Expected SystemID=sys-hash-1, got %s", result[0].SystemID)
	}
	if result[0].UserName != "user1" {
		t.Errorf("Expected UserName=user1, got %s", result[0].UserName)
	}
	if result[1].SystemID != "sys-hash-2" {
		t.Errorf("Expected SystemID=sys-hash-2, got %s", result[1].SystemID)
	}
}

func TestSatelliteSystemSetToList_Empty(t *testing.T) {
	result := SatelliteSystemSetToList([]SatelliteSystemSetEntry{})
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d entries", len(result))
	}
}

// ---- GetSystemDownSetting ----

func TestGetSystemDownSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		statusCode   int
		responseBody string
		expectError  bool
		expectNil    bool
		expectedEnable bool
	}{
		{
			name:       "successful get",
			systemID:   "sys-abc123",
			statusCode: 200,
			responseBody: `[{
				"enableSystemDownEmailAlert": true,
				"emailDampeningPeriod": 3600000,
				"emailSet": ["ops@example.com"]
			}]`,
			expectError:    false,
			expectNil:      false,
			expectedEnable: true,
		},
		{
			name:         "not found 404",
			systemID:     "sys-notfound",
			statusCode:   404,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "no content 204",
			systemID:     "sys-empty",
			statusCode:   204,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "empty array",
			systemID:     "sys-abc123",
			statusCode:   200,
			responseBody: `[]`,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				params := r.URL.Query()
				if params.Get("customerName") == "" {
					t.Error("Expected customerName query param")
				}
				if params.Get("systemIds") == "" {
					t.Error("Expected systemIds query param")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			setting, err := c.GetSystemDownSetting(tt.systemID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && setting != nil {
				t.Error("Expected nil but got non-nil setting")
			}
			if !tt.expectNil && setting == nil && !tt.expectError {
				t.Error("Expected non-nil setting but got nil")
			}
			if setting != nil && setting.EnableSystemDownEmailAlert != tt.expectedEnable {
				t.Errorf("Expected EnableSystemDownEmailAlert=%v, got %v", tt.expectedEnable, setting.EnableSystemDownEmailAlert)
			}
		})
	}
}

func TestSetSystemDownSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		setting      *SystemDownSetting
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name:     "successful set",
			systemID: "sys-abc123",
			setting: &SystemDownSetting{
				EnableSystemDownEmailAlert: true,
				EmailDampeningPeriod:       3600000,
				EmailSet:                   []string{"ops@example.com"},
			},
			statusCode:   200,
			responseBody: `{"success": true}`,
			expectError:  false,
		},
		{
			name:         "api returns success false",
			systemID:     "sys-abc123",
			setting:      &SystemDownSetting{},
			statusCode:   200,
			responseBody: `{"success": false, "message": "not found"}`,
			expectError:  true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			setting:      &SystemDownSetting{},
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("Failed to parse form: %v", err)
				}
				if r.FormValue("customerName") == "" {
					t.Error("Expected customerName in form data")
				}
				if r.FormValue("settingModels") == "" {
					t.Error("Expected settingModels in form data")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			err := c.SetSystemDownSetting(tt.systemID, tt.setting)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// ---- GetInsightsReportSetting ----

func TestGetInsightsReportSetting(t *testing.T) {
	tests := []struct {
		name           string
		systemID       string
		statusCode     int
		responseBody   string
		expectError    bool
		expectNil      bool
		expectedEnable bool
	}{
		{
			name:       "successful get",
			systemID:   "sys-abc123",
			statusCode: 200,
			responseBody: `[{
				"enableDailyInsightsReport": true,
				"emailSet": ["reports@example.com"],
				"weeklyEmailSet": ["exec@example.com"],
				"enableWeeklyInsightsReport": false
			}]`,
			expectError:    false,
			expectNil:      false,
			expectedEnable: true,
		},
		{
			name:         "not found 404",
			systemID:     "sys-notfound",
			statusCode:   404,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "empty array",
			systemID:     "sys-abc123",
			statusCode:   200,
			responseBody: `[]`,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				params := r.URL.Query()
				if params.Get("customerName") == "" {
					t.Error("Expected customerName query param")
				}
				if params.Get("systemIds") == "" {
					t.Error("Expected systemIds query param")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			setting, err := c.GetInsightsReportSetting(tt.systemID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && setting != nil {
				t.Error("Expected nil but got non-nil setting")
			}
			if !tt.expectNil && setting == nil && !tt.expectError {
				t.Error("Expected non-nil setting but got nil")
			}
			if setting != nil && setting.EnableDailyInsightsReport != tt.expectedEnable {
				t.Errorf("Expected EnableDailyInsightsReport=%v, got %v", tt.expectedEnable, setting.EnableDailyInsightsReport)
			}
		})
	}
}

func TestSetInsightsReportSetting(t *testing.T) {
	tests := []struct {
		name         string
		systemID     string
		emailSet     []string
		enable       bool
		isDaily      bool
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name:         "successful set daily",
			systemID:     "sys-abc123",
			emailSet:     []string{"reports@example.com"},
			enable:       true,
			isDaily:      true,
			statusCode:   200,
			responseBody: `{"success": true}`,
			expectError:  false,
		},
		{
			name:         "successful set weekly",
			systemID:     "sys-abc123",
			emailSet:     []string{"exec@example.com"},
			enable:       false,
			isDaily:      false,
			statusCode:   200,
			responseBody: `{"success": true}`,
			expectError:  false,
		},
		{
			name:         "api returns success false",
			systemID:     "sys-abc123",
			emailSet:     []string{},
			enable:       false,
			isDaily:      true,
			statusCode:   200,
			responseBody: `{"success": false, "message": "invalid system"}`,
			expectError:  true,
		},
		{
			name:         "server error",
			systemID:     "sys-abc123",
			emailSet:     []string{},
			enable:       false,
			isDaily:      true,
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("Failed to parse form: %v", err)
				}
				if r.FormValue("customerName") == "" {
					t.Error("Expected customerName in form data")
				}
				if r.FormValue("settingModels") == "" {
					t.Error("Expected settingModels in form data")
				}
				if r.FormValue("isDaily") == "" {
					t.Error("Expected isDaily in form data")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			err := c.SetInsightsReportSetting(tt.systemID, tt.emailSet, tt.enable, tt.isDaily)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// ---- GetInstanceDownSetting ----

func TestGetInstanceDownSetting(t *testing.T) {
	tests := []struct {
		name            string
		projectName     string
		statusCode      int
		responseBody    string
		expectError     bool
		expectNil       bool
		expectedEnable  bool
	}{
		{
			name:        "successful get",
			projectName: "my-project",
			statusCode:  200,
			responseBody: `{
				"data": {
					"projectModelAllJSON": {
						"projectName": "my-project",
						"instanceDownEnable": true,
						"instanceDownDampening": 1800000,
						"instanceDownThreshold": 300000,
						"instanceDownReportNumber": 2,
						"instanceDownEmails": ["oncall@example.com"]
					}
				},
				"success": true
			}`,
			expectError:    false,
			expectNil:      false,
			expectedEnable: true,
		},
		{
			name:        "success false",
			projectName: "bad-project",
			statusCode:  200,
			responseBody: `{
				"data": {},
				"success": false
			}`,
			expectError: true,
			expectNil:   true,
		},
		{
			name:         "not found 404",
			projectName:  "notfound",
			statusCode:   404,
			responseBody: ``,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "server error",
			projectName:  "my-project",
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				params := r.URL.Query()
				if params.Get("projectName") == "" {
					t.Error("Expected projectName query param")
				}
				if params.Get("customerName") == "" {
					t.Error("Expected customerName query param")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			setting, err := c.GetInstanceDownSetting(tt.projectName)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectNil && setting != nil {
				t.Error("Expected nil but got non-nil setting")
			}
			if !tt.expectNil && setting == nil && !tt.expectError {
				t.Error("Expected non-nil setting but got nil")
			}
			if setting != nil && setting.InstanceDownEnable != tt.expectedEnable {
				t.Errorf("Expected InstanceDownEnable=%v, got %v", tt.expectedEnable, setting.InstanceDownEnable)
			}
		})
	}
}

func TestSetInstanceDownSetting(t *testing.T) {
	tests := []struct {
		name         string
		setting      *InstanceDownSetting
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name: "successful set",
			setting: &InstanceDownSetting{
				ProjectName:              "my-project",
				InstanceDownEnable:       true,
				InstanceDownDampening:    1800000,
				InstanceDownThreshold:    300000,
				InstanceDownReportNumber: 2,
				InstanceDownEmails:       []string{"oncall@example.com"},
			},
			statusCode:   200,
			responseBody: `{"alertSetting": {"success": true}, "instanceDownThreshold": {"success": true}}`,
			expectError:  false,
		},
		{
			name: "sub-key success false",
			setting: &InstanceDownSetting{
				ProjectName:        "my-project",
				InstanceDownEmails: []string{},
			},
			statusCode:   200,
			responseBody: `{"alertSetting": {"success": false}}`,
			expectError:  true,
		},
		{
			name: "server error",
			setting: &InstanceDownSetting{
				ProjectName:        "my-project",
				InstanceDownEmails: []string{},
			},
			statusCode:   500,
			responseBody: `{"error":"internal error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("Failed to parse form: %v", err)
				}
				if r.FormValue("projectName") == "" {
					t.Error("Expected projectName in form data")
				}
				if r.FormValue("operation") == "" {
					t.Error("Expected operation in form data")
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c, _ := NewClient(server.URL, "test_user", "test_key")
			err := c.SetInstanceDownSetting(tt.setting)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
