/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"strconv"
)

// TestAccRollbarProject tests creation and deletion of a Rollbar project.
func (s *AccSuite) TestAccRollbarProjectAccessToken() {
	rn := "rollbar_project_access_token.test"

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configResourceRollbarProjectAccessToken(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "access_token"),
					s.checkRollbarProjectAccessTokenExists(rn),
					//s.checkRollbarProjectAccessTokenInTokenList(rn),
				),
			},
		},
	})
}

func (s *AccSuite) configResourceRollbarProjectAccessToken() string {
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

// checkRollbarProjectAccessTokenExists tests that the newly created project exists
func (s *AccSuite) checkRollbarProjectAccessTokenExists(resourceName string) resource.TestCheckFunc {
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
		name := rs.Primary.Attributes["name"]
		scopes := rs.Primary.Attributes["scopes"]
		c := s.provider.Meta().(*client.RollbarApiClient)
		pat, err := c.ReadProjectAccessToken(projectID, accessToken)
		if err != nil {
			return err
		}
		if pat.AccessToken != accessToken {
			return fmt.Errorf("access token from API does not match access token in Terraform config")
		}
		if pat.Name != name {
			return fmt.Errorf("token name from API does not match token name in Terraform config")
		}
		if !assert.ObjectsAreEqual(pat.Scopes, scopes) {
			return fmt.Errorf("token scopes from API do not match token scopes in Terraform config")

		}
		return nil
	}
}

// checkRollbarProjectAccessTokenInProjectList tests that the newly created project is
// present in the list of all projects.
func (s *AccSuite) checkRollbarProjectAccessTokenInTokenList(rn string) resource.TestCheckFunc {
	log.Fatal().Msg("Not yet implemented")
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		if err != nil {
			return err
		}
		c := s.provider.Meta().(*client.RollbarApiClient)
		projList, err := c.ListProjects()
		if err != nil {
			return err
		}
		found := false
		for _, proj := range projList {
			if proj.Id == id {
				found = true
			}
		}
		if !found {
			msg := "Project not found in project list"
			log.Debug().Int("id", id).Msg(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}
