// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceNowResource_BasicAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServiceNowResourceConfigBasicAuth(
					"test-account",
					"test.service-now.com",
					"testuser",
					"testpass",
					[]string{"system1", "system2"},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "account", "test-account"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "service_host", "test.service-now.com"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "password", "testpass"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "system_names.#", "2"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "system_names.0", "system1"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "system_names.1", "system2"),
					resource.TestCheckResourceAttrSet("insightfinder_servicenow.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "insightfinder_servicenow.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore password as it's not returned by API
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccServiceNowResourceConfigBasicAuth(
					"test-account",
					"test.service-now.com",
					"testuser",
					"newpassword",
					[]string{"system1", "system2", "system3"},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "system_names.#", "3"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "password", "newpassword"),
				),
			},
		},
	})
}

func TestAccServiceNowResource_OAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNowResourceConfigOAuth(
					"test-oauth-account",
					"test-oauth.service-now.com",
					"app-id-123",
					"app-key-secret",
					[]string{"oauth-system1"},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "account", "test-oauth-account"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "auth_type", "oauth"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "app_id", "app-id-123"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "app_key", "app-key-secret"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "system_names.#", "1"),
				),
			},
		},
	})
}

func TestAccServiceNowResource_WithProxy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNowResourceConfigWithProxy(
					"test-proxy-account",
					"test.service-now.com",
					"testuser",
					"testpass",
					"http://proxy.example.com:8080",
					[]string{"proxy-system"},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "account", "test-proxy-account"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "proxy", "http://proxy.example.com:8080"),
				),
			},
		},
	})
}

func TestAccServiceNowResource_WithProjectConfigs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNowResourceConfigWithProjectConfigs(
					"test-account",
					"test.service-now.com",
					"testpass",
					[]string{"system1"},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "account", "test-account"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "project_configs.my-project.enable_ticket_creation", "true"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "project_configs.my-project.enable_ticket_update", "true"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "project_configs.my-project.enable_incident_consolidation_info_update", "false"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "project_configs.another-project.enable_ticket_creation", "false"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "project_configs.another-project.enable_ticket_update", "false"),
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "project_configs.another-project.enable_incident_consolidation_info_update", "true"),
				),
			},
		},
	})
}

func TestAccServiceNowResource_WithDampeningPeriod(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNowResourceConfigWithDampening(
					"test-dampening-account",
					"test.service-now.com",
					"testuser",
					"testpass",
					30,
					[]string{"dampening-system"},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_servicenow.test", "dampening_period", "30"),
				),
			},
		},
	})
}

func testAccServiceNowResourceConfigBasicAuth(account, serviceHost, username, password string, systemNames []string) string {
	systemNamesStr := ""
	for _, name := range systemNames {
		systemNamesStr += fmt.Sprintf("%q,", name)
	}
	if len(systemNamesStr) > 0 {
		systemNamesStr = systemNamesStr[:len(systemNamesStr)-1] // Remove trailing comma
	}

	return fmt.Sprintf(`
resource "insightfinder_servicenow" "test" {
  account       = %[1]q
  service_host  = %[2]q
  password      = %[4]q
  system_names  = [%[5]s]

  options = [
    "send_incident",
    "sync_status"
  ]

  content_option = [
    "root_cause",
    "related_incidents"
  ]
}
`, account, serviceHost, username, password, systemNamesStr)
}

func testAccServiceNowResourceConfigOAuth(account, serviceHost, appID, appKey string, systemNames []string) string {
	systemNamesStr := ""
	for _, name := range systemNames {
		systemNamesStr += fmt.Sprintf("%q,", name)
	}
	if len(systemNamesStr) > 0 {
		systemNamesStr = systemNamesStr[:len(systemNamesStr)-1]
	}

	return fmt.Sprintf(`
resource "insightfinder_servicenow" "test" {
  account      = %[1]q
  service_host = %[2]q
  auth_type    = "oauth"
  app_id       = %[3]q
  app_key      = %[4]q
  system_names = [%[5]s]

  options = [
    "send_incident"
  ]

  content_option = [
    "root_cause"
  ]
}
`, account, serviceHost, appID, appKey, systemNamesStr)
}

func testAccServiceNowResourceConfigWithProxy(account, serviceHost, username, password, proxy string, systemNames []string) string {
	systemNamesStr := ""
	for _, name := range systemNames {
		systemNamesStr += fmt.Sprintf("%q,", name)
	}
	if len(systemNamesStr) > 0 {
		systemNamesStr = systemNamesStr[:len(systemNamesStr)-1]
	}

	return fmt.Sprintf(`
resource "insightfinder_servicenow" "test" {
  account      = %[1]q
  service_host = %[2]q
  password     = %[4]q
  proxy        = %[5]q
  system_names = [%[6]s]

  options = [
    "send_incident"
  ]

  content_option = [
    "root_cause"
  ]
}
`, account, serviceHost, username, password, proxy, systemNamesStr)
}

func testAccServiceNowResourceConfigWithDampening(account, serviceHost, username, password string, dampeningPeriod int, systemNames []string) string {
	systemNamesStr := ""
	for _, name := range systemNames {
		systemNamesStr += fmt.Sprintf("%q,", name)
	}
	if len(systemNamesStr) > 0 {
		systemNamesStr = systemNamesStr[:len(systemNamesStr)-1]
	}

	return fmt.Sprintf(`
resource "insightfinder_servicenow" "test" {
  account          = %[1]q
  service_host     = %[2]q
  password         = %[4]q
  dampening_period = %[5]d
  system_names     = [%[6]s]

  options = [
    "send_incident"
  ]

  content_option = [
    "root_cause"
  ]
}
`, account, serviceHost, username, password, dampeningPeriod, systemNamesStr)
}

func testAccServiceNowResourceConfigWithProjectConfigs(account, serviceHost, password string, systemNames []string) string {
	systemNamesStr := ""
	for _, name := range systemNames {
		systemNamesStr += fmt.Sprintf("%q,", name)
	}
	if len(systemNamesStr) > 0 {
		systemNamesStr = systemNamesStr[:len(systemNamesStr)-1]
	}

	return fmt.Sprintf(`
resource "insightfinder_servicenow" "test" {
  account          = %[1]q
  service_host     = %[2]q
  password         = %[3]q
  dampening_period = 7200000
  system_names     = [%[4]s]

  options = [
    "Root Cause",
    "Predicted Incident",
    "Detected Incident with RCA"
  ]

  content_option = ["SUMMARY", "RECOMMENDATION"]

  trigger_window_in_mills = 604800000
  enable_feedback_collect = true
  service_now_field       = "u_probable_cause"
  content_source          = "work_notes"

  ticket_created_by_source_key   = "activity_source"
  ticket_created_by_source_value = "abccccdddd"
  configuration_item             = "aaglskdfjasdl"

  project_configs = {
    "my-project" = {
      enable_ticket_creation                    = true
      enable_ticket_update                      = true
      enable_incident_consolidation_info_update = false
    }
    "another-project" = {
      enable_ticket_creation                    = false
      enable_ticket_update                      = false
      enable_incident_consolidation_info_update = true
    }
  }
}
`, account, serviceHost, password, systemNamesStr)
}
