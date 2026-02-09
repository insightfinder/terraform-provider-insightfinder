// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Holiday represents a holiday configuration
type Holiday struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// HolidayResponse represents the API response for holiday operations
type HolidayResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// GetHolidays retrieves all holidays for a project
func (c *Client) GetHolidays(projectName string) (map[string]string, error) {
	params := url.Values{}
	params.Add("projectName", projectName)
	params.Add("customerName", c.Username)

	path := fmt.Sprintf("/api/external/v1/holiday?%s", params.Encode())
	body, statusCode, err := c.DoRequestWithLicenseKeyHeader("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return make(map[string]string), nil // No holidays exist
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get holidays: HTTP %d - %s", statusCode, string(body))
	}

	// Response format: {"holidayName": "startDate,endDate", ...}
	var holidays map[string]string
	if err := json.Unmarshal(body, &holidays); err != nil {
		return nil, fmt.Errorf("failed to parse holiday response: %w", err)
	}

	return holidays, nil
}

// CreateHoliday creates a new holiday for a project
func (c *Client) CreateHoliday(projectName string, holiday *Holiday) error {
	params := url.Values{}
	params.Add("projectName", projectName)
	params.Add("holidayName", holiday.Name)
	params.Add("startTime", holiday.StartDate)
	params.Add("endTime", holiday.EndDate)

	path := fmt.Sprintf("/api/external/v1/holiday?%s", params.Encode())
	body, statusCode, err := c.DoRequestWithLicenseKeyHeader("POST", path, nil)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to create holiday: HTTP %d - %s", statusCode, string(body))
	}

	var response HolidayResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("failed to create holiday: %s", response.Message)
	}

	return nil
}

// DeleteHolidays deletes holidays from a project
func (c *Client) DeleteHolidays(projectName string, holidayNames []string) error {
	if len(holidayNames) == 0 {
		return nil // Nothing to delete
	}

	// Convert holiday names to JSON array
	holidayListJSON, err := json.Marshal(holidayNames)
	if err != nil {
		return fmt.Errorf("failed to marshal holiday names: %w", err)
	}

	params := url.Values{}
	params.Add("projectName", projectName)
	params.Add("holidayNameList", string(holidayListJSON))

	path := fmt.Sprintf("/api/external/v1/holiday?%s", params.Encode())
	body, statusCode, err := c.DoRequestWithLicenseKeyHeader("DELETE", path, nil)
	if err != nil {
		return err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil // Already deleted or doesn't exist
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to delete holidays: HTTP %d - %s", statusCode, string(body))
	}

	var response HolidayResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("failed to delete holidays: %s", response.Message)
	}

	return nil
}
