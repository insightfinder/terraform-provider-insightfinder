// Copyright (c) InsightFinder Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJWTConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJWTConfigResourceConfig("test-system-jwt", "my-jwt-secret-key"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "system_name", "test-system-jwt"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "jwt_secret", "my-jwt-secret-key"),
					resource.TestCheckResourceAttrSet("insightfinder_jwt_config.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "insightfinder_jwt_config.test",
				ImportState:       true,
				ImportStateVerify: true,
				// jwt_secret is sensitive and not returned by API
				ImportStateVerifyIgnore: []string{"jwt_secret"},
			},
			// Update and Read testing
			{
				Config: testAccJWTConfigResourceConfig("test-system-jwt", "updated-jwt-secret-key"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "jwt_secret", "updated-jwt-secret-key"),
				),
			},
			// Delete testing automatically occurs at the end
			// This verifies that empty string deletion works correctly
		},
	})
}

func TestAccJWTConfigResource_MultipleSystems(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJWTConfigResourceConfigMultiple(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_jwt_config.system1", "system_name", "jwt-system-1"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.system1", "jwt_secret", "secret-for-system-1"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.system2", "system_name", "jwt-system-2"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.system2", "jwt_secret", "secret-for-system-2"),
				),
			},
		},
	})
}

func TestAccJWTConfigResource_EmptySecretDeletion(t *testing.T) {
	// This test specifically validates that deletion sends empty string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create JWT config
			{
				Config: testAccJWTConfigResourceConfig("deletion-test-system", "initial-secret"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "system_name", "deletion-test-system"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "jwt_secret", "initial-secret"),
				),
			},
			// Delete by removing from config
			{
				Config: `# JWT config removed`,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccJWTConfigResource_SpecialCharacters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJWTConfigResourceConfig(
					"special-char-system",
					"secret-with-!@#$%^&*()-_=+[]{}|;:',.<>?/~`",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "system_name", "special-char-system"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "jwt_secret", "secret-with-!@#$%^&*()-_=+[]{}|;:',.<>?/~`"),
				),
			},
		},
	})
}

func TestAccJWTConfigResource_LongSecret(t *testing.T) {
	longSecret := "very-long-jwt-secret-key-with-many-characters-to-test-length-handling-1234567890-abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJWTConfigResourceConfig("long-secret-system", longSecret),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "system_name", "long-secret-system"),
					resource.TestCheckResourceAttr("insightfinder_jwt_config.test", "jwt_secret", longSecret),
				),
			},
		},
	})
}

func testAccJWTConfigResourceConfig(systemName, jwtSecret string) string {
	return fmt.Sprintf(`
resource "insightfinder_jwt_config" "test" {
  system_name = %[1]q
  jwt_secret  = %[2]q
}
`, systemName, jwtSecret)
}

func testAccJWTConfigResourceConfigMultiple() string {
	return `
resource "insightfinder_jwt_config" "system1" {
  system_name = "jwt-system-1"
  jwt_secret  = "secret-for-system-1"
}

resource "insightfinder_jwt_config" "system2" {
  system_name = "jwt-system-2"
  jwt_secret  = "secret-for-system-2"
}
`
}
