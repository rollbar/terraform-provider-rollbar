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
	"strconv"
)

// TestAccRollbarProject tests creation and deletion of a Rollbar project.
func (s *Suite) TestAccRollbarProject() {
	rn := "rollbar_project.foo"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.testAccRollbarProjectConfig(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", s.projectName),
					s.testAccRollbarProjectExists(rn, s.projectName),
					s.testAccRollbarProjectInProjectList(rn),
				),
			},
		},
	})
}

func (s *Suite) testAccRollbarProjectConfig() string {
	return fmt.Sprintf(`
		resource "rollbar_project" "foo" {
		  name         = "%s"
		}
	`, s.projectName)
}

// testAccRollbarProjectExists tests that the newly created project exists
func (s *Suite) testAccRollbarProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		// Check terraform config is sane
		rs, ok := ts.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not Found: %s", rn)
		}
		idString := rs.Primary.ID
		if idString == "" {
			return fmt.Errorf("No project ID is set")
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return err
		}

		// Check that project exists
		c := s.provider.Meta().(*client.RollbarApiClient)
		proj, err := c.ReadProject(id)
		if err != nil {
			return err
		}
		if proj.Name != name {
			return fmt.Errorf("project name from API does not match project name in Terraform config")
		}

		// Success
		return nil
	}
}

// testAccRollbarProjectInProjectList tests that the newly created project is
// present in the list of all projects.
func (s *Suite) testAccRollbarProjectInProjectList(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		// Check terraform config is sane
		rs, ok := ts.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not Found: %s", rn)
		}
		idString := rs.Primary.ID
		if idString == "" {
			return fmt.Errorf("No project ID is set")
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return err
		}

		// Check that project exists
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

		// Success
		return nil
	}
}
