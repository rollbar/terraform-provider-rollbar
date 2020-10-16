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
)

// TestAccRollbarProject tests creation and deletion of a Rollbar project.
func (s *AccSuite) TestAccRollbarProject() {
	rn := "rollbar_project.foo"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configResourceRollbarProject(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", s.projectName),
					s.checkRollbarProjectExists(rn, s.projectName),
					s.checkRollbarProjectInProjectList(rn),
				),
			},
		},
	})
}

func (s *AccSuite) configResourceRollbarProject() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "foo" {
		  name         = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

// checkRollbarProjectExists tests that the newly created project exists
func (s *AccSuite) checkRollbarProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(*client.RollbarApiClient)
		proj, err := c.ReadProject(id)
		s.Nil(err)
		s.Equal(name, proj.Name, "project name from API does not match project name in Terraform config")
		return nil
	}
}

// checkRollbarProjectInProjectList tests that the newly created project is
// present in the list of all projects.
func (s *AccSuite) checkRollbarProjectInProjectList(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(*client.RollbarApiClient)
		projList, err := c.ListProjects()
		s.Nil(err)
		found := false
		for _, proj := range projList {
			if proj.Id == id {
				found = true
			}
		}
		s.True(found, "Project not found in project list")
		return nil
	}
}
