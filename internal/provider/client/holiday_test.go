// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHolidays(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/api/external/v1/holiday" {
			t.Errorf("Expected path /api/external/v1/holiday, got %s", r.URL.Path)
		}

		// Verify query parameters
		projectName := r.URL.Query().Get("projectName")
		if projectName != "test-project" {
			t.Errorf("Expected projectName=test-project, got %s", projectName)
		}

		// Return mock response
		response := map[string]string{
			"christmas": "12-25,12-26",
			"new_year":  "01-01,01-01",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := &Client{
		BaseURL:    server.URL,
		Username:   "test-user",
		LicenseKey: "test-key",
		HTTPClient: &http.Client{},
	}

	// Test GetHolidays
	holidays, err := client.GetHolidays("test-project")
	if err != nil {
		t.Fatalf("GetHolidays failed: %v", err)
	}

	if len(holidays) != 2 {
		t.Errorf("Expected 2 holidays, got %d", len(holidays))
	}

	if holidays["christmas"] != "12-25,12-26" {
		t.Errorf("Expected christmas to be '12-25,12-26', got '%s'", holidays["christmas"])
	}

	if holidays["new_year"] != "01-01,01-01" {
		t.Errorf("Expected new_year to be '01-01,01-01', got '%s'", holidays["new_year"])
	}
}

func TestCreateHoliday(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/api/external/v1/holiday" {
			t.Errorf("Expected path /api/external/v1/holiday, got %s", r.URL.Path)
		}

		// Verify query parameters
		projectName := r.URL.Query().Get("projectName")
		if projectName != "test-project" {
			t.Errorf("Expected projectName=test-project, got %s", projectName)
		}

		holidayName := r.URL.Query().Get("holidayName")
		if holidayName != "test_holiday" {
			t.Errorf("Expected holidayName=test_holiday, got %s", holidayName)
		}

		startTime := r.URL.Query().Get("startTime")
		if startTime != "12-25" {
			t.Errorf("Expected startTime=12-25, got %s", startTime)
		}

		endTime := r.URL.Query().Get("endTime")
		if endTime != "12-26" {
			t.Errorf("Expected endTime=12-26, got %s", endTime)
		}

		// Return success response
		response := HolidayResponse{
			Success: true,
			Message: "Success",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := &Client{
		BaseURL:    server.URL,
		Username:   "test-user",
		LicenseKey: "test-key",
		HTTPClient: &http.Client{},
	}

	// Test CreateHoliday
	holiday := &Holiday{
		Name:      "test_holiday",
		StartDate: "12-25",
		EndDate:   "12-26",
	}

	err := client.CreateHoliday("test-project", holiday)
	if err != nil {
		t.Fatalf("CreateHoliday failed: %v", err)
	}
}

func TestDeleteHolidays(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/api/external/v1/holiday" {
			t.Errorf("Expected path /api/external/v1/holiday, got %s", r.URL.Path)
		}

		// Verify query parameters
		projectName := r.URL.Query().Get("projectName")
		if projectName != "test-project" {
			t.Errorf("Expected projectName=test-project, got %s", projectName)
		}

		holidayNameList := r.URL.Query().Get("holidayNameList")
		var names []string
		if err := json.Unmarshal([]byte(holidayNameList), &names); err != nil {
			t.Errorf("Failed to parse holidayNameList: %v", err)
		}

		if len(names) != 2 {
			t.Errorf("Expected 2 holiday names, got %d", len(names))
		}

		// Return success response
		response := HolidayResponse{
			Success: true,
			Message: "Success",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := &Client{
		BaseURL:    server.URL,
		Username:   "test-user",
		LicenseKey: "test-key",
		HTTPClient: &http.Client{},
	}

	// Test DeleteHolidays
	holidayNames := []string{"holiday1", "holiday2"}
	err := client.DeleteHolidays("test-project", holidayNames)
	if err != nil {
		t.Fatalf("DeleteHolidays failed: %v", err)
	}
}

func TestGetHolidays_Empty(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create client
	client := &Client{
		BaseURL:    server.URL,
		Username:   "test-user",
		LicenseKey: "test-key",
		HTTPClient: &http.Client{},
	}

	// Test GetHolidays with 404 response
	holidays, err := client.GetHolidays("test-project")
	if err != nil {
		t.Fatalf("GetHolidays failed: %v", err)
	}

	if len(holidays) != 0 {
		t.Errorf("Expected 0 holidays for 404 response, got %d", len(holidays))
	}
}
