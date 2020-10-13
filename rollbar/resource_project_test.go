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

func (s *Suite) configResourceRollbarProject() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "foo" {
		  name         = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

// checkRollbarProjectExists tests that the newly created project exists
func (s *Suite) checkRollbarProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		if err != nil {
			return err
		}
		c := s.provider.Meta().(*client.RollbarApiClient)
		proj, err := c.ReadProject(id)
		if err != nil {
			return err
		}
		if proj.Name != name {
			return fmt.Errorf("project name from API does not match project name in Terraform config")
		}
		return nil
	}
}

// checkRollbarProjectInProjectList tests that the newly created project is
// present in the list of all projects.
func (s *Suite) checkRollbarProjectInProjectList(rn string) resource.TestCheckFunc {
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
