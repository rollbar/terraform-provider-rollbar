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
	"github.com/stretchr/testify/assert"
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
				Config: s.configResourceProjectAccessToken(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "access_token"),
					s.checkProjectAccessToken(rn),
					s.checkProjectAccessTokenInTokenList(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: importIdProjectAccessToken(rn),
				ImportStateVerify: true,
			},
			{
				Config: s.configResourceProjectAccessTokenUpdatedRateLimit(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "access_token"),
					s.checkProjectAccessToken(rn),
				),
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
			rate_limit_window_size = 500
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
			return fmt.Errorf("token scopesCount from API do not match token scopesCount in Terraform config")

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
