// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// JiraFieldUpdateRule represents a single post-creation Jira field-update rule: when a ticket's
// root-cause project/metric/component match, fieldId is PATCHed to value after ticket creation
// (reaching fields that are not on the issue type's Create screen).
type JiraFieldUpdateRule struct {
	Projects   []string `json:"projects,omitempty"`
	Metrics    []string `json:"metrics,omitempty"`
	Components []string `json:"components,omitempty"`
	FieldID    string   `json:"fieldId,omitempty"`
	FieldName  string   `json:"fieldName,omitempty"`
	ValueType  string   `json:"valueType,omitempty"`
	Value      string   `json:"value,omitempty"`
	Enabled    bool     `json:"enabled"`
}

// JiraConfig represents Jira integration configuration.
type JiraConfig struct {
	Account              string
	HostURL              string
	APIToken             string
	SystemIDs            []string
	Options              []string
	ContentOption        []string
	JiraProjectKey       string
	JiraReporterID       string
	JiraAssigneeID       string
	JiraIssueFields      map[string]string
	JiraCloseStatusID    string
	DescriptionTemplate  string
	TemplateFieldMapping map[string]string
	FieldUpdateRules     []JiraFieldUpdateRule
}

// JiraResponse represents the API response for Jira operations.
type JiraResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// jiraConfigsPayload mirrors JiraIntegrationConfig (InfrastructureCore) field-for-field. Unlike
// Slack/ServiceNow, the list endpoint's per-entry serialization (Customer.getAllExtServices) has
// no Jira-specific flattening block, so every field below only ever arrives inside the entry's
// generic "configs" JSON blob — never as flattened top-level keys.
type jiraConfigsPayload struct {
	SystemIDs                []string              `json:"systemIds"`
	Account                  string                `json:"account"`
	HostURL                  string                `json:"hostUrl"`
	EncodedToken             string                `json:"encodedToken"`
	JiraProjectKey           string                `json:"jiraProjectKey"`
	JiraReporterID           string                `json:"jiraReporterId"`
	ResolveTransitionPathIDs []string              `json:"resolveTransitionPathIds"`
	ContentOption            []string              `json:"contentOption"`
	JiraAssigneeID           string                `json:"jiraAssigneeId"`
	JiraIssueFields          map[string]string     `json:"jiraIssueFields"`
	JiraCloseStatusID        string                `json:"jiraCloseStatusId"`
	DescriptionTemplate      string                `json:"descriptionTemplate"`
	TemplateFieldMapping     map[string]string     `json:"templateFieldMapping"`
	FieldUpdateRules         []JiraFieldUpdateRule `json:"fieldUpdateRules"`
}

// stringSliceJSONOrDefault marshals v, falling back to empty (a valid JSON literal like "[]")
// when v is nil/empty. Several Jira fields are read on the backend via helpers that choke on a
// missing/null JSON value but treat "[]"/"{}" as a normal empty value — Go's json.Marshal of a
// nil slice/map produces the string "null", which would trip that, so callers must never send
// json.Marshal's raw output for a possibly-nil collection.
func stringSliceJSONOrDefault(v []string, empty string) string {
	if len(v) == 0 {
		return empty
	}
	b, err := json.Marshal(v)
	if err != nil {
		return empty
	}
	return string(b)
}

func stringMapJSONOrDefault(v map[string]string, empty string) string {
	if len(v) == 0 {
		return empty
	}
	b, err := json.Marshal(v)
	if err != nil {
		return empty
	}
	return string(b)
}

func fieldUpdateRulesJSONOrDefault(v []JiraFieldUpdateRule, empty string) string {
	if len(v) == 0 {
		return empty
	}
	b, err := json.Marshal(v)
	if err != nil {
		return empty
	}
	return string(b)
}

// ensureTrailingSlash appends "/" unless it's already there.
func ensureTrailingSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

// jiraAPIBaseURL mirrors JiraHelper.getJiraBaseURL() on the backend, which appends "rest/api/"
// to the configured host and stores THAT as the integration's hostUrl. Anything addressing an
// existing row by key — notably the delete service_id — has to use this form, not the site root.
func jiraAPIBaseURL(host string) string {
	h := ensureTrailingSlash(strings.TrimSpace(host))
	if strings.HasSuffix(h, "rest/api/") {
		return h
	}
	return h + "rest/api/"
}

// GetJiraConfig retrieves Jira integration configuration.
func (c *Client) GetJiraConfig(account, hostURL, username string) (*JiraConfig, error) {
	params := url.Values{}
	params.Add("serviceProvider", "AtlassianJira")
	params.Add("tzOffset", "-14400000")

	path := fmt.Sprintf("/api/external/v1/system/externalServlies/list?%s", params.Encode())
	body, statusCode, err := c.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 || statusCode == 204 {
		return nil, nil // Configuration doesn't exist
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get Jira config: HTTP %d", statusCode)
	}

	var response struct {
		ExtServiceAllInfo []map[string]interface{} `json:"extServiceAllInfo"`
		Success           bool                     `json:"success"`
		Message           string                   `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return nil, nil
	}

	// Find the entry matching account and service_host.
	// The backend stores hostUrl after JiraHelper.getJiraBaseURL() has appended "rest/api/",
	// so a stored entry reads "https://site/rest/api/" while the resource is configured with
	// the site root. Strip that suffix on both sides or the lookup never matches, and Read
	// would report the integration as deleted on every refresh.
	normalizeHost := func(h string) string {
		h = strings.TrimRight(strings.TrimSpace(h), "/")
		h = strings.TrimSuffix(h, "/rest/api")
		return strings.TrimRight(h, "/")
	}
	var entry map[string]interface{}
	for _, info := range response.ExtServiceAllInfo {
		entryAccount, _ := info["account"].(string)
		entryHost, _ := info["service_host"].(string)
		if strings.EqualFold(strings.TrimSpace(entryAccount), strings.TrimSpace(account)) &&
			normalizeHost(entryHost) == normalizeHost(hostURL) {
			entry = info
			break
		}
	}

	if entry == nil {
		return nil, nil // Not found
	}

	config := &JiraConfig{
		Account: account,
		HostURL: hostURL,
	}

	// "options" (notification trigger types) is a generic column populated the same way for
	// every provider, separate from the Jira-specific fields below.
	if optionsStr, ok := entry["options"].(string); ok && optionsStr != "" {
		var options []string
		if err := json.Unmarshal([]byte(optionsStr), &options); err == nil {
			config.Options = options
		}
	}

	configsStr, ok := entry["configs"].(string)
	if !ok || configsStr == "" {
		return config, nil
	}

	var payload jiraConfigsPayload
	if err := json.Unmarshal([]byte(configsStr), &payload); err != nil {
		return nil, fmt.Errorf("failed to parse Jira configs: %w", err)
	}

	config.SystemIDs = payload.SystemIDs
	config.JiraProjectKey = payload.JiraProjectKey
	config.JiraReporterID = payload.JiraReporterID
	config.ContentOption = payload.ContentOption
	config.JiraAssigneeID = payload.JiraAssigneeID
	config.JiraIssueFields = payload.JiraIssueFields
	config.JiraCloseStatusID = payload.JiraCloseStatusID
	config.DescriptionTemplate = payload.DescriptionTemplate
	config.TemplateFieldMapping = payload.TemplateFieldMapping
	config.FieldUpdateRules = payload.FieldUpdateRules

	return config, nil
}

// CreateOrUpdateJiraConfig creates or updates a Jira integration.
func (c *Client) CreateOrUpdateJiraConfig(config *JiraConfig, username string, verify bool) error {
	formData := url.Values{}
	if verify {
		formData.Set("verify", "true")
	}
	formData.Set("operation", "AtlassianJira")
	formData.Set("customerName", username)
	formData.Set("account", config.Account)
	// Wire field is "password" (shared with ServiceNow's real account password) but for Jira it
	// carries the API token — the backend base64-encodes account+token into its stored
	// encodedToken itself, so the provider never needs to do that encoding.
	formData.Set("password", config.APIToken)
	// JiraHelper.getJiraBaseURL (backend) appends "rest/api/..." to this value via plain string
	// concatenation, with no separator — a host without a trailing slash produces a broken
	// hostname (e.g. "company.atlassian.netrest") and every Jira API call then fails as if the
	// credentials were wrong. Always send a trailing slash so users don't have to know this.
	formData.Set("service_host", ensureTrailingSlash(config.HostURL))
	// account/service_host are RequiresReplace in the schema, so on every Update call this
	// resource actually makes, they're identical to what's already stored — sending them back
	// unconditionally just keeps parity with the backend's key-change-detection contract without
	// this provider ever needing to exercise its "key changed" branch (Terraform already handles
	// a real key change via destroy+recreate).
	//
	// storedHost must be the *derived* base URL, not the raw host: the backend compares it
	// against JiraHelper.getJiraBaseURL(service_host) both to look up the existing row
	// (existingLookupHost) and to decide keyChanged. Sending the site root instead makes that
	// lookup miss on every call — so the backend sees keyChanged on a plain in-place edit and
	// fires a delete against a key that does not exist, and its "absent field preserves the
	// stored value" branch silently loses the value it was meant to preserve.
	formData.Set("stored_account", config.Account)
	formData.Set("storedHost", jiraAPIBaseURL(config.HostURL))
	formData.Set("options", stringSliceJSONOrDefault(config.Options, "[]"))

	if !verify {
		formData.Set("systemIds", stringSliceJSONOrDefault(config.SystemIDs, "[]"))
		formData.Set("jiraProjectKey", config.JiraProjectKey)
		formData.Set("jiraReporterId", config.JiraReporterID)
		// These six fields all share the backend's "absent preserves / blank clears / present
		// overwrites" semantics (resolveOptionalField / resolveContentOptions in
		// ExtServiceIntegrationServlet.addJiraCredentials). We always send them explicitly —
		// using the JSON-empty form ("", "{}", "[]") when unset — so a value removed from
		// Terraform config actually lands on the "clear" branch instead of silently falling
		// through to "absent means preserve whatever the backend already has".
		formData.Set("jiraAssigneeId", config.JiraAssigneeID)
		formData.Set("jiraCloseStatusId", config.JiraCloseStatusID)
		formData.Set("descriptionTemplate", config.DescriptionTemplate)
		formData.Set("jiraIssueFields", stringMapJSONOrDefault(config.JiraIssueFields, "{}"))
		formData.Set("templateFieldMapping", stringMapJSONOrDefault(config.TemplateFieldMapping, "{}"))
		formData.Set("fieldUpdateRules", fieldUpdateRulesJSONOrDefault(config.FieldUpdateRules, "[]"))
		formData.Set("contentOption", stringSliceJSONOrDefault(config.ContentOption, "[]"))
	}

	path := "/api/external/v1/service-integration"
	body, statusCode, err := c.DoFormRequest("POST", path, formData)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to configure Jira: HTTP %d - %s", statusCode, string(body))
	}

	var response JiraResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// If we can't parse the response but got 200, assume success
		return nil
	}

	if !response.Success {
		if response.Message != "" {
			return fmt.Errorf("Jira configuration failed: %s", response.Message)
		}
		return fmt.Errorf("Jira configuration failed")
	}

	return nil
}

// DeleteJiraConfig removes a Jira integration.
func (c *Client) DeleteJiraConfig(account, hostURL, username string) error {
	hostURL = strings.TrimSpace(hostURL)
	if hostURL == "" {
		return fmt.Errorf("service_host is required for deletion")
	}

	serviceID := fmt.Sprintf("AtlassianJira:%s:%s", account, jiraAPIBaseURL(hostURL))

	formData := url.Values{}
	formData.Set("serviceProvider", "AtlassianJira")
	formData.Set("operation", "delete")
	formData.Set("service_id", serviceID)
	formData.Set("serviceOwner", username)
	formData.Set("customerName", username)

	path := "/api/external/v1/service-integration"
	body, statusCode, err := c.DoFormRequest("POST", path, formData)
	if err != nil {
		return err
	}

	// 200 or 404 are both acceptable for deletion
	if statusCode != 200 && statusCode != 404 {
		return fmt.Errorf("failed to delete Jira config: HTTP %d - %s", statusCode, string(body))
	}

	return nil
}
