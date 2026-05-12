// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// L2MDerivedValueModel represents the derived value model within a JSON parser
type L2MDerivedValueModel struct {
	BaseValue     *string  `json:"baseValue,omitempty"`
	ActualValue   *string  `json:"actualValue,omitempty"`
	Operation     *int     `json:"operation,omitempty"`
	MappingIDList []string `json:"mappingIdList,omitempty"`
}

// L2MRegexParser represents a single regex parser in L2M settings
type L2MRegexParser struct {
	MetricNameRegex     *string `json:"metricNameRegex,omitempty"`
	MetricValueRegex    *string `json:"metricValueRegex,omitempty"`
	BaseValueKey        *string `json:"baseValueKey,omitempty"`
	InstanceNameRegex   *string `json:"instanceNameRegex,omitempty"`
	ContainerNameRegex  *string `json:"containerNameRegex,omitempty"`
	TimestampRegex      *string `json:"timestampRegex,omitempty"`
	TimestampFormat     *string `json:"timestampFormat,omitempty"`
	DataFilter          *string `json:"dataFilter,omitempty"`
	Operation           *int    `json:"operation,omitempty"`
	AggregationMode     *int    `json:"aggregationMode,omitempty"`
	MetricName          *string `json:"metricName,omitempty"`
	GroupingByComponent *bool   `json:"groupingByComponent,omitempty"`
	AggregationPeriod   *int    `json:"aggregationPeriod,omitempty"`
	ContainerType       *int    `json:"containerType,omitempty"`
}

// L2MJsonParser represents a single JSON parser in L2M settings
type L2MJsonParser struct {
	MetricValueKey       *string               `json:"metricValueKey,omitempty"`
	BaseValueKey         *string               `json:"baseValueKey,omitempty"`
	InstanceNameKey      *string               `json:"instanceNameKey,omitempty"`
	ContainerNameKey     *string               `json:"containerNameKey,omitempty"`
	TimestampKey         *string               `json:"timestampKey,omitempty"`
	TimestampFormat      *string               `json:"timestampFormat,omitempty"`
	Operation            *int                  `json:"operation,omitempty"`
	AdditionalMetricName *string               `json:"additionalMetricName,omitempty"`
	AggregationMode      *int                  `json:"aggregationMode,omitempty"`
	GroupingByComponent  *bool                 `json:"groupingByComponent,omitempty"`
	AggregationPeriod    *int                  `json:"aggregationPeriod,omitempty"`
	ContainerType        *int                  `json:"containerType,omitempty"`
	DerivedValueModel    *L2MDerivedValueModel `json:"derivedValueModel,omitempty"`
}

// L2MSetting represents a single L2M configuration entry (one per target metric project)
type L2MSetting struct {
	MetricProjectName string           `json:"metricProjectName"`
	JsonFlag          *bool            `json:"jsonFlag,omitempty"`
	EnableMapping     *bool            `json:"enableMapping,omitempty"`
	Regexs            []L2MRegexParser `json:"regexs,omitempty"`
	JsonParsers       []L2MJsonParser  `json:"jsonParsers,omitempty"`
}

// GetL2MSettings retrieves all L2M settings for a project
func (c *Client) GetL2MSettings(projectName string) ([]L2MSetting, error) {
	path := fmt.Sprintf("/api/external/v1/logtometricsetting?projectName=%s", url.QueryEscape(projectName))

	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return []L2MSetting{}, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get L2M settings: HTTP %d - %s", statusCode, string(body))
	}

	var settings []L2MSetting
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse L2M settings: %w", err)
	}

	return settings, nil
}

// UpdateL2MSetting creates or replaces the L2M configuration for one target metric project
func (c *Client) UpdateL2MSetting(projectName string, setting *L2MSetting) error {
	regexInfoJSON, err := json.Marshal(setting)
	if err != nil {
		return fmt.Errorf("failed to marshal L2M setting: %w", err)
	}

	formData := url.Values{}
	formData.Set("projectName", projectName)
	formData.Set("regexInfo", string(regexInfoJSON))

	path := fmt.Sprintf("/api/external/v1/logtometricsetting?projectName=%s", url.QueryEscape(projectName))
	body, statusCode, err := c.DoFormRequest("POST", path, formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to update L2M setting: HTTP %d - %s", statusCode, string(body))
	}

	return nil
}
