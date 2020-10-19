/*
 * Copyright (c) 2020 Rollbar, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strconv"
)

// TestAccProjectAccessToken tests creation and deletion of a Rollbar project.
func (s *AccSuite) TestAccProjectAccessToken() {
	rn := "rollbar_project_access_token.test" // Resource name

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Info().Msg("Test create project access token with non-existent project ID")
				},
				ExpectError: regexp.MustCompile("Not found"),
				Config:      s.configResourceProjectAccessTokenNonExistentProject(),
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test invalid project access token scopes")
				},
				Config:      s.configResourceProjectAccessTokenInvalidScopes(),
				ExpectError: regexp.MustCompile("invalid scope"),
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test invalid project access token status")
				},
				Config:      s.configResourceProjectAccessTokenInvalidStatus(),
				ExpectError: regexp.MustCompile("invalid status"),
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test creating project access token")
				},
				Config: s.configResourceProjectAccessToken(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "access_token"),
					s.checkProjectAccessToken(rn),
					s.checkProjectAccessTokenInTokenList(rn),
					resource.TestCheckResourceAttr(rn, "rate_limit_window_size", "0"),
					resource.TestCheckResourceAttr(rn, "rate_limit_window_count", "0"),
					resource.TestCheckResourceAttr(rn, "scopes.#", `1`),
					resource.TestCheckResourceAttr(rn, "scopes.0", "read"),
				),
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test updating project access token rate limit")
				},
				Config: s.configResourceProjectAccessTokenUpdatedRateLimit(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					// Confirm the update produced the expected values
					resource.TestCheckResourceAttr(rn, "rate_limit_window_size", "60"),
					resource.TestCheckResourceAttr(rn, "rate_limit_window_count", "500"),
					s.checkProjectAccessToken(rn),
				),
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test updating project access token scopes")
				},
				Config: s.configResourceProjectAccessTokenUpdatedScopes(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "scopes.#", `1`),
					resource.TestCheckResourceAttr(rn, "scopes.0", "post_server_item"),
					s.checkProjectAccessToken(rn),
				),
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test importing a project access token")
				},
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: importIdProjectAccessToken(rn),
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					log.Info().Msg("Test invalid ID format when importing a project access token")
				},
				ExpectError:       regexp.MustCompile("Unexpected format of ID"),
				ResourceName:      rn,
				ImportState:       true,
				ImportStateId:     "wrong format",
				ImportStateVerify: true,
			},
		},
	})
}

func (s *AccSuite) configResourceProjectAccessToken() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test-token"
			scopes = ["read"]
			status = "enabled"
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

func (s *AccSuite) configResourceProjectAccessTokenUpdatedRateLimit() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test-token"
			scopes = ["read"]
			status = "enabled"
			rate_limit_window_size = 60
			rate_limit_window_count = 500
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

func (s *AccSuite) configResourceProjectAccessTokenUpdatedScopes() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test-token"
			scopes = ["post_server_item"]
			status = "enabled"
			rate_limit_window_size = 60
			rate_limit_window_count = 500
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

func (s *AccSuite) configResourceProjectAccessTokenNonExistentProject() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project_access_token" "test" {
			project_id = 1234567890123457890
			name = "test-token"
			scopes = ["read"]
			status = "enabled"
			rate_limit_window_size = 60
			rate_limit_window_count = 500
		}
	`
	return fmt.Sprintf(tmpl)
}

func (s *AccSuite) configResourceProjectAccessTokenInvalidScopes() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test-token"
			scopes = ["avocado"]
			status = "enabled"
			rate_limit_window_size = 60
			rate_limit_window_count = 500
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

func (s *AccSuite) configResourceProjectAccessTokenInvalidStatus() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test-token"
			scopes = ["post_server_item"]
			status = "avocado"
			rate_limit_window_size = 60
			rate_limit_window_count = 500
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

// checkProjectAccessToken tests that the newly created project exists.
func (s *AccSuite) checkProjectAccessToken(resourceName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		accessToken, err := s.getResourceIDString(ts, resourceName)
		if err != nil {
			return err
		}
		rs, ok := ts.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		projectIDString := rs.Primary.Attributes["project_id"]
		projectID, err := strconv.Atoi(projectIDString)
		if err != nil {
			return err
		}
		c := s.provider.Meta().(*client.RollbarApiClient)
		pat, err := c.ReadProjectAccessToken(projectID, accessToken)
		if err != nil {
			return err
		}
		if pat.AccessToken != accessToken {
			return fmt.Errorf("access token from API does not match access token in Terraform config")
		}
		name := rs.Primary.Attributes["name"]
		if pat.Name != name {
			return fmt.Errorf("token name from API does not match token name in Terraform config")
		}
		scopesCount, err := strconv.Atoi(rs.Primary.Attributes["scopes.#"])
		if err != nil {
			return err
		}
		var scopes []client.Scope
		for i := 0; i < scopesCount; i++ {
			attr := "scopes." + strconv.Itoa(i)
			scopeString := rs.Primary.Attributes[attr]
			s := client.Scope(scopeString)
			scopes = append(scopes, s)
		}
		if !assert.ObjectsAreEqual(pat.Scopes, scopes) {
			return fmt.Errorf("token scopes from API do not match token scopes in Terraform config")

		}
		sizeStr := rs.Primary.Attributes["rate_limit_window_size"]
		size, err := strconv.Atoi(sizeStr)
		s.Nil(err)
		if pat.RateLimitWindowSize != size {
			return fmt.Errorf("token rate_limit_window_size from API does not match token rate_limit_window_size in Terraform config")
		}
		countStr := rs.Primary.Attributes["rate_limit_window_count"]
		count, err := strconv.Atoi(countStr)
		s.Nil(err)
		if pat.RateLimitWindowCount != count {
			return fmt.Errorf("token rate_limit_window_count from API does not match token rate_limit_window_count in Terraform config")
		}
		return nil
	}
}

// checkProjectAccessTokenInTokenList tests that the newly created Rollbar
// project access token is present in the list of all project access tokens.
func (s *AccSuite) checkProjectAccessTokenInTokenList(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		accessToken, err := s.getResourceIDString(ts, rn)
		s.Nil(err)
		projectID, err := s.getResourceAttrInt(ts, rn, "project_id")
		s.Nil(err)
		c := s.provider.Meta().(*client.RollbarApiClient)
		pats, err := c.ListProjectAccessTokens(projectID)
		s.Nil(err)
		found := false
		for _, t := range pats {
			if t.AccessToken == accessToken {
				found = true
			}
		}
		s.True(found, "project access token not found in project access token list")
		return nil
	}
}

func importIdProjectAccessToken(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		projectId := rs.Primary.Attributes["project_id"]
		accessToken := rs.Primary.ID

		return fmt.Sprintf("%s/%s", projectId, accessToken), nil
	}
}
