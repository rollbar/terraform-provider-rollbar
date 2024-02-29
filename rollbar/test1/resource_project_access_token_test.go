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

package test1

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func (s *AccSuite) TestAccTokenImportInvalidID() {
	rn := "rollbar_project_access_token.test" // Resource name
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
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ExpectError:       regexp.MustCompile("unexpected format of ID"),
				ResourceName:      rn,
				ImportState:       true,
				ImportStateId:     "wrong format",
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTokenImport tests importing a Rollbar project access token.
func (s *AccSuite) TestAccTokenImport() {
	rn := "rollbar_project_access_token.test" // Resource name
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
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: importIdProjectAccessToken(rn),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTokenUpdateScope tests updating the scope on a Rollbar project access
// token.
func (s *AccSuite) TestAccTokenUpdateScope() {
	rn := "rollbar_project_access_token.test" // Resource name
	// language=hcl
	tmpl1 := `
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
	config1 := fmt.Sprintf(tmpl1, s.randName)
	// language=hcl
	tmpl2 := `
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
	config2 := fmt.Sprintf(tmpl2, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "scopes.#", `1`),
					resource.TestCheckResourceAttr(rn, "scopes.0", "read"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "scopes.#", `1`),
					resource.TestCheckResourceAttr(rn, "scopes.0", "post_server_item"),
					s.checkProjectAccessToken(rn),
				),
			},
		},
	})
}

// TestAccTokenUpdateRateLimit tests updating the rate limit on a Rollbar
// project access token.
func (s *AccSuite) TestAccTokenUpdateRateLimit() {
	s.T().Skip("Upgrade account to configure rate limits")
	rn := "rollbar_project_access_token.test" // Resource name
	// language=hcl
	tmpl1 := `
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
	config1 := fmt.Sprintf(tmpl1, s.randName)
	// language=hcl
	tmpl2 := `
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
	config2 := fmt.Sprintf(tmpl2, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "rate_limit_window_size", "0"),
					resource.TestCheckResourceAttr(rn, "rate_limit_window_count", "0"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					// Confirm the update produced the expected values
					resource.TestCheckResourceAttr(rn, "rate_limit_window_size", "60"),
					resource.TestCheckResourceAttr(rn, "rate_limit_window_count", "500"),
					s.checkProjectAccessToken(rn),
				),
			},
		},
	})
}

// TestAccTokenCreate tests creating a project access token.
func (s *AccSuite) TestAccTokenCreate() {
	projectResourceName := "rollbar_project.test"
	tokenResourceName := "rollbar_project_access_token.test"
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
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(tokenResourceName),
					resource.TestCheckResourceAttrSet(tokenResourceName, "access_token"),
					s.checkProjectAccessToken(tokenResourceName),
					s.checkProjectAccessTokenInTokenList(tokenResourceName),
					s.checkNoUnexpectedTokens(projectResourceName, []string{"test-token"}),
					resource.TestCheckResourceAttr(tokenResourceName, "rate_limit_window_size", "0"),
					resource.TestCheckResourceAttr(tokenResourceName, "rate_limit_window_count", "0"),
					resource.TestCheckResourceAttr(tokenResourceName, "scopes.#", `1`),
					resource.TestCheckResourceAttr(tokenResourceName, "scopes.0", "read"),
				),
			},
		},
	})
}

// TestAccTokenDelete tests deleting a project access token.
func (s *AccSuite) TestAccTokenDelete() {
	projectResourceName := "rollbar_project.test"
	// language=hcl
	tmpl1 := `
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
	config1 := fmt.Sprintf(tmpl1, s.randName)
	// language=hcl
	tmpl2 := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
	`
	config2 := fmt.Sprintf(tmpl2, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					// Project should have zero tokens
					s.checkNoUnexpectedTokens(projectResourceName, []string{}),
				),
			},
		},
	})
}

// TestAccTokenInvalidScope tests creating a project access token with an
// invalid scope.
func (s *AccSuite) TestAccTokenInvalidScope() {
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
	config := fmt.Sprintf(tmpl, s.randName)

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("invalid scope"),
			},
		},
	})
}

// TestAccTokenCreateWithNonExistentProjectID tests creating a project access
// token with a non-existent project ID.
func (s *AccSuite) TestAccTokenCreateWithNonExistentProjectID() {
	// language=hcl
	config := `
		resource "rollbar_project_access_token" "test" {
			project_id = 1234567890123457890
			name = "test-token"
			scopes = ["read"]
			status = "enabled"
			rate_limit_window_size = 60
			rate_limit_window_count = 500
		}
	`
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("not found"),
				Config:      config,
			},
		},
	})
}

// TestAccTokenDeleteOnAPIBeforeApply tests creating a Rollbar project access
// token with Terraform; then deleting it via API, before re-applying Terraform
// configuration.
func (s *AccSuite) TestAccTokenDeleteOnAPIBeforeApply() {
	projectResourceName := "rollbar_project.test"
	tokenResourceName := "rollbar_project_access_token.test"
	// FIXME: Why does adding this suffix to s.randName make this test pass,
	//  while using bare s.randName causes failure?  Could it be a Testify
	//  issue, where s.randName is not actually unique for each test in the
	//  suite?
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/160
	projectName := s.randName + "-0"
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
	config := fmt.Sprintf(tmpl, projectName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Initial create
			{
				Config: config,
			},
			// Before running Terraform, delete the token on Rollbar but not in local state
			{
				PreConfig: func() {
					c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
					var projectID int
					projects, err := c.ListProjects()
					s.Nil(err)
					for _, p := range projects {
						if p.Name == projectName {
							projectID = p.ID
						}
					}
					s.NotZero(projectID)
					tokens, err := c.ListProjectAccessTokens(projectID)
					s.Nil(err)
					for _, t := range tokens {
						if t.Name == "test-token" {
							err = c.DeleteProjectAccessToken(projectID, t.AccessToken)
							s.Nil(err)
							log.Info().
								Str("token", t.AccessToken).
								Msg("Deleted token from API before re-applying Terraform config")
						}
					}
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(projectResourceName),
					s.checkResourceStateSanity(tokenResourceName),
					s.checkProjectAccessToken(tokenResourceName),
					s.checkProjectAccessTokenInTokenList(tokenResourceName),
				),
			},
		},
	})
}

// checkProjectAccessToken checks that a PAT exists and has the expected
// properties.
func (s *AccSuite) checkProjectAccessToken(resourceName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		accessToken, err := s.getResourceIDString(ts, resourceName)
		if err != nil {
			return err
		}
		rs, ok := ts.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectIDString := rs.Primary.Attributes["project_id"]
		projectID, err := strconv.Atoi(projectIDString)
		if err != nil {
			return err
		}
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
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
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
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
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		accessToken := rs.Primary.ID

		return fmt.Sprintf("%s/%s", projectID, accessToken), nil
	}
}

// checkNoUnexpectedTokens checks that a project does not have any unexpected
// access tokens.
func (s *AccSuite) checkNoUnexpectedTokens(projectResourceName string, expectedTokenNames []string) resource.TestCheckFunc {
	l := log.With().
		Strs("expected_token_names", expectedTokenNames).
		Logger()
	expected := make(map[string]bool)
	for _, name := range expectedTokenNames {
		expected[name] = true
	}
	return func(ts *terraform.State) error {
		projectID, err := s.getResourceIDInt(ts, projectResourceName)
		if err != nil {
			return err
		}
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
		tokens, err := c.ListProjectAccessTokens(projectID)
		s.Nil(err)
		for _, t := range tokens {
			if !expected[t.Name] {
				msg := fmt.Sprintf("unexpected token name: %s", t.Name)
				l.Error().Str("unexpected_name", t.Name).Msg(msg)
				err = fmt.Errorf(msg)
				return err
			}
		}
		return nil
	}
}
